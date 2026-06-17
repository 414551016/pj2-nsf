package nsfContext

import (
	"nsf/internal/db"
	"nsf/logger"
)


// TODO: You may need to add some fields
type serviceContextIE struct {
	db.DbIf

	*logger.NsfLogger
}

// TODO: You may need to add some fields
type serviceContext struct {
	db.DbIf

	*logger.NsfLogger
}

// TODO: Finish the initialization of the serviceContext instance
func newServiceContext(serviceContextIE *serviceContextIE) *serviceContext {

	return &serviceContext{
		DbIf: serviceContextIE.DbIf,

		NsfLogger: serviceContextIE.NsfLogger,
	}
}

// TODO: Finish the release of the serviceContext instance
func (s *serviceContext) release() {
	s.NsfLogger.SvcLog.Infoln("Release serviceContext...")

	s.NsfLogger.SvcLog.Infoln("serviceContext released")
}

// TODO: Finish the implementation of the registerServices method with an error return type
func (s *serviceContext) registerServices(nfName, sbiAddr string, sbiPort int, services []string) error {
	s.NsfLogger.SvcLog.Tracef("Registering %s services", nfName)

}

// TODO: Finish the implementation of the deregisterServices method with an error return type
func (s *serviceContext) deregisterServices(nfName string, services []string) error {
	s.NsfLogger.SvcLog.Tracef("De-registering %s services", nfName)

}

// TODO: Finish the implementation of the discoverService method with these return values:
// (string, string, int, error) representing the NF name, SBI address, SBI port, and an error if any
func (s *serviceContext) discoverService(service string) (string, string, int, error) {
	s.NsfLogger.SvcLog.Tracef("Discovering service %s", service)

}
