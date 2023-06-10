package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/postgres"
	"log"
	"os"
	"reflect"
)

func main() {
	// connect DB
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}
	execute(db, "CapitalsSeed")
	os.Exit(0)
}

func execute(db *sql.DB, seedMethodNames ...string) {
	s := Seeder{db, postgres.NewWriter(db)}

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

type Seeder struct {
	db *sql.DB
	w  storage.Writer
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

func (s Seeder) CapitalsSeed() {
	lesson := models.Lesson{
		Name: "Capitals",
	}

	err := s.w.StoreLesson(&lesson)
	if err != nil {
		panic(err)
	}

	exercises := models.Exercises{
		models.Exercise{
			Lesson:   &models.Lesson{Id: 1},
			Question: "Poland",
			Answer:   "Warsaw",
		},
		models.Exercise{
			Lesson:   &models.Lesson{Id: 1},
			Question: "Germany",
			Answer:   "Berlin",
		},
		models.Exercise{
			Lesson:   &models.Lesson{Id: 1},
			Question: "France",
			Answer:   "Paris",
		},
		models.Exercise{
			Lesson:   &models.Lesson{Id: 1},
			Question: "Netherlands",
			Answer:   "Amsterdam",
		},
		models.Exercise{
			Lesson:   &models.Lesson{Id: 1},
			Question: "Spain",
			Answer:   "Madrid",
		},
	}

	for _, exercise := range exercises {
		err := s.w.StoreExercise(exercise)
		if err != nil {
			panic(err)
		}
	}
}
