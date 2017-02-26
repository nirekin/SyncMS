package main

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"time"
)

var (
	TraceActivity *log.Logger
	TraceTime     *log.Logger
	TraceError    *log.Logger
	TraceLoaded   *log.Logger
)

func initLogger() {
	l := &lumberjack.Logger{
		Filename:   "MassShooting.log",
		MaxSize:    50, // size in mb
		MaxBackups: 5,
		MaxAge:     60, //  age in days
	}

	TraceLoaded = log.New(l, "Load: ", log.Ldate|log.Ltime|log.Lshortfile)
	TraceError = log.New(l, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
	TraceTime = log.New(l, "Time: ", log.Ldate|log.Ltime|log.Lshortfile)
	TraceActivity = log.New(l, "Activity: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func traceTime(label string) func() {
	t := time.Now()
	return func() {
		TraceTime.Printf("execution of %s took %s \n", label, (time.Since(t)))
	}
}
