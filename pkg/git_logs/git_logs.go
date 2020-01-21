package git_logs

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	gitDir   = ".git"
	headsDir = "refs/heads"
	objDir   = "objects"
)

type GitLogs interface {
	Parse() ([]Commit, error)
}

type gitLogs struct {
	path string
}

type Commit struct {
	Author    string
	Committer string
	Date      time.Time
	Message   string
	Parent    string
}

func (g gitLogs) Parse() ([]Commit, error) {
	if _, err := os.Stat(g.path); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s not exist", g.path)
	}

	commits := make([]Commit, 0)

	// assumes it's a valid git dir, exists dir of objects & refs
	head, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s", g.path, headsDir, "master"))
	if nil != err {
		return nil, fmt.Errorf("read parent refs with error: %s", err)
	}

	parent := strings.TrimSuffix(string(head), "\n")

	for parent != "" {
		raw, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s/%s", g.path, objDir, string(parent[0:2]), string(parent[2:])))
		if nil != err {
			return nil, fmt.Errorf("read object with error: %s", err)
		}

		c, err := parseGitObject(raw)
		if nil != err {
			return nil, fmt.Errorf("parse git object with error: %s", err)
		}
		commits = append(commits, c)
		parent = c.Parent
	}

	return commits, nil
}

func parseGitObject(raw []byte) (Commit, error) {
	b := bytes.NewReader(raw)
	r, _ := zlib.NewReader(b)

	c := Commit{}
	var content bytes.Buffer
	w := bufio.NewWriter(&content)

	_, err := io.Copy(w, r)
	if nil != err {
		return c, fmt.Errorf("write to buffer with error: %s", err)
	}

	return parseCommit(content.Bytes()), nil
}

func parseCommit(data []byte) Commit {
	lines := strings.Split(string(data), "\n")

	var parent, author, committer string
	var date time.Time
	var msgStart int
	for idx, l := range lines {
		// parent
		if i := strings.Index(l, "parent"); i != -1 {
			start := strings.Index(l, "parent") + 7
			parent = l[start:]
			continue
		}

		// author
		if i := strings.Index(l, "author"); i != -1 {
			start := strings.Index(l, "author ") + 7
			end := strings.Index(l, " <")
			author = l[start:end]
			continue
		}

		// committer
		if i := strings.Index(l, "committer"); i != -1 {
			start := strings.Index(l, "committer ") + 10
			end := strings.Index(l, " <")
			committer = l[start:end]

			start = strings.Index(l, "> ") + 2
			end = strings.Index(l, " +")

			sec, _ := strconv.Atoi(l[start:end])
			date = time.Unix(int64(sec), 0)

			msgStart = idx + 2
			continue
		}
	}
	return Commit{
		Author:    author,
		Committer: committer,
		Date:      date,
		Message:   lines[msgStart],
		Parent:    parent,
	}
}

func NewGitLogs(path string) GitLogs {
	return &gitLogs{path: fmt.Sprintf("%s/%s", path, gitDir)}
}
