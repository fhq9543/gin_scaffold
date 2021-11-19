package model

type Demo struct {
	Model
	Describe string `json:"describe" search:"describe"`
}
