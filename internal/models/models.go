package models

import (
	"strconv"
	"time"
)

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
	SlugName  string    `json:"slug_name"`
	Assigment bool      `json:"assigment"`
	DT        time.Time `json:"dt"`
}

func (a *Archive) Write() []string {
	return []string{
		strconv.Itoa(a.UserId),
		a.SlugName,
		strconv.FormatBool(a.Assigment),
		a.DT.Format("1/_2/2006 15:04:05")}
}

type AvitoUser struct {
	Id int `json:"id"`
}

type AddDeleteSlugs struct {
	AddSlugs     []Slug `json:"add_slugs"`
	DeleteSlugs  []Slug `json:"delete_slugs"`
	SlugDuration int    `json:"duration_minutes,omitempty"`
}

type AddUsers struct {
	Count int `json:"user_cnt"`
}

type AddSlug struct {
	Name    string `json:"name"`
	Percent int    `json:"user_percent,omitempty"`
}
