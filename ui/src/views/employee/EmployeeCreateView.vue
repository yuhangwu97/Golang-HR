<template>
  <div class="employee-create">
    <!-- 页面头部 -->
    <div class="page-header slide-up">
      <div class="page-title">
        <div class="title-icon">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 12C14.21 12 16 10.21 16 8C16 5.79 14.21 4 12 4C9.79 4 8 5.79 8 8C8 10.21 9.79 12 12 12ZM12 14C9.33 14 4 15.34 4 18V20H20V18C20 15.34 14.67 14 12 14Z" fill="currentColor"/>
          </svg>
        </div>
        <div class="title-content">
          <h1>添加员工</h1>
          <p>创建新员工档案信息</p>
        </div>
      </div>
      <div class="page-actions">
        <el-button @click="handleCancel" class="action-btn btn-animate">
          <i class="el-icon-arrow-left"></i>
          <span>返回</span>
        </el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading" class="action-btn btn-animate primary-btn">
          <i class="el-icon-check"></i>
          <span>保存</span>
        </el-button>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="create-content fade-in">
      <el-form
        ref="employeeForm"
        :model="formData"
        :rules="rules"
        label-width="120px"
        class="employee-form"
        @submit.native.prevent="handleSubmit"
      >
        <!-- 步骤进度 -->
        <div class="form-steps">
          <el-steps :active="currentStep" finish-status="success" align-center>
            <el-step title="基本信息" icon="el-icon-user"></el-step>
            <el-step title="职位信息" icon="el-icon-office-building"></el-step>
            <el-step title="入职信息" icon="el-icon-date"></el-step>
            <el-step title="其他信息" icon="el-icon-more"></el-step>
          </el-steps>
        </div>

        <!-- 表单内容 -->
        <div class="form-content">
          <!-- 第一步：基本信息 -->
          <div v-show="currentStep === 0" class="form-step zoom-in">
            <div class="step-header">
              <h3>基本信息</h3>
              <p>请填写员工的基础个人信息</p>
            </div>
            
            <!-- 头像上传 -->
            <div class="avatar-section">
              <el-upload
                class="avatar-uploader"
                action="#"
                :show-file-list="false"
                :before-upload="beforeAvatarUpload"
                :http-request="handleAvatarUpload"
              >
                <img v-if="formData.avatar" :src="formData.avatar" class="avatar">
                <div v-else class="avatar-placeholder">
                  <i class="el-icon-plus avatar-uploader-icon"></i>
                  <div class="avatar-text">上传头像</div>
                </div>
              </el-upload>
            </div>

            <el-row :gutter="24">
              <el-col :span="8">
                <el-form-item label="姓名" prop="name" required>
                  <el-input
                    v-model="formData.name"
                    placeholder="请输入员工姓名"
                    prefix-icon="el-icon-user"
                    maxlength="50"
                    show-word-limit
                  />
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="工号" prop="employee_id">
                  <el-input
                    v-model="formData.employee_id"
                    placeholder="自动生成或手动输入"
                    prefix-icon="el-icon-postcard"
                    maxlength="20"
                  >
                    <el-button slot="append" @click="generateEmployeeId">生成</el-button>
                  </el-input>
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="邮箱" prop="email" required>
                  <el-input
                    v-model="formData.email"
                    placeholder="请输入邮箱地址"
                    prefix-icon="el-icon-message"
                    maxlength="100"
                  />
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="手机号" prop="phone">
                  <el-input
                    v-model="formData.phone"
                    placeholder="请输入手机号码"
                    prefix-icon="el-icon-phone"
                    maxlength="11"
                  />
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="性别" prop="gender">
                  <el-select v-model="formData.gender" placeholder="请选择性别" style="width: 100%">
                    <el-option
                      v-for="option in genderOptions"
                      :key="option.value"
                      :label="option.label"
                      :value="option.value"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="出生日期" prop="birthday">
                  <el-date-picker
                    v-model="formData.birthday"
                    type="date"
                    placeholder="请选择出生日期"
                    style="width: 100%"
                    format="yyyy-MM-dd"
                    value-format="yyyy-MM-dd"
                  />
                </el-form-item>
              </el-col>
              
              <el-col :span="12">
                <el-form-item label="身份证号" prop="id_card">
                  <el-input
                    v-model="formData.id_card"
                    placeholder="请输入身份证号码"
                    prefix-icon="el-icon-postcard"
                    maxlength="18"
                  />
                </el-form-item>
              </el-col>
              
              <el-col :span="12">
                <el-form-item label="地址" prop="address">
                  <el-input
                    v-model="formData.address"
                    placeholder="请输入详细地址"
                    prefix-icon="el-icon-location"
                    maxlength="255"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <!-- 第二步：职位信息 -->
          <div v-show="currentStep === 1" class="form-step zoom-in">
            <div class="step-header">
              <h3>职位信息</h3>
              <p>请设置员工的组织关系和职位信息</p>
            </div>
            
            <el-row :gutter="24">
              <el-col :span="8">
                <el-form-item label="部门" prop="department_id" required>
                  <DepartmentTreeSelect
                    v-model="formData.department_id"
                    :tree-data="departmentTree"
                    placeholder="请选择部门"
                    @change="handleDepartmentChange"
                    :show-code="true"
                    :show-employee-count="true"
                    :default-expand-all="false"
                  />
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="职位" prop="position_id">
                  <el-select
                    v-model="formData.position_id"
                    placeholder="请先选择部门"
                    style="width: 100%"
                    :disabled="!formData.department_id || loadingPositions"
                    :loading="loadingPositions"
                    filterable
                    clearable
                  >
                    <el-option
                      v-for="pos in positions"
                      :key="pos.id"
                      :label="pos.name"
                      :value="pos.id"
                    >
                      <span style="float: left">{{ pos.name }}</span>
                      <span style="float: right; color: #8492a6; font-size: 13px">{{ pos.code }}</span>
                    </el-option>
                    <div v-if="loadingPositions" class="position-loading">
                      <i class="el-icon-loading"></i>
                      <span>加载职位中...</span>
                    </div>
                    <div v-else-if="formData.department_id && positions.length === 0" class="no-positions">
                      <i class="el-icon-warning"></i>
                      <span>该部门暂无可用职位</span>
                    </div>
                  </el-select>
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="职级" prop="job_level_id">
                  <el-select v-model="formData.job_level_id" placeholder="请选择职级" style="width: 100%">
                    <el-option
                      v-for="level in jobLevels"
                      :key="level.id"
                      :label="level.name"
                      :value="level.id"
                    >
                      <span style="float: left">{{ level.name }}</span>
                      <span style="float: right; color: #8492a6; font-size: 13px">L{{ level.level }}</span>
                    </el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              
              <el-col :span="12">
                <el-form-item label="直接上级" prop="manager_id">
                  <el-select
                    v-model="formData.manager_id"
                    placeholder="请选择直接上级"
                    style="width: 100%"
                    filterable
                  >
                    <el-option
                      v-for="manager in managers"
                      :key="manager.id"
                      :label="manager.name"
                      :value="manager.id"
                    >
                      <span style="float: left">{{ manager.name }}</span>
                      <span style="float: right; color: #8492a6; font-size: 13px">{{ manager.department?.name }}</span>
                    </el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              
              <el-col :span="12">
                <el-form-item label="员工状态" prop="status">
                  <el-select v-model="formData.status" placeholder="请选择状态" style="width: 100%">
                    <el-option label="在职" value="active"/>
                    <el-option label="试用期" value="probation"/>
                    <el-option label="离职" value="inactive"/>
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <!-- 第三步：入职信息 -->
          <div v-show="currentStep === 2" class="form-step zoom-in">
            <div class="step-header">
              <h3>入职信息</h3>
              <p>请填写员工的入职相关信息</p>
            </div>
            
            <el-row :gutter="24">
              <el-col :span="8">
                <el-form-item label="入职日期" prop="hire_date" required>
                  <el-date-picker
                    v-model="formData.hire_date"
                    type="date"
                    placeholder="请选择入职日期"
                    style="width: 100%"
                    format="yyyy-MM-dd"
                    value-format="yyyy-MM-dd"
                  />
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="试用期结束" prop="probation_end_date">
                  <el-date-picker
                    v-model="formData.probation_end_date"
                    type="date"
                    placeholder="请选择试用期结束日期"
                    style="width: 100%"
                    format="yyyy-MM-dd"
                    value-format="yyyy-MM-dd"
                  />
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="合同类型" prop="contract_type">
                  <el-select v-model="formData.contract_type" placeholder="请选择合同类型" style="width: 100%">
                    <el-option
                      v-for="option in contractTypeOptions"
                      :key="option.value"
                      :label="option.label"
                      :value="option.value"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="合同开始" prop="contract_start_date">
                  <el-date-picker
                    v-model="formData.contract_start_date"
                    type="date"
                    placeholder="请选择合同开始日期"
                    style="width: 100%"
                    format="yyyy-MM-dd"
                    value-format="yyyy-MM-dd"
                  />
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="合同结束" prop="contract_end_date">
                  <el-date-picker
                    v-model="formData.contract_end_date"
                    type="date"
                    placeholder="请选择合同结束日期"
                    style="width: 100%"
                    format="yyyy-MM-dd"
                    value-format="yyyy-MM-dd"
                  />
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="基本薪资" prop="base_salary">
                  <el-input-number
                    v-model="formData.base_salary"
                    :min="0"
                    :precision="2"
                    placeholder="请输入基本薪资"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <!-- 第四步：其他信息 -->
          <div v-show="currentStep === 3" class="form-step zoom-in">
            <div class="step-header">
              <h3>其他信息</h3>
              <p>请填写员工的教育背景和紧急联系人信息</p>
            </div>
            
            <el-row :gutter="24">
              <!-- 教育背景 -->
              <el-col :span="24">
                <div class="section-title">教育背景</div>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="学历" prop="education">
                  <el-select v-model="formData.education" placeholder="请选择学历" style="width: 100%">
                    <el-option
                      v-for="option in educationOptions"
                      :key="option.value"
                      :label="option.label"
                      :value="option.value"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="毕业学校" prop="school">
                  <el-input
                    v-model="formData.school"
                    placeholder="请输入毕业学校"
                    maxlength="100"
                  />
                </el-form-item>
              </el-col>
              
              <el-col :span="8">
                <el-form-item label="专业" prop="major">
                  <el-input
                    v-model="formData.major"
                    placeholder="请输入专业"
                    maxlength="100"
                  />
                </el-form-item>
              </el-col>
              
              <!-- 紧急联系人 -->
              <el-col :span="24">
                <div class="section-title">紧急联系人</div>
              </el-col>
              
              <el-col :span="12">
                <el-form-item label="联系人姓名" prop="emergency_contact">
                  <el-input
                    v-model="formData.emergency_contact"
                    placeholder="请输入紧急联系人姓名"
                    prefix-icon="el-icon-user"
                    maxlength="50"
                  />
                </el-form-item>
              </el-col>
              
              <el-col :span="12">
                <el-form-item label="联系人电话" prop="emergency_phone">
                  <el-input
                    v-model="formData.emergency_phone"
                    placeholder="请输入紧急联系人电话"
                    prefix-icon="el-icon-phone"
                    maxlength="11"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </div>
        </div>

        <!-- 步骤控制按钮 -->
        <div class="form-controls">
          <el-button
            v-if="currentStep > 0"
            @click="prevStep"
            class="control-btn"
          >
            <i class="el-icon-arrow-left"></i>
            上一步
          </el-button>
          
          <el-button
            v-if="currentStep < 3"
            type="primary"
            @click="nextStep"
            class="control-btn"
          >
            下一步
            <i class="el-icon-arrow-right"></i>
          </el-button>
          
          <el-button
            v-if="currentStep === 3"
            type="primary"
            @click="handleSubmit"
            :loading="loading"
            class="control-btn submit-btn"
          >
            <i class="el-icon-check"></i>
            提交保存
          </el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script>
