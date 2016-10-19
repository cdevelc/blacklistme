package api

//import "html/template"
import "q29"
import "blacklistme/model/emaddr"
import "net/http"
import "encoding/json"

type JSONE struct {
	Email     string `json:"email"`
	Blacklist string `json:"blacklist"`
}

func BeforeFilter(q *q29.ReqRsp) bool {
	return true
}

func Index(q *q29.ReqRsp) {	
	var emad emaddr.Emaddr
	var jresponse JSONE
	
	emad.Email = q.R.URL.Query().Get("email")
	jresponse.Email = emad.Email
	jresponse.Blacklist = "no"
	found := emaddr.Find(q.M, emad.Email, &emad)
	if found == true { jresponse.Blacklist = "yes" }
	js, err := json.Marshal(jresponse)
	if err != nil {
		http.Error(q.W, err.Error(), http.StatusInternalServerError)
		return
	}
	q.W.Header().Set("Content-Type", "application/json")
	q.W.Write(js)
}
