package tests

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"

	"github.com/codemodifyprivate/go-exercise-csv-1/helpers"
)

func TestParseOneFile(t *testing.T) {
	f, _ := os.Open("../data/simple/1.csv")
	defer f.Close()

	csvFile := helpers.ProcessCSVs([]*csv.Reader{
		csv.NewReader(f),
	})

	for key, val := range *csvFile {
		fmt.Println(key, val)
	}
}
