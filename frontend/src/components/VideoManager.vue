<template>
  <div class="video-manager">
    <div class="download-section">
      <div class="form-group">
        <textarea v-model="videoUrls" placeholder="输入多个视频链接，每行一个"></textarea>
      </div>

      <button class="margin-0 save-button" @click="downloadVideos">新建下载任务</button>
    </div>

    <div class="video-list">
      <h4 class="title-4">下载队列</h4>
      <DownloadProgress v-for="video in activeDownloads" :key="video.id" :video="video" />
    </div>
  </div>
</template>

<script>
import { GetVideoList, Download, DeleteVideo, OpenDirectory } from '../../wailsjs/go/main/App'
import DownloadProgress from './DownloadProgress.vue'

export default {
  components: {
    DownloadProgress
  },
  data() {
    return {
      videoUrls: '',
      activeDownloads: []
    }
  },
  methods: {
    async downloadVideos() {
      if (!this.videoUrls) return
      const urls = this.videoUrls.split('\n').map(url => url.trim());
      for (const url of urls) {
        if (url) {
          try {
            await Download(url);
          } catch (error) {
            console.error('下载视频失败:', error);
          }
        }
      }
      this.videoUrls = '';
      this.updateLists();
    },
    async deleteVideo(id) {
      await DeleteVideo(id);
      this.updateLists();
    },
    async updateLists() {
      this.activeDownloads = await GetVideoList();
    },
    async openDirectory(filePath) {
      try {
        await OpenDirectory(filePath);
      } catch (error) {
        console.error('打开目录失败:', error);
      }
    }
  },
  async mounted() {
    this.updateLists();
    setInterval(this.updateLists, 1000);
  }
}
</script>

<style>
.video-manager {
  padding: 20px;
}

.form-group {
  margin: 10px 0;
}

.margin-0 {
  margin: 0;
}

.title-4 {
  border-bottom: 1px solid #ccc;
}

textarea {
  width: 100%;
  height: 100px; /* 设置文本框高度 */
  padding: 10px;
  border-radius: 4px;
  border: 1px solid #ccc;
  resize: none; /* 禁止调整大小 */
  box-sizing: border-box;
}
</style> 