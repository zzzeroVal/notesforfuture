package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
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

func SearchServer(w http.ResponseWriter, r *http.Request) {
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

	query := r.URL.Query().Get("query")
	var acceptUsers []User

	for _, user := range users {
		if strings.Contains(user.Name, query) || strings.Contains(user.About, query) {
			acceptUsers = append(acceptUsers, user)
		}
	}
	orderField := r.URL.Query().Get("order_field")
	orderBy := r.URL.Query().Get("order_by")

	orderByInt, err := strconv.Atoi(orderBy)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch orderField {
	case "Age":
		switch orderByInt {
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
		switch orderByInt {
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
		switch orderByInt {
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
		// доделать
	}

	offset := r.URL.Query().Get("offset")

	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limit := r.URL.Query().Get("limit")

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	end := limitInt + offsetInt

	if end > len(acceptUsers) {
		end = len(acceptUsers)
	}
	if offsetInt < len(acceptUsers) {
		acceptUsers = acceptUsers[offsetInt:end]
	} else {
		acceptUsers = nil
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(acceptUsers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func main() {
	http.HandleFunc("/searchServer", SearchServer)

	if err := http.ListenAndServe(":9090", nil); err != nil {
		fmt.Println("Ошибка при работе с HTTP сервером", err)
	}
}
