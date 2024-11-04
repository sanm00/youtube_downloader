# youtube downloader

## 项目功能

该项目是一个基于Wails的桌面应用，前端使用Vue.js，主要功能包括：

- **视频下载**：更具YouTube web 链接下载视频。
- **下载进度管理**：实时显示视频下载的进度。
- **视频列表管理**：用户可以查看已上传的视频列表，并进行管理操作。
- **历史记录**：记录用户的下载历史，方便用户查看。

## 使用方法

### 环境要求

- Go >= 1.16
- Node.js >= 12.x
- Vue.js >= 2.x
- Wails >= 1.x

### 安装步骤

1. 克隆项目到本地：
   ```bash
   git@github.com:sanm00/youtube_downloader.git
   cd youtube_downloader
   ```

2. 安装前端依赖：
   ```bash
   cd frontend
   npm install
   ```

3. 安装Wails：
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

### 启动项目

1. 启动开发模式：
   ```bash
   wails dev
   ```

2. 构建项目：
   ```bash
   wails build
   ```

### 访问项目

在开发模式下，打开 `http://localhost:34115` 。

## 贡献

欢迎任何形式的贡献！请提交问题或拉取请求。

## 许可证

该项目使用 MIT 许可证，详细信息请查看 [LICENSE](frontend/node_modules/magic-string/LICENSE) 文件。