package main

/*
 */

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

type Config struct {
	Database struct {
		Server       string `json:"server"`
		ServerIp     string `json:"serverIp"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		DatabaseName string `json:"databaseName"`
		DatabaseType string `json:"databaseType"`
	} `json:"database"`
}

type Customer struct {
	ID       int
	FNAME    string
	LNAME    string
	PHONE    string
	ENTRY_DT string
}

func PrintLine(txt string) { fmt.Println(txt) }

func LoadConfig(filename string) (Config, error) {
	var config Config
	var configFile, err = os.Open(filename)
	defer configFile.Close() // defer closing until the last line in function
	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func ReadCustomers(db *sql.DB, ids []int) (outCustomers []Customer, err error) {
	//LoadConfig("sqlQueries").sqlStatements.customer_read
	qry := "SELECT * from [websuppo_swadmin].[CUSTOMERS]"
	if len(ids) > 0 {
		filters := arrayToString(ids, ",")
		qry += fmt.Sprintf("ID in (%s);", filters)
	}

	rows, err := db.Query(qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var FNAME, LNAME, PHONE, ENTRY_DT string
		var ID int
		err := rows.Scan(&ID, &FNAME, &LNAME, &PHONE, &ENTRY_DT)
		if err != nil {
			return nil, err
		}

		outCustomers = append(outCustomers, Customer{ID: ID, FNAME: FNAME, LNAME: LNAME, PHONE: PHONE, ENTRY_DT: ENTRY_DT})
	}

	return outCustomers, nil
}

func main() {
	PrintLine("************************")
	PrintLine("Starting app")
	PrintLine("")

	config, _ := LoadConfig("config.json")

	PrintLine("")
	PrintLine("opening connection")

	connstr := fmt.Sprintf("server=%s;user id=%s;password=%s;", config.Database.Server, config.Database.Username, config.Database.Password)
	conn, errdb := sql.Open("mssql", connstr)
	if errdb != nil {
		fmt.Println(" Error open db:", errdb.Error())
	}
	defer conn.Close()

	var listOfCustomers, err = ReadCustomers(conn, nil)
	if err != nil {
		log.Fatal("Reading Customers Failed:", err.Error())
	}

	fmt.Printf("there are %d customers ", len(listOfCustomers))

	PrintLine("")
	PrintLine("************************")
	PrintLine("End of app")
}
