package main

import (
	"fmt"
	"my-pingbot/workerpool"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	INTERVAL        = time.Second * 10
	REQUEST_TIMEOUT = time.Second * 2
	WORKERS_COUNT   = 3
)

var urls = []string{
    "https://go.dev/",
    "https://go.dev/learn/",
    "https://ru.wikipedia.org/wiki/Go",
    "https://habr.com/ru/companies/yandex_praktikum/articles/754682/",
    "https://tproger.ru/translations/golang-basics",
}

func main() {
    results := make(chan workerpool.Result)
    workerPool := workerpool.New(WORKERS_COUNT, REQUEST_TIMEOUT, results)

    go generateJobs(workerPool)
    go processResults(results)

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

    <-quit

    workerPool.Stop()
}

func processResults(results chan workerpool.Result) {
    go func() {
        for result := range results {
            fmt.Println(result.Info())
        }
    }()
}

func generateJobs(wp *workerpool.Pool) {
    for {
        for _, url := range urls {
            wp.Push(workerpool.Job{URL: url})
        }

        time.Sleep(INTERVAL)
    }
}