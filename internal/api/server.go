package api

import (
	"github.com/rtrzebinski/simple-memorizer-go/internal/api/methods"
	"github.com/rtrzebinski/simple-memorizer-go/internal/storage"
	"log"
	"net/http"
)

func ListenAndServe(r storage.Reader, port string) {
	http.Handle(RandomExercise, methods.NewExercisesHandler(r))

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
