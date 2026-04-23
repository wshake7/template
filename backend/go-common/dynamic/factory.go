package dynamic

import (
	"fmt"
	"sync"
)

// FactoryFunc 是用于创建 Engine 实例的工厂函数类型。
type FactoryFunc func() (Engine, error)

var (
	factoryMu sync.RWMutex
	factories = make(map[EngineType]FactoryFunc)
)

// NewEngine 使用已注册的工厂函数创建一个 Engine 实例。
func NewEngine(typ EngineType) (Engine, error) {
	f, ok := GetFactory(typ)
	if !ok {
		return nil, fmt.Errorf("script engine factory %s not registered", typ)
	}
	return f()
}

// Register registers a FactoryFunc for a given Type.
func Register(typ EngineType, f FactoryFunc) error {
	factoryMu.Lock()
	defer factoryMu.Unlock()
	if _, ok := factories[typ]; ok {
		return fmt.Errorf("script engine factory %s already registered", typ)
	}
	factories[typ] = f
	return nil
}

// GetFactory returns a registered FactoryFunc for a given Type and whether it existed.
func GetFactory(typ EngineType) (FactoryFunc, bool) {
	factoryMu.RLock()
	defer factoryMu.RUnlock()
	f, ok := factories[typ]
	return f, ok
}

// ListFactories returns a slice of currently registered Types.
func ListFactories() []EngineType {
	factoryMu.RLock()
	defer factoryMu.RUnlock()
	res := make([]EngineType, 0, len(factories))
	for k := range factories {
		res = append(res, k)
	}
	return res
}

// Unregister removes a registered factory by Type. It returns true if a factory was removed.
func Unregister(typ EngineType) bool {
	factoryMu.Lock()
	defer factoryMu.Unlock()
	if _, ok := factories[typ]; ok {
		delete(factories, typ)
		return true
	}
	return false
}
