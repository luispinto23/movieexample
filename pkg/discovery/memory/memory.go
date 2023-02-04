package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/luispinto23/movieexample/pkg/discovery"
)

type serviceName string
type instanceID string

// Registry defines an in-memory service registry.
type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

// NewRegistry creates a new in-memory service registry instance.
func NewRegistry() *Registry {
	return &Registry{
		serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{},
	}
}

// Register creates a service record in the registry.
func (r *Registry) Register(ctx context.Context, instID string, svcName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName(svcName)]; !ok {
		r.serviceAddrs[serviceName(svcName)] = map[instanceID]*serviceInstance{}
	}
	r.serviceAddrs[serviceName(svcName)][instanceID(instID)] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}

	return nil
}

// Deregister removes a service record from the registry
func (r *Registry) Deregister(ctx context.Context, instID string, svcName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName(svcName)]; !ok {
		return nil
	}
	delete(r.serviceAddrs[serviceName(svcName)], instanceID(instID))
	return nil
}

// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(instID string, svcName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName(svcName)]; !ok {
		return errors.New("service is not registered")
	}

	if _, ok := r.serviceAddrs[serviceName(svcName)][instanceID(instID)]; !ok {
		return errors.New("service instance is not registered")
	}
	r.serviceAddrs[serviceName(svcName)][instanceID(instID)].lastActive = time.Now()

	return nil
}

// ServiceAddresses returns the list of addresses of active instances of the given service.
func (r *Registry) ServiceAddresses(ctx context.Context, svcName string) ([]string, error) {
	r.Lock()
	defer r.Unlock()

	if len(r.serviceAddrs[serviceName(svcName)]) == 0 {
		return nil, discovery.ErrNotFound
	}

	var res []string

	for _, i := range r.serviceAddrs[serviceName(svcName)] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}
