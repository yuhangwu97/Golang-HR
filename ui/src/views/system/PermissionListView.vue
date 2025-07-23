<template>
  <div class="permission-list">
    <div class="page-header">
      <h2>权限管理</h2>
      <div class="actions">
        <el-button type="primary" @click="showCreateDialog = true">
          <i class="el-icon-plus"></i>
          新增权限
        </el-button>
        <el-button @click="showTreeView = !showTreeView">
          <i class="el-icon-menu"></i>
          {{ showTreeView ? '列表视图' : '树形视图' }}
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
              <div class="stat-label">总权限数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ Object.keys((statistics && statistics.by_resource) || {}).length }}</div>
              <div class="stat-label">资源类型</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ Object.keys((statistics && statistics.by_action) || {}).length }}</div>
              <div class="stat-label">操作类型</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ activePermissions }}</div>
              <div class="stat-label">活跃权限</div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 搜索过滤 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="权限名称">
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入权限名称或编码"
            clearable
            @keyup.enter.native="handleSearch"
          />
        </el-form-item>
        <el-form-item label="资源类型">
          <el-select v-model="searchForm.resource" placeholder="选择资源类型" clearable>
            <el-option
              v-for="(label, value) in resourceLabels"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="操作类型">
          <el-select v-model="searchForm.action" placeholder="选择操作类型" clearable>
            <el-option
              v-for="(label, value) in actionLabels"
              :key="value"
              :label="label"
              :value="value"
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

    <!-- 内容区域 -->
    <el-card>
      <!-- 树形视图 -->
      <div v-if="showTreeView" v-loading="treeLoading" class="tree-view">
        <div
          v-for="node in permissionTree"
          :key="node.resource"
          class="resource-group"
        >
          <div class="resource-header">
            <h3>{{ resourceLabels[node.resource] || node.resource }}</h3>
            <span class="permission-count">{{ node.permissions.length }} 个权限</span>
          </div>
          <div class="permissions-grid">
            <el-card
              v-for="permission in node.permissions"
              :key="permission.id"
              class="permission-card"
              shadow="hover"
            >
              <div class="permission-info">
                <div class="permission-header">
                  <span class="permission-name">{{ permission.name }}</span>
                  <el-tag
                    :type="permission.status === 'active' ? 'success' : 'danger'"
                    size="small"
                  >
                    {{ permission.status === 'active' ? '活跃' : '非活跃' }}
                  </el-tag>
                </div>
                <div class="permission-details">
                  <div class="permission-code">编码: {{ permission.code }}</div>
                  <div class="permission-action">
                    操作: {{ actionLabels[permission.action] || permission.action }}
                  </div>
                  <div class="permission-desc">{{ permission.description }}</div>
                </div>
                <div class="permission-actions">
                  <el-button size="small" @click="handleEdit(permission)">编辑</el-button>
                  <el-button
                    size="small"
                    type="danger"
                    @click="handleDelete(permission)"
                  >
                    删除
                  </el-button>
                </div>
              </div>
            </el-card>
          </div>
        </div>
      </div>

      <!-- 列表视图 -->
      <div v-else>
        <el-table
          v-loading="loading"
          :data="permissions"
          stripe
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="name" label="权限名称" width="200" />
          <el-table-column prop="code" label="权限编码" width="150" />
          <el-table-column prop="resource" label="资源类型" width="120">
            <template slot-scope="{ row }">
              {{ resourceLabels[row.resource] || row.resource }}
            </template>
          </el-table-column>
          <el-table-column prop="action" label="操作类型" width="120">
            <template slot-scope="{ row }">
              <el-tag type="info" size="small">
                {{ actionLabels[row.action] || row.action }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="权限描述" show-overflow-tooltip />
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
          <el-table-column label="操作" width="180" fixed="right">
            <template slot-scope="{ row }">
              <el-button size="small" @click="handleEdit(row)">编辑</el-button>
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
      </div>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <PermissionDialog
      :visible.sync="showCreateDialog"
      :permission="editPermission"
      @success="handleDialogSuccess"
    />
  </div>
</template>

<script>
import { permissionApi } from '@/services/systemApi'
import PermissionDialog from '@/components/system/PermissionDialog.vue'
import { formatDate } from '@/utils/date'

export default {
  name: 'PermissionListView',
  components: {
    PermissionDialog
  },
  data() {
    return {
      loading: false,
      treeLoading: false,
      showStatistics: false,
      showTreeView: false,
      showCreateDialog: false,
      permissions: [],
      permissionTree: [],
      statistics: null,
      editPermission: null,
      selectedPermissions: [],
      searchForm: {
        keyword: '',
        resource: '',
        action: '',
        status: '',
        page: 1,
        page_size: 20
      },
      pagination: {
        page: 1,
        page_size: 20,
        total_items: 0,
        total_pages: 0
      },
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
    activePermissions() {
      return (this.statistics && this.statistics.by_status && this.statistics.by_status.active) || 0
    }
  },
  mounted() {
    this.fetchPermissions()
    this.fetchPermissionTree()
    this.fetchStatistics()
  },
  methods: {
    formatDate,
    async fetchPermissions() {
      this.loading = true
      try {
        const response = await permissionApi.getPermissions({
          ...this.searchForm,
          page: this.pagination.page,
          page_size: this.pagination.page_size
        })
        this.permissions = response.data
        this.pagination.total_items = response.total_items
        this.pagination.total_pages = response.total_pages
      } catch (error) {
        this.$message.error('获取权限列表失败')
      } finally {
        this.loading = false
      }
    },
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
    async fetchStatistics() {
      try {
        this.statistics = await permissionApi.getPermissionStatistics()
      } catch (error) {
        this.$message.error('获取统计信息失败')
      }
    },
    handleSearch() {
      this.pagination.page = 1
      this.fetchPermissions()
    },
    handleReset() {
      Object.assign(this.searchForm, {
        keyword: '',
        resource: '',
        action: '',
        status: '',
        page: 1,
        page_size: 20
      })
      this.pagination.page = 1
      this.fetchPermissions()
    },
    handleEdit(permission) {
      this.editPermission = permission
      this.showCreateDialog = true
    },
    async handleDelete(permission) {
      try {
        await this.$confirm(
          `确定要删除权限"${permission.name}"吗？`,
          '确认删除',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        await permissionApi.deletePermission(permission.id)
        this.$message.success('删除成功')
        this.fetchPermissions()
        if (this.showTreeView) {
          this.fetchPermissionTree()
        }
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除失败')
        }
      }
    },
    handleSelectionChange(selected) {
      this.selectedPermissions = selected
    },
    handleSizeChange(size) {
      this.pagination.page_size = size
      this.pagination.page = 1
      this.fetchPermissions()
    },
    handleCurrentChange(page) {
      this.pagination.page = page
      this.fetchPermissions()
    },
    handleDialogSuccess() {
      this.showCreateDialog = false
      this.editPermission = null
      this.fetchPermissions()
      if (this.showTreeView) {
        this.fetchPermissionTree()
      }
      if (this.showStatistics) {
        this.fetchStatistics()
      }
    }
  }
}
</script>

<style scoped>
.permission-list {
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

/* 树形视图样式 */
.tree-view {
  padding: 20px 0;
}

.resource-group {
  margin-bottom: 30px;
}

.resource-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 2px solid #e0e0e0;
}

.resource-header h3 {
  margin: 0;
  color: #303133;
}

.permission-count {
  color: #909399;
  font-size: 14px;
}

.permissions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 15px;
}

.permission-card {
  transition: all 0.3s ease;
}

.permission-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.permission-info {
  padding: 10px;
}

.permission-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.permission-name {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.permission-details {
  margin-bottom: 15px;
}

.permission-code,
.permission-action {
  font-size: 12px;
  color: #606266;
  margin-bottom: 4px;
}

.permission-desc {
  font-size: 13px;
  color: #909399;
  line-height: 1.4;
  margin-top: 8px;
}

.permission-actions {
  display: flex;
  gap: 8px;
}
</style>