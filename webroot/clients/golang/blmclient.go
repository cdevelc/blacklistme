package main

import "encoding/json"
import "fmt"
import "log"
import "net/http"

type BlackListMeResponse struct {
	Email string     `json:"email"`
	Blacklist string `json:"blacklist"`
	Error string     `json:"error"`	
}

func main() {
	server  := "www.dgears.com:1337"
	uri     := "/blacklistme/api"
	apikey  := "8078ac0a68877d828efccb68e91dabc1e720a716acb92d71e4f278428ec2a311"
  email   := "alane@fueledcafe.com"

	url := fmt.Sprintf("http://%s%s?email=%s&apikey=%s",server, uri, email, apikey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("http NewRequest failed: ", err)
		return
	}
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		log.Fatal("http GET failed: ", err)
		return
	}
	defer rsp.Body.Close()
	var jresponse BlackListMeResponse
	err = json.NewDecoder(rsp.Body).Decode(&jresponse)
	if err != nil {
		log.Fatal("json decode failed: ", err)
		return
	}
	if jresponse.Error != "" {
		log.Fatal("api query failed: ", jresponse.Error)
		return
	}
	fmt.Printf("response: %v\n", jresponse)
}
