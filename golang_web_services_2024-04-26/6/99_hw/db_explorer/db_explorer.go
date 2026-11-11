package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

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
	db   *sql.DB
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {

	info := Handler{Info: make(map[string][]Column),
		db: db,
	}
	var tableNames []string

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
		tableNames = append(tableNames, nameTable)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	for _, name := range tableNames {
		columnRows, err := db.Query("SHOW FULL COLUMNS FROM " + "`" + name + "`")
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
			info.Info[name] = append(info.Info[name], column)
		}
		err = columnRows.Err()
		if err != nil {
			columnRows.Close()
			return nil, err
		}
		columnRows.Close()
	}
	return info, nil
}

func (info Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var namesTable []string
	if r.Method == http.MethodGet && r.URL.Path == "/" {
		for k, _ := range info.Info {
			namesTable = append(namesTable, k)
		}
		sort.Strings(namesTable)
		data := map[string][]string{
			"tables": namesTable,
		}
		response := map[string]map[string][]string{
			"response": data,
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		return
	}
	tableName := strings.TrimPrefix(r.URL.Path, "/")
	_, ok := info.Info[tableName]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		response := map[string]string{
			"error": "unknown table",
		}
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		return
	}

	query, err := info.db.Query("SELECT * FROM " + tableName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer query.Close()
	records := make([]map[string]interface{}, 0)

	for query.Next() {
		rowColumnValues := make([]interface{}, len(info.Info[tableName]))
		scanDestinations := make([]interface{}, len(info.Info[tableName]))
		for i, _ := range rowColumnValues {
			scanDestinations[i] = &rowColumnValues[i]
		}
		err := query.Scan(scanDestinations...)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		record := make(map[string]interface{})
		columns := info.Info[tableName]
		for i, column := range columns {
			columnName := column.Field
			if rowColumnValues[i] == nil {
				record[columnName] = nil
				continue
			}
			data, ok := rowColumnValues[i].([]byte)
			if !ok {
				record[columnName] = rowColumnValues[i]
				continue
			}
			switch {
			case strings.HasPrefix(column.Type, "int"):
				record[columnName], err = strconv.ParseInt(string(data), 10, 64)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			case strings.HasPrefix(column.Type, "float"), strings.HasPrefix(column.Type, "double"):
				record[columnName], err = strconv.ParseFloat(string(data), 64)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			default:
				record[columnName] = string(data)
			}
		}
		records = append(records, record)
	}
	err = query.Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data := map[string][]map[string]interface{}{
		"records": records,
	}
	response := map[string]interface{}{
		"response": data,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
}
