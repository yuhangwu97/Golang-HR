import { apiService } from './api'

class SalaryService {
  async getSalaries(params = {}) {
    return apiService.get('/salaries', { params })
  }

  async getSalary(id) {
    return apiService.get(`/salaries/${id}`)
  }

  async createSalary(salaryData) {
    return apiService.post('/salaries', salaryData)
  }

  async updateSalary(id, salaryData) {
    return apiService.put(`/salaries/${id}`, salaryData)
  }

  async deleteSalary(id) {
    return apiService.delete(`/salaries/${id}`)
  }

  async getMySalary(month) {
    return apiService.get('/salaries/my-salary', { params: { month } })
  }

  async calculateSalary(employeeId, month) {
    return apiService.post('/salaries/calculate', {
      employee_id: employeeId,
      month
    })
  }

  async batchCalculateSalary(month, departmentId, employeeIds) {
    return apiService.post('/salaries/batch-calculate', {
      month,
      department_id: departmentId,
      employee_ids: employeeIds
    })
  }

  async approveSalary(id, status, remark) {
    return apiService.put(`/salaries/${id}/approve`, {
      status,
      remark
    })
  }

  async getSalaryStatistics(params = {}) {
    return apiService.get('/salaries/statistics', { params })
  }

  async exportSalaries(format = 'excel', params = {}) {
    const response = await apiService.get('/salaries/export', {
      params: { format, ...params },
      responseType: 'blob'
    })
    return response
  }

  async processPayroll(salaryId, paymentMethod, bankAccount) {
    return apiService.post(`/salaries/${salaryId}/process-payroll`, {
      payment_method: paymentMethod,
      bank_account: bankAccount
    })
  }

  async getPayrollRecords(salaryId) {
    return apiService.get(`/salaries/${salaryId}/payroll-records`)
  }

  async updatePayrollStatus(id, status, remark) {
    return apiService.put(`/payroll/${id}/status`, {
      status,
      remark
    })
  }

  // Enhanced salary methods for PayrollWorkflow component
  async getEnhancedSalaries(params = {}) {
    return apiService.get('/salaries/enhanced', { params })
  }

  async batchCalculateSalaries(params) {
    return apiService.post('/salaries/enhanced/batch-calculate', params)
  }

  async bulkApproveSalaries(params) {
    return apiService.post('/salaries/enhanced/bulk-approve', params)
  }

  async finalApproveSalaries(params) {
    return apiService.post('/salaries/enhanced/final-approve', params)
  }

  async createPaymentBatch(params) {
    return apiService.post('/salaries/payment-batches', params)
  }

  async processPaymentBatch(batchId) {
    return apiService.post(`/salaries/payment-batches/${batchId}/process`)
  }

  async processPayrollSingle(salaryId, params) {
    return apiService.post(`/salaries/${salaryId}/payroll-single`, params)
  }

  async approveSalary(salaryId, params) {
    return apiService.post(`/salaries/${salaryId}/approve`, params)
  }

  async getPayrollPeriods(params = {}) {
    return apiService.get('/salaries/payroll-periods', { params })
  }

  async getEmployees(params = {}) {
    return apiService.get('/employees', { params })
  }

  async getDepartments(params = {}) {
    return apiService.get('/departments', { params })
  }

  async getWorkflowHistory(periodId) {
    return apiService.get(`/salaries/workflow-history/${periodId || ''}`)
  }
}

export const salaryService = new SalaryService()
export default salaryService