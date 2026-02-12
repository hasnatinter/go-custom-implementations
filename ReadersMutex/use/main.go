package main

import (
	"app/ReadersMutex"
	"fmt"
	"math/rand"
	"time"
)

const WriterWait time.Duration = 200 * time.Millisecond
const ReadersWait time.Duration = 10 * time.Microsecond

type CustomLog struct {
	level   string
	message string
}

func GenerateLogs(logsList *[]CustomLog, rm *ReadersMutex.RMutex) {
	sampleLogs := [5]CustomLog{
		{level: "Info", message: "An info message"},
		{level: "Warn", message: "An warn message"},
		{level: "Info", message: "Another info message"},
		{level: "Warn", message: "Another warn message"},
		{level: "Error", message: "An Error message"},
	}
	for {
		rm.WriteLock()
		newLog := sampleLogs[rand.Intn(5)]
		*logsList = append(*logsList, newLog)
		rm.WriteUnlock()
		time.Sleep(WriterWait)
	}
}

func ReportLogs(logsList *[]CustomLog, rm *ReadersMutex.RMutex) {
	processedCount := 0
	for {
		rm.ReadLock()
		if processedCount != len(*logsList) {
			copiedLogs := ExtractNewLogs(*logsList, processedCount)
			for _, log := range copiedLogs {
				fmt.Printf("%d - %s: %s\n", processedCount, log.level, log.message)
			}
			processedCount = len(*logsList)
		}
		rm.ReadUnlock()
		time.Sleep(ReadersWait)
	}
}

func ReportErrorLogs(logsList *[]CustomLog, rm *ReadersMutex.RMutex) {
	processedCount := 0
	for {
		rm.ReadLock()
		if processedCount != len(*logsList) {
			copiedLogs := ExtractNewLogs(*logsList, processedCount)
			for _, log := range copiedLogs {
				if log.level == "Error" {
					fmt.Printf("%d - ***An error log*** %s: %s\n", processedCount, log.level, log.message)
				}
			}
			processedCount = len(*logsList)
		}
		rm.ReadUnlock()
		time.Sleep(ReadersWait)
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
	readerMutex := ReadersMutex.New()
	go ReportLogs(&logsList, readerMutex)
	go ReportErrorLogs(&logsList, readerMutex)
	GenerateLogs(&logsList, readerMutex)
}
