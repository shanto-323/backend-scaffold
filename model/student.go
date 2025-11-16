package model

import (
	"fmt"

	"github.com/google/uuid"
)

type Student struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Roll int       `json:"roll"`
}

func (s *Student) Validate() error{
	if s.Name == "" || s.Roll < 0 {
		return fmt.Errorf("missing fields")
	}

	return  nil
}
