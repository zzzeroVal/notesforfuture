package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

type Item struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Updated     *string `json:"updated"`
}

type User struct {
	UserID   int     `json:"user_id"`
	Login    string  `json:"login"`
	Password string  `json:"password"`
	Email    string  `json:"email"`
	Info     string  `json:"info"`
	Updated  *string `json:"updated"`
}

type Column struct {
	Field   string
	Null    string
	Type    string
	Key     string
	Extra   string
	Default *string
}

type Handler struct {
	Info map[string][]Column
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {

	info := Handler{Info: make(map[string][]Column)}

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var nameTable string
		err := rows.Scan(&nameTable)
		if err != nil {
			return nil, err
		}
		columnRows, err := db.Query("SHOW FULL COLUMNS FROM " + "`" + nameTable + "`")
		if err != nil {
			return nil, err
		}
		for columnRows.Next() {
			var column Column
			var noUseColumn *string

			err := columnRows.Scan(&column.Field, &column.Type, &noUseColumn, &column.Null, &column.Key,
				&column.Default, &column.Extra, &noUseColumn, &noUseColumn)
			if err != nil {
				columnRows.Close()
				return nil, err
			}
			info.Info[nameTable] = append(info.Info[nameTable], column)
		}
		err = columnRows.Err()
		if err != nil {
			columnRows.Close()
			return nil, err
		}
		columnRows.Close()
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return info, nil
}
func (info Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var tablesName []string
	if r.Method == http.MethodGet && r.URL.Path == "/" {
		for k, _ := range info.Info {
			tablesName = append(tablesName, k)
		}
	}
	fmt.Println(tablesName)
}
