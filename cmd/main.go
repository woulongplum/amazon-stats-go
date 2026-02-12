package main

import (
	"fmt"
	"log"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"github.com/woulongplum/amazon-stats-go/internal/analyzer"
	"github.com/woulongplum/amazon-stats-go/internal/loader"
)

func main() {

	// CSVファイルを読み込む
	records, err := loader.LoadCSV("../data/records.csv")
	if err != nil {
		log.Fatal(err)
	}

	// 月別集計を表示する
	monthly := analyzer.Monthly(records)
	for month, total := range monthly {
		fmt.Println(month, total)
	}
	// 月別集計のグラフを作成する
	createMonthlyGraph(monthly)

	// 曜日別集計を表示する
	weekly := analyzer.Weekly(records)
	for day, total := range weekly {
		fmt.Println(day, total)
	}

	// 曜日別集計のグラフを作成する
	createWeeklyGraph(weekly)

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
	weekOrder := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	for _, month := range months {
		fmt.Println("月:", month)
		weekData := MonthlyAndWeekly[month]

		for _, day := range weekOrder {
			if total, ok := weekData[day]; ok {
				fmt.Printf(" %s: %d\n", jpWeek[day], total)
			}
		}
	}

	// 月×曜日クロス集計のグラフを作成する
	createMonthlyWeeklyCrossGraph(MonthlyAndWeekly)

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

	// 月×曜日比率クロス集計のグラフを作成する
	MonthlyWeeklyRatioGraph(MonthlyWeeklyRatios,"../output/month_week_ratio.png")

}

// 月別集計のグラフを作成する関数
func createMonthlyGraph(monthly map[string]int) {

	// 月のキーをソートして取得
	months := make([]string, 0, len(monthly))
	for m := range monthly {
		months = append(months, m)
	}
	// ソート
	sort.Strings(months)
	// グラフの値を作成
	values := make(plotter.Values, len(months))
	for i, m := range months {
		values[i] = float64(monthly[m])
	}
	// グラフの作成
	p := plot.New()
	p.Title.Text = "月別集計"
	p.Y.Label.Text = "合計金額"

	// Y軸の最大値を設定
	maxValue := 0
	for _, v := range monthly {
		if v > maxValue {
			maxValue = v
		}
	}

	p.Y.Max = float64(maxValue) * 1.2

	// Y軸の目盛りを設定
	ticks := []plot.Tick{}
	for v := 0.0; v <= p.Y.Max; v += 100 {
		ticks = append(ticks, plot.Tick{Value: v, Label: fmt.Sprintf("%.0f", v)})
	}

	p.Y.Tick.Marker = plot.ConstantTicks(ticks)

	// 棒グラフの作成
	bar, err := plotter.NewBarChart(values, vg.Points(20))
	if err != nil {
		log.Fatal(err)
	}

	p.NominalX(months...)
	p.Add(bar)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, "../output/monthly.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("月別集計のグラフを output/monthly.png に保存しました。")
}

func createWeeklyGraph(weekly map[string]int) {
	// 曜日の順番を定義
	weekOrder := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	values := make(plotter.Values, len(weekOrder))
	for i, w := range weekOrder {
		values[i] = float64(weekly[w])
	}

	p := plot.New()
	p.Title.Text = "曜日別集計"
	p.Y.Label.Text = "合計金額"

	// Y軸の最大値を設定
	maxValue := 0
	for _, v := range weekly {
		if v > maxValue {
			maxValue = v
		}
	}

	p.Y.Max = float64(maxValue) * 1.2

	// Y軸の目盛りを設定
	ticks := []plot.Tick{}
	for v := 0.0; v <= p.Y.Max; v += 100 {
		ticks = append(ticks, plot.Tick{Value: v, Label: fmt.Sprintf("%.0f", v)})
	}

	p.Y.Tick.Marker = plot.ConstantTicks(ticks)

	bar, err := plotter.NewBarChart(values, vg.Points(20))
	if err != nil {
		log.Fatal(err)
	}

	p.NominalX(weekOrder...)
	p.Add(bar)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, "../output/weekly.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("曜日別集計のグラフを output/weekly.png に保存しました。")
}

func createMonthlyWeeklyCrossGraph(cross map[string]map[string]int) {
	months := make([]string, 0, len(cross))
	for m := range cross {
		months = append(months, m)
	}
	sort.Strings(months)

	weekOrder := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	p := plot.New()
	p.Title.Text = "月×曜日クロス集計"
	p.Y.Label.Text = "合計金額"

	// Y軸の最大値を設定
	maxValue := 0
	for _, weekData := range cross {
		for _, v := range weekData {
			if v > maxValue {
				maxValue = v
			}
		}
	}
	p.Y.Max = float64(maxValue) * 1.2

	// Y軸の目盛りを設定
	ticks := []plot.Tick{}
	for v := 0.0; v <= p.Y.Max; v += 100 {
		ticks = append(ticks, plot.Tick{Value: v, Label: fmt.Sprintf("%.0f", v)})
	}
	p.Y.Tick.Marker = plot.ConstantTicks(ticks)

	barWidth := vg.Points(12)
	for wi, w := range weekOrder {
		vals := make(plotter.Values, len(months))
		for mi, m := range months {
			vals[mi] = float64(cross[m][w])
		}

		bar, err := plotter.NewBarChart(vals, barWidth)
		if err != nil {
			log.Fatal(err)
		}
		bar.Offset = barWidth * vg.Length(wi)
		bar.Color = plotutil.Color(wi)
		p.Add(bar)
	}

	p.NominalX(months...)

	if err := p.Save(12*vg.Inch, 6*vg.Inch, "../output/month_week_cross.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("月×曜日クロス集計のグラフを output/month_week_cross.png に保存しました。")
}

func MonthlyWeeklyRatioGraph(ratioCross map[string]map[string]float64, output string) error {
	months := make([]string, 0, len(ratioCross))
	for m := range ratioCross {
		months = append(months, m)
	}
	sort.Strings(months)
	p := plot.New()
	p.Title.Text = "月×曜日 比率クロス集計"
	p.Y.Label.Text = "比率 (%)"
	p.Y.Min = 0
	p.Y.Max = 50

	// Y軸の目盛りを設定
	ticks := []plot.Tick{}
	for v := 0.0; v <= 50; v += 10 {
		ticks = append(ticks, plot.Tick{Value: v, Label: fmt.Sprintf("%.0f%%", v)})
	}
	p.Y.Tick.Marker = plot.ConstantTicks(ticks)

	// 棒グラフの横幅を設定
	barWidth := vg.Points(20)
	weekOrder := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for wi, w := range weekOrder {
		vals := make(plotter.Values, len(months))
		for mi, m := range months {
			vals[mi] = ratioCross[m][w]
		}
		bar, err := plotter.NewBarChart(vals, barWidth)
		if err != nil {
			return err
		}

		// 曜日ごとに棒を横へずらして、横並びにするため
		bar.Offset = barWidth * vg.Length(wi)
		// 棒グラフの色を設定
		bar.Color = plotutil.Color(wi)
		p.Add(bar)
	}
	p.NominalX(months...)
	return p.Save(12*vg.Inch, 6*vg.Inch, output)
}
