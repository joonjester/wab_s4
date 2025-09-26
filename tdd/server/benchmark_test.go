package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"tdd/task"
	"testing"
)

func BenchmarkAddTaskHandler(b *testing.B) {
	tm = task.NewTaskManager()

	payload := task.Task{Title: "BenchTask", Description: "desc", Labels: []string{"UNI"}}
	data, _ := json.Marshal(payload)

	b.ResetTimer() // reset before measurement
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(data))
		w := httptest.NewRecorder()

		addTaskHandler(w, req)

		if w.Code != http.StatusCreated {
			b.Fatalf("unexpected status: got %d", w.Code)
		}
	}
}

func BenchmarkGetTaskHandler(b *testing.B) {
	tm = task.NewTaskManager()
	// Fülle TaskManager einmal vorab
	tm.AddTask("BenchTask", "desc", []string{"UNI"})

	req := httptest.NewRequest(http.MethodGet, "/get", nil)

	b.ResetTimer()
	for b.Loop() {
		w := httptest.NewRecorder()
		getTaskHandler(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("unexpected status: got %d", w.Code)
		}
	}
}

func BenchmarkUpdateTaskHandler(b *testing.B) {
	tm = task.NewTaskManager()
	t := tm.AddTask("BenchTask", "desc", []string{"UNI"})

	payload := task.Task{ID: t.ID, Status: "done"}
	data, _ := json.Marshal(payload)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewReader(data))
		w := httptest.NewRecorder()

		updateTaskHandler(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("unexpected status: got %d", w.Code)
		}
	}
}

func BenchmarkDeleteTaskHandler(b *testing.B) {
	tm = task.NewTaskManager()
	t := tm.AddTask("BenchTask", "desc", []string{"UNI"})

	url := "/delete?id=" + strconv.Itoa(t.ID)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		w := httptest.NewRecorder()

		deleteTaskHandler(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("unexpected status: got %d", w.Code)
		}
	}
}

func BenchmarkFullApplication(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/add", addTaskHandler)
	mux.HandleFunc("/get", getTaskHandler)
	mux.HandleFunc("/delete", deleteTaskHandler)
	mux.HandleFunc("/update", updateTaskHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	todo := task.Task{Title: "Test", Description: "Test the test", Labels: []string{"Uni"}}
	todoInProgress := task.Task{ID: 1, Status: task.InProgress}
	todoData, _ := json.Marshal(todo)
	todoJsonInProgress, _ := json.Marshal(todoInProgress)

	b.ResetTimer()
	for b.Loop() {
		resp, _ := http.Post(server.URL+"/add", "application/json", bytes.NewBuffer(todoData))
		resp.Body.Close()

		req, _ := http.NewRequest(http.MethodPut, server.URL+"/update", bytes.NewBuffer(todoJsonInProgress))
		req.Header.Set("Content-Type", "application/json")
		resp, _ = http.DefaultClient.Do(req)
		resp.Body.Close()

		resp, _ = http.Get(server.URL + "/get")
		resp.Body.Close()

		req, _ = http.NewRequest(http.MethodDelete, server.URL+"/delete?id=1", nil)
		resp, _ = http.DefaultClient.Do(req)
		resp.Body.Close()

	}
}
