package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func logger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 格式
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 输出
	syncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))

	//
	fileCore := zapcore.NewCore(encoder, syncer, zapcore.DebugLevel)

	return zap.New(fileCore, zap.AddCaller())
}

func data(values ...interface{}) []zap.Field {
	log := zap.NewExample()
	if len(values) == 0 || len(values)%2 != 0 {
		log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", values))
		return nil
	}
	// Zap.Field is used when keyvals pairs appear
	var fields []zap.Field
	for i := 0; i < len(values); i += 2 {
		fields = append(fields, zap.Any(fmt.Sprint(values[i]), values[i+1]))
	}
	return fields
}
