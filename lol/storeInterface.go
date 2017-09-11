package lol

type lolStorer interface {
	GetGame(gameID int64, currentPlatformID string) (Game, error)
	SaveGame(game Game, currentPlatformID string) error
	StorePlayer(p Player) error
	Stats()
	Close()
	LoadAllGameIDS()
	lolCache
}

type lolCache interface {
	HaveGame(gameID int64) bool
	AddGame(gameID int64)
	Player(accountID int64)
	VisitedPlayer(accountID int64)
	HaveVisitedPlayer(accountID int64) bool
	GetPlayerToVisit() int64
}
