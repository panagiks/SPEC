package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/panagiks/SPEC/storage"
)

// Adapted to the needs of this project from:
// https://golang.org/src/testing/example.go
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

func TestExecutorAdd(t *testing.T) {
	// Test missing arguments
	expected := "Must select what to add; Options {project, task}!"
	output := captureOutput(func() {
		executorAdd([]string{})
	})
	if output != expected {
		t.Errorf("\nExpexted: \t'%s'\nGot: \t\t'%s'\n", expected, output)
	}
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

func TestExecutorAddTask(t *testing.T) {
	// Test missing arguments
	expected := "Usage: \"add task {project} {name} {estimation} {confidence}\"!"
	dataStore = storage.NewStorage()
	output := captureOutput(func() {
		executorAddTask([]string{})
	})
	if output != expected {
		t.Errorf("\nExpexted: \t'%s'\nGot: \t\t'%s'\n", expected, output)
	}
	// Test invalid project name
	output = captureOutput(func() {
		executorAddTask([]string{"ProjectName", "TaskName", "1.0", "0.5"})
	})
	expected = "No project with name: \"ProjectName\"! Please create one with \"add project\" command!"
	if output != expected {
		t.Errorf("\nExpexted: \t'%s'\nGot: \t\t'%s'\n", expected, output)
	}
	// Test task addition
	dataStore["ProjectName"] = storage.NewProject()
	expected = "TaskName"
	output = captureOutput(func() {
		executorAddTask([]string{"ProjectName", "TaskName", "1.0", "0.5"})
	})
	found := false
	for k := range dataStore["ProjectName"] {
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
