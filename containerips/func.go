package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/comail/colog"
	"github.com/digitalocean/godo"
	"github.com/nats-io/go-nats"
	"golang.org/x/crypto/ssh"
	"golang.org/x/oauth2"
	mgo "gopkg.in/mgo.v2"
)

var ipdb *mgo.Collection

var DOToken string

type natsWriter struct {
	topic, host string
	conn        *nats.Conn
}

func (nw *natsWriter) Write(data []byte) (int, error) {
	err := nw.conn.Publish(nw.topic, append([]byte(nw.host), data...))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to write:", err)
		nw.conn.Publish(nw.topic, []byte(err.Error()))
	}
	time.Sleep(time.Millisecond * 50)
	return len(data), err
}
func conReconnect(con *nats.Conn) {
	var err error
	con, err = nats.Connect(natsconnString)
	if err != nil {
		fmt.Fprintln(os.Stdout, "lost connection to nats.")
		os.Exit(32)
	}
}

const natsconnString = "nats://linode.jhrb.us:4222"

func main() {
	conn, err := nats.Connect(natsconnString)
	conn.SetDisconnectHandler(conReconnect)
	colog.Register()
	colog.SetFlags(log.Lshortfile)
	if err == nil {
		host, _ := os.Hostname()
		colog.SetOutput(&natsWriter{"logs", host, conn})
	} else {
		log.Fatalln(err)
	}
	if len(os.Args) != 2 {
		log.Fatalln("server or client")
	}

	sess, err := mgo.DialWithTimeout("dev.jhrb.us:27217", time.Second*2)
	ipdb = sess.DB("ips").C("ips_farming")
	if err != nil {
		log.Println("err:", err)
		os.Exit(3)
	}
	ipdb.EnsureIndex(mgo.Index{
		Key:      []string{"ip"},
		DropDups: true,
		Unique:   true,
	})
	ipdb.DropCollection()

	switch strings.ToLower(os.Args[1]) {
	case "logs":
		udpconn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 42335})
		if err != nil {
			log.Fatalln(err)
		}
		go func() {
			buff := make([]byte, 2000)
			for {
				n, addr, err := udpconn.ReadFromUDP(buff)
				if err == nil {
					log.Println(addr.IP, string(buff[:n]))
				}
				udpconn.WriteToUDP([]byte("ok\n"), addr)
			}
		}()
		conlogger := colog.NewCoLog(os.Stdout, "", 0)
		conn.Subscribe("logs", func(msg *nats.Msg) {
			conlogger.Write(msg.Data)
		})
		count := 0
		log.Println(conn.Publish("logs", []byte("Hello")))
		for {
			count++
			fmt.Fprintf(os.Stdout, "watching logs on nats %d\r", count)
			time.Sleep(time.Second * 5)
		}
	case "server":
		for {
			time.Sleep(time.Second)
			getIPStoreAndPrint()
		}
	case "service":
		ip := createDroplet(DOToken)
		if ip == "" {
			log.Println("No ip got. Bad drop?")
			os.Exit(5)
		}
		startService(ip)

	case "client":
		getIPStoreAndPrint()
	default:
		log.Fatalln("server or client")
	}

}

func getIPStoreAndPrint() {
	ip := getIP()
	if ip != "" {
		storeIP(ip)

	}
	printIPs()
}
func printIPs() {
	var ips []string
	ipdb.Find(nil).Distinct("ip", &ips)
	log.Println(getIP())
	log.Println(len(ips))
}

func storeIP(ip string) {
	err := ipdb.Insert(struct{ Ip string }{ip})
	if err != nil && !strings.HasPrefix(err.Error(), "E11000") {
		log.Println("err:", err)
	}
}

func getIP() string {
	resp, err := http.Get("http://canihazip.com/s")
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return ""
	}
	var buff bytes.Buffer
	io.Copy(&buff, resp.Body)
	return buff.String()
}

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func createDroplet(pat string) string { //pat is personal Authentiation token for DO.
	tokenSource := &TokenSource{
		AccessToken: pat,
	}
	ctx := context.Background()
	oauthClient := oauth2.NewClient(ctx, tokenSource)
	client := godo.NewClient(oauthClient)
	drops, resp, err := client.Droplets.List(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}
	allImages := []godo.Image{}
	opts := &godo.ListOptions{PerPage: 200, Page: 0}
	images, resp, _ := client.Images.List(ctx, opts)
	allImages = images
	for resp.Links != nil && !resp.Links.IsLastPage() {
		log.Printf("next page %##v", resp.Links.Pages)
		images, resp, err = client.Images.List(ctx, nil)
		if err != nil {
			log.Println(err)
			break
		}
		allImages = append(allImages, images...)
		opts.Page++
	}
	var createImage godo.DropletCreateImage
	var latest time.Time
	for _, image := range images {
		if strings.Contains(strings.ToLower(image.Name), "docker") {
			// log.Println(image.Distribution, image.ID, image.Name, image.Type, image.Slug, image.Public, image.Created)
			t, err := time.Parse(time.RFC3339, image.Created)
			if err != nil {
				log.Fatalln(err)
			}
			if t.After(latest) {
				latest = t
				createImage.ID = image.ID
				createImage.Slug = image.Slug
			}
		}
	}
	// log.Println(createImage)
	sizes, _, _ := client.Sizes.List(ctx, nil)
	var sizeUse godo.Size
	for _, size := range sizes {
		if sizeUse.PriceHourly > size.PriceHourly || sizeUse.PriceHourly == 0 {
			sizeUse = size
		}
	}
	// log.Println(sizeUse)

	opts.Page = 0
	sshkeys, resp, err := client.Keys.List(ctx, opts)
	// log.Println(sshkeys, resp, err)
	if err != nil {
		log.Fatalln(err)
	}
	keysToUse := []godo.DropletCreateSSHKey{}
	for _, key := range sshkeys {
		keysToUse = append(keysToUse, godo.DropletCreateSSHKey{
			ID:          key.ID,
			Fingerprint: key.Fingerprint,
		})
	}
	if len(drops) > 0 {
		log.Println(drops)
		ip, err := drops[0].PublicIPv4()
		if err != nil {
			log.Fatalln(err)
		}
		return ip
	}
	name := fmt.Sprintf("%s-%s-%d", sizeUse.Slug, createImage.Slug, time.Now().Unix())
	// log.Println(name)
	drop, resp, err := client.Droplets.Create(ctx, &godo.DropletCreateRequest{
		Name:    name,
		Size:    sizeUse.Slug,
		Region:  "sfo1",
		Image:   createImage,
		SSHKeys: keysToUse,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(drop)
	for drop.Locked {
		time.Sleep(time.Second)
		drop, _, _ := client.Droplets.Get(ctx, drop.ID)
		ip, err := drop.PublicIPv4()
		if err != nil {
			log.Println(err)
			continue
		}
		return ip
	}
	return ""
}

func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

func startService(ip string) {
	publicKey, err := PublicKeyFile(os.Getenv("HOME") + `/.ssh/id_rsa`)
	if err != nil {
		log.Println(err)
		return
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			publicKey,
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			log.Println(hostname)
			return nil
		},
	}
	conn, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	defer conn.Close()

	s, _ := conn.NewSession()
	defer s.Close()
	results, err := s.CombinedOutput("docker run -d hibooboo2/leguefarmer client")
	log.Println(string(results), err)
}
