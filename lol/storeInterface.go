package lol

type lolStorer interface {
	GetGame(gameID int64, currentPlatformID string) (Game, error)
	SaveGame(game Game, currentPlatformID string) error
	GetGames(accountID int64, currentPlatformID string) ([]Game, error)
	StorePlayer(p Player, gotMatches bool) error
	UpdatePlayer(p PlayerWithVisited) error
	GetPlayersToVisit() ([]PlayerWithVisited, error)
	Close()
}