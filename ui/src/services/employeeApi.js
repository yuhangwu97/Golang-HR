import { apiService } from './api'

class EmployeeApiService {
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

  async searchEmployees(query) {
    return apiService.get('/employees/search', { params: { q: query } })
  }

  async getEmployeeHistory(id) {
    return apiService.get(`/employees/${id}/history`)
  }

  async updateEmployeeStatus(id, status) {
    return apiService.put(`/employees/${id}/status`, { status })
  }

  async getEmployeeContracts(id) {
    return apiService.get(`/employees/${id}/contracts`)
  }

  async addEmployeeContract(id, contractData) {
    return apiService.post(`/employees/${id}/contracts`, contractData)
  }

  async updateEmployeeContract(employeeId, contractId, contractData) {
    return apiService.put(`/employees/${employeeId}/contracts/${contractId}`, contractData)
  }

  async deleteEmployeeContract(employeeId, contractId) {
    return apiService.delete(`/employees/${employeeId}/contracts/${contractId}`)
  }

  // 获取部门的完整员工列表（包含部门负责人和所有下属，含子部门员工）
  async getDepartmentEmployeesWithHierarchy(departmentId, params = {}) {
    return apiService.get(`/departments/${departmentId}/employees/hierarchy`, { 
      params: {
        include_manager: true,
        include_subordinates: true,
        include_subdepartments: true,
        ...params
      }
    })
  }

  // 获取部门直接员工（不包含子部门）
  async getDepartmentDirectEmployees(departmentId, params = {}) {
    return apiService.get(`/departments/${departmentId}/employees`, { 
      params: {
        include_manager: true,
        ...params
      }
    })
  }

  // 获取员工的所有下级（递归）
  async getEmployeeSubordinates(employeeId, params = {}) {
    return apiService.get(`/employees/${employeeId}/subordinates`, { 
      params: {
        recursive: true,
        ...params
      }
    })
  }
}

export const employeeApiService = new EmployeeApiService()
export default employeeApiService