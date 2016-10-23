package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/adamryman/rello"
)

func main() {
	port := os.Getenv("WEBLOGGERPORT")
	_, _ = ioutil.TempFile("", "")
	http.HandleFunc("/", HandleWebhook)
	_ = http.ListenAndServe(":"+port, nil)
}

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var update rello.ChecklistWebhook
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&update)
	if err != nil {
		fmt.Println(err)
	}
	a := update.Action
	t := a.Type
	c := a.Data.CheckItem

	if t == "updateCheckItemStateOnCard" {
		fmt.Printf("%s is now %s\n", c.Name, c.State)
	} else {
		fmt.Println(t)
	}
}
