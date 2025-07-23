<template>
  <div class="payroll-workflow">
    <!-- 工作流步骤 -->
    <el-card class="workflow-header">
      <div class="workflow-title">
        <h3>薪资发放流程</h3>
        <el-tag :type="getWorkflowStatusType(workflowStatus)">{{ getWorkflowStatusText(workflowStatus) }}</el-tag>
      </div>
      <el-steps :active="activeStep" finish-status="success" align-center>
        <el-step title="计算薪资" description="系统自动计算薪资"></el-step>
        <el-step title="初审" description="HR初步审核"></el-step>
        <el-step title="终审" description="管理员最终审批"></el-step>
        <el-step title="发放" description="薪资发放处理"></el-step>
        <el-step title="完成" description="发放完成确认"></el-step>
      </el-steps>
    </el-card>

    <!-- 当前步骤详情 -->
    <el-card class="current-step-detail">
      <div slot="header" class="clearfix">
        <span>当前步骤: {{ getCurrentStepTitle() }}</span>
        <el-button style="float: right; padding: 3px 0" type="text" @click="showWorkflowHistory = !showWorkflowHistory">
          {{ showWorkflowHistory ? '隐藏' : '查看' }}历史记录
        </el-button>
      </div>
      
      <!-- 步骤1: 计算薪资 -->
      <div v-if="activeStep === 0" class="step-content">
        <el-alert title="需要计算薪资" type="info" show-icon>
          <p>请先为员工计算薪资，然后提交审核。</p>
        </el-alert>
        <div class="step-actions">
          <el-button type="primary" @click="handleCalculateStep">批量计算薪资</el-button>
          <el-button @click="showCalculateDialog = true">单个计算薪资</el-button>
        </div>
      </div>

      <!-- 步骤2: 初审 -->
      <div v-if="activeStep === 1" class="step-content">
        <el-alert title="待初审" type="warning" show-icon>
          <p>薪资已计算完成，等待HR进行初步审核。</p>
        </el-alert>
        <div class="approval-section">
          <el-form :model="approvalForm" label-width="80px">
            <el-form-item label="审核意见">
              <el-input
                v-model="approvalForm.comments"
                type="textarea"
                :rows="3"
                placeholder="请输入审核意见"
              />
            </el-form-item>
          </el-form>
          <div class="step-actions">
            <el-button type="success" @click="handleApproval(true)">通过</el-button>
            <el-button type="danger" @click="handleApproval(false)">拒绝</el-button>
            <el-button @click="showBatchApprovalDialog = true">批量审批</el-button>
          </div>
        </div>
      </div>

      <!-- 步骤3: 终审 -->
      <div v-if="activeStep === 2" class="step-content">
        <el-alert title="待终审" type="warning" show-icon>
          <p>初审已通过，等待管理员进行最终审批。</p>
        </el-alert>
        <div class="approval-section" v-if="isAdmin">
          <el-form :model="finalApprovalForm" label-width="80px">
            <el-form-item label="终审意见">
              <el-input
                v-model="finalApprovalForm.comments"
                type="textarea"
                :rows="3"
                placeholder="请输入最终审批意见"
              />
            </el-form-item>
          </el-form>
          <div class="step-actions">
            <el-button type="success" @click="handleFinalApproval(true)">最终批准</el-button>
            <el-button type="danger" @click="handleFinalApproval(false)">最终拒绝</el-button>
          </div>
        </div>
        <div v-else class="waiting-approval">
          <el-alert title="等待管理员审批" type="info" show-icon />
        </div>
      </div>

      <!-- 步骤4: 发放 -->
      <div v-if="activeStep === 3" class="step-content">
        <el-alert title="准备发放" type="success" show-icon>
          <p>薪资审批已完成，可以进行发放处理。</p>
        </el-alert>
        <div class="payroll-section">
          <el-button type="primary" @click="showPayrollBatchDialog = true">创建发放批次</el-button>
          <el-button @click="showPayrollSingleDialog = true">单个发放</el-button>
          <el-button @click="viewPayrollBatches">查看发放批次</el-button>
        </div>
      </div>

      <!-- 步骤5: 完成 -->
      <div v-if="activeStep === 4" class="step-content">
        <el-alert title="发放完成" type="success" show-icon>
          <p>薪资已成功发放，流程完成。</p>
        </el-alert>
        <div class="completion-summary">
          <el-descriptions title="发放汇总" :column="2">
            <el-descriptions-item label="发放人数">{{ payrollSummary.totalEmployees }}</el-descriptions-item>
            <el-descriptions-item label="发放总额">{{ formatCurrency(payrollSummary.totalAmount) }}</el-descriptions-item>
            <el-descriptions-item label="成功发放">{{ payrollSummary.successCount }}</el-descriptions-item>
            <el-descriptions-item label="发放失败">{{ payrollSummary.failedCount }}</el-descriptions-item>
            <el-descriptions-item label="发放时间">{{ formatDate(payrollSummary.completedAt) }}</el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
    </el-card>

    <!-- 工作流历史记录 -->
    <el-card v-if="showWorkflowHistory" class="workflow-history">
      <div slot="header" class="clearfix">
        <span>工作流历史</span>
      </div>
      <el-timeline>
        <el-timeline-item
          v-for="record in workflowHistory"
          :key="record.id"
          :timestamp="formatDate(record.createdAt)"
          :type="getTimelineType(record.action)"
        >
          <el-card>
            <h4>{{ record.title }}</h4>
            <p>操作人: {{ record.operator }}</p>
            <p v-if="record.comments">备注: {{ record.comments }}</p>
            <p>状态: <el-tag :type="getActionType(record.action)">{{ record.action }}</el-tag></p>
          </el-card>
        </el-timeline-item>
      </el-timeline>
    </el-card>

    <!-- 薪资列表 -->
    <el-card class="salary-list-card">
      <div slot="header" class="clearfix">
        <span>薪资列表</span>
        <div style="float: right;">
          <el-button @click="refreshSalaryList">刷新</el-button>
          <el-dropdown @command="handleBulkAction">
            <el-button type="primary">
              批量操作<i class="el-icon-arrow-down el-icon--right"></i>
            </el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item command="bulkApprove" :disabled="selectedSalaries.length === 0">
                批量审批
              </el-dropdown-item>
              <el-dropdown-item command="bulkReject" :disabled="selectedSalaries.length === 0">
                批量拒绝
              </el-dropdown-item>
              <el-dropdown-item command="bulkPayroll" :disabled="selectedSalaries.length === 0">
                批量发放
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
      </div>
      
      <el-table
        v-loading="salaryListLoading"
        :data="salaryList"
        @selection-change="handleSelectionChange"
        stripe
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="employee.name" label="员工姓名" width="120" />
        <el-table-column prop="employee.department.name" label="部门" width="120" />
        <el-table-column prop="month" label="月份" width="100" />
        <el-table-column prop="net_salary" label="实发薪资" width="120">
          <template slot-scope="{ row }">
            {{ formatCurrency(row.net_salary) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template slot-scope="{ row }">
            <el-tag :type="getSalaryStatusType(row.status)">
              {{ getSalaryStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reviewer" label="审核人" width="100" />
        <el-table-column prop="approver" label="批准人" width="100" />
        <el-table-column label="操作" width="200" fixed="right">
          <template slot-scope="{ row }">
            <el-button size="small" @click="viewSalaryDetail(row)">详情</el-button>
            <el-button
              v-if="canApprove(row.status)"
              size="small"
              type="success"
              @click="approveSingle(row)"
            >
              审批
            </el-button>
            <el-button
              v-if="canPayroll(row.status)"
              size="small"
              type="warning"
              @click="payrollSingle(row)"
            >
              发放
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination">
        <el-pagination
          :current-page.sync="salaryPagination.page"
          :page-size.sync="salaryPagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="salaryPagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 计算薪资对话框 -->
    <el-dialog
      title="计算薪资"
      :visible.sync="showCalculateDialog"
      width="600px"
    >
      <el-form :model="calculateForm" label-width="100px">
        <el-form-item label="薪资周期">
          <el-select v-model="calculateForm.periodId" placeholder="选择薪资周期" style="width: 100%">
            <el-option
              v-for="period in payrollPeriods"
              :key="period.id"
              :label="period.name"
              :value="period.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="员工">
          <el-select v-model="calculateForm.employeeIds" multiple placeholder="选择员工" style="width: 100%">
            <el-option
              v-for="employee in employees"
              :key="employee.id"
              :label="employee.name"
              :value="employee.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="部门">
          <el-select v-model="calculateForm.departmentId" placeholder="按部门计算(可选)" style="width: 100%" clearable>
            <el-option
              v-for="dept in departments"
              :key="dept.id"
              :label="dept.name"
              :value="dept.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="showCalculateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCalculateSalary">开始计算</el-button>
      </span>
    </el-dialog>

    <!-- 批量审批对话框 -->
    <el-dialog
      title="批量审批"
      :visible.sync="showBatchApprovalDialog"
      width="500px"
    >
      <el-form :model="batchApprovalForm" label-width="100px">
        <el-form-item label="审批结果">
          <el-radio-group v-model="batchApprovalForm.approved">
            <el-radio :label="true">通过</el-radio>
            <el-radio :label="false">拒绝</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="审批意见">
          <el-input
            v-model="batchApprovalForm.comments"
            type="textarea"
            :rows="4"
            placeholder="请输入审批意见"
          />
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="showBatchApprovalDialog = false">取消</el-button>
        <el-button type="primary" @click="handleBatchApproval">确认审批</el-button>
      </span>
    </el-dialog>

    <!-- 创建发放批次对话框 -->
    <el-dialog
      title="创建发放批次"
      :visible.sync="showPayrollBatchDialog"
      width="600px"
    >
      <el-form :model="payrollBatchForm" label-width="100px">
        <el-form-item label="批次名称">
          <el-input v-model="payrollBatchForm.name" placeholder="输入批次名称" />
        </el-form-item>
        <el-form-item label="薪资周期">
          <el-select v-model="payrollBatchForm.periodId" placeholder="选择薪资周期" style="width: 100%">
            <el-option
              v-for="period in payrollPeriods"
              :key="period.id"
              :label="period.name"
              :value="period.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="计划发放日期">
          <el-date-picker
            v-model="payrollBatchForm.scheduledDate"
            type="datetime"
            placeholder="选择发放时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="包含薪资">
          <el-checkbox-group v-model="payrollBatchForm.salaryIds">
            <el-checkbox
              v-for="salary in approvedSalaries"
              :key="salary.id"
              :label="salary.id"
            >
              {{ salary.employee.name }} - {{ formatCurrency(salary.net_salary) }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="showPayrollBatchDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreatePayrollBatch">创建批次</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { salaryService } from '@/services/salary'
import { formatDate } from '@/utils/date'

export default {
  name: 'PayrollWorkflow',
  props: {
    periodId: {
      type: [Number, String],
      default: null
    }
  },
  data() {
    return {
      activeStep: 0,
      workflowStatus: 'in_progress',
      showWorkflowHistory: false,
      salaryListLoading: false,
      
      // 表单数据
      approvalForm: {
        comments: ''
      },
      finalApprovalForm: {
        comments: ''
      },
      calculateForm: {
        periodId: null,
        employeeIds: [],
        departmentId: null
      },
      batchApprovalForm: {
        approved: true,
        comments: ''
      },
      payrollBatchForm: {
        name: '',
        periodId: null,
        scheduledDate: null,
        salaryIds: []
      },
      
      // 数据列表
      salaryList: [],
      selectedSalaries: [],
      workflowHistory: [],
      payrollPeriods: [],
      employees: [],
      departments: [],
      approvedSalaries: [],
      
      // 分页
      salaryPagination: {
        page: 1,
        pageSize: 20,
        total: 0
      },
      
      // 对话框状态
      showCalculateDialog: false,
      showBatchApprovalDialog: false,
      showPayrollBatchDialog: false,
      showPayrollSingleDialog: false,
      
      // 汇总数据
      payrollSummary: {
        totalEmployees: 0,
        totalAmount: 0,
        successCount: 0,
        failedCount: 0,
        completedAt: null
      }
    }
  },
  computed: {
    isAdmin() {
      return this.$store.getters.roles.includes('admin')
    },
    isHR() {
      return this.$store.getters.roles.includes('hr') || this.isAdmin
    }
  },
  mounted() {
    this.initializeWorkflow()
  },
  methods: {
    formatDate,
    async initializeWorkflow() {
      await this.fetchSalaryList()
      await this.fetchPayrollPeriods()
      await this.fetchEmployees()
      await this.fetchDepartments()
      await this.fetchWorkflowHistory()
      this.updateWorkflowStep()
    },
    
    async fetchSalaryList() {
      this.salaryListLoading = true
      try {
        const params = {
          period_id: this.periodId,
          page: this.salaryPagination.page,
          page_size: this.salaryPagination.pageSize
        }
        const response = await salaryService.getEnhancedSalaries(params)
        this.salaryList = Array.isArray(response.data) ? response.data : []
        this.salaryPagination.total = response.total_items || 0
      } catch (error) {
        this.$message.error('获取薪资列表失败')
        this.salaryList = [] // 确保是数组
      } finally {
        this.salaryListLoading = false
      }
    },
    
    async fetchPayrollPeriods() {
      try {
        const response = await salaryService.getPayrollPeriods()
        this.payrollPeriods = Array.isArray(response.data) ? response.data : []
      } catch (error) {
        this.$message.error('获取薪资周期失败')
        this.payrollPeriods = []
      }
    },
    
    async fetchEmployees() {
      try {
        const response = await salaryService.getEmployees()
        this.employees = Array.isArray(response.data) ? response.data : []
      } catch (error) {
        this.$message.error('获取员工列表失败')
        this.employees = []
      }
    },
    
    async fetchDepartments() {
      try {
        const response = await salaryService.getDepartments()
        this.departments = Array.isArray(response.data) ? response.data : []
      } catch (error) {
        this.$message.error('获取部门列表失败')
        this.departments = []
      }
    },
    
    async fetchWorkflowHistory() {
      try {
        const response = await salaryService.getWorkflowHistory(this.periodId)
        this.workflowHistory = Array.isArray(response.data) ? response.data : []
      } catch (error) {
        console.error('获取工作流历史失败', error)
        this.workflowHistory = []
      }
    },
    
    updateWorkflowStep() {
      if (!Array.isArray(this.salaryList) || this.salaryList.length === 0) {
        this.activeStep = 0
        return
      }
      
      const hasCalculated = this.salaryList.some(s => s.status !== 'draft')
      const hasReviewed = this.salaryList.some(s => ['approved', 'paid'].includes(s.status))
      const hasApproved = this.salaryList.some(s => s.status === 'approved')
      const hasPaid = this.salaryList.every(s => s.status === 'paid')
      
      if (hasPaid) {
        this.activeStep = 4
        this.workflowStatus = 'completed'
      } else if (hasApproved) {
        this.activeStep = 3
        this.workflowStatus = 'ready_for_payroll'
      } else if (hasReviewed) {
        this.activeStep = 2
        this.workflowStatus = 'pending_final_approval'
      } else if (hasCalculated) {
        this.activeStep = 1
        this.workflowStatus = 'pending_review'
      } else {
        this.activeStep = 0
        this.workflowStatus = 'pending_calculation'
      }
    },
    
    getCurrentStepTitle() {
      const titles = ['计算薪资', '初审', '终审', '发放', '完成']
      return titles[this.activeStep] || '未知步骤'
    },
    
    // 处理计算薪资
    async handleCalculateStep() {
      this.showCalculateDialog = true
    },
    
    async handleCalculateSalary() {
      if (!this.calculateForm.periodId) {
        this.$message.error('请选择薪资周期')
        return
      }
      
      try {
        this.$loading({
          lock: true,
          text: '正在计算薪资...',
          spinner: 'el-icon-loading',
          background: 'rgba(0, 0, 0, 0.7)'
        })
        
        const params = {
          period_id: this.calculateForm.periodId,
          department_id: this.calculateForm.departmentId,
          employee_ids: this.calculateForm.employeeIds
        }
        
        const response = await salaryService.batchCalculateSalaries(params)
        
        this.$message.success(`批量计算完成: 成功 ${response.success} 个，失败 ${response.failed} 个`)
        this.showCalculateDialog = false
        await this.fetchSalaryList()
        this.updateWorkflowStep()
      } catch (error) {
        this.$message.error('计算薪资失败')
      } finally {
        this.$loading().close()
      }
    },
    
    // 处理审批
    async handleApproval(approved) {
      if (!this.approvalForm.comments && !approved) {
        this.$message.error('拒绝时必须填写审批意见')
        return
      }
      
      try {
        const selectedIds = this.selectedSalaries.map(s => s.id)
        if (selectedIds.length === 0) {
          this.$message.error('请选择要审批的薪资记录')
          return
        }
        
        await salaryService.bulkApproveSalaries({
          salary_ids: selectedIds,
          approved,
          comments: this.approvalForm.comments
        })
        
        this.$message.success(approved ? '审批通过' : '审批拒绝')
        this.approvalForm.comments = ''
        await this.fetchSalaryList()
        this.updateWorkflowStep()
      } catch (error) {
        this.$message.error('审批失败')
      }
    },
    
    // 处理最终审批
    async handleFinalApproval(approved) {
      if (!this.finalApprovalForm.comments && !approved) {
        this.$message.error('拒绝时必须填写审批意见')
        return
      }
      
      try {
        const selectedIds = this.selectedSalaries.map(s => s.id)
        if (selectedIds.length === 0) {
          this.$message.error('请选择要最终审批的薪资记录')
          return
        }
        
        await salaryService.finalApproveSalaries({
          salary_ids: selectedIds,
          approved,
          comments: this.finalApprovalForm.comments
        })
        
        this.$message.success(approved ? '最终批准' : '最终拒绝')
        this.finalApprovalForm.comments = ''
        await this.fetchSalaryList()
        this.updateWorkflowStep()
      } catch (error) {
        this.$message.error('最终审批失败')
      }
    },
    
    // 处理批量操作
    async handleBulkAction(command) {
      if (this.selectedSalaries.length === 0) {
        this.$message.error('请选择要操作的薪资记录')
        return
      }
      
      switch (command) {
        case 'bulkApprove':
          this.showBatchApprovalDialog = true
          break
        case 'bulkReject':
          this.batchApprovalForm.approved = false
          this.showBatchApprovalDialog = true
          break
        case 'bulkPayroll':
          this.showPayrollBatchDialog = true
          break
      }
    },
    
    async handleBatchApproval() {
      try {
        const selectedIds = this.selectedSalaries.map(s => s.id)
        await salaryService.bulkApproveSalaries({
          salary_ids: selectedIds,
          approved: this.batchApprovalForm.approved,
          comments: this.batchApprovalForm.comments
        })
        
        this.$message.success('批量审批完成')
        this.showBatchApprovalDialog = false
        this.batchApprovalForm = { approved: true, comments: '' }
        await this.fetchSalaryList()
        this.updateWorkflowStep()
      } catch (error) {
        this.$message.error('批量审批失败')
      }
    },
    
    // 创建发放批次
    async handleCreatePayrollBatch() {
      if (!this.payrollBatchForm.name || !this.payrollBatchForm.periodId) {
        this.$message.error('请填写完整信息')
        return
      }
      
      try {
        const params = {
          name: this.payrollBatchForm.name,
          payroll_period_id: this.payrollBatchForm.periodId,
          scheduled_date: this.payrollBatchForm.scheduledDate,
          salary_ids: this.payrollBatchForm.salaryIds
        }
        
        const response = await salaryService.createPaymentBatch(params)
        
        this.$message.success('发放批次创建成功')
        this.showPayrollBatchDialog = false
        this.payrollBatchForm = { name: '', periodId: null, scheduledDate: null, salaryIds: [] }
        
        // 自动处理批次
        await this.processPayrollBatch(response.id)
      } catch (error) {
        this.$message.error('创建发放批次失败')
      }
    },
    
    async processPayrollBatch(batchId) {
      try {
        this.$loading({
          lock: true,
          text: '正在处理发放批次...',
          spinner: 'el-icon-loading',
          background: 'rgba(0, 0, 0, 0.7)'
        })
        
        const response = await salaryService.processPaymentBatch(batchId)
        
        this.$message.success('发放批次处理完成')
        await this.fetchSalaryList()
        this.updateWorkflowStep()
        this.updatePayrollSummary(response)
      } catch (error) {
        this.$message.error('处理发放批次失败')
      } finally {
        this.$loading().close()
      }
    },
    
    updatePayrollSummary(batchResult) {
      this.payrollSummary = {
        totalEmployees: batchResult.total_records || 0,
        totalAmount: batchResult.total_amount || 0,
        successCount: batchResult.success_records || 0,
        failedCount: batchResult.failed_records || 0,
        completedAt: batchResult.completed_date || new Date()
      }
    },
    
    // 工具方法
    handleSelectionChange(selection) {
      this.selectedSalaries = selection
    },
    
    refreshSalaryList() {
      this.fetchSalaryList()
    },
    
    handleSizeChange(size) {
      this.salaryPagination.pageSize = size
      this.salaryPagination.page = 1
      this.fetchSalaryList()
    },
    
    handleCurrentChange(page) {
      this.salaryPagination.page = page
      this.fetchSalaryList()
    },
    
    formatCurrency(amount) {
      return new Intl.NumberFormat('zh-CN', {
        style: 'currency',
        currency: 'CNY'
      }).format(amount || 0)
    },
    
    // 状态相关方法
    getWorkflowStatusType(status) {
      const types = {
        pending_calculation: 'info',
        pending_review: 'warning',
        pending_final_approval: 'warning',
        ready_for_payroll: 'success',
        in_progress: 'primary',
        completed: 'success'
      }
      return types[status] || 'info'
    },
    
    getWorkflowStatusText(status) {
      const texts = {
        pending_calculation: '待计算',
        pending_review: '待初审',
        pending_final_approval: '待终审',
        ready_for_payroll: '待发放',
        in_progress: '进行中',
        completed: '已完成'
      }
      return texts[status] || status
    },
    
    getSalaryStatusType(status) {
      const types = {
        draft: 'info',
        calculated: 'warning',
        reviewed: 'primary',
        approved: 'success',
        paid: 'success',
        rejected: 'danger'
      }
      return types[status] || 'info'
    },
    
    getSalaryStatusText(status) {
      const texts = {
        draft: '草稿',
        calculated: '已计算',
        reviewed: '已初审',
        approved: '已批准',
        paid: '已发放',
        rejected: '已拒绝'
      }
      return texts[status] || status
    },
    
    getTimelineType(action) {
      const types = {
        calculated: 'primary',
        reviewed: 'warning',
        approved: 'success',
        paid: 'success',
        rejected: 'danger'
      }
      return types[action] || 'primary'
    },
    
    getActionType(action) {
      const types = {
        calculated: 'primary',
        reviewed: 'warning',
        approved: 'success',
        paid: 'success',
        rejected: 'danger'
      }
      return types[action] || 'info'
    },
    
    canApprove(status) {
      return this.isHR && ['calculated', 'reviewed'].includes(status)
    },
    
    canPayroll(status) {
      return this.isHR && status === 'approved'
    },
    
    // 单个操作
    async approveSingle(salary) {
      try {
        await this.$confirm(`确定要审批员工 ${salary.employee.name} 的薪资吗？`, '确认审批')
        
        await salaryService.approveSalary(salary.id, {
          approved: true,
          comments: '单个审批通过'
        })
        
        this.$message.success('审批成功')
        await this.fetchSalaryList()
        this.updateWorkflowStep()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('审批失败')
        }
      }
    },
    
    async payrollSingle(salary) {
      try {
        await this.$confirm(`确定要发放员工 ${salary.employee.name} 的薪资吗？`, '确认发放')
        
        await salaryService.processPayrollSingle(salary.id, {
          payment_method: 'bank_transfer',
          bank_account: salary.employee.bank_account || ''
        })
        
        this.$message.success('发放成功')
        await this.fetchSalaryList()
        this.updateWorkflowStep()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('发放失败')
        }
      }
    },
    
    viewSalaryDetail(salary) {
      this.$router.push(`/salary/detail/${salary.id}`)
    },
    
    viewPayrollBatches() {
      this.$router.push('/salary/payroll-batches')
    }
  }
}
</script>

<style scoped>
.payroll-workflow {
  padding: 20px;
}

.workflow-header {
  margin-bottom: 20px;
}

.workflow-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.workflow-title h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.current-step-detail {
  margin-bottom: 20px;
}

.step-content {
  padding: 20px 0;
}

.step-actions {
  margin-top: 20px;
}

.step-actions .el-button {
  margin-right: 10px;
}

.approval-section {
  margin-top: 20px;
}

.payroll-section {
  margin-top: 20px;
}

.payroll-section .el-button {
  margin-right: 10px;
}

.waiting-approval {
  text-align: center;
  padding: 20px;
}

.completion-summary {
  margin-top: 20px;
}

.workflow-history {
  margin-bottom: 20px;
}

.salary-list-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.dialog-footer {
  text-align: right;
}

.clearfix:before,
.clearfix:after {
  content: "";
  display: table;
}

.clearfix:after {
  clear: both;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .payroll-workflow {
    padding: 10px;
  }
  
  .workflow-title {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .workflow-title h3 {
    margin-bottom: 10px;
  }
  
  .step-actions .el-button {
    display: block;
    margin-bottom: 10px;
    margin-right: 0;
    width: 100%;
  }
}
</style>