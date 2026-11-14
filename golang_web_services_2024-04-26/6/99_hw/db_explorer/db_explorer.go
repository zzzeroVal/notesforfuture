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

func (info Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var namesTable []string
	if req.Method == http.MethodGet && req.URL.Path == "/" {
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
		resp.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(resp).Encode(response)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		return
	}
	tableName := strings.TrimPrefix(req.URL.Path, "/")
	tablePaths := strings.Split(tableName, "/")
	tableName = tablePaths[0]
	var recordID string

	if len(tablePaths) > 1 {
		recordID = tablePaths[1]
	}
	_, ok := info.Info[tableName]
	if !ok {
		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusNotFound)
		response := map[string]string{
			"error": "unknown table",
		}
		err := json.NewEncoder(resp).Encode(response)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		return
	}

	var primaryKey string //  id || user_id
	for _, val := range info.Info[tableName] {
		if val.Key == "PRI" {
			primaryKey = val.Field
		}
	}
	limit := req.URL.Query().Get("limit")
	offset := req.URL.Query().Get("offset")
	var err error
	limitDeff := 5
	offsetDeff := 0

	if limit != "" {
		parsedLimit, err := strconv.Atoi(limit) // временная переменная сохраняет default значение от перезаписи
		if err == nil {
			limitDeff = parsedLimit
		}
	}
	if offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err == nil {
			offsetDeff = parsedOffset
		}
	}

	var query *sql.Rows
	if recordID != "" {
		query, err = info.db.Query("SELECT * FROM "+tableName+" WHERE "+primaryKey+"= ?", recordID)
	} else {
		query, err = info.db.Query("SELECT * FROM "+tableName+" LIMIT ? OFFSET ?"+"",
			limitDeff, offsetDeff,
		)
	}
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
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
			resp.WriteHeader(http.StatusInternalServerError)
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
					resp.WriteHeader(http.StatusInternalServerError)
					return
				}
			case strings.HasPrefix(column.Type, "float"), strings.HasPrefix(column.Type, "double"):
				record[columnName], err = strconv.ParseFloat(string(data), 64)
				if err != nil {
					resp.WriteHeader(http.StatusInternalServerError)
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
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	if recordID != "" && len(records) == 0 {
		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusNotFound)
		response := map[string]string{
			"error": "record not found",
		}
		err = json.NewEncoder(resp).Encode(response)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	var response map[string]interface{}

	if recordID != "" {
		data := map[string]interface{}{
			"record": records[0],
		}
		response = map[string]interface{}{
			"response": data,
		}
	} else {
		data := map[string][]map[string]interface{}{
			"records": records,
		}
		response = map[string]interface{}{
			"response": data,
		}
	}
	err = json.NewEncoder(resp).Encode(response)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
}
