package lol

import "golang.org/x/sync/syncmap"

type memCache struct {
	games   syncmap.Map
	toVisit syncmap.Map
	visit   chan int64
	visited syncmap.Map
}

func (c *memCache) HaveGame(gameID int64) bool {
	_, ok := c.games.Load(gameID)
	return ok
}

func (c *memCache) AddGame(gameID int64) {
	c.games.Store(gameID, struct{}{})
}

func (c *memCache) Player(accountID int64) {
	_, ok := c.visited.Load(accountID)
	if ok {
		return
	}
	c.toVisit.Store(accountID, struct{}{})
}

func (c *memCache) VisitedPlayer(accountID int64) {
	c.toVisit.Delete(accountID)
	c.visited.Store(accountID, struct{}{})
}

func (c *memCache) GetPlayerToVisit() int64 {
	var id int64
	var ok bool
	c.toVisit.Range(func(key interface{}, value interface{}) bool {
		id, ok = key.(int64)
		return false
	})
	if id == 0 {
		return 0
	}
	c.VisitedPlayer(id)
	return id
}

func (c *memCache) HaveVisitedPlayer(accountID int64) bool {
	_, ok := c.visited.Load(accountID)
	return ok
}
