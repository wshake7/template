package dynamic

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

// EnginePool 管理多个独立 Engine 实例以支持并发执行。
// NewEnginePool 需要提供一个 factory 用于创建单个 Engine 实例。
type EnginePool struct {
	pool        chan Engine
	size        int
	mu          sync.Mutex
	closed      bool
	initialized bool
	typ         EngineType
}

// NewEnginePool 创建并初始化一个包含 size 个 Engine 的池。
// factory 用于创建单个 Engine（例如 newLuaEngine）。
func NewEnginePool(size int, typ EngineType) (*EnginePool, error) {
	if size < 1 {
		return nil, errors.New("pool size must be >= 1")
	}

	p := &EnginePool{
		pool: make(chan Engine, size),
		size: size,
		typ:  typ,
	}

	// 创建并初始化子 engine
	created := make([]Engine, 0, size)
	for i := 0; i < size; i++ {
		eng, err := NewEngine(typ)
		if err != nil {
			// 清理已创建的 engines
			for _, e := range created {
				_ = e.Close()
			}
			return nil, fmt.Errorf("factory failed: %w", err)
		}

		// 调用 Init，失败则清理并返回
		if initErr := eng.Init(context.Background()); initErr != nil {
			_ = eng.Close()
			for _, e := range created {
				_ = e.Close()
			}
			return nil, fmt.Errorf("init failed: %w", initErr)
		}

		created = append(created, eng)
	}

	for _, e := range created {
		p.pool <- e
	}

	return p, nil
}

// Acquire 从池中获取一个 Engine（会阻塞直到有可用的）。
func (p *EnginePool) Acquire() (Engine, error) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil, errors.New("engine pool closed")
	}
	p.mu.Unlock()

	eng, ok := <-p.pool
	if !ok {
		return nil, errors.New("engine pool closed")
	}
	return eng, nil
}

// Release 将 Engine 放回池中；若池已关闭则关闭该 Engine。
func (p *EnginePool) Release(e Engine) {
	if e == nil {
		return
	}
	p.mu.Lock()
	closed := p.closed
	p.mu.Unlock()

	if closed {
		_ = e.Close()
		return
	}

	// 捕获并发 Close 导致的 send-on-closed panic
	defer func() {
		if r := recover(); r != nil {
			_ = e.Close()
		}
	}()

	select {
	case p.pool <- e:
	default:
		_ = e.Close()
	}
}

func (p *EnginePool) SafeAcquire(fn func(Engine) error) error {
	e, err := p.Acquire()
	if err != nil {
		return err
	}
	defer p.Release(e)
	return fn(e)
}

func (p *EnginePool) GetType() EngineType {
	return p.typ
}

// IsClosed 返回池是否已关闭。
func (p *EnginePool) IsClosed() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.closed
}

func (p *EnginePool) Init(ctx context.Context) error {
	// 尝试获取池中所有实例
	engines := make([]Engine, 0, p.size)
	for i := 0; i < p.size; i++ {
		eng, err := p.Acquire()
		if err != nil {
			for _, e := range engines {
				_ = e.Close()
			}
			return err
		}
		engines = append(engines, eng)
	}

	// 对每个实例执行 Init()
	for _, eng := range engines {
		if err := eng.Init(ctx); err != nil {
			for _, e := range engines {
				_ = e.Close()
			}
			return fmt.Errorf("init failed: %w", err)
		}
	}

	// 释放回池
	for _, eng := range engines {
		p.Release(eng)
	}
	p.initialized = true
	return nil
}

func (p *EnginePool) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	close(p.pool)
	p.mu.Unlock()

	var lastErr error
	for eng := range p.pool {
		if err := eng.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (p *EnginePool) IsInitialized() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.initialized
}

func (p *EnginePool) LoadString(ctx context.Context, sources ...string) error {
	//TODO implement me
	panic("implement me")
}

func (p *EnginePool) LoadFile(ctx context.Context, filePath ...string) error {
	//TODO implement me
	panic("implement me")
}

func (p *EnginePool) LoadReader(ctx context.Context, reader io.Reader, name string) error {
	//TODO implement me
	panic("implement me")
}

func (p *EnginePool) ExecuteLoaded(ctx context.Context) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (p *EnginePool) ExecuteString(ctx context.Context, source ...string) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (p *EnginePool) ExecuteFile(ctx context.Context, filePath ...string) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (p *EnginePool) RegisterGlobal(name string, value any) error {
	//TODO implement me
	panic("implement me")
}

func (p *EnginePool) GetGlobal(name string) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (p *EnginePool) RegisterFunction(name string, fn any) error {
	//TODO implement me
	panic("implement me")
}

func (p *EnginePool) CallFunction(ctx context.Context, name string, args ...any) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (p *EnginePool) RegisterModule(name string, module any) error {
	//TODO implement me
	panic("implement me")
}

func (p *EnginePool) GetLastError() error {
	//TODO implement me
	panic("implement me")
}

func (p *EnginePool) ClearError() {
	//TODO implement me
	panic("implement me")
}
