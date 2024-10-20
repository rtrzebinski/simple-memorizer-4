package main

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage/postgres"
	"log"
	"os"
	"reflect"
	"time"
)

type config struct {
	Db struct {
		Driver string `envconfig:"DB_DRIVER" default:"postgres"`
		DSN    string `envconfig:"DB_DSN" default:"postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable"`
	}
}

func main() {
	// Configuration
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// connect DB
	db, err := sql.Open(cfg.Db.Driver, cfg.Db.DSN)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}
	execute(db, "CapitalsSeed")
	//execute(db, "LargeLessonSeed")
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

func randomString() string {
	return uuid.NewString()
}

type Seeder struct {
	db *sql.DB
	w  *postgres.Writer
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
		Name:        "Capitals",
		Description: "What is the capital of given country?",
	}

	err := s.w.UpsertLesson(&lesson)
	if err != nil {
		panic(err)
	}

	exercises := models.Exercises{
		models.Exercise{
			Lesson:   &models.Lesson{Id: lesson.Id},
			Question: "Poland",
			Answer:   "Warsaw",
		},
		models.Exercise{
			Lesson:   &models.Lesson{Id: lesson.Id},
			Question: "Germany",
			Answer:   "Berlin",
		},
		models.Exercise{
			Lesson:   &models.Lesson{Id: lesson.Id},
			Question: "France",
			Answer:   "Paris",
		},
		models.Exercise{
			Lesson:   &models.Lesson{Id: lesson.Id},
			Question: "Netherlands",
			Answer:   "Amsterdam",
		},
		models.Exercise{
			Lesson:   &models.Lesson{Id: lesson.Id},
			Question: "Spain",
			Answer:   "Madrid",
		},
	}

	for _, exercise := range exercises {
		err := s.w.UpsertExercise(&exercise)
		if err != nil {
			panic(err)
		}
	}
}

func (s Seeder) LargeLessonSeed() {
	exercisesCount := 100
	answersCount := 100

	lesson := models.Lesson{
		Name:        "Large lesson",
		Description: "This lesson has plenty of exercises and answers",
	}

	err := s.w.UpsertLesson(&lesson)
	if err != nil {
		panic(err)
	}

	exercises := models.Exercises{}

	for i := exercisesCount; i > 0; i-- {
		exercises = append(exercises, models.Exercise{
			Lesson:   &models.Lesson{Id: lesson.Id},
			Question: randomString(),
			Answer:   randomString(),
		})
	}

	for k := range exercises {
		err := s.w.UpsertExercise(&exercises[k])
		if err != nil {
			panic(err)
		}
	}

	for i := range exercises {
		for j := answersCount; j > 0; j-- {
			answer := &models.Result{
				Exercise: &exercises[i],
			}

			currentTime := time.Now().UnixNano() / int64(time.Millisecond)
			randomBool := currentTime%2 == 0

			if randomBool == true {
				answer.Type = models.Good
			} else {
				answer.Type = models.Bad
			}

			err := s.w.StoreResult(answer)
			if err != nil {
				panic(err)
			}
		}
	}
}
