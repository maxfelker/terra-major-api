package models

import "time"

type Instance struct {
	PrefabName string    `json:"prefabName"`
	OwnerId    int       `json:"ownerId"`
	Created    time.Time `json:"created"`
	Modified   time.Time `json:"modified"`
	Health     int       `json:"health"`
	Position   Vector3   `json:"position"`
	Rotation   Vector3   `json:"rotation"`
}
