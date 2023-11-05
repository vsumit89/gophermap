package services

import (
	"bufio"
	"fmt"
	"gophermap/internal/db"
	"os"
)

type TransactionLogger interface {
	Init() error
	WritePut(key, value string) // WritePut writes a key/value pair to the transaction log
	WriteDelete(key string)     // WriteDelete writes a delete operation to the transaction log
	Run()
	Err() <-chan error
	ReadEvents() (<-chan LogEvent, <-chan error)
}

func NewTransactionLogger(persistentType string) TransactionLogger {
	if persistentType == string(file) {
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

// PersistenceType is the type of persistence that is used for the transaction log
// PersistenceType is an enum with two values: file and db
type PersistenceType string

const (
	file     PersistenceType = "logfile"
	database PersistenceType = "database"
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

func (l *FileTransactionLogger) Init() error {
	file, err := os.OpenFile("../data/ds.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	l.file = file
	l.lastSequenceNumber = 0
	l.Run()
	return nil
}

func (l *FileTransactionLogger) WritePut(key, value string) {
	if l.events != nil {
		l.events <- LogEvent{
			EventType: EventPut,
			Key:       key,
			Value:     value,
		}
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

func (l *FileTransactionLogger) ReadEvents() (<-chan LogEvent, <-chan error) {
	scanner := bufio.NewScanner(l.file) // Create a Scanner for l.file
	outEvent := make(chan LogEvent)     // An unbuffered LogEvent channel
	outError := make(chan error, 1)     // A buffered error channel

	go func() {
		var e LogEvent

		defer close(outEvent) // Close the channels when the
		defer close(outError) // goroutine ends

		count := 0
		for scanner.Scan() {
			count++
			line := scanner.Text()
			_, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s", &e.SequenceNumber, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				outError <- fmt.Errorf("input parse error: %w", err)
			}

			// Sanity check! Are the sequence numbers in increasing order?
			if l.lastSequenceNumber >= e.SequenceNumber {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			l.lastSequenceNumber = e.SequenceNumber // Update last used sequence #
			outEvent <- e                           // Send the event along
		}

		fmt.Println("sequqnce number", l.lastSequenceNumber)
		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()

	return outEvent, outError
}

// Database based transaction logger
type DatabaseTransactionLogger struct {
	db db.IDatabase
}

func (l *DatabaseTransactionLogger) Init() error {
	l.db.Connect()
	return nil
}
func (l *DatabaseTransactionLogger) WritePut(key, value string) {

}

func (l *DatabaseTransactionLogger) WriteDelete(key string) {
}

func (l *DatabaseTransactionLogger) Run() {

}

func (l *DatabaseTransactionLogger) Err() <-chan error {
	return nil
}

func (l *DatabaseTransactionLogger) ReadEvents() (<-chan LogEvent, <-chan error) {
	return nil, nil
}
