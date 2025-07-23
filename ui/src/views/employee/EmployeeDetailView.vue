<template>
  <div class="employee-detail">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-title">
        <div class="title-icon">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 12C14.21 12 16 10.21 16 8C16 5.79 14.21 4 12 4C9.79 4 8 5.79 8 8C8 10.21 9.79 12 12 12ZM12 14C9.33 14 4 15.34 4 18V20H20V18C20 15.34 14.67 14 12 14Z" fill="currentColor"/>
          </svg>
        </div>
        <div class="title-content">
          <h1>{{ employee ? employee.name : '员工详情' }}</h1>
          <p>查看员工详细信息</p>
        </div>
      </div>
      <div class="page-actions">
        <el-button @click="$router.back()" class="action-btn">
          <i class="el-icon-arrow-left"></i>
          <span>返回</span>
        </el-button>
        <el-button type="primary" @click="editEmployee" class="action-btn primary-btn" v-if="employee">
          <i class="el-icon-edit"></i>
          <span>编辑</span>
        </el-button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-card>
        <div class="loading-content">
          <i class="el-icon-loading"></i>
          <span>加载中...</span>
        </div>
      </el-card>
    </div>

    <!-- 员工信息内容 -->
    <div v-else-if="employee" class="employee-content">
      <el-row :gutter="24">
        <!-- 左侧主要信息 -->
        <el-col :span="18">
          <!-- 员工基本信息卡片 -->
          <el-card class="info-card">
            <div slot="header" class="card-header">
              <div class="employee-summary">
                <el-avatar :size="60" :src="employee.avatar">
                  {{ (employee.name && employee.name.charAt(0).toUpperCase()) || '' }}
                </el-avatar>
                <div class="summary-info">
                  <h2>{{ employee.name }}</h2>
                  <p class="position">{{ employee.position ? employee.position.name : '未设置职位' }}</p>
                  <p class="department">
                    <i class="el-icon-office-building"></i>
                    {{ employee.department ? employee.department.name : '未分配部门' }}
                  </p>
                  <el-tag :type="employee.status === 'active' ? 'success' : 'danger'" size="medium">
                    {{ employee.status === 'active' ? '在职' : '离职' }}
                  </el-tag>
                </div>
              </div>
            </div>

            <!-- 基本信息表格 -->
            <div class="info-section">
              <h3>基本信息</h3>
              <el-table :data="basicInfoData" :show-header="false" border>
                <el-table-column prop="label" width="120" class-name="label-column">
                  <template slot-scope="scope">
                    <strong>{{ scope.row.label }}</strong>
                  </template>
                </el-table-column>
                <el-table-column prop="value" class-name="value-column">
                  <template slot-scope="scope">
                    <span v-html="scope.row.value"></span>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-card>

          <!-- 工作信息 -->
          <el-card class="info-card">
            <div slot="header" class="card-header">
              <span><i class="el-icon-suitcase"></i> 工作信息</span>
            </div>
            <div class="info-section">
              <el-table :data="workInfoData" :show-header="false" border>
                <el-table-column prop="label" width="120" class-name="label-column">
                  <template slot-scope="scope">
                    <strong>{{ scope.row.label }}</strong>
                  </template>
                </el-table-column>
                <el-table-column prop="value" class-name="value-column">
                  <template slot-scope="scope">
                    <span v-html="scope.row.value"></span>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-card>

          <!-- 联系信息 -->
          <el-card class="info-card" v-if="hasContactInfo">
            <div slot="header" class="card-header">
              <span><i class="el-icon-phone"></i> 联系信息</span>
            </div>
            <div class="info-section">
              <el-table :data="contactInfoData" :show-header="false" border>
                <el-table-column prop="label" width="120" class-name="label-column">
                  <template slot-scope="scope">
                    <strong>{{ scope.row.label }}</strong>
                  </template>
                </el-table-column>
                <el-table-column prop="value" class-name="value-column">
                  <template slot-scope="scope">
                    <span v-html="scope.row.value"></span>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-card>

          <!-- 教育背景 -->
          <el-card class="info-card" v-if="hasEducationInfo">
            <div slot="header" class="card-header">
              <span><i class="el-icon-school"></i> 教育背景</span>
            </div>
            <div class="info-section">
              <el-table :data="educationInfoData" :show-header="false" border>
                <el-table-column prop="label" width="120" class-name="label-column">
                  <template slot-scope="scope">
                    <strong>{{ scope.row.label }}</strong>
                  </template>
                </el-table-column>
                <el-table-column prop="value" class-name="value-column">
                  <template slot-scope="scope">
                    <span v-html="scope.row.value"></span>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-card>
        </el-col>

        <!-- 右侧操作面板 -->
        <el-col :span="6">
          <!-- 快捷操作 -->
          <el-card class="actions-card">
            <div slot="header" class="card-header">
              <span><i class="el-icon-setting"></i> 快捷操作</span>
            </div>
            <div class="quick-actions">
              <el-button type="primary" @click="editEmployee" block class="action-button">
                <i class="el-icon-edit"></i>
                编辑员工
              </el-button>
              <el-button @click="viewAttendance" block class="action-button">
                <i class="el-icon-time"></i>
                查看考勤
              </el-button>
              <el-button @click="viewSalary" block class="action-button">
                <i class="el-icon-money"></i>
                查看薪资
              </el-button>
              <el-button @click="viewPerformance" block class="action-button">
                <i class="el-icon-data-analysis"></i>
                查看绩效
              </el-button>
              <el-button @click="changeStatus" block class="action-button">
                <i class="el-icon-switch-button"></i>
                修改状态
              </el-button>
            </div>
          </el-card>

          <!-- 统计信息 -->
          <el-card class="stats-card">
            <div slot="header" class="card-header">
              <span><i class="el-icon-data-line"></i> 统计信息</span>
            </div>
            <div class="stats-content">
              <div class="stat-item">
                <div class="stat-value">{{ workDays }}</div>
                <div class="stat-label">工作天数</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ employee.base_salary ? `¥${employee.base_salary.toLocaleString()}` : '-' }}</div>
                <div class="stat-label">基础薪资</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ employee.job_level ? employee.job_level.name : '-' }}</div>
                <div class="stat-label">当前职级</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 错误状态 -->
    <div v-else class="error-container">
      <el-card>
        <div class="error-content">
          <el-result
            icon="warning"
            title="员工不存在"
            sub-title="抱歉，您访问的员工信息不存在或已被删除。"
          >
            <template slot="extra">
              <el-button type="primary" @click="$router.push('/employees')">
                返回员工列表
              </el-button>
            </template>
          </el-result>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script>
