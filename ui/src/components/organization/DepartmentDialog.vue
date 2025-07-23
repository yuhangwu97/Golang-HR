<template>
  <el-dialog
    :visible.sync="dialogVisible"
    :title="dialogTitle"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="部门名称" prop="name">
        <el-input
          v-model="form.name"
          placeholder="请输入部门名称"
          maxlength="100"
        />
      </el-form-item>
      
      <el-form-item label="部门编码" prop="code">
        <el-input
          v-model="form.code"
          placeholder="请输入部门编码"
          maxlength="20"
        />
      </el-form-item>
      
      <el-form-item label="上级部门" prop="parent_id">
        <el-cascader
          v-model="form.parent_id"
          :options="departmentTree"
          :props="cascaderProps"
          placeholder="选择上级部门（可选）"
          clearable
          style="width: 100%"
        />
        <div v-if="parentDepartment" class="form-tip">
          将在"{{ parentDepartment.name }}"下创建子部门
        </div>
      </el-form-item>
      
      <el-form-item label="部门负责人" prop="manager_id">
        <el-select
          v-model="form.manager_id"
          placeholder="选择部门负责人（可选）"
          style="width: 100%"
          clearable
          filterable
        >
          <el-option
            v-for="employee in employees"
            :key="employee.id"
            :label="`${employee.name} (${employee.employee_id})`"
            :value="employee.id"
          />
        </el-select>
      </el-form-item>
      
      <el-form-item label="部门描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="请输入部门描述"
          maxlength="500"
        />
      </el-form-item>
      
      <el-form-item label="排序权重" prop="sort">
        <el-input-number
          v-model="form.sort"
          :min="0"
          :max="999"
          placeholder="排序权重，数字越大越靠前"
          style="width: 100%"
        />
        <div class="form-tip">数字越大，部门在列表中越靠前显示</div>
      </el-form-item>
      
      <el-form-item label="状态" prop="status">
        <el-radio-group v-model="form.status">
          <el-radio label="active">活跃</el-radio>
          <el-radio label="inactive">非活跃</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>

    <span slot="footer" class="dialog-footer">
      <el-button @click="handleClose">取消</el-button>
      <el-button
        type="primary"
        :loading="loading"
        @click="handleSubmit"
      >
        {{ isEdit ? '更新' : '创建' }}
      </el-button>
    </span>
  </el-dialog>
</template>

<script>
import { departmentApi } from '@/services/departmentApi'
import { employeeService } from '@/services/employee'

export default {
  name: 'DepartmentDialog',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    department: {
      type: Object,
      default: null
    },
    parentDepartment: {
      type: Object,
      default: null
    }
  },
  data() {
    return {
      loading: false,
      departmentTree: [],
      employees: [],
      form: {
        name: '',
        code: '',
        parent_id: undefined,
        manager_id: undefined,
        description: '',
        sort: 0,
        status: 'active'
      },
      cascaderProps: {
        value: 'id',
        label: 'name',
        children: 'children',
        emitPath: false,
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
    },
    isEdit() {
      return !!this.department
    },
    dialogTitle() {
      if (this.parentDepartment) {
        return `在"${this.parentDepartment.name}"下添加子部门`
      }
      return this.isEdit ? '编辑部门' : '新增部门'
    },
    rules() {
      return {
        name: [
          { required: true, message: '请输入部门名称', trigger: 'blur' },
          { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
        ],
        code: [
          { required: true, message: '请输入部门编码', trigger: 'blur' },
          { min: 1, max: 20, message: '长度在 1 到 20 个字符', trigger: 'blur' },
          { pattern: /^[A-Za-z0-9_-]+$/, message: '编码只能包含字母、数字、下划线和连字符', trigger: 'blur' }
        ],
        description: [
          { max: 500, message: '长度不能超过 500 个字符', trigger: 'blur' }
        ],
        sort: [
          { type: 'number', min: 0, max: 999, message: '排序权重必须在 0 到 999 之间', trigger: 'blur' }
        ],
        status: [
          { required: true, message: '请选择状态', trigger: 'change' }
        ]
      }
    }
  },
  watch: {
    visible(val) {
      if (val) {
        this.initForm()
        this.fetchDepartmentTree()
        this.fetchEmployees()
      }
    }
  },
  mounted() {
    this.fetchDepartmentTree()
    this.fetchEmployees()
  },
  methods: {
    async fetchDepartmentTree() {
      try {
        const response = await departmentApi.getDepartmentTree()
        const responseData = response.data || response
        this.departmentTree = responseData.data || responseData || []
      } catch (error) {
        console.error('获取部门树失败:', error)
        this.$message.error('获取部门树失败')
      }
    },
    async fetchEmployees() {
      try {
        const response = await employeeService.getEmployees({ page: 1, pageSize: 100, status: 'active' })
        const responseData = response.data || response
        this.employees = responseData.data || responseData || []
      } catch (error) {
        console.error('获取员工列表失败:', error)
        this.$message.error('获取员工列表失败')
      }
    },
    initForm() {
      if (this.department) {
        Object.assign(this.form, {
          name: this.department.name,
          code: this.department.code,
          parent_id: this.department.parent_id,
          manager_id: this.department.manager_id,
          description: this.department.description,
          sort: this.department.sort || 0,
          status: this.department.status
        })
      } else {
        Object.assign(this.form, {
          name: '',
          code: '',
          parent_id: this.parentDepartment?.id,
          manager_id: undefined,
          description: '',
          sort: 0,
          status: 'active'
        })
      }
    },
    async handleSubmit() {
      if (!this.$refs.formRef) return
      
      try {
        await this.$refs.formRef.validate()
        this.loading = true
        
        // 调用实际API
        if (this.isEdit && this.department) {
          await departmentApi.updateDepartment(this.department.id, this.form)
          this.$message.success('更新成功')
        } else {
          await departmentApi.createDepartment(this.form)
          this.$message.success('创建成功')
        }
        
        this.$emit('success')
        this.handleClose()
      } catch (error) {
        console.error('操作失败:', error)
        if (error.response?.data?.message) {
          this.$message.error(error.response.data.message)
        } else {
          this.$message.error(this.isEdit ? '更新失败' : '创建失败')
        }
      } finally {
        this.loading = false
      }
    },
    handleClose() {
      this.$refs.formRef.clearValidate()
      this.dialogVisible = false
    }
  }
}
</script>

<style scoped>
.dialog-footer {
  text-align: right;
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}
</style>