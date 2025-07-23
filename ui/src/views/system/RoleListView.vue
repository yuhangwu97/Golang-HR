<template>
  <div class="role-list">
    <div class="page-header">
      <h2>角色管理</h2>
      <div class="actions">
        <el-button type="primary" @click="showCreateDialog = true">
          <i class="el-icon-plus"></i>
          新增角色
        </el-button>
        <el-button @click="showStatistics = !showStatistics">
          <i class="el-icon-s-data"></i>
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
              <div class="stat-label">总角色数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ activeRoles }}</div>
              <div class="stat-label">活跃角色</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ inactiveRoles }}</div>
              <div class="stat-label">非活跃角色</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ totalUsers }}</div>
              <div class="stat-label">用户总数</div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 搜索过滤 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="角色名称">
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入角色名称或编码"
            clearable
            @keyup.enter.native="handleSearch"
          />
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

    <!-- 角色表格 -->
    <el-card>
      <el-table
        v-loading="loading"
        :data="roles"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="角色名称" width="200" />
        <el-table-column prop="code" label="角色编码" width="150" />
        <el-table-column prop="description" label="角色描述" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template slot-scope="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? '活跃' : '非活跃' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template slot-scope="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template slot-scope="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="warning" @click="handlePermissions(row)">
              权限配置
            </el-button>
            <el-button size="small" type="info" @click="handleUsers(row)">
              用户列表
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
          :current-page.sync="pagination.page"
          :page-size.sync="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total_items"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <RoleDialog
      :visible.sync="showCreateDialog"
      :role="editRole"
      @success="handleDialogSuccess"
    />

    <!-- 权限配置对话框 -->
    <RolePermissionDialog
      :visible.sync="showPermissionDialog"
      :role="selectedRole"
      @success="handlePermissionSuccess"
    />

    <!-- 用户列表对话框 -->
    <RoleUserDialog
      :visible.sync="showUserDialog"
      :role="selectedRole"
    />
  </div>
</template>

<script>
import { roleApi } from '@/services/systemApi'
import RoleDialog from '@/components/system/RoleDialog.vue'
import RolePermissionDialog from '@/components/system/RolePermissionDialog.vue'
import RoleUserDialog from '@/components/system/RoleUserDialog.vue'
import { formatDate } from '@/utils/date'

export default {
  name: 'RoleListView',
  components: {
    RoleDialog,
    RolePermissionDialog,
    RoleUserDialog
  },
  data() {
    return {
      loading: false,
      showStatistics: false,
      showCreateDialog: false,
      showPermissionDialog: false,
      showUserDialog: false,
      roles: [],
      statistics: null,
      editRole: null,
      selectedRole: null,
      selectedRoles: [],
      searchForm: {
        keyword: '',
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
    activeRoles() {
      return (this.statistics && this.statistics.by_status && this.statistics.by_status.active) || 0
    },
    inactiveRoles() {
      return (this.statistics && this.statistics.by_status && this.statistics.by_status.inactive) || 0
    },
    totalUsers() {
      if (!this.statistics || !this.statistics.user_counts) return 0
      return Object.values(this.statistics.user_counts).reduce((sum, count) => sum + count, 0)
    }
  },
  mounted() {
    this.fetchRoles()
    this.fetchStatistics()
  },
  methods: {
    formatDate,
    async fetchRoles() {
      this.loading = true
      try {
        const response = await roleApi.getRoles({
          ...this.searchForm,
          page: this.pagination.page,
          page_size: this.pagination.page_size
        })
        this.roles = response.data
        this.pagination.total_items = response.total_items
        this.pagination.total_pages = response.total_pages
      } catch (error) {
        this.$message.error('获取角色列表失败')
      } finally {
        this.loading = false
      }
    },
    async fetchStatistics() {
      try {
        this.statistics = await roleApi.getRoleStatistics()
      } catch (error) {
        this.$message.error('获取统计信息失败')
      }
    },
    handleSearch() {
      this.pagination.page = 1
      this.fetchRoles()
    },
    handleReset() {
      Object.assign(this.searchForm, {
        keyword: '',
        status: '',
        page: 1,
        page_size: 20
      })
      this.pagination.page = 1
      this.fetchRoles()
    },
    handleEdit(role) {
      this.editRole = role
      this.showCreateDialog = true
    },
    handlePermissions(role) {
      this.selectedRole = role
      this.showPermissionDialog = true
    },
    handleUsers(role) {
      this.selectedRole = role
      this.showUserDialog = true
    },
    async handleDelete(role) {
      try {
        await this.$confirm(
          `确定要删除角色"${role.name}"吗？`,
          '确认删除',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        await roleApi.deleteRole(role.id)
        this.$message.success('删除成功')
        this.fetchRoles()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除失败')
        }
      }
    },
    handleSelectionChange(selected) {
      this.selectedRoles = selected
    },
    handleSizeChange(size) {
      this.pagination.page_size = size
      this.pagination.page = 1
      this.fetchRoles()
    },
    handleCurrentChange(page) {
      this.pagination.page = page
      this.fetchRoles()
    },
    handleDialogSuccess() {
      this.showCreateDialog = false
      this.editRole = null
      this.fetchRoles()
      if (this.showStatistics) {
        this.fetchStatistics()
      }
    },
    handlePermissionSuccess() {
      this.showPermissionDialog = false
      this.selectedRole = null
    }
  }
}
</script>

<style scoped>
.role-list {
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
</style>