<template>
  <el-dialog
    :visible.sync="dialogVisible"
    :title="isEdit ? '编辑员工分配' : '新增员工分配'"
    width="800px"
    @close="handleClose"
    :close-on-click-modal="false"
  >
    <el-form
      ref="assignmentForm"
      :model="formData"
      :rules="formRules"
      label-width="120px"
      size="small"
    >
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="员工" prop="employee_id">
            <el-select
              v-model="formData.employee_id"
              placeholder="请选择员工"
              style="width: 100%"
              filterable
              remote
              :remote-method="searchEmployees"
              :loading="employeeLoading"
              @change="handleEmployeeChange"
            >
              <el-option
                v-for="employee in availableEmployees"
                :key="employee.id"
                :label="employee.name"
                :value="employee.id"
              >
                <div style="display: flex; justify-content: space-between;">
                  <span>{{ employee.name }}</span>
                  <span style="color: #8492a6; font-size: 12px;">{{ employee.employee_id }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="组织单元" prop="organization_unit_id">
            <el-cascader
              v-model="formData.organization_unit_id"
              :options="organizationTree"
              :props="cascaderProps"
              placeholder="请选择组织单元"
              style="width: 100%"
              clearable
              :show-all-levels="false"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="分配类型" prop="assignment_type">
            <el-select
              v-model="formData.assignment_type"
              placeholder="请选择分配类型"
              style="width: 100%"
            >
              <el-option
                v-for="type in assignmentTypes"
                :key="type.value"
                :label="type.label"
                :value="type.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="管理类型" prop="management_type">
            <el-select
              v-model="formData.management_type"
              placeholder="请选择管理类型"
              style="width: 100%"
            >
              <el-option
                v-for="type in managementTypes"
                :key="type.value"
                :label="type.label"
                :value="type.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="工作占比" prop="work_percentage">
            <el-input-number
              v-model="formData.work_percentage"
              :min="1"
              :max="100"
              :precision="2"
              style="width: 100%"
              controls-position="right"
            />
            <span style="margin-left: 8px; color: #909399;">%</span>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="职位" prop="position_id">
            <el-select
              v-model="formData.position_id"
              placeholder="请选择职位"
              style="width: 100%"
              filterable
            >
              <el-option
                v-for="position in availablePositions"
                :key="position.id"
                :label="position.name"
                :value="position.id"
              />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="直接上级" prop="direct_manager_id">
            <el-select
              v-model="formData.direct_manager_id"
              placeholder="请选择直接上级"
              style="width: 100%"
              filterable
              remote
              :remote-method="searchManagers"
              :loading="managerLoading"
              clearable
            >
              <el-option
                v-for="manager in availableManagers"
                :key="manager.id"
                :label="manager.name"
                :value="manager.id"
              >
                <div style="display: flex; justify-content: space-between;">
                  <span>{{ manager.name }}</span>
                  <span style="color: #8492a6; font-size: 12px;">{{ manager.position_name }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="功能上级" prop="functional_manager_id">
            <el-select
              v-model="formData.functional_manager_id"
              placeholder="请选择功能上级"
              style="width: 100%"
              filterable
              remote
              :remote-method="searchManagers"
              :loading="managerLoading"
              clearable
            >
              <el-option
                v-for="manager in availableManagers"
                :key="manager.id"
                :label="manager.name"
                :value="manager.id"
              >
                <div style="display: flex; justify-content: space-between;">
                  <span>{{ manager.name }}</span>
                  <span style="color: #8492a6; font-size: 12px;">{{ manager.position_name }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="生效日期" prop="effective_date">
            <el-date-picker
              v-model="formData.effective_date"
              type="date"
              placeholder="请选择生效日期"
              style="width: 100%"
              value-format="yyyy-MM-dd"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="失效日期" prop="expiration_date">
            <el-date-picker
              v-model="formData.expiration_date"
              type="date"
              placeholder="请选择失效日期"
              style="width: 100%"
              value-format="yyyy-MM-dd"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="是否主要分配">
        <el-switch
          v-model="formData.is_primary"
          active-text="是"
          inactive-text="否"
        />
        <span style="margin-left: 12px; color: #909399; font-size: 12px;">
          主要分配将成为员工的默认组织关系
        </span>
      </el-form-item>

      <el-form-item label="分配原因" prop="assignment_reason">
        <el-input
          v-model="formData.assignment_reason"
          type="textarea"
          placeholder="请输入分配原因"
          :rows="3"
          maxlength="500"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="备注" prop="notes">
        <el-input
          v-model="formData.notes"
          type="textarea"
          placeholder="请输入备注信息"
          :rows="2"
          maxlength="200"
          show-word-limit
        />
      </el-form-item>
    </el-form>

    <div slot="footer" class="dialog-footer">
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
        {{ isEdit ? '更新' : '创建' }}
      </el-button>
    </div>
  </el-dialog>
</template>

<script>
import { organizationApi } from '@/services/organizationApi'
import { employeeApi } from '@/services/employeeApi'
import { positionApi } from '@/services/positionApi'

export default {
  name: 'EmployeeAssignmentDialog',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    organizationUnit: {
      type: Object,
      default: null
    },
    selectedDepartmentInfo: {
      type: Object,
      default: null
    },
    assignment: {
      type: Object,
      default: null
    }
  },
  data() {
    return {
      submitLoading: false,
      employeeLoading: false,
      managerLoading: false,
      availableEmployees: [],
      availableManagers: [],
      availablePositions: [],
      organizationTree: [],
      
      formData: {
        employee_id: null,
        organization_unit_id: null,
        assignment_type: 'primary',
        management_type: 'line',
        work_percentage: 100,
        position_id: null,
        direct_manager_id: null,
        functional_manager_id: null,
        effective_date: null,
        expiration_date: null,
        is_primary: false,
        assignment_reason: '',
        notes: ''
      },
      
      formRules: {
        employee_id: [
          { required: true, message: '请选择员工', trigger: 'change' }
        ],
        organization_unit_id: [
          { required: true, message: '请选择组织单元', trigger: 'change' }
        ],
        assignment_type: [
          { required: true, message: '请选择分配类型', trigger: 'change' }
        ],
        management_type: [
          { required: true, message: '请选择管理类型', trigger: 'change' }
        ],
        work_percentage: [
          { required: true, message: '请输入工作占比', trigger: 'blur' },
          { type: 'number', min: 1, max: 100, message: '工作占比必须在1-100之间', trigger: 'blur' }
        ],
        effective_date: [
          { required: true, message: '请选择生效日期', trigger: 'change' }
        ]
      },
      
      cascaderProps: {
        value: 'id',
        label: 'name',
        children: 'children',
        checkStrictly: true
      },
      
      assignmentTypes: [
        { value: 'primary', label: '主要分配' },
        { value: 'additional', label: '额外分配' },
        { value: 'temporary', label: '临时分配' },
        { value: 'project', label: '项目分配' }
      ],
      
      managementTypes: [
        { value: 'line', label: '直线管理' },
        { value: 'matrix', label: '矩阵管理' },
        { value: 'functional', label: '功能管理' }
      ]
    }
  },
  computed: {
    dialogVisible: {
      get() {
        return this.visible
      },
      set(value) {
        this.$emit('update:visible', value)
      }
    },
    
    isEdit() {
      return !!this.assignment
    }
  },
  watch: {
    visible(val) {
      if (val) {
        this.initDialog()
      }
    }
  },
  methods: {
    async initDialog() {
      await Promise.all([
        this.loadOrganizationTree(),
        this.loadPositions(),
        this.loadInitialEmployees(),
        this.loadInitialManagers()
      ])
      
      if (this.isEdit && this.assignment) {
        this.formData = { ...this.assignment }
      } else {
        this.resetForm()
        // 优先使用 selectedDepartmentInfo，然后使用 organizationUnit
        const departmentInfo = this.selectedDepartmentInfo || this.organizationUnit
        if (departmentInfo) {
          this.formData.organization_unit_id = departmentInfo.id
          
          // 从部门信息预填充相关字段
          if (departmentInfo.manager) {
            this.formData.direct_manager_id = departmentInfo.manager.id
          }
          
          // 设置默认生效日期为今天
          if (!this.formData.effective_date) {
            const today = new Date()
            this.formData.effective_date = today.toISOString().split('T')[0]
          }
          
          // 设置默认分配类型为主要分配
          this.formData.assignment_type = 'primary'
          this.formData.is_primary = true
          
          // 设置默认工作占比为100%
          this.formData.work_percentage = 100
        }
      }
    },
    
    async loadOrganizationTree() {
      try {
        const response = await organizationApi.getOrganizationTree()
        this.organizationTree = response.data
      } catch (error) {
        console.error('Failed to load organization tree:', error)
      }
    },
    
    async loadPositions() {
      try {
        const response = await positionApi.getPositionList()
        this.availablePositions = response.data
      } catch (error) {
        console.error('Failed to load positions:', error)
      }
    },
    
    async loadInitialEmployees() {
      try {
        const response = await employeeApi.getEmployeeList({ page: 1, size: 50 })
        this.availableEmployees = response.data.items || []
      } catch (error) {
        console.error('Failed to load employees:', error)
      }
    },
    
    async loadInitialManagers() {
      try {
        const response = await employeeApi.getEmployeeList({ 
          page: 1, 
          size: 50,
          is_manager: true 
        })
        this.availableManagers = response.data.items || []
      } catch (error) {
        console.error('Failed to load managers:', error)
      }
    },
    
    async searchEmployees(query) {
      if (!query) {
        this.loadInitialEmployees()
        return
      }
      
      this.employeeLoading = true
      try {
        const response = await employeeApi.searchEmployees({ keyword: query, limit: 20 })
        this.availableEmployees = response.data || []
      } catch (error) {
        console.error('Failed to search employees:', error)
      } finally {
        this.employeeLoading = false
      }
    },
    
    async searchManagers(query) {
      if (!query) {
        this.loadInitialManagers()
        return
      }
      
      this.managerLoading = true
      try {
        const response = await employeeApi.searchEmployees({ 
          keyword: query, 
          limit: 20,
          is_manager: true 
        })
        this.availableManagers = response.data || []
      } catch (error) {
        console.error('Failed to search managers:', error)
      } finally {
        this.managerLoading = false
      }
    },
    
    handleEmployeeChange(employeeId) {
      const employee = this.availableEmployees.find(e => e.id === employeeId)
      if (employee && employee.position_id) {
        this.formData.position_id = employee.position_id
      }
    },
    
    resetForm() {
      this.formData = {
        employee_id: null,
        organization_unit_id: null,
        assignment_type: 'primary',
        management_type: 'line',
        work_percentage: 100,
        position_id: null,
        direct_manager_id: null,
        functional_manager_id: null,
        effective_date: null,
        expiration_date: null,
        is_primary: false,
        assignment_reason: '',
        notes: ''
      }
      
      this.$nextTick(() => {
        this.$refs.assignmentForm?.clearValidate()
      })
    },
    
    async handleSubmit() {
      try {
        await this.$refs.assignmentForm.validate()
        
        this.submitLoading = true
        
        if (this.isEdit) {
          await organizationApi.updateAssignment(this.assignment.id, this.formData)
        } else {
          await organizationApi.assignEmployee(this.formData)
        }
        
        this.$emit('success')
        this.handleClose()
        this.$message.success(this.isEdit ? '更新成功' : '分配成功')
      } catch (error) {
        if (error.response) {
          this.$message.error(error.response.data.message || '操作失败')
        } else {
          console.error('Form validation failed:', error)
        }
      } finally {
        this.submitLoading = false
      }
    },
    
    handleClose() {
      this.dialogVisible = false
      this.resetForm()
    }
  }
}
</script>

<style scoped>
.dialog-footer {
  text-align: right;
}

.el-form-item {
  margin-bottom: 16px;
}
</style>