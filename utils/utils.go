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

func (u *Utils) IsHigherRole(role, roleToCheck string) bool {
	rolePos := u.RolePosition(role)
	checkPos := u.RolePosition(roleToCheck)
	if rolePos == -1 || checkPos == -1 {
		return false
	}
	return rolePos < checkPos
}

func (u *Utils) IsLowerRole(role, roleToCheck string) bool {
	return !u.IsHigherRole(role, roleToCheck)
}

func (u *Utils) CommandPrefix() string {
	log, _ := server.NewServer(u.iw4m).RecentAuditLog()
	return string(log.Data[0])
}
