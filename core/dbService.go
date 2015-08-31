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
	GetDb() (*mgo.Session, *mgo.Database)
	GetCollection() (*mgo.Session, *mgo.Collection)
	Insert(doc bson.M) error
	Update(q, u bson.M) error
	FindId(id, fields interface{}) (bson.M, error)
	Find(query, fields interface{}) ([]bson.M, error)
	RemoveId(id interface{}) error
	UpdateId(id, u interface{}) error
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

		_, collection := d.GetCollection()

		if len(d.index.Key) > 0 {
			if err := collection.EnsureIndex(d.index); err != nil {
				log.Fatal(err)
			}
		}
	}

	return storedSession.Copy()
}

func (d *dbService) GetDb() (*mgo.Session, *mgo.Database) {
	session := d.GetSession()
	db := session.DB(d.dbName)
	return session, db
}

func (d *dbService) GetCollection() (*mgo.Session, *mgo.Collection) {
	session, db := d.GetDb()
	col := db.C(d.colName)
	return session, col
}

func (d *dbService) Insert(doc bson.M) error {
	if _, ok := doc["_id"]; !ok {
		doc["_id"] = utils.GetCleanUUID()
	}

	session, collection := d.GetCollection()

	defer session.Close()
	return collection.Insert(doc)
}

func (d *dbService) Update(q, u bson.M) error {
	session, collection := d.GetCollection()

	defer session.Close()
	return collection.Update(q, u)
}

func (d *dbService) FindId(id, fields interface{}) (bson.M, error) {
	result := bson.M{}

	session, collection := d.GetCollection()

	defer session.Close()
	err := collection.FindId(id).Select(fields).One(&result)

	return result, err;
}

func (d *dbService) Find(query, fields interface{}) ([]bson.M, error) {
	result := []bson.M{}
	session, collection := d.GetCollection()

	defer session.Close()
	err := collection.Find(query).Select(fields).All(&result)

	return result, err;
}

func (d *dbService) RemoveId(id interface{}) error {
	session, collection := d.GetCollection()

	defer session.Close()
	return collection.RemoveId(id)
}

func (d *dbService) UpdateId(id, u interface{}) error {
	session, collection := d.GetCollection()
	defer session.Close()

	return collection.UpdateId(id, u)
}