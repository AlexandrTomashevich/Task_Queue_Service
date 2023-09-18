package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"task_queue/queue"
)

var taskStorage []*queue.Task
var mut sync.RWMutex

func AddTask(w http.ResponseWriter, r *http.Request) {
	var task queue.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Println("Error decoding task", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mut.Lock()
	taskAdded := queue.AddTask(&task)
	taskStorage = append(taskStorage, taskAdded)
	mut.Unlock()

	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Println("Error encoding task", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ListTask(w http.ResponseWriter, r *http.Request) {
	mut.RLock()
	defer mut.RUnlock()
	jsonData, err := json.MarshalIndent(taskStorage, "", "    ")
	if err != nil {
		log.Println("Error encoding task list", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
