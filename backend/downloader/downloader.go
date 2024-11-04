package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/kkdai/youtube/v2"

	"movie_downloader/backend/config"
)

type VideoInfo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	FilePath    string    `json:"filePath"`
	Status      string    `json:"status"` // pending, downloading, completed, error
	Progress    float64   `json:"progress"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"createdAt"`
	CompletedAt time.Time `json:"completedAt,omitempty"`
	Size        int64     `json:"size"`
	Exists      bool      `json:"exists"`
}

type DownloadManager struct {
	videos          map[string]*VideoInfo // 当前下载任务
	downloadHistory []*VideoInfo          // 下载历史记录
	client          *youtube.Client
	semaphore       chan struct{}
	mu              sync.Mutex
	logger          *log.Logger
	historyFile     string // 添加历史记录文件路径
}

func NewDownloadManager(maxConcurrent int) *DownloadManager {
	// 创建日志文件
	logFile, err := os.OpenFile("download.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	dm := &DownloadManager{
		videos:          make(map[string]*VideoInfo),
		downloadHistory: make([]*VideoInfo, 0),
		client:          &youtube.Client{},
		semaphore:       make(chan struct{}, maxConcurrent),
		logger:          log.New(logFile, "", log.LstdFlags),
		historyFile:     "download_history.json",
	}

	// 加载历史记录
	dm.loadHistory()
	return dm
}

// 加载下载历史
func (dm *DownloadManager) loadHistory() {
	data, err := os.ReadFile(dm.historyFile)
	if err != nil {
		if !os.IsNotExist(err) {
			dm.logger.Printf("读取历史记录失败: %v", err)
		}
		return
	}

	var history []*VideoInfo
	if err := json.Unmarshal(data, &history); err != nil {
		dm.logger.Printf("解析历史记录失败: %v", err)
		return
	}

	for _, video := range history {
		if video.Status == "completed" {
			// 检查文件是否存在，如果不存在则标记为已删除
			exists, _ := checkFileExists(video.FilePath)
			video.Exists = exists
		} else {
			video.Exists = false
		}
	}

	dm.downloadHistory = history
}

func (dm *DownloadManager) addToHistory(video *VideoInfo) {
	dm.mu.Lock()
	dm.downloadHistory = append(dm.downloadHistory, video)
	delete(dm.videos, video.ID)
	dm.mu.Unlock()

	dm.saveHistory()
}

// 保存下载历史
func (dm *DownloadManager) saveHistory() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	data, err := json.MarshalIndent(dm.downloadHistory, "", "  ")
	if err != nil {
		dm.logger.Printf("保存历史记录失败: %v", err)
		return err
	}

	if err := os.WriteFile(dm.historyFile, data, 0644); err != nil {
		dm.logger.Printf("保存历史记录失败: %v", err)
		return err
	}

	return nil
}

