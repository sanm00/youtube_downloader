<template>
  <div id="app">
    <header>
      <nav>
        <router-link to="/videos" :class="{ active: isActive('/videos') }">下载视频</router-link>
        <router-link to="/settings" :class="{ active: isActive('/settings') }">下载设置</router-link>
        <router-link to="/download-history" :class="{ active: isActive('/download-history') }">下载历史</router-link>
      </nav>
    </header>
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  name: 'App',
  methods: {
    isActive(route) {
      return this.$route.path === route;
    }
  }
}
</script>

<style>
header {
  background-color: #f8f9fa;
  padding: 10px;
  border-bottom: 1px solid #ccc;
}

nav {
  display: flex;
  justify-content: space-around;
}

nav a {
  text-decoration: none;
  color: #333;
  font-weight: bold;
}

nav a:hover {
  color: #007bff;
}

nav a.active {
  color: #007bff;
  font-weight: bold;
}
</style> 