import { employeeService } from '@/services/employee'
import { departmentApi } from '@/services/departmentApi'
import { positionApi } from '@/services/positionApi'
import DepartmentTreeSelect from '@/components/common/DepartmentTreeSelect.vue'

export default {
  name: 'EmployeeCreateView',
  components: {
    DepartmentTreeSelect
  },
  data() {
    return {
      loading: false,
      currentStep: 0,
      loadingPositions: false,
      
      // 表单数据
      formData: {
        name: '',
        employee_id: '',
        email: '',
        phone: '',
        gender: '',
        birthday: '',
        id_card: '',
        avatar: '',
        address: '',
        department_id: null,
        position_id: null,
        job_level_id: null,
        manager_id: null,
        status: 'active',
        hire_date: '',
        probation_end_date: '',
        contract_type: '',
        contract_start_date: '',
        contract_end_date: '',
        base_salary: null,
        education: '',
        school: '',
        major: '',
        emergency_contact: '',
        emergency_phone: ''
      },
      
      // 选项数据
      departments: [],
      departmentTree: [],
      positions: [],
      jobLevels: [],
      managers: [],
      
      // 表单验证规则
      rules: {
        name: [
          { required: true, message: '请输入员工姓名', trigger: 'blur' },
          { min: 2, max: 50, message: '姓名长度在 2 到 50 个字符', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '请输入邮箱地址', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ],
        phone: [
          { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
        ],
        id_card: [
          { pattern: /^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$/, message: '请输入正确的身份证号码', trigger: 'blur' }
        ],
        department_id: [
          { required: true, message: '请选择部门', trigger: 'change' }
        ],
        hire_date: [
          { required: true, message: '请选择入职日期', trigger: 'change' }
        ]
      }
    }
  },
  
  computed: {
    // 性别选项
    genderOptions() {
      return [
        { label: '男', value: 'male' },
        { label: '女', value: 'female' }
      ]
    },
    
    // 合同类型选项
    contractTypeOptions() {
      return [
        { label: '全职', value: 'full_time' },
        { label: '兼职', value: 'part_time' },
        { label: '实习', value: 'intern' },
        { label: '合同工', value: 'contract' }
      ]
    },
    
    // 学历选项
    educationOptions() {
      return [
        { label: '高中', value: 'high_school' },
        { label: '大专', value: 'college' },
        { label: '本科', value: 'bachelor' },
        { label: '硕士', value: 'master' },
        { label: '博士', value: 'doctor' }
      ]
    }
  },
  
  methods: {
    // 生成员工工号
    generateEmployeeId() {
      const timestamp = Date.now().toString().slice(-6)
      const random = Math.floor(Math.random() * 100).toString().padStart(2, '0')
      this.formData.employee_id = `EMP${timestamp}${random}`
    },
    
    // 头像上传前验证
    beforeAvatarUpload(file) {
      const isJPG = file.type === 'image/jpeg' || file.type === 'image/png'
      const isLt2M = file.size / 1024 / 1024 < 2

      if (!isJPG) {
        this.$message.error('上传头像图片只能是 JPG/PNG 格式!')
        return false
      }
      if (!isLt2M) {
        this.$message.error('上传头像图片大小不能超过 2MB!')
        return false
      }
      return true
    },
    
    // 处理头像上传
    handleAvatarUpload(param) {
      // 这里可以实现实际的上传逻辑
      const reader = new FileReader()
      reader.onload = (e) => {
        this.formData.avatar = e.target.result
      }
      reader.readAsDataURL(param.file)
    },
    
    // 处理部门变化
    async handleDepartmentChange(departmentId, departmentData) {
      // 清空当前选择的职位和直接上级
      this.formData.position_id = null
      this.formData.manager_id = null
      
      if (departmentId) {
        // 显示加载状态
        this.loadingPositions = true
        await this.fetchPositionsByDepartment(departmentId)
        this.loadingPositions = false
        
        // 如果只有一个职位，自动选择
        if (this.positions.length === 1) {
          this.formData.position_id = this.positions[0].id
          this.$message.success(`已自动选择职位：${this.positions[0].name}`)
        }

        // 自动选择部门领导作为直接上级
        await this.setDepartmentManagerAsSupervisor(departmentId, departmentData)
      } else {
        this.positions = []
      }
    },
    
    // 获取部门职位
    async fetchPositionsByDepartment(departmentId) {
      try {
        const response = await positionApi.getPositionsByDepartment(departmentId)
        this.positions = response.data || []
        
        if (this.positions.length === 0) {
          this.$message.warning('该部门暂无可用职位，请联系管理员添加职位')
        }
      } catch (error) {
        console.error('获取职位列表失败:', error)
        this.positions = []
        this.$message.error('获取职位列表失败，请重试')
      }
    },

    // 设置部门领导为直接上级
    async setDepartmentManagerAsSupervisor(departmentId, departmentData) {
      try {
        // 如果departmentData中包含manager信息，直接使用
        if (departmentData && departmentData.manager_id && departmentData.manager) {
          this.formData.manager_id = departmentData.manager_id
          this.$message.success(`已自动选择部门领导 ${departmentData.manager.name} 为直接上级`)
          return
        }

        // 否则通过API获取部门详细信息
        const response = await departmentApi.getDepartment(departmentId)
        const department = response.data || response
        
        if (department.manager_id && department.manager) {
          this.formData.manager_id = department.manager_id
          this.$message.success(`已自动选择部门领导 ${department.manager.name} 为直接上级`)
        } else {
          this.$message.info('该部门暂未设置直属领导，请手动选择直接上级')
        }
      } catch (error) {
        console.error('获取部门领导信息失败:', error)
        this.$message.warning('无法获取部门领导信息，请手动选择直接上级')
      }
    },
    
    // 加载表单数据
    async loadFormData() {
      try {
        const [departmentsRes, departmentTreeRes, jobLevelsRes, managersRes] = await Promise.all([
          employeeService.getDepartments(),
          departmentApi.getDepartmentTree(),
          employeeService.getJobLevels(),
          employeeService.getManagers()
        ])
        
        this.departments = departmentsRes.data || []
        
        // 处理部门树数据
        const treeData = departmentTreeRes.data?.data || departmentTreeRes.data || []
        this.departmentTree = this.processDepartmentTree(treeData)
        
        this.jobLevels = jobLevelsRes.data || []
        this.managers = managersRes.data || []
      } catch (error) {
        console.error('加载表单数据失败:', error)
        this.$message.error('加载表单数据失败')
      }
    },
    
    // 处理部门树数据
    processDepartmentTree(treeData) {
      const processNode = (node) => {
        const children = []
        if (node.children && node.children.length > 0) {
          for (const child of node.children) {
            children.push(processNode(child))
          }
        }
        
        // 使用后端提供的真实数据
        const directEmployeeCount = node.employee_count || 0
        const hierarchicalEmployeeCount = this.calculateHierarchicalEmployeeCount(node)
        
        return {
          id: node.id,
          name: node.name,
          code: node.code,
          type: node.type || 'department',
          level: node.level || 1,
          description: node.description,
          manager_id: node.manager_id,
          employeeCount: directEmployeeCount,
          hierarchicalEmployeeCount: hierarchicalEmployeeCount,
          children: children
        }
      }
      
      const result = []
      for (const node of treeData) {
        result.push(processNode(node))
      }
      
      return result
    },
    
    // 计算层级员工数量（包括所有子部门的员工）
    calculateHierarchicalEmployeeCount(node) {
      let totalCount = node.employee_count || 0
      
      if (node.children && node.children.length > 0) {
        for (const child of node.children) {
          totalCount += this.calculateHierarchicalEmployeeCount(child)
        }
      }
      
      return totalCount
    },
    
    // 下一步
    async nextStep() {
      // 验证当前步骤的字段
      const fieldsToValidate = this.getFieldsForStep(this.currentStep)
      
      try {
        if (fieldsToValidate.length > 0) {
          await this.$refs.employeeForm.validateField(fieldsToValidate)
        }
        if (this.currentStep < 3) {
          this.currentStep++
        }
      } catch (error) {
        console.log('验证失败:', error)
      }
    },
    
    // 上一步
    prevStep() {
      if (this.currentStep > 0) {
        this.currentStep--
      }
    },
    
    // 获取当前步骤需要验证的字段
    getFieldsForStep(step) {
      const stepFields = {
        0: ['name', 'email'], // 基本信息
        1: ['department_id'], // 职位信息
        2: ['hire_date'], // 入职信息
        3: [] // 其他信息
      }
      return stepFields[step] || []
    },
    
    // 提交表单
    async handleSubmit() {
      try {
        await this.$refs.employeeForm.validate()
        
        this.loading = true
        
        // 处理表单数据
        const submitData = { ...this.formData }
        
        // 移除空值
        Object.keys(submitData).forEach(key => {
          if (submitData[key] === '' || submitData[key] === null) {
            delete submitData[key]
          }
        })
        
        await employeeService.createEmployee(submitData)
        
        this.$message.success('员工创建成功')
        this.$router.push('/employees')
      } catch (error) {
        console.error('创建员工失败:', error)
        const message = error.response?.data?.message || error.message || '创建失败'
        this.$message.error(message)
      } finally {
        this.loading = false
      }
    },
    
    // 取消操作
    handleCancel() {
      this.$confirm('确定要离开吗？未保存的数据将丢失', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$router.push('/employees')
      }).catch(() => {})
    },

    // 处理来自组织架构的查询参数预填充
    async handleQueryParameters() {
      const query = this.$route.query
      
      // 如果有部门信息，预填充部门相关字段
      if (query.department_id) {
        console.log('从组织架构传入部门信息:', query)
        
        // 设置部门ID
        this.formData.department_id = parseInt(query.department_id)
        
        // 设置直接上级
        if (query.direct_manager_id) {
          this.formData.manager_id = parseInt(query.direct_manager_id)
        }
        
        // 处理部门职位信息
        if (query.department_positions) {
          try {
            const positions = JSON.parse(query.department_positions)
            // 更新可选职位列表，只显示该部门的职位
            this.positionOptions = positions
            
            // 如果只有一个职位，自动选择
            if (positions.length === 1) {
              this.formData.position_id = positions[0].id
            }
          } catch (error) {
            console.error('解析部门职位数据失败:', error)
          }
        }
        
        // 触发部门变更逻辑，确保相关数据加载正确
        this.handleDepartmentChange(this.formData.department_id)
        
        // 显示预填充提示
        this.$message.success(`已自动填充部门信息：${query.department_name || ''}`)
      }
    }
  },
  
  async mounted() {
    await this.loadFormData()
    this.generateEmployeeId()
    // 处理来自组织架构的预填充参数
    this.handleQueryParameters()
  }
}
</script>

<style scoped>
/* 整体容器 */
.employee-create {
  padding: 0;
  min-height: 100vh;
  background: var(--background-color);
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 32px;
  padding: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.2);
  position: relative;
  overflow: hidden;
}

