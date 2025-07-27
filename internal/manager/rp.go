package manager

import "sync"

type RPManager struct {
	context map[string]string
	mu      sync.RWMutex
}

func NewRPManager() *RPManager {
	return &RPManager{}
}

func (m *RPManager) AddRP(domainName string, port string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.context[domainName] = port
}

func (m *RPManager) RemoveRP(domainName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.context, domainName)
}

func (m *RPManager) GetContext() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.context
}

func (m *RPManager) LoadContext(value map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.context = value
}
