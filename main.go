package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/codemodifyprivate/go-exercise-csv-1/helpers"
)

var flagDir = flag.String("dir", "", "directory containing CSVs")

func main() {
	flag.Parse()

	if *flagDir == "" {
		fmt.Println("missing -dir flag")
		os.Exit(1)
	}

	var csvs []*csv.Reader
	files, err := ioutil.ReadDir(*flagDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range files {
		if filepath.Ext(fi.Name()) == ".csv" {
			f, err := os.Open(filepath.Join(*flagDir, fi.Name()))
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			csvs = append(csvs, csv.NewReader(f))
		}
	}

	// read and consolidate csvs...
	oneBigCSV := helpers.ProcessCSVs(csvs)
	helpers.PrintCSVFile(oneBigCSV)
}