func (dm *DownloadManager) Download(videoURL, downloadDir string) error {
	video := &VideoInfo{
		ID:        videoURL,
		Status:    "pending",
		Progress:  0,
		CreatedAt: time.Now(),
	}

	videoID := extractVideoID(videoURL)
	if videoID == "" {
		dm.logger.Printf("无法从URL中提取视频ID: %s", videoURL)
		dm.mu.Lock()
		video.Message = "无法从URL中提取视频ID"
		video.Status = "error"
		dm.mu.Unlock()

		dm.addToHistory(video)
		return nil
	}

	dm.mu.Lock()
	video.ID = videoID

	if video, exists := dm.videos[videoID]; exists {
		if video.Status == "downloading" {
			dm.mu.Unlock()
			return fmt.Errorf("视频已在下载列表中")
		}
	}
	dm.mu.Unlock()

	go func() {
		dm.semaphore <- struct{}{}
		defer func() { <-dm.semaphore }()

		// 记录开始下载日志
		dm.logger.Printf("开始下载视频: %s", videoID)
		video.Status = "downloading"
		video.Message = "获取视频信息..."

		// 添加重试机制
		var videoData *youtube.Video
		var err error

		for retries := 0; retries < config.AppConfig.RetryTime+1; retries++ {
			videoData, err = dm.client.GetVideo(videoID)
			if err == nil {
				break
			}
			video.Message = fmt.Sprintf("重试获取视频信息 (%d/3)...", retries+1)
			time.Sleep(time.Second * 2)
		}

		if err != nil {
			video.Status = "error"
			video.Message = "获取视频信息失败: " + err.Error()
			dm.logger.Printf("获取视频信息失败: %v", err)
			dm.addToHistory(video)
			return
		}

		// 获取最佳质量的格式
		formats := videoData.Formats.WithAudioChannels()
		if len(formats) == 0 {
			video.Status = "error"
			video.Message = "没有可用的视频格式"
			dm.logger.Printf("视频 %s 没有可用的格式", videoID)
			dm.addToHistory(video)
			return
		}

		// 选择最佳质量的格式
		bestFormat := dm.selectBestFormat(formats)
		video.Message = fmt.Sprintf("开始下载视频流... (比特率: %d kbps)", bestFormat.Bitrate/1024)

		// 添加重试机制获取视频流
		var stream io.ReadCloser
		var size int64
		for retries := 0; retries < config.AppConfig.RetryTime+1; retries++ {
			stream, size, err = dm.client.GetStream(videoData, bestFormat)
			if err == nil {
				break
			}
			video.Message = fmt.Sprintf("重试获取视频流 (%d/3)...", retries+1)
			time.Sleep(time.Second * 2)
		}

		if err != nil {
			video.Status = "error"
			video.Message = "获取视频流失败: " + err.Error()
			dm.logger.Printf("获取视频流失败: %v", err)
			dm.addToHistory(video)
			return
		}
		defer stream.Close()

		// 确保下载目录存在
		if err := os.MkdirAll(downloadDir, 0755); err != nil {
			video.Status = "error"
			video.Message = "创建下载目录失败: " + err.Error()
			dm.logger.Printf("创建下载目录失败: %v", err)
			dm.addToHistory(video)
			return
		}

		// 处理文件名中的非法字符
		safeTitle := sanitizeFileName(videoData.Title)
		fileName := safeTitle + ".mp4"
		filePath := filepath.Join(downloadDir, fileName)

		file, err := os.Create(filePath)
		if err != nil {
			video.Status = "error"
			video.Message = "创建文件失败: " + err.Error()
			dm.logger.Printf("创建文件失败: %v", err)
			dm.addToHistory(video)
			return
		}
		defer file.Close()

		// 使用带进度的复制
		video.Message = "正在下载..."
		written, err := dm.copyWithProgress(file, stream, size, video)
		if err != nil {
			video.Status = "error"
			video.Message = "下载过程中出错: " + err.Error()
			dm.logger.Printf("下载过程中出错: %v", err)
			dm.logger.Printf("video: %+v", video)
			dm.addToHistory(video)
			return
		}

		video.Title = videoData.Title
		video.FilePath = filePath
		video.Status = "completed"
		video.Progress = 100
		video.Message = fmt.Sprintf("下载完成 (大小: %.2f MB)", float64(written)/(1024*1024))
		video.Size = written
		dm.logger.Printf("视频下载完成: %s, 大小: %d bytes", videoID, written)
		// 在下载完成或失败时更新历史记录
		video.CompletedAt = time.Now()

		dm.addToHistory(video)
	}()

	return nil
}

// 选择最佳视频格式
func (dm *DownloadManager) selectBestFormat(formats []youtube.Format) *youtube.Format {
	var bestFormat *youtube.Format
	var maxBitrate int

	for i, format := range formats {
		// 检查是否同时包含视频和音频
		hasVideo := format.ItagNo != 0 && format.Quality != ""
		hasAudio := format.AudioChannels > 0

		if hasVideo && hasAudio {
			// 使用比特率作为质量判断标准
			bitrate := format.Bitrate
			if bitrate > maxBitrate {
				maxBitrate = bitrate
				bestFormat = &formats[i]
			}
		}
	}

	// 如果没有找到同时包含视频和音频的格式，使用第一个格式
	if bestFormat == nil && len(formats) > 0 {
		// 尝试找到最高质量的格式
		maxBitrate = 0
		for i, format := range formats {
			if format.Bitrate > maxBitrate {
				maxBitrate = format.Bitrate
				bestFormat = &formats[i]
			}
		}
	}

	// 如果还是没找到，使用第一个格式
	if bestFormat == nil && len(formats) > 0 {
		bestFormat = &formats[0]
	}

	return bestFormat
}

// 处理文件名中的非法字符
func sanitizeFileName(fileName string) string {
	// 替换 Windows 文件名中的非法字符
	illegal := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	result := fileName

	for _, char := range illegal {
		result = strings.ReplaceAll(result, char, "_")
	}

	// 限制文件名长度
	if len(result) > 200 {
		result = result[:200]
	}

	return result
}

