import { apiService } from './api'

class PositionApiService {
  async getPositions(params = {}) {
    return apiService.get('/positions', { params })
  }

  async getPosition(id) {
    return apiService.get(`/positions/${id}`)
  }

  async createPosition(positionData) {
    return apiService.post('/positions', positionData)
  }

  async updatePosition(id, positionData) {
    return apiService.put(`/positions/${id}`, positionData)
  }

  async deletePosition(id) {
    return apiService.delete(`/positions/${id}`)
  }

  async getPositionsByDepartment(departmentId) {
    return apiService.get(`/positions/department/${departmentId}`)
  }

  async getPositionEmployees(id) {
    return apiService.get(`/positions/${id}/employees`)
  }

  async getPositionStatistics() {
    return apiService.get('/positions/statistics')
  }

  async searchPositions(query) {
    return apiService.get('/positions/search', { params: { q: query } })
  }

  async getPositionTree() {
    return apiService.get('/positions/tree')
  }

  async getAllPositions() {
    return apiService.get('/positions/all')
  }
}

export const positionApiService = new PositionApiService()
export const positionApi = positionApiService
export default positionApiService