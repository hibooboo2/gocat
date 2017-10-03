package lol

// RiotClient is the interface that defines the interactions with riots api.
type RiotClient interface {
	Mastery() *champMastery
	Spectator() *spectator
	StaticData() *staticData
	// /lol/league/v3/challengerleagues/by-queue/{queue}
	// /lol/league/v3/leagues/by-summoner/{summonerId}
	// /lol/league/v3/masterleagues/by-queue/{queue}
	// /lol/league/v3/positions/by-summoner/{summonerId}
	// /lol/match/v3/matches/{matchId}
	// /lol/match/v3/matches/{matchId}/by-tournament-code/{tournamentCode}
	// /lol/match/v3/matches/by-tournament-code/{tournamentCode}/ids
	// /lol/match/v3/matchlists/by-account/{accountId}
	// /lol/match/v3/matchlists/by-account/{accountId}/recent
	// /lol/match/v3/timelines/by-match/{matchId}
	// /lol/platform/v3/champions
	// /lol/platform/v3/champions/{id}
	// /lol/platform/v3/masteries/by-summoner/{summonerId}
	// /lol/platform/v3/masteries/by-summoner/{summonerId}
	// /lol/platform/v3/runes/by-summoner/{summonerId}
	// /lol/static-data/v3/masteries
	// /lol/static-data/v3/masteries/{id}
	// /lol/static-data/v3/profile-icons
	// /lol/static-data/v3/runes
	// /lol/static-data/v3/runes/{id}
	// /lol/static-data/v3/summoner-spells
	// /lol/static-data/v3/summoner-spells/{id}
	// /lol/static-data/v3/realms

	// /lol/status/v3/shard-data
	Summoners() *summoners
	GetAllGames(accountID int64, platformID string) ([]Game, error)
}

func (c *client) Mastery() *champMastery {
	return &champMastery{c}
}

func (c *client) Summoners() *summoners {
	return &summoners{c}
}

func (c *client) Spectator() *spectator {
	return &spectator{c}
}

func (c *client) StaticData() *staticData {
	return &staticData{c}
}
