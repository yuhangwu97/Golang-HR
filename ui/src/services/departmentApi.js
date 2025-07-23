import { apiService } from './api'

class DepartmentApiService {
  async getDepartments(params = {}) {
    return apiService.get('/departments', { params })
  }

  async getDepartment(id) {
    return apiService.get(`/departments/${id}`)
  }

  async createDepartment(departmentData) {
    return apiService.post('/departments', departmentData)
  }

  async updateDepartment(id, departmentData) {
    return apiService.put(`/departments/${id}`, departmentData)
  }

  async deleteDepartment(id) {
    return apiService.delete(`/departments/${id}`)
  }

  async getDepartmentTree() {
    return apiService.get('/departments/tree')
  }

  async moveDepartment(id, parentId) {
    return apiService.put(`/departments/${id}/move`, { parentId })
  }

  async getDepartmentEmployees(id) {
    return apiService.get(`/departments/${id}/employees`)
  }

  async getDepartmentStatistics() {
    return apiService.get(`/departments/statistics`)
  }
}

export const departmentApiService = new DepartmentApiService()
export const departmentApi = departmentApiService
export default departmentApiService