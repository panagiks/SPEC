package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/panagiks/SPEC/storage"
)

func captureOutput(f func()) (got string) {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	outC := make(chan string)
	go func() {
		var buf strings.Builder
		io.Copy(&buf, r)
		r.Close()
		outC <- buf.String()
	}()
	defer func() {
		// Close pipe, restore stdout, get output.
		w.Close()
		os.Stdout = stdout
		out := <-outC
		got = strings.TrimSpace(out)
	}()
	f()
	return
}

func TestExecutorAddProject(t *testing.T) {
	// Test missing arguments
	expected := "Usage: \"add project {name}\"!"
	dataStore = storage.NewStorage()
	output := captureOutput(func() {
		executorAddProject([]string{})
	})
	if output != expected {
		t.Errorf("\nExpexted: \t'%s'\nGot: \t\t'%s'\n", expected, output)
	}
	// Test successful project creation
	expected = "ProjectName"
	output = captureOutput(func() {
		executorAddProject([]string{"ProjectName"})
	})
	found := false
	for k := range dataStore {
		if k == expected {
			found = true
		}
	}
	if output != "" {
		t.Errorf("\nExpexted: no output\nGot: \t\t'%s'\n", output)
	}
	if !found {
		t.Errorf("\nExpexted: \t'%s' project not found\n", expected)
	}
}
