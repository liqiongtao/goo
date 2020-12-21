package goo

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"time"
)

type DBConfig struct {
	Driver      string   `yaml:"driver"`
	Master      string   `yaml:"master"`
	Slaves      []string `yaml:"slaves"`
	LogModel    bool     `yaml:"log_model"`
	MaxIdle     int      `yaml:"max_idle"`
	MaxOpen     int      `yaml:"max_open"`
	AutoPing    bool     `yaml:"auto_ping"`
	logFilePath string   `yaml:"log_file_path"`
	logFileName string   `yaml:"log_file_name"`
}

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
	db.new(conf.logFilePath, conf.logFileName)
	if conf.AutoPing {
		AsyncFunc(db.ping)
	}
	return db
}

func (db *gooDB) new(logFilePath, logFileName string) {
	conns := []string{db.conf.Master}
	if n := len(db.conf.Slaves); n > 0 {
		conns = append(conns, db.conf.Slaves...)
	}

	var err error
	db.orm, err = xorm.NewEngineGroup(db.conf.Driver, conns)
	if err != nil {
		panic(err.Error())
	}

	if logFilePath == "" {
		logFilePath = "logs/"
	}
	if logFileName == "" {
		logFileName = "sql.log"
	}
	db.orm.SetLogger(newDBLogger(logFilePath, logFileName))

	db.orm.ShowSQL(db.conf.LogModel)
	db.orm.SetMaxIdleConns(db.conf.MaxIdle)
	db.orm.SetMaxOpenConns(db.conf.MaxOpen)
}

func (db *gooDB) ping() {
	dur := 60 * time.Second
	ti := time.NewTimer(dur)
	defer ti.Stop()
	for {
		select {
		case <-db.ctx.Done():
			return
		case <-ti.C:
			if err := db.orm.Ping(); err != nil {
				Log.Error(err.Error())
			}
			ti.Reset(dur)
		}
	}
}

var __db = map[string]*gooDB{}

func DB(names ...string) *xorm.EngineGroup {
	if l := len(names); l == 0 || names[0] == "" {
		return __db["default"].orm
	}
	return __db[names[0]].orm
}

func DBInit(conf DBConfig) {
	__db["default"] = NewDB(Context, conf)
}

func DBSInit(confs map[string]DBConfig) {
	for name, conf := range confs {
		__db[name] = NewDB(Context, conf)
	}
}
