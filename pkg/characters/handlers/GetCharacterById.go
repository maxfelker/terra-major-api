package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	models "github.com/mw-felker/centerpoint-instance-api/pkg/characters/models"
	utils "github.com/mw-felker/centerpoint-instance-api/pkg/utils"
)

func getCharacter(characterId string) models.Character {
	log.Println(characterId)
	var response models.Character
	jsonData, err := ioutil.ReadFile("./pkg/characters/handlers/mock-characters.json")
	if err != nil {
		utils.ErrorHandler(err)
	}
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		utils.ErrorHandler(err)
	}

	return response
}

func GetCharacterById(writer http.ResponseWriter, request *http.Request) {
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
