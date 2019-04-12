package shard

import "github.com/xormplus/xorm"

type User struct {
	ID    int64  `xorm:"'ID'  notnull"`
	Name  string `xorm:"'Name'  varchar(20) notnull"`
	Level int16  `xorm:"'Level'  notnull"`
}

type Item struct {
	ID    int64  `xorm:"'ID'  notnull"`
	PID   int64  `xorm:"'PID'  notnull"`
	Type  string `xorm:"'Type'  varchar(20) notnull"`
	Count int16  `xorm:"'Count'  notnull"`
}

type UserDB struct {
	db *xorm.Engine
	sess *xorm.Session
}

func (userDB *UserDB)GetByID(ID int64) *User {
	var user User

	return &user
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

}

func (userDB *UserDB)Del(user *User) {

}
