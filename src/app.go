package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var APP_LOG_FILE string = "/var/log/app.log"
var APP_ACCESS_LOG_FILE string = "/var/log/app_access.log"
var APP_ERROR_LOG_FILE string = "/var/log/app_error.log"
var APP_PID_FILE string = "/var/run/app.pid"
var APP_PORT string = "10099"
var DB_ADDRESS string = "127.0.0.1"
var DB_PORT string = "3306"
var DB_NAME string = ""
var DB_USERNAME string = ""
var DB_PASSWORD string = ""
var DB_ROOT_PASSWORD string = ""
var DB_TABLE string = "counter" // predefined
var APP_COMMON_MESSAGE string = "Welcome!"

var conn *sqlx.DB
var err error


func readConfiguration() {
	var content = []byte("")
	content, errNearConfFile := ioutil.ReadFile("app.conf")
	if errNearConfFile != nil {
		content, _ = ioutil.ReadFile("/etc/app/app.conf")
	}
	if string(content) != "" {
		for _, configLine := range strings.Split(string(content), "\n") {
			configLineArray := strings.Split(configLine, "=")
			if len(configLineArray) == 2 {
				switch strings.TrimSpace(configLineArray[0]) {
				case "log_file":
					APP_LOG_FILE = strings.TrimSpace(configLineArray[1])
				case "access_log_file":
					APP_ACCESS_LOG_FILE = strings.TrimSpace(configLineArray[1])
				case "error_log_file":
					APP_ERROR_LOG_FILE = strings.TrimSpace(configLineArray[1])
				case "pid_file":
					APP_PID_FILE = strings.TrimSpace(configLineArray[1])
				case "port":
					APP_PORT = strings.TrimSpace(configLineArray[1])
				case "db_address":
					DB_ADDRESS = strings.TrimSpace(configLineArray[1])
				case "db_port":
					DB_PORT = strings.TrimSpace(configLineArray[1])
				case "db_name":
					DB_NAME = strings.TrimSpace(configLineArray[1])
				case "db_username":
					DB_USERNAME = strings.TrimSpace(configLineArray[1])
				case "db_password":
					DB_PASSWORD = strings.TrimSpace(configLineArray[1])
				case "db_root_password":
					DB_ROOT_PASSWORD = strings.TrimSpace(configLineArray[1])
				case "message":
					APP_COMMON_MESSAGE = strings.TrimSpace(configLineArray[1])
				}
			}
		}
	}

	if os.Getenv("APP_LOG_FILE") != "" {
		APP_LOG_FILE = os.Getenv("APP_LOG_FILE")
	}
	if os.Getenv("APP_ACCESS_LOG_FILE") != "" {
		APP_ACCESS_LOG_FILE = os.Getenv("APP_ACCESS_LOG_FILE")
	}
	if os.Getenv("APP_ERROR_LOG_FILE") != "" {
		APP_ERROR_LOG_FILE = os.Getenv("APP_ERROR_LOG_FILE")
	}
	if os.Getenv("APP_PID_FILE") != "" {
		APP_PID_FILE = os.Getenv("APP_PID_FILE")
	}
	if os.Getenv("APP_PORT") != "" {
		APP_PORT = os.Getenv("APP_PORT")
	}
	if os.Getenv("DB_ADDRESS") != "" {
		DB_ADDRESS = os.Getenv("DB_ADDRESS")
	}
	if os.Getenv("DB_PORT") != "" {
		DB_PORT = os.Getenv("DB_PORT")
	}
	if os.Getenv("DB_NAME") != "" {
		DB_NAME = os.Getenv("DB_NAME")
	}
	if os.Getenv("DB_USERNAME") != "" {
		DB_USERNAME = os.Getenv("DB_USERNAME")
	}
	if os.Getenv("DB_PASSWORD") != "" {
		DB_PASSWORD = os.Getenv("DB_PASSWORD")
	}
	if os.Getenv("DB_ROOT_PASSWORD") != "" {
		DB_ROOT_PASSWORD = os.Getenv("DB_ROOT_PASSWORD")
	}
	if os.Getenv("APP_COMMON_MESSAGE") != "" {
		APP_COMMON_MESSAGE = os.Getenv("APP_COMMON_MESSAGE")
	}

	logRecordMaker(APP_LOG_FILE, "a", fmt.Sprintf("%s - %s", "Common Log File", APP_LOG_FILE))
	logRecordMaker(APP_LOG_FILE, "a", fmt.Sprintf("%s - %s", "Access Log File", APP_ACCESS_LOG_FILE))
	logRecordMaker(APP_LOG_FILE, "a", fmt.Sprintf("%s - %s", "Error Log File", APP_ERROR_LOG_FILE))
	logRecordMaker(APP_LOG_FILE, "a", fmt.Sprintf("%s - %s", "PID File", APP_PID_FILE))
	logRecordMaker(APP_LOG_FILE, "a", fmt.Sprintf("%s - %s", "Serving Port", APP_PORT))
	logRecordMaker(APP_LOG_FILE, "a", fmt.Sprintf("%s - %s", "Database Address", DB_ADDRESS))
	logRecordMaker(APP_LOG_FILE, "a", fmt.Sprintf("%s - %s", "Database Port", DB_PORT))
	logRecordMaker(APP_LOG_FILE, "a", fmt.Sprintf("%s - %s", "Database Name", DB_NAME))
	logRecordMaker(APP_LOG_FILE, "a", fmt.Sprintf("%s - %s", "Database User Name", DB_USERNAME))
	logRecordMaker(APP_LOG_FILE, "a", fmt.Sprintf("%s - %s", "Database Use Password", DB_PASSWORD))
	logRecordMaker(APP_LOG_FILE, "a", fmt.Sprintf("%s - %s", "Database Root Password", DB_ROOT_PASSWORD))
	logRecordMaker(APP_LOG_FILE, "a", fmt.Sprintf("%s - %s", "Message", APP_COMMON_MESSAGE))
}