.page-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="%23ffffff" fill-opacity="0.1"/><circle cx="75" cy="75" r="1" fill="%23ffffff" fill-opacity="0.1"/></pattern></defs><rect width="100%" height="100%" fill="url(%23grain)"/></svg>');
  opacity: 0.3;
  pointer-events: none;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 16px;
  position: relative;
  z-index: 1;
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
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.title-content p {
  margin: 8px 0 0;
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}

.page-actions {
  display: flex;
  gap: 12px;
  position: relative;
  z-index: 1;
}

.action-btn {
  padding: 12px 24px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.1);
  color: white;
  backdrop-filter: blur(10px);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
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

/* 主要内容 */
.create-content {
  background: white;
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.05);
  margin-bottom: 24px;
}

/* 步骤进度 */
.form-steps {
  margin-bottom: 40px;
  padding: 24px;
  background: #f8f9fa;
  border-radius: 8px;
}

.form-steps :deep(.el-steps) {
  margin: 0;
}

.form-steps :deep(.el-step__title) {
  font-size: 14px;
  font-weight: 500;
}

/* 表单内容 */
.form-content {
  min-height: 400px;
}

.form-step {
  animation-duration: 0.3s;
}

.step-header {
  margin-bottom: 32px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.step-header h3 {
  margin: 0 0 8px;
  font-size: 20px;
  font-weight: 600;
  color: #262626;
}

.step-header p {
  margin: 0;
  color: #8c8c8c;
  font-size: 14px;
}

/* 头像上传 */
.avatar-section {
  display: flex;
  justify-content: center;
  margin-bottom: 32px;
}

.avatar-uploader {
  display: block;
}

.avatar-uploader :deep(.el-upload) {
  border: 2px dashed #d9d9d9;
  border-radius: 50%;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  width: 120px;
  height: 120px;
  transition: all 0.3s;
}

.avatar-uploader :deep(.el-upload:hover) {
  border-color: #409eff;
  transform: scale(1.05);
}

.avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 120px;
  height: 120px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #8c8c8c;
  background: #fafafa;
  border-radius: 50%;
}

