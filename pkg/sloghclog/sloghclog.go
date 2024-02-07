// Transform slog.Logger into hclog.Logger instance
package sloghclog

import (
	"context"
	"io"
	"log"
	"log/slog"

	"github.com/hashicorp/go-hclog"
)

func convertLevel(level hclog.Level) slog.Level {
	if level == hclog.NoLevel {
		return slog.LevelInfo
	}
	return slog.Level((level - 3) * 4)
}

func Adapt(logger *slog.Logger) hclog.Logger {
	return &Adapter{
		origLogger: logger,
		logger:     logger,
	}
}

type Adapter struct {
	origLogger *slog.Logger
	logger     *slog.Logger
	name       string
	with       []any
}

const levelTrace slog.Level = slog.Level(-8)

// Args are alternating key, val pairs
// keys must be strings
// vals can be any type, but display is implementation specific
// Emit a message and key/value pairs at a provided log level
func (a *Adapter) Log(level hclog.Level, msg string, args ...interface{}) {
	a.logger.Log(context.Background(), convertLevel(level), msg, args...)
}

// Emit a message and key/value pairs at the TRACE level
func (a *Adapter) Trace(msg string, args ...interface{}) {
	a.logger.Log(context.Background(), levelTrace, msg, args...)
}

// Emit a message and key/value pairs at the DEBUG level
func (a *Adapter) Debug(msg string, args ...interface{}) {
	a.logger.Debug(msg, args...)
}

// Emit a message and key/value pairs at the INFO level
func (a *Adapter) Info(msg string, args ...interface{}) {
	a.logger.Info(msg, args...)
}

// Emit a message and key/value pairs at the WARN level
func (a *Adapter) Warn(msg string, args ...interface{}) {
	a.logger.Warn(msg, args...)
}

// Emit a message and key/value pairs at the ERROR level
func (a *Adapter) Error(msg string, args ...interface{}) {
	a.logger.Error(msg, args...)
}

// Indicate if TRACE logs would be emitted. This and the other Is* guards
// are used to elide expensive logging code based on the current level.
func (a *Adapter) IsTrace() bool {
	return a.logger.Enabled(context.Background(), levelTrace)
}

// Indicate if DEBUG logs would be emitted. This and the other Is* guards
func (a *Adapter) IsDebug() bool {
	return a.logger.Enabled(context.Background(), slog.LevelDebug)
}

// Indicate if INFO logs would be emitted. This and the other Is* guards
func (a *Adapter) IsInfo() bool {
	return a.logger.Enabled(context.Background(), slog.LevelInfo)
}

// Indicate if WARN logs would be emitted. This and the other Is* guards
func (a *Adapter) IsWarn() bool {
	return a.logger.Enabled(context.Background(), slog.LevelWarn)
}

// Indicate if ERROR logs would be emitted. This and the other Is* guards
func (a *Adapter) IsError() bool {
	return a.logger.Enabled(context.Background(), slog.LevelError)
}

// ImpliedArgs returns With key/value pairs
func (a *Adapter) ImpliedArgs() []interface{} {
	return a.with
}

// Creates a sublogger that will always have the given key/value pairs
func (a *Adapter) With(args ...interface{}) hclog.Logger {
	dup := &Adapter{
		origLogger: a.origLogger,
		name:       a.name,
	}

	length := len(a.with) + len(args)
	if a.name != "" {
		length += 2
	}

	with := make([]any, 0, length)
	if a.name != "" {
		with = append(with, "name", a.name)
	}
	with = append(with, a.with...)
	with = append(with, args...)

	if a.name != "" {
		dup.with = with[2:]
	} else {
		dup.with = with
	}

	if len(with) != 0 {
		a.logger = a.origLogger.With(with...)
	} else {
		a.logger = a.origLogger
	}
	return dup
}

func (a *Adapter) applyName() {
	var with []any
	if a.name != "" {
		with = make([]any, 0, len(a.with)+2)
		with = append(with, "name", a.name)
		with = append(with, a.with...)
	} else {
		with = a.with
	}
	if len(with) != 0 {
		a.logger = a.origLogger.With(with...)
	} else {
		a.logger = a.origLogger
	}
}

// Returns the Name of the logger
func (a *Adapter) Name() string {
	return a.name
}

// Create a logger that will prepend the name string on the front of all messages.
// If the logger already has a name, the new value will be appended to the current
// name. That way, a major subsystem can use this to decorate all it's own logs
// without losing context.
func (a *Adapter) Named(name string) hclog.Logger {
	dup := *a
	if dup.name != "" {
		dup.name = dup.name + "." + name
	} else {
		dup.name = name
	}
	dup.applyName()
	return &dup
}

// Create a logger that will prepend the name string on the front of all messages.
// This sets the name of the logger to the value directly, unlike Named which honor
// the current name as well.
func (a *Adapter) ResetNamed(name string) hclog.Logger {
	dup := *a
	dup.name = name
	dup.applyName()
	return &dup
}

// Updates the level. This should affect all sub-loggers as well. If an
// implementation cannot update the level on the fly, it should no-op.
func (a *Adapter) SetLevel(level hclog.Level) {
	// no-op: we can't be sure that we can update the level
}

// Return a value that conforms to the stdlib log.Logger interface
func (a *Adapter) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	if opts == nil {
		opts = &hclog.StandardLoggerOptions{}
	}
	// cant infer levels

	return log.New(a.StandardWriter(opts), "", 0)
}

// Return a value that conforms to io.Writer, which can be passed into log.SetOutput()
func (a *Adapter) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return &stdlogAdapter{
		log:   a,
		level: convertLevel(opts.ForceLevel),
	}
}

type stdlogAdapter struct {
	log   *Adapter
	level slog.Level
}

// Write implements io.Writer.
func (sa *stdlogAdapter) Write(data []byte) (n int, err error) {
	sa.log.logger.Log(context.Background(), sa.level, string(data))
	return len(data), nil
}

var _ hclog.Logger = (*Adapter)(nil)
