<template>
  <el-dialog
    :visible.sync="dialogVisible"
    :title="`配置用户角色 - ${user ? user.username : ''}`"
    width="600px"
    @close="handleClose"
  >
    <div v-if="user" class="user-role-config">
      <div class="user-info">
        <el-row :gutter="20">
          <el-col :span="12">
            <div class="info-item">
              <label>用户名:</label>
              <span>{{ user.username }}</span>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="info-item">
              <label>邮箱:</label>
              <span>{{ user.email }}</span>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="info-item">
              <label>状态:</label>
              <el-tag :type="user.status === 'active' ? 'success' : 'danger'">
                {{ user.status === 'active' ? '活跃' : '非活跃' }}
              </el-tag>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="info-item">
              <label>关联员工:</label>
              <span v-if="user.employee">{{ user.employee.name }}</span>
              <span v-else class="text-muted">未关联</span>
            </div>
          </el-col>
        </el-row>
      </div>

      <el-divider />

      <div class="role-section">
        <div class="section-header">
          <h4>角色配置</h4>
          <div class="role-actions">
            <el-button @click="handleSelectAll" size="small">全选</el-button>
            <el-button @click="handleSelectNone" size="small">全不选</el-button>
          </div>
        </div>

        <div v-loading="rolesLoading" class="roles-list">
          <el-checkbox-group v-model="selectedRoles" @change="handleRoleChange">
            <div
              v-for="role in allRoles"
              :key="role.id"
              class="role-item"
            >
              <el-checkbox :label="role.id">
                <div class="role-content">
                  <div class="role-header">
                    <span class="role-name">{{ role.name }}</span>
                    <el-tag
                      :type="role.status === 'active' ? 'success' : 'danger'"
                      size="small"
                    >
                      {{ role.status === 'active' ? '活跃' : '非活跃' }}
                    </el-tag>
                  </div>
                  <div class="role-details">
                    <div class="role-code">编码: {{ role.code }}</div>
                    <div class="role-desc">{{ role.description }}</div>
                  </div>
                </div>
              </el-checkbox>
            </div>
          </el-checkbox-group>
        </div>
      </div>

      <div v-if="userPermissions.length > 0" class="permissions-preview">
        <el-divider />
        <div class="section-header">
          <h4>用户权限预览 ({{ userPermissions.length }})</h4>
          <el-button @click="handleRefreshPermissions" size="small" type="primary">
            刷新权限
          </el-button>
        </div>
        <div class="permissions-list">
          <el-tag
            v-for="permission in userPermissions"
            :key="permission.id"
            class="permission-tag"
            size="small"
            type="info"
          >
            {{ permission.name }}
          </el-tag>
        </div>
      </div>
    </div>

    <span slot="footer" class="dialog-footer">
      <el-button @click="handleClose">取消</el-button>
      <el-button
        type="primary"
        :loading="loading"
        @click="handleSave"
      >
        保存配置
      </el-button>
    </span>
  </el-dialog>
</template>

<script>
export default {
  name: 'UserRoleDialog',
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
      rolesLoading: false,
      allRoles: [],
      selectedRoles: [],
      userPermissions: []
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
    async visible(val) {
      if (val && this.user) {
        await this.fetchAllRoles()
        await this.fetchUserRoles()
        await this.fetchUserPermissions()
      }
    }
  },
  methods: {
    async fetchAllRoles() {
      this.rolesLoading = true
      try {
        // Mock API call - replace with actual API
        this.allRoles = []
        this.$message.success && this.$message.success('获取角色列表成功')
      } catch (error) {
        this.$message.error && this.$message.error('获取角色列表失败')
      } finally {
        this.rolesLoading = false
      }
    },
    async fetchUserRoles() {
      if (!this.user) return
      
      try {
        // Mock API call - replace with actual API
        this.selectedRoles = []
      } catch (error) {
        this.$message.error && this.$message.error('获取用户角色失败')
      }
    },
    async fetchUserPermissions() {
      if (!this.user) return
      
      try {
        // Mock API call - replace with actual API
        this.userPermissions = []
      } catch (error) {
        this.$message.error && this.$message.error('获取用户权限失败')
      }
    },
    handleRoleChange() {
      // This will be called automatically when checkbox-group changes
    },
    handleSelectAll() {
      this.selectedRoles = this.allRoles
        .filter(role => role.status === 'active')
        .map(role => role.id)
    },
    handleSelectNone() {
      this.selectedRoles = []
    },
    async handleRefreshPermissions() {
      await this.fetchUserPermissions()
    },
    async handleSave() {
      if (!this.user) return
      
      this.loading = true
      try {
        // Mock API call - replace with actual API
        // await userApi.assignRoles(this.user.id, this.selectedRoles)
        this.$message.success && this.$message.success('角色配置保存成功')
        
        // 刷新用户权限
        await this.fetchUserPermissions()
        
        this.$emit('success')
        this.handleClose()
      } catch (error) {
        if (error.response?.data?.error) {
          this.$message.error && this.$message.error(error.response.data.error)
        } else {
          this.$message.error && this.$message.error('保存角色配置失败')
        }
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
.user-role-config {
  max-height: 600px;
  overflow-y: auto;
}

.user-info {
  margin-bottom: 20px;
}

.role-section {
  margin-bottom: 20px;
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

.role-actions {
  display: flex;
  gap: 8px;
}

.roles-list {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  padding: 15px;
}

.role-item {
  margin-bottom: 15px;
  padding-bottom: 15px;
  border-bottom: 1px solid #f0f0f0;
}

.role-item:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.role-content {
  margin-left: 24px;
}

.role-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.role-name {
  font-weight: 500;
  font-size: 16px;
  color: #303133;
}

.role-details {
  color: #666;
  font-size: 14px;
}

.role-code {
  margin-bottom: 4px;
}

.role-desc {
  font-style: italic;
}

.permissions-preview {
  margin-top: 20px;
}

.permissions-list {
  max-height: 150px;
  overflow-y: auto;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 10px;
  background-color: #f9f9f9;
  border-radius: 6px;
}

.permission-tag {
  margin: 0;
}

.text-muted {
  color: #999;
  font-style: italic;
}

.dialog-footer {
  text-align: right;
}

:deep(.el-checkbox-group) {
  width: 100%;
}

:deep(.el-checkbox) {
  width: 100%;
  margin-right: 0;
  white-space: normal;
}
</style>