func (dm *DownloadManager) copyWithProgress(dst io.Writer, src io.Reader, size int64, video *VideoInfo) (int64, error) {
	buf := make([]byte, 32*1024)
	var written int64

	for {
		nr, err := src.Read(buf)
		if nr > 0 {
			nw, err := dst.Write(buf[0:nr])
			if err != nil {
				return written, err
			}
			if nr != nw {
				return written, io.ErrShortWrite
			}
			written += int64(nw)

			// 更新进度
			if size > 0 {
				video.Progress = float64(written) / float64(size) * 100
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return written, err
		}
	}
	return written, nil
}

func (dm *DownloadManager) GetVideoList() []*VideoInfo {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	videos := make([]*VideoInfo, 0, len(dm.videos))
	for _, v := range dm.videos {
		videos = append(videos, v)
	}
	return videos
}

func (dm *DownloadManager) DeleteVideo(videoID string) error {
	// dm.mu.Lock()
	// defer dm.mu.Unlock()

	// 先检查当前下载列表
	if video, exists := dm.videos[videoID]; exists {
		if video.Status == "downloading" {
			return fmt.Errorf("无法删除正在下载的视频")
		}
		delete(dm.videos, videoID)
		if video.FilePath != "" && video.Exists {
			if err := os.Remove(video.FilePath); err != nil {
				return err
			}
		}
		return nil
	}

	dm.loadHistory()

	// 检查历史记录
	for i, video := range dm.downloadHistory {
		if video.ID == videoID {
			// 从历史记录中移除
			dm.downloadHistory = append(dm.downloadHistory[:i], dm.downloadHistory[i+1:]...)
			if video.FilePath != "" && video.Exists {
				if err := os.Remove(video.FilePath); err != nil {
					return err
				}
			}
			// 保存更新后的历史记录
			return dm.saveHistory()
		}
	}

	return nil
}

func extractVideoID(url string) string {
	// 支持的 URL 格式：
	// - https://www.youtube.com/watch?v=VIDEO_ID
	// - https://youtu.be/VIDEO_ID
	// - https://www.youtube.com/embed/VIDEO_ID
	// - https://youtube.com/shorts/VIDEO_ID

	// 使用正则表达式匹配视频ID
	patterns := []string{
		`(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/embed/|youtube\.com/shorts/)([^&?/\s]{11})`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	// 如果输入的已经是11位的视频ID，直接返回
	if len(url) == 11 {
		return url
	}

	return ""
}

// 获取下载历史记录
func (dm *DownloadManager) GetDownloadHistory() []*VideoInfo {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.loadHistory()
	return dm.downloadHistory
}

// 重新下载视频
func (dm *DownloadManager) RedownloadVideo(videoID string) error {
	dm.mu.Lock()
	// 检查是否已经在下载中
	if _, exists := dm.videos[videoID]; exists {
		dm.mu.Unlock()
		return fmt.Errorf("视频已在下载列表中")
	}
	dm.mu.Unlock()

	// 创建新的下载任务
	return dm.Download("https://www.youtube.com/watch?v="+videoID, config.AppConfig.DownloadDir)
}

// 添加清理历史记录的方法
func (dm *DownloadManager) ClearHistory() error {
	dm.mu.Lock()
	dm.downloadHistory = make([]*VideoInfo, 0)
	dm.mu.Unlock()

	return dm.saveHistory()
}

// 添加导出历史记录的方法
func (dm *DownloadManager) ExportHistory(filePath string) error {
	dm.mu.Lock()
	data, err := json.MarshalIndent(dm.downloadHistory, "", "  ")
	dm.mu.Unlock()

	if err != nil {
		return fmt.Errorf("序列化历史记录失败: %v", err)
	}

	return os.WriteFile(filePath, data, 0644)
}

// 添加导入历史记录的方法
func (dm *DownloadManager) ImportHistory(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取历史记录文件失败: %v", err)
	}

	var history []*VideoInfo
	if err := json.Unmarshal(data, &history); err != nil {
		return fmt.Errorf("解析历史记录失败: %v", err)
	}

	dm.mu.Lock()
	dm.downloadHistory = history
	dm.mu.Unlock()

	return dm.saveHistory()
}

func checkFileExists(filePath string) (bool, error) {
	if filePath == "" {
		return false, fmt.Errorf("文件路径为空")
	}

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
