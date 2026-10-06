package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strings"
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

	query := ""
	var acceptUsers []User

	for _, user := range users {
		if strings.Contains(user.Name, query) || strings.Contains(user.About, query) {
			acceptUsers = append(acceptUsers, user)
		}
	}
	orderField := "Id"
	orderBy := OrderByAsc

	switch orderField {
	case "Age":
		switch orderBy {
		case OrderByAsc:
			sort.Slice(acceptUsers, func(i, j int) bool {
				return acceptUsers[i].Age < acceptUsers[j].Age
			})
		case OrderByAsIs:
		case OrderByDesc:
			sort.Slice(acceptUsers, func(i, j int) bool {
				return acceptUsers[i].Age > acceptUsers[j].Age
			})
		}
	case "Id":
		switch orderBy {
		case OrderByAsc:
			sort.Slice(acceptUsers, func(i, j int) bool {
				return acceptUsers[i].Id < acceptUsers[j].Id
			})
		case OrderByAsIs:
		case OrderByDesc:
			sort.Slice(acceptUsers, func(i, j int) bool {
				return acceptUsers[i].Id > acceptUsers[j].Id
			})
		}
	case "Name", "":
		switch orderBy {
		case OrderByAsc:
			sort.Slice(acceptUsers, func(i, j int) bool {
				return acceptUsers[i].Name < acceptUsers[j].Name
			})
		case OrderByAsIs:
		case OrderByDesc:
			sort.Slice(acceptUsers, func(i, j int) bool {
				return acceptUsers[i].Name > acceptUsers[j].Name
			})
		}
	default:
		fmt.Printf("error: %v", orderField)
	}

	offset := 35
	limit := 4
	end := offset + limit
	if end > len(acceptUsers) {
		end = len(acceptUsers)
	}
	if offset < len(acceptUsers) {
		acceptUsers = acceptUsers[offset:end]
	} else {
		acceptUsers = nil
	}
	fmt.Println(len(acceptUsers))
}
