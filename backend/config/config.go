package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DownloadDir   string `json:"downloadDir"`
	MaxConcurrent int    `json:"maxConcurrent"`
	RetryTime     int    `json:"retryTime"`
}

var AppConfig = &Config{
	DownloadDir:   "~/.youtube-downloader/downloads",
	MaxConcurrent: 3,
	RetryTime:     0,
}

const configFilePath = "config.json"

// 保存配置到文件
func SaveConfig(downloadDir string, maxConcurrent int, retryTime int) error {
	AppConfig.DownloadDir = downloadDir
	AppConfig.MaxConcurrent = maxConcurrent
	AppConfig.RetryTime = retryTime

	data, err := json.MarshalIndent(AppConfig, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFilePath, data, 0644)
}

// 从文件加载配置
func LoadConfig() error {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return SaveConfig(AppConfig.DownloadDir, AppConfig.MaxConcurrent, AppConfig.RetryTime)
		}
		return err
	}

	return json.Unmarshal(data, AppConfig)
}
