package nsfContext

import (
	"encoding/json"
	"fmt"
	"nsf/internal/db"
	"nsf/logger"
	"sync"
)

type serviceContextIE struct {
	db.DbIf
	*logger.NsfLogger
}

// The database is the source of truth. No separate cache needs reconstruction
// after restart; each successful mutation is persisted before returning.
type serviceContext struct {
	db.DbIf
	*logger.NsfLogger
	mu     sync.RWMutex
	closed bool
}

type serviceProvider struct {
	NfName     string `json:"nf_name"`
	SbiAddress string `json:"sbi_address"`
	SbiPort    int    `json:"sbi_port"`
}

const serviceKeyPrefix = "nsf/service/"

func newServiceContext(ie *serviceContextIE) *serviceContext {
	return &serviceContext{DbIf: ie.DbIf, NsfLogger: ie.NsfLogger}
}

func (s *serviceContext) release() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	// dbContext owns the database handle and closes it after this method.
	// Do not remove persisted registrations during shutdown.
	s.NsfLogger.SvcLog.Infoln("serviceContext released; registrations preserved")
}

func (s *serviceContext) providers(service string) ([]serviceProvider, error) {
	key := serviceKeyPrefix + service
	exists, err := s.DbIf.Exist(key)
	if err != nil {
		return nil, fmt.Errorf("check service %q: %w", service, err)
	}
	if !exists {
		return nil, nil
	}
	value, err := s.DbIf.Load(key)
	if err != nil {
		return nil, fmt.Errorf("load service %q: %w", service, err)
	}
	var providers []serviceProvider
	if err := json.Unmarshal([]byte(value), &providers); err != nil {
		return nil, fmt.Errorf("decode service %q: %w", service, err)
	}
	return providers, nil
}

func (s *serviceContext) registerServices(nfName, sbiAddr string, sbiPort int, services []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return fmt.Errorf("service context is closed")
	}
	data := make(map[string]string, len(services))
	for _, service := range services {
		providers, err := s.providers(service)
		if err != nil {
			return err
		}
		provider := serviceProvider{NfName: nfName, SbiAddress: sbiAddr, SbiPort: sbiPort}
		found := false
		for i := range providers {
			if providers[i].NfName == nfName {
				providers[i] = provider
				found = true
				break
			}
		}
		if !found {
			providers = append(providers, provider)
		}
		value, err := json.Marshal(providers)
		if err != nil {
			return fmt.Errorf("encode service %q: %w", service, err)
		}
		data[serviceKeyPrefix+service] = string(value)
	}
	if len(data) == 0 {
		return nil
	}
	if err := s.DbIf.SaveAll(data); err != nil {
		return fmt.Errorf("save registrations: %w", err)
	}
	return nil
}

func (s *serviceContext) deregisterServices(nfName string, services []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return fmt.Errorf("service context is closed")
	}
	for _, service := range services {
		providers, err := s.providers(service)
		if err != nil {
			return err
		}
		remaining := make([]serviceProvider, 0, len(providers))
		removed := false
		for _, provider := range providers {
			if provider.NfName == nfName {
				removed = true
			} else {
				remaining = append(remaining, provider)
			}
		}
		if !removed {
			continue
		}
		key := serviceKeyPrefix + service
		if len(remaining) == 0 {
			if err := s.DbIf.Remove(key); err != nil {
				return fmt.Errorf("remove service %q: %w", service, err)
			}
		} else {
			value, err := json.Marshal(remaining)
			if err != nil {
				return fmt.Errorf("encode service %q: %w", service, err)
			}
			if err := s.DbIf.Save(key, string(value)); err != nil {
				return fmt.Errorf("save service %q: %w", service, err)
			}
		}
	}
	return nil
}

func (s *serviceContext) discoverService(service string) (string, string, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return "", "", 0, fmt.Errorf("service context is closed")
	}
	providers, err := s.providers(service)
	if err != nil {
		return "", "", 0, err
	}
	if len(providers) == 0 {
		return "", "", 0, fmt.Errorf("service %q not found", service)
	}
	// Stable first-registered selection; other providers survive deregistration.
	p := providers[0]
	return p.NfName, p.SbiAddress, p.SbiPort, nil
}
