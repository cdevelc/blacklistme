package blacklistme

import "net/http"
import "q29"
import "blacklistme/controller/blist"
import "blacklistme/controller/ulist"
import "blacklistme/controller/dlist"
import "blacklistme/controller/account"
import "blacklistme/controller/support"
import "blacklistme/controller/services"
import "blacklistme/controller/api"
import "blacklistme/controller/payment"
import "os"
import "fmt"
import	"github.com/golang/glog"
import "flag"

func Init() {
	q29.AddSite("blacklistme","blacklistme", "localhost:27017/blacklistme", Dispatch)
}
func init() {
	glog.Info("Gratuitious message to allow me to import the glog pkg first \n");
	glog.Flush()
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARN|FATAL] -log_dir=[string]\n", )
	flag.PrintDefaults()
	os.Exit(2)
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
		case "upload":         blist.Upload(q)
		default:
			http.Error(q.W, q.R.URL.Path+" "+http.StatusText(404), 404)
		}

	case "services":
		services.Index(q)
		
	case "ulist":
		if !ulist.BeforeFilter(q) { return }
		switch (q.Action) {
		case "index":          ulist.Dashboard(q)
		case "dashboard":      ulist.Dashboard(q)
		case "profile":        ulist.Profile(q)
		case "apikey":         ulist.Apikey(q)
		case "apikeyregen":    ulist.ApikeyRegen(q)
		case "plist":          ulist.Plist(q)
		case "plistadd":       ulist.PlistAdd(q)
		case "plistdel":       ulist.PlistDel(q)
		default:
			http.Error(q.W, q.R.URL.Path+" "+http.StatusText(404), 404)
		}

	case "dlist":
		if !dlist.BeforeFilter(q) { return }
		switch (q.Action) {
		case "index":          dlist.Index(q)
		case "add":            dlist.Add(q)
		case "addconfirm":     dlist.AddConfirm(q)			
		case "del":            dlist.Del(q)
		case "elist":          dlist.Elist(q)
		case "elistadd":       dlist.ElistAdd(q)
		case "elistdel":       dlist.ElistDel(q)
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
		case "password":       account.Password(q)
		case "email":          account.Email(q)
		case "rename":         account.Rename(q)
		case "forgot":         account.Forgot(q)
		case "recovery":       account.Recovery(q)
		default:
			http.Error(q.W, q.R.URL.Path+" "+http.StatusText(404), 404)
		}

	case "payment":
		if !payment.BeforeFilter(q) { return }
		switch (q.Action) {
		case "enroll":         payment.Enroll(q)
		case "confirm":        payment.Confirm(q)			
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
		case "status":       support.Status(q)
		default:
			http.Error(q.W, q.R.URL.Path+" "+http.StatusText(404), 404)
		}

		
	default: 
		http.Error(q.W, q.R.URL.Path+" "+http.StatusText(404), 404)
	}
}
