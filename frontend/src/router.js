import { createRouter, createWebHashHistory } from 'vue-router'
import Settings from './components/Settings.vue'
import VideoManager from './components/VideoManager.vue'
import History from './components/DownloadHistory.vue'

const routes = [
  { path: '/', redirect: '/videos' },
  { path: '/settings', component: Settings },
  { path: '/videos', component: VideoManager },
  { path: '/download-history', component: History }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router 