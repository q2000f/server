package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/xormplus/core"
	"github.com/xormplus/xorm"
	"hope.lib/sql/shard"
	"log"
	"net/http"
	"testing"
	_"net/http/pprof"
)

//Player
type Player struct {
	ID    int `gorm:"primary" xorm:"pk"`
	Name  string `gorm:"char(100) not null" xorm:"char(100) notnull"`
	Level int8   `gorm:"tinyint not null default 1" xorm:"tinyint notnull default 1"`
}
var  engine *xorm.Engine
var db *gorm.DB
func init() {
	engine = InitXorm()
	db = InitGorm()
}

func InitGorm() *gorm.DB {
	db, err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/game?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		panic("connect shard failed")
	}
	//db.LogMode(true)
	//defer db.Close()
	return db
}

func InitXorm() *xorm.Engine {
	engine, err := xorm.NewEngine(xorm.MYSQL_DRIVER,
		"root:123456@tcp(127.0.0.1:3306)/game?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
		return nil
	}
	engine.DB().SetMaxOpenConns(10)
	engine.DB().SetMaxIdleConns(10)
	engine.SetMapper(core.SameMapper{})
	engine.SetTableMapper(core.SameMapper{})
	engine.SetColumnMapper(core.SameMapper{})
	//engine.ShowSQL(true)
	return engine
}

func InitSqlDB() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/game?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return db
}
var count = 2

func TestMysql(t *testing.T) {
	db := InitSqlDB()
	db.Exec("delete from players")
	player := Player{Name: "qq", Level : 1}
	for i := 1; i < count; i++ {
		player.ID = i
		stmt, _ := db.Prepare("insert into players values(?, ?, ?)")
		_, err := stmt.Exec(player.ID, player.Name, player.Level)
		if err != nil {
			log.Println(err)
		}
	}
}

func TestReflect(t *testing.T) {
	db := InitSqlDB()
	mdb := ModeDB{shard: db}
	db.Exec("delete from players")
	player := Player{Name: "qq", Level : 1}
	for i := 1; i < count; i++ {
		player.ID = i
		err := mdb.Add(&player, "players")
		if err != nil {
			log.Println(err)
		}
	}
}

func TestXormInsert(t *testing.T) {
	engine.DropTables(Player{})
	err := engine.Sync2(&Player{})
	err = engine.Sync2(&shard.Item{})
	if err != nil {
		t.Fatal(err)
		return
	}

	player := Player{Name: "qq", Level : 1}
	for i := 1; i < count; i++ {
		player.ID = i
		engine.Insert(player)
	}
	engine.Close()
}

func TestGormInsert(t *testing.T) {
	db := InitGorm()
	//db.LogMode(true)
	db.DropTable(Player{})
	db.AutoMigrate(Player{})
	db.DB().SetMaxOpenConns(10)
	db.DB().SetMaxIdleConns(10)

	defer db.Close()

	player := Player{ Name: "qq", Level : 1}
	for i := 1; i < count; i++ {
		player.ID = i
		db.Create(&player)
	}
	http.ListenAndServe(":911", nil)
}
