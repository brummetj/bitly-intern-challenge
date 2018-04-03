package bitly

import (
	"os"
	"log"
	"net/http"
	"encoding/json"
)

func GetBitlyToken() (string) {

	var token = os.Getenv("BITLY_ACCESS_TOKEN")
	if token == "" {
		log.Println("Request token variable not found")
		return ""
	}
	return token
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
