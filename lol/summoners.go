package lol

import (
	"fmt"
	"net/url"
	"strconv"
)

// /lol/summoner/v3/summoners/{summonerId}
// /lol/summoner/v3/summoners/by-account/{accountId}
// /lol/summoner/v3/summoners/by-name/{summonerName}
type summoners struct {
	c *client
}

func (s *summoners) Arg(arg string) (*Summoner, error) {
	sum := s.ByName(arg)
	if sum != nil {
		return sum, nil
	}
	id, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return nil, err
	}
	sum = s.ByAccountID(id)
	sum2 := s.ID(id)
	switch {
	case sum != nil && sum2 != nil:
		return nil, fmt.Errorf("Not sure if using accountID or summonerID: %s", arg)
	case sum != nil:
		return sum, nil
	case sum2 != nil:
		return sum2, nil
	default:
		return nil, fmt.Errorf("Cannot locate summoner: %s", arg)
	}
}

func (s *summoners) ID(summonerID int64) *Summoner {
	return nil
}
func (s *summoners) ByAccountID(accountID int64) *Summoner {
	return nil
}
func (s *summoners) ByName(summonerName string) *Summoner {
	var sum Summoner
	err := s.c.GetObjRiot(fmt.Sprintf("/lol/summoner/v3/summoners/by-name/%s", url.PathEscape(summonerName)), &sum)
	if err != nil {
		logger.Println("err: ", err)
		return nil
	}
	return &sum
}

type Summoner struct {
	ID            int64  `json:"id"`
	AccountID     int64  `json:"accountId"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

type Player struct {
	PlatformID        string `json:"platformId"`
	AccountID         int64  `json:"accountId"`
	SummonerName      string `json:"summonerName"`
	SummonerID        int64  `json:"summonerId"`
	CurrentPlatformID string `json:"currentPlatformId"`
	CurrentAccountID  int64  `json:"currentAccountId"`
	MatchHistoryURI   string `json:"matchHistoryUri"`
	ProfileIcon       int    `json:"profileIcon"`
}

func (p Player) ToSummoner() Summoner {
	return Summoner{
		ID:            p.SummonerID,
		Name:          p.SummonerName,
		AccountID:     p.AccountID,
		ProfileIconID: p.ProfileIcon,
	}
}
