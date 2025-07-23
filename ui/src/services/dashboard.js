import { apiService } from './api'

class DashboardService {
  // 获取员工统计数据
  async getEmployeeStatistics() {
    try {
      const response = await apiService.get('/employees/statistics')
      return {
        success: true,
        data: response.data?.data || response.data || {}
      }
    } catch (error) {
      console.error('获取员工统计数据失败:', error)
      return {
        success: false,
        data: {},
        error: error.message
      }
    }
  }

  // 获取部门统计数据
  async getDepartmentStatistics() {
    try {
      const response = await apiService.get('/departments/statistics')
      return {
        success: true,
        data: response.data?.data || response.data || {}
      }
    } catch (error) {
      console.error('获取部门统计数据失败:', error)
      return {
        success: false,
        data: {},
        error: error.message
      }
    }
  }

  // 获取薪资统计数据
  async getSalaryStatistics(month = '', departmentId = null) {
    try {
      const params = {}
      if (month) params.month = month
      if (departmentId) params.departmentId = departmentId
      
      const response = await apiService.get('/salaries/statistics', { params })
      return {
        success: true,
        data: response.data?.data || response.data || {}
      }
    } catch (error) {
      console.error('获取薪资统计数据失败:', error)
      return {
        success: false,
        data: {},
        error: error.message
      }
    }
  }

  // 获取考勤统计数据
  async getAttendanceStatistics() {
    try {
      const response = await apiService.get('/attendance/statistics')
      return {
        success: true,
        data: response.data?.data || response.data || {}
      }
    } catch (error) {
      console.error('获取考勤统计数据失败:', error)
      return {
        success: false,
        data: {},
        error: error.message
      }
    }
  }

  // 获取部门列表
  async getDepartments() {
    try {
      const response = await apiService.get('/departments')
      return {
        success: true,
        data: response.data?.data || response.data || []
      }
    } catch (error) {
      console.error('获取部门列表失败:', error)
      return {
        success: false,
        data: [],
        error: error.message
      }
    }
  }

  // 获取最近入职员工
  async getRecentHires(limit = 10) {
    try {
      const response = await apiService.get('/employees', {
        params: {
          sortBy: 'hire_date',
          sortOrder: 'desc',
          page: 1,
          pageSize: limit
        }
      })
      return {
        success: true,
        data: response.data?.data || response.data || []
      }
    } catch (error) {
      console.error('获取最近入职员工失败:', error)
      return {
        success: false,
        data: [],
        error: error.message
      }
    }
  }

  // 获取综合仪表盘数据
  async getDashboardData() {
    try {
      const [
        employeeStats,
        departmentStats,
        salaryStats,
        attendanceStats,
        departments,
        recentHires
      ] = await Promise.all([
        this.getEmployeeStatistics(),
        this.getDepartmentStatistics(),
        this.getSalaryStatistics(),
        this.getAttendanceStatistics(),
        this.getDepartments(),
        this.getRecentHires(20)
      ])

      // 合并员工数据包含最近入职
      const employeeData = {
        ...employeeStats.data,
        recent_hires: recentHires.data
      }

      return {
        success: true,
        data: {
          employee: employeeData,
          department: departmentStats.data,
          salary: salaryStats.data,
          attendance: attendanceStats.data,
          departments: departments.data
        }
      }
    } catch (error) {
      console.error('获取仪表盘数据失败:', error)
      return {
        success: false,
        data: {},
        error: error.message
      }
    }
  }

  // 导出数据
  async exportData(type = 'employees', format = 'excel') {
    try {
      const endpoint = type === 'employees' ? '/employees/export' : '/salaries/export'
      const response = await apiService.get(endpoint, {
        params: { format },
        responseType: 'blob'
      })
      return response
    } catch (error) {
      console.error('导出数据失败:', error)
      throw error
    }
  }
}

export const dashboardService = new DashboardService()
export default dashboardService