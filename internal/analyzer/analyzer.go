package analyzer

import (
	"github.com/woulongplum/amazon-stats-go/internal/loader"
)

//月別集計を作る関数
func Monthly(records []loader.Record) map[string]int {
	result := make(map[string]int)

	for _, r := range records {
		key := r.Date.Format("2006-01")
		result[key] += r.Amount

	}

	return result
}

//曜日別集計を作る関数
func Weekly(records []loader.Record) map[string]int {
	result := make(map[string]int)

	for _ , r := range records {
		key := r.Date.Weekday().String()
		result[key] += r.Amount
	}
	return result
}

// 月 × 曜日クロス集計を作る関数
func MonthlyWeeklyCross(records []loader.Record) map[string]map[string]int {
	result := make(map[string]map[string]int)
	
	for _,r := range records {
		month := r.Date.Format("2006-01")
		weekday := r.Date.Weekday().String()

		if _, ok := result[month]; !ok {
			result[month] = make(map[string]int)
			
		}
		result[month][weekday] += r.Amount
	}
	return result
}


//月×曜日の比率のクロス集計を作る関数
func MonthlyWeeklyRatio(records []loader.Record) map[string]map[string]float64 {
	result := make(map[string]map[string]float64)
	
	cross := MonthlyWeeklyCross(records)
	MonthlyTotals := make(map[string]int)

	for month,weekData := range cross {
		for _ , total := range weekData {
			MonthlyTotals[month] += total
		}

		if _, ok := result[month]; !ok {
			result[month] = make(map[string]float64)
		}

		for weekday,total := range weekData {
			ratio := float64(total) / float64(MonthlyTotals[month]) * 100
			result[month][weekday] = ratio
		}
	}
	  return result
	}
