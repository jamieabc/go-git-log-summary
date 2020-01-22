package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

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
	data map[string][]int
}

func (s stat) String() string {
	var sb strings.Builder

	keys := make([]string, 0)
	for v := range s.data {
		keys = append(keys, v)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s\ttotal %d\t(%v)\n", k, len(s.data[k]), s.data[k]))
	}

	return sb.String()
}

func main() {
	g := git_logs.NewGitLogs(repoPath)
	logs, err := g.Parse()
	if nil != err {
		panic(err)
	}

	result := stat{
		data: make(map[string][]int),
	}

	for _, l := range logs {
		if count, nums := extractNumber(l.Message); count > 0 {
			t := l.Date.Format("2006-Jan-2")
			if _, ok := result.data[t]; !ok {
				result.data[t] = nums
			} else {
				result.data[t] = append(result.data[t], nums...)
			}
		}
	}

	fmt.Println(result)
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
