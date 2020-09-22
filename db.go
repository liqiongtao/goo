package goo

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"time"
)

type gooDB struct {
	ctx  context.Context
	conf DBConfig
	orm  *xorm.EngineGroup
}

func NewDB(ctx context.Context, conf DBConfig) *gooDB {
	db := &gooDB{
		ctx:  ctx,
		conf: conf,
	}
	db.new()
	go db.ping()
	return db
}

func (db *gooDB) new() {
	conns := []string{db.conf.Master}
	if n := len(db.conf.Slaves); n > 0 {
		conns = append(conns, db.conf.Slaves...)
	}

	var err error
	db.orm, err = xorm.NewEngineGroup(db.conf.Driver, conns)
	if err != nil {
		panic(err.Error())
	}

	db.orm.SetLogger(&DBLogger{})

	db.orm.ShowSQL(db.conf.LogModel)
	db.orm.SetMaxIdleConns(db.conf.MaxIdle)
	db.orm.SetMaxOpenConns(db.conf.MaxOpen)
}

func (db *gooDB) ping() {
	dur := 5 * time.Second
	ti := time.NewTimer(dur)
	defer ti.Stop()
	for {
		select {
		case <-db.ctx.Done():
			return
		case <-ti.C:
			if err := db.orm.Ping(); err != nil {
				log.Println("[db-ping]", err.Error())
			}
			ti.Reset(dur)
		}
	}
}

var __db *gooDB

func DB() *xorm.EngineGroup {
	return __db.orm
}

func DBInit(conf DBConfig) {
	__db = NewDB(ctx, conf)
}
