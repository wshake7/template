package dynamic

import (
	"context"
	"errors"
	"sync"
)

// Manager 管理多个 Engine 实例的生命周期与访问。
// - 适用于需要多个引擎实例、统一 Init/Close、或按 name 获取的场景。
// - 若项目只需要单个全局 Engine，可不使用 Manager。
type Manager struct {
	mu          sync.RWMutex
	engines     map[EngineType]Engine
	defaultName EngineType
}

// NewManager 创建 Manager。
func NewManager() *Manager {
	return &Manager{
		engines: make(map[EngineType]Engine),
	}
}

// Register 注册一个 Engine（不初始化）。
// 若 name 已存在返回错误。
func (m *Manager) Register(typ EngineType, e Engine) error {
	if typ == "" || e == nil {
		return errors.New("invalid name or engine")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.engines[typ]; ok {
		return errors.New("engine already registered")
	}
	m.engines[typ] = e
	return nil
}

// Get 返回已注册的 Engine。
func (m *Manager) Get(typ EngineType) (Engine, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	eng, ok := m.engines[typ]
	return eng, ok
}

func (m *Manager) Init(ctx context.Context, typ EngineType) error {
	if typ == "" {
		return errors.New("engine name is empty")
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if e, ok := m.engines[typ]; ok {
		return e.Init(ctx)
	}
	return nil
}

// InitAll 对所有已注册引擎执行 Init。
func (m *Manager) InitAll(ctx context.Context) error {
	m.mu.RLock()
	list := make([]Engine, 0, len(m.engines))
	for _, e := range m.engines {
		list = append(list, e)
	}
	m.mu.RUnlock()

	for _, e := range list {
		if err := e.Init(ctx); err != nil {
			return err
		}
	}
	return nil
}

// CloseAll 关闭所有已注册引擎（并忽略单个 Close 错误，返回最后一个错误）。
func (m *Manager) CloseAll() error {
	m.mu.Lock()
	list := make([]Engine, 0, len(m.engines))
	for _, e := range m.engines {
		list = append(list, e)
	}
	// 清空注册表以防重复 Close
	m.engines = make(map[EngineType]Engine)
	m.mu.Unlock()

	var lastErr error
	for _, e := range list {
		if err := e.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Remove 注销并可选择关闭该 Engine（若 closeIfExists 为 true）。
func (m *Manager) Remove(typ EngineType, closeIfExists bool) {
	m.mu.Lock()
	e, ok := m.engines[typ]
	if ok {
		delete(m.engines, typ)
	}
	m.mu.Unlock()

	if ok && closeIfExists {
		_ = e.Close()
	}
}

// SetDefault 设置默认引擎名，便于不指定 name 时使用。
func (m *Manager) SetDefault(typ EngineType) {
	m.mu.Lock()
	m.defaultName = typ
	m.mu.Unlock()
}

// GetDefault 获取默认引擎。
func (m *Manager) GetDefault() (Engine, bool) {
	return m.Get(m.defaultName)
}
