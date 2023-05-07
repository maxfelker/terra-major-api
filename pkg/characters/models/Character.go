package models

type Character struct {
	Name         string `json:"name"`
	Bio          string `json:"bio"`
	Age          int    `json:"age"`
	Strength     int    `json:"strength"`
	Intelligence int    `json:"intelligence"`
	Endurance    int    `json:"endurance"`
	Agility      int    `json:"agility"`
}
