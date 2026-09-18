package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

// this part is for log level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l LogLevel) GetLogType() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

// this part is for feild and log record

// feild is structured key-value pair attached to a log line
type Feild struct {
	Key   string
	Value interface{}
}

// logrecord is immuatable data passed from loggger -> formatter -> writer
type LogRecord struct {
	TimeStamp time.Time
	Level     LogLevel
	Message   string
	Feilds    []Feild
}

// formatter (strategy design pattern)

// fromatter turns a LogRecord into bytes ready to be written
type Formatter interface {
	Format(rec LogRecord) ([]byte, error)
}

// now solid implementations of the interface

// JSONFormatter render the record as a single-line JSON object
type JSONFormatter struct{}

func (JSONFormatter) Format(rec LogRecord) ([]byte, error) {
	out := map[string]interface{}{
		"ts":    rec.TimeStamp.Format(time.RFC3339),
		"level": rec.Level,
		"msg":   rec.Message,
	}
	for _, f := range rec.Feilds {
		out[f.Key] = f.Value
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

// textFormatter renders "timstamp level messahge key=value..."
type TextFormatter struct{}

func (TextFormatter) Format(rec LogRecord) ([]byte, error) {
	line := fmt.Sprintf("%s [%s] %s", rec.TimeStamp.Format(time.RFC3339), rec.Level, rec.Message)
	for _, f := range rec.Feilds {
		line += fmt.Sprintf(" %s = %v", f.Key, f.Value)
	}
	return []byte(line + "\n"), nil
}

// ............ writer/sink (strategy) ..............

type Writer interface {
	Write(p []byte) (n int, err error)
}

type ConsoleWriter struct{}

func (ConsoleWriter) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

type FileWriter struct {
	mu   sync.Mutex
	file *os.File
}

func NewFileWriter(path string) (*FileWriter, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &FileWriter{file: f}, nil
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.file.Write(p)
}

// rotatingfile decorates the filewriter with size based rotation
// It satisfies the same Writer interface, so it's a drop-in replacement
// for FileWriter anywhere one is used — no change needed to AsyncLogger.
type RotatingFileWriter struct {
	mu      sync.Mutex
	path    string
	maxByte int64
	written int64
	current *FileWriter
}

func NewRotatingWriter(path string, maxByte int64) (*RotatingFileWriter, error) {
	fw, err := NewFileWriter(path)
	if err != nil {
		return nil, err
	}
	return &RotatingFileWriter{path: path, maxByte: maxByte, current: fw}, nil
}

func (w *RotatingFileWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.written+int64(len(p)) > w.maxByte {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}
	n, err := w.current.Write(p)
	w.written += int64(n)
	return n, err
}

func (w *RotatingFileWriter) rotate() error {
	rotatedName := fmt.Sprintf("%s.%d", w.path, time.Now().Unix())
	if err := os.Rename(w.path, rotatedName); err != nil {
		return err
	}
	fw, err := NewFileWriter((w.path))
	if err != nil {
		return err
	}
	w.current = fw
	w.written = 0
	return nil
}

// ---------- Sampler (Strategy — optional, controls log volume) ----------
// Sampler decides whether a given record should actually be emitted.
// A nil Sampler on AsyncLogger means "log everything that passes the
// level filter" — sampling is strictly opt-in.
type Sampler interface {
	ShouldLog(rec LogRecord) bool
}

type RateSampler struct {
	rate float64
}

func NewRateSampler(r float64) *RateSampler {
	return &RateSampler{rate: r}
}

func (s *RateSampler) ShouldLog(rec LogRecord) bool {
	if rec.Level >= ERROR {
		return true
	}
	return rand.Float64() < s.rate
}

// ---------- Logger (interface the rest of the app depends on) ----------

type Logger interface {
	Debug(msg string, feilds ...Feild)
	Info(msg string, feilds ...Feild)
	Warn(msg string, feilds ...Feild)
	Error(msg string, feilds ...Feild)
	With(feilds ...Feild) Logger
	Close()
}

// ---------- AsyncLogger (concrete implementation) ----------
// asyncLogger buffers records on a channel and writes them from a single
// background goroutine, so the calling goroutine never blocks on slow I/O.

type asyncLogger struct {
	minLevel   LogLevel
	formatter  Formatter
	writers    []Writer
	sampler    Sampler
	baseFeilds []Feild

	recordCh  chan LogRecord
	wg        sync.WaitGroup
	closeOnce sync.Once
}

func (l *asyncLogger) log(level LogLevel, msg string, feilds ...Feild) {
	if level < l.minLevel {
		return
	}

	rec := LogRecord{
		TimeStamp: time.Now(),
		Level:     level,
		Message:   msg,
		Feilds:    append(append([]Feild{}, l.baseFeilds...), feilds...),
	}

	if l.sampler != nil && !l.sampler.ShouldLog(rec) {
		return
	}

	select {
	case l.recordCh <- rec:
	default:
		// Channel full: drop the record instead of blocking the caller.
		// A production system might increment a "dropped_logs" metric here.
		fmt.Fprintln(os.Stderr, "logger: buffer full")
	}
}

func (l *asyncLogger) Debug(msg string, feilds ...Feild) { l.log(DEBUG, msg, feilds...) }
func (l *asyncLogger) Info(msg string, feilds ...Feild)  { l.log(INFO, msg, feilds...) }
func (l *asyncLogger) Warn(msg string, feilds ...Feild)  { l.log(WARN, msg, feilds...) }
func (l *asyncLogger) Error(msg string, feilds ...Feild) { l.log(ERROR, msg, feilds...) }

// With returns a *new* logger that shares the same writers/formatter/channel
// but carries extra baseFields — cheap, no locking needed since fields are
// copied, not mutated.
func (l *asyncLogger) With(feilds ...Feild) Logger {
	return &asyncLogger{
		minLevel:   l.minLevel,
		formatter:  l.formatter,
		writers:    l.writers,
		sampler:    l.sampler,
		baseFeilds: append(append([]Feild{}, l.baseFeilds...), feilds...),
		recordCh:   l.recordCh,
	}
}

// run is the single background worker draining recordCh. Only the ROOT
// logger starts this goroutine; children created via With() share it.
func (l *asyncLogger) run() {
	defer l.wg.Done()
	for rec := range l.recordCh {
		data, err := l.formatter.Format(rec)
		if err != nil {
			fmt.Fprintln(os.Stderr, "logger: format error:", err)
			continue
		}
		for _, w := range l.writers {
			if _, err := w.Write(data); err != nil {
				// One writer failing must not stop the others.
				fmt.Fprintln(os.Stderr, "logger: write error:", err)
			}
		}
	}
}

// Close flushes remaining records and stops the background goroutine.
// Safe to call multiple times.
func (l *asyncLogger) Close() {
	l.closeOnce.Do(func() {
		close(l.recordCh)
		l.wg.Wait()
	})
}

// ---------- Builder ----------

type LoggerBuilder struct {
	minLevel   LogLevel
	Formatter  Formatter
	writers    []Writer
	sampler    Sampler
	bufferSize int
}

func NewLoggerBuilder() *LoggerBuilder {
	return &LoggerBuilder{
		minLevel:   INFO,
		Formatter:  TextFormatter{},
		bufferSize: 1024,
	}
}

func (b *LoggerBuilder) WithLevel(l LogLevel) *LoggerBuilder {
	b.minLevel = l
	return b
}

func (b *LoggerBuilder) WithFormatter(f Formatter) *LoggerBuilder {
	b.Formatter = f
	return b
}

func (b *LoggerBuilder) AddWriter(w Writer) *LoggerBuilder {
	b.writers = append(b.writers, w)
	return b
}

func (b *LoggerBuilder) WithSampler(s Sampler) *LoggerBuilder {
	b.sampler = s
	return b
}

func (b *LoggerBuilder) Build() Logger {
	l := &asyncLogger{
		minLevel:  b.minLevel,
		formatter: b.Formatter,
		writers:   b.writers,
		sampler:   b.sampler,
		recordCh:  make(chan LogRecord, b.bufferSize),
	}
	l.wg.Add(1)
	go l.run()
	return l
}
