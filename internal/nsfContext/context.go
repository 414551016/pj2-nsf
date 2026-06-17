package nsfContext

import "nsf/logger"

type NsfContextIE struct {
	DbType string

	*logger.NsfLogger
}

type NsfContext struct {
	*serviceContext
	*dbContext

	*logger.NsfLogger
}

func NewNsfContext(nsfContextIE *NsfContextIE) *NsfContext {
	dbContext, err := newDbContext(&dbContextIE{
		dbType: nsfContextIE.DbType,

		NsfLogger: nsfContextIE.NsfLogger,
	})
	if err != nil {
		nsfContextIE.NsfLogger.CtxLog.Errorf("Failed to create dbContext: %v", err)
		return nil
	}
	serviceContext := newServiceContext(&serviceContextIE{
		DbIf: dbContext.db,

		NsfLogger: nsfContextIE.NsfLogger,
	})

	return &NsfContext{
		serviceContext: serviceContext,
		dbContext:      dbContext,

		NsfLogger: nsfContextIE.NsfLogger,
	}
}

func (c *NsfContext) Release() {
	c.NsfLogger.CtxLog.Infoln("Release NsfContext...")

	c.serviceContext.release()
	c.dbContext.release()

	c.NsfLogger.CtxLog.Infoln("NsfContext released")
}

func (c *NsfContext) RegisterServices(nfName, address string, port int, services []string) error {
	return c.serviceContext.registerServices(nfName, address, port, services)
}

func (c *NsfContext) DeregisterServices(nfName string, services []string) error {
	return c.serviceContext.deregisterServices(nfName, services)
}

func (c *NsfContext) DiscoverService(service string) (string, string, int, error) {
	return c.serviceContext.discoverService(service)
}
