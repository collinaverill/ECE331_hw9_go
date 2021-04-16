package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"math/rand"
	"net/http"
)

const (
	listen = ":9090"
)

// type DataPoint represents a single value to be plotted on an x-y coordinate axis
// Value will be on the y axis
// Time  will be on the x axis
type DataPoint struct {
	Value float64 `json:"y"`
	Time  int64   `json:"x"`
}

func main() {

	reload := flag.Bool("r", false, "pass this flag to have the template reloaded on each request.")
	flag.Parse()

	log.Print("Opening Database")
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Print("Done Opening Database")

	tmpl, err := template.ParseFiles("./template.html")
	if err != nil {
		panic(err)
	}

	m := http.NewServeMux()
	m.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if *reload {
			tmpl, err = template.ParseFiles("./template.html")
			if err != nil {
				log.Printf("problem loading template: %s", err)
				writer.WriteHeader(500)
				return
			}
		}
		err := tmpl.Execute(writer, nil)
		if err != nil {
			log.Printf("problem executing template: %s", err)
			writer.WriteHeader(500)
			return
		}
	})

	m.HandleFunc("/data", func(writer http.ResponseWriter, request *http.Request) {

		rows, err := db.Query("select * from data")
		if err != nil {
			log.Printf("unable to query data - %s", err)
			writer.WriteHeader(500)
			return
		}
		defer rows.Close()

		var points []DataPoint

		i := 0
		for rows.Next() {
			var date string
			var time string
			var reconfirmed int
			err = rows.Scan(&date, &time, &reconfirmed)
			if err != nil {
				log.Printf("unable to scan row result %d - %s", i, err)
				writer.WriteHeader(500)
				return
			}
			points = append(points, DataPoint{
				Value: float64(reconfirmed),
				Time:  int64(i),
			})
			i++
			//fmt.Println(date, time, reconfirmed)
		}
		err = rows.Err()
		if err != nil {
			log.Printf("row error - %s", err)
			writer.WriteHeader(500)
			return
		}

		jsonBytes, err := json.MarshalIndent(points, "", "    ")
		if err != nil {
			writer.WriteHeader(500)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err = fmt.Fprint(writer, string(jsonBytes))
		if err != nil {
			writer.WriteHeader(500)
			return
		}

	})

	fmt.Printf("Server listening - %s\n", listen)
	err = http.ListenAndServe(listen, m)
	if err != nil {
		panic(err)
	}
}

// getValue produces a random 64bit floating point number on the range [-plusMinus, plusMinus)
func getValue(plusMinus float64) float64 {
	return (rand.Float64() * plusMinus * 2) - plusMinus
}
