package logger

import (
	"github.com/kylin-ops/slog/kafka"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var _kafka *kafka.Produce

type Options struct {
	Level    string   `yaml:"level" json:"level"`
	Path     string   `yaml:"path" json:"path"`
	Console  bool     `yaml:"console" json:"console"`
	File     bool     `yaml:"file" json:"file"`
	RollTime int      `yaml:"roll_time" json:"roll_time"`
	Count    int      `yaml:"count" json:"count"`
	Kafka    bool     `yaml:"kafka" json:"kafka"`
	Topic    string   `yaml:"topic" json:"topic"`
	Addrs    []string `yaml:"addrs" json:"addrs"`
}

type logConfig struct {
	level    string
	path     string
	console  bool
	file     bool
	rollTime int
	count    int
	kafka    bool
	topic    string
	addrs    []string
}

func setLogger(config *logConfig) (*logrus.Logger, error) {
	var err error
	var log = logrus.New()
	log.SetReportCaller(true)
	log.SetLevel(setLogLevel(config.level))
	nullWriter, err := os.OpenFile(os.DevNull, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	log.SetOutput(nullWriter)

	if config.file && config.path != "" {
		writer, err := rotatelogs.New(
			config.path+".%Y%m%d%H%M%S",
			rotatelogs.WithLinkName(config.path),
			rotatelogs.WithRotationTime(time.Duration(config.rollTime)*time.Second*86400),
			rotatelogs.WithRotationCount(uint(config.count)),
		)
		if err != nil {
			config.console = true
			log.Error(err)
		}
		lfsHook := NewLfsHook(writer, stdoutFormatter{})
		log.AddHook(lfsHook)
	}

	if config.console {
		writer, _ := os.OpenFile(os.Stdout.Name(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		lfsHook := NewLfsHook(writer, stdoutFormatter{})
		log.AddHook(lfsHook)
	}

	if config.kafka && config.addrs != nil && config.topic != "" {
		lfsHook := NewLfsHook(nullWriter, jsonFormatter{})
		log.AddHook(lfsHook)
	}
	return log, nil
}

type logger struct {
	logger *logrus.Logger
}

func (l *logger) Infoln(args ...interface{}) {
	l.logger.Infoln(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.logger.Debugln(args...)
}

func NewLogger(option *Options) *logger {
	config := &logConfig{
		level:    option.Level,
		path:     option.Path,
		console:  option.Console,
		file:     option.File,
		rollTime: option.RollTime,
		count:    option.Count,
		kafka:    option.Kafka,
		topic:    option.Topic,
		addrs:    option.Addrs,
	}

	if config.kafka && config.addrs != nil && config.topic != "" {
		k, err := kafka.NewProduce(config.addrs, config.topic)
		if err != nil {
			panic(err)
		}
		_kafka = k
	}

	_logger, err := setLogger(config)
	if err != nil {
		panic(err)
	}
	return &logger{
		logger: _logger,
	}
}