import { employeeApiService } from '@/services/employeeApi'
import dayjs from 'dayjs'

export default {
  name: 'EmployeeDetailView',
  data() {
    return {
      loading: false,
      employee: null
    }
  },
  computed: {
    basicInfoData() {
      if (!this.employee) return []
      return [
        { label: '工号', value: this.employee.employee_id || '-' },
        { label: '姓名', value: this.employee.name || '-' },
        { label: '邮箱', value: this.employee.email || '-' },
        { label: '手机号', value: this.employee.phone || '-' },
        { label: '性别', value: this.formatGender(this.employee.gender) },
        { label: '生日', value: this.formatDate(this.employee.birthday) },
        { label: '身份证号', value: this.employee.id_card || '-' },
      ]
    },
    workInfoData() {
      if (!this.employee) return []
      return [
        { label: '部门', value: this.employee.department ? this.employee.department.name : '-' },
        { label: '职位', value: this.employee.position ? this.employee.position.name : '-' },
        { label: '职级', value: this.employee.job_level ? this.employee.job_level.name : '-' },
        { label: '直接上级', value: this.employee.manager ? this.employee.manager.name : '-' },
        { label: '入职日期', value: this.formatDate(this.employee.hire_date) },
        { label: '试用期结束', value: this.formatDate(this.employee.probation_end_date) },
        { label: '合同类型', value: this.formatContractType(this.employee.contract_type) },
        { label: '基本薪资', value: this.employee.base_salary ? `¥${this.employee.base_salary.toLocaleString()}` : '-' },
      ]
    },
    contactInfoData() {
      if (!this.employee) return []
      return [
        { label: '地址', value: this.employee.address || '-' },
        { label: '紧急联系人', value: this.employee.emergency_contact || '-' },
        { label: '紧急联系电话', value: this.employee.emergency_phone || '-' },
      ]
    },
    educationInfoData() {
      if (!this.employee) return []
      return [
        { label: '学历', value: this.formatEducation(this.employee.education) },
        { label: '毕业学校', value: this.employee.school || '-' },
        { label: '专业', value: this.employee.major || '-' },
      ]
    },
    hasContactInfo() {
      return this.employee && (this.employee.address || this.employee.emergency_contact || this.employee.emergency_phone)
    },
    hasEducationInfo() {
      return this.employee && (this.employee.education || this.employee.school || this.employee.major)
    },
    workDays() {
      if (!this.employee || !this.employee.hire_date) return '-'
      const hireDate = dayjs(this.employee.hire_date)
      const today = dayjs()
      return today.diff(hireDate, 'day') + '天'
    }
  },
  async mounted() {
    await this.fetchEmployee()
  },
  methods: {
    async fetchEmployee() {
      const id = this.$route.params.id
      if (!id) {
        this.$message.error('员工ID无效')
        return
      }

      this.loading = true
      try {
        const response = await employeeApiService.getEmployee(id)
        const responseData = response.data || response
        this.employee = responseData.data || responseData
      } catch (error) {
        console.error('获取员工详情失败:', error)
        this.$message.error('获取员工信息失败: ' + (error.response?.data?.message || error.message || '未知错误'))
      } finally {
        this.loading = false
      }
    },
    editEmployee() {
      this.$router.push(`/employees/${this.employee.id}/edit`)
    },
    viewAttendance() {
      this.$message.info('考勤功能开发中')
    },
    viewSalary() {
      this.$message.info('薪资功能开发中')
    },
    viewPerformance() {
      this.$message.info('绩效功能开发中')
    },
    changeStatus() {
      this.$message.info('状态修改功能开发中')
    },
    formatDate(date) {
      return date ? dayjs(date).format('YYYY-MM-DD') : '-'
    },
    formatGender(gender) {
      const genderMap = {
        male: '男',
        female: '女'
      }
      return genderMap[gender] || '-'
    },
    formatContractType(type) {
      const typeMap = {
        full_time: '全职',
        part_time: '兼职',
        intern: '实习',
        contract: '合同工',
        temporary: '临时工'
      }
      return typeMap[type] || '-'
    },
    formatEducation(education) {
      const educationMap = {
        primary: '小学',
        junior: '初中',
        senior: '高中',
        associate: '大专',
        bachelor: '本科',
        master: '硕士',
        doctorate: '博士'
      }
      return educationMap[education] || '-'
    }
  }
}
</script>

