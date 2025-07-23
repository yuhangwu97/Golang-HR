<template>
  <div class="employee-edit">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-title">
        <div class="title-icon">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 12C14.21 12 16 10.21 16 8C16 5.79 14.21 4 12 4C9.79 4 8 5.79 8 8C8 10.21 9.79 12 12 12ZM12 14C9.33 14 4 15.34 4 18V20H20V18C20 15.34 14.67 14 12 14Z" fill="currentColor"/>
          </svg>
        </div>
        <div class="title-content">
          <h1>编辑员工</h1>
          <p>修改员工档案信息</p>
        </div>
      </div>
      <div class="page-actions">
        <el-button @click="$router.back()" class="action-btn">
          <i class="el-icon-arrow-left"></i>
          <span>返回</span>
        </el-button>
        <el-button type="primary" @click="handleSubmit" class="action-btn primary-btn" :loading="loading">
          <i class="el-icon-check"></i>
          <span>保存</span>
        </el-button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading && !formData.id" class="loading-container">
      <el-card>
        <div class="loading-content">
          <i class="el-icon-loading"></i>
          <span>加载中...</span>
        </div>
      </el-card>
    </div>

    <!-- 表单内容 -->
    <div v-else-if="formData.id" class="form-container">
      <el-card>
        <el-form
          ref="formRef"
          :model="formData"
          :rules="rules"
          label-position="top"
          @submit.native.prevent="handleSubmit"
        >
          <!-- 基本信息 -->
          <div class="section-title">
            <h3><i class="el-icon-user"></i> 基本信息</h3>
          </div>
          
          <el-row :gutter="24">
            <el-col :span="8">
              <el-form-item label="姓名" prop="name" required>
                <el-input 
                  v-model="formData.name" 
                  placeholder="请输入姓名"
                  clearable
                />
              </el-form-item>
            </el-col>
            
            <el-col :span="8">
              <el-form-item label="邮箱" prop="email" required>
                <el-input 
                  v-model="formData.email" 
                  placeholder="请输入邮箱"
                  clearable
                />
              </el-form-item>
            </el-col>
            
            <el-col :span="8">
              <el-form-item label="手机号" prop="phone">
                <el-input 
                  v-model="formData.phone" 
                  placeholder="请输入手机号"
                  clearable
                />
              </el-form-item>
            </el-col>
            
            <el-col :span="8">
              <el-form-item label="性别" prop="gender">
                <el-select 
                  v-model="formData.gender" 
                  placeholder="请选择性别"
                  clearable
                >
                  <el-option label="男" value="male" />
                  <el-option label="女" value="female" />
                </el-select>
              </el-form-item>
            </el-col>
            
            <el-col :span="8">
              <el-form-item label="生日" prop="birthday">
                <el-date-picker
                  v-model="formData.birthday"
                  type="date"
                  placeholder="请选择生日"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            
            <el-col :span="8">
              <el-form-item label="状态" prop="status">
                <el-select 
                  v-model="formData.status" 
                  placeholder="请选择状态"
                >
                  <el-option label="在职" value="active" />
                  <el-option label="离职" value="inactive" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <!-- 工作信息 -->
          <div class="section-title">
            <h3><i class="el-icon-suitcase"></i> 工作信息</h3>
          </div>
          
          <el-row :gutter="24">
            <el-col :span="8">
              <el-form-item label="部门" prop="department_id">
                <el-select 
                  v-model="formData.department_id" 
                  placeholder="请选择部门"
                  clearable
                  filterable
                >
                  <el-option
                    v-for="dept in departments"
                    :key="dept.id"
                    :label="dept.name"
                    :value="dept.id"
                  />
                </el-select>
              </el-form-item>
            </el-col>
            
            <el-col :span="8">
              <el-form-item label="职位" prop="position_id">
                <el-select 
                  v-model="formData.position_id" 
                  placeholder="请选择职位"
                  clearable
                  filterable
                >
                  <el-option
                    v-for="position in positions"
                    :key="position.id"
                    :label="position.name"
                    :value="position.id"
                  />
                </el-select>
              </el-form-item>
            </el-col>
            
            <el-col :span="8">
              <el-form-item label="职级" prop="job_level_id">
                <el-select 
                  v-model="formData.job_level_id" 
                  placeholder="请选择职级"
                  clearable
                  filterable
                >
                  <el-option
                    v-for="level in jobLevels"
                    :key="level.id"
                    :label="level.name"
                    :value="level.id"
                  />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-row :gutter="24">
            <el-col :span="8">
              <el-form-item label="直属领导" prop="manager_id">
                <el-select 
                  v-model="formData.manager_id" 
                  placeholder="请选择直属领导"
                  clearable
                  filterable
                >
                  <el-option
                    v-for="manager in availableManagers"
                    :key="manager.id"
                    :label="`${manager.name} (${manager.employee_id})`"
                    :value="manager.id"
                  />
                </el-select>
              </el-form-item>
            </el-col>
            
            <el-col :span="8">
              <el-form-item label="入职日期" prop="hire_date">
                <el-date-picker
                  v-model="formData.hire_date"
                  type="date"
                  placeholder="请选择入职日期"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            
            <el-col :span="8">
              <el-form-item label="基本薪资" prop="base_salary">
                <el-input-number
                  v-model="formData.base_salary"
                  placeholder="请输入基本薪资"
                  :min="0"
                  :precision="2"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>
            
            <el-col :span="8">
              <el-form-item label="合同类型" prop="contract_type">
                <el-select 
                  v-model="formData.contract_type" 
                  placeholder="请选择合同类型"
                  clearable
                >
                  <el-option label="全职" value="full_time" />
                  <el-option label="兼职" value="part_time" />
                  <el-option label="实习" value="intern" />
                  <el-option label="合同工" value="contract" />
                  <el-option label="临时工" value="temporary" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <!-- 联系信息 -->
          <div class="section-title">
            <h3><i class="el-icon-phone"></i> 联系信息</h3>
          </div>
          
          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="地址" prop="address">
                <el-input 
                  v-model="formData.address" 
                  placeholder="请输入地址"
                  type="textarea"
                  :rows="3"
                />
              </el-form-item>
            </el-col>
            
            <el-col :span="6">
              <el-form-item label="紧急联系人" prop="emergency_contact">
                <el-input 
                  v-model="formData.emergency_contact" 
                  placeholder="请输入紧急联系人"
                  clearable
                />
              </el-form-item>
            </el-col>
            
            <el-col :span="6">
              <el-form-item label="紧急联系电话" prop="emergency_phone">
                <el-input 
                  v-model="formData.emergency_phone" 
                  placeholder="请输入紧急联系电话"
                  clearable
                />
              </el-form-item>
            </el-col>
          </el-row>

          <!-- 教育背景 -->
          <div class="section-title">
            <h3><i class="el-icon-school"></i> 教育背景</h3>
          </div>
          
          <el-row :gutter="24">
            <el-col :span="8">
              <el-form-item label="学历" prop="education">
                <el-select 
                  v-model="formData.education" 
                  placeholder="请选择学历"
                  clearable
                >
                  <el-option label="小学" value="primary" />
                  <el-option label="初中" value="junior" />
                  <el-option label="高中" value="senior" />
                  <el-option label="大专" value="associate" />
                  <el-option label="本科" value="bachelor" />
                  <el-option label="硕士" value="master" />
                  <el-option label="博士" value="doctorate" />
                </el-select>
              </el-form-item>
            </el-col>
            
            <el-col :span="8">
              <el-form-item label="毕业学校" prop="school">
                <el-input 
                  v-model="formData.school" 
                  placeholder="请输入毕业学校"
                  clearable
                />
              </el-form-item>
            </el-col>
            
            <el-col :span="8">
              <el-form-item label="专业" prop="major">
                <el-input 
                  v-model="formData.major" 
                  placeholder="请输入专业"
                  clearable
                />
              </el-form-item>
            </el-col>
          </el-row>

          <!-- 表单按钮 -->
          <div class="form-actions">
            <el-button @click="$router.back()" size="large">
              <i class="el-icon-arrow-left"></i>
              返回
            </el-button>
            <el-button type="primary" @click="handleSubmit" size="large" :loading="loading">
              <i class="el-icon-check"></i>
              保存
            </el-button>
          </div>
        </el-form>
      </el-card>
    </div>

    <!-- 错误状态 -->
    <div v-else class="error-container">
      <el-card>
        <div class="error-content">
          <el-result
            icon="warning"
            title="员工不存在"
            sub-title="抱歉，您要编辑的员工信息不存在或已被删除。"
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
import { departmentApi } from '@/services/departmentApi'
import { positionApi } from '@/services/positionApi'
import { jobLevelApi } from '@/services/jobLevelApi'
import dayjs from 'dayjs'

