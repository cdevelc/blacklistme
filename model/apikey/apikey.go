package apikey

import "log"
import "fmt"
import "time"
import "crypto/sha256"
import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"

type Apikey struct {
	Id        bson.ObjectId "_id,omitempty"
	APIkey    string
	UserId    bson.ObjectId ",omitempty"
	Created   string
	Requests  int
}

var collection string = "apikeylist"

func Upsert(m *mgo.Session, ap *Apikey) bool {
  var selector bson.M 

	if ap.Id == "" { //new apikey, first insertion
		selector = bson.M{"_id": nil}

	} else {
		selector = bson.M{"_id": ap.Id}	
	}
	ap.Created = time.Now().Format("2006-01-02 15:04:05")
	str := ap.Created + fmt.Sprintf("%x",[]byte(ap.UserId))
	ap.APIkey = fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
	_,err := m.DB("").C(collection).Upsert(selector, ap)
	if err == nil {
		//		if ap.Id == "" { ap.Id = changeInfo.UpsertedId.(bson.ObjectId) }
		// failing on new mongo for some reason
		return true
	}
	log.Printf("%s Upsert: %s\n", collection, err)
	return false
}

func Find(m *mgo.Session, apikey string, ap *Apikey) bool {
	err := m.DB("").C(collection).Find(bson.M{"apikey": apikey}).One(ap)
	if err == nil {
		ap.Requests = ap.Requests + 1
		m.DB("").C(collection).Update(bson.M{"_id":ap.Id}, ap)		
		return true		
	}
	if err.Error() != "not found" {
		log.Printf("%s Find: %s\n", collection, err)
	}
	return false
}

func FindByUserId(m *mgo.Session, userid bson.ObjectId, ap *Apikey) bool {
	err := m.DB("").C(collection).Find(bson.M{"userid": userid}).One(ap)
	if err == nil {
		return true		
	}
	if err.Error() != "not found" {
		log.Printf("%s Find: %s\n", collection, err)
	}
	return false
}

