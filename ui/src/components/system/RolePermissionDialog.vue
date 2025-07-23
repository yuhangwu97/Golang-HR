<template>
  <el-dialog
    :visible.sync="dialogVisible"
    :title="role && role.name ? ('配置角色权限 - ' + role.name) : '配置角色权限'"
    width="800px"
    @close="handleClose"
  >
    <div v-if="role" class="permission-config">
      <div class="permission-actions">
        <el-button @click="handleSelectAll">全选</el-button>
        <el-button @click="handleSelectNone">全不选</el-button>
        <el-button type="primary" :loading="loading" @click="handleSave">
          保存配置
        </el-button>
      </div>

      <el-divider />

      <div v-loading="treeLoading" class="permission-tree">
        <div
          v-for="node in permissionTree"
          :key="node.resource"
          class="resource-group"
        >
          <div class="resource-header">
            <el-checkbox
              :model-value="isResourceSelected(node.resource)"
              :indeterminate="isResourceIndeterminate(node.resource)"
              @change="handleResourceChange(node.resource, $event)"
            >
              <strong>{{ resourceLabels[node.resource] || node.resource }}</strong>
            </el-checkbox>
          </div>
          <div class="permissions-list">
            <el-checkbox-group
              v-model="selectedPermissions"
              @change="handlePermissionChange"
            >
              <el-checkbox
                v-for="permission in node.permissions"
                :key="permission.id"
                :value="permission.id"
                :label="permission.id"
              >
                <span class="permission-item">
                  <span class="permission-name">{{ permission.name }}</span>
                  <span class="permission-action">{{ actionLabels[permission.action] || permission.action }}</span>
                  <span class="permission-desc">{{ permission.description }}</span>
                </span>
              </el-checkbox>
            </el-checkbox-group>
          </div>
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
import { roleApi, permissionApi } from '@/services/systemApi'

export default {
  name: 'RolePermissionDialog',
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
      treeLoading: false,
      permissionTree: [],
      selectedPermissions: [],
      // 资源和操作的中文标签
      resourceLabels: {
        employee: '员工管理',
        department: '部门管理',
        position: '职位管理',
        job_level: '职级管理',
        salary: '薪资管理',
        attendance: '考勤管理',
        system: '系统管理',
        user: '用户管理',
        role: '角色管理',
        permission: '权限管理'
      },
      actionLabels: {
        create: '创建',
        read: '查看',
        update: '更新',
        delete: '删除',
        manage: '管理',
        assign: '分配',
        approve: '审批'
      }
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
        if (visible) {
          await this.fetchPermissionTree()
          await this.fetchRolePermissions()
        }
      },
      immediate: true
    }
  },
  methods: {
    async fetchPermissionTree() {
      this.treeLoading = true
      try {
        this.permissionTree = await permissionApi.getPermissionTree()
      } catch (error) {
        this.$message.error('获取权限树失败')
      } finally {
        this.treeLoading = false
      }
    },
    async fetchRolePermissions() {
      if (!this.role) return
      
      try {
        const permissions = await roleApi.getRolePermissions(this.role.id)
        this.selectedPermissions = permissions.map(p => p.id)
      } catch (error) {
        this.$message.error('获取角色权限失败')
      }
    },
    isResourceSelected(resource) {
      const resourcePermissions = this.permissionTree
        .find(node => node.resource === resource)
      if (!resourcePermissions || !resourcePermissions.permissions) return false
      const selectedIds = new Set(this.selectedPermissions)
      return resourcePermissions.permissions.every(p => selectedIds.has(p.id))
    },
    isResourceIndeterminate(resource) {
      const resourcePermissions = this.permissionTree
        .find(node => node.resource === resource)
      if (!resourcePermissions || !resourcePermissions.permissions) return false
      const selectedIds = new Set(this.selectedPermissions)
      const selectedCount = resourcePermissions.permissions.filter(p => selectedIds.has(p.id)).length
      return selectedCount > 0 && selectedCount < resourcePermissions.permissions.length
    },
    handleResourceChange(resource, checked) {
      const resourcePermissions = this.permissionTree
        .find(node => node.resource === resource)
      if (!resourcePermissions || !resourcePermissions.permissions) return
      
      if (checked) {
        // 添加该资源下的所有权限
        const newPermissions = resourcePermissions.permissions.map(p => p.id)
        const currentSet = new Set(this.selectedPermissions)
        newPermissions.forEach(id => currentSet.add(id))
        this.selectedPermissions = Array.from(currentSet)
      } else {
        // 移除该资源下的所有权限
        const removeIds = new Set(resourcePermissions.permissions.map(p => p.id))
        this.selectedPermissions = this.selectedPermissions.filter(id => !removeIds.has(id))
      }
    },
    handlePermissionChange() {
      // This will be called automatically when checkbox-group changes
    },
    handleSelectAll() {
      const allPermissionIds = []
      this.permissionTree.forEach(node => {
        if (node.permissions) {
          node.permissions.forEach(permission => {
            allPermissionIds.push(permission.id)
          })
        }
      })
      this.selectedPermissions = allPermissionIds
    },
    handleSelectNone() {
      this.selectedPermissions = []
    },
    async handleSave() {
      if (!this.role) return
      
      this.loading = true
      try {
        await roleApi.assignPermissions(this.role.id, this.selectedPermissions)
        this.$message.success('权限配置保存成功')
        this.$emit('success')
      } catch (error) {
        if (error.response && error.response.data && error.response.data.error) {
          this.$message.error(error.response.data.error)
        } else {
          this.$message.error('保存权限配置失败')
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
.permission-config {
  max-height: 600px;
  overflow-y: auto;
}

.permission-actions {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}

.resource-group {
  margin-bottom: 20px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  padding: 15px;
}

.resource-header {
  margin-bottom: 10px;
  padding-bottom: 10px;
  border-bottom: 1px solid #f0f0f0;
}

.permissions-list {
  padding-left: 20px;
}

.permission-item {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 300px;
}

.permission-name {
  font-weight: 500;
  min-width: 100px;
}

.permission-action {
  color: #409eff;
  background: #ecf5ff;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  min-width: 50px;
  text-align: center;
}

.permission-desc {
  color: #666;
  font-size: 12px;
  flex: 1;
}

.dialog-footer {
  text-align: right;
}

:deep(.el-checkbox-group) {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

:deep(.el-checkbox) {
  margin-right: 0;
  white-space: nowrap;
}
</style>