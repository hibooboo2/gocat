package lol

type lolStorer interface {
	GetGame(gameID int64, currentPlatformID string) (Game, error)
	SaveGame(game Game, currentPlatformID string) error
	GetGames(accountID int64, currentPlatformID string) ([]Game, error)
	StorePlayer(p Player) error
	VisitPlayer(p Player) error
	GetPlayersToVisit() ([]Player, error)
	GetPlayerToVisit() (Player, error)
	CheckGameStored(gameID int64) bool
	Stats()
	Close()
}
