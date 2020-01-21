package git_logs_test

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/jamieabc/go-git-log-summary/pkg/git_logs"
	"github.com/stretchr/testify/assert"
)

const (
	fixtureGitDir = "./fixtures/.git"
	fixtureDir    = "./fixtures/git"
)

func TestGitLogs_ParseWhenDirNotExist(t *testing.T) {
	g := git_logs.NewGitLogs("not_exist_path")
	_, err := g.Parse()
	assert.NotNil(t, err, "parse not existing dir w/o error")
}

func setupFixtures() {
	if _, err := os.Stat(fixtureGitDir); os.IsExist(err) {
		exec.Command("rm", "-rf", fixtureGitDir)
	}
	_ = exec.Command("cp", "-r", fixtureDir, fixtureGitDir).Run()
}

func teardownFixtures() {
	if _, err := os.Stat(fixtureGitDir); os.IsExist(err) {
		_ = exec.Command("rm", "-rf", fixtureGitDir).Run()
	}
}

func TestGitLogs_Parse(t *testing.T) {
	setupFixtures()
	defer teardownFixtures()

	g := git_logs.NewGitLogs("./fixtures")
	commits, err := g.Parse()
	assert.Nil(t, err, "parse with error")
	assert.Equal(t, 2, len(commits), "wrong number of commits")

	c := commits[0]
	commitTime, _ := time.Parse("Mon Jan 2 15:04:05 2006 -0700", "Tue Jan 21 14:08:51 2020 +0800")

	assert.Equal(t, "second test file", c.Message, "wrong commit message")
	assert.Equal(t, "Aaron Chang", c.Author, "wrong author")
	assert.Equal(t, "Aaron Chang", c.Committer, "wrong committer")
	assert.Equal(t, commitTime, c.Date, "wrong commit date")
	assert.Equal(t, "35bbb173992d2078f8869b508a910f2b1ee1bd20", c.Parent, "wrong parent")

	c = commits[1]
	commitTime, _ = time.Parse("Mon Jan 2 15:04:05 2006 -0700", "Mon Jan 20 14:26:50 2020 +0800")

	assert.Equal(t, "this is a test commit", c.Message, "wrong commit message")
	assert.Equal(t, "Aaron Chang", c.Author, "wrong author")
	assert.Equal(t, "Aaron Chang", c.Committer, "wrong committer")
	assert.Equal(t, commitTime, c.Date, "wrong commit date")
	assert.Equal(t, "", c.Parent, "wrong parent")
}
