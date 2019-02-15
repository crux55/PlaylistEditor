package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type Entry struct {
	Uid          int	`json:"uid"`
	PlaylistName string `json:"name"`
	PlaylistUrl  string `json:"url"`
	Enabled      string	`json:"enabled"`
}

var ctx = context.Background()
var db *sql.DB

// our main function
func main() {

	var err error
	db, err = sql.Open("mysql", "username:*****@tcp(host:port)/database")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/playlisteditor", PlaylistEditor).Methods("POST")
	router.HandleFunc("/playlisteditor/add", Add).Methods("POST")
	router.HandleFunc("/playlisteditor/read", Read).Methods("GET")
	log.Fatal(http.ListenAndServe(":4100", router))
}

// The handler for the login endpoint
func PlaylistEditor(w http.ResponseWriter, r *http.Request) {

}

func Add(w http.ResponseWriter, r *http.Request) {
	entry := Entry{1, "new name", "new url", "no"}
	params := mux.Vars(r)
	entry.PlaylistUrl = params["url"]
	entry.PlaylistName = params["name"]
	entry.Enabled = params["enabled"]

	// replace all spaces with underscores
	if strings.Contains(entry.PlaylistName, " "){
		strings.Replace(entry.PlaylistName, " ", "_", -1)
	}
	if strings.Contains(entry.PlaylistUrl, " "){
		strings.Replace(entry.PlaylistUrl, " ", "_", -1)
	}

	strings.Trim(entry.PlaylistUrl, " ")
	strings.Trim(entry.PlaylistName, " ")

	if !strings.EqualFold(entry.Enabled, "yes") || !strings.EqualFold(entry.Enabled, "no"){
		entry.Enabled = "no"
	}

	query := fmt.Sprintf(`INSERT into playlists (name, url, enabled) values ("%s", "%s", "%s")`,
		entry.PlaylistName, entry.PlaylistUrl, entry.Enabled)

	row, err2 := db.QueryContext(ctx, query)
	if err2 != nil {
		panic(err2)
	}
	row.Close()
}

func Read(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(getAll())
}

func getAll() []Entry{
	rows, err := db.QueryContext(ctx, "SELECT * FROM playlists")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	entries := make([]Entry, 0)
	for rows.Next() {
		var name string
		var uid int
		var url string
		var enabled string
		if err := rows.Scan(&uid, &name, &url, &enabled); err != nil {
			log.Fatal(err)
		}
		entry := Entry{uid, name, url, enabled}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return entries
}
