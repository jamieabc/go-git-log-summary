package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/jamieabc/go-git-log-summary/internal/analyzer"

	"github.com/jamieabc/go-git-log-summary/pkg/git_logs"
)

var (
	repoPath string
)

func init() {
	flag.StringVar(&repoPath, "p", "", "repo path")
	flag.Parse()
	if repoPath == "" {
		fmt.Println("Usage: -p repo path (required)")
		os.Exit(1)
	}
	if !path.IsAbs(repoPath) {

	}
}

type stat struct {
	data map[time.Time][]int
}

func main() {
	g := git_logs.NewGitLogs(repoPath)
	logs, err := g.Parse()
	if nil != err {
		panic(err)
	}

	result := stat{
		data: make(map[time.Time][]int),
	}

	for _, l := range logs {
		if count, nums := extractNumber(l.Message); count > 0 {
			if _, ok := result.data[l.Date]; !ok {
				result.data[l.Date] = nums
			} else {
				result.data[l.Date] = append(result.data[l.Date], nums...)
			}
		}
	}

	a := analyzer.NewAnalyzer(result.data)
	analyzed := a.Analyze()

	fmt.Printf("monthly\n")
	for k, v := range analyzed["monthly"] {
		fmt.Printf("\t%s: total %d (%v)\n\n", k, len(v), v)
	}

	fmt.Printf("daily\n")
	for k, v := range analyzed["daily"] {
		fmt.Printf("\t%s: total %d (%v)\n\n", k, len(v), v)
	}
}

func extractNumber(str string) (int, []int) {
	strs := strings.Split(str, " ")
	count := 0
	result := make([]int, 0)

	var i int
	for _, s := range strs {
		// remove trailing :
		s1 := strings.TrimSuffix(s, ":")

		// skip leading #:/
		for i = 0; i < len(s1); i++ {
			if s[i] != '#' && s[i] != ':' && s[i] != '/' {
				break
			}
		}

		if s1[i] >= '0' && s1[i] <= '9' {
			count++
			num, err := strconv.Atoi(s1)
			if nil != err {
				fmt.Printf("convert string %s to interger with error: %s", s, err)
				continue
			}
			result = append(result, num)
		}
	}

	return count, result
}
