package server

import (
	"awesomeProject1/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/search", searchHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, %s!", r.URL.Path[1:])
}

type test_struct struct {
	Test string
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	tfmQuery := models.TfmSearchQuery{}
	err = json.Unmarshal([]byte(body), &tfmQuery)
	if err != nil {
		panic(err)
	}
	fmt.Println(tfmQuery.Destination)
}
