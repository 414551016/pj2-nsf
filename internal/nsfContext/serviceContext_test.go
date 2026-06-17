package nsfContext_test

import "testing"

/*
This test file is designed to validate the functionality of the service context (ServiceContext) in the nsfContext package.
It includes tests for service registration, discovery, and deregistration.
Each test case checks for the correct behavior of the ServiceContext methods and ensures that services are correctly registered, discovered, and deregistered within the context.
*/

var testServiceContextCases = []struct {
	name       string
	nfName     string
	sbiAddress string
	sbiPort    int
	services   []string
}{
	{
		name:       "Service 1",
		nfName:     "NF1",
		sbiAddress: "127.0.0.1",
		sbiPort:    8080,
		services:   []string{"service1", "service2"},
	},
	{
		name:       "Service 2",
		nfName:     "NF2",
		sbiAddress: "127.0.0.2",
		sbiPort:    8081,
		services:   []string{"service3", "service4"},
	},
}

func TestServiceContext(t *testing.T) {
	t.Run("Service Registration", testServiceRegistration)
	t.Run("Service Discovery", testServiceDiscovery)
	t.Run("Service Deregistration", testServiceDeregistration)
}

func testServiceRegistration(t *testing.T) {
	for _, tc := range testServiceContextCases {
		if err := nsfCtx.RegisterServices(tc.nfName, tc.sbiAddress, tc.sbiPort, tc.services); err != nil {
			t.Errorf("Failed to register services for %s: %v", tc.name, err)
		}
	}
}

func testServiceDiscovery(t *testing.T) {
	for _, tc := range testServiceContextCases {
		for _, service := range tc.services {
			nfName, sbiAddr, sbiPort, err := nsfCtx.DiscoverService(service)
			if err != nil {
				t.Errorf("Failed to discover service %s: %v", service, err)
			} else if nfName != tc.nfName || sbiAddr != tc.sbiAddress || sbiPort != tc.sbiPort {
				t.Errorf("Discovered service %s has incorrect details: got (%s, %s, %d), want (%s, %s, %d)", service, nfName, sbiAddr, sbiPort, tc.nfName, tc.sbiAddress, tc.sbiPort)
			}
		}
	}
}

func testServiceDeregistration(t *testing.T) {
	for _, tc := range testServiceContextCases {
		if err := nsfCtx.DeregisterServices(tc.nfName, tc.services); err != nil {
			t.Errorf("Failed to deregister services for %s: %v", tc.name, err)
		}
	}
}
