package api

import (
	"github.com/rtrzebinski/simple-memorizer-go/internal/api/handlers"
	"github.com/rtrzebinski/simple-memorizer-go/internal/storage"
	"log"
	"net/http"
)

func ListenAndServe(r storage.Reader, port string) {
	http.Handle(Exercises, handlers.NewExercisesHandler(r))

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
