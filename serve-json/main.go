package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type response struct {
	Error bool                   `json:"error"`
	Data  map[string]interface{} `json:"data"`
}

func main() {
	port := 3000

	http.HandleFunc("/json/user", jsonHandler)

	log.Println(fmt.Sprintf("Server listening on port %v", port))
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	resp := response{
		Error: false,
		Data: map[string]interface{}{
			"fname": "Lon",
			"lname": "Ilesanmi",
		},
	}

	// Using json.Marshal
	j, err := json.Marshal(resp)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(j))

	// Using j*son.Encoder to write directly to the response writer
	decoder := json.NewEncoder(w)
	decoder.Encode(&resp)
}
