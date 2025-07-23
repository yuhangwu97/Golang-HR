<template>
  <div class="position-list">
    <div class="page-header">
      <h2>职位管理</h2>
      <div class="actions">
        <el-button type="primary" @click="showCreateDialog = true">
          <i class="el-icon-plus"></i>
          新增职位
        </el-button>
        <el-button type="success" @click="showOrganizationModal = true">
          <i class="el-icon-s-cooperation"></i>
          组织管理
        </el-button>
        <el-button @click="toggleViewMode">
          <i :class="viewMode === 'tree' ? 'el-icon-menu' : 'el-icon-share'"></i>
          {{ viewMode === 'tree' ? '列表视图' : '树形视图' }}
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
              <div class="stat-label">总职位数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ (statistics && statistics.active_positions) || 0 }}</div>
              <div class="stat-label">活跃职位</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ (statistics && statistics.inactive_positions) || 0 }}</div>
              <div class="stat-label">非活跃职位</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ Object.keys((statistics && statistics.by_department) || {}).length }}</div>
              <div class="stat-label">涉及部门</div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 搜索过滤 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="职位名称">
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入职位名称或编码"
            clearable
            @keyup.enter.native="handleSearch"
          />
        </el-form-item>
        <el-form-item label="所属部门">
          <el-select v-model="searchForm.department_id" placeholder="选择部门" clearable>
            <el-option
              v-for="dept in departments"
              :key="dept.id"
              :label="dept.name"
              :value="dept.id"
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

    <!-- 职位列表/树 -->
    <el-card>
      <!-- 树形视图 -->
      <div v-if="viewMode === 'tree'" v-loading="treeLoading" class="tree-view">
        <el-tree
          :data="positionTree"
          :props="treeProps"
          :expand-on-click-node="false"
          node-key="id"
          default-expand-all
          class="position-tree"
        >
          <div class="tree-node" slot-scope="{ node, data }">
            <div class="node-info">
              <div class="node-main">
                <span class="node-name">{{ data.name }}</span>
                <el-tag size="mini" :type="data.status === 'active' ? 'success' : 'danger'">
                  {{ data.status === 'active' ? '活跃' : '非活跃' }}
                </el-tag>
              </div>
              <div class="node-details">
                <span class="node-code">编码: {{ data.code }}</span>
                <span class="node-department">部门: {{ data.department?.name || '未分配' }}</span>
                <span class="node-level">层级: {{ data.level }}</span>
              </div>
            </div>
            <div class="node-actions">
              <el-button size="mini" @click="addChild(data)">添加下级</el-button>
              <el-button size="mini" @click="handleEdit(data)">编辑</el-button>
              <el-button size="mini" type="danger" @click="handleDelete(data)">删除</el-button>
            </div>
          </div>
        </el-tree>
      </div>

      <!-- 列表视图 -->
      <div v-else>
        <el-table
          v-loading="loading"
          :data="positions"
          stripe
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="name" label="职位名称" width="200" />
          <el-table-column prop="code" label="职位编码" width="150" />
          <el-table-column prop="department.name" label="所属部门" width="120" />
          <el-table-column label="上级职位" width="120">
            <template slot-scope="{ row }">
              <span v-if="row.parent">{{ row.parent.name }}</span>
              <span v-else class="text-muted">无</span>
            </template>
          </el-table-column>
          <el-table-column prop="level" label="层级" width="80" />
          <el-table-column prop="description" label="职位描述" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="100">
            <template slot-scope="{ row }">
              <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
                {{ row.status === 'active' ? '活跃' : '非活跃' }}
              </el-tag>
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
    <PositionDialog
      :visible.sync="showCreateDialog"
      :position="editPosition"
      :parent-position="parentPosition"
      :departments="departments"
      @success="handleDialogSuccess"
    />
  </div>
</template>

<script>
import { positionApi } from '@/services/positionApi'
import { departmentApi } from '@/services/departmentApi'
import PositionDialog from '@/components/organization/PositionDialog.vue'
import { formatDate } from '@/utils/date'

