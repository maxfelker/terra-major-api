package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	models "github.com/mw-felker/centerpoint-instance-api/pkg/characters/models"
	utils "github.com/mw-felker/centerpoint-instance-api/pkg/utils"
)

func getCharacters() []models.Character {
	var response []models.Character
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

func GetCharactersHandler(writer http.ResponseWriter, request *http.Request) {
	var characters = getCharacters()
	response, e := json.Marshal(characters)
	if e != nil {
		log.Panic(e)
	}

	//update content type
	writer.Header().Set("Content-Type", "application/json")

	//specify HTTP status code
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)

}
