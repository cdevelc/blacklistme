package blacklistme

import "testing"
import "log"
import "net/http"

func TestInit(t *testing.T) {
	Init()
	log.Panic(http.ListenAndServe(":1337", nil))
}
