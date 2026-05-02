package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
)

// RecordsToBytes converts records to a slice of bytes
func RecordsToBytes(records [][]string) ([]byte, error) {
	if records == nil || len(records) == 0 {
		return nil, errors.New("records cannot be nil or empty")
	}

	var buf bytes.Buffer

	csvWriter := csv.NewWriter(&buf)

	err := csvWriter.WriteAll(records)
	if err != nil {
		return nil, err
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
