package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
)

// ReadAll converts a slice of bytes to records
func ReadAll(data []byte) ([][]string, error) {
	reader := bytes.NewReader(data)
	csvReader := csv.NewReader(reader)

	var res [][]string
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return res, fmt.Errorf("failed to read CSV: %w", err)
		}

		res = append(res, rec)
	}

	return res, nil
}
