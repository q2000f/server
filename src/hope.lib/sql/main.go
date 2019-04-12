package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	_ "github.com/xormplus/xorm"
	"log"
	"reflect"
	"strings"
	"sync"
	"util"
)

type Field struct {
	Name      string
	Type      string
	Default   string
	DBType    string
	IsPrimary bool
	IsUnique  bool
	IndexName string
}

type Index []string
type Query struct {
	Param Index
	Name  string
}

var z sync.Map

type Table struct {
	Name       string
	Fields     []Field
	PrimaryKey Index
	ShardKey   Index
	Indexes    []Index
	Uniques    []Index
	Gets       []Query
	Finds      []Query
}

func ToGormField(field Field) string {
	var tag string
	if field.DBType != "" {
		tag += ";type:" + field.DBType
	}
	if field.IsPrimary {
		tag += ";primary_key"
	}
	if field.IsUnique {
		tag += ";unique"
	}
	if field.IndexName != "" {
		tag += ";index:" + field.IndexName
	}
	if tag != "" {
		tag = fmt.Sprintf("`gorm:\"column:%s%s\"`", field.Name, tag)
	} else {
		tag = fmt.Sprintf("`gorm:\"column:%s\"`", field.Name)
	}
	return fmt.Sprintf("%s %s %s", field.Name, field.Type, tag)
}

func ToXormField(field Field) string {
	var tag string
	if field.DBType != "" {
		tag += " " + field.DBType
	}
	tag += " notnull"
	if field.IsPrimary {
		tag += " pk"
	}
	if field.IsUnique {
		tag += " unique"
	}
	if field.IndexName != "" {
		tag += " index(" + field.IndexName + ")"
	}

	tag = fmt.Sprintf("`xorm:\"'%s' %s\"`", field.Name, tag)

	return fmt.Sprintf("%s %s %s", field.Name, field.Type, tag)
}

func ToStruct(table Table) string {
	fields := []string{}
	for _, field := range table.Fields {
		fields = append(fields, ToXormField(field))
	}
	return fmt.Sprintf("type %s struct {\n\t%s\n}", table.Name, strings.Join(fields, "\n\t"))
}

func ToPackage(name string, tables []Table) string {
	res := []string{}
	for _, table := range tables {
		res = append(res, ToStruct(table))
	}
	txt := fmt.Sprintf("package %s\n%s", name, strings.Join(res, "\n\n"))
	util.SaveGoFile(name, "struct.go", []byte(txt))
	return txt
}

func TestOrm() {
	var tables []Table
	err := util.LoadJsonFile("mysql.json", &tables)
	if err != nil {
		log.Println(err)
		stackErr := errors.Wrap(err, "test err!")
		log.Println(stackErr)
		return
	}

	packName := "shard"
	txt := ToPackage(packName, tables)

	util.SaveGoFile(packName, "struct.go", []byte(txt))
}

//gorm:"primary_key;index;unique_index;type:varchar(100);AUTO_INCREMENT;size:255; not null"
type User struct {
	ID    int64  `gorm:"column:ID;primary_key"`
	Name  string `gorm:"column:Name"`
	Level int16  `gorm:"column:Level"`
}

func (User) TableName() string {
	return "users"
}

func init() {
	db, err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/game?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		panic("connect shard failed")
	}
	defer db.Close()

	db.AutoMigrate(&User{})
}

func GetCondition(conditions interface{}) (where string, params []interface{}) {
	conditionsValue := reflect.Indirect(reflect.ValueOf(conditions))
	paramNames := []string{}
	params = make([]interface{}, 0)
	for i := 0; i < conditionsValue.Type().NumField(); i++ {
		paramNames = append(paramNames, fmt.Sprint(conditionsValue.Type().Field(i).Name, "=?"))
		params = append(params, conditionsValue.Field(i).Addr().Interface())
	}
	where = strings.Join(paramNames, " and ")
	return where, params
}

type ModeDB struct {
	shard *sql.DB
}

func (db *ModeDB) Get(data interface{}, tableName, where string, params ...interface{}) error {
	dataValue := reflect.Indirect(reflect.ValueOf(data))

	fieldNames := []string{}
	fieldDatas := make([]interface{}, 0)
	for i := 0; i < dataValue.Type().NumField(); i++ {
		fieldNames = append(fieldNames, dataValue.Type().Field(i).Name)
		fieldDatas = append(fieldDatas, dataValue.Field(i).Addr().Interface())
	}

	s := fmt.Sprintf("select %s from %s", strings.Join(fieldNames, ", "), tableName)
	if where != "" {
		s += " where " + where
		log.Println(s)
	}
	stmt, err := db.shard.Prepare(s)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(params...).Scan(fieldDatas...)
	if err != nil {
		return err
	}
	return nil
}

