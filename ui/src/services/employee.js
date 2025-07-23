import { apiService } from './api'

class EmployeeService {
  async getEmployees(params = {}) {
    return apiService.get('/employees', { params })
  }

  async getEmployee(id) {
    return apiService.get(`/employees/${id}`)
  }

  async createEmployee(employeeData) {
    return apiService.post('/employees', employeeData)
  }

  async updateEmployee(id, employeeData) {
    return apiService.put(`/employees/${id}`, employeeData)
  }

  async deleteEmployee(id) {
    return apiService.delete(`/employees/${id}`)
  }

  async getEmployeeStatistics() {
    return apiService.get('/employees/statistics')
  }

  async exportEmployees(format = 'excel') {
    const response = await apiService.get('/employees/export', {
      params: { format },
      responseType: 'blob'
    })
    return response
  }

  async importEmployees(file) {
    const formData = new FormData()
    formData.append('file', file)
    return apiService.post('/employees/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }

  async bulkUpdateEmployees(updates) {
    return apiService.put('/employees/bulk', { updates })
  }

  async getEmployeesByDepartment(departmentId) {
    return apiService.get(`/employees/department/${departmentId}`)
  }

  async getEmployeesByStatus(status) {
    return apiService.get('/employees', { params: { status } })
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

  // 获取职级列表
  async getJobLevels() {
    try {
      const response = await apiService.get('/job_levels')
      return {
        success: true,
        data: response.data?.data || response.data || []
      }
    } catch (error) {
      console.error('获取职级列表失败:', error)
      return {
        success: false,
        data: [],
        error: error.message
      }
    }
  }

  // 获取经理列表
  async getManagers() {
    try {
      const response = await apiService.get('/employees', { 
        params: { 
          status: 'active',
          page: 1,
          pageSize: 100
        } 
      })
      return {
        success: true,
        data: response.data?.data || response.data || []
      }
    } catch (error) {
      console.error('获取经理列表失败:', error)
      return {
        success: false,
        data: [],
        error: error.message
      }
    }
  }

  // 根据部门获取职位
  async getPositionsByDepartment(departmentId) {
    try {
      const response = await apiService.get(`/positions/department/${departmentId}`)
      return {
        success: true,
        data: response.data?.data || response.data || []
      }
    } catch (error) {
      console.error('获取职位列表失败:', error)
      return {
        success: false,
        data: [],
        error: error.message
      }
    }
  }
}

export const employeeService = new EmployeeService()
export default employeeService