package blist

import "net/http"
import "strings"
import "regexp"
import "q29"
import "q29/user"
import "blacklistme/model/emaddr"

const regexEmailValue string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"

var rxEmail = regexp.MustCompile(regexEmailValue)

func BeforeFilter(q *q29.ReqRsp) bool {
	return true
}

func Index(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
	}
	page.Vw.Template = "blist/index"
	q29.Render(q, &page)
}

func Inquire(q *q29.ReqRsp) {
	var page struct {
		Email       string
		Found       bool		
		EmAddr      emaddr.Emaddr
		Vw          q29.View
	}
	page.Email = q.R.URL.Query().Get("email")
	if rxEmail.MatchString(page.Email) == false {
		http.Error(q.W, q.R.URL.Path+" invalid email address", 404)
		return
	}
	page.Found = emaddr.Find(q.M, page.Email, &page.EmAddr)
	q29.Render(q, &page)	
}

func AddRem(q *q29.ReqRsp) {
	var email   string
	var found   bool		
	var emAddr  emaddr.Emaddr
	
	q.R.ParseForm()
	email = q.R.FormValue("email")
	if rxEmail.MatchString(email) == false {
		http.Error(q.W, q.R.URL.Path+" invalid email address", 404)
		return
	}
	found = emaddr.Find(q.M, email, &emAddr)
	if found == true {
		emaddr.Delete(q.M, emAddr.Id);
		q29.SetFlash(q, "The email address "+email+" has been removed from the global blacklistme database.")
		q29.Redirect(q, "blist/complete/remove")		
	} else {
		emAddr.Email = email;
		emaddr.Upsert(q.M, &emAddr);
		q29.SetFlash(q, "The email address "+email+" has been added to the global blacklistme database.")
		q29.Redirect(q, "blist/complete")			
	}
}

func Complete(q *q29.ReqRsp) {
	var page struct {
		Remove bool
		Vw q29.View
	}
	page.Remove = false
	chunks := strings.Split(q.R.URL.Path, "/")
	if chunks[len(chunks)-1] == "remove" {
		page.Remove = true
	}
	q29.Render(q, &page)
}

func Dump(q *q29.ReqRsp) {
	var page struct {
		Em []emaddr.Emaddr
		Us []user.User
		Vw q29.View
	}
	emaddr.List(q.M, &page.Em)
	user.List(q.M, &page.Us)
	q29.Render(q, &page)
}
