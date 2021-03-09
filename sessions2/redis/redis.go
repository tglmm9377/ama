package redis

import "github.com/gin-contrib/sessions"

type Store interface {
	sessions.Store
}

