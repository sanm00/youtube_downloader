<template>
  <div class="settings">
    <div class="form-group">
      <label>下载目录：</label>
      <input v-model="config.downloadDir" type="text" readonly>
      <button @click="selectDownloadDirectory">选择目录</button>
    </div>
    <div class="form-group">
      <label>同时下载数量：</label>
      <input v-model.number="config.maxConcurrent" type="number" min="1" max="10">
    </div>
    <div class="form-group">
      <label>重试次数：</label>
      <input v-model="config.retryTime" type="number" min="0" placeholder="失败重试次数">
    </div>
    <button @click="saveConfig">保存设置</button>
  </div>
</template>

<script>
// 尝试使用完整路径
import {SaveConfig, GetConfig, SelectDownloadDirectory} from '../../wailsjs/go/main/App';

export default {
  data() {
    return {
      config: {
        downloadDir: '',
        maxConcurrent: 3,
        retryTime: 0
      }
    }
  },
  methods: {
    async saveConfig() {
      try {
        await SaveConfig(this.config.downloadDir, this.config.maxConcurrent, this.config.retryTime)
      } catch (error) {
        console.error('保存配置失败:', error)
      }
    },
    async loadConfig() {
      try {
        const currentConfig = await GetConfig()
        this.config.downloadDir = currentConfig.downloadDir
        this.config.maxConcurrent = currentConfig.maxConcurrent
        this.config.retryTime = currentConfig.retryTime
      } catch (error) {
        console.error('获取配置失败:', error)
      }
    },
    async selectDownloadDirectory() {
      try {
        const dir = await SelectDownloadDirectory()
        if (dir) {
          this.config.downloadDir = dir
        }
      } catch (error) {
        console.error('选择目录失败:', error)
      }
    }
  },
  async mounted() {
    await this.loadConfig()
  }
}
</script>

<style>
.settings {
  padding: 20px;
}

.form-group {
  margin: 10px 0;
}
</style> 