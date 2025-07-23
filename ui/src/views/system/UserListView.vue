<template>
  <div class="user-list">
    <div class="page-header">
      <h2>用户管理</h2>
      <div class="actions">
        <el-button type="primary" @click="showCreateDialog = true">
          <i class="el-icon-plus"></i>
          新增用户
        </el-button>
        <el-button @click="showStatistics = !showStatistics">
          <i class="el-icon-data-analysis"></i>
          统计信息
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div v-if="showStatistics" class="statistics-cards">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ (statistics && statistics.total) || 0 }}</div>
              <div class="stat-label">总用户数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ activeUsers }}</div>
              <div class="stat-label">活跃用户</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ inactiveUsers }}</div>
              <div class="stat-label">非活跃用户</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ Object.keys((statistics && statistics.role_counts) || {}).length }}</div>
              <div class="stat-label">角色类型</div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 搜索过滤 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="用户名称">
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入用户名或邮箱"
            clearable
            @keyup.enter.native="handleSearch"
          />
        </el-form-item>
        <el-form-item label="角色筛选">
          <el-select v-model="searchForm.role_id" placeholder="选择角色" clearable>
            <el-option
              v-for="role in roles"
              :key="role.id"
              :label="role.name"
              :value="role.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="选择状态" clearable>
            <el-option label="活跃" value="active" />
            <el-option label="非活跃" value="inactive" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <i class="el-icon-search"></i>
            搜索
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 用户表格 -->
    <el-card>
      <el-table
        v-loading="loading"
        :data="users"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="employee" label="关联员工" width="120">
          <template slot-scope="{ row }">
            <el-tag v-if="row.employee" type="success" size="small">
              {{ row.employee.name }}
            </el-tag>
            <span v-else class="text-muted">未关联</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template slot-scope="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? '活跃' : '非活跃' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_login_at" label="最后登录" width="180">
          <template slot-scope="{ row }">
            <span v-if="row.last_login_at">
              {{ formatDate(row.last_login_at) }}
            </span>
            <span v-else class="text-muted">从未登录</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template slot-scope="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template slot-scope="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="warning" @click="handleRoles(row)">
              角色配置
            </el-button>
            <el-button size="small" type="info" @click="handleResetPassword(row)">
              重置密码
            </el-button>
            <el-button
              size="small"
              type="danger"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          :current-page="pagination.page"
          :page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total_items"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <UserDialog
      :visible.sync="showCreateDialog"
      :user="editUser"
      @success="handleDialogSuccess"
    />

    <!-- 角色配置对话框 -->
    <UserRoleDialog
      :visible.sync="showRoleDialog"
      :user="selectedUser"
      @success="handleRoleSuccess"
    />

    <!-- 重置密码对话框 -->
    <ResetPasswordDialog
      :visible.sync="showResetDialog"
      :user="selectedUser"
      @success="handleResetSuccess"
    />
  </div>
</template>

<script>
import { systemApiService } from '@/services/systemApi'
import UserDialog from '@/components/system/UserDialog.vue'
import UserRoleDialog from '@/components/system/UserRoleDialog.vue'
import ResetPasswordDialog from '@/components/system/ResetPasswordDialog.vue'

export default {
  name: 'UserListView',
  components: {
    UserDialog,
    UserRoleDialog,
    ResetPasswordDialog
  },
  data() {
    return {
      loading: false,
      showStatistics: false,
      showCreateDialog: false,
      showRoleDialog: false,
      showResetDialog: false,
      users: [],
      roles: [],
      statistics: null,
      editUser: null,
      selectedUser: null,
      selectedUsers: [],
      searchForm: {
        keyword: '',
        role_id: undefined,
        status: '',
        page: 1,
        page_size: 20
      },
      pagination: {
        page: 1,
        page_size: 20,
        total_items: 0,
        total_pages: 0
      }
    }
  },
  computed: {
    activeUsers() {
      return (this.statistics && this.statistics.by_status && this.statistics.by_status.active) || 0
    },
    inactiveUsers() {
      return (this.statistics && this.statistics.by_status && this.statistics.by_status.inactive) || 0
    }
  },
  mounted() {
    this.fetchUsers()
    this.fetchRoles()
    this.fetchStatistics()
  },
  methods: {
    formatDate(date) {
      if (!date) return ''
      return new Date(date).toLocaleString('zh-CN')
    },
    
    async fetchUsers() {
      this.loading = true
      try {
        const response = await systemApiService.getUsers({
          ...this.searchForm,
          page: this.pagination.page,
          page_size: this.pagination.page_size
        })
        this.users = response.data || []
        this.pagination.total_items = response.total_items || 0
        this.pagination.total_pages = response.total_pages || 0
      } catch (error) {
        this.$message.error('获取用户列表失败')
      } finally {
        this.loading = false
      }
    },

    async fetchRoles() {
      try {
        const response = await systemApiService.getRoles({ page: 1, page_size: 100 })
        this.roles = response.data || []
      } catch (error) {
        this.$message.error('获取角色列表失败')
      }
    },

    async fetchStatistics() {
      try {
        this.statistics = await systemApiService.getUserStatistics()
      } catch (error) {
        this.$message.error('获取统计信息失败')
      }
    },

    handleSearch() {
      this.pagination.page = 1
      this.fetchUsers()
    },

    handleReset() {
      Object.assign(this.searchForm, {
        keyword: '',
        role_id: undefined,
        status: '',
        page: 1,
        page_size: 20
      })
      this.pagination.page = 1
      this.fetchUsers()
    },

    handleEdit(user) {
      this.editUser = user
      this.showCreateDialog = true
    },

    handleRoles(user) {
      this.selectedUser = user
      this.showRoleDialog = true
    },

    handleResetPassword(user) {
      this.selectedUser = user
      this.showResetDialog = true
    },

    async handleDelete(user) {
      try {
        await this.$confirm(
          `确定要删除用户"${user.username}"吗？`,
          '确认删除',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        await systemApiService.deleteUser(user.id)
        this.$message.success('删除成功')
        this.fetchUsers()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除失败')
        }
      }
    },

    handleSelectionChange(selected) {
      this.selectedUsers = selected
    },

    handleSizeChange(size) {
      this.pagination.page_size = size
      this.pagination.page = 1
      this.fetchUsers()
    },

    handleCurrentChange(page) {
      this.pagination.page = page
      this.fetchUsers()
    },

    handleDialogSuccess() {
      this.showCreateDialog = false
      this.editUser = null
      this.fetchUsers()
      if (this.showStatistics) {
        this.fetchStatistics()
      }
    },

    handleRoleSuccess() {
      this.showRoleDialog = false
      this.selectedUser = null
    },

    handleResetSuccess() {
      this.showResetDialog = false
      this.selectedUser = null
    }
  }
}
</script>

<style scoped>
.user-list {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

.actions {
  display: flex;
  gap: 10px;
}

.statistics-cards {
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.search-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}

.text-muted {
  color: #999;
  font-style: italic;
}
</style>