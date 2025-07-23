<template>
  <div class="login-container">
    <div class="login-box">
      <h1>HR管理系统</h1>
      <el-form
        ref="loginForm"
        :model="form"
        :rules="rules"
        class="login-form"
        @submit.native.prevent="handleSubmit"
      >
        <el-form-item prop="email">
          <el-input
            v-model="form.email"
            placeholder="邮箱"
            size="large"
            prefix-icon="el-icon-user"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            size="large"
            prefix-icon="el-icon-lock"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-button"
            @click="handleSubmit"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'

export default {
  name: 'LoginView',
  data() {
    return {
      form: {
        email: '',
        password: ''
      },
      rules: {
        email: [
          { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' }
        ]
      }
    }
  },
  computed: {
    ...mapGetters('auth', ['loading'])
  },
  methods: {
    ...mapActions('auth', ['login']),
    
    async handleSubmit() {
      this.$refs.loginForm.validate(async (valid) => {
        if (valid) {
          try {
            await this.login(this.form)
            this.$message.success('登录成功')
            
            // 确保路由跳转在下一个tick执行，避免重定向冲突
            this.$nextTick(() => {
              this.$router.push('/dashboard')
            })
          } catch (error) {
            console.error('登录失败:', error)
            this.$message.error('登录失败: ' + (error.response?.data?.message || error.message || '未知错误'))
          }
        }
      })
    }
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-box {
  width: 400px;
  padding: 40px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.login-box h1 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
  font-size: 24px;
}

.login-form {
  width: 100%;
}

.login-button {
  width: 100%;
}
</style>