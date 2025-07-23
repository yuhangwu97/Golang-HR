// 路由使用示例
import { routeApi } from '@/api/routes'
import { routePermissionManager } from '@/utils/routePermissions'
import { routeConfig, generateNavigation } from '@/config/routes'

// 示例1: 获取所有路由
export async function getAllRoutesExample() {
  try {
    const response = await routeApi.getAllRoutes()
    console.log('所有路由:', response.data)
    return response.data
  } catch (error) {
    console.error('获取路由失败:', error)
  }
}

// 示例2: 根据用户角色获取路由
export async function getRoutesForUserExample(userRoles = ['employee']) {
  try {
    const response = await routeApi.getRoutesForUser(userRoles)
    console.log(`角色 ${userRoles.join(', ')} 的路由:`, response.data)
    return response.data
  } catch (error) {
    console.error('获取用户路由失败:', error)
  }
}

// 示例3: 获取导航菜单
export async function getNavigationExample(userRoles = ['hr']) {
  try {
    const response = await routeApi.getNavigation(userRoles)
    console.log(`角色 ${userRoles.join(', ')} 的导航菜单:`, response.data)
    return response.data
  } catch (error) {
    console.error('获取导航菜单失败:', error)
  }
}

// 示例4: 检查路由权限
export async function checkRoutePermissionExample(routePath = '/system/users', userRoles = ['admin']) {
  try {
    const response = await routeApi.checkRoutePermission(routePath, userRoles)
    console.log(`路由 ${routePath} 权限检查:`, response.data)
    return response.data.hasPermission
  } catch (error) {
    console.error('检查路由权限失败:', error)
  }
}

// 示例5: 生成面包屑导航
export async function generateBreadcrumbExample(routePath = '/employees/create') {
  try {
    const response = await routeApi.generateBreadcrumb(routePath)
    console.log(`路由 ${routePath} 的面包屑:`, response.data)
    return response.data
  } catch (error) {
    console.error('生成面包屑失败:', error)
  }
}

// 示例6: 使用权限管理器
export async function permissionManagerExample() {
  const userRoles = ['hr', 'employee']
  
  // 检查特定路由权限
  const hasPermission = await routePermissionManager.checkPermission('/employees', userRoles)
  console.log('是否有员工管理权限:', hasPermission)
  
  // 获取用户可访问的路由
  const accessibleRoutes = await routePermissionManager.getAccessibleRoutes(userRoles)
  console.log('可访问的路由:', accessibleRoutes)
  
  return { hasPermission, accessibleRoutes }
}

// 示例7: 直接从配置文件获取路由
export function getRoutesFromConfigExample() {
  console.log('路由配置:', routeConfig)
  
  // 生成特定角色的导航
  const adminNavigation = generateNavigation(['admin'])
  console.log('管理员导航:', adminNavigation)
  
  const employeeNavigation = generateNavigation(['employee'])
  console.log('员工导航:', employeeNavigation)
  
  return { adminNavigation, employeeNavigation }
}

// 示例8: 从外部JSON文件获取路由
export async function getRoutesFromJsonExample() {
  try {
    const response = await fetch('/api/routes.json')
    const routesData = await response.json()
    console.log('从JSON文件获取的路由:', routesData)
    return routesData
  } catch (error) {
    console.error('从JSON文件获取路由失败:', error)
  }
}

// 示例9: 过滤用户可访问的导航菜单
export function filterNavigationByRoleExample(userRoles = ['hr']) {
  const navigation = [
    { key: 'dashboard', title: '仪表盘', roles: ['admin', 'hr', 'employee'] },
    { key: 'employees', title: '员工管理', roles: ['admin', 'hr'] },
    { key: 'system', title: '系统管理', roles: ['admin'] },
    { key: 'profile', title: '个人资料', roles: ['admin', 'hr', 'employee'] }
  ]
  
  const filteredNavigation = navigation.filter(item => {
    return item.roles.some(role => userRoles.includes(role))
  })
  
  console.log(`角色 ${userRoles.join(', ')} 的过滤后导航:`, filteredNavigation)
  return filteredNavigation
}

// 示例10: 动态路由权限检查
export function dynamicRoutePermissionExample(currentRoute, userRoles) {
  const permissionMap = {
    '/dashboard': ['admin', 'hr', 'employee'],
    '/employees': ['admin', 'hr'],
    '/system/users': ['admin'],
    '/salary': ['admin', 'hr'],
    '/my-salary': ['admin', 'hr', 'employee'],
    '/profile': ['admin', 'hr', 'employee']
  }
  
  const requiredRoles = permissionMap[currentRoute]
  if (!requiredRoles) {
    console.log(`路由 ${currentRoute} 未找到权限配置`)
    return false
  }
  
  const hasPermission = requiredRoles.some(role => userRoles.includes(role))
  console.log(`路由 ${currentRoute} 权限检查结果:`, hasPermission)
  
  return hasPermission
}

// 导出所有示例函数
export default {
  getAllRoutesExample,
  getRoutesForUserExample,
  getNavigationExample,
  checkRoutePermissionExample,
  generateBreadcrumbExample,
  permissionManagerExample,
  getRoutesFromConfigExample,
  getRoutesFromJsonExample,
  filterNavigationByRoleExample,
  dynamicRoutePermissionExample
}