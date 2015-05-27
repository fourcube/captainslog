package main

import (
	"log"
	"strings"
	"time"
)

type LogEntry struct {
	time.Time
	Lines []string
}

func Parse(raw string) []LogEntry {
	fullLog := strings.Split(raw, "##")
	entries := make([]LogEntry, 0)

	for _, text := range fullLog {
		logEntry := LogEntry{time.Now(), make([]string, 0)}
		validEntry := false

		for i, entry := range strings.Split(text, "\n") {

			if i == 0 {
				entry := strings.TrimSpace(entry)
				if len(entry) < 1 {
					continue
				}

				t, err := time.Parse(TimeLayout, entry)
				if err != nil {
					log.Println("Couldn't parse %s as time, %v", entry, err)
					return entries
				}

				validEntry = true
				logEntry.Time = t
				continue
			}

			logEntry.Lines = append(logEntry.Lines, entry)
		}

		if validEntry {
			entries = append(entries, logEntry)
		}
	}

	return entries
}
