<template>
  <el-dialog
    :visible.sync="dialogVisible"
    :title="isEdit ? '编辑组织单元' : '新建组织单元'"
    width="600px"
    @close="handleClose"
    :close-on-click-modal="false"
  >
    <el-form
      ref="unitForm"
      :model="formData"
      :rules="formRules"
      label-width="120px"
      size="small"
    >
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="组织单元名称" prop="name">
            <el-input
              v-model="formData.name"
              placeholder="请输入组织单元名称"
              maxlength="200"
              show-word-limit
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="组织编码" prop="code">
            <el-input
              v-model="formData.code"
              placeholder="请输入组织编码"
              maxlength="50"
              show-word-limit
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="组织类型" prop="type">
            <el-select
              v-model="formData.type"
              placeholder="请选择组织类型"
              style="width: 100%"
            >
              <el-option
                v-for="type in unitTypes"
                :key="type.value"
                :label="type.label"
                :value="type.value"
              >
                <div style="display: flex; align-items: center; gap: 8px;">
                  <i :class="type.icon"></i>
                  <span>{{ type.label }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="简称" prop="short_name">
            <el-input
              v-model="formData.short_name"
              placeholder="请输入简称"
              maxlength="50"
              show-word-limit
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="上级组织" prop="parent_id">
        <el-cascader
          v-model="formData.parent_id"
          :options="organizationTree"
          :props="cascaderProps"
          placeholder="请选择上级组织"
          style="width: 100%"
          clearable
          :show-all-levels="false"
          :disabled="isEdit && unit && unit.id"
        />
      </el-form-item>

      <el-form-item label="描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          placeholder="请输入描述"
          :rows="3"
          maxlength="500"
          show-word-limit
        />
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="负责人" prop="manager_id">
            <el-select
              v-model="formData.manager_id"
              placeholder="请选择负责人"
              style="width: 100%"
              filterable
              remote
              :remote-method="searchEmployees"
              :loading="employeeLoading"
              clearable
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
          <el-form-item label="功能负责人" prop="functional_manager_id">
            <el-select
              v-model="formData.functional_manager_id"
              placeholder="请选择功能负责人"
              style="width: 100%"
              filterable
              remote
              :remote-method="searchEmployees"
              :loading="employeeLoading"
              clearable
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
      </el-row>

      <el-row :gutter="20">
        <el-col :span="8">
          <el-form-item label="成本中心" prop="cost_center">
            <el-input
              v-model="formData.cost_center"
              placeholder="请输入成本中心编码"
              maxlength="50"
            />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="国家代码" prop="country_code">
            <el-input
              v-model="formData.country_code"
              placeholder="如: CN"
              maxlength="10"
            />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="货币代码" prop="currency_code">
            <el-input
              v-model="formData.currency_code"
              placeholder="如: CNY"
              maxlength="10"
            />
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

      <el-row :gutter="20">
        <el-col :span="24">
          <el-form-item label="联系信息">
            <el-collapse>
              <el-collapse-item title="详细信息" name="contact">
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="电话" prop="phone">
                      <el-input
                        v-model="formData.phone"
                        placeholder="请输入电话号码"
                        maxlength="50"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="邮箱" prop="email">
                      <el-input
                        v-model="formData.email"
                        placeholder="请输入邮箱"
                        maxlength="100"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-form-item label="地址" prop="address">
                  <el-input
                    v-model="formData.address"
                    placeholder="请输入地址"
                    maxlength="500"
                  />
                </el-form-item>
                <el-form-item label="网站" prop="website">
                  <el-input
                    v-model="formData.website"
                    placeholder="请输入网站地址"
                    maxlength="200"
                  />
                </el-form-item>
              </el-collapse-item>
            </el-collapse>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="8">
          <el-form-item label="状态设置">
            <el-switch
              v-model="formData.is_active"
              active-text="激活"
              inactive-text="停用"
            />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="总部标识">
            <el-switch
              v-model="formData.is_headquarters"
              active-text="是"
              inactive-text="否"
            />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="允许子单元">
            <el-switch
              v-model="formData.allow_subunits"
              active-text="允许"
              inactive-text="不允许"
            />
          </el-form-item>
        </el-col>
      </el-row>
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

export default {
  name: 'OrganizationUnitDialog',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    unit: {
      type: Object,
      default: null
    },
    parentUnit: {
      type: Object,
      default: null
    },
    parentDepartmentInfo: {
      type: Object,
      default: null
    },
    unitTypes: {
      type: Array,
      default: () => []
    },
    isEdit: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      submitLoading: false,
      employeeLoading: false,
      availableEmployees: [],
      organizationTree: [],
      
      formData: {
        name: '',
        code: '',
        type: '',
        parent_id: null,
        description: '',
        short_name: '',
        manager_id: null,
        functional_manager_id: null,
        cost_center: '',
        country_code: '',
        currency_code: '',
        phone: '',
        email: '',
        address: '',
        website: '',
        effective_date: null,
        expiration_date: null,
        is_active: true,
        is_headquarters: false,
        allow_subunits: true
      },
      
      formRules: {
        name: [
          { required: true, message: '请输入组织单元名称', trigger: 'blur' },
          { min: 2, max: 200, message: '长度在 2 到 200 个字符', trigger: 'blur' }
        ],
        code: [
          { required: true, message: '请输入组织编码', trigger: 'blur' },
          { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
        ],
        type: [
          { required: true, message: '请选择组织类型', trigger: 'change' }
        ],
        email: [
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ]
      },
      
      cascaderProps: {
        value: 'id',
        label: 'name',
        children: 'children',
        checkStrictly: true
      }
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
      await this.loadOrganizationTree()
      await this.loadInitialEmployees()
      
      if (this.isEdit && this.unit) {
        this.formData = { ...this.unit }
      } else {
        this.resetForm()
        // 优先使用 parentDepartmentInfo，然后使用 parentUnit
        const parentInfo = this.parentDepartmentInfo || this.parentUnit
        if (parentInfo) {
          this.formData.parent_id = parentInfo.id
          
          // 从父部门继承一些默认值
          if (parentInfo.cost_center) {
            this.formData.cost_center = parentInfo.cost_center
          }
          if (parentInfo.country_code) {
            this.formData.country_code = parentInfo.country_code
          }
          if (parentInfo.currency_code) {
            this.formData.currency_code = parentInfo.currency_code
          }
          
          // 设置默认层级（父级层级+1）
          if (parentInfo.level) {
            // 注意：这里不直接设置level字段，因为level通常由后端计算
            // 但可以用于生成默认的组织编码
            if (!this.formData.code && parentInfo.code) {
              // 生成建议的子部门编码，例如：PARENT001 -> PARENT001-01
              const parentCode = parentInfo.code
              const nextSubCode = String(Math.floor(Math.random() * 99) + 1).padStart(2, '0')
              this.formData.code = `${parentCode}-${nextSubCode}`
            }
          }
          
          // 设置默认生效日期为今天
          if (!this.formData.effective_date) {
            const today = new Date()
            this.formData.effective_date = today.toISOString().split('T')[0]
          }
          
          // 设置默认的组织类型（通常子部门为 department）
          if (!this.formData.type) {
            this.formData.type = 'department'
          }
          
          // 如果父部门有管理者，可以作为功能负责人的候选
          if (parentInfo.manager_id && !this.formData.functional_manager_id) {
            this.formData.functional_manager_id = parentInfo.manager_id
          }
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
    
    async loadInitialEmployees() {
      try {
        const response = await employeeApi.getEmployeeList({ page: 1, size: 50 })
        this.availableEmployees = response.data.items || []
      } catch (error) {
        console.error('Failed to load employees:', error)
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
    
    resetForm() {
      this.formData = {
        name: '',
        code: '',
        type: '',
        parent_id: null,
        description: '',
        short_name: '',
        manager_id: null,
        functional_manager_id: null,
        cost_center: '',
        country_code: '',
        currency_code: '',
        phone: '',
        email: '',
        address: '',
        website: '',
        effective_date: null,
        expiration_date: null,
        is_active: true,
        is_headquarters: false,
        allow_subunits: true
      }
      
      this.$nextTick(() => {
        this.$refs.unitForm?.clearValidate()
      })
    },
    
    async handleSubmit() {
      try {
        await this.$refs.unitForm.validate()
        
        this.submitLoading = true
        
        if (this.isEdit) {
          await organizationApi.updateOrganizationUnit(this.unit.id, this.formData)
        } else {
          await organizationApi.createOrganizationUnit(this.formData)
        }
        
        this.$emit('success')
        this.handleClose()
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

.el-collapse {
  border: none;
}

.el-collapse :deep(.el-collapse-item__header) {
  background: #f5f7fa;
  border-bottom: 1px solid #e4e7ed;
  padding-left: 16px;
}

.el-collapse :deep(.el-collapse-item__content) {
  padding: 16px;
}
</style>