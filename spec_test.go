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

func TestAddProject(t *testing.T) {
	// Test successful project creation
	dataStore = storage.NewStorage()
	expected := "ProjectName"
	output := captureOutput(func() {
		addProject("ProjectName")
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

func TestSelectProject(t *testing.T) {
	dataStore = storage.NewStorage()
	// Test invalid project name
	output := captureOutput(func() {
		selectProject("ProjectName")
	})
	expected := "No project with name: \"ProjectName\"! Please create one with \"add\" command!"
	if output != expected {
		t.Errorf("\nExpexted: \t'%s'\nGot: \t\t'%s'\n", expected, output)
	}
	expected = "ProjectName"
	dataStore["ProjectName"] = storage.NewProject()
	// dataStore["ProjectName"]["TaskName"] = storage.NewTask("TaskName", 2, 0.5)
	output = captureOutput(func() {
		selectProject("ProjectName")
	})
	found := false
	for k := range cli {
		if cli[k].project == expected {
			found = true
		}
	}
	if output != "" {
		t.Errorf("\nExpexted: no output\nGot: \t\t'%s'\n", output)
	}
	if !found {
		t.Errorf("\nExpexted: \t'%s' task not found\n", expected)
	}
}

func TestAddTask(t *testing.T) {
	dataStore = storage.NewStorage()
	// Test invalid project name
	output := captureOutput(func() {
		addTask("ProjectName", "TaskName", "1.0", "0.5")
	})
	expected := "No project with name: \"ProjectName\"! Please create one with \"add\" command!"
	if output != expected {
		t.Errorf("\nExpexted: \t'%s'\nGot: \t\t'%s'\n", expected, output)
	}
	// Test task addition
	dataStore["ProjectName"] = storage.NewProject()
	expected = "TaskName"
	output = captureOutput(func() {
		addTask("ProjectName", "TaskName", "1.0", "0.5")
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

func TestUpdateTask(t *testing.T) {
	dataStore = storage.NewStorage()
	// Test invalid project name
	output := captureOutput(func() {
		updateTask("ProjectName", "TaskName", "1.0")
	})
	expected := "No project with name: \"ProjectName\"! Please create one!"
	if output != expected {
		t.Errorf("\nExpexted: \t'%s'\nGot: \t\t'%s'\n", expected, output)
	}
	// Test invalid task name
	dataStore["ProjectName"] = storage.NewProject()
	output = captureOutput(func() {
		updateTask("ProjectName", "TaskName", "1.0")
	})
	expected = "No task with name: \"TaskName\"! Please create one with \"add\" command!"
	if output != expected {
		t.Errorf("\nExpexted: \t'%s'\nGot: \t\t'%s'\n", expected, output)
	}
}

func TestListTask(t *testing.T) {
	dataStore = storage.NewStorage()
	// Test invalid project name
	output := captureOutput(func() {
		listTask("ProjectName")
	})
	expected := "No project with name: \"ProjectName\"! Please create one!"
	if output != expected {
		t.Errorf("\nExpexted: \t'%s'\nGot: \t\t'%s'\n", expected, output)
	}
}

func TestSelectTask(t *testing.T) {
	dataStore = storage.NewStorage()
	// Test invalid project name
	output := captureOutput(func() {
		selectTask("ProjectName", "TaskName")
	})
	expected := "No task with name: \"TaskName\"! Please create one with \"add\" command!"
	if output != expected {
		t.Errorf("\nExpexted: \t'%s'\nGot: \t\t'%s'\n", expected, output)
	}
	expected = "TaskName"
	dataStore["ProjectName"] = storage.NewProject()
	dataStore["ProjectName"]["TaskName"] = storage.NewTask("TaskName", 2, 0.5)
	output = captureOutput(func() {
		selectTask("ProjectName", "TaskName")
	})
	found := false
	for k := range cli {
		if cli[k].task == expected {
			found = true
		}
	}
	if output != "" {
		t.Errorf("\nExpexted: no output\nGot: \t\t'%s'\n", output)
	}
	if !found {
		t.Errorf("\nExpexted: \t'%s' task not found\n", expected)
	}
}
