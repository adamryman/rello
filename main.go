package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

func init() {
	_ = errors.Wrap
}

func main() {
	port := os.Getenv("WEBLOGGERPORT")
	_, _ = ioutil.TempFile("", "")
	http.HandleFunc("/", WriteToFile)
	_ = http.ListenAndServe(":"+port, nil)
}

func WriteToFile(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
}
