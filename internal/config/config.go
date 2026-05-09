package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Config struct {
	TargetPath         string        `yaml:"target_path"`
	Workspace          string        `yaml:"workspace"`
	ExecutionTimeout   time.Duration `yaml:"execution_timeout"`
	MemoryTraceEnabled bool          `yaml:"memory_trace_enabled"`
	CommandTrace       bool          `yaml:"command_trace"`
	TaskTrace          bool          `yaml:"task_trace"`
	MaxFileSizeMB      int64         `yaml:"max_file_size_mb"`
}

func Default(target string) Config {
	workRoot := filepath.Join(os.TempDir(), "vigtem")
	_ = os.MkdirAll(workRoot, 0o755)

	return Config{
		TargetPath:         target,
		Workspace:          filepath.Join(workRoot, time.Now().Format("20060102_150405")),
		ExecutionTimeout:   30 * time.Second,
		MemoryTraceEnabled: true,
		CommandTrace:       true,
		TaskTrace:          true,
		MaxFileSizeMB:      32,
	}
}

func HostFingerprint() map[string]string {
	name, _ := os.Hostname()
	return map[string]string{
		"goos":   runtime.GOOS,
		"goarch": runtime.GOARCH,
		"host":   name,
	}
}
