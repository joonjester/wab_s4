package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"tdd/task"
	"testing"
)

func setup() {
	tm = task.NewTaskManager()
}

func TestAddTaskHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
	}{
		{"Success", http.MethodPost, `{"Title":"Test","Description":"desc","Labels":["UNI"]}`, http.StatusCreated},
		{"Wrong Method", http.MethodGet, ``, http.StatusMethodNotAllowed},
		{"Invalid JSON", http.MethodPost, `{bad-json`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			req := httptest.NewRequest(tt.method, "/add", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			addTaskHandler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestGetTaskHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{"Success", http.MethodGet, http.StatusOK},
		{"Wrong Method", http.MethodPost, http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			tm.AddTask("Test", "desc", []string{"UNI"})
			req := httptest.NewRequest(tt.method, "/get", nil)
			w := httptest.NewRecorder()

			getTaskHandler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestUpdateTaskHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       any
		wantStatus int
	}{
		{"Success", http.MethodPut, task.Task{ID: 1, Status: "done"}, http.StatusOK},
		{"Wrong Method", http.MethodGet, nil, http.StatusMethodNotAllowed},
		{"Invalid JSON", http.MethodPut, "{bad-json", http.StatusBadRequest},
		{"Not Found", http.MethodPut, task.Task{ID: 999, Status: "done"}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			tm.AddTask("Test", "desc", []string{"UNI"})

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			case task.Task:
				bodyBytes, _ = json.Marshal(v)
			case nil:
				bodyBytes = nil
			}

			req := httptest.NewRequest(tt.method, "/update", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			updateTaskHandler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestDeleteTaskHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		id         string
		wantStatus int
	}{
		{"Success", http.MethodDelete, "1", http.StatusOK},
		{"Wrong Method", http.MethodGet, "1", http.StatusMethodNotAllowed},
		{"Invalid ID", http.MethodDelete, "abc", http.StatusBadRequest},
		{"Not Found", http.MethodDelete, "999", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			tm.AddTask("Test", "desc", []string{"UNI"})

			req := httptest.NewRequest(tt.method, "/delete?id="+tt.id, nil)
			w := httptest.NewRecorder()

			deleteTaskHandler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
