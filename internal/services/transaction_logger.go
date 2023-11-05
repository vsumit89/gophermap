package services

import (
	"fmt"
	"os"
)

type TransactionLogger interface {
	WritePut(key, value string) // WritePut writes a key/value pair to the transaction log
	WriteDelete(key string)     // WriteDelete writes a delete operation to the transaction log

}

func NewTransactionLogger(persistentType bool) TransactionLogger {
	if persistentType {
		return &FileTransactionLogger{}
	} else {
		return &DatabaseTransactionLogger{}
	}
}

// EventType is the type of event that is logged in the transaction log
// EventType is an enum with two values: EventPut and EventDelete
type EventType byte

const (
	EventPut EventType = iota + 1
	EventDelete
)

// LogEvent is the event that is logged in the transaction log
// LogEvent is a struct with the following fields:
// SequenceNumber: the sequence number of the event
// EventType: the type of the event
// Key: the key of the event
// Value: the value of the event
type LogEvent struct {
	SequenceNumber int64
	EventType      EventType
	Key            string
	Value          string
}

// file based transaction logger
// the file will be an append only file created in data dir in the root of the project
// the file will be named ds.log
type FileTransactionLogger struct {
	events             chan<- LogEvent // Write-only channel for sending events to the file logger goroutine
	errors             <-chan error    // Read-only channel for receiving errors from the file logger goroutine
	lastSequenceNumber int64
	file               *os.File
}

func (l *FileTransactionLogger) WritePut(key, value string) {
	l.events <- LogEvent{
		EventType: EventPut,
		Key:       key,
		Value:     value,
	}
}

func (l *FileTransactionLogger) WriteDelete(key string) {
	l.events <- LogEvent{
		EventType: EventDelete,
		Key:       key,
	}
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

const (
	MAX_BUFFER_SIZE = 16
)

func (l *FileTransactionLogger) Run() {
	events := make(chan LogEvent, MAX_BUFFER_SIZE)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	// go routine to write events to the file
	go func() {
		for e := range events {
			l.lastSequenceNumber++
			_, err := fmt.Fprintf(l.file, "%d\t%d\t%s\t%s\n", l.lastSequenceNumber, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

// Database based transaction logger
type DatabaseTransactionLogger struct {
}

func (l *DatabaseTransactionLogger) WritePut(key, value string) {

}

func (l *DatabaseTransactionLogger) WriteDelete(key string) {
}
