<template>
  <div class="download-history">
    <div class="video-list">
      <DownloadProgress 
        v-for="video in downloadHistory" 
        :key="video.id" 
        :video="video" 
        @open-directory="openDirectory" 
        @delete-video="deleteVideo" 
      />
    </div>
  </div>
</template>

<script>
import { GetDownloadHistory, OpenDirectory, DeleteVideo } from '../../wailsjs/go/main/App'
import DownloadProgress from './DownloadProgress.vue'

export default {
  components: {
    DownloadProgress
  },
  data() {
    return {
      downloadHistory: []
    }
  },
  async mounted() {
    console.log('DownloadHistory mounted')
    this.downloadHistory = await GetDownloadHistory();
  },
  methods: {
    async openDirectory(filePath) {
      try {
        await OpenDirectory(filePath);
      } catch (error) {
        console.error('打开目录失败:', error);
      }
    },
    async deleteVideo(id) {
      try {
        console.log('Deleting video with ID:', id);
        await DeleteVideo(id); // 调用删除视频的API
        this.downloadHistory = this.downloadHistory.filter(video => video.id !== id); // 更新下载历史
      } catch (error) {
        console.error('删除视频失败:', error);
      }
    }
  }
}
</script>

<style>
.download-history {
  padding: 20px;
}
</style> 