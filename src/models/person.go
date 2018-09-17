package models

import "github.com/icrowley/fake"

// A non-player character
type Person struct {
	Name string
}

func GenPerson() Person {
	return Person{
		Name: fake.FirstName(),
	}
}
