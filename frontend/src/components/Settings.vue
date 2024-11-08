<template>
  <div class="settings">
    <div class="form-group">
      <label for="downloadDir">下载目录：</label>
      <input id="downloadDir" v-model="config.downloadDir" @click="selectDownloadDirectory" type="text" readonly>
      <!-- <button @click="selectDownloadDirectory"  aria-label="选择下载目录">选择目录</button> -->
    </div>
    <div class="form-group">
      <label for="maxConcurrent">同时下载数量：</label>
      <input id="maxConcurrent" v-model.number="config.maxConcurrent" type="number" min="1" max="10">
    </div>
    <div class="form-group">
      <label for="retryTime">重试次数：</label>
      <input id="retryTime" v-model="config.retryTime" type="number" min="0" placeholder="失败重试次数">
    </div>
    <button class="save-button" @click="saveConfig">保存设置</button>
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
:root {
  --primary-color: #4CAF50;
  --secondary-color: #f1f1f1;
  --font-color: #333;
}

.settings {
  padding: 20px;
  background-color: var(--secondary-color);
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.form-group {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.form-group label {
  flex: 1;
  font-size: 16px;
  color: var(--font-color);
  margin-right: 10px;
}

.form-group input {
  flex: 2;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
}

.form-group button {
  flex: 1;
  padding: 8px 16px;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.form-group button:hover {
  background-color: #45a049;
}

.save-button {
  width: 100%;
  padding: 10px;
  margin: 0;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.save-button:hover {
  background-color: #45a049;
}

#downloadDir {
  width: 100%;
  padding: 10px;
  margin: 0;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

@media (max-width: 768px) {
  .settings {
    padding: 15px;
  }

  .form-group {
    flex-direction: column;
    align-items: flex-start;
  }

  .form-group label {
    margin-bottom: 5px;
  }

  .form-group input,
  .form-group button {
    width: 100%;
    margin-bottom: 10px;
  }

  .save-button {
    margin-top: 15px;
  }
}
</style> 