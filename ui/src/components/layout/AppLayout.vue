<template>
  <a-layout class="app-layout">
    <a-layout-sider
      v-model:collapsed="collapsed"
      :trigger="null"
      collapsible
      class="custom-sider"
      :width="256"
      :collapsed-width="80"
    >
      <div class="logo" :class="{ collapsed: collapsed }">
        <div class="logo-icon">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 12C14.21 12 16 10.21 16 8C16 5.79 14.21 4 12 4C9.79 4 8 5.79 8 8C8 10.21 9.79 12 12 12ZM12 14C9.33 14 4 15.34 4 18V20H20V18C20 15.34 14.67 14 12 14Z" fill="currentColor"/>
          </svg>
        </div>
        <transition name="fade-slide" mode="out-in">
          <div v-if="!collapsed" class="logo-text">
            <h2>HR管理系统</h2>
            <span>Human Resource</span>
          </div>
          <div v-else class="logo-text-mini">
            <h2>HR</h2>
          </div>
        </transition>
      </div>
      
      <a-menu
        v-model:selectedKeys="selectedKeys"
        class="custom-menu"
        mode="inline"
        @click="handleMenuClick"
      >
        <a-menu-item key="/dashboard" class="menu-item">
          <div class="menu-item-content">
            <div class="menu-icon">
              <DashboardOutline />
            </div>
            <span class="menu-text">仪表盘</span>
          </div>
        </a-menu-item>
        
        <a-menu-item key="/employees" class="menu-item">
          <div class="menu-item-content">
            <div class="menu-icon">
              <TeamOutline />
            </div>
            <span class="menu-text">员工管理</span>
          </div>
        </a-menu-item>
        
        <a-sub-menu key="organization" class="menu-submenu">
          <template #title>
            <div class="menu-item-content">
              <div class="menu-icon">
                <ApartmentOutline />
              </div>
              <span class="menu-text">组织管理</span>
            </div>
          </template>
          <a-menu-item key="/organization/chart" class="submenu-item">
            <div class="submenu-item-content">
              <div class="submenu-dot"></div>
              <span>组织架构</span>
            </div>
          </a-menu-item>
          <a-menu-item key="/organization/departments" class="submenu-item">
            <div class="submenu-item-content">
              <div class="submenu-dot"></div>
              <span>部门管理</span>
            </div>
          </a-menu-item>
          <a-menu-item key="/organization/positions" class="submenu-item">
            <div class="submenu-item-content">
              <div class="submenu-dot"></div>
              <span>职位管理</span>
            </div>
          </a-menu-item>
          <a-menu-item key="/organization/job-levels" class="submenu-item">
            <div class="submenu-item-content">
              <div class="submenu-dot"></div>
              <span>职级管理</span>
            </div>
          </a-menu-item>
        </a-sub-menu>
        
        <a-sub-menu key="salary" class="menu-submenu">
          <template #title>
            <div class="menu-item-content">
              <div class="menu-icon">
                <DollarOutline />
              </div>
              <span class="menu-text">薪资管理</span>
            </div>
          </template>
          <a-menu-item key="/salary" class="submenu-item">
            <div class="submenu-item-content">
              <div class="submenu-dot"></div>
              <span>薪资列表</span>
            </div>
          </a-menu-item>
          <a-menu-item key="/my-salary" class="submenu-item">
            <div class="submenu-item-content">
              <div class="submenu-dot"></div>
              <span>我的薪资</span>
            </div>
          </a-menu-item>
        </a-sub-menu>
        
        <a-sub-menu key="system" v-if="hasSystemAccess" class="menu-submenu">
          <template #title>
            <div class="menu-item-content">
              <div class="menu-icon">
                <SettingOutline />
              </div>
              <span class="menu-text">系统管理</span>
            </div>
          </template>
          <a-menu-item key="/system/users" class="submenu-item">
            <div class="submenu-item-content">
              <div class="submenu-dot"></div>
              <span>用户管理</span>
            </div>
          </a-menu-item>
          <a-menu-item key="/system/roles" class="submenu-item">
            <div class="submenu-item-content">
              <div class="submenu-dot"></div>
              <span>角色管理</span>
            </div>
          </a-menu-item>
          <a-menu-item key="/system/permissions" class="submenu-item">
            <div class="submenu-item-content">
              <div class="submenu-dot"></div>
              <span>权限管理</span>
            </div>
          </a-menu-item>
        </a-sub-menu>
        
        <a-menu-item key="/profile" class="menu-item">
          <div class="menu-item-content">
            <div class="menu-icon">
              <UserOutline />
            </div>
            <span class="menu-text">个人资料</span>
          </div>
        </a-menu-item>
      </a-menu>
    </a-layout-sider>
    
    <a-layout>
      <a-layout-header class="app-header">
        <div class="header-left">
          <div class="trigger-container">
            <transition name="flip" mode="out-in">
              <MenuUnfoldOutlined
                v-if="collapsed"
                key="unfold"
                class="trigger"
                @click="toggleSidebar"
              />
              <MenuFoldOutlined
                v-else
                key="fold"
                class="trigger"
                @click="toggleSidebar"
              />
            </transition>
          </div>
          
          <div class="breadcrumb-container">
            <a-breadcrumb class="header-breadcrumb">
              <a-breadcrumb-item>
                <HomeOutlined />
              </a-breadcrumb-item>
              <a-breadcrumb-item v-for="(crumb, index) in breadcrumbs" :key="index">
                <span v-if="index === breadcrumbs.length - 1" class="current-page">{{ crumb.name }}</span>
                <a v-else @click="$router.push(crumb.path)">{{ crumb.name }}</a>
              </a-breadcrumb-item>
            </a-breadcrumb>
          </div>
        </div>
        
        <div class="header-right">
          <div class="header-actions">
            <a-tooltip title="搜索">
              <div class="action-item search-btn" @click="showSearch = !showSearch">
                <SearchOutlined />
              </div>
            </a-tooltip>
            
            <a-tooltip title="通知">
              <a-badge :count="notificationCount" :offset="[2, -2]">
                <div class="action-item">
                  <BellOutlined />
                </div>
              </a-badge>
            </a-tooltip>
            
            <a-tooltip title="设置">
              <div class="action-item">
                <SettingOutlined />
              </div>
            </a-tooltip>
            
            <div class="user-info">
              <a-dropdown>
                <div class="user-dropdown-trigger">
                  <a-avatar class="user-avatar" :src="user?.avatar">
                    {{ user?.name?.charAt(0)?.toUpperCase() }}
                  </a-avatar>
                  <div class="user-details">
                    <div class="user-name">{{ user?.name || '未知用户' }}</div>
                    <div class="user-role">{{ getUserRole() }}</div>
                  </div>
                  <DownOutlined class="dropdown-arrow" />
                </div>
                
                <template #overlay>
                  <a-menu class="user-menu">
                    <a-menu-item @click="handleProfile">
                      <UserOutlined />
                      个人资料
                    </a-menu-item>
                    <a-menu-item>
                      <SettingOutlined />
                      账户设置
                    </a-menu-item>
                    <a-menu-divider />
                    <a-menu-item @click="handleLogout" class="logout-item">
                      <LogoutOutlined />
                      退出登录
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>
          </div>
        </div>
        
        <!-- 搜索框 -->
        <transition name="slide-down">
          <div v-if="showSearch" class="global-search">
            <a-input-search
              v-model:value="searchValue"
              placeholder="搜索菜单、员工、部门..."
              style="width: 400px"
              @search="handleGlobalSearch"
              allow-clear
            />
          </div>
        </transition>
      </a-layout-header>
      
      <a-layout-content class="app-content">
        <div class="content-wrapper">
          <transition name="fade-slide" mode="out-in">
            <router-view />
          </transition>
        </div>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script>
