<template>
  <div class="assignment-view">
    <div class="view-header">
      <h3>人员分配管理</h3>
      <div class="header-actions">
        <el-button size="small" @click="refreshData">
          <i class="el-icon-refresh"></i>
          刷新
        </el-button>
      </div>
    </div>

    <div v-if="!selectedUnit" class="empty-state">
      <div class="empty-content">
        <i class="el-icon-user"></i>
        <p>请从左侧选择一个组织单元查看人员分配</p>
      </div>
    </div>

    <div v-else class="assignment-content">
      <!-- 分配统计卡片 -->
      <el-row :gutter="20" class="stats-cards">
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-item">
              <div class="stat-value">{{ assignmentStats.total || 0 }}</div>
              <div class="stat-label">总分配数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-item">
              <div class="stat-value primary">{{ assignmentStats.primary || 0 }}</div>
              <div class="stat-label">主要分配</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-item">
              <div class="stat-value additional">{{ assignmentStats.additional || 0 }}</div>
              <div class="stat-label">额外分配</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-item">
              <div class="stat-value temporary">{{ assignmentStats.temporary || 0 }}</div>
              <div class="stat-label">临时分配</div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 筛选工具栏 -->
      <el-card class="filter-card">
        <div class="filter-toolbar">
          <div class="filter-left">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索员工姓名或工号"
              size="small"
              style="width: 200px"
              prefix-icon="el-icon-search"
              @change="handleSearch"
              clearable
            />
            <el-select
              v-model="filterAssignmentType"
              placeholder="分配类型"
              size="small"
              style="width: 120px"
              @change="handleFilter"
              clearable
            >
              <el-option label="主要" value="primary" />
              <el-option label="额外" value="additional" />
              <el-option label="临时" value="temporary" />
              <el-option label="项目" value="project" />
            </el-select>
            <el-select
              v-model="filterStatus"
              placeholder="状态"
              size="small"
              style="width: 100px"
              @change="handleFilter"
              clearable
            >
              <el-option label="激活" value="active" />
              <el-option label="停用" value="inactive" />
            </el-select>
          </div>
          <div class="filter-right">
            <el-button size="small" @click="exportAssignments">
              <i class="el-icon-download"></i>
              导出
            </el-button>
            <el-button size="small" type="primary" @click="showAssignDialog = true">
              <i class="el-icon-plus"></i>
              新增分配
            </el-button>
          </div>
        </div>
      </el-card>

      <!-- 分配列表 -->
      <el-card class="assignments-card">
        <div slot="header" class="card-header">
          <span>分配列表 ({{ filteredAssignments.length }})</span>
          <div class="header-actions">
            <el-button size="mini" type="danger" @click="handleBatchRemove">批量移除</el-button>
          </div>
        </div>

        <el-table
          :data="paginatedAssignments"
          v-loading="loading"
          size="small"
          @selection-change="handleSelectionChange"
          @sort-change="handleSortChange"
        >
          <el-table-column type="selection" width="55" />
          
          <el-table-column
            prop="employee_name"
            label="员工姓名"
            width="120"
            sortable="custom"
          >
            <template slot-scope="scope">
              <div class="employee-info">
                <el-avatar
                  :size="28"
                  :src="scope.row.employee_avatar"
                  :alt="scope.row.employee_name"
                >
                  {{ scope.row.employee_name?.charAt(0) }}
                </el-avatar>
                <span class="employee-name">{{ scope.row.employee_name }}</span>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column
            prop="employee_id"
            label="工号"
            width="100"
            sortable="custom"
          />
          
          <el-table-column
            prop="position_name"
            label="职位"
            width="150"
            show-overflow-tooltip
          />
          
          <el-table-column
            prop="assignment_type"
            label="分配类型"
            width="100"
            align="center"
          >
            <template slot-scope="scope">
              <el-tag
                :type="getAssignmentTypeColor(scope.row.assignment_type)"
                size="mini"
              >
                {{ getAssignmentTypeName(scope.row.assignment_type) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column
            prop="management_type"
            label="管理类型"
            width="100"
            align="center"
          >
            <template slot-scope="scope">
              <el-tag
                :type="getManagementTypeColor(scope.row.management_type)"
                size="mini"
              >
                {{ getManagementTypeName(scope.row.management_type) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column
            prop="work_percentage"
            label="工作占比"
            width="100"
            align="center"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.work_percentage }}%</span>
            </template>
          </el-table-column>
          
          <el-table-column
            prop="direct_manager_name"
            label="直接上级"
            width="120"
            show-overflow-tooltip
          />
          
          <el-table-column
            prop="effective_date"
            label="生效日期"
            width="120"
            sortable="custom"
          >
            <template slot-scope="scope">
              {{ formatDate(scope.row.effective_date) }}
            </template>
          </el-table-column>
          
          <el-table-column
            prop="expiration_date"
            label="失效日期"
            width="120"
          >
            <template slot-scope="scope">
              {{ formatDate(scope.row.expiration_date) || '无限期' }}
            </template>
          </el-table-column>
          
          <el-table-column
            prop="status"
            label="状态"
            width="80"
            align="center"
          >
            <template slot-scope="scope">
              <el-tag
                :type="scope.row.status === 'active' ? 'success' : 'danger'"
                size="mini"
              >
                {{ scope.row.status === 'active' ? '激活' : '停用' }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="150" fixed="right">
            <template slot-scope="scope">
              <el-button
                size="mini"
                type="text"
                @click="handleEditAssignment(scope.row)"
              >
                编辑
              </el-button>
              <el-button
                size="mini"
                type="text"
                @click="handleEditAssignment(scope.row)"
              >
                编辑
              </el-button>
              <el-button
                size="mini"
                type="text"
                @click="handleRemoveAssignment(scope.row)"
              >
                移除
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-container">
          <el-pagination
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page="pagination.current"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="pagination.size"
            layout="total, sizes, prev, pager, next, jumper"
            :total="filteredAssignments.length"
          />
        </div>
      </el-card>
    </div>

    <!-- 员工分配对话框 -->
    <EmployeeAssignmentDialog
      :visible.sync="showAssignDialog"
      :organization-unit="selectedUnit"
      :assignment="editingAssignment"
      @success="handleAssignmentSuccess"
    />


  </div>
</template>

<script>
import dayjs from 'dayjs'
import EmployeeAssignmentDialog from './EmployeeAssignmentDialog.vue'

export default {
  name: 'EmployeeAssignmentView',
  components: {
    EmployeeAssignmentDialog,
  },
  props: {
    selectedUnit: {
      type: Object,
      default: null
    },
    assignments: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      loading: false,
      searchKeyword: '',
      filterAssignmentType: '',
      filterStatus: '',
      selectedAssignments: [],
      editingAssignment: null,
      
      // 对话框状态
      showAssignDialog: false,
      
      // 分页
      pagination: {
        current: 1,
        size: 20
      },
      
      // 排序
      sortField: '',
      sortOrder: '',
      
      // 统计数据
      assignmentStats: {
        total: 0,
        primary: 0,
        additional: 0,
        temporary: 0
      },
      
      // 类型配置
      assignmentTypes: [
        { value: 'primary', label: '主要', color: 'success' },
        { value: 'additional', label: '额外', color: 'warning' },
        { value: 'temporary', label: '临时', color: 'info' },
        { value: 'project', label: '项目', color: 'primary' }
      ],
      
      managementTypes: [
        { value: 'line', label: '直线', color: 'success' },
        { value: 'matrix', label: '矩阵', color: 'warning' },
        { value: 'functional', label: '功能', color: 'info' }
      ]
    }
  },
  computed: {
    filteredAssignments() {
      let result = [...this.assignments]
      
      // 搜索过滤
      if (this.searchKeyword) {
        const keyword = this.searchKeyword.toLowerCase()
        result = result.filter(item =>
          item.employee_name?.toLowerCase().includes(keyword) ||
          item.employee_id?.toLowerCase().includes(keyword)
        )
      }
      
      // 分配类型过滤
      if (this.filterAssignmentType) {
        result = result.filter(item => item.assignment_type === this.filterAssignmentType)
      }
      
      // 状态过滤
      if (this.filterStatus) {
        result = result.filter(item => item.status === this.filterStatus)
      }
      
      // 排序
      if (this.sortField) {
        result.sort((a, b) => {
          const aVal = a[this.sortField]
          const bVal = b[this.sortField]
          
          if (this.sortOrder === 'ascending') {
            return aVal > bVal ? 1 : -1
          } else {
            return aVal < bVal ? 1 : -1
          }
        })
      }
      
      return result
    },
    
    paginatedAssignments() {
      const start = (this.pagination.current - 1) * this.pagination.size
      const end = start + this.pagination.size
      return this.filteredAssignments.slice(start, end)
    }
  },
  watch: {
    selectedUnit: {
      handler(newVal) {
        if (newVal) {
          this.loadAssignments()
        }
      },
      immediate: true
    },
    
    assignments: {
      handler() {
        this.updateStats()
      },
      immediate: true
    }
  },
  methods: {
    // 加载分配数据
    async loadAssignments() {
      if (!this.selectedUnit) return
      
      this.loading = true
      try {
        // TODO: 调用API加载分配数据
        this.$emit('load-assignments', this.selectedUnit.id)
      } catch (error) {
        console.error('Failed to load assignments:', error)
      } finally {
        this.loading = false
      }
    },
    
    // 更新统计数据
    updateStats() {
      this.assignmentStats = {
        total: this.assignments.length,
        primary: this.assignments.filter(a => a.assignment_type === 'primary').length,
        additional: this.assignments.filter(a => a.assignment_type === 'additional').length,
        temporary: this.assignments.filter(a => a.assignment_type === 'temporary').length
      }
    },
    
    // 刷新数据
    refreshData() {
      this.loadAssignments()
    },
    
    // 搜索处理
    handleSearch() {
      this.pagination.current = 1
    },
    
    // 过滤处理
    handleFilter() {
      this.pagination.current = 1
    },
    
    // 选择变更
    handleSelectionChange(selection) {
      this.selectedAssignments = selection
    },
    
    // 排序变更
    handleSortChange({ column, prop, order }) {
      this.sortField = prop
      this.sortOrder = order
    },
    
    // 分页处理
    handleSizeChange(val) {
      this.pagination.size = val
      this.pagination.current = 1
    },
    
    handleCurrentChange(val) {
      this.pagination.current = val
    },
    
    // 编辑分配
    handleEditAssignment(assignment) {
      this.editingAssignment = assignment
      this.showAssignDialog = true
    },
    
    
    // 移除分配
    async handleRemoveAssignment(assignment) {
      try {
        await this.$confirm(
          `确定要移除员工"${assignment.employee_name}"的分配吗？`,
          '确认移除',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        // TODO: 调用API移除分配
        this.$emit('remove-assignment', assignment)
        this.$message.success('移除成功')
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('移除失败')
        }
      }
    },
    
    
    // 批量移除
    async handleBatchRemove() {
      if (this.selectedAssignments.length === 0) {
        this.$message.warning('请选择要移除的分配')
        return
      }
      
      try {
        await this.$confirm(
          `确定要移除选中的 ${this.selectedAssignments.length} 个分配吗？`,
          '确认移除',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        // TODO: 调用API批量移除
        this.$emit('batch-remove-assignments', this.selectedAssignments)
        this.$message.success('批量移除成功')
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('批量移除失败')
        }
      }
    },
    
    // 导出分配
    exportAssignments() {
      // TODO: 实现导出功能
      this.$message.info('导出功能开发中')
    },
    
    // 分配成功
    handleAssignmentSuccess() {
      this.showAssignDialog = false
      this.editingAssignment = null
      this.loadAssignments()
    },
    
    
    
    // 获取分配类型颜色
    getAssignmentTypeColor(type) {
      const assignmentType = this.assignmentTypes.find(t => t.value === type)
      return assignmentType ? assignmentType.color : ''
    },
    
    // 获取分配类型名称
    getAssignmentTypeName(type) {
      const assignmentType = this.assignmentTypes.find(t => t.value === type)
      return assignmentType ? assignmentType.label : type
    },
    
    // 获取管理类型颜色
    getManagementTypeColor(type) {
      const managementType = this.managementTypes.find(t => t.value === type)
      return managementType ? managementType.color : ''
    },
    
    // 获取管理类型名称
    getManagementTypeName(type) {
      const managementType = this.managementTypes.find(t => t.value === type)
      return managementType ? managementType.label : type
    },
    
    // 格式化日期
    formatDate(date) {
      return date ? dayjs(date).format('YYYY-MM-DD') : ''
    }
  }
}
</script>

<style scoped>
.assignment-view {
  height: 100%;
  overflow: auto;
}

.view-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e4e7ed;
}

.view-header h3 {
  margin: 0;
  color: #303133;
  font-size: 18px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 400px;
  color: #909399;
}

.empty-content {
  text-align: center;
}

.empty-content i {
  font-size: 48px;
  margin-bottom: 16px;
  color: #c0c4cc;
}

.empty-content p {
  margin: 0;
  font-size: 14px;
}

.assignment-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.stats-cards {
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;
}

.stat-item {
  padding: 16px 0;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #409eff;
  margin-bottom: 8px;
}

.stat-value.primary {
  color: #67c23a;
}

.stat-value.additional {
  color: #e6a23c;
}

.stat-value.temporary {
  color: #909399;
}

.stat-label {
  font-size: 12px;
  color: #909399;
}

.filter-card {
  margin-bottom: 20px;
}

.filter-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-left {
  display: flex;
  gap: 12px;
  align-items: center;
}

.filter-right {
  display: flex;
  gap: 8px;
}

.assignments-card {
  flex: 1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.employee-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.employee-name {
  font-weight: 500;
  color: #303133;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}
</style>