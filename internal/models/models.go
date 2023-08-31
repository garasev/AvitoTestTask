package models

import "time"

type Slug struct {
	Name string `json:"name"`
}

type UserSlug struct {
	UserId int       `json:"user_id"`
	SlugId int       `json:"slug_id"`
	DTEnd  time.Time `json:"dt_end"`
}

type Archive struct {
	UserId    int       `json:"user_id"`
	SlugId    int       `json:"slug_id"`
	Assigment bool      `json:"assigment"`
	DT        time.Time `json:"dt"`
}

type AvitoUser struct {
	Id int `json:"id"`
}
