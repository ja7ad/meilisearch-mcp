package logger

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ShortStringer interface {
	ShortString() string
}

var globalInst *logger

type logger struct {
	config *Config
	subs   map[string]*SubLogger
	writer io.Writer
}

type SubLogger struct {
	logger zerolog.Logger
	name   string
	obj    fmt.Stringer
}

func getLoggersInst() *logger {
	if globalInst == nil {
		conf := &Config{
			Levels: make(map[string]string),
		}

		conf.Levels["default"] = "debug"
		conf.Levels["_transport"] = "debug"
		conf.Levels["_protocol"] = "debug"
		conf.Levels["_pool"] = "debug"

		globalInst = &logger{
			config: conf,
			subs:   make(map[string]*SubLogger),
			writer: zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"},
		}
		log.Logger = zerolog.New(globalInst.writer).With().Timestamp().Logger()
	}

	return globalInst
}

func InitGlobalLogger(conf *Config) {
	if globalInst != nil {
		return
	}

	writers := []io.Writer{}

	consoleWriter := &zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04:05",
	}
	writers = append(writers, consoleWriter)

	globalInst = &logger{
		config: conf,
		subs:   make(map[string]*SubLogger),
		writer: io.MultiWriter(writers...),
	}
	log.Logger = zerolog.New(globalInst.writer).With().Timestamp().Logger()

	lvl, err := zerolog.ParseLevel(conf.Levels["default"])
	if err != nil {
		Warn("invalid default log level", "error", err)
	}
	log.Logger = log.Logger.Level(lvl)
}

func addFields(event *zerolog.Event, keyvals ...any) *zerolog.Event {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "!MISSING-VALUE!")
	}
	for index := 0; index < len(keyvals); index += 2 {
		key, ok := keyvals[index].(string)
		if !ok {
			key = "!INVALID-KEY!"
		}

		value := keyvals[index+1]
		switch typ := value.(type) {
		case fmt.Stringer:
			if isNil(typ) {
				event.Any(key, typ)
			} else {
				event.Stringer(key, typ)
			}
		case ShortStringer:
			if isNil(typ) {
				event.Any(key, typ)
			} else {
				event.Str(key, typ.ShortString())
			}
		case error:
			event.AnErr(key, typ)
		case []byte:
			event.Str(key, hex.EncodeToString(typ))
		default:
			event.Any(key, typ)
		}
	}

	return event
}

func NewSubLogger(name string, obj fmt.Stringer) *SubLogger {
	inst := getLoggersInst()
	sub := &SubLogger{
		logger: zerolog.New(inst.writer).With().Timestamp().Logger(),
		name:   name,
		obj:    obj,
	}

	lvlStr := inst.config.Levels[name]
	if lvlStr == "" {
		lvlStr = inst.config.Levels["default"]
	}

	lvl, err := zerolog.ParseLevel(lvlStr)
	if err != nil {
		Warn("invalid log level", "error", err, "name", name)
	}
	sub.logger = sub.logger.Level(lvl)

	inst.subs[name] = sub

	return sub
}

func (sl *SubLogger) logObj(event *zerolog.Event, msg string, keyvals ...any) {
	if sl.obj != nil {
		event = event.Str(sl.name, sl.obj.String())
	}

	addFields(event, keyvals...).Msg(msg)
}

func (sl *SubLogger) Trace(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Trace(), msg, keyvals...)
}

func (sl *SubLogger) Debug(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Debug(), msg, keyvals...)
}

func (sl *SubLogger) Info(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Info(), msg, keyvals...)
}

func (sl *SubLogger) Warn(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Warn(), msg, keyvals...)
}

func (sl *SubLogger) Error(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Error().Caller(2), msg, keyvals...)
}

func (sl *SubLogger) Fatal(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Fatal().Caller(2), msg, keyvals...)
}

func (sl *SubLogger) Panic(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Panic().Caller(2), msg, keyvals...)
}

func Trace(msg string, keyvals ...any) {
	addFields(log.Trace(), keyvals...).Msg(msg)
}

func Debug(msg string, keyvals ...any) {
	addFields(log.Debug(), keyvals...).Msg(msg)
}

func Info(msg string, keyvals ...any) {
	addFields(log.Info(), keyvals...).Msg(msg)
}

func Warn(msg string, keyvals ...any) {
	addFields(log.Warn(), keyvals...).Msg(msg)
}

func Error(msg string, keyvals ...any) {
	addFields(log.Error().Caller(2), keyvals...).Msg(msg)
}

func Fatal(msg string, keyvals ...any) {
	addFields(log.Fatal().Caller(2), keyvals...).Msg(msg)
}

func Panic(msg string, keyvals ...any) {
	addFields(log.Panic().Caller(2), keyvals...).Msg(msg)
}

func isNil(i any) bool {
	if i == nil {
		return true
	}
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		return reflect.ValueOf(i).IsNil()
	}

	return false
}

func (sl *SubLogger) SetObj(obj fmt.Stringer) {
	sl.obj = obj
}
