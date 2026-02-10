package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/woulongplum/amazon-stats-go/internal/loader"
	"github.com/woulongplum/amazon-stats-go/internal/analyzer"
)

func main() {

	// CSVファイルを読み込む
	records,err := loader.LoadCSV("../data/records.csv")
	if err != nil {
		log.Fatal(err)
	}

	// 月別集計を表示する
	monthly := analyzer.Monthly(records)
	for month, total := range monthly { 
		fmt.Println(month, total) 
	}

	// 曜日別集計を表示する
	weekly := analyzer.Weekly(records)
	for day, total := range weekly {
		fmt.Println(day, total)
	}
	
	// 月×曜日クロス集計を表示する
	MonthlyAndWeekly := analyzer.MonthlyWeeklyCross(records)

	months := make([]string, 0, len(MonthlyAndWeekly))
	for month := range MonthlyAndWeekly {
		months = append(months, month)
	}

	sort.Strings(months)

	jpWeek := map[string]string{ 
		"Monday": "月曜日", "Tuesday": "火曜日", "Wednesday": "水曜日", "Thursday": "木曜日", "Friday": "金曜日", "Saturday": "土曜日", "Sunday": "日曜日", 
	} 
		
		// 曜日の並び順（英語のキーで並べる） 
	weekOrder := []string{ "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday", }

	for _, month := range months { 
		fmt.Println("月:", month) 
		weekData := MonthlyAndWeekly[month] 
		
		for _, day := range weekOrder { 
			if total, ok := weekData[day]; ok { 
				fmt.Printf(" %s: %d\n", jpWeek[day], total) 
			} 
		}
	}

	// 月×曜日比率クロス集計を表示する
	MonthlyWeeklyRatios := analyzer.MonthlyWeeklyRatio(records)

	fmt.Println("月×曜日比率クロス集計:")

	// ① 月のキーを集める 
	
	months = make([]string, 0, len(MonthlyWeeklyRatios))

	for month := range MonthlyWeeklyRatios { 
		months = append(months, month) 
	} 

	sort.Strings(months)

	for _, month := range months { 
		fmt.Println("月:", month) 
		weekData := MonthlyWeeklyRatios[month] 
		
		for _, day := range weekOrder { 
			if ratio, ok := weekData[day]; ok { 
				fmt.Printf(" %s: %.2f%%\n", jpWeek[day], ratio)
		  } 
		}
		fmt.Println()
	}
}	
