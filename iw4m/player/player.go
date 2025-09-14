package player

import (
	"encoding/json"
	"fmt"

	"github.com/Yallamaztar/iw4m-go/iw4m"
)

type Player struct {
	iw4m *iw4m.IW4MWrapper
}

func NewPlayer(iw4m *iw4m.IW4MWrapper) *Player {
	return &Player{iw4m: iw4m}
}

func (p *Player) Stats(clientID string) (Stats, error) {
	url := fmt.Sprintf("/api/stats/%s", clientID)
	res, err := p.iw4m.DoRequest(url)
	if err != nil {
		return Stats{}, err
	}

	body, err := readBody(res)
	if err != nil {
		return Stats{}, err
	}

	var statsSlice []Stats
	if err := json.Unmarshal(body, &statsSlice); err != nil {
		return Stats{}, err
	}

	return statsSlice[0], nil
}
