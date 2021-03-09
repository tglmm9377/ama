package cookie

import (
	"github.com/gin-contrib/sessions"
	gsessions "github.com/gorilla/sessions"
)

type Store interface {
	sessions.Store
}


func NewStore(keyPairs ...[]byte) Store{
	return &store{gsessions.NewCookieStore(keyPairs...)}
}

type store struct{
	*gsessions.CookieStore
}

func  (c *store) Options(options sessions.Options){
	c.CookieStore.Options = options.ToGorillaOptions()
}
