package dlist

import "log"
//import "time"
import "strings"
import "net/http"
import "q29"
import "q29/validfield"
import "q29/session"
import "blacklistme/model/domain"
import "blacklistme/model/emaddr"
import "blacklistme/util/domaintoemail"
import "blacklistme/util/mailgo"

func BeforeFilter(q *q29.ReqRsp) bool {
	if q.U == nil && q.Action == "addconfirm" { return true } //allow Domain Confirm without login
	if q.U != nil { return true }
	q29.Redirect(q, "account/login")
	return false
}

func Index(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
		Dlist []domain.Domain
		DlistCount int
		FlashMsg string
	}
	domain.ListByUid(q.M, q.U.Id, &page.Dlist)
	page.DlistCount = len(page.Dlist)
	q29.Render(q, &page)
}

func sendEmailToDomainContact(dom string, vos string, q *q29.ReqRsp) {
	//seems to take a little while, so we will call here as a goroutine and let the web user move on
	emaddress, _ := domaintoemail.Get(dom)
	if emaddress != "" {
		s := mailgo.Session {
			Email: emaddress,
			URL: "http://"+q.R.Host+q29.AssetURL(q, "dlist/addconfirm?vos="+vos),
		}
		log.Printf("emurl = %s", s.URL)
		mailgo.ConfirmDomainControlEmail(&s)
	}	
}

func AddConfirm(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
		Dm domain.Domain
	}
	
	found := domain.FindByVos(q.M, q.R.URL.Query().Get("vos"), &page.Dm)
	if found == false {
		http.Error(q.W, q.R.URL.Path+" invalid request", 404)
		return
	}
	if page.Dm.Confirm == false {
		page.Dm.Confirm = true
		domain.Update(q.M, &page.Dm)
	}
	q29.Render(q, &page)
}

func Add(q *q29.ReqRsp) {
	var dm domain.Domain

	q.R.ParseForm()
	domname := q.R.FormValue("domain")
	emsg := validfield.Domain(validfield.F{"Domain Name", domname, 0, 0, true})
	if emsg != "" {
		q29.SetFlash(q, "The domain name "+domname+" is invalid.")
		q29.Redirect(q, "dlist/index")
		return
	}
	found := domain.Find(q.M, domname, &dm)
	if found == true {
		q29.SetFlash(q, "The domain "+domname+" is already under BlackList control.")				
		q29.Redirect(q, "dlist/index")
		return
	}
	if found == false {
		//nothing to do with "session", ShakeSalt is just a handy utility for generating
		vos := session.ShakeSalt(domname) 		
		//random, unique number for lookup later during confirm callback

		go sendEmailToDomainContact(domname, vos, q)		
		dm.Id = ""
		dm.Domain = domname
		dm.UserId = q.U.Id
		dm.Vos = vos
		dm.Confirm = false
		domain.Upsert(q.M, &dm)
	}
	q29.SetFlash(q, "The domain "+domname+" has been added.")	
	q29.Redirect(q, "dlist/index")
}

func Del(q *q29.ReqRsp) {
	var dm domain.Domain

	q.R.ParseForm()
	domname := q.R.FormValue("domain")
	emsg := validfield.Domain(validfield.F{"Domain Name", domname, 0, 0, true})
	if emsg == "" {
		found := domain.Find(q.M, domname, &dm)
		if found == true && dm.UserId == q.U.Id {
			emaddr.DeleteByDomainId(q.M, "domainaddrs", dm.Id)
			emaddr.DeleteByDomainId(q.M, "blacklist", dm.Id)
			domain.Delete(q.M, dm.Id)
		}
	}
	q29.SetFlash(q, "The domain "+domname+" has been removed.")
	q29.Redirect(q, "dlist/index")
}

func Elist(q *q29.ReqRsp) {
	var page struct {
		Vw q29.View
		Dname string
		Domain domain.Domain
		Elist []emaddr.Emaddr
		ElistCount int
		FlashMsg string
	}
	page.Dname = q.R.URL.Query().Get("dname")
	found := domain.Find(q.M, page.Dname, &page.Domain)
	if found == false || page.Domain.UserId != q.U.Id {
		http.Error(q.W, q.R.URL.Path+" invalid domain", 404)
		return
	}
	emaddr.ListByDid(q.M, "domainaddrs", page.Domain.Id, &page.Elist)
	page.ElistCount = len(page.Elist)
	q29.Render(q, &page)
}

func ElistAdd(q *q29.ReqRsp) {
	var emAddr emaddr.Emaddr
	var dm domain.Domain

	q.R.ParseForm()
	dname := q.R.FormValue("dname")
  email := q.R.FormValue("email")
	emsg := validfield.Email(validfield.F{"Email address", email, 0, 0, true})
	if emsg != "" {
		q29.SetFlash(q, "That email address was invalid.")
	} else {
		found := domain.Find(q.M, dname, &dm)
		emdom := strings.Split(email, "@")
		if found == false || dname != emdom[1] || dm.UserId != q.U.Id {
			http.Error(q.W, q.R.URL.Path+" invalid domain", 404)
			return
		}
		found = emaddr.FindByDid(q.M, "domainaddrs", dm.Id, email, &emAddr)
		if found == false {
			emAddr.Id = ""
			emAddr.Email = email
			emAddr.UserId = q.U.Id
			emAddr.DomainId = dm.Id
			emaddr.Upsert(q.M, "domainaddrs", &emAddr)
			emAddr.Id = ""			
			emaddr.Upsert(q.M, "blacklist", &emAddr)
		}
		q29.SetFlash(q, "The email address "+email+" has been added to this Domain BlackList.")
	}
	q29.Redirect(q, "dlist/elist?dname="+dname)
}

func ElistDel(q *q29.ReqRsp) {
	var emAddr emaddr.Emaddr
	var dm domain.Domain	
	q.R.ParseForm()
	dname := q.R.FormValue("dname")
  email := q.R.FormValue("email")
	emsg := validfield.Email(validfield.F{"Email address", email, 0, 0, true})
	if emsg != "" {
		q29.SetFlash(q, "That email address was invalid.")		
	} else {
		found := domain.Find(q.M, dname, &dm)
		if found == false || dm.UserId != q.U.Id {
			http.Error(q.W, q.R.URL.Path+" invalid domain", 404)
			return
		}
		found = emaddr.FindByDid(q.M, "domainaddrs", dm.Id, email, &emAddr)
		if found == false {
			q29.SetFlash(q, "That email address was not found.");
		} else {
			emaddr.Delete(q.M, "domainaddrs", emAddr.Id)
			q29.SetFlash(q, "The email address "+email+" has been removed from this Domain BlackList.")		
		}
		found = emaddr.Find(q.M, "blacklist", email, &emAddr)
		if found == true {
			emaddr.Delete(q.M, "blacklist", emAddr.Id)
		}		
	}
	q29.Redirect(q, "dlist/elist?dname="+dname)
}
