package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"vigtem/internal/analyzer"
	"vigtem/internal/config"
	"vigtem/internal/debugger"
	rtm "vigtem/internal/runtime"
)

type Report struct {
	StartedAt   time.Time          `json:"started_at"`
	EndedAt     time.Time          `json:"ended_at"`
	Target      string             `json:"target"`
	Workspace   string             `json:"workspace"`
	HostMeta    map[string]string  `json:"host_meta"`
	Findings    []analyzer.Finding `json:"findings"`
	TraceEvents int                `json:"trace_events"`
}

type Engine struct {
	cfg    config.Config
	tracer *debugger.Tracer
	probe  rtm.MemoryProbe
}

func New(cfg config.Config, tracer *debugger.Tracer) *Engine {
	return &Engine{cfg: cfg, tracer: tracer}
}

func (e *Engine) Run() (Report, error) {
	rep := Report{StartedAt: time.Now().UTC(), Target: e.cfg.TargetPath, Workspace: e.cfg.Workspace, HostMeta: config.HostFingerprint()}
	e.tracer.Record("INIT", "sandbox", "creating isolated workspace")
	if err := os.MkdirAll(e.cfg.Workspace, 0o755); err != nil {
		return rep, err
	}

	files, err := e.ingestTarget(e.cfg.TargetPath)
	if err != nil {
		return rep, err
	}
	for _, f := range files {
		e.tracer.Record("SCAN", filepath.Base(f), "static signature inspection")
		finding, ferr := analyzer.ScanFile(f, e.cfg.MaxFileSizeMB)
		if ferr != nil {
			e.tracer.Record("WARN", filepath.Base(f), ferr.Error())
			continue
		}
		rep.Findings = append(rep.Findings, finding)
		e.tracer.Record("MEMORY", filepath.Base(f), "buffer address "+e.probe.AddressOfString(finding.Hash))
		if e.cfg.CommandTrace {
			e.simulateCommands(f)
		}
	}

	rep.EndedAt = time.Now().UTC()
	rep.TraceEvents = len(e.tracer.Snapshot())
	return rep, e.writeReport(rep)
}

func (e *Engine) ingestTarget(target string) ([]string, error) {
	e.tracer.Record("INGEST", "filesystem", "cloning target artifacts")
	var copied []string
	err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(target, path)
		dst := filepath.Join(e.cfg.Workspace, "sample", rel)
		if mkErr := os.MkdirAll(filepath.Dir(dst), 0o755); mkErr != nil {
			return mkErr
		}
		b, rdErr := os.ReadFile(path)
		if rdErr != nil {
			return rdErr
		}
		if wrErr := os.WriteFile(dst, b, 0o644); wrErr != nil {
			return wrErr
		}
		copied = append(copied, dst)
		e.tracer.Record("FILE", rel, fmt.Sprintf("copied %d bytes addr=%s", len(b), e.probe.AddressOfBytes(b)))
		return nil
	})
	return copied, err
}

func (e *Engine) simulateCommands(file string) {
	ctx, cancel := context.WithTimeout(context.Background(), e.cfg.ExecutionTimeout)
	defer cancel()
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/c", "echo", "sandbox probe", file)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", "echo sandbox probe \""+strings.ReplaceAll(file, "\"", "")+"\"")
	}
	out, err := cmd.CombinedOutput()
	detail := strings.TrimSpace(string(out))
	if err != nil {
		detail = detail + " err=" + err.Error()
	}
	e.tracer.Record("CMD", filepath.Base(file), detail)
}

func (e *Engine) writeReport(rep Report) error {
	path := filepath.Join(e.cfg.Workspace, "analysis_report.json")
	buf, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return err
	}
	e.tracer.Record("REPORT", filepath.Base(path), "serialized analysis report")
	return os.WriteFile(path, buf, 0o644)
}
