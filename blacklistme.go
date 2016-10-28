package blacklistme

import "net/http"
import "q29"
import "blacklistme/controller/blist"
import "blacklistme/controller/ulist"
import "blacklistme/controller/account"
import "blacklistme/controller/support"
import "blacklistme/controller/api"

func Init() {
	q29.AddSite("blacklistme","blacklistme", "localhost:27017/blacklistme", Dispatch)
}

func Dispatch(q *q29.ReqRsp) {

	switch (q.Controller) {
  case "home": fallthrough
	case "blist":
		if !blist.BeforeFilter(q) { return }
		switch (q.Action) {
		case "index":          blist.Index(q)
		case "inquire":        blist.Inquire(q)
		case "addrem":         blist.AddRem(q)
		case "complete":       blist.Complete(q)
		case "dump":           blist.Dump(q)
		default:
			http.Error(q.W, q.R.URL.Path+" "+http.StatusText(404), 404)
		}

	case "ulist":
		if !ulist.BeforeFilter(q) { return }
		switch (q.Action) {
		case "index":          ulist.Dashboard(q)
		case "dashboard":      ulist.Dashboard(q)
		case "profile":        ulist.Profile(q)
		case "apikey":         ulist.Apikey(q)
		case "apikeyregen":    ulist.ApikeyRegen(q)
		default:
			http.Error(q.W, q.R.URL.Path+" "+http.StatusText(404), 404)
		}

	case "account":
		if !account.BeforeFilter(q) { return }
		switch (q.Action) {
		case "login":          account.Login(q)
		case "logout":         account.Logout(q)
		case "register":       account.Register(q)
		case "thanks":         account.Thanks(q)
		case "confirm":        account.Confirm(q)
		case "forgot":         account.Forgot(q)
		case "password":       account.Password(q)
		case "email":          account.Email(q)
		case "rename":         account.Rename(q)			
		default:
			http.Error(q.W, q.R.URL.Path+" "+http.StatusText(404), 404)
		}
		
	case "api":
		if !api.BeforeFilter(q) { return }
		switch (q.Action) {
		case "index":          api.Index(q)
		default:
			http.Error(q.W, q.R.URL.Path+" "+http.StatusText(404), 404)
		}

	case "support":
		if !support.BeforeFilter(q) { return }
		switch (q.Action) {
		case "faq":          support.FAQ(q)
		default:
			http.Error(q.W, q.R.URL.Path+" "+http.StatusText(404), 404)
		}

		
	default: 
		http.Error(q.W, q.R.URL.Path+" "+http.StatusText(404), 404)
	}
}
