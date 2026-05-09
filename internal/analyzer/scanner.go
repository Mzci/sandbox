package analyzer

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Finding struct {
	File       string
	Hash       string
	Suspicious []string
	Size       int64
}

var suspiciousPatterns = []string{
	"powershell -enc", "cmd.exe /c", "rundll32", "CreateRemoteThread", "VirtualAlloc",
	"curl http", "wget http", "/bin/sh", "reg add", "schtasks", "net user", "Mimikatz",
}

func ScanFile(path string, maxMB int64) (Finding, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return Finding{}, err
	}
	if fi.IsDir() {
		return Finding{}, errors.New("is directory")
	}
	if fi.Size() > maxMB*1024*1024 {
		return Finding{File: path, Size: fi.Size()}, fmt.Errorf("file exceeds max size")
	}

	f, err := os.Open(path)
	if err != nil {
		return Finding{}, err
	}
	defer f.Close()

	h := sha256.New()
	buf := bufio.NewScanner(io.TeeReader(f, h))
	buf.Buffer(make([]byte, 1024), 4*1024*1024)
	found := map[string]struct{}{}
	for buf.Scan() {
		line := strings.ToLower(buf.Text())
		for _, p := range suspiciousPatterns {
			if strings.Contains(line, strings.ToLower(p)) {
				found[p] = struct{}{}
			}
		}
	}

	patterns := make([]string, 0, len(found))
	for k := range found {
		patterns = append(patterns, k)
	}

	return Finding{
		File:       filepath.Clean(path),
		Hash:       hex.EncodeToString(h.Sum(nil)),
		Suspicious: patterns,
		Size:       fi.Size(),
	}, nil
}
