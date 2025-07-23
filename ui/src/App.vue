<template>
  <div id="app">
    <transition name="page-fade" mode="out-in">
      <router-view />
    </transition>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
/* 现代化的CSS变量定义 */
:root {
  --primary-color: #1890ff;
  --primary-color-hover: #40a9ff;
  --primary-color-active: #096dd9;
  --success-color: #52c41a;
  --warning-color: #faad14;
  --error-color: #f5222d;
  --text-color: #262626;
  --text-color-secondary: #8c8c8c;
  --text-color-disabled: #bfbfbf;
  --background-color: #f0f2f5;
  --component-background: #ffffff;
  --border-color: #d9d9d9;
  --border-radius: 8px;
  --box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  --box-shadow-hover: 0 4px 16px rgba(0, 0, 0, 0.12);
  --transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 全局样式重置 */
* {
  box-sizing: border-box;
}

body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', 'Helvetica Neue', Helvetica, Arial, sans-serif;
  font-size: 14px;
  line-height: 1.5715;
  color: var(--text-color);
  background-color: var(--background-color);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  min-height: 100vh;
  position: relative;
}

/* 页面过渡动画 */
.page-fade-enter-active, .page-fade-leave-active {
  transition: opacity 0.3s ease-in-out;
}

.page-fade-enter, .page-fade-leave-to {
  opacity: 0;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
  transition: background 0.3s;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* 通用动画类 */
.fade-in {
  animation: fadeIn 0.3s ease-in-out;
}

.slide-up {
  animation: slideUp 0.3s ease-out;
}

.zoom-in {
  animation: zoomIn 0.3s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes zoomIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

/* 卡片悬停效果 */
.card-hover {
  transition: var(--transition);
  cursor: pointer;
}

.card-hover:hover {
  box-shadow: var(--box-shadow-hover);
  transform: translateY(-2px);
}

/* 按钮动画增强 */
.btn-animate {
  position: relative;
  overflow: hidden;
  transition: var(--transition);
}

.btn-animate::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.btn-animate:hover::before {
  left: 100%;
}

/* 响应式设计 */
@media (max-width: 768px) {
  body {
    font-size: 12px;
  }
}

/* 主题色调整 */
.ant-btn-primary {
  background: linear-gradient(135deg, var(--primary-color), var(--primary-color-hover));
  border: none;
  box-shadow: 0 2px 4px rgba(24, 144, 255, 0.2);
  transition: var(--transition);
}

.ant-btn-primary:hover {
  background: linear-gradient(135deg, var(--primary-color-hover), var(--primary-color));
  box-shadow: 0 4px 8px rgba(24, 144, 255, 0.3);
  transform: translateY(-1px);
}

/* 表格样式增强 */
.ant-table {
  background: var(--component-background);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
}

.ant-table-thead > tr > th {
  background: #fafafa;
  border-bottom: 2px solid var(--border-color);
  font-weight: 600;
}

.ant-table-tbody > tr:hover > td {
  background: #f5f5f5;
}

/* 卡片样式增强 */
.ant-card {
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  border: none;
  transition: var(--transition);
}

.ant-card:hover {
  box-shadow: var(--box-shadow-hover);
}

/* 输入框样式增强 */
.ant-input, .ant-select-selector {
  border-radius: var(--border-radius);
  transition: var(--transition);
}

.ant-input:focus, .ant-select-selector:focus {
  box-shadow: 0 0 0 3px rgba(24, 144, 255, 0.1);
}
</style>