.avatar-uploader-icon {
  font-size: 28px;
  margin-bottom: 8px;
}

.avatar-text {
  font-size: 12px;
}

/* 分组标题 */
.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;
}

/* 表单控制按钮 */
.form-controls {
  margin-top: 40px;
  padding-top: 24px;
  border-top: 1px solid #f0f0f0;
  text-align: right;
}

.control-btn {
  margin-left: 12px;
  padding: 10px 24px;
  border-radius: 6px;
  font-weight: 500;
  transition: all 0.3s;
}

.control-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.submit-btn {
  background: linear-gradient(135deg, #52c41a, #389e0d);
  border: none;
  color: white;
}

.submit-btn:hover {
  background: linear-gradient(135deg, #389e0d, #52c41a);
}

/* 表单项优化 */
.employee-form :deep(.el-form-item__label) {
  font-weight: 500;
  color: #262626;
}

.employee-form :deep(.el-input__inner) {
  border-radius: 6px;
  transition: all 0.3s;
}

.employee-form :deep(.el-input__inner:focus) {
  box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
}

.employee-form :deep(.el-select) {
  width: 100%;
}

.employee-form :deep(.el-date-editor) {
  width: 100%;
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
  
  .create-content {
    padding: 16px;
  }
  
  .form-steps {
    padding: 16px;
  }
  
  .form-controls {
    text-align: center;
  }
  
  .control-btn {
    margin: 8px;
  }
}

/* 动画效果 */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes zoomIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.fade-in {
  animation: fadeIn 0.3s ease-in-out;
}

.slide-up {
  animation: slideUp 0.3s ease-out;
}

.zoom-in {
  animation: zoomIn 0.3s ease-out;
}

/* 职位选择相关样式 */
.position-loading,
.no-positions {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  color: #909399;
  font-size: 13px;
}

.position-loading i {
  animation: rotation 1s linear infinite;
}

.no-positions {
  color: #f56c6c;
}

@keyframes rotation {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>