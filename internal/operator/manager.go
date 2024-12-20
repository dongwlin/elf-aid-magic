package operator

import "sync"

type Manager struct {
	operators map[string]*Operator
	mutex     sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		operators: make(map[string]*Operator, 1),
	}
}

func (m *Manager) AddOperator(operator *Operator) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.operators[operator.ID]; exists {
		return false
	}

	m.operators[operator.ID] = operator
	return true
}

func (m *Manager) RemoveOperatorByID(id string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.operators[id]; exists {
		return false
	}

	delete(m.operators, id)
	return true
}

func (m *Manager) GetOperatorByID(id string) (*Operator, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	operator, exists := m.operators[id]
	return operator, exists
}
