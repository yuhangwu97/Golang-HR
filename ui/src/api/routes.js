import { apiService } from '@/services/api'
import { routeConfig, getAllRoutes, getRoutesByRole, generateNavigation } from '@/config/routes'

// 路由API服务
export class RouteApiService {
  // 获取所有路由配置
  static getAllRoutes() {
    return Promise.resolve({
      success: true,
      data: getAllRoutes(),
      message: '获取路由配置成功'
    })
  }

  // 根据用户权限获取路由
  static getRoutesForUser(userRoles = []) {
    const routes = getRoutesByRole(userRoles)
    return Promise.resolve({
      success: true,
      data: routes,
      message: '获取用户路由成功'
    })
  }

  // 获取导航菜单
  static getNavigation(userRoles = []) {
    const navigation = generateNavigation(userRoles)
    return Promise.resolve({
      success: true,
      data: navigation,
      message: '获取导航菜单成功'
    })
  }

  // 获取路由权限映射
  static getRoutePermissions() {
    const permissions = {}
    const allRoutes = getAllRoutes()
    
    allRoutes.forEach(route => {
      permissions[route.path] = {
        requiresAuth: route.meta.requiresAuth || false,
        requiresRole: route.meta.requiresRole || [],
        title: route.meta.title,
        group: route.meta.group
      }
    })

    return Promise.resolve({
      success: true,
      data: permissions,
      message: '获取路由权限映射成功'
    })
  }

  // 检查用户是否有访问特定路由的权限
  static checkRoutePermission(routePath, userRoles = []) {
    const allRoutes = getAllRoutes()
    const route = allRoutes.find(r => r.path === routePath)
    
    if (!route) {
      return Promise.resolve({
        success: false,
        data: { hasPermission: false },
        message: '路由不存在'
      })
    }

    let hasPermission = true

    // 检查是否需要认证
    if (route.meta.requiresAuth && !userRoles.length) {
      hasPermission = false
    }

    // 检查角色权限
    if (route.meta.requiresRole && route.meta.requiresRole.length) {
      hasPermission = route.meta.requiresRole.some(role => userRoles.includes(role))
    }

    return Promise.resolve({
      success: true,
      data: { 
        hasPermission,
        route: route
      },
      message: hasPermission ? '有访问权限' : '无访问权限'
    })
  }

  // 生成路由面包屑
  static generateBreadcrumb(routePath) {
    const allRoutes = getAllRoutes()
    const route = allRoutes.find(r => r.path === routePath)
    
    if (!route) {
      return Promise.resolve({
        success: false,
        data: [],
        message: '路由不存在'
      })
    }

    const breadcrumb = []
    
    // 添加根路径
    breadcrumb.push({ title: '首页', path: '/dashboard' })
    
    // 添加父级路径
    if (route.meta.parent) {
      const parentRoute = allRoutes.find(r => r.path === route.meta.parent)
      if (parentRoute) {
        breadcrumb.push({ 
          title: parentRoute.meta.title, 
          path: parentRoute.path 
        })
      }
    }
    
    // 添加当前路径
    if (route.path !== '/dashboard') {
      breadcrumb.push({ 
        title: route.meta.title, 
        path: route.path,
        current: true
      })
    }

    return Promise.resolve({
      success: true,
      data: breadcrumb,
      message: '生成面包屑成功'
    })
  }
}

// 导出服务实例
export const routeApi = RouteApiService

// 如果需要与后端API集成，可以使用以下方法
export const routeApiBackend = {
  // 从后端获取路由配置
  async getRoutesFromBackend() {
    return apiService.get('/api/routes')
  },

  // 从后端获取用户路由权限
  async getUserRoutePermissions(userId) {
    return apiService.get(`/api/users/${userId}/routes`)
  },

  // 从后端获取导航菜单
  async getNavigationFromBackend(userId) {
    return apiService.get(`/api/users/${userId}/navigation`)
  },

  // 同步前端路由配置到后端
  async syncRoutesToBackend() {
    const routes = getAllRoutes()
    return apiService.post('/api/routes/sync', { routes })
  }
}