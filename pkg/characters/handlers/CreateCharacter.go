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

func CreateCharacter(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var newCharacter models.Character
	err := decoder.Decode(&newCharacter)
	//var newCharacter = createCharacter("test")
	//response, e := json.Marshal(newCharacter)
	if err != nil {
		log.Panic(err)
	}

	response, e := json.Marshal(newCharacter)

	if e != nil {
		log.Panic(e)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	writer.Write(response)
}