var startTime time.Time = time.Now()

// LOGS Making
func logRecordMaker(logFilePath string, fileMode string, content string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log := fmt.Sprintf("[%s]: %s\n", timestamp, content)
	fmt.Print(log)
	writeToFile(logFilePath, log, fileMode)
}

func writeToFile(path string, content string, write_mode string) { // "w" or "a"
	var flag int
	switch write_mode {
	case "a":
		flag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	case "append":
		flag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	case "w":
		flag = os.O_RDWR | os.O_CREATE
	case "write":
		flag = os.O_RDWR | os.O_CREATE
	}

	f, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		log.Println(err)
	}
}


// HANDLERS Part
func webHandlerRoot(w http.ResponseWriter, r *http.Request) {
	if len(selectFromTable(conn, "counter", 1)) == 0 { // there is no record
		insertIntoTable(conn, "counter", "Count", 1)
	} else {
		old_visiter_number := selectFromTable(conn, "counter", 1)[0].Count
		updateTable(conn, "counter", 1, "Count", old_visiter_number + 1)
	}
	
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	visiter_number := selectFromTable(conn, "counter", 1)[0].Count
	log.Println("You're visiter number", visiter_number)
	fmt.Fprintf(w, fmt.Sprintf("Hi!\nYou came from %s\nYou're %d visiter.\nWelcome!\n", ip, visiter_number))
}



func webHandlerDrop(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/drop" && r.URL.Path != "/drop/" {
		http.Error(w, r.URL.Path, http.StatusNotFound)
		return
	}
	dropTable(conn, "counter")
	conn.Exec(schema)
	log.Println("Back to the root!..")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type Counter struct {
	ID    int    `db:"id"`
	Count int    `db:"count"`
}

var schema string = "CREATE TABLE `counter` (`id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY, `count` integer)"



func selectFromTable(conn *sqlx.DB, table_name string, id ...int) []Counter {
	var counters []Counter
	var rows *sqlx.Rows
	var err error
	if len(id) > 0 {
		rows, err = conn.Queryx(fmt.Sprintf("SELECT * FROM %s WHERE id=%d;", table_name, id[0]))
	} else {
		rows, err = conn.Queryx(fmt.Sprintf("SELECT * FROM %s;", table_name))
	}
	for rows.Next() {
		var counter Counter
		err = rows.StructScan(&counter)
		if err != nil {
			log.Println("error in rows scanning")
		} 
		log.Println("Selected {id, count}", counter)
		counters = append(counters, counter)
	}
	return counters
}


func insertIntoTable(conn *sqlx.DB, table_name string, field string, new_value int) error {
	_, err := conn.Queryx(fmt.Sprintf("INSERT INTO %s (%s) VALUES(%d);", table_name, field, new_value))
	if err != nil {
		log.Println("insert into table was failed:", err)
	}
	log.Println("insert into table was successful")
	return err
}

func updateTable(conn *sqlx.DB, table_name string, id int, field string, new_value int) error {
	_, err := conn.Queryx(fmt.Sprintf("UPDATE %s SET %s=%d where id=%d;", table_name, field, new_value, id))
	if err != nil {
		log.Println("update table was failed:", err)
	}
	log.Println("update table was successful")
	return err
}

func dropTable(conn *sqlx.DB, table_name string) error {
	_, err := conn.Queryx(fmt.Sprintf("DROP TABLE %s;", table_name))
	if err != nil {
		log.Println("drop table failed:", err)
	}
	log.Println("table was successfully dropped..")
	return err
}


func main() {
	writeToFile(APP_LOG_FILE, "", "w") // cleaning common log file
	readConfiguration()

	conn, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DB_USERNAME, DB_PASSWORD, DB_ADDRESS, DB_PORT, DB_NAME))
	if err != nil {
		log.Println("Can't connect to the database:", err)
		os.Exit(1)
	}
	defer conn.Close()

	conn.Exec(schema) // create counter table

	m := mux.NewRouter()
	m.HandleFunc("/", webHandlerRoot)
	m.HandleFunc("/{drop:drop(?:\\/)?}", webHandlerDrop)

	log.Println(fmt.Sprintf("Starting server on :%s port ... ", APP_PORT))
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", APP_PORT), m); err != nil {
		log.Fatal(err)
	}
}
