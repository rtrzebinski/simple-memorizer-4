package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/csv"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/validation"
)

type ExportLessonCsvHandler struct {
	s Service
}

func NewExportLessonCsvHandler(s Service) *ExportLessonCsvHandler {
	return &ExportLessonCsvHandler{s: s}
}

func (h *ExportLessonCsvHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	lessonId, err := strconv.Atoi(req.URL.Query().Get("lesson_id"))
	if err != nil {
		log.Print(fmt.Errorf("failed to get a lesson_id: %w", err))

		// validate empty lesson if lesson_id is not present, this is for error messages consistency
		validator := validation.ValidateLessonIdentified(backend.Lesson{})

		if validator.Failed() {
			log.Print(fmt.Errorf("invalid input: %w", validator))

			res.WriteHeader(http.StatusBadRequest)

			encoded, err := json.Marshal(validator.Error())
			if err != nil {
				log.Print(fmt.Errorf("failed to encode ExportLessonCsvHandler HTTP response: %w", err))
				res.WriteHeader(http.StatusInternalServerError)

				return
			}

			_, err = res.Write(encoded)
			if err != nil {
				log.Print(fmt.Errorf("failed to write ExportLessonCsvHandler HTTP response: %w", err))
				res.WriteHeader(http.StatusInternalServerError)

				return
			}

			return
		}
	}

	// Hydrate lesson
	lesson := backend.Lesson{Id: lessonId}
	err = h.s.HydrateLesson(ctx, &lesson, "")
	if err != nil {
		log.Print(fmt.Errorf("failed to hydrate lesson: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	// Set the oldest exercise ID to 1, as we are exporting all exercises
	oldestExerciseID := 1

	// Fetch exercises of the lesson
	exercises, err := h.s.FetchExercises(ctx, backend.Lesson{Id: lessonId}, oldestExerciseID, "")
	if err != nil {
		log.Print(fmt.Errorf("failed to fetch exercises: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	// Sort by ID asc
	slices.SortFunc(exercises, func(a, b backend.Exercise) int {
		if a.Id < b.Id {
			return -1
		} else if a.Id > b.Id {
			return 1
		}
		return 0
	})

	// Create CSV records
	var records [][]string
	for _, exercise := range exercises {
		var record []string
		record = append(record, exercise.Question)
		record = append(record, exercise.Answer)
		records = append(records, record)
	}

	// Create CSV file content from records
	fileContent, err := csv.WriteAll(records)
	if err != nil {
		log.Print(fmt.Errorf("failed to create a CSV from exercises: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	// Set the appropriate headers for the HTTP response
	res.Header().Set("Content-Disposition", "attachment; filename="+lesson.Name+".csv")
	res.Header().Set("Content-Type", "application/octet-stream")
	res.Header().Set("Content-Length", strconv.Itoa(len(fileContent)))

	// Write the file content to the response
	_, err = res.Write(fileContent)
	if err != nil {
		log.Print(fmt.Errorf("failed to write ExportLessonCsvHandler HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
