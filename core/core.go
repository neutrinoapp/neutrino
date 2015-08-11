package realbase

import "gopkg.in/mgo.v2"

func Connect(host string) *mgo.Session {
	session, err := mgo.Dial(host)

	if err != nil {
		panic(err)
	}

	return session
}
