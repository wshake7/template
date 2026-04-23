package dynamic

import (
	"context"
	"io"
	"time"
)

type EngineType string

const (
	// EngineLuaType Lua Script Engine Type
	EngineLuaType EngineType = "lua"

	// EngineJavaScriptType JavaScript Script Engine Type
	EngineJavaScriptType EngineType = "javascript"

	// EnginePythonType Python Script Engine Type
	EnginePythonType EngineType = "python"
)

type Engine interface {
	// GetType get the type of the script engine
	GetType() EngineType
	// Init initialize the script engine
	Init(ctx context.Context) error
	// Close the script engine and release resources
	Close() error
	// IsInitialized check if the engine is initialized
	IsInitialized() bool

	// LoadString load multiple scripts from string sources
	LoadString(ctx context.Context, sources ...string) error
	// LoadFile load script from file path
	LoadFile(ctx context.Context, filePaths ...string) error
	// LoadReader load script from io.Reader
	LoadReader(ctx context.Context, reader io.Reader, name string) error

	// ExecuteLoaded execute the previously loaded script(s)
	ExecuteLoaded(ctx context.Context) (any, error)
	// ExecuteString execute script from string source
	ExecuteString(ctx context.Context, source string) (any, error)
	// ExecuteFile execute script from file path
	ExecuteFile(ctx context.Context, filePath string) (any, error)
	// ExecuteStrings execute multiple scripts from string sources (immediate execution)
	ExecuteStrings(ctx context.Context, sources []string) ([]any, error)
	// ExecuteFiles execute multiple scripts from file paths (immediate execution)
	ExecuteFiles(ctx context.Context, filePaths []string) ([]any, error)

	// RegisterGlobal register a global variable
	RegisterGlobal(name string, value any) error
	// GetGlobal get a global variable
	GetGlobal(name string) (any, error)

	// RegisterFunction register a function with the given name
	RegisterFunction(name string, fn any) error
	// CallFunction call a function with the given name and arguments
	CallFunction(ctx context.Context, name string, args ...any) (any, error)

	// RegisterModule register a module with the given name
	RegisterModule(name string, module any) error

	// GetLastError get the last error occurred in the engine
	GetLastError() error
	// ClearError clear the last error
	ClearError()
}

// CallResult 函数调用结果
type CallResult struct {
	Values []any
	Error  error
}

// ExecuteOptions 执行选项
type ExecuteOptions struct {
	Timeout  time.Duration
	Globals  map[string]any
	MaxStack int
}
