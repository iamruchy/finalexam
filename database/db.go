package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Conn() *sql.DB {
	if db != nil {
		return db
	}
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Fail connect database")
	}
	return db
}

func CreateTable() {
	createTb := `
	CREATE TABLE IF NOT EXISTS customers(
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
`
	_, err := Conn().Exec(createTb)
	if err != nil {
		log.Fatal("create table fail : ", err)

	}
	fmt.Println("create table success")
}

func InsertCustomer(name, email, status string) (int, error) {
	stmt, err := Conn().Prepare("insert into customers (name,email,status) values ($1,$2,$3) returning id")
	if err != nil {
		return -1, err
	}
	row := stmt.QueryRow(name, email, status)
	var id int
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func GetCustomer() (*sql.Rows, error) {
	stmt, err := Conn().Prepare("select id,name,email,status from customers")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func GetCustomerByID(id int) (*sql.Row, error) {
	stmt, err := Conn().Prepare("select id,name,email,status from customers where id = $1")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(id)
	return row, nil

}

func UpdateCustomer(id int, name, email, status string) error {
	stmt, err := Conn().Prepare("update customers set name = $1,email = $2,status = $3 where id = $4")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, email, status, id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCustomer(id int) error {
	stmt, err := Conn().Prepare("delete from customers where id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
