package main

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/lib/pq"
)

type Query struct {
	Text string
}

func (q *Query) Expand(command string) {
	q.Text += command
}

func (q *Query) Execute(DB *sql.DB) error {
	var result string
	err := DB.QueryRow(q.Text + `;`).Scan(&result)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	fmt.Println(result)

	return nil
}

type Item struct {
	ItemID      int    `json:"item_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Price       int    `json:"price,omitempty"`
	PhotoPath   string `json:"photo_path,omitempty"`
	IsAvailable bool   `json:"is_available,omitempty"`
}

func main() {

	// Подключаемся к БД
	DB, err := sql.Open("postgres", "host=::1 port=5432 user=postgres password=1111 dbname=portal_KD sslmode=disable")
	if err != nil {
		fmt.Errorf("%s", err)
	}

	item := Item{
		ItemID: 1,
	}

	var query Query
	query.Expand(SELECT(`name`, getType(item)))
	query.Expand(WHERE(`item_id`, `1`))
	fmt.Println(query.Text)
	query.Execute(DB)
}

func getType(table interface{}) string {
	if t := reflect.TypeOf(table); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func SELECT(field, table string) string {
	return ` SELECT ` + field + ` FROM ` + table
}

func WHERE(field, value string) string {
	return ` WHERE ` + field + ` = ` + value
}
