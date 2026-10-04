package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Person struct {
	Id        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

type Root struct {
	Row []Person `xml:"row"`
}

func main() {
	data, err := os.ReadFile("dataset.xml")
	if err != nil {
		panic(err)
	}
	var root Root

	err = xml.Unmarshal(data, &root)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	var users []User

	for _, person := range root.Row {
		user := User{
			Id:     person.Id,
			Name:   person.FirstName + " " + person.LastName,
			Age:    person.Age,
			About:  person.About,
			Gender: person.Gender,
		}
		users = append(users, user)
	}
	fmt.Println(len(users))
	fmt.Println(users[0])
}