import { message } from 'ant-design-vue'
import {
  MenuUnfoldOutlined,
  MenuFoldOutlined,
  DashboardOutline,
  TeamOutline,
  UserOutlined,
  LogoutOutlined,
  ApartmentOutline,
  DollarOutline,
  SettingOutlined,
  HomeOutlined,
  SearchOutlined,
  BellOutlined,
  DownOutlined
} from '@ant-design/icons-vue'
import { mapGetters, mapActions } from 'vuex'

export default {
  name: 'AppLayout',
  components: {
    MenuUnfoldOutlined,
    MenuFoldOutlined,
    DashboardOutline,
    TeamOutline,
    UserOutlined,
    LogoutOutlined,
    ApartmentOutline,
    DollarOutline,
    SettingOutlined,
    HomeOutlined,
    SearchOutlined,
    BellOutlined,
    DownOutlined
  },
  data() {
    return {
      collapsed: false,
      selectedKeys: [this.$route.path],
      showSearch: false,
      searchValue: '',
      notificationCount: 3
    }
  },
  computed: {
    ...mapGetters('auth', ['user', 'isAuthenticated']),
    hasSystemAccess() {
      return this.user?.role === 'admin' || this.user?.roles?.some((role) => role.code === 'admin' || role.code === 'hr')
    },
    breadcrumbs() {
      const path = this.$route.path
      const breadcrumbs = []
      
      if (path === '/dashboard') {
        breadcrumbs.push({ name: '仪表盘', path: '/dashboard' })
      } else if (path === '/employees') {
        breadcrumbs.push({ name: '员工管理', path: '/employees' })
      } else if (path.startsWith('/organization')) {
        breadcrumbs.push({ name: '组织管理', path: '/organization' })
        if (path.includes('departments')) {
          breadcrumbs.push({ name: '部门管理', path: '/organization/departments' })
        } else if (path.includes('positions')) {
          breadcrumbs.push({ name: '职位管理', path: '/organization/positions' })
        } else if (path.includes('job-levels')) {
          breadcrumbs.push({ name: '职级管理', path: '/organization/job-levels' })
        } else if (path.includes('chart')) {
          breadcrumbs.push({ name: '组织架构', path: '/organization/chart' })
        }
      } else if (path.startsWith('/salary')) {
        breadcrumbs.push({ name: '薪资管理', path: '/salary' })
      } else if (path.startsWith('/system')) {
        breadcrumbs.push({ name: '系统管理', path: '/system' })
      } else if (path === '/profile') {
        breadcrumbs.push({ name: '个人资料', path: '/profile' })
      }
      
      return breadcrumbs
    }
  },
  watch: {
    '$route.path'(newPath) {
      this.selectedKeys = [newPath]
    }
  },
  methods: {
    ...mapActions('auth', ['logout']),
    handleMenuClick({ key }) {
      this.$router.push(key)
    },
    handleProfile() {
      this.$router.push('/profile')
    },
    handleLogout() {
      this.logout()
      message.success('退出登录成功')
      this.$router.push('/login')
    },
    toggleSidebar() {
      this.collapsed = !this.collapsed
    },
    handleGlobalSearch(value) {
      console.log('搜索:', value)
      this.showSearch = false
    },
    getUserRole() {
      if (this.user?.role === 'admin') return '管理员'
      if (this.user?.role === 'hr') return '人事专员'
      return '员工'
    }
  }
}
</script>

