package loader 

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

type Record struct {
	Date   time.Time
	Amount int
}

func LoadCSV(path string) ([]Record, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var records []Record
	for i, row := range rows {
		if i == 0 {
			continue 
		}
		date, err := time.Parse("2006-01-02", row[0])
		if err != nil {
			return nil, err
		}
		amount, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, err
		}
		records = append(records, Record{
			Date:   date,
			Amount: amount,
		})
	}
	return records, nil
}
