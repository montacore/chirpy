package main

import (
	"log"
	"github.com/google/uuid"
	"time"
	"encoding/json"
)

type User struct {
	ID				uuid.UUID		`json:"id"`
	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
	Email			string			`json:"email"`
}


