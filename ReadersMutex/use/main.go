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

func ReportLogs(cLog *[]CustomLog, rm *ReadersMutex.RMutex) {
	total := 0
	for {
		rm.ReadLock()
		if total != len(*cLog) {
			copiedLogs := ExtractLogs(*cLog, total)
			for _, log := range copiedLogs {
				fmt.Printf("%d - %s: %s\n", total, log.level, log.message)
			}
			total = len(*cLog)
		}
		rm.ReadUnlock()
		time.Sleep(ReadersWait)
	}
}

func ReportErrorLogs(cLog *[]CustomLog, rm *ReadersMutex.RMutex) {
	total := 0
	for {
		rm.ReadLock()
		if total != len(*cLog) {
			copiedLogs := ExtractLogs(*cLog, total)
			for _, log := range copiedLogs {
				if log.level == "Error" {
					fmt.Printf("%d - ***An error log*** %s: %s\n", total, log.level, log.message)
				}
			}
			total = len(*cLog)
		}
		rm.ReadUnlock()
		time.Sleep(ReadersWait)
	}
}

func ExtractLogs(cLog []CustomLog, total int) []CustomLog {
	copiedCLog := make([]CustomLog, 0)
	for i := total; i < len(cLog); i++ {
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
