package config

import (
	"os"
	"sync"
	"time"
)

type ConfigWatcher struct {
	configPath string
	onChange   func(*Config)
	mu         sync.RWMutex
	config     *Config
	stopCh     chan struct{}
}

func NewConfigWatcher(configPath string, onChange func(*Config)) *ConfigWatcher {
	return &ConfigWatcher{
		configPath: configPath,
		onChange:   onChange,
		stopCh:     make(chan struct{}),
	}
}

func (w *ConfigWatcher) Start() error {
	config, err := Load(w.configPath)
	if err != nil {
		return err
	}

	w.mu.Lock()
	w.config = config
	w.mu.Unlock()

	go w.watch()
	return nil
}

func (w *ConfigWatcher) Stop() {
	close(w.stopCh)
}

func (w *ConfigWatcher) GetConfig() *Config {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.config
}

func (w *ConfigWatcher) watch() {
	var lastModTime time.Time

	info, err := os.Stat(w.configPath)
	if err == nil {
		lastModTime = info.ModTime()
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopCh:
			return
		case <-ticker.C:
			info, err := os.Stat(w.configPath)
			if err != nil {
				continue
			}

			if info.ModTime().After(lastModTime) {
				lastModTime = info.ModTime()

				newConfig, err := Load(w.configPath)
				if err != nil {
					continue
				}

				w.mu.Lock()
				w.config = newConfig
				w.mu.Unlock()

				if w.onChange != nil {
					w.onChange(newConfig)
				}
			}
		}
	}
}