export default {
  name: 'EmployeeEditView',
  data() {
    return {
      loading: false,
      formData: {
        id: null,
        name: '',
        email: '',
        phone: '',
        gender: '',
        birthday: null,
        status: 'active',
        department_id: null,
        position_id: null,
        job_level_id: null,
        manager_id: null,
        hire_date: null,
        base_salary: null,
        contract_type: '',
        address: '',
        emergency_contact: '',
        emergency_phone: '',
        education: '',
        school: '',
        major: ''
      },
      departments: [],
      positions: [],
      jobLevels: [],
      availableManagers: [],
      rules: {
        name: [
          { required: true, message: '请输入姓名', trigger: 'blur' },
          { min: 2, max: 50, message: '姓名长度在 2 到 50 个字符', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '请输入邮箱', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
        ],
        phone: [
          { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
        ],
        status: [
          { required: true, message: '请选择状态', trigger: 'change' }
        ]
      }
    }
  },
  async mounted() {
    await this.fetchEmployee()
    await this.fetchSelectOptions()
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
        const employee = responseData.data || responseData

        // 格式化数据
        this.formData = {
          id: employee.id,
          name: employee.name || '',
          email: employee.email || '',
          phone: employee.phone || '',
          gender: employee.gender || '',
          birthday: employee.birthday ? new Date(employee.birthday) : null,
          status: employee.status || 'active',
          department_id: employee.department_id || null,
          position_id: employee.position_id || null,
          job_level_id: employee.job_level_id || null,
          manager_id: employee.manager_id || null,
          hire_date: employee.hire_date ? new Date(employee.hire_date) : null,
          base_salary: employee.base_salary || null,
          contract_type: employee.contract_type || '',
          address: employee.address || '',
          emergency_contact: employee.emergency_contact || '',
          emergency_phone: employee.emergency_phone || '',
          education: employee.education || '',
          school: employee.school || '',
          major: employee.major || ''
        }
      } catch (error) {
        console.error('获取员工详情失败:', error)
        this.$message.error('获取员工信息失败: ' + (error.response?.data?.message || error.message || '未知错误'))
      } finally {
        this.loading = false
      }
    },
    async fetchSelectOptions() {
      try {
        // 并行获取选项数据
        const [deptResponse, posResponse, levelResponse, empResponse] = await Promise.all([
          departmentApi.getDepartments(),
          positionApi.getPositions(),
          jobLevelApi.getJobLevels(),
          employeeApiService.getAllEmployees()
        ])

        // 处理部门数据
        const deptData = deptResponse.data || deptResponse
        this.departments = deptData.data || deptData || []

        // 处理职位数据
        const posData = posResponse.data || posResponse
        this.positions = posData.data || posData || []

        // 处理职级数据
        const levelData = levelResponse.data || levelResponse
        this.jobLevels = levelData.data || levelData || []

        // 处理员工数据，用于选择直属领导
        const empData = empResponse.data || empResponse
        const allEmployees = empData.data || empData || []
        // 排除自己作为直属领导选项
        this.availableManagers = allEmployees.filter(emp => emp.id !== this.formData.id)
      } catch (error) {
        console.error('获取选项数据失败:', error)
        this.$message.error('获取选项数据失败')
      }
    },
    async handleSubmit() {
      try {
        // 表单验证
        await this.$refs.formRef.validate()
        
        this.loading = true
        
        // 准备提交数据
        const submitData = { ...this.formData }
        
        // 格式化日期
        if (submitData.birthday) {
          submitData.birthday = dayjs(submitData.birthday).format('YYYY-MM-DD')
        }
        if (submitData.hire_date) {
          submitData.hire_date = dayjs(submitData.hire_date).format('YYYY-MM-DD')
        }
        
        // 提交更新
        await employeeApiService.updateEmployee(this.formData.id, submitData)
        
        this.$message.success('员工信息更新成功')
        this.$router.push('/employees')
      } catch (error) {
        console.error('更新员工失败:', error)
        this.$message.error('更新失败: ' + (error.response?.data?.message || error.message || '未知错误'))
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.employee-edit {
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

/* 表单容器 */
.form-container {
  margin-top: 24px;
}

.form-container :deep(.el-card) {
  border-radius: 12px;
  border: none;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.form-container :deep(.el-card__body) {
  padding: 32px;
}

/* 区域标题 */
.section-title {
  margin: 32px 0 24px;
  padding-bottom: 16px;
  border-bottom: 2px solid #f0f2f5;
}

.section-title:first-child {
  margin-top: 0;
}

.section-title h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-title h3 i {
  color: #409eff;
  font-size: 20px;
}

/* 表单项样式 */
.form-container :deep(.el-form-item) {
  margin-bottom: 24px;
}

.form-container :deep(.el-form-item__label) {
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.form-container :deep(.el-input__inner),
.form-container :deep(.el-textarea__inner) {
  border-radius: 8px;
  border: 1px solid #dcdfe6;
  transition: all 0.3s;
}

.form-container :deep(.el-input__inner:focus),
.form-container :deep(.el-textarea__inner:focus) {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
}

.form-container :deep(.el-select) {
  width: 100%;
}

.form-container :deep(.el-date-editor.el-input) {
  width: 100%;
}

/* 表单按钮 */
.form-actions {
  margin-top: 48px;
  text-align: center;
  padding-top: 24px;
  border-top: 1px solid #f0f2f5;
}

.form-actions .el-button {
  margin: 0 8px;
  padding: 12px 32px;
  border-radius: 8px;
  font-weight: 600;
  transition: all 0.3s;
}

.form-actions .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
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
  
  .form-container :deep(.el-card__body) {
    padding: 20px;
  }
  
  .section-title {
    margin: 24px 0 16px;
  }
  
  .form-actions .el-button {
    display: block;
    width: 100%;
    margin: 8px 0;
  }
}
</style>