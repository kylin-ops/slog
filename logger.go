package logger

import (
	"net"
	"os"
	"time"

	"github.com/kylin-ops/slog/kafka"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var _kafka *kafka.Produce

var (
	_appName  string
	_ipaddr   string
	_hostname string
	_options  *Options
)

type Options struct {
	// 定义日志级别
	Level string `yaml:"level" json:"level"`
	// 定义是否打印日志文件函数行号
	Caller bool `yaml:"caller" json:"caller"`
	// 定义是否打印到console
	Console bool `yaml:"console" json:"console"`
	// 定义日志是否打印到文件
	File bool `yaml:"file" json:"file"`
	// 定义日志路径
	Path string `yaml:"path" json:"path"`
	// 定义输出到文件的日志是否为json格式
	Json bool `yaml:"format" json:"format"`
	// 定义日志滚动时间(天)
	RollTime int `yaml:"roll_time" json:"roll_time"`
	// 定义日志保存数量
	Count int `yaml:"count" json:"count"`
	// 定义日志是否发送到kafka
	Kafka bool `yaml:"kafka" json:"kafka"`
	// 定义kafka的topic
	Topic string `yaml:"topic" json:"topic"`
	// 定义kafka的地址
	Addrs []string `yaml:"addrs" json:"addrs"`
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
	caller   bool
}

func _getIpFromNetIf() string {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() {
			ipv4 := ipNet.IP.To4()
			if ipv4 != nil {
				return ipv4.String()
			}
		}
	}
	return ""
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
		if _options.Json {
			lfsHook := NewLfsHook(writer, jsonFormatter{})
			log.AddHook(lfsHook)
		} else {
			lfsHook := NewLfsHook(writer, stdoutFormatter{})
			log.AddHook(lfsHook)
		}
	}

	if config.console {
		writer, _ := os.OpenFile(os.Stdout.Name(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		lfsHook := NewLfsHook(writer, stdoutFormatter{})
		log.AddHook(lfsHook)
	}

	if config.kafka && config.addrs != nil && config.topic != "" {
		lfsHook := NewLfsHook(nullWriter, kafkaFormatter{})
		log.AddHook(lfsHook)
	}
	return log, nil
}

type slogger struct {
	*logrus.Logger
}

type Logger interface {
	Infoln(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Debugln(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorln(args ...interface{})
	Errorf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnln(args ...interface{})
	Warnf(format string, args ...interface{})
	Trace(args ...interface{})
	Traceln(args ...interface{})
	Tracef(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalln(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicln(args ...interface{})
	Panicf(format string, args ...interface{})
}

func NewLogger(option *Options, appName ...string) Logger {
	_hostname, _ = os.Hostname()
	_ipaddr = _getIpFromNetIf()
	_options = option
	if appName != nil {
		_appName = appName[0]
	}

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
	return &slogger{
		Logger: _logger,
	}
}
