package api

//import "html/template"
import "q29"
import "blacklistme/model/emaddr"
import "blacklistme/model/apikey"
import "net/http"
import "encoding/json"

type JSONE struct {
	Email     string `json:"email"`
	Blacklist string `json:"blacklist"`
	Error     string `json:"error,omitempty"`
}

func BeforeFilter(q *q29.ReqRsp) bool {
	return true
}

func Index(q *q29.ReqRsp) {	
	var emad emaddr.Emaddr
	var apik apikey.Apikey
	var jresponse JSONE
	
	email := q.R.URL.Query().Get("email")
	apiky := q.R.URL.Query().Get("apikey")	
	jresponse.Email = email
	jresponse.Blacklist = "no"
	if apiky == "" {
		jresponse.Error = "apikey argument required"
	}	else if email == "" {		
		jresponse.Error = "email argument required"		
	}	else {
		foundApiKey := apikey.Find(q.M, apiky, &apik)
		if foundApiKey == false {
			jresponse.Error = "apikey invalid"			
		} else {
			foundEmail := emaddr.FindByUid(q.M, "blacklistprivate", apik.UserId, email, &emad)
			if foundEmail == true {
				jresponse.Blacklist = "yes"
			} else {
				foundEmail = emaddr.Find(q.M, "blacklist", email, &emad)
				if foundEmail == true {
					jresponse.Blacklist = "yes"
				}
			}
		}
	}
	js, err := json.Marshal(jresponse)
	if err != nil {
		http.Error(q.W, err.Error(), http.StatusInternalServerError)
		return
	}
	q.W.Header().Set("Content-Type", "application/json")
	q.W.Write(js)
}
