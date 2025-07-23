// 前端路由配置暴露
export const routeConfig = {
  // 公共路由（不需要认证）
  public: [
    {
      path: '/login',
      name: 'Login',
      title: '登录',
      component: 'LoginView',
      meta: {
        requiresAuth: false,
        title: '登录',
        icon: 'login'
      }
    }
  ],

  // 需要认证的路由
  protected: [
    {
      path: '/dashboard',
      name: 'Dashboard',
      title: '仪表盘',
      component: 'DashboardView',
      meta: {
        requiresAuth: true,
        title: '仪表盘',
        icon: 'dashboard'
      }
    },
    {
      path: '/employees',
      name: 'Employees',
      title: '员工管理',
      component: 'EmployeeListView',
      meta: {
        requiresAuth: true,
        title: '员工管理',
        icon: 'team'
      }
    },
    {
      path: '/employees/create',
      name: 'EmployeeCreate',
      title: '创建员工',
      component: 'EmployeeCreateView',
      meta: {
        requiresAuth: true,
        title: '创建员工',
        icon: 'user-add',
        parent: '/employees'
      }
    },
    {
      path: '/employees/:id',
      name: 'EmployeeDetail',
      title: '员工详情',
      component: 'EmployeeDetailView',
      meta: {
        requiresAuth: true,
        title: '员工详情',
        icon: 'user',
        parent: '/employees'
      }
    },
    {
      path: '/employees/:id/edit',
      name: 'EmployeeEdit',
      title: '编辑员工',
      component: 'EmployeeEditView',
      meta: {
        requiresAuth: true,
        title: '编辑员工',
        icon: 'edit',
        parent: '/employees'
      }
    }
  ],

  // 组织管理路由
  organization: [
    {
      path: '/organization/departments',
      name: 'Departments',
      title: '部门管理',
      component: 'DepartmentListView',
      meta: {
        requiresAuth: true,
        title: '部门管理',
        icon: 'apartment',
        group: 'organization'
      }
    },
    {
      path: '/organization/positions',
      name: 'Positions',
      title: '职位管理',
      component: 'PositionListView',
      meta: {
        requiresAuth: true,
        title: '职位管理',
        icon: 'solution',
        group: 'organization'
      }
    },
    {
      path: '/organization/job-levels',
      name: 'JobLevels',
      title: '职级管理',
      component: 'JobLevelListView',
      meta: {
        requiresAuth: true,
        title: '职级管理',
        icon: 'trophy',
        group: 'organization'
      }
    }
  ],

  // 系统管理路由
  system: [
    {
      path: '/system/users',
      name: 'Users',
      title: '用户管理',
      component: 'UserListView',
      meta: {
        requiresAuth: true,
        requiresRole: ['admin'],
        title: '用户管理',
        icon: 'user',
        group: 'system'
      }
    },
    {
      path: '/system/roles',
      name: 'Roles',
      title: '角色管理',
      component: 'RoleListView',
      meta: {
        requiresAuth: true,
        requiresRole: ['admin'],
        title: '角色管理',
        icon: 'team',
        group: 'system'
      }
    },
    {
      path: '/system/permissions',
      name: 'Permissions',
      title: '权限管理',
      component: 'PermissionListView',
      meta: {
        requiresAuth: true,
        requiresRole: ['admin'],
        title: '权限管理',
        icon: 'safety',
        group: 'system'
      }
    }
  ],

  // 薪资管理路由
  salary: [
    {
      path: '/salary',
      name: 'Salary',
      title: '薪资管理',
      component: 'SalaryListView',
      meta: {
        requiresAuth: true,
        title: '薪资管理',
        icon: 'dollar',
        group: 'salary'
      }
    },
    {
      path: '/my-salary',
      name: 'MySalary',
      title: '我的薪资',
      component: 'MySalaryView',
      meta: {
        requiresAuth: true,
        title: '我的薪资',
        icon: 'wallet',
        group: 'salary'
      }
    }
  ],

  // 个人中心路由
  profile: [
    {
      path: '/profile',
      name: 'Profile',
      title: '个人资料',
      component: 'ProfileView',
      meta: {
        requiresAuth: true,
        title: '个人资料',
        icon: 'user'
      }
    }
  ]
}

// 获取所有路由的扁平数组
export const getAllRoutes = () => {
  const allRoutes = []
  Object.values(routeConfig).forEach(group => {
    allRoutes.push(...group)
  })
  return allRoutes
}

// 根据权限获取路由
export const getRoutesByRole = (roles = []) => {
  const allRoutes = getAllRoutes()
  return allRoutes.filter(route => {
    if (!route.meta.requiresRole) return true
    return route.meta.requiresRole.some(role => roles.includes(role))
  })
}

// 生成导航菜单数据
export const generateNavigation = (userRoles = []) => {
  const routes = getRoutesByRole(userRoles)
  const navigation = []

  // 分组处理
  const groups = {
    dashboard: { title: '仪表盘', icon: 'dashboard', routes: [] },
    employees: { title: '员工管理', icon: 'team', routes: [] },
    organization: { title: '组织管理', icon: 'apartment', routes: [] },
    system: { title: '系统管理', icon: 'setting', routes: [] },
    salary: { title: '薪资管理', icon: 'dollar', routes: [] },
    profile: { title: '个人中心', icon: 'user', routes: [] }
  }

  routes.forEach(route => {
    const group = route.meta.group || route.name.toLowerCase()
    if (groups[group]) {
      groups[group].routes.push(route)
    }
  })

  // 只返回有路由的分组
  Object.entries(groups).forEach(([key, group]) => {
    if (group.routes.length > 0) {
      navigation.push({
        key,
        ...group
      })
    }
  })

  return navigation
}

// 导出路由配置给外部使用
export default routeConfig