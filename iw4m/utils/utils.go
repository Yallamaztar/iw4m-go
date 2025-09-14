package utils

import (
	"strings"

	"github.com/Yallamaztar/iw4m-go/iw4m"
	"github.com/Yallamaztar/iw4m-go/iw4m/server"
)

type Utils struct {
	iw4m *iw4m.IW4MWrapper
}

// Create a new Utils wrapper
func NewUtils(iw4m *iw4m.IW4MWrapper) *Utils {
	return &Utils{iw4m: iw4m}
}

func (u *Utils) RoleExists(role string) bool {
	roles, _ := server.NewServer(u.iw4m).Roles()
	for _, r := range roles {
		if strings.EqualFold(r, role) {
			return true
		}
	}
	return false
}

func (u *Utils) RolePosition(role string) int {
	roles, _ := server.NewServer(u.iw4m).Roles()

	for i, r := range roles {
		if strings.EqualFold(r, role) {
			return i
		}
	}

	return -1
}

func (u *Utils) IsHigherRole(roleToCheck, role string) bool {
	if strings.ToLower(roleToCheck) == "creator" {
		return true
	}

	rolePos := u.RolePosition(role)
	checkPos := u.RolePosition(roleToCheck)

	if rolePos == -1 || checkPos == -1 {
		return false
	}

	return rolePos < checkPos
}

func (u *Utils) IsLowerRole(roleToCheck, role string) bool {
	return !u.IsHigherRole(role, roleToCheck)
}

func (u *Utils) IsPlayerOnline(player string) bool {
	players, _ := server.NewServer(u.iw4m).ListPlayers()

	for _, p := range players {
		if p.Name == player {
			return true
		}
	}

	return false
}
