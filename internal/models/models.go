package models

import "time"

type Slug struct {
	Name string
}

type UserSlug struct {
	UserId int
	SlugId int
	DTEnd  time.Time
}

type Archive struct {
	UserId    int
	SlugId    int
	Assigment bool
	DT        time.Time
}
