package database

import (
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
	"log"
)

type RdbSess struct {
	session   *rethinkdb.Session
	dbName    string
	tableName string
}

func Rethink() RdbSess {
	var rdbSession RdbSess
	RethinkDB = &rdbSession
	session := getRethinkSession()
	rdbSession.session = session
	rdbSession.dbName = "sdock"
	rdbSession.tableName = "devices"
	return rdbSession
}

func getRethinkSession() *rethinkdb.Session {
	session, err := rethinkdb.Connect(rethinkdb.ConnectOpts{
		Address: "127.0.0.1", // endpoint without http
	})
	if err != nil {
		log.Fatalln(err)
	}
	return session
}

func (RethinkDB *RdbSess) createRethink() error {
	RethinkDB.session = getRethinkSession()
	RethinkDB.dbName = "sdock"
	RethinkDB.tableName = "devices"
	_, err := rethinkdb.DBCreate(RethinkDB.dbName).RunWrite(RethinkDB.session)
	if err != nil {
		return err
	}

	return nil
}

func (RethinkDB *RdbSess) CreateOrUpdate(data interface{}) error {
	_, err := rethinkdb.DB(RethinkDB.dbName).Table(RethinkDB.tableName).Get(data.(map[string]interface{})["device_id"]).Update(data).RunWrite(RethinkDB.session)
	if err != nil {
		return err
	}
	return nil
}

func (RethinkDB *RdbSess) Insert(data interface{}) error {
	_, err := rethinkdb.DB(RethinkDB.dbName).Table(RethinkDB.tableName).Insert(data).RunWrite(RethinkDB.session)
	if err != nil {
		return err
	}

	return nil
}

func (RethinkDB *RdbSess) update(id string, data interface{}) error {
	_, err := rethinkdb.DB(RethinkDB.dbName).Table(RethinkDB.tableName).Get(id).Update(data).RunWrite(RethinkDB.session)
	if err != nil {
		return err
	}
	return nil
}

func (RethinkDB *RdbSess) close() {
	err := RethinkDB.session.Close()
	if err != nil {
		return
	}
}
