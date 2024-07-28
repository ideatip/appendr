package appenders

import (
	"fmt"
	"go.ideatip.dev/appendr/models"
	"go.ideatip.dev/appendr/utils"
	"os"
	"sync"
	"time"
)

type ConsoleAppender struct{}

func (c *ConsoleAppender) Append(level models.LogLevel, message string, fields []models.Field) {
	fmt.Printf("[%s] %s %s\n", level, message, utils.FieldsToString(fields))
}

type FileAppender struct {
	mu            sync.Mutex
	file          *os.File
	filename      string
	maxSize       int64 // maximum size of the log file before rotation (in bytes)
	currentSize   int64
	rotationCount int
}

func NewFileAppender(filename string, maxSize int64) (*FileAppender, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &FileAppender{
		file:          file,
		filename:      filename,
		maxSize:       maxSize,
		currentSize:   info.Size(),
		rotationCount: 0,
	}, nil
}

func (fa *FileAppender) Append(level models.LogLevel, message string, fields []models.Field) {
	fa.mu.Lock()
	defer fa.mu.Unlock()

	logEntry := fmt.Sprintf("[%s] %s %s - %s\n",
		time.Now().Format(time.RFC3339),
		level,
		message,
		utils.FieldsToString(fields))

	entrySize := int64(len(logEntry))

	if fa.currentSize+entrySize > fa.maxSize {
		fa.rotate()
	}

	n, err := fa.file.WriteString(logEntry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to log file: %v\n", err)
		return
	}

	fa.currentSize += int64(n)
}

func (fa *FileAppender) rotate() {
	fa.file.Close()

	// Rename the current file
	newName := fmt.Sprintf("%s.%d", fa.filename, fa.rotationCount)
	err := os.Rename(fa.filename, newName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rotating log file: %v\n", err)
		return
	}

	// Open a new file
	file, err := os.OpenFile(fa.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating new log file: %v\n", err)
		return
	}

	fa.file = file
	fa.currentSize = 0
	fa.rotationCount++

	// Delete old log files if there are too many
	fa.cleanOldLogs()
}

func (fa *FileAppender) cleanOldLogs() {
	const maxOldLogs = 5
	if fa.rotationCount <= maxOldLogs {
		return
	}

	oldestLog := fa.rotationCount - maxOldLogs
	oldLogName := fmt.Sprintf("%s.%d", fa.filename, oldestLog)
	err := os.Remove(oldLogName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error removing old log file: %v\n", err)
	}
}

func (fa *FileAppender) Close() error {
	fa.mu.Lock()
	defer fa.mu.Unlock()
	return fa.file.Close()
}
