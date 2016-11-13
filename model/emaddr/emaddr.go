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

func List(m *mgo.Session, collection string, em *[]Emaddr) {
	m.DB("").C(collection).Find(nil).Sort("email").All(em)
}

func ListByUid(m *mgo.Session, collection string, uid bson.ObjectId, em *[]Emaddr) {
	m.DB("").C(collection).Find(bson.M{"userid": uid}).All(em)
}
func ListByUidCount(m *mgo.Session, collection string, uid bson.ObjectId) int {
	cnt, _ := m.DB("").C(collection).Find(bson.M{"userid": uid}).Count()
	return cnt
}
func ListByDid(m *mgo.Session, collection string, did bson.ObjectId, em *[]Emaddr) {
	m.DB("").C(collection).Find(bson.M{"domainid": did}).All(em)
}

func Upsert(m *mgo.Session, collection string, em *Emaddr) bool {
  var selector bson.M 

	if em.Id == "" { //new email addr, first insertion
		selector = bson.M{"_id": nil}
		em.Created = time.Now().Format("2006-01-02 15:04:05")		
		em.Sha256 = fmt.Sprintf("%x", sha256.Sum256([]byte(em.Email)))
	} else {
		selector = bson.M{"_id": em.Id}	
	}
	changeInfo,err := m.DB("").C(collection).Upsert(selector, em)
	if err == nil {
		if em.Id == "" { em.Id = changeInfo.UpsertedId.(bson.ObjectId) }
		return true
	}
	log.Printf("%s Upsert: %s\n", collection, err)
	return false
}

func Delete(m *mgo.Session, collection string, id bson.ObjectId) bool {
	if id != "" {
		err := m.DB("").C(collection).Remove(bson.M{"_id": id})
		if err == nil { return true }
		log.Printf("%s Delete: %s\n", collection, err)
	}
	return false
}

func DeleteByEmail(m *mgo.Session, collection string, e string) bool {
	if e != "" {
		err := m.DB("").C(collection).Remove(bson.M{"email": e})
		if err == nil { return true }
		log.Printf("%s DeleteByEmail: %s\n", collection, err)
	}
	return false
}
func DeleteByDomainId(m *mgo.Session, collection string, did bson.ObjectId) bool {
	if did != "" {
		_,err := m.DB("").C(collection).RemoveAll(bson.M{"domainid": did})
		if err == nil { return true }
		log.Printf("%s DeleteAllDid: %s\n", collection, err)
	}
	return false
}

func Find(m *mgo.Session, collection string, e string, em *Emaddr) bool {
	err := m.DB("").C(collection).Find(bson.M{"email": e}).One(em)
	if err == nil { return true }
	if err.Error() != "not found" {
		log.Printf("%s Find: %s\n", collection, err)
	}
	return false
}

func FindBySig(m *mgo.Session, collection string, sig string, em *Emaddr) bool {
	err := m.DB("").C(collection).Find(bson.M{"sha256": sig}).One(em)
	if err == nil { return true }
	if err.Error() != "not found" {
		log.Printf("%s Find: %s\n", collection, err)
	}
	return false
}

func FindByUid(m *mgo.Session, collection string, uid bson.ObjectId, emaddr string, em *Emaddr) bool {
	err := m.DB("").C(collection).Find(bson.M{"userid": uid, "email": emaddr}).One(em)
	if err == nil { return true }
	if err.Error() != "not found" {
		log.Printf("%s Find: %s\n", collection, err)
	}
	return false
}
func FindByDid(m *mgo.Session, collection string, did bson.ObjectId, emaddr string, em *Emaddr) bool {
	err := m.DB("").C(collection).Find(bson.M{"domainid": did, "email": emaddr}).One(em)
	if err == nil { return true }
	if err.Error() != "not found" {
		log.Printf("%s Find: %s\n", collection, err)
	}
	return false
}
