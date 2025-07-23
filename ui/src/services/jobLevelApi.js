import { apiService } from './api'

class JobLevelApiService {
  async getJobLevels(params = {}) {
    return apiService.get('/job_levels', { params })
  }

  async getJobLevel(id) {
    return apiService.get(`/job_levels/${id}`)
  }

  async createJobLevel(jobLevelData) {
    return apiService.post('/job_levels', jobLevelData)
  }

  async updateJobLevel(id, jobLevelData) {
    return apiService.put(`/job_levels/${id}`, jobLevelData)
  }

  async deleteJobLevel(id) {
    return apiService.delete(`/job_levels/${id}`)
  }

  async getJobLevelEmployees(id) {
    return apiService.get(`/job_levels/${id}/employees`)
  }

  async getJobLevelStatistics() {
    return apiService.get('/job_levels/statistics')
  }

  async searchJobLevels(query) {
    return apiService.get('/job_levels/search', { params: { q: query } })
  }
}

export const jobLevelApiService = new JobLevelApiService()
export const jobLevelApi = jobLevelApiService
export default jobLevelApiService