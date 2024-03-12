package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/shasderias/whalebrew/packages"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	defer os.RemoveAll(".whalebrew-tests")
	os.Setenv("RUN_WHALEBREW", "true")
	os.Setenv("WHALEBREW_INSTALL_PATH", ".whalebrew-tests")
	assert.NoError(t, os.MkdirAll(".whalebrew-tests", 0777))
	if err := exec.Command("docker", "pull", "alpine").Run(); err != nil {
		t.Skipf("Unable to pull alpine image: %v, this test needs a running docker daemon", err)
		return
	}
	c := exec.Command(os.Args[0], "install", "-y", "-f", "whalebrew/awscli")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	assert.NoError(t, c.Run())
	assert.NoError(t, packages.NewPackageManager(".whalebrew-tests").ForceInstall(&packages.Package{
		Image:      "alpine",
		Name:       "alpine",
		Entrypoint: []string{"sh", "-c"},
	}))

	wd, err := os.Getwd()
	assert.NoError(t, err)

	// running in whalebrew performs an exec, we need to fork to be able
	// to execute the rest of the tests
	assert.NoError(t, exec.Command("docker", "pull", "alpine").Run())
	c = exec.Command(os.Args[0], wd+"/.whalebrew-tests/alpine", "pwd;ls -al .")
	stdout := bytes.NewBuffer(nil)
	c.Stdout = stdout
	c.Stderr = os.Stderr
	assert.NoError(t, c.Run())
	assert.Contains(t, stdout.String(), "main.go")
	assert.True(t, strings.HasPrefix(stdout.String(), "/workdir"))
}

func TestMain(m *testing.M) {
	if os.Getenv("RUN_WHALEBREW") == "true" {
		main()
		os.Exit(0)
	}
	os.Exit(m.Run())
}
