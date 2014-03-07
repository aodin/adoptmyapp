package robot

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

// Wraps an existing connection
type DB interface {
	Get(string) (*Meta, error)
	Save(string) (*Meta, error)
}

type Mongo struct {
	name    string
	session *mgo.Session
}

// Get the repository from the database using the normalized URL. If no
// repository exists, return nil and an error.
func (db *Mongo) Get(url string) (meta *Meta, err error) {
	c := db.session.DB(db.name).C("repositories")
	err = c.Find(bson.M{"url": url}).One(&meta)
	return
}

func (db *Mongo) Save(url string) (*Meta, error) {
	meta := &Meta{URL: url, Added: time.Now()}
	c := db.session.DB(db.name).C("repositories")
	// TODO test for existence first?
	err := c.Insert(meta)
	return meta, err
}

func NewMongoDB(session *mgo.Session, name string) *Mongo {
	return &Mongo{name: name, session: session}
}

type Memory struct {
	repositories map[string]*Meta
}

func (db *Memory) Get(url string) (*Meta, error) {
	meta, exists := db.repositories[url]
	if exists {
		return meta, nil
	}
	return nil, fmt.Errorf("The repository %s does not exist", url)
}

func (db *Memory) Save(url string) (*Meta, error) {
	meta := &Meta{URL: url, Added: time.Now()}
	db.repositories[url] = meta
	return meta, nil
}

func NewMemoryDB() *Memory {
	return &Memory{repositories: make(map[string]*Meta)}
}
