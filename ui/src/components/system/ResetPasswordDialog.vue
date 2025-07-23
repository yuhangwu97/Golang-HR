<template>
  <el-dialog
    :visible.sync="dialogVisible"
    :title="`重置密码 - ${user ? user.username : ''}`"
    width="500px"
    @close="handleClose"
  >
    <div v-if="user" class="reset-password">
      <div class="user-info">
        <el-alert
          :title="`即将重置用户 ${user.username} 的密码`"
          type="warning"
          :closable="false"
          show-icon
        />
      </div>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        style="margin-top: 20px;"
      >
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="form.newPassword"
            type="password"
            placeholder="请输入新密码"
            maxlength="50"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            maxlength="50"
            show-password
          />
        </el-form-item>
      </el-form>

      <div class="password-tips">
        <el-alert
          title="密码安全提示"
          type="info"
          :closable="false"
          show-icon
        >
          <ul slot="default">
            <li>密码长度至少6个字符</li>
            <li>建议包含大小写字母、数字和特殊字符</li>
            <li>避免使用简单的密码如123456、password等</li>
            <li>重置后请及时通知用户修改密码</li>
          </ul>
        </el-alert>
      </div>
    </div>

    <span slot="footer" class="dialog-footer">
      <el-button @click="handleClose">取消</el-button>
      <el-button
        type="primary"
        :loading="loading"
        @click="handleSubmit"
      >
        确认重置
      </el-button>
    </span>
  </el-dialog>
</template>

<script>
export default {
  name: 'ResetPasswordDialog',
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
      form: {
        newPassword: '',
        confirmPassword: ''
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
    rules() {
      return {
        newPassword: [
          { required: true, message: '请输入新密码', trigger: 'blur' },
          { min: 6, max: 50, message: '长度在 6 到 50 个字符', trigger: 'blur' }
        ],
        confirmPassword: [
          { required: true, message: '请确认密码', trigger: 'blur' },
          { validator: this.validateConfirmPassword, trigger: 'blur' }
        ]
      }
    }
  },
  watch: {
    visible(val) {
      if (val) {
        this.initForm()
      }
    }
  },
  methods: {
    validateConfirmPassword(rule, value, callback) {
      if (value !== this.form.newPassword) {
        callback(new Error('两次输入的密码不一致'))
      } else {
        callback()
      }
    },
    initForm() {
      Object.assign(this.form, {
        newPassword: '',
        confirmPassword: ''
      })
    },
    async handleSubmit() {
      if (!this.$refs.formRef || !this.user) return
      
      try {
        await this.$refs.formRef.validate()
        this.loading = true
        
        // Mock API call - replace with actual API
        // await userApi.resetPassword(this.user.id, this.form.newPassword)
        this.$message.success && this.$message.success('密码重置成功')
        
        this.$emit('success')
        this.handleClose()
      } catch (error) {
        if (error.response?.data?.error) {
          this.$message.error && this.$message.error(error.response.data.error)
        } else {
          this.$message.error && this.$message.error('密码重置失败')
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
.reset-password {
  padding: 10px 0;
}

.user-info {
  margin-bottom: 20px;
}

.password-tips {
  margin-top: 20px;
}

.password-tips ul {
  margin: 10px 0 0 0;
  padding-left: 20px;
}

.password-tips li {
  margin-bottom: 5px;
  font-size: 13px;
  color: #666;
}

.dialog-footer {
  text-align: right;
}

:deep(.el-alert__content) {
  font-size: 14px;
}
</style>