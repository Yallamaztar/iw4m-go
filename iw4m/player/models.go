package player

type Stats struct {
	Name               string  `json:"name"`
	Ranking            int     `json:"ranking"`
	Kills              int     `json:"kills"`
	Deaths             int     `json:"deaths"`
	Performance        float64 `json:"performance"`
	ScorePerMinute     float64 `json:"scorePerMinute"`
	LastPlayed         string  `json:"lastPlayed"`
	TotalSecondsPlayed int     `json:"totalSecondsPlayed"`
	ServerName         string  `json:"serverName"`
	ServerGame         string  `json:"serverGame"`
}
