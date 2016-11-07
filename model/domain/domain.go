package domain

import "log"
import "time"
import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"

type Domain struct {
	Id        bson.ObjectId "_id,omitempty"
	Domain    string
	UserId    bson.ObjectId
	Created   string
	Confirm   bool
}

var collection string = "domainlist"

func Upsert(m *mgo.Session, dm *Domain) bool {
  var selector bson.M 

	if dm.Id == "" { //new domain, first insertion
		selector = bson.M{"_id": nil}

	} else {
		selector = bson.M{"_id": dm.Id}	
	}
	dm.Created = time.Now().Format("2006-01-02 15:04:05")
	changeInfo,err := m.DB("").C(collection).Upsert(selector, dm)
	if err == nil {
		if dm.Id == "" { dm.Id = changeInfo.UpsertedId.(bson.ObjectId) }
		return true
	}
	log.Printf("%s Upsert: %s\n", collection, err)
	return false
}

func Update(m *mgo.Session, dm *Domain) {
	m.DB("").C(collection).Update(bson.M{"_id": dm.Id}, dm)
}

func Find(m *mgo.Session, domain string, dm *Domain) bool {
	err := m.DB("").C(collection).Find(bson.M{"domain": domain}).One(dm)
	if err == nil {
		return true		
	}
	if err.Error() != "not found" {
		log.Printf("%s Find: %s\n", collection, err)
	}
	return false
}

func ListByUid(m *mgo.Session, userid bson.ObjectId, dm *[]Domain) {
	m.DB("").C(collection).Find(bson.M{"userid": userid}).All(dm)
}
func ListByUidCount(m *mgo.Session, userid bson.ObjectId) int {
	cnt, _ := m.DB("").C(collection).Find(bson.M{"userid": userid}).Count()
	return cnt
}