export default {
  name: 'PositionListView',
  components: {
    PositionDialog
  },
  data() {
    return {
      loading: false,
      treeLoading: false,
      showStatistics: false,
      showCreateDialog: false,
      viewMode: 'tree', // 'tree' or 'list'
      positions: [],
      positionTree: [],
      departments: [],
      statistics: null,
      editPosition: null,
      selectedPositions: [],
      parentPosition: null, // 用于添加下级职位
      searchForm: {
        keyword: '',
        department_id: undefined,
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
      treeProps: {
        children: 'children',
        label: 'name'
      }
    }
  },
  mounted() {
    this.fetchPositions()
    this.fetchPositionTree()
    this.fetchDepartments()
    this.fetchStatistics()
  },
  methods: {
    formatDate,
    toggleViewMode() {
      this.viewMode = this.viewMode === 'tree' ? 'list' : 'tree'
      if (this.viewMode === 'tree') {
        this.fetchPositionTree()
      }
    },
    async fetchPositionTree() {
      this.treeLoading = true
      try {
        const response = await positionApi.getPositionTree()
        this.positionTree = response.data || []
      } catch (error) {
        this.$message.error('获取职位树失败')
      } finally {
        this.treeLoading = false
      }
    },
    addChild(parentPosition) {
      this.parentPosition = parentPosition
      this.editPosition = null
      this.showCreateDialog = true
    },
    async fetchPositions() {
      this.loading = true
      try {
        const response = await positionApi.getPositions({
          ...this.searchForm,
          page: this.pagination.page,
          page_size: this.pagination.page_size
        })
        
        // 处理API响应数据结构
        const responseData = response.data || response
        this.positions = responseData.data || []
        this.pagination.total_items = responseData.total_items || 0
        this.pagination.total_pages = responseData.total_pages || 0
      } catch (error) {
        console.error('获取职位列表失败:', error)
        this.$message.error('获取职位列表失败')
      } finally {
        this.loading = false
      }
    },
    async fetchDepartments() {
      try {
        const response = await departmentApi.getDepartments()
        // 处理API响应数据结构
        const responseData = response.data || response
        this.departments = responseData.data || responseData || []
      } catch (error) {
        console.error('获取部门列表失败:', error)
        this.$message.error('获取部门列表失败')
      }
    },
    async fetchStatistics() {
      try {
        const response = await positionApi.getPositionStatistics()
        // 处理API响应数据结构
        const responseData = response.data || response
        this.statistics = responseData.data || responseData || {}
      } catch (error) {
        console.error('获取统计信息失败:', error)
        // 如果统计接口失败，使用基础数据计算统计信息
        this.statistics = this.calculateBasicStatistics()
      }
    },
    handleSearch() {
      this.pagination.page = 1
      this.fetchPositions()
    },
    handleReset() {
      Object.assign(this.searchForm, {
        keyword: '',
        department_id: undefined,
        status: '',
        page: 1,
        page_size: 20
      })
      this.pagination.page = 1
      this.fetchPositions()
    },
    handleEdit(position) {
      this.editPosition = position
      this.showCreateDialog = true
    },
    async handleDelete(position) {
      try {
        await this.$confirm(
          `确定要删除职位"${position.name}"吗？`,
          '确认删除',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        await positionApi.deletePosition(position.id)
        this.$message.success('删除成功')
        this.fetchPositions()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除失败')
        }
      }
    },
    handleSelectionChange(selected) {
      this.selectedPositions = selected
    },
    handleSizeChange(size) {
      this.pagination.page_size = size
      this.pagination.page = 1
      this.fetchPositions()
    },
    handleCurrentChange(page) {
      this.pagination.page = page
      this.fetchPositions()
    },
    handleDialogSuccess() {
      this.showCreateDialog = false
      this.editPosition = null
      this.parentPosition = null
      this.fetchPositions()
      if (this.viewMode === 'tree') {
        this.fetchPositionTree()
      }
      if (this.showStatistics) {
        this.fetchStatistics()
      }
    },
    // 计算基础统计信息
    calculateBasicStatistics() {
      const total = this.positions.length
      const activePositions = this.positions.filter(pos => pos.status === 'active').length
      const inactivePositions = total - activePositions
      const byDepartment = {}
      
      this.positions.forEach(pos => {
        if (pos.department) {
          byDepartment[pos.department.name] = (byDepartment[pos.department.name] || 0) + 1
        }
      })

      return {
        total,
        active_positions: activePositions,
        inactive_positions: inactivePositions,
        by_department: byDepartment
      }
    }
  }
}
</script>

<style scoped>
.tree-view {
  padding: 20px;
  background-color: #fafafa;
  border-radius: 8px;
}

.position-tree {
  max-width: 100%;
  background-color: white;
  border-radius: 6px;
  padding: 16px;
}

.position-tree :deep(.el-tree-node__content) {
  height: auto;
  min-height: 60px;
  padding: 0;
  border-bottom: 1px solid #f0f0f0;
}

.position-tree :deep(.el-tree-node__content:hover) {
  background-color: transparent;
}

.position-tree :deep(.el-tree-node__expand-icon) {
  padding: 8px;
  font-size: 14px;
}

.position-tree :deep(.el-tree-node__children) {
  margin-left: 20px;
  padding-left: 20px;
  border-left: 2px solid #e4e7ed;
}

.tree-node {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  width: 100%;
  padding: 12px 16px;
  min-height: 60px;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.tree-node:hover {
  background-color: #f5f7fa;
}

.node-info {
  flex: 1;
  min-width: 0;
  margin-right: 16px;
}

.node-main {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 6px;
  flex-wrap: wrap;
}

.node-name {
  font-weight: 500;
  color: #303133;
  font-size: 14px;
  flex-shrink: 0;
}

.node-details {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: #909399;
  flex-wrap: wrap;
  line-height: 1.5;
}

.node-details > span {
  white-space: nowrap;
}

.node-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
  align-items: flex-start;
}

.node-actions .el-button {
  padding: 4px 8px;
  font-size: 12px;
}

.text-muted {
  color: #999;
  font-style: italic;
}

.position-list {
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

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .actions {
    flex-wrap: wrap;
    justify-content: center;
  }
  
  .tree-node {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .node-info {
    margin-right: 0;
  }
  
  .node-actions {
    justify-content: flex-start;
    flex-wrap: wrap;
  }
  
  .node-details {
    flex-direction: column;
    gap: 8px;
  }
  
  .statistics-cards .el-col {
    margin-bottom: 10px;
  }
}

@media (max-width: 480px) {
  .position-list {
    padding: 10px;
  }
  
  .tree-view {
    padding: 10px;
  }
  
  .position-tree {
    padding: 12px;
  }
  
  .node-actions .el-button {
    padding: 6px 12px;
    font-size: 12px;
  }
}
</style>