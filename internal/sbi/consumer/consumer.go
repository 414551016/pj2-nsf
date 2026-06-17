package consumer

import "nsf/logger"

type ConsumerIE struct {
	*logger.NsfLogger
}

type Consumer struct {
	*logger.NsfLogger
}

func NewConsumer(consumerIE *ConsumerIE) *Consumer {
	return &Consumer{
		NsfLogger: consumerIE.NsfLogger,
	}
}

func (c *Consumer) Release() {
	c.NsfLogger.ConsLog.Infoln("Release consumer...")
	c.NsfLogger.ConsLog.Infoln("Consumer released")
}
