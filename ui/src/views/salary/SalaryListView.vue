<template>
  <div class="salary-list">
    <div class="page-header">
      <h2>薪资管理</h2>
      <div class="page-actions">
        <el-button type="primary" @click="activeTab = 'workflow'" v-if="activeTab !== 'workflow'">
          <i class="el-icon-s-operation"></i>
          薪资发放流程
        </el-button>
        <el-button @click="activeTab = 'list'" v-if="activeTab !== 'list'">
          <i class="el-icon-s-grid"></i>
          薪资列表
        </el-button>
      </div>
    </div>

    <!-- 标签页切换 -->
    <el-tabs v-model="activeTab" type="card" class="salary-tabs">
      <el-tab-pane label="薪资列表" name="list">
        <!-- 原有薪资列表内容 -->
        <div class="salary-list-content">
          <!-- 悬浮操作按钮 -->
          <div class="floating-actions">
            <el-dropdown trigger="click" placement="top">
              <el-button type="primary" circle size="large" class="float-button">
                <i class="el-icon-plus"></i>
              </el-button>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item @click.native="showCalculateDialog = true">
                  <i class="el-icon-s-operation"></i>
                  单个计算薪资
                </el-dropdown-item>
                <el-dropdown-item @click.native="showBatchCalculateDialog = true">
                  <i class="el-icon-s-data"></i>
                  批量计算薪资
                </el-dropdown-item>
                <el-dropdown-item @click.native="showStatistics = !showStatistics">
                  <i class="el-icon-pie-chart"></i>
                  显示统计
                </el-dropdown-item>
                <el-dropdown-item @click.native="exportReport">
                  <i class="el-icon-download"></i>
                  导出报表
                </el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
          </div>
        </div>
      </el-tab-pane>
      
      <el-tab-pane label="薪资发放流程" name="workflow">
        <!-- 薪资发放流程组件 -->
        <payroll-workflow :period-id="selectedPeriodId" />
      </el-tab-pane>
    </el-tabs>

    <!-- 统计卡片 -->
    <div v-if="showStatistics" class="statistics-cards">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ statistics.total_employees || 0 }}</div>
              <div class="stat-label">总员工数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ formatCurrency(statistics.total_salary || 0) }}</div>
              <div class="stat-label">总薪资</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ formatCurrency(statistics.average_salary || 0) }}</div>
              <div class="stat-label">平均薪资</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ formatCurrency(statistics.max_salary || 0) }}</div>
              <div class="stat-label">最高薪资</div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 搜索过滤 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="月份">
          <el-date-picker
            v-model="searchForm.month"
            type="month"
            placeholder="选择月份"
            value-format="yyyy-MM"
            @change="handleSearch"
          />
        </el-form-item>
        <el-form-item label="员工ID">
          <el-input
            v-model="searchForm.employee_id"
            placeholder="输入员工ID"
            clearable
            @keyup.enter.native="handleSearch"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="选择状态" clearable>
            <el-option label="已计算" value="calculated" />
            <el-option label="已审批" value="approved" />
            <el-option label="已发放" value="paid" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <i class="el-icon-search"></i>
            搜索
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 薪资表格 -->
    <el-card>
      <el-table
        v-loading="loading"
        :data="salaries"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="employee.name" label="员工姓名" width="120" />
        <el-table-column prop="employee.department.name" label="部门" width="120" />
        <el-table-column prop="month" label="月份" width="100" />
        <el-table-column prop="base_salary" label="基本薪资" width="120">
          <template slot-scope="{ row }">
            {{ formatCurrency(row.base_salary) }}
          </template>
        </el-table-column>
        <el-table-column prop="gross_salary" label="应发薪资" width="120">
          <template slot-scope="{ row }">
            {{ formatCurrency(row.gross_salary) }}
          </template>
        </el-table-column>
        <el-table-column prop="net_salary" label="实发薪资" width="120">
          <template slot-scope="{ row }">
            {{ formatCurrency(row.net_salary) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template slot-scope="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template slot-scope="{ row }">
            <el-button size="small" @click="viewDetail(row.id)">查看</el-button>
            <el-button
              v-if="row.status === 'calculated'"
              size="small"
              type="success"
              @click="handleApprove(row.id)"
            >
              审批
            </el-button>
            <el-button
              v-if="row.status === 'approved'"
              size="small"
              type="warning"
              @click="openPayrollDialog(row)"
            >
              发放
            </el-button>
            <el-button
              v-if="['calculated'].includes(row.status)"
              size="small"
              type="danger"
              @click="handleDelete(row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          :current-page.sync="pagination.page"
          :page-size.sync="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total_items"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>


    <!-- 单个计算薪资对话框 -->
    <el-dialog
      title="计算薪资"
      :visible.sync="showCalculateDialog"
      width="500px"
    >
      <el-form ref="calculateForm" :model="calculateForm" label-width="100px">
        <el-form-item label="员工" prop="employee_id">
          <el-select
            v-model="calculateForm.employee_id"
            placeholder="请选择员工"
            style="width: 100%"
            filterable
          >
            <el-option
              v-for="employee in employees"
              :key="employee.id"
              :label="employee.name"
              :value="employee.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="月份" prop="month">
          <el-date-picker
            v-model="calculateForm.month"
            type="month"
            placeholder="选择月份"
            value-format="yyyy-MM"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="showCalculateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCalculate">计算</el-button>
      </span>
    </el-dialog>

    <!-- 批量计算薪资对话框 -->
    <el-dialog
      title="批量计算薪资"
      :visible.sync="showBatchCalculateDialog"
      width="500px"
    >
      <el-form ref="batchCalculateForm" :model="batchCalculateForm" label-width="100px">
        <el-form-item label="月份" prop="month">
          <el-date-picker
            v-model="batchCalculateForm.month"
            type="month"
            placeholder="选择月份"
            value-format="yyyy-MM"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="部门" prop="department_id">
          <el-select
            v-model="batchCalculateForm.department_id"
            placeholder="请选择部门（可选）"
            style="width: 100%"
            clearable
          >
            <el-option
              v-for="department in departments"
              :key="department.id"
              :label="department.name"
              :value="department.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="showBatchCalculateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleBatchCalculate">批量计算</el-button>
      </span>
    </el-dialog>

    <!-- 发放薪资对话框 -->
    <el-dialog
      title="发放薪资"
      :visible.sync="showPayrollDialog"
      width="500px"
    >
      <div v-if="selectedSalary" class="salary-info">
        <el-descriptions title="薪资信息" :column="1">
          <el-descriptions-item label="员工">{{ selectedSalary.employee?.name }}</el-descriptions-item>
          <el-descriptions-item label="月份">{{ selectedSalary.month }}</el-descriptions-item>
          <el-descriptions-item label="实发薪资">{{ formatCurrency(selectedSalary.net_salary) }}</el-descriptions-item>
        </el-descriptions>
      </div>
      <el-form ref="payrollForm" :model="payrollForm" label-width="100px">
        <el-form-item label="发放方式" prop="payment_method">
          <el-select v-model="payrollForm.payment_method" style="width: 100%">
            <el-option label="银行转账" value="bank" />
            <el-option label="现金" value="cash" />
            <el-option label="支票" value="check" />
          </el-select>
        </el-form-item>
        <el-form-item label="银行账户" prop="bank_account">
          <el-input v-model="payrollForm.bank_account" placeholder="请输入银行账户" />
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="showPayrollDialog = false">取消</el-button>
        <el-button type="primary" @click="handleProcessPayroll">发放</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { salaryService } from '@/services/salary'
import { employeeApi } from '@/services/employeeApi'
import { departmentApi } from '@/services/departmentApi'
import { formatDate } from '@/utils/date'
import PayrollWorkflow from '@/components/salary/PayrollWorkflow.vue'

export default {
  name: 'SalaryListView',
  components: {
    PayrollWorkflow
  },
  data() {
    return {
      activeTab: 'list', // 默认显示薪资列表
      selectedPeriodId: null, // 选中的薪资周期ID
      loading: false,
      showStatistics: false,
      showCalculateDialog: false,
      showBatchCalculateDialog: false,
      showPayrollDialog: false,
      salaries: [],
      employees: [],
      departments: [],
      statistics: {},
      selectedSalary: null,
      selectedSalaries: [],
      searchForm: {
        month: '',
        employee_id: '',
        status: '',
        page: 1,
        page_size: 20
      },
      calculateForm: {
        employee_id: '',
        month: ''
      },
      batchCalculateForm: {
        month: '',
        department_id: null
      },
      payrollForm: {
        payment_method: 'bank',
        bank_account: ''
      },
      pagination: {
        page: 1,
        page_size: 20,
        total_items: 0,
        total_pages: 0
      }
    }
  },
  mounted() {
    this.fetchSalaries()
    this.fetchEmployees()
    this.fetchDepartments()
    this.fetchStatistics()
  },
  methods: {
    formatDate,
    async fetchSalaries() {
      this.loading = true
      try {
        const response = await salaryService.getSalaries({
          ...this.searchForm,
          page: this.pagination.page,
          page_size: this.pagination.page_size
        })
        // 确保数据是数组格式，防止indexOf错误
        this.salaries = Array.isArray(response.data) ? response.data : []
        this.pagination.total_items = response.total_items || 0
        this.pagination.total_pages = response.total_pages || 0
      } catch (error) {
        this.$message.error('获取薪资列表失败')
        // 确保在错误情况下也设置为空数组
        this.salaries = []
      } finally {
        this.loading = false
      }
    },
    async fetchEmployees() {
      try {
        const response = await employeeApi.getAllEmployees()
        this.employees = Array.isArray(response.data) ? response.data : []
      } catch (error) {
        this.$message.error('获取员工列表失败')
        this.employees = []
      }
    },
    async fetchDepartments() {
      try {
        const response = await departmentApi.getAllDepartments()
        this.departments = Array.isArray(response.data) ? response.data : []
      } catch (error) {
        this.$message.error('获取部门列表失败')
        this.departments = []
      }
    },
    async fetchStatistics() {
      try {
        const response = await salaryService.getSalaryStatistics({
          month: this.searchForm.month
        })
        this.statistics = response.data || {}
      } catch (error) {
        this.$message.error('获取统计信息失败')
      }
    },
    handleSearch() {
      this.pagination.page = 1
      this.fetchSalaries()
      if (this.showStatistics) {
        this.fetchStatistics()
      }
    },
    handleReset() {
      Object.assign(this.searchForm, {
        month: '',
        employee_id: '',
        status: '',
        page: 1,
        page_size: 20
      })
      this.pagination.page = 1
      this.fetchSalaries()
    },
    viewDetail(id) {
      this.$router.push(`/salary/${id}`)
    },
    async handleCalculate() {
      if (!this.calculateForm.employee_id || !this.calculateForm.month) {
        this.$message.error('请填写完整信息')
        return
      }
      
      try {
        await salaryService.calculateSalary(this.calculateForm.employee_id, this.calculateForm.month)
        this.$message.success('薪资计算成功')
        this.showCalculateDialog = false
        this.fetchSalaries()
      } catch (error) {
        this.$message.error('薪资计算失败')
      }
    },
    async handleBatchCalculate() {
      if (!this.batchCalculateForm.month) {
        this.$message.error('请选择月份')
        return
      }
      
      try {
        const result = await salaryService.batchCalculateSalary(
          this.batchCalculateForm.month,
          this.batchCalculateForm.department_id,
          []
        )
        this.$message.success(`批量计算完成: 成功 ${result.success} 个，失败 ${result.failed} 个`)
        this.showBatchCalculateDialog = false
        this.fetchSalaries()
      } catch (error) {
        this.$message.error('批量计算失败')
      }
    },
    async handleApprove(id) {
      try {
        await this.$confirm('确定要审批通过该薪资记录吗？', '确认审批', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        await salaryService.approveSalary(id, 'approved', '审批通过')
        this.$message.success('审批成功')
        this.fetchSalaries()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('审批失败')
        }
      }
    },
    async handleDelete(id) {
      try {
        await this.$confirm('确定要删除这条薪资记录吗？', '确认删除', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        await salaryService.deleteSalary(id)
        this.$message.success('删除成功')
        this.fetchSalaries()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除失败')
        }
      }
    },
    openPayrollDialog(salary) {
      this.selectedSalary = salary
      this.showPayrollDialog = true
    },
    async handleProcessPayroll() {
      if (!this.selectedSalary || !this.payrollForm.payment_method || !this.payrollForm.bank_account) {
        this.$message.error('请填写完整信息')
        return
      }
      
      try {
        await salaryService.processPayroll(
          this.selectedSalary.id,
          this.payrollForm.payment_method,
          this.payrollForm.bank_account
        )
        this.$message.success('薪资发放处理成功')
        this.showPayrollDialog = false
        this.selectedSalary = null
        this.payrollForm = { payment_method: 'bank', bank_account: '' }
        this.fetchSalaries()
      } catch (error) {
        this.$message.error('薪资发放处理失败')
      }
    },
    async exportReport() {
      try {
        const response = await salaryService.exportSalaries('excel', this.searchForm)
        const blob = new Blob([response.data], {
          type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
        })
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `salary_report_${new Date().toISOString().slice(0, 10)}.xlsx`
        a.click()
        window.URL.revokeObjectURL(url)
      } catch (error) {
        this.$message.error('导出报表失败')
      }
    },
    handleSelectionChange(selected) {
      this.selectedSalaries = selected
    },
    handleSizeChange(size) {
      this.pagination.page_size = size
      this.pagination.page = 1
      this.fetchSalaries()
    },
    handleCurrentChange(page) {
      this.pagination.page = page
      this.fetchSalaries()
    },
    formatCurrency(amount) {
      return new Intl.NumberFormat('zh-CN', {
        style: 'currency',
        currency: 'CNY'
      }).format(amount || 0)
    },
    getStatusType(status) {
      const statusTypes = {
        calculated: 'info',
        approved: 'success',
        paid: 'primary',
        rejected: 'danger'
      }
      return statusTypes[status] || 'info'
    },
    getStatusText(status) {
      const statusTexts = {
        calculated: '已计算',
        approved: '已审批',
        paid: '已发放',
        rejected: '已拒绝'
      }
      return statusTexts[status] || status
    }
  }
}
</script>

<style scoped>
.salary-list {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

.page-actions {
  display: flex;
  gap: 10px;
}

.salary-tabs {
  margin-bottom: 20px;
}

.salary-tabs .el-tabs__header {
  margin-bottom: 20px;
}

.salary-list-content {
  position: relative;
}

.floating-actions {
  position: fixed;
  bottom: 30px;
  right: 30px;
  z-index: 1000;
}

.float-button {
  width: 60px;
  height: 60px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.statistics-cards {
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.search-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.salary-info {
  margin-bottom: 20px;
}

.dialog-footer {
  text-align: right;
}
</style>