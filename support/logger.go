package support

import (
	"io"
	"log"
)

// Logger is a multi-level logger, supporting trace, info, warning and error
// levels that can be enabled/disabled independently.
type Logger struct {
	trace   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger

	traceEnabled bool
}

// NewLogger returns a new logger instance configured with the given destination
// for each level. Setting any destination to nil will disable the logging (log
// will be sent to ioutil.Discard). Common destinations are os.Stdout and
// os.Stderr.
func NewLogger(traceDest, infoDest, warningDest, errorDest io.Writer) *Logger {
	l := &Logger{}
	if traceDest == nil {
		traceDest = io.Discard
	} else {
		l.traceEnabled = true
	}
	if infoDest == nil {
		infoDest = io.Discard
	}
	if warningDest == nil {
		warningDest = io.Discard
	}
	if errorDest == nil {
		errorDest = io.Discard
	}

	l.trace = log.New(traceDest, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.info = log.New(infoDest, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.warning = log.New(warningDest, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.error = log.New(errorDest, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return l
}

// TRACE logs trace messages, used mostly for debugging.
func (l *Logger) TRACE(format string, a ...interface{}) {
	l.trace.Printf(format, a...)
}

// INFO logs informational messages, used for reporting back informational
// messages that are generally useful.
func (l *Logger) INFO(format string, a ...interface{}) {
	l.info.Printf(format, a...)
}

// WARNING logs warning messages, used to report anything that might be
// problematic but is not considered an error.
func (l *Logger) WARNING(format string, a ...interface{}) {
	l.warning.Printf(format, a...)
}

// ERROR logs error messages, used to report actual issues that might need to
// be fixed.
func (l *Logger) ERROR(format string, a ...interface{}) {
	l.error.Printf(format, a...)
}

// TraceEnabled returns true if tracing was enabled. Used to enable output for
// external code (like unitybridge).
func (l *Logger) TraceEnabled() bool {
	return l.traceEnabled
}