<style scoped>
/* 整体布局 */
.app-layout {
  min-height: 100vh;
  background: var(--background-color);
}

/* 自定义侧边栏 */
.custom-sider {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 1000;
}

.custom-sider::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(24, 144, 255, 0.1) 0%, rgba(135, 208, 104, 0.1) 100%);
  pointer-events: none;
}

/* Logo区域 */
.logo {
  height: 80px;
  padding: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  margin-bottom: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  position: relative;
  overflow: hidden;
}

.logo::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
  transition: left 0.5s;
}

.logo:hover::before {
  left: 100%;
}

.logo-icon {
  width: 32px;
  height: 32px;
  margin-right: 12px;
  color: #1890ff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.logo.collapsed .logo-icon {
  margin-right: 0;
}

.logo-text h2 {
  color: white;
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  line-height: 1.2;
}

.logo-text span {
  color: rgba(255, 255, 255, 0.7);
  font-size: 12px;
  margin-top: 2px;
  display: block;
}

.logo-text-mini h2 {
  color: white;
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

/* 动画效果 */
.fade-slide-enter-active, .fade-slide-leave-active {
  transition: all 0.3s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateX(-10px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateX(10px);
}

/* 菜单样式 */
.custom-menu {
  background: transparent;
  border: none;
  padding: 0 8px;
}

.custom-menu :deep(.ant-menu-item),
.custom-menu :deep(.ant-menu-submenu-title) {
  height: 48px;
  line-height: 48px;
  margin: 4px 0;
  border-radius: 8px;
  padding: 0 16px;
  color: rgba(255, 255, 255, 0.85);
  border: none;
  background: transparent;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.custom-menu :deep(.ant-menu-item::after),
.custom-menu :deep(.ant-menu-submenu-title::after) {
  display: none;
}

.custom-menu :deep(.ant-menu-item:hover),
.custom-menu :deep(.ant-menu-submenu-title:hover),
.custom-menu :deep(.ant-menu-item-selected) {
  background: rgba(24, 144, 255, 0.15);
  color: #fff;
  transform: translateX(4px);
}

.custom-menu :deep(.ant-menu-item-selected) {
  background: linear-gradient(135deg, rgba(24, 144, 255, 0.2) 0%, rgba(135, 208, 104, 0.2) 100%);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
}

.menu-item-content {
  display: flex;
  align-items: center;
  width: 100%;
}

.menu-icon {
  width: 20px;
  height: 20px;
  margin-right: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: inherit;
  font-size: 16px;
}

.menu-text {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
}

/* 子菜单样式 */
.custom-menu :deep(.ant-menu-sub) {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
  margin: 4px 0;
  padding: 8px 0;
}

.submenu-item-content {
  display: flex;
  align-items: center;
  padding-left: 12px;
}

.submenu-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.5);
  margin-right: 12px;
  transition: all 0.3s;
}

.custom-menu :deep(.ant-menu-item-selected) .submenu-dot {
  background: #1890ff;
}

/* 头部样式 */
.app-header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  position: relative;
  z-index: 999;
}

