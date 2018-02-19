package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "nas1g0r3ng"
	DB_NAME     = "excite_mobile_dev"
)


//at first lets open/close database for each hit
func GetAllUser() []Person{

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	fmt.Println("# Querying")
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	var people []Person
	for rows.Next() {
		var uid int
		var username string
		var password string
		var firstname string
		var lastname string
		var address string
		var created time.Time
		var modified time.Time
		err = rows.Scan(&uid, &username, &password ,&firstname, &lastname, &address, &created, &modified)
		checkErr(err)
		fmt.Println("uid | username | password | firstname | lastname | address | created | modified ")
		fmt.Printf("%3v | %8v | %8v | %8v | %8v | %8v | %8v| %8v\n", uid, username, password, firstname, lastname, address, created, modified)

		//lets put the value inside person
		var person Person
		person.ID = username
		person.Password = password
		person.Firstname = firstname
		person.Lastname = lastname
		person.Address.City = address
		person.Address.State = address
		people = append(people, person)
	}
	return people
}

func insertUser(username string, password string, firstname string, lastname string, address string, dateCreated string){
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	fmt.Println("# Inserting values "+ username + " " +password + " " +firstname + " " +lastname + " " +address + " " +dateCreated)

	var lastInsertId int
	err = db.QueryRow("INSERT INTO userinfo(username,password,firstname, lastname,address,created,modified) VALUES($1,$2,$3,$4,$5,$6,$7) returning uid;",
		username, password, firstname, lastname, address, dateCreated, dateCreated).Scan(&lastInsertId)
	checkErr(err)
	fmt.Println("last inserted id =", lastInsertId)
}

func MainAuthentication() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	fmt.Println("# Inserting values")

	var lastInsertId int
	err = db.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "研发部门", "2012-12-09").Scan(&lastInsertId)
	checkErr(err)
	fmt.Println("last inserted id =", lastInsertId)

	fmt.Println("# Updating")
	stmt, err := db.Prepare("update userinfo set username=$1 where uid=$2")
	checkErr(err)

	res, err := stmt.Exec("astaxieupdate", lastInsertId)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect, "rows changed")

	fmt.Println("# Deleting")
	stmt, err = db.Prepare("delete from userinfo where uid=$1")
	checkErr(err)

	res, err = stmt.Exec(lastInsertId)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect, "rows changed")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
