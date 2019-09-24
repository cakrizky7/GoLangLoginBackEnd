package models

import (
	"time"

	"github.com/google/uuid"
)

func NewUsers() *Users {
	m := new(Users)
	m.Id = uuid.Must(uuid.NewRandom())
	return m
}

func (m *Users) TableName() string {
	return "users"
}

type Users struct {
	Id            uuid.UUID `json:"Id" validate:"required"`
	Username      string    `json:"Username" validate:"required"`
	Password      string    `json:"Password" validate:"required"`
	Nama_lengkap  string    `json:"Nama_lengkap" validate:"required"`
	Token         string    `json:"Token"`
	Token_expired time.Time `json:"Token_expired"`
	Created_at    time.Time `json:"Created_at" validate:"required"`
	Updated_at    time.Time `json:"Updated_at" validate:"required"`
}
