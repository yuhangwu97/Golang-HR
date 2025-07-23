<template>
  <div class="register-container">
    <div class="register-form-wrapper">
      <div class="register-header">
        <h1>HR管理系统</h1>
        <p>创建新账号</p>
      </div>
      
      <a-form
        :model="formData"
        :rules="rules"
        @finish="handleRegister"
        layout="vertical"
        class="register-form"
      >
        <a-form-item name="name" label="姓名">
          <a-input
            v-model="formData.name"
            placeholder="请输入姓名"
            size="large"
          >
            <UserOutlined slot="prefix" />
          </a-input>
        </a-form-item>
        
        <a-form-item name="email" label="邮箱">
          <a-input
            v-model="formData.email"
            placeholder="请输入邮箱"
            size="large"
          >
            <MailOutlined slot="prefix" />
          </a-input>
        </a-form-item>
        
        <a-form-item name="password" label="密码">
          <a-input-password
            v-model="formData.password"
            placeholder="请输入密码"
            size="large"
          >
            <LockOutlined slot="prefix" />
          </a-input-password>
        </a-form-item>
        
        <a-form-item name="confirmPassword" label="确认密码">
          <a-input-password
            v-model="formData.confirmPassword"
            placeholder="请再次输入密码"
            size="large"
          >
            <LockOutlined slot="prefix" />
          </a-input-password>
        </a-form-item>
        
        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            block
            :loading="authStore.loading"
          >
            注册
          </a-button>
        </a-form-item>
        
        <div class="register-footer">
          <router-link to="/login">已有账号？立即登录</router-link>
        </div>
      </a-form>
    </div>
  </div>
</template>

<script>
import { message } from 'ant-design-vue'
import { UserOutlined, MailOutlined, LockOutlined } from '@ant-design/icons-vue'

export default {
  name: 'RegisterView',
  components: {
    UserOutlined,
    MailOutlined,
    LockOutlined
  },
  data() {
    return {
      formData: {
        name: '',
        email: '',
        password: '',
        confirmPassword: ''
      }
    }
  },
  computed: {
    authStore() {
      return this.$store.state.auth || {}
    },
    rules() {
      return {
        name: [
          { required: true, message: '请输入姓名' },
          { min: 2, message: '姓名至少2个字符' }
        ],
        email: [
          { required: true, message: '请输入邮箱' },
          { type: 'email', message: '邮箱格式不正确' }
        ],
        password: [
          { required: true, message: '请输入密码' },
          { min: 6, message: '密码至少6位' }
        ],
        confirmPassword: [
          { required: true, message: '请确认密码' },
          {
            validator: (_, value) => {
              if (value && value !== this.formData.password) {
                return Promise.reject('两次输入的密码不一致')
              }
              return Promise.resolve()
            }
          }
        ]
      }
    }
  },
  methods: {
    async handleRegister() {
      try {
        const { confirmPassword, ...registerData } = this.formData
        // Mock API call - replace with actual store action
        // await this.$store.dispatch('auth/register', registerData)
        message.success('注册成功')
        this.$router.push('/dashboard')
      } catch (error) {
        message.error(error.response?.data?.message || '注册失败')
      }
    }
  }
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-form-wrapper {
  width: 400px;
  background: white;
  border-radius: 8px;
  padding: 40px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.register-header h1 {
  color: #1890ff;
  margin-bottom: 8px;
}

.register-header p {
  color: #666;
  margin: 0;
}

.register-footer {
  text-align: center;
  margin-top: 16px;
}

.register-footer a {
  color: #1890ff;
}
</style>