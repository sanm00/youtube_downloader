<template>
  <div class="video-item" :class="itemClass(video)">
    <div class="video-info">
      <div class="">
        <span class="video-id">{{ video.id }}</span>
        <span class="video-title">{{ video.title }}</span>
      </div>
      <div class="video-status">
        <span class="status-text">{{ video.status }}</span>
        <span class="download-time">{{ formatTime(video.downloadTime) }}</span>
        <!-- <span class="file-size">{{ formatFileSize(video.size) }}</span> -->
        <span class="message">{{ video.message }}</span>
      </div>
      <div v-if="video.status === 'downloading'" class="progress-bar">
        <div class="progress" :style="{ width: video.progress + '%' }"></div>
      </div>
    </div>
    <div class="actions">
      <button v-if="video.status === 'completed' && video.Exists" @click="openDirectory">打开文件夹</button>
      <button @click="deleteVideo">删除</button>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    video: {
      type: Object,
      required: true
    }
  },
  methods: {
    openDirectory() {
      this.$emit('open-directory', this.video.filePath); // 触发事件，传递文件路径
    },
    deleteVideo() {
      this.$emit('delete-video', this.video.id); // 触发事件，传递视频ID
    },
    formatTime(timestamp) {
      if (!timestamp) return '';
      return new Date(timestamp).toLocaleString();
    },
    formatFileSize(size) {
      if (!size) return '未知大小';
      const units = ['B', 'KB', 'MB', 'GB'];
      let index = 0;
      let formattedSize = size;
      while (formattedSize >= 1024 && index < units.length - 1) {
        formattedSize /= 1024;
        index++;
      }
      return `${formattedSize.toFixed(2)} ${units[index]}`;
    },
    itemClass(video) {
      if (!video.exists && video.status == 'completed') {
        return ''
      }

      if (video.status === 'completed') {
        return 'highlight';
      } else if (video.status === 'error') {
        return 'error';
      }

      return '';
    }
  },
}
</script>

<style>
.video-item {
  background: #f5f5f5;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 10px;
  display: flex;
  justify-content: space-between; /* 左右分布 */
  align-items: center; /* 垂直居中 */
}

.video-item.highlight {
  background: #d4edda; /* 高亮背景颜色 */
  border-left: 5px solid #28a745; /* 左侧边框颜色 */
}
.video-item.error {
  background: #f3d7d8;
  border-left: 5px solid #dc3545; /* 左侧边框颜色 */
}

.video-info {
  flex: 1; /* 占据剩余空间 */
}

.video-id {
  margin-right: 5px;
  font-weight: bold;
  margin-bottom: 8px;
}

.video-status {
  display: flex;
  align-items: center;
  margin-bottom: 5px;
}

.status-text {
  margin-right: 10px;
  min-width: 80px;
}

.download-time {
  margin-right: 10px;
}

.file-size {
  margin-right: 10px;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: #ddd;
  border-radius: 4px;
  overflow: hidden;
  margin: 0 10px;
}

.progress {
  height: 100%;
  background: #4CAF50;
  transition: width 0.3s ease;
}

.actions {
  display: flex;
  align-items: center;
}

button {
  margin-left: 10px; /* 按钮之间的间距 */
}
</style> 