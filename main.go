package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	metricpkg "main/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Status struct {
	Stat string `json:"status"`
}

func getEnv(key string) string {
	value := os.Getenv(key)
	return value
}

func getLogHandler(metric *metricpkg.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		LOG_FILE := getEnv("LOG_FILE")

		if _, err := os.Stat(LOG_FILE); os.IsNotExist(err) {
			emptyFile, createErr := os.Create(LOG_FILE)

			if createErr != nil {
				http.Error(w, fmt.Sprintf("[ERROR] Could not create log file: %v", createErr), http.StatusInternalServerError)
				return
			}
			emptyFile.Close()
		}

		logs, err := os.ReadFile(LOG_FILE)
		if err != nil {
			http.Error(w, fmt.Sprintf("[ERROR] Read file error %v", err), http.StatusBadRequest)
			return
		}

		_, _ = w.Write(logs)

		duration := time.Since(start).Seconds()
		metric.ReqDuration.Observe(duration)
	}
}

func statusHandler(metric *metricpkg.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		status := Status{
			Stat: "ok",
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(status)

		duration := time.Since(start).Seconds()
		metric.ReqDuration.Observe(duration)
	}
}

func simpleHandler(metric *metricpkg.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		WELCOME_MESSAGE := getEnv("WELCOME_MESSAGE")
		fmt.Fprintf(w, WELCOME_MESSAGE)

		duration := time.Since(start).Seconds()
		metric.ReqDuration.Observe(duration)
	}
}

func postLogHandler(metric *metricpkg.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metric.LogCalls.Inc()
		start := time.Now()
		LOG_FILE := getEnv("LOG_FILE")

		if _, err := os.Stat(LOG_FILE); os.IsNotExist(err) {
			emptyFile, createErr := os.Create(LOG_FILE)

			if createErr != nil {
				metric.FailedLogs.Inc()
				http.Error(w, fmt.Sprintf("[ERROR] Could not create log file: %v", createErr), http.StatusInternalServerError)
				return
			}
			emptyFile.Close()
		}

		logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			metric.FailedLogs.Inc()
			http.Error(w, fmt.Sprintf("[ERROR] Read file error %v", err), http.StatusBadRequest)
			return
		}
		defer logFile.Close()

		log.SetOutput(logFile)
		log.SetFlags(log.Lshortfile | log.LstdFlags)

		body, readErr := io.ReadAll(r.Body)

		if readErr != nil {
			metric.FailedLogs.Inc()
			http.Error(w, fmt.Sprintf("[ERROR] Read file error %v", err), http.StatusBadRequest)
			return
		}

		log.Print(string(body))

		duration := time.Since(start).Seconds()
		metric.ReqDuration.Observe(duration)
		metric.SuccesLogs.Inc()
	}
}

func main() {
	metric, err := metricpkg.NewMetrics()
	if err != nil {
		log.Fatal(fmt.Sprintf("Can not create metrics%v", err))
	}

	fmt.Println("Server start")

	http.HandleFunc("/status", statusHandler(metric))
	http.HandleFunc("/log", postLogHandler(metric))
	http.HandleFunc("/logs", getLogHandler(metric))
	http.HandleFunc("/", simpleHandler(metric))
	http.Handle("/metrics", promhttp.HandlerFor(metric.Registry, promhttp.HandlerOpts{}))

	http.ListenAndServe("0.0.0.0:8081", nil)
}
