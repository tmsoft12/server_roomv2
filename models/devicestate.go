package models

type StateDev struct {
	ID   int    `json:"id"`
	Door string `json:"door"`
	Fire string `json:"fire"`
	Pir  string `json:"pir"`
	Temp string `json:"temp"`
	Hum  string `json:"hum"`
}
