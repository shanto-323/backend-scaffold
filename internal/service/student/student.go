package student

import "fmt"

type student struct{}

func NewService() Service {
	return &student{}
}

func (s *student) Create() {
	fmt.Println("This is Create Student Function")
}
