// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"time"
)

type Objectid struct {
	ID         int32        `json:"id"`
	TemplateID int32        `json:"template_id"`
	Serial     int32        `json:"serial"`
	ObjectID   string       `json:"object_id"`
	CreatedAt  sql.NullTime `json:"created_at"`
	UpdatedAt  sql.NullTime `json:"updated_at"`
}

type Template struct {
	ID              int32     `json:"id"`
	TemplateName    string    `json:"template_name"`
	Description     string    `json:"description"`
	ParseType       int32     `json:"parse_type"`
	AutoGenObjectID bool      `json:"auto_gen_object_id"`
	OriginObjectID  string    `json:"origin_object_id"`
	ValueExtract    string    `json:"value_extract"`
	GoStruct        string    `json:"go_struct"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
