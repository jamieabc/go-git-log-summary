package analyzer

import (
	"time"
)

type Analyzer interface {
	Analyze() map[string]map[string][]int
}

type analyzer struct {
	data map[time.Time][]int
}

func (a analyzer) Analyze() map[string]map[string][]int {
	keys := make([]time.Time, 0)
	for t := range a.data {
		keys = append(keys, t)
	}

	result := make(map[string]map[string][]int)
	result["monthly"] = monthly(a, keys)
	result["daily"] = daily(a, keys)

	return result
}

func monthly(a analyzer, keys []time.Time) map[string][]int {
	result := make(map[string][]int)
	for _, t := range keys {
		str := t.Format("2006-Jan")
		if nums, ok := result[str]; ok {
			result[str] = append(nums, a.data[t]...)
		} else {
			result[str] = a.data[t]
		}
	}
	return result
}

func daily(a analyzer, keys []time.Time) map[string][]int {
	result := make(map[string][]int)
	for _, t := range keys {
		str := t.Format("2006-Jan-2")
		if nums, ok := result[str]; ok {
			result[str] = append(nums, a.data[t]...)
		} else {
			result[str] = a.data[t]
		}
	}
	return result
}
func NewAnalyzer(data map[time.Time][]int) Analyzer {
	return &analyzer{data: data}
}
