import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '@/store'

Vue.use(VueRouter)

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/components/layout/AppLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/DashboardView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'employees',
        name: 'Employees',
        component: () => import('@/views/employee/EmployeeListView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'employees/create',
        name: 'EmployeeCreate',
        component: () => import('@/views/employee/EmployeeCreateView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'employees/test-create',
        name: 'TestEmployeeCreate',
        component: () => import('@/views/employee/TestCreateView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'employees/:id',
        name: 'EmployeeDetail',
        component: () => import('@/views/employee/EmployeeDetailView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'employees/:id/edit',
        name: 'EmployeeEdit',
        component: () => import('@/views/employee/EmployeeEditView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'organization/departments',
        name: 'Departments',
        component: () => import('@/views/organization/DepartmentListView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'organization/positions',
        name: 'Positions',
        component: () => import('@/views/organization/PositionListView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'organization/job-levels',
        name: 'JobLevels',
        component: () => import('@/views/organization/JobLevelListView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'organization/chart',
        name: 'OrganizationChart',
        component: () => import('@/views/organization/OrganizationChartView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'system/users',
        name: 'Users',
        component: () => import('@/views/system/UserListView.vue'),
        meta: { requiresAuth: true, requiresRole: ['admin'] }
      },
      {
        path: 'system/roles',
        name: 'Roles',
        component: () => import('@/views/system/RoleListView.vue'),
        meta: { requiresAuth: true, requiresRole: ['admin'] }
      },
      {
        path: 'system/permissions',
        name: 'Permissions',
        component: () => import('@/views/system/PermissionListView.vue'),
        meta: { requiresAuth: true, requiresRole: ['admin'] }
      },
      {
        path: 'salary',
        name: 'Salary',
        component: () => import('@/views/salary/SalaryListView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'my-salary',
        name: 'MySalary',
        component: () => import('@/views/salary/MySalaryView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/ProfileView.vue'),
        meta: { requiresAuth: true }
      }
    ]
  }
]

const router = new VueRouter({
  mode: 'hash', // 临时使用hash模式避免404问题
  base: process.env.BASE_URL,
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const hasToken = store.state.auth.token
  
  // 避免重定向循环：如果已经在目标页面，直接通过
  if (to.name === 'Login' && hasToken) {
    // 防止从dashboard重定向回来时再次重定向
    if (from.path === '/dashboard' || from.path === '/') {
      return next(false)
    }
    return next({ name: 'Dashboard' })
  }
  
  // 如果需要认证但没有token，重定向到登录页
  if (requiresAuth && !hasToken) {
    // 防止重复重定向到登录页
    if (to.path === '/login') {
      return next()
    }
    return next({ name: 'Login' })
  }
  
  // 其他情况直接通过
  next()
})

export default router