<style scoped>
.employee-detail {
  padding: 0;
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  padding: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.2);
}

.page-title {
  display: flex;
  align-items: center;
  gap: 16px;
}

.title-icon {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.title-content h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 600;
  color: white;
}

.title-content p {
  margin: 8px 0 0;
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}

.page-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  padding: 12px 24px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.1);
  color: white;
  backdrop-filter: blur(10px);
  transition: all 0.3s;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.primary-btn {
  background: rgba(255, 255, 255, 0.95);
  color: #667eea;
  border: none;
}

.primary-btn:hover {
  background: white;
  color: #667eea;
}

/* 加载和错误状态 */
.loading-container,
.error-container {
  margin-top: 24px;
}

.loading-content,
.error-content {
  text-align: center;
  padding: 60px 0;
}

.loading-content i {
  font-size: 32px;
  color: #409eff;
  margin-right: 12px;
}

.loading-content span {
  font-size: 16px;
  color: #606266;
}

/* 信息卡片 */
.info-card,
.actions-card,
.stats-card {
  margin-bottom: 24px;
  border-radius: 12px;
  border: none;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #303133;
}

.card-header i {
  color: #409eff;
}

/* 员工摘要 */
.employee-summary {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 20px 0;
}

.summary-info h2 {
  margin: 0 0 8px;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.summary-info .position {
  margin: 0 0 8px;
  font-size: 16px;
  color: #409eff;
  font-weight: 500;
}

.summary-info .department {
  margin: 0 0 12px;
  font-size: 14px;
  color: #606266;
  display: flex;
  align-items: center;
  gap: 6px;
}

/* 信息表格 */
.info-section {
  margin-top: 20px;
}

.info-section h3 {
  margin: 0 0 16px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  padding-bottom: 8px;
  border-bottom: 2px solid #f0f2f5;
}

.info-section :deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;
}

.info-section :deep(.label-column) {
  background: #fafafa;
  font-weight: 600;
  color: #303133;
}

.info-section :deep(.value-column) {
  color: #606266;
}

/* 快捷操作 */
.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-button {
  justify-content: flex-start;
  padding: 12px 16px;
  border-radius: 8px;
  transition: all 0.3s;
}

.action-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.action-button i {
  margin-right: 8px;
}

/* 统计信息 */
.stats-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.stat-item {
  text-align: center;
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
  border-left: 4px solid #409eff;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #909399;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .page-actions {
    justify-content: center;
  }
  
  .employee-summary {
    flex-direction: column;
    text-align: center;
  }
  
  .info-section :deep(.el-table) {
    font-size: 14px;
  }
  
  .action-button {
    padding: 10px 12px;
  }
}
</style>