package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	models "github.com/mw-felker/centerpoint-instance-api/pkg/characters/models"
)

func createCharacter(name string) models.Character {

	return models.Character{
		Name:         name,
		Bio:          "",
		Age:          35,
		Strength:     5,
		Intelligence: 5,
		Endurance:    5,
		Agility:      5,
	}
}

func CreateCharacterHandler(writer http.ResponseWriter, request *http.Request) {
	var newCharacter = createCharacter("random")
	response, e := json.Marshal(newCharacter)
	if e != nil {
		log.Panic(e)
	}

	//update content type
	writer.Header().Set("Content-Type", "application/json")

	//specify HTTP status code
	writer.WriteHeader(http.StatusCreated)
	writer.Write(response)

}
