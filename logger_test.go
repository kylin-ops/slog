package slog

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	InitLogger(&Options{
		Level:    "debug",
		Path:     "/tmp/log_test.log",
		Console:  true,
		File:     true,
		RollTime: 1,
		Count:    5,
		// Json: true,
	})
	Info("aaaaaaa")
	Infoln("aaaaaaa")
	Infof("aaaaaaa")
	Error("aaaaaaa")
	Errorln("aaaaaaa")
	Errorf("aaaaaaa")
	Debug("aaaaaaa")
	Debugln("aaaaaaa")
	Debugf("aaaaaaa")
	Trace("aaaaaaa")
	Traceln("aaaaaaa")
	Tracef("aaaaaaa")
	Warn("aaaaaaa")
	Warnln("aaaaaaa")
	Warnf("aaaaaaa")
	Panic("aaaaaaa")
	Panicln("aaaaaaa")
	Panicf("aaaaaaa")
	Fatal("aaaaaaa")
	Fatalln("aaaaaaa")
	Fatalf("aaaaaaa")
}
