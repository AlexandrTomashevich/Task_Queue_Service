package queue

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type TaskStatus string

const (
	TaskStatusQueued     TaskStatus = "Queued"
	TaskStatusProcessing TaskStatus = "Processing"
	TaskStatusCompleted  TaskStatus = "Completed"
	TaskStatusFailed     TaskStatus = "Failed"
)

type Task struct {
	Id          int64
	N           int
	D, N1, I    float64
	TTL         time.Duration
	StartTime   time.Time
	EndTime     time.Time
	CreatedTime time.Time
	Status      TaskStatus
}

type Result struct {
	task        *Task
	sum         float64
	expiredTime time.Time
}

// Очередь задач и результатов выполнения.
type TaskQueue struct {
	mut sync.Mutex
	q   []*Task

	results   map[int64]*Result
	resultMut sync.Mutex
}

// Очистка устаревших результатов задач
func ClearExpiredResult(q *TaskQueue) {
	go func() {
		timer := time.NewTicker(time.Second)
		for {
			select {
			case <-timer.C:
				q.resultMut.Lock()
				for _, r := range q.results {
					if time.Now().After(r.expiredTime) {
						delete(q.results, r.task.Id)
					}
				}
				q.resultMut.Unlock()
			}
		}
	}()
}

// Инициализация пустой карты результатов и запуск ClearExpiredResult
func NewTaskQueue() *TaskQueue {
	results := make(map[int64]*Result)

	qq := &TaskQueue{
		results: results,
	}

	ClearExpiredResult(qq)

	return qq
}

// Глобальная переменная
var q = NewTaskQueue()

// Добавление задачи в очередь
func (t *TaskQueue) Add(task *Task) *Task {
	t.mut.Lock()
	t.q = append(t.q, task)
	log.Println("task added", task)
	t.mut.Unlock()
	return task
}

// Извлечение задачи из очереди
func (t *TaskQueue) Fetch() (*Task, bool) {
	t.mut.Lock()
	defer t.mut.Unlock()
	if len(t.q) == 0 {
		return nil, false
	}
	task := t.q[0]
	t.q = t.q[1:]
	log.Println("task fetched", task)
	return task, true
}

var taskID int64

// Генерируем уникальный id
func AddTask(task *Task) *Task {
	task.Id = atomic.AddInt64(&taskID, 1)
	task.Status = TaskStatusQueued
	task.CreatedTime = time.Now()
	return q.Add(task)
}

// Извлекаем задачу из очереди и меняем статус
func FetchTask() (*Task, bool) {
	task, ok := q.Fetch()
	if ok {
		task.StartTime = time.Now()
		task.Status = TaskStatusProcessing
	}
	return task, ok
}

func AddResult(sum float64, task *Task, ttl time.Duration) {
	q.resultMut.Lock()
	task.Status = TaskStatusCompleted
	task.EndTime = time.Now()
	q.results[task.Id] = &Result{
		task:        task,
		sum:         sum,
		expiredTime: time.Now().Add(ttl),
	}
	q.resultMut.Unlock()
}
