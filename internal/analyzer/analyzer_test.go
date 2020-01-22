package analyzer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jamieabc/go-git-log-summary/internal/analyzer"
)

func TestAnalyzer_Analyze(t *testing.T) {
	data := make(map[time.Time][]int)
	now := time.Now()

	yesterday := now.AddDate(0, 0, -1)
	nextMonth := now.AddDate(0, 1, 0)

	data[now] = []int{1, 2}
	data[yesterday] = []int{3, 4}
	data[nextMonth] = []int{7, 8, 9}

	a := analyzer.NewAnalyzer(data)
	monthlyFormat := "2006-Jan"
	monthly := map[string][]int{
		now.Format(monthlyFormat):       []int{1, 2, 3, 4},
		nextMonth.Format(monthlyFormat): []int{7, 8, 9},
	}

	dailyFormat := "2006-Jan-2"
	daily := map[string][]int{
		now.Format(dailyFormat):       []int{1, 2},
		yesterday.Format(dailyFormat): []int{3, 4},
		nextMonth.Format(dailyFormat): []int{7, 8, 9},
	}

	expected := map[string]map[string][]int{
		"monthly": monthly,
		"daily":   daily,
	}

	assert.Equal(t, expected, a.Analyze(), "wrong result")
}
