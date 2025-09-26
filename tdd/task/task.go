package task

import "errors"

type Status string

const (
	Open       Status = "open"
	InProgress Status = "in progress"
	Done       Status = "done"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      Status
	Labels      []string
}

type TaskManager struct {
	tasks  []Task
	nextID int
}

func NewTaskManager() *TaskManager {
	return &TaskManager{tasks: []Task{}, nextID: 1}
}

func (tm *TaskManager) AddTask(title, description string, lables []string) Task {
	task := Task{
		ID:          tm.nextID,
		Title:       title,
		Description: description,
		Status:      Open,
		Labels:      lables,
	}

	tm.tasks = append(tm.tasks, task)
	tm.nextID++
	return task
}

func (tm *TaskManager) UpdateStatus(id int, status Status) error {
	for i := range tm.tasks {
		if tm.tasks[i].ID == id {
			tm.tasks[i].Status = status
			return nil
		}
	}

	return errors.New("task not found")
}

func (tm *TaskManager) DeleteTask(id int) error {
	for i := range tm.tasks {
		if tm.tasks[i].ID == id {
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}

func (tm *TaskManager) GetTasks() []Task {
	return tm.tasks
}
