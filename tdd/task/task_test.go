package task

import (
	"testing"
)

func TestAddTask(t *testing.T) {
	tm := NewTaskManager()

	tests := []struct {
		title       string
		description string
		labels      []string
		wantTitle   string
	}{
		{"Task 1", "Beschreibung 1", []string{"Uni"}, "Task 1"},
		{"Task 2", "Beschreibung 2", []string{"Privat", "Arbeit"}, "Task 2"},
	}

	for _, tt := range tests {
		task := tm.AddTask(tt.title, tt.description, tt.labels)
		if task.Title != tt.wantTitle {
			t.Errorf("got %v, want %v", task.Title, tt.wantTitle)
		}
	}
}

func TestUpdateStatus(t *testing.T) {
	tm := NewTaskManager()

	tm.AddTask("Test", "Beschreibung", []string{"Uni"})
	tm.AddTask("Test2", "Beschreibung", []string{"Uni"})
	tm.AddTask("Test3", "Beschreibung", []string{"Uni"})

	tests := map[string]struct {
		id      int
		status  Status
		want    Status
		wantErr bool
	}{
		"Happy Path: Done": {
			id:      1,
			status:  Done,
			want:    Done,
			wantErr: false,
		},
		"Happy Path: In Progress": {
			id:      1,
			status:  InProgress,
			want:    InProgress,
			wantErr: false,
		},
		"Unhappy Path: task not found": {
			id:      5,
			status:  Done,
			want:    Done,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			err := tm.UpdateStatus(tt.id, tt.status)

			if tt.wantErr {
				if err == nil {
					t.Error("wanted error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			got := tm.GetTasks()[tt.id-1].Status
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := map[string]struct {
		id      int
		wantErr bool
	}{
		"Happy Path: could find the task": {
			id:      1,
			wantErr: false,
		},
		"Unhappy Path: couldn't find the task": {
			id:      2,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tm := NewTaskManager()
			_ = tm.AddTask("Test", "Beschreibung", []string{"Uni"})

			err := tm.DeleteTask(tt.id)
			if tt.wantErr {
				if err == nil {
					t.Error("wanted error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(tm.GetTasks()) != 0 {
				t.Errorf("expected no tasks, got %d", len(tm.GetTasks()))
				return
			}
		})
	}
}
