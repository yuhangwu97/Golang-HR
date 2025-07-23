<template>
  <div class="test-view">
    <h1>Vue 2 应用运行正常！</h1>
    <p>当前时间: {{ currentTime }}</p>
    <el-button type="primary" @click="showMessage">点击测试</el-button>
  </div>
</template>

<script>
export default {
  name: 'TestView',
  data() {
    return {
      currentTime: new Date().toLocaleString()
    }
  },
  methods: {
    showMessage() {
      this.$message.success('Vue 2 + Element UI 工作正常！')
    }
  }
}
</script>

<style scoped>
.test-view {
  padding: 50px;
  text-align: center;
}

h1 {
  color: #409EFF;
  margin-bottom: 20px;
}
</style>