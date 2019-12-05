package mysql

import (
	s "magic/stock/service/middleware/session"
	"magic/stock/service/middleware/session/storage/gormstore"
	"magic/stock/service/middleware/session/storage/sessions"

	"github.com/jinzhu/gorm"
)

type GormStore interface {
	s.Store
}

func NewGormStore(db *gorm.DB, keyPairs ...[]byte) GormStore {
	store := gormstore.New(db, keyPairs...)
	return &gormStore{store}
}

type gormStore struct {
	*gormstore.Store
}

func (c *gormStore) Options(options s.Options) {
	c.Store.SessionOpts = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
