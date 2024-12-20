package util

import (
	"io"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger   zerolog.Logger
	errorMap map[string]string
}

var defaultErrorMap = map[string]string{
	"argumentError":                  "An argument given to this functionality has either missing or wrong.",
	"passwordIncorrect":              "Entered password is not correct.",
	"argumentErrorMissing":           "An argument is missing from the call",
	"argumentErrorUnknownStartPoint": "Given start point is not implemented or possible.",
	"ioErrorReadWriteMismatch":       "Read byte count is not matching written byte count",
	"notImplementedOwnerType":        "Given owner type is not implemented",
	"argumentErrorKeyEmpty":          "The argument \"Key\" is missing from the call",
	"argumentErrorKeyNotFound":       "The given key is not found.",
	"argumentErrorIDNotFound":        "The given id is not found",
	"argumentErrorMissingFromID":     "The language id for \"From\" is missing from the call",
	"argumentErrorNote":              "Note argument is missing from the call",
	"argumentErrorLanguage":          "Language argument is missing from the call",
	"configNotLoadedProperly":        "App config is not read or loaded correctly.\n Terminating",
}

func NewLogger(errorMap map[string]string, w io.Writer) *Logger {
	if len(errorMap) == 0 {
		errorMap = defaultErrorMap
	}
	return &Logger{
		errorMap: errorMap,
		logger:   zerolog.New(w),
	}
}

func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

func (l *Logger) Err(e error) string {
	msg := e.Error()
	log, ok := l.errorMap[msg]
	if !ok {
		log = "UnknownError"
	}
	l.logger.Err(e).Msg(log)
	return log
}
