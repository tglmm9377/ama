package memstore

import (
	"github.com/gin-contrib/sessions"
	"github.com/quasoft/memstore"
)

type Store interface {
	sessions.Store
}

func NewStore(keyPairs ...[]byte)Store{
	return &store{memstore.NewMemStore(keyPairs...)}
}


type store struct {
	*memstore.MemStore
}

func (c *store) Options(options sessions.Options){
	c.MemStore.Options = options.ToGorillaOptions()
}