.header-left {
  display: flex;
  align-items: center;
  flex: 1;
}

.trigger-container {
  margin-right: 24px;
}

.trigger {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.3s;
  color: var(--text-color);
}

.trigger:hover {
  background: var(--primary-color);
  color: white;
  transform: scale(1.1);
}

.flip-enter-active, .flip-leave-active {
  transition: all 0.3s;
}

.flip-enter-from {
  transform: rotateY(90deg);
  opacity: 0;
}

.flip-leave-to {
  transform: rotateY(-90deg);
  opacity: 0;
}

.breadcrumb-container {
  flex: 1;
}

.header-breadcrumb {
  font-size: 14px;
}

.header-breadcrumb :deep(.ant-breadcrumb-link) {
  color: var(--text-color-secondary);
  transition: color 0.3s;
}

.header-breadcrumb :deep(.ant-breadcrumb-link:hover) {
  color: var(--primary-color);
}

.current-page {
  color: var(--text-color);
  font-weight: 500;
}

/* 头部右侧操作区 */
.header-right {
  display: flex;
  align-items: center;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-item {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 16px;
  color: var(--text-color-secondary);
}

.action-item:hover {
  background: var(--primary-color);
  color: white;
  transform: translateY(-2px);
}

.search-btn.action-item {
  background: var(--primary-color);
  color: white;
}

/* 用户信息 */
.user-info {
  margin-left: 16px;
}

.user-dropdown-trigger {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  gap: 12px;
}

.user-dropdown-trigger:hover {
  background: var(--background-color);
}

.user-avatar {
  background: linear-gradient(135deg, var(--primary-color), var(--primary-color-hover));
  border: 2px solid rgba(255, 255, 255, 0.8);
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.3);
  width: 36px;
  height: 36px;
}

.user-details {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color);
  line-height: 1.2;
}

.user-role {
  font-size: 12px;
  color: var(--text-color-secondary);
  line-height: 1.2;
}

.dropdown-arrow {
  color: var(--text-color-secondary);
  font-size: 12px;
  transition: transform 0.3s;
}

.user-dropdown-trigger:hover .dropdown-arrow {
  transform: translateY(-1px);
}

/* 用户菜单 */
.user-menu {
  min-width: 160px;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.user-menu :deep(.ant-menu-item) {
  padding: 12px 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s;
}

.user-menu :deep(.ant-menu-item:hover) {
  background: var(--background-color);
}

.logout-item {
  color: var(--error-color);
}

.logout-item:hover {
  background: rgba(245, 34, 45, 0.1);
}

/* 全局搜索 */
.global-search {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  padding: 16px 24px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  display: flex;
  justify-content: center;
}

.slide-down-enter-active, .slide-down-leave-active {
  transition: all 0.3s ease;
}

.slide-down-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* 内容区域 */
.app-content {
  background: var(--background-color);
  padding: 24px;
  min-height: calc(100vh - 64px);
}

.content-wrapper {
  background: var(--component-background);
  border-radius: 12px;
  padding: 24px;
  min-height: calc(100vh - 112px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  position: relative;
  overflow: hidden;
}

.content-wrapper::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--primary-color), var(--success-color));
  border-radius: 12px 12px 0 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-header {
    padding: 0 16px;
  }
  
  .header-actions {
    gap: 4px;
  }
  
  .user-details {
    display: none;
  }
  
  .breadcrumb-container {
    display: none;
  }
  
  .app-content {
    padding: 16px;
  }
  
  .content-wrapper {
    padding: 16px;
  }
}

/* 滚动条优化 */
.custom-sider ::-webkit-scrollbar {
  width: 4px;
}

.custom-sider ::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
}

.custom-sider ::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 2px;
}

.custom-sider ::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.5);
}
</style>