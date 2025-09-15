package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.SugaredLogger, error) {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	baseEncoder := zapcore.NewJSONEncoder(encoderCfg)
	prettyCore := zapcore.NewCore(&prettyJSONEncoder{baseEncoder}, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	logger := zap.New(prettyCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger.Sugar(), nil
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

type prettyJSONEncoder struct {
	zapcore.Encoder
}

func (p *prettyJSONEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := p.Encoder.EncodeEntry(ent, fields)
	if err != nil {
		return nil, err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		return buf, nil
	}

	var pretty bytes.Buffer
	enc := json.NewEncoder(&pretty)
	enc.SetIndent("", "    ")
	if err := enc.Encode(raw); err != nil {
		return buf, nil
	}

	finalBuf := buffer.NewPool().Get()
	finalBuf.AppendString(pretty.String())
	buf.Free()
	return finalBuf, nil
}
