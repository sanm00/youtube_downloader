package main

import (
	"context"
	"path/filepath"

	"movie_downloader/backend/config"
	"movie_downloader/backend/downloader"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
	dm  *downloader.DownloadManager
}

func NewApp() *App {
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}

	return &App{
		dm: downloader.NewDownloadManager(config.AppConfig.MaxConcurrent),
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Download(url string) error {
	return a.dm.Download(url, config.AppConfig.DownloadDir)
}

func (a *App) GetVideoList() []*downloader.VideoInfo {
	return a.dm.GetVideoList()
}

func (a *App) DeleteVideo(id string) error {
	return a.dm.DeleteVideo(id)
}

func (a *App) SaveConfig(downloadDir string, maxConcurrent int, retryTime int) error {
	err := config.SaveConfig(downloadDir, maxConcurrent, retryTime)
	if err != nil {
		return err
	}

	a.dm = downloader.NewDownloadManager(maxConcurrent)
	return nil
}

func (a *App) GetDownloadHistory() []*downloader.VideoInfo {
	return a.dm.GetDownloadHistory()
}

func (a *App) RedownloadVideo(id string) error {
	return a.dm.RedownloadVideo(id)
}

func (a *App) GetConfig() *config.Config {
	return config.AppConfig
}

func (a *App) ExportHistory(filePath string) error {
	return a.dm.ExportHistory(filePath)
}

func (a *App) ImportHistory(filePath string) error {
	return a.dm.ImportHistory(filePath)
}

func (a *App) ClearHistory() error {
	return a.dm.ClearHistory()
}

func (a *App) OpenDirectory(filePath string) {
	dir := filepath.Dir(filePath)
	dirURL := "file://" + dir
	runtime.BrowserOpenURL(a.ctx, dirURL)
}

func (a *App) SelectDownloadDirectory() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择下载目录",
	})
	if err != nil {
		return "", err
	}
	return dir, nil
}
