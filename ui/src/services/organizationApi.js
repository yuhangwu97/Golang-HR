import { apiService } from './api'

class OrganizationApiService {
  // 组织架构树形数据
  async getTree() {
    return apiService.get('/departments/tree')
  }

  // 获取组织单位列表
  async getUnits(params = {}) {
    return apiService.get('/departments', { params })
  }

  // 获取单个组织单位
  async getUnit(id) {
    return apiService.get(`/departments/${id}`)
  }

  // 创建组织单位
  async createUnit(unitData) {
    return apiService.post('/departments', unitData)
  }

  // 更新组织单位
  async updateUnit(id, unitData) {
    return apiService.put(`/departments/${id}`, unitData)
  }

  // 删除组织单位
  async deleteUnit(id) {
    return apiService.delete(`/departments/${id}`)
  }

  // 移动组织单位
  async moveUnit(id, newParentId) {
    return apiService.put(`/departments/${id}/move`, { parent_id: newParentId })
  }

  // 获取组织架构变更历史
  async getChangeHistory(params = {}) {
    return apiService.get('/organization/changes', { params })
  }

  // 获取组织架构快照
  async getSnapshot(date) {
    return apiService.get('/organization/snapshot', { params: { date } })
  }

  // 获取员工分配
  async getEmployeeAssignments(params = {}) {
    return apiService.get('/organization/assignments', { params })
  }

  // 分配员工到组织单位
  async assignEmployee(assignmentData) {
    return apiService.post('/organization/assignments', assignmentData)
  }

  // 更新员工分配
  async updateAssignment(id, assignmentData) {
    return apiService.put(`/organization/assignments/${id}`, assignmentData)
  }

  // 删除员工分配
  async removeAssignment(id) {
    return apiService.delete(`/organization/assignments/${id}`)
  }

  // 批量分配员工
  async bulkAssignEmployees(assignments) {
    return apiService.post('/organization/assignments/bulk', { assignments })
  }

  // 获取组织架构统计数据
  async getStatistics() {
    return apiService.get('/departments/statistics')
  }

  // 获取组织架构图表数据
  async getChart() {
    return apiService.get('/departments/chart')
  }

  // 获取组织架构层级信息
  async getHierarchy(id) {
    return apiService.get(`/departments/${id}/hierarchy`)
  }

  // 获取组织架构变更统计
  async getChangesStatistics(params = {}) {
    return apiService.get('/organization/statistics/changes', { params })
  }

  // 导出组织架构
  async exportStructure(format = 'excel') {
    return apiService.get('/organization/export', { 
      params: { format },
      responseType: 'blob'
    })
  }

  // 导入组织架构
  async importStructure(file) {
    const formData = new FormData()
    formData.append('file', file)
    return apiService.post('/organization/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  }

  // 搜索组织单位
  async searchUnits(query) {
    return apiService.get('/departments/search', { params: { q: query } })
  }

  // 获取组织单位路径
  async getUnitPath(id) {
    return apiService.get(`/departments/${id}/path`)
  }

  // 批量更新组织架构排序
  async updateSort(updates) {
    return apiService.put('/organization/sort', { updates })
  }

  // 获取组织架构模板
  async getTemplates() {
    return apiService.get('/organization/templates')
  }

  // 应用组织架构模板
  async applyTemplate(templateId) {
    return apiService.post(`/organization/templates/${templateId}/apply`)
  }

  // 验证组织架构变更
  async validateChanges(changes) {
    return apiService.post('/organization/validate', { changes })
  }

  // 获取组织架构分析报告
  async getAnalysisReport(params = {}) {
    return apiService.get('/organization/analysis', { params })
  }

  // 获取子部门
  async getChildren(id) {
    return apiService.get(`/departments/${id}/subdepartments`)
  }
}

export const organizationApiService = new OrganizationApiService()
export default organizationApiService