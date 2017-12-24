package output

import (
	"encoding/csv"
	"os"
)

// Write to csv file
func Write(filepath string, segment [][]string) error {
	f, err := os.Create(filepath)
	defer f.Close()
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	for _, val := range segment {
		if err := w.Write(val); err != nil {
			break
		}
	}

	w.Flush()
	return w.Error()
}
