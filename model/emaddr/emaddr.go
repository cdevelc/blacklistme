package emaddr

import "log"
import "fmt"
import "time"
import "crypto/sha256"
import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"

type Emaddr struct {
	Id        bson.ObjectId "_id,omitempty"
	Email     string
	DomainId  bson.ObjectId ",omitempty"
	UserId    bson.ObjectId ",omitempty"
	Created   string
	Sha256    string
}

func List(m *mgo.Session,em *[]Emaddr) {
	m.DB("").C("blacklist").Find(nil).Sort("email").All(em)
}

func Upsert(m *mgo.Session, em *Emaddr) bool {
  var selector bson.M 

	if em.Id == "" { //new email addr, first insertion
		selector = bson.M{"_id": nil}
		em.Created = time.Now().Format("2006-01-02 15:04:05")		
		em.Sha256 = fmt.Sprintf("%x", sha256.Sum256([]byte(em.Email)))
	} else {
		selector = bson.M{"_id": em.Id}	
	}
	changeInfo,err := m.DB("").C("blacklist").Upsert(selector, em)
	if err == nil {
		if em.Id == "" { em.Id = changeInfo.UpsertedId.(bson.ObjectId) }
		return true
	}
	log.Printf("blacklist Upsert: %s\n", err)
	return false
}

func Delete(m *mgo.Session, id bson.ObjectId) bool {
	if id != "" {
		err := m.DB("").C("blacklist").Remove(bson.M{"_id": id})
		if err == nil { return true }
		log.Printf("blacklist Delete: %s\n", err)
	}
	return false
}

func Find(m *mgo.Session, e string, em *Emaddr) bool {
	err := m.DB("").C("blacklist").Find(bson.M{"email": e}).One(em)
	if err == nil { return true }
	if err.Error() != "not found" {
		log.Printf("blacklist Find: %s\n", err)
	}
	return false
}
