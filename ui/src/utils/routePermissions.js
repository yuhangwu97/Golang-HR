import { routeApi } from '@/api/routes'
import store from '@/store'

// 路由权限管理工具
export class RoutePermissionManager {
  constructor() {
    this.permissions = new Map()
    this.initialized = false
  }

  // 初始化权限映射
  async initialize() {
    if (this.initialized) return
    
    try {
      const response = await routeApi.getRoutePermissions()
      if (response.success) {
        Object.entries(response.data).forEach(([path, permission]) => {
          this.permissions.set(path, permission)
        })
        this.initialized = true
      }
    } catch (error) {
      console.error('初始化路由权限失败:', error)
    }
  }

  // 检查路由权限
  async checkPermission(routePath, userRoles = []) {
    await this.initialize()
    
    const permission = this.permissions.get(routePath)
    if (!permission) return false

    // 检查是否需要认证
    if (permission.requiresAuth && !userRoles.length) {
      return false
    }

    // 检查角色权限
    if (permission.requiresRole && permission.requiresRole.length) {
      return permission.requiresRole.some(role => userRoles.includes(role))
    }

    return true
  }

  // 获取用户可访问的路由
  async getAccessibleRoutes(userRoles = []) {
    await this.initialize()
    
    const accessibleRoutes = []
    for (const [path, permission] of this.permissions) {
      const hasPermission = await this.checkPermission(path, userRoles)
      if (hasPermission) {
        accessibleRoutes.push({ path, ...permission })
      }
    }
    
    return accessibleRoutes
  }

  // 清除权限缓存
  clearCache() {
    this.permissions.clear()
    this.initialized = false
  }
}

// 创建全局权限管理器实例
export const routePermissionManager = new RoutePermissionManager()

// 路由权限检查中间件
export const routePermissionMiddleware = async (to, from, next) => {
  const userRoles = store.getters['auth/user']?.roles || []
  const hasPermission = await routePermissionManager.checkPermission(to.path, userRoles)
  
  if (!hasPermission) {
    // 无权限，重定向到403页面或首页
    next('/403')
    return
  }
  
  next()
}

// 动态生成路由权限映射
export const generateRoutePermissionMap = () => {
  const routes = [
    // 公共路由
    { path: '/login', roles: [], auth: false },
    { path: '/404', roles: [], auth: false },
    { path: '/403', roles: [], auth: false },
    
    // 基础路由
    { path: '/dashboard', roles: ['admin', 'hr', 'employee'], auth: true },
    { path: '/profile', roles: ['admin', 'hr', 'employee'], auth: true },
    
    // 员工管理
    { path: '/employees', roles: ['admin', 'hr'], auth: true },
    { path: '/employees/create', roles: ['admin', 'hr'], auth: true },
    { path: '/employees/:id', roles: ['admin', 'hr'], auth: true },
    { path: '/employees/:id/edit', roles: ['admin', 'hr'], auth: true },
    
    // 组织管理
    { path: '/organization/departments', roles: ['admin', 'hr'], auth: true },
    { path: '/organization/positions', roles: ['admin', 'hr'], auth: true },
    { path: '/organization/job-levels', roles: ['admin', 'hr'], auth: true },
    
    // 系统管理
    { path: '/system/users', roles: ['admin'], auth: true },
    { path: '/system/roles', roles: ['admin'], auth: true },
    { path: '/system/permissions', roles: ['admin'], auth: true },
    
    // 薪资管理
    { path: '/salary', roles: ['admin', 'hr'], auth: true },
    { path: '/my-salary', roles: ['admin', 'hr', 'employee'], auth: true }
  ]
  
  const permissionMap = {}
  routes.forEach(route => {
    permissionMap[route.path] = {
      requiresAuth: route.auth,
      requiresRole: route.roles,
      accessible: true
    }
  })
  
  return permissionMap
}

// 路由权限验证装饰器
export const requirePermission = (roles = []) => {
  return (target, propertyKey, descriptor) => {
    const originalMethod = descriptor.value
    
    descriptor.value = async function(...args) {
      const userRoles = store.getters['auth/user']?.roles || []
      const hasPermission = roles.some(role => userRoles.includes(role))
      
      if (!hasPermission) {
        throw new Error('无权限访问')
      }
      
      return originalMethod.apply(this, args)
    }
    
    return descriptor
  }
}

// 导出默认权限映射
export default generateRoutePermissionMap()