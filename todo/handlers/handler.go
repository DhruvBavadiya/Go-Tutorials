package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type TODO struct {
	TaskName    string `json:"task_name"` // JSON key will be task_name
	TaskID      int    `json:"task_id"`   // We'll assign this manually, so it's not expected in JSON
	IsCompleted bool   `json:"is_completed"`
}

var (
	tasks  = make(map[int]TODO)
	nextID = 1
	todoMu sync.Mutex
)

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/add/todo" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	var newTask TODO
	// unmarshal json to byte
	if err := json.Unmarshal(body, &newTask); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	// core logic
	todoMu.Lock()
	newTask.TaskID = nextID
	nextID++
	tasks[newTask.TaskID] = newTask
	todoMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET request received")
	if r.URL.Path != "/get/all" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	todoMu.Lock()
	defer todoMu.Unlock()

	data, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET one request received")

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/get/todo/")
	fmt.Println(idStr)
	// convert string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	todoMu.Lock()
	defer todoMu.Unlock()

	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	res, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Could not get task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func EditByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT one request received")
	if r.Method != "PUT" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/edit/todo/")
	fmt.Println(idStr)
	// convert string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	// get body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	todoMu.Lock()
	defer todoMu.Unlock()
	var Todo TODO

	if err := json.Unmarshal(body, &Todo); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	task.TaskName = Todo.TaskName
	task.IsCompleted = Todo.IsCompleted
	tasks[id] = task

	res, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Could not get task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func DeletebyID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete one request received")
	if r.Method != "DELETE" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/delete/todo/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, exists := tasks[id]
	if !exists {
		http.Error(w, "Record is not present", http.StatusBadRequest)
		return
	}
	todoMu.Lock()
	defer todoMu.Unlock()
	delete(tasks, id)

	msgStr := "record with given id is deleted successfully"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msgStr))
}
