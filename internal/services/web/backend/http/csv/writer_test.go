package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteAll(t *testing.T) {
	var records [][]string
	records = append(records, []string{"foo", "bar"})

	res, err := WriteAll(records)

	assert.NoError(t, err)
	assert.Equal(t, "foo,bar\n", string(res))
}
