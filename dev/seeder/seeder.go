package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/auth"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/worker"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/postgres"
)

type config struct {
	Db struct {
		Driver string `envconfig:"DB_DRIVER" default:"postgres"`
		DSN    string `envconfig:"DB_DSN" default:"postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable&timezone=Europe/Warsaw"`
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
	kc := gocloak.NewClient("http://localhost:8180")
	authSvc := auth.NewService(kc, auth.Config{
		Realm:        "realm-dev",
		ClientID:     "client-id-dev",
		ClientSecret: "client-secret-dev",
	})

	s := Seeder{db, postgres.NewWebWriter(db), postgres.NewWorkerWriter(db), authSvc}

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
	db            *sql.DB
	backendWriter *postgres.WebWriter
	workerWriter  *postgres.WorkerWriter
	authSvc       *auth.Service
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
	ctx := context.Background()

	uEmail := "test.user.seeded@example.com"
	uPassword := "password"
	uFirstName := "Test User"
	uLastName := "Seeded"

	t, err := s.authSvc.Register(ctx, uFirstName, uLastName, uEmail, uPassword)
	if err != nil {
		panic(err)
	}

	uID := userID(t.AccessToken)

	lesson := backend.Lesson{
		Name:        "Capitals",
		Description: "What is the capital of given country?",
	}

	err = s.backendWriter.UpsertLesson(ctx, &lesson, uID)
	if err != nil {
		panic(err)
	}

	exercises := backend.Exercises{
		backend.Exercise{
			Lesson:   &backend.Lesson{Id: lesson.Id},
			Question: "Poland",
			Answer:   "Warsaw",
		},
		backend.Exercise{
			Lesson:   &backend.Lesson{Id: lesson.Id},
			Question: "Germany",
			Answer:   "Berlin",
		},
		backend.Exercise{
			Lesson:   &backend.Lesson{Id: lesson.Id},
			Question: "France",
			Answer:   "Paris",
		},
		backend.Exercise{
			Lesson:   &backend.Lesson{Id: lesson.Id},
			Question: "Netherlands",
			Answer:   "Amsterdam",
		},
		backend.Exercise{
			Lesson:   &backend.Lesson{Id: lesson.Id},
			Question: "Spain",
			Answer:   "Madrid",
		},
	}

	for _, exercise := range exercises {
		err := s.backendWriter.UpsertExercise(ctx, &exercise, uID)
		if err != nil {
			panic(err)
		}
	}
}

func (s Seeder) LargeLessonSeed() {
	ctx := context.Background()

	exercisesCount := 100
	answersCount := 100

	lesson := backend.Lesson{
		Name:        "Large lesson",
		Description: "This lesson has plenty of exercises and answers",
	}

	err := s.backendWriter.UpsertLesson(ctx, &lesson, "1")
	if err != nil {
		panic(err)
	}

	exercises := backend.Exercises{}

	for i := exercisesCount; i > 0; i-- {
		exercises = append(exercises, backend.Exercise{
			Lesson:   &backend.Lesson{Id: lesson.Id},
			Question: randomString(),
			Answer:   randomString(),
		})
	}

	for k := range exercises {
		err := s.backendWriter.UpsertExercise(ctx, &exercises[k], "1")
		if err != nil {
			panic(err)
		}
	}

	for i := range exercises {
		for j := answersCount; j > 0; j-- {
			r := worker.Result{
				ExerciseId: exercises[i].Id,
			}

			currentTime := time.Now().UnixNano() / int64(time.Millisecond)
			randomBool := currentTime%2 == 0

			if randomBool == true {
				r.Type = worker.Good
			} else {
				r.Type = worker.Bad
			}

			err := s.workerWriter.StoreResult(ctx, r)
			if err != nil {
				panic(err)
			}
		}
	}
}

func userID(accessToken string) string {
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		panic(err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		panic(err)
	}

	return claims["sub"].(string)
}
