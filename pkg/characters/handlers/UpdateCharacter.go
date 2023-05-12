package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func UpdateCharacter(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	characterId := vars["id"]
	var character = getCharacter(characterId)
	response, e := json.Marshal(character)
	if e != nil {
		log.Panic(e)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}