func (db *ModeDB) Add(data interface{}, tableName string) error {
	dataValue := reflect.Indirect(reflect.ValueOf(data))
	fieldNames := []string{}
	valueDatas := make([]interface{}, 0)
	valueTokens := []string{}
	for i := 0; i < dataValue.Type().NumField(); i++ {
		fieldNames = append(fieldNames, dataValue.Type().Field(i).Name)
		valueDatas = append(valueDatas, dataValue.Field(i).Addr().Interface())
		valueTokens = append(valueTokens, "?")
	}

	sqlPrepare := fmt.Sprintf("insert into %s(%s) values(%s)", tableName,
		strings.Join(fieldNames, ","),
		strings.Join(valueTokens, ","))
	stmt, err := db.shard.Prepare(sqlPrepare)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(valueDatas...)
	if err != nil {
		return err
	}
	return nil
}

func (db *ModeDB) Set(data interface{}, old interface{}, tableName, where string, params ...interface{}) error {
	dataValue := reflect.Indirect(reflect.ValueOf(data))
	oldValue := reflect.Indirect(reflect.ValueOf(old))

	updateNames := []string{}
	valueDatas := make([]interface{}, 0)
	for i := 0; i < dataValue.Type().NumField(); i++ {
		fieldName := dataValue.Type().Field(i).Name
		v := oldValue.FieldByName(fieldName)
		if !v.IsValid() {
			log.Println("can not find:", fieldName)
			continue
		}

		if v.Type().Kind() != dataValue.Type().Field(i).Type.Kind() {
			log.Println("invalid type:", v.Type().Kind(), "!=", dataValue.Type().Field(i).Type.Kind())
			continue
		}
		if v.Addr().Interface() == dataValue.Field(i).Addr().Interface() {
			continue
		}

		log.Println("diff:", fieldName)

		v.Set(dataValue.Field(i))

		updateNames = append(updateNames, fieldName+"=?")
		valueDatas = append(valueDatas, dataValue.Field(i).Addr().Interface())
	}
	valueDatas = append(valueDatas, params...)

	sqlPrepare := fmt.Sprintf("update %s set %s ", tableName,
		strings.Join(updateNames, ","))
	if where != "" {
		sqlPrepare += " where " + where
	}
	stmt, err := db.shard.Prepare(sqlPrepare)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(valueDatas...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (db *ModeDB) Del(tableName string, where string, params ...interface{}) error {
	sql := "delete from " + tableName
	if where != "" {
		sql += " where " + where
	}
	stmt, err := db.shard.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(params...)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	TestOrm()
	return
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/game?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		return
	}
	userDB := ModeDB{shard: db}
	user := User{}
	find := struct{ ID int }{100}

	userDB.Get(&find, "users", "")
	log.Println(user)
	//id := userDB.Insert(&User{ Name: "vv", Level: 3})
	//userDB.findOne()
	//if id > 0 {
	//	userDB.Delete(id)
	//}

	var tmp interface{}
	tmp = &user
	rt := reflect.TypeOf(user)
	log.Println("user name:", rt.Name())
	log.Println("*user name:", reflect.TypeOf(&user).Name())
	log.Println("interface name:", reflect.TypeOf(tmp).Name())
	log.Println("interface name:", reflect.ValueOf(tmp))
	log.Println("indirect name:", reflect.TypeOf(reflect.Indirect(reflect.ValueOf(tmp))).Name())

	//userValue := reflect.Indirect(reflect.ValueOf(tmp))

	//vT := reflect.TypeOf(&user)
	//reflect.Indirect(vT)
	//sValue := reflect.Indirect(reflect.ValueOf(tmp))

	user2 := User{ID: 100}
	sValue := reflect.ValueOf(user2)

	st := sValue.Type()
	for i := 0; i < st.NumField(); i++ {
		fieldS := st.Field(i)
		z := sValue.Field(i).IsValid()
		log.Println(fieldS.Name, "111", z, sValue.Field(i))
	}
}
