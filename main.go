package main

import (
	"flag"
	"log"
	"net/http"

	"task_queue/handlers"
	"task_queue/worker"
)

func main() {
	var numOfWorkers int
	flag.IntVar(&numOfWorkers, "N", 5, "Number of workers")
	flag.Parse()

	http.HandleFunc("/add", handlers.AddTask)
	http.HandleFunc("/list", handlers.ListTask)

	for i := 0; i < numOfWorkers; i++ {
		worker.Start()
	}

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Println("http.ListenAndServe error: ", err)
	}
}
