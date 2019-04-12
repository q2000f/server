package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/core"
	"github.com/xormplus/xorm"
	"log"
	"github.com/json-iterator/go"
)

//Account
type Account struct {
	ID int64
	Name string `xorm:"varchar(255) not null unique"`
	Level int8 `xorm:"tinyint not null default 1"`
}

func main3() {
	engine, err := xorm.NewEngine(xorm.MYSQL_DRIVER,
		"root:123456@tcp(127.0.0.1:3306)/game?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		return
	}
	engine.SetMapper(core.SameMapper{})
	engine.SetTableMapper(core.SameMapper{})
	engine.SetColumnMapper(core.SameMapper{})
	engine.ShowSQL(true)
	engine.DropTables(Account{})

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	z := &Account{}
	json.Marshal(z)
	json.Unmarshal([]byte{}, &z)
	err = engine.Sync2(&Account{})
	if err != nil {
		log.Println(err)
		return
	}
	acc := Account{}
	acc.Name = "aa3"
	acc.Level = 2
	engine.Insert(&acc)
	log.Println("id", acc.ID)

	var acc2 = Account{}
	engine.ID(1).Get(&acc2)
	log.Println(fmt.Sprintf("get: %v", acc2))

	//acc.Level = 100
	//engine.ID(1).Update(&acc)
	//
	//acc2 := Account{ID: 1}
	//engine.ID(&acc2).Get(&acc2)
	//log.Println(fmt.Sprintf("affter update get: %v", acc2))
	//
	//engine.ID(1).Delete(&acc)
}
