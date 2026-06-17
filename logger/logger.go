package logger

import (
	loggergo "github.com/Alonza0314/logger-go/v2"
	loggergoModel "github.com/Alonza0314/logger-go/v2/model"
	loggergoUtil "github.com/Alonza0314/logger-go/v2/util"
)

type NsfLogger struct {
	*loggergo.Logger

	AppLog  loggergoModel.LoggerInterface
	SbiLog  loggergoModel.LoggerInterface
	ProcLog loggergoModel.LoggerInterface
	ConsLog loggergoModel.LoggerInterface
	CtxLog  loggergoModel.LoggerInterface
	SvcLog  loggergoModel.LoggerInterface
	DbLog   loggergoModel.LoggerInterface
}

func NewNsfLogger(level loggergoUtil.LogLevelString) NsfLogger {
	logger := loggergo.NewLogger("/tmp/nsf.log", false)
	logger.SetLevel(level)

	return NsfLogger{
		Logger: logger,

		AppLog:  logger.WithTags(LOG_NSF_TAG, LOG_APP_TAG),
		SbiLog:  logger.WithTags(LOG_NSF_TAG, LOG_SBI_TAG),
		ProcLog: logger.WithTags(LOG_NSF_TAG, LOG_SBI_TAG, LOG_PROCESSOR_TAG),
		ConsLog: logger.WithTags(LOG_NSF_TAG, LOG_SBI_TAG, LOG_CONSUMER_TAG),
		CtxLog:  logger.WithTags(LOG_NSF_TAG, LOG_CONTEXT_TAG),
		SvcLog:  logger.WithTags(LOG_NSF_TAG, LOG_CONTEXT_TAG, LOG_CONTEXT_SERVICE_TAG),
		DbLog:   logger.WithTags(LOG_NSF_TAG, LOG_CONTEXT_TAG, LOG_CONTEXT_DB_TAG),
	}
}
