package htmx

import (
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/auth"
)

type Component interface {
	SetUser(u *auth.User)
	User() *auth.User
	Lock()
	IsLocked() bool
}

type BaseComponent struct {
	IsOutOfBand  bool
	user         *auth.User
	requiresAuth bool
}

func (c *BaseComponent) Lock() {
	c.requiresAuth = true
}

func (c *BaseComponent) IsLocked() bool {
	return c.requiresAuth
}

func (c *BaseComponent) SetUser(u *auth.User) {
	if u != nil {
		c.user = u
	}
}

func (c *BaseComponent) User() *auth.User {
	return c.user
}

func (c *BaseComponent) IsAuthenticated() bool {
	return c.user != nil
}
