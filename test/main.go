package main

import (
	"fmt"

	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/utils"
)

type Person struct {
	Name string
	age  int
}

type House struct {
	Person *Person
	street string
}

type Family struct {
	Person   *Person
	children int
}

func NewFamily(person *Person, child int) *Family {
	return &Family{
		Person:   person,
		children: child,
	}
}

func testHashValue() {
	var pepe utils.HashValue
	pepe.Set("0c515910c21c81b00d899705c2da2afc70db2d0c5b29d4293f5e698fd5afa5c0")
	fmt.Println(pepe.String())
}
func testFiles() {
	logger.CreateLogger("file", "0.0.0.0", true)
	// f, err := file.NewFileFromLocalSync("test.png")
	// err = f.Reconstruct()
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func main() {
	house := &House{
		Person: &Person{Name: "adri", age: 20},
		street: "calle",
	}

	family := NewFamily(house.Person, 2)

	fmt.Println(*house.Person)

	fmt.Println(*family.Person)

	house.Person.Name = "Pepe"
	fmt.Println(*house.Person)
	fmt.Println(*family.Person)

}
