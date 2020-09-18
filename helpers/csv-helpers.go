package helpers

import (
	"encoding/csv"
	"sort"
	"strings"
	"sync"

	"fmt"
	"os"
)

func trimmedString(val string) string {
	return strings.TrimSpace(val)
}

// CSVLineHeader - colPosition -> colName
type CSVLineHeader map[string]bool

// CSVLine - colName -> colValue
type CSVLine map[string]string

// CSVFile - colID -> line
type CSVFile map[string]CSVLine

// ProcessCSVs -
func ProcessCSVs(csvReaders []*csv.Reader) *CSVFile {
	syncMutex := sync.Mutex{}
	syncWG := sync.WaitGroup{}
	oneBigCSV := CSVFile{}
	knownColumnName := CSVLineHeader{}

	updateKnownColumnName := func(line CSVLine) {
		for columnName := range line {
			knownColumnName[columnName] = true
		}
	}
	createNewTemplateLine := func() CSVLine {
		templateCSVLine := CSVLine{}
		for colName := range knownColumnName {
			templateCSVLine[colName] = "__blank__"
		}

		return templateCSVLine
	}
	setTemplateLineFromData := func(template *CSVLine, line CSVLine) {
		for columnName, columnValue := range line {
			(*template)[columnName] = columnValue
		}
	}

	for _, r := range csvReaders {
		// MAP
		syncWG.Add(1)
		go func(csvReader *csv.Reader) {
			lines, err := csvReader.ReadAll()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			var columnPositions = map[int]string{}
			var oneCSV = map[string]CSVLine{}

			for lineIndex, line := range lines {

				if lineIndex == 0 {

					//
					// header line
					//
					for columnIndex, columnName := range line { // scan column position
						columnName = trimmedString(columnName)
						columnPositions[columnIndex] = columnName
					}
				} else {

					//
					// data line
					//
					recordID := ""
					for columnIndex, columnValue := range line { // set column value by name
						if columnIndex == 0 { // ID column
							recordID = columnValue
							oneCSV[columnValue] = CSVLine{}
						} else {
							columnName := columnPositions[columnIndex]
							oneCSV[recordID][columnName] = columnValue
						}
					}
				}
			}

			// REDUCE
			for recordID, line := range oneCSV {
				syncMutex.Lock()

				if existinLine, ok := oneBigCSV[recordID]; ok { // check if a key is in dictionary

					updateKnownColumnName(line)

					templateCSVLine := createNewTemplateLine()
					setTemplateLineFromData(&templateCSVLine, existinLine)
					setTemplateLineFromData(&templateCSVLine, line)

					oneBigCSV[recordID] = templateCSVLine

				} else {

					updateKnownColumnName(line)

					templateCSVLine := createNewTemplateLine()
					setTemplateLineFromData(&templateCSVLine, line)

					oneBigCSV[recordID] = templateCSVLine
				}

				syncMutex.Unlock()
			}

			syncWG.Done()
		}(r)
	}

	syncWG.Wait()

	return &oneBigCSV
}

// PrintCSVFile -
func PrintCSVFile(csvFile *CSVFile) {
	if len(*csvFile) <= 0 {
		fmt.Println("EMPTY CSV")
	}

	// print header
	columnsInOrder := []string{}
	sb := strings.Builder{}
	sb.WriteString("ID | ")
	for _, columnData := range *csvFile {
		for key := range columnData {
			columnsInOrder = append(columnsInOrder, key)
			sb.WriteString(key + " | ")
		}

		break // not an actual loop, just get to the first line
	}
	csvHeader := sb.String()
	if len(csvHeader) > 0 {
		csvHeader = csvHeader[:len(csvHeader)-2]
	}
	fmt.Println(csvHeader)

	// sort keys
	sortedKeys := []string{}
	for id := range *csvFile {
		sortedKeys = append(sortedKeys, id)
	}
	sort.Strings(sortedKeys)

	// print data
	for _, sortedKey := range sortedKeys {
		csvLine := (*csvFile)[sortedKey]

		sb := strings.Builder{}
		for _, columName := range columnsInOrder {
			sb.WriteString(csvLine[columName] + " | ")
		}
		csvData := sb.String()
		if len(csvData) > 0 {
			csvData = csvData[:len(csvData)-2]
		}
		fmt.Printf("%s | %s\n", sortedKey, csvData)
	}

}
