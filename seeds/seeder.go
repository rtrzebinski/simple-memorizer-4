package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"reflect"
)

type Seeder struct {
	db *sql.DB
}

func seed(s Seeder, seedMethodName string) {
	// Get the reflection value of the method
	m := reflect.ValueOf(s).MethodByName(seedMethodName)
	// Exit if the method doesn't exist
	if !m.IsValid() {
		log.Fatal("No method called ", seedMethodName)
	}
	// Execute the method
	m.Call(nil)
}

func (s Seeder) LessonSeed() {
	//prepare the statement
	stmt, err := s.db.Prepare(`INSERT INTO lesson(name) VALUES ($1)`)
	if err != nil {
		panic(err)
	}

	lessons := []string{"Capitals"}

	for _, n := range lessons {
		// execute query
		_, err = stmt.Exec(n)
		if err != nil {
			panic(err)
		}
	}
}

func (s Seeder) ExerciseSeed() {
	//prepare the statement
	stmt, err := s.db.Prepare(`INSERT INTO exercise(question, answer, lesson_id) VALUES ($1, $2, $3)`)
	if err != nil {
		panic(err)
	}

	capitals := map[string]string{
		"Poland":      "Warsaw",
		"Germany":     "Berlin",
		"France":      "Paris",
		"Netherlands": "Amsterdam",
		//"Spain":       "Madrid",
		//"Greece":      "Athens",
		//"Slovakia":    "Bratislava",
		//"Hungary":     "Budapest",
		//"Slovenia":    "Ljubljana",
		//"Cyprus":      "Nicosia",
		//"Iceland":     "Reykjavik",
		//"Latvia":      "Riga",
		//"Bulgaria":    "Sofia",
	}

	for q, a := range capitals {
		// execute query
		_, err = stmt.Exec(q, a, 1)
		if err != nil {
			panic(err)
		}
	}
}

func execute(db *sql.DB, seedMethodNames ...string) {
	s := Seeder{db}

	seedType := reflect.TypeOf(s)

	// Execute all seeders if no method name is given
	if len(seedMethodNames) == 0 {
		log.Println("Running all seeders...")
		// We are looping over the method on a Seeder struct
		for i := 0; i < seedType.NumMethod(); i++ {
			// Get the method in the current iteration
			method := seedType.Method(i)
			// Execute seeder
			seed(s, method.Name)
		}
	}

	// Execute only the given method names
	for _, item := range seedMethodNames {
		seed(s, item)
	}
}

func main() {
	// connect DB
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}
	execute(db, "LessonSeed", "ExerciseSeed")
	os.Exit(0)
}
