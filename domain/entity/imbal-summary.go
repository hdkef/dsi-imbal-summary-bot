package entity

import (
	"fmt"
	"sort"
	"strings"
)

type MonthlySummaryResult struct {
	Key                string
	total              float64
	DailySummaryResult []DailySummaryResult
}

func (m MonthlySummaryResult) GetTotal() float64 {
	totalPerMonth := 0.0
	for _, v := range m.DailySummaryResult {
		totalPerMonth += v.Total
	}
	return totalPerMonth
}

type DailySummaryResult struct {
	Key    string
	Imbals []Imbal
	Total  float64
}

type ImbalSummaryDto struct {
	token      string
	startMonth *int
	endMonth   *int
	year       *int
}

// Define a slice of Data structs
type DateSlice []struct {
	key   string
	value uint32
}

// Implement the sort.Interface interface for DateSlice
func (s DateSlice) Len() int           { return len(s) }
func (s DateSlice) Less(i, j int) bool { return s[i].value < s[j].value }
func (s DateSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type ImbalSummary struct {
	imbals            map[string][]Imbal //the key is date-month-year
	allAvailableDates DateSlice
	allAvailableMonth DateSlice
	allTotalDates     map[string]float64
}

func (i *ImbalSummaryDto) SetToken(token string) {
	i.token = token
}

func (i *ImbalSummaryDto) GetToken() string {
	return i.token
}

func (i *ImbalSummaryDto) SetStartMonth(startMonth *int) {
	i.startMonth = startMonth
}

func (i *ImbalSummaryDto) GetStartMonth() *int {
	return i.startMonth
}

func (i *ImbalSummaryDto) SetEndMonth(endMonth *int) {
	i.endMonth = endMonth
}

func (i *ImbalSummaryDto) GetEndMonth() *int {
	return i.endMonth
}

func (i *ImbalSummaryDto) SetYear(year *int) {
	i.year = year
}

func (i *ImbalSummaryDto) GetYear() *int {
	return i.year
}

func (i *ImbalSummary) SetImbal(imbal Imbal) {

	key := fmt.Sprintf("%d-%d-%d", imbal.GetDate().Day(), imbal.GetDate().Month(), imbal.GetDate().Year())

	if i.imbals == nil {
		i.imbals = make(map[string][]Imbal)
	}

	// if not exist, append to available dates
	if _, ok := i.imbals[key]; !ok {
		i.allAvailableDates = append(i.allAvailableDates, struct {
			key   string
			value uint32
		}{
			key:   key,
			value: uint32(imbal.GetDate().Unix()),
		})
	}

	i.imbals[key] = append(i.imbals[key], imbal)

	if i.allTotalDates == nil {
		i.allTotalDates = make(map[string]float64)
	}

	i.allTotalDates[key] += imbal.GetAmount()
}

func (i *ImbalSummary) GetImbal() ([]MonthlySummaryResult, float64) {

	// sort available dates asc
	sort.Sort(i.allAvailableDates)

	MonthlyResults := make(map[string]MonthlySummaryResult)
	grandTotal := 0.0

	// append values based on sorted available dates and month
	for _, v := range i.allAvailableDates {
		splitted := strings.Split(v.key, "-")
		month := splitted[1]
		year := splitted[2]

		key := fmt.Sprintf("%s-%s", month, year)

		m, exist := MonthlyResults[key]
		if exist {
			MonthlyResults[key] = MonthlySummaryResult{
				Key: key,
				DailySummaryResult: append(m.DailySummaryResult, DailySummaryResult{
					Key:    v.key,
					Imbals: i.imbals[v.key],
					Total:  i.allTotalDates[v.key],
				}),
			}
		} else {
			MonthlyResults[key] = MonthlySummaryResult{
				Key: key,
				DailySummaryResult: []DailySummaryResult{
					{
						Key:    v.key,
						Imbals: i.imbals[v.key],
						Total:  i.allTotalDates[v.key],
					},
				},
			}
			i.allAvailableMonth = append(i.allAvailableMonth, struct {
				key   string
				value uint32
			}{
				key: key,
			})
		}
	}

	results := []MonthlySummaryResult{}

	for _, v := range i.allAvailableMonth {
		totalPerMonth := MonthlyResults[v.key].GetTotal()
		grandTotal += totalPerMonth
		results = append(results, MonthlySummaryResult{
			Key:                MonthlyResults[v.key].Key,
			total:              totalPerMonth,
			DailySummaryResult: MonthlyResults[v.key].DailySummaryResult,
		})
	}

	return results, grandTotal
}
