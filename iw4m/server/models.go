package server

type ServerStatus struct {
	ID             int            `json:"id"`
	IsOnline       bool           `json:"isOnline"`
	Name           string         `json:"name"`
	MaxPlayers     int            `json:"maxPlayers"`
	CurrentPlayers int            `json:"currentPlayers"`
	Map            mapStatus      `json:"map"`
	GameMode       string         `json:"gameMode"`
	ListenAddress  string         `json:"listenAddress"`
	ListenPort     int            `json:"listenPort"`
	Game           string         `json:"game"`
	Players        []playerStatus `json:"players"`
}

type mapStatus struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

type playerStatus struct {
	Name           string `json:"name"`
	Score          int    `json:"score"`
	Ping           int    `json:"ping"`
	State          string `json:"state"`
	ClientNumber   int    `json:"clientNumber"`
	ConnectionTime int    `json:"connectionTime"`
	Level          string `json:"level"`
}

type ServerInfo struct {
	TotalConnectedClients int                      `json:"totalConnectedClients"`
	TotalClientSlots      int                      `json:"totalClientSlots"`
	TotalTrackedClients   int                      `json:"totalTrackedClients"`
	TotalRecentClients    totalRecentClientsInfo   `json:"totalRecentClients"`
	MaxConcurrentClients  maxConcurrentClientsInfo `json:"maxConcurrentClients"`
}

type totalRecentClientsInfo struct {
	Value   int    `json:"value"`
	Time    string `json:"time"`
	StartAt string `json:"startAt"`
	EndAt   string `json:"endAt"`
}

type maxConcurrentClientsInfo struct {
	Value   int    `json:"value"`
	Time    string `json:"time"`
	StartAt string `json:"startAt"`
	EndAt   string `json:"endAt"`
}

type Report struct {
	Origin    string
	Reason    string
	Target    string
	Timestamp string
}

type Help struct {
	Sections map[string]HelpSection `json:"sections"`
}

type HelpSection struct {
	Title    string             `json:"title"`
	Commands map[string]Command `json:"commands"`
}

type Command struct {
	Alias          string `json:"alias"`
	Description    string `json:"description"`
	RequiresTarget string `json:"requires_target"`
	Syntax         string `json:"syntax"`
	MinLevel       string `json:"min_level"`
}

type ServerID struct {
	Name string `json:"server"`
	ID   string `json:"id"`
}

type Chat struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type FindPlayerResponse struct {
	TotalFoundClients int          `json:"totalFoundClients"`
	Clients           []FindPlayer `json:"clients"`
}

type FindPlayer struct {
	Name     string `json:"name"`
	XUID     string `json:"xuid"`
	ClientId int    `json:"clientId"`
}

type Players struct {
	Role     string `json:"role"`
	Name     string `json:"name"`
	ClientId string `json:"clientId"`
	URL      string `json:"url"`
}

type RecentClient struct {
	Name      string `json:"name"`
	Link      string `json:"link"`
	Country   string `json:"country,omitempty"`
	IPAddress string `json:"ip_address"`
	LastSeen  string `json:"last_seen"`
}

type AuditLog struct {
	Type       string `json:"type"`
	Origin     string `json:"origin"`
	OriginRank string `json:"origin_rank,omitempty"`
	Href       string `json:"href"`
	Target     string `json:"target"`
	Data       string `json:"data"`
	Time       string `json:"time"`
}

type Admin struct {
	Name          string `json:"name"`
	Role          string `json:"role"`
	Game          string `json:"game"`
	LastConnected string `json:"last_connected"`
}

type TopPlayer struct {
	Rank   string            `json:"rank"`
	Name   string            `json:"name"`
	Link   string            `json:"link"`
	Rating string            `json:"rating"`
	Stats  map[string]string `json:"stats"`
}
