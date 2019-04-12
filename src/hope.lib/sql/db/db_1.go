package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/core"
	"github.com/xormplus/xorm"
	"hope.lib/sql/shard"
)

type User struct {
	ID    int64  `gorm:"column:ID;primary_key"`
	Name  string `gorm:"column:Name;type:varchar(20);unique"`
	Level int16  `gorm:"column:Level;index:idx_level"`
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
	engine.NewSession()
	return engine
}

var shardDBs []*xorm.Engine

func (User) TableName() string {
	return "users"
}

func syncUserDB(engine *xorm.Engine) error {
	err := engine.Sync2(&shard.User{})
	if err != nil {
		return err
	}
	return nil
}

type UserDB struct {
	db *xorm.Engine
	sess *xorm.Session
}

func NewUserDB(shardID int) *UserDB {
	num := len(shardDBs)
	if num == 0 {
		panic("error no db")
	} else if num == 1 {
		return &UserDB{db: shardDBs[0]}
	}
	return &UserDB{db: shardDBs[shardID % num]}
}


//ToDo:
func (userDB *UserDB) GetDB(id int64) *xorm.Engine {
	return shardDBs[0]
}

func (userDB *UserDB)GetByID(ID int64) *User {
	user := &User{ID: ID}
	db := userDB.GetDB(ID)
	db.Get(user)
	return user
}

func (userDB *UserDB) Begin() {
	userDB.sess = userDB.db.NewSession()
}

func (userDB *UserDB) RollBack() {
	userDB.sess.Rollback()
}

func (userDB *UserDB) Commit() {
	userDB.sess.Commit()
}

func (userDB *UserDB)Update(user *User) {
	db := userDB.GetDB(user.ID)
	db.ID(user.ID).Update(user)
}

func (userDB *UserDB)Del(user *User) {
	db := userDB.GetDB(user.ID)
	db.Where("ID=?", user.ID).Delete(&User{})
}
