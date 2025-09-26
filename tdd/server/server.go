package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tdd/task"
)

var tm *task.TaskManager

func isTaskEmpty(t task.Task) bool {
	return t.ID == 0 &&
		t.Title == "" &&
		t.Description == "" &&
		t.Status == "" &&
		len(t.Labels) == 0
}

func addTaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var t task.Task
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	newTask := tm.AddTask(t.Title, t.Description, t.Labels)
	if isTaskEmpty(newTask) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func updateTaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var t task.Task
	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	if err := tm.UpdateStatus(t.ID, t.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteTaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	idStr := req.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "the id is not a number", http.StatusBadRequest)
		return
	}

	if err := tm.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tm.GetTasks())
}

func Rounter() {
	tm = task.NewTaskManager()

	http.HandleFunc("/add", addTaskHandler)
	http.HandleFunc("/get", getTaskHandler)
	http.HandleFunc("/update", updateTaskHandler)
	http.HandleFunc("/delete", deleteTaskHandler)

	http.ListenAndServe(":8080", nil)
}
