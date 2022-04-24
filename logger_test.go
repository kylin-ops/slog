package logger

import "testing"

func TestNewLogger(t *testing.T) {
	_logger := NewLogger(&Options{
		Level:    "debug",
		Path:     "/tmp/log_test.log",
		Console:  true,
		File:     true,
		RollTime: 1,
		Count:    5,
	})

	_logger.Infoln("this is info")
	_logger.Debug("this is debug")
}
