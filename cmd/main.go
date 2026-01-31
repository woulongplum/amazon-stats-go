package main

import (
	"fmt"
	"log"
	
	"github.com/woulongplum/amazon-stats-go/internal/loader"

)

func main() {
	records,err := loader.LoadCSV("../data/records.csv")
	if err != nil {
		log.Fatal(err)
	}
	for _ , r := range records {
		fmt.Println(r)
	}
}
