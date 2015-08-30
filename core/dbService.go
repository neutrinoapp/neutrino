package realbase

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"log"
	"github.com/go-realbase/realbase/utils"
)

var connectionPool map[string]*mgo.Session

type DbService interface {
	GetSettings() map[string]string
	GetSession() *mgo.Session
	GetDb() *mgo.Database
	GetCollection() *mgo.Collection
	Insert(doc bson.M) error
	Update(q, u bson.M) error
	FindId(id, fields interface{}) (bson.M, error)
	Find(query, fields interface{}) ([]bson.M, error)
	RemoveId(id interface{}) error
}

type dbService struct {
	connectionString, dbName, colName string
	index mgo.Index
}

func NewDbService(dbName, colName string, index mgo.Index) DbService {
	connectionString := GetConfig().GetConnectionString()
	d := dbService{connectionString, dbName, colName, index}
	return &d
}

func NewUsersDbService() DbService {
	return NewDbService(Constants.DatabaseName(), Constants.UsersCollection(), mgo.Index{})
}

func NewTypeDbService(appId, typeName string) DbService {
	return NewDbService(Constants.DatabaseName(), appId + "." + typeName, mgo.Index{})
}

func NewApplicationsDbService(user string) DbService {
	index := mgo.Index{
		Key: []string{"$text:name"},
		Unique: true,
		DropDups: true,
		Background: true,
		Sparse: false,
	}

	return NewDbService(Constants.DatabaseName(), user + "." + Constants.ApplicationsCollection(), index)
}

func (d *dbService) GetSettings() map[string]string {
	m := make(map[string]string)
	m["ConnectionString"] = d.connectionString
	m["DbName"] = d.dbName
	m["ColName"] = d.colName

	return m
}

func (d *dbService) GetSession() *mgo.Session {
	if connectionPool == nil {
		connectionPool = make(map[string]*mgo.Session)
	}

	storedSession := connectionPool[d.connectionString]

	if storedSession == nil {
		session, err := mgo.Dial(d.connectionString)
		if err != nil {
			log.Fatal(err)
		}

		connectionPool[d.connectionString] = session
		storedSession = session

		if len(d.index.Key) > 0 {
			if err := d.GetCollection().EnsureIndex(d.index); err != nil {
				log.Fatal(err)
			}
		}
	}

	return storedSession.Copy()
}

func (d *dbService) GetDb() *mgo.Database {
	db := d.GetSession().DB(d.dbName)
	return db
}

func (d *dbService) GetCollection() *mgo.Collection {
	col := d.GetDb().C(d.colName)
	return col
}

func (d *dbService) Insert(doc bson.M) error {
	if _, ok := doc["_id"]; !ok {
		doc["_id"] = utils.GetCleanUUID()
	}

	return d.GetCollection().Insert(doc)
}

func (d *dbService) Update(q, u bson.M) error {
	return d.GetCollection().Update(q, u)
}

func (d *dbService) FindId(id, fields interface{}) (bson.M, error) {
	result := bson.M{}

	err := d.GetCollection().FindId(id).Select(fields).One(&result)

	return result, err;
}

func (d *dbService) Find(query, fields interface{}) ([]bson.M, error) {
	result := []bson.M{}
	err := d.GetCollection().Find(query).Select(fields).All(&result)

	return result, err;
}

func (d *dbService) RemoveId(id interface{}) error {
	return d.GetCollection().RemoveId(id)
}