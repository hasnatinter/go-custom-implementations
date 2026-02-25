package main

import (
	Writersmutex "app/WritersMutex"
	"fmt"
	"math/rand"
	"time"
)

const WriterWait time.Duration = 200 * time.Millisecond
const ReadersProcessingTime time.Duration = 100 * time.Microsecond

type CustomLog struct {
	level   string
	message string
}

func GenerateUsereLogs(logsList *[]CustomLog, wm *Writersmutex.WPMutex) {
	sampleLogs := [5]CustomLog{
		{level: "Info", message: "A new user logged-in"},
		{level: "Warn", message: "User entered incorrect credentials"},
		{level: "Info", message: "User added a new order"},
		{level: "Warn", message: "User entered order with no product"},
		{level: "Error", message: "User could not register order"},
	}
	for {
		wm.WriteLock()
		newLog := sampleLogs[rand.Intn(5)]
		*logsList = append(*logsList, newLog)
		wm.WriteUnlock()
		time.Sleep(WriterWait)
	}
}

func GenerateRequestLogs(logsList *[]CustomLog, wm *Writersmutex.WPMutex) {
	sampleLogs := [3]CustomLog{
		{level: "Info", message: "A new request : req-123"},
		{level: "Warn", message: "Order form validation failed: req-456"},
		{level: "Error", message: "Failed process new order: req-159"},
	}
	for {
		wm.WriteLock()
		newLog := sampleLogs[rand.Intn(3)]
		*logsList = append(*logsList, newLog)
		wm.WriteUnlock()
		time.Sleep(WriterWait)
	}
}

func ReportLogs(logsList *[]CustomLog, wm *Writersmutex.WPMutex) {
	processedCount := 0
	for {
		wm.ReadLock()
		if processedCount != len(*logsList) {
			copiedLogs := ExtractNewLogs(*logsList, processedCount)
			processedCount = len(*logsList)
			for _, log := range copiedLogs {
				fmt.Printf("%d - %s: %s\n", processedCount, log.level, log.message)
			}
		}
		time.Sleep(ReadersProcessingTime)
		wm.ReadUnlock()
	}
}

func ReportErrorLogs(logsList *[]CustomLog, wm *Writersmutex.WPMutex) {
	processedCount := 0
	for {
		wm.ReadLock()
		if processedCount != len(*logsList) {
			copiedLogs := ExtractNewLogs(*logsList, processedCount)
			processedCount = len(*logsList)
			for _, log := range copiedLogs {
				if log.level == "Error" {
					fmt.Printf("%d - ***An error log*** %s: %s\n", processedCount, log.level, log.message)
				}
			}
		}
		time.Sleep(ReadersProcessingTime)
		wm.ReadUnlock()
	}
}

func ExtractNewLogs(cLog []CustomLog, processedCount int) []CustomLog {
	copiedCLog := make([]CustomLog, 0)
	for i := processedCount; i < len(cLog); i++ {
		copiedCLog = append(copiedCLog, cLog[i])
	}
	return copiedCLog
}

func main() {
	logsList := make([]CustomLog, 0)
	readerMutex := Writersmutex.New()
	go ReportLogs(&logsList, readerMutex)
	go ReportErrorLogs(&logsList, readerMutex)
	go GenerateUsereLogs(&logsList, readerMutex)
	GenerateRequestLogs(&logsList, readerMutex)
}
