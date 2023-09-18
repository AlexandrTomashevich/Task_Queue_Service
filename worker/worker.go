package worker

import (
	"log"
	"task_queue/queue"
	"time"
)

// Выполнение задач из очереди
func doWork() {
	task, ok := queue.FetchTask()
	if !ok {
		time.Sleep(100 * time.Millisecond)
		return
	}

	sum := task.N1
	delay := time.Duration(int(task.I*1000)) * time.Millisecond
	log.Println("task started", task)

	for i := 0; i < task.N; i++ {
		sum += task.D
		time.Sleep(delay)
	}

	log.Println("task completed successfully", task, sum)

	queue.AddResult(sum, task, task.TTL)
}

func Start() {
	go func() {
		for {
			doWork()
		}
	}()
}
