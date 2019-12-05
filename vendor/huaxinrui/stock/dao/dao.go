package dao

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"time"
	"log"
)


type ORM struct {
	Store  string
}

type backend struct {
	DB    *gorm.DB
}

func (o *ORM) InitDal(auto ...interface{}) (backend *backend) {
	Backend.DB =  o.initStore(auto...)
	//Backend.Redis = o.initRedis()
	log.Println("init dal success")
	return
}

func (o *ORM) initStore(auto ...interface{}) *gorm.DB{
	var db *gorm.DB
	var err error
	if o.Store == "" {
		return nil
	}
	db, err = gorm.Open( "mysql", o.Store)
	//db.LogMode(true)
	if err != nil {
		log.Fatalln("[-] open database error.", err)
		return nil
	}
	db.DB().SetConnMaxLifetime(60 * time.Second)
	db.DB().SetMaxOpenConns(30)

	for _, i := range auto {
		db.AutoMigrate(i)
	}
	return db
}

var Backend backend






