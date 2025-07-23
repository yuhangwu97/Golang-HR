<template>
  <el-dialog
    :visible.sync="dialogVisible"
    :title="isEdit ? '编辑用户' : '新增用户'"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="用户名" prop="username">
        <el-input
          v-model="form.username"
          placeholder="请输入用户名"
          maxlength="50"
        />
      </el-form-item>
      
      <el-form-item label="邮箱" prop="email">
        <el-input
          v-model="form.email"
          placeholder="请输入邮箱"
          maxlength="100"
        />
      </el-form-item>
      
      <el-form-item v-if="!isEdit" label="密码" prop="password">
        <el-input
          v-model="form.password"
          type="password"
          placeholder="请输入密码"
          maxlength="50"
          show-password
        />
      </el-form-item>
      
      <el-form-item label="关联员工" prop="employee_id">
        <el-select
          v-model="form.employee_id"
          placeholder="请选择关联的员工（可选）"
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
import { systemApiService } from '@/services/systemApi'
import { employeeApi } from '@/services/employeeApi'

export default {
  name: 'UserDialog',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    user: {
      type: Object,
      default: null
    }
  },
  data() {
    return {
      loading: false,
      employees: [],
      form: {
        username: '',
        email: '',
        password: '',
        employee_id: undefined,
        status: 'active'
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
      return !!this.user
    },
    rules() {
      return {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' },
          { min: 3, max: 50, message: '长度在 3 到 50 个字符', trigger: 'blur' },
          { pattern: /^[A-Za-z0-9_-]+$/, message: '用户名只能包含字母、数字、下划线和连字符', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '请输入邮箱', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
          { max: 100, message: '长度不能超过 100 个字符', trigger: 'blur' }
        ],
        password: [
          { required: !this.isEdit, message: '请输入密码', trigger: 'blur' },
          { min: 6, max: 50, message: '长度在 6 到 50 个字符', trigger: 'blur' }
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
        this.fetchEmployees()
      }
    }
  },
  mounted() {
    this.fetchEmployees()
  },
  methods: {
    async fetchEmployees() {
      try {
        const response = await employeeApi.getAllEmployees()
        this.employees = response.data || []
      } catch (error) {
        this.$message.error('获取员工列表失败')
      }
    },
    initForm() {
      if (this.user) {
        Object.assign(this.form, {
          username: this.user.username,
          email: this.user.email,
          password: '', // 编辑时不显示密码
          employee_id: this.user.employee_id,
          status: this.user.status
        })
      } else {
        Object.assign(this.form, {
          username: '',
          email: '',
          password: '',
          employee_id: undefined,
          status: 'active'
        })
      }
    },
    async handleSubmit() {
      try {
        await this.$refs.formRef.validate()
        this.loading = true
        
        const submitData = { ...this.form }
        // 如果是编辑且密码为空，删除密码字段
        if (this.isEdit && !submitData.password) {
          delete submitData.password
        }
        
        if (this.isEdit && this.user) {
          await systemApiService.updateUser(this.user.id, submitData)
          this.$message.success('更新成功')
        } else {
          await systemApiService.createUser(submitData)
          this.$message.success('创建成功')
        }
        
        this.$emit('success')
        this.handleClose()
      } catch (error) {
        if (error.response?.data?.error) {
          this.$message.error(error.response.data.error)
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
</style>