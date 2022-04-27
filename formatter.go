package slog

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type logMessage struct {
	AppName   interface{} `json:"appName,omitempty"`
	Class     interface{} `json:"class,omitempty"`
	Timestamp interface{} `json:"timestamp,omitempty"`
	Level     interface{} `json:"level,omitempty"`
	Message   interface{} `json:"message,omitempty"`
	Type      interface{} `json:"Type,omitempty"`
	TraceId   interface{} `json:"traceId,omitempty"`
	SpanId    interface{} `json:"spanId,omitempty"`
	ParentId  interface{} `json:"parentId,omitempty"`
	Host      interface{} `json:"host,omitempty"`
	Ip        interface{} `json:"ip,omitempty"`
}

type jsonFormatter struct{}

func (s jsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var msg = logMessage{
		Timestamp: time.Now().Format("2006-01-02T15:04:05.999+0800"),
		Level:     strings.ToUpper(entry.Level.String()),
		Message:   entry.Message,
		Type:      entry.Data["type"],
		TraceId:   entry.Data["traceID"],
		SpanId:    entry.Data["spanID"],
		ParentId:  entry.Data["parentID"],
	}
	if _hostname != "" {
		msg.Host = _hostname
	}
	if _ipaddr != "" {
		msg.Ip = _ipaddr
	}
	if _appName != "" {
		msg.AppName = _appName
	}
	//if _options.Caller {
	//	msg.Class = entry.Caller.File + ":" + entry.Caller.Function + ":" + strconv.Itoa(entry.Caller.Line)
	//}

	m, _ := json.Marshal(msg)
	m = append(m, '\n')
	return m, nil
}

type stdoutFormatter struct{}

func (s stdoutFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var msg string
	_level := entry.Level.String()
	_timeStr := time.Now().Format("2006-01-02T15:04:05.999+0800")
	msg = fmt.Sprintf("%-30s %-9s %-8s\n", _timeStr, strings.ToUpper(_level), entry.Message)
	//if _options.Caller {
	//	_caller := entry.Caller.File + ":" + entry.Caller.Function + ":" + strconv.Itoa(entry.Caller.Line)
	//	msg = fmt.Sprintf("%-30s %-9s %s   %s\n", _timeStr, strings.ToUpper(_level), _caller, entry.Message)
	//}

	return []byte(msg), nil
}

type kafkaFormatter struct{}

func (s kafkaFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var msg = logMessage{
		Timestamp: time.Now().Format("2006-01-02T15:04:05.999+0800"),
		Level:     strings.ToUpper(entry.Level.String()),
		Message:   entry.Message,
		Type:      entry.Data["type"],
		TraceId:   entry.Data["traceID"],
		SpanId:    entry.Data["spanID"],
		ParentId:  entry.Data["parentID"],
	}
	if _hostname != "" {
		msg.Host = _hostname
	}
	if _ipaddr != "" {
		msg.Ip = _ipaddr
	}
	if _appName != "" {
		msg.AppName = _appName
	}
	//if _options.Caller {
	//	msg.Class = entry.Caller.File + ":" + entry.Caller.Function + ":" + strconv.Itoa(entry.Caller.Line)
	//}
	m, _ := json.Marshal(msg)
	err := _kafka.SendSingleTopicMessage("log", m)
	return m, err
}
