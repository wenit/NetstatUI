package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/zwb/network-ports/services/kill"
	"github.com/zwb/network-ports/services/process"
	"github.com/zwb/network-ports/services/system"
)

type AppService struct {
	cache *process.Cache
}

func NewAppService(cache *process.Cache) *AppService {
	return &AppService{cache: cache}
}

func (a *AppService) ServiceName() string { return "app" }

func (a *AppService) KillProcess(pid uint32) kill.Result {
	return kill.Process(pid)
}

func (a *AppService) GetProcessDetail(pid uint32) (process.Info, error) {
	info, err := a.cache.Get(pid)
	if err != nil {
		return process.Info{}, err
	}
	return info, nil
}

func (a *AppService) OpenProcessFolder(pid uint32) error {
	info, err := a.cache.Get(pid)
	if err != nil {
		return err
	}
	if info.Path == "" {
		return fmt.Errorf("process path unavailable (pid=%d)", pid)
	}
	dir := filepath.Dir(info.Path)
	switch runtime.GOOS {
	case "windows":
		c := exec.Command("explorer.exe", "/select,", info.Path)
		return c.Start()
	default:
		c := exec.Command("xdg-open", dir)
		return c.Start()
	}
}

func (a *AppService) GetSystemLocale() string {
	return system.GetSystemLocale()
}
