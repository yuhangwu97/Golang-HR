<template>
  <el-dialog
    :visible.sync="dialogVisible"
    :title="role && role.name ? ('角色用户列表 - ' + role.name) : '角色用户列表'"
    width="700px"
    @close="handleClose"
  >
    <div v-if="role" class="role-users">
      <div class="header-info">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="角色名称">{{ role.name }}</el-descriptions-item>
          <el-descriptions-item label="角色编码">{{ role.code }}</el-descriptions-item>
          <el-descriptions-item label="角色描述" :span="2">{{ role.description }}</el-descriptions-item>
        </el-descriptions>
      </div>

      <el-divider />

      <div class="users-section">
        <div class="section-header">
          <h4>拥有此角色的用户 ({{ users.length }})</h4>
          <el-button
            v-loading="loading"
            @click="fetchUsers"
            size="small"
            type="primary"
          >
            刷新
          </el-button>
        </div>

        <div v-loading="loading" class="users-list">
          <el-empty v-if="users.length === 0" description="暂无用户拥有此角色" />
          <div v-else class="user-cards">
            <el-card
              v-for="user in users"
              :key="user.id"
              class="user-card"
              shadow="hover"
            >
              <div class="user-info">
                <div class="user-basic">
                  <div class="user-name">{{ user.username }}</div>
                  <div class="user-email">{{ user.email }}</div>
                </div>
                <div class="user-meta">
                  <el-tag
                    :type="user.status === 'active' ? 'success' : 'danger'"
                    size="small"
                  >
                    {{ user.status === 'active' ? '活跃' : '非活跃' }}
                  </el-tag>
                  <div class="user-dates">
                    <div class="created-date">
                      创建时间: {{ formatDateShort(user.created_at) }}
                    </div>
                    <div v-if="user.last_login_at" class="login-date">
                      最后登录: {{ formatDateShort(user.last_login_at) }}
                    </div>
                  </div>
                </div>
              </div>
              <div class="user-employee" v-if="user.employee">
                <el-tag type="info" size="small">
                  关联员工: {{ user.employee.name }}
                </el-tag>
              </div>
            </el-card>
          </div>
        </div>
      </div>
    </div>

    <span slot="footer" class="dialog-footer">
      <el-button @click="handleClose">关闭</el-button>
    </span>
  </el-dialog>
</template>

<script>
import { roleApi } from '@/services/systemApi'
import { formatDateShort } from '@/utils/date'

export default {
  name: 'RoleUserDialog',
  props: {
    value: {
      type: Boolean,
      default: false
    },
    role: {
      type: Object,
      default: null
    }
  },
  data() {
    return {
      loading: false,
      users: []
    }
  },
  computed: {
    dialogVisible: {
      get() {
        return this.value
      },
      set(value) {
        this.$emit('input', value)
      }
    }
  },
  watch: {
    value: {
      async handler(visible) {
        if (visible && this.role) {
          await this.fetchUsers()
        }
      },
      immediate: true
    }
  },
  methods: {
    formatDateShort,
    async fetchUsers() {
      if (!this.role) return
      
      this.loading = true
      try {
        this.users = await roleApi.getRoleUsers(this.role.id)
      } catch (error) {
        this.$message.error('获取角色用户列表失败')
      } finally {
        this.loading = false
      }
    },
    handleClose() {
      this.dialogVisible = false
    }
  }
}
</script>

<style scoped>
.role-users {
  max-height: 600px;
  overflow-y: auto;
}

.header-info {
  margin-bottom: 20px;
}

.users-section {
  margin-top: 20px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.section-header h4 {
  margin: 0;
  color: #303133;
}

.users-list {
  min-height: 200px;
}

.user-cards {
  display: grid;
  gap: 12px;
}

.user-card {
  transition: all 0.3s ease;
}

.user-card:hover {
  transform: translateY(-2px);
}

.user-info {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.user-basic {
  flex: 1;
}

.user-name {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.user-email {
  font-size: 14px;
  color: #606266;
}

.user-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
}

.user-dates {
  text-align: right;
  font-size: 12px;
  color: #909399;
}

.created-date,
.login-date {
  margin-bottom: 2px;
}

.user-employee {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid #f0f0f0;
}

.dialog-footer {
  text-align: right;
}
</style>