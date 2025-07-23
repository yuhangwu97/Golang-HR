<template>
  <div class="department-list">
    <div class="page-header">
      <h2>部门管理</h2>
      <div class="actions">
        <el-button type="primary" @click="showCreateDialog = true">
          <i class="el-icon-plus"></i>
          新增部门
        </el-button>
        <el-button type="success" @click="showOrganizationModal = true">
          <i class="el-icon-s-cooperation"></i>
          组织管理
        </el-button>
        <el-button @click="showTreeView = !showTreeView">
          <i class="el-icon-menu"></i>
          {{ showTreeView ? '列表视图' : '树形视图' }}
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
              <div class="stat-label">总部门数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ (statistics && statistics.with_manager) || 0 }}</div>
              <div class="stat-label">有负责人部门</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ Object.keys((statistics && statistics.by_level) || {}).length }}</div>
              <div class="stat-label">层级数量</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ totalEmployees }}</div>
              <div class="stat-label">总员工数</div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 搜索过滤 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="部门名称">
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入部门名称或编码"
            clearable
            @keyup.enter.native="handleSearch"
          />
        </el-form-item>
        <el-form-item label="上级部门">
          <el-cascader
            v-model="searchForm.parent_id"
            :options="departmentTree"
            :props="cascaderProps"
            placeholder="选择上级部门"
            clearable
            style="width: 200px"
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

    <!-- 内容区域 -->
    <el-card>
      <!-- 树形视图 -->
      <div v-if="showTreeView">
        <el-tree
          ref="treeRef"
          v-loading="treeLoading"
          :data="departmentTree"
          :props="treeProps"
          node-key="id"
          show-checkbox
          default-expand-all
          :expand-on-click-node="false"
          @node-click="handleNodeClick"
        >
          <span slot-scope="{ node, data }" class="tree-node">
            <div class="node-content">
              <span class="node-label">{{ data.name }}</span>
              <el-tag v-if="data.code" size="small" type="info">{{ data.code }}</el-tag>
              <el-tag v-if="data.manager" size="small" type="success">
                负责人: {{ data.manager.name }}
              </el-tag>
              <el-tag :type="data.status === 'active' ? 'success' : 'danger'" size="small">
                {{ data.status === 'active' ? '活跃' : '非活跃' }}
              </el-tag>
            </div>
            <div class="node-actions">
              <el-button size="small" @click.stop="handleEdit(data)">编辑</el-button>
              <el-button size="small" type="primary" @click.stop="handleAddChild(data)">
                添加子部门
              </el-button>
              <el-button size="small" type="danger" @click.stop="handleDelete(data)">
                删除
              </el-button>
            </div>
          </span>
        </el-tree>
      </div>

      <!-- 列表视图 -->
      <div v-else>
        <el-table
          v-loading="loading"
          :data="departments"
          stripe
          row-key="id"
          :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="name" label="部门名称" width="200" />
          <el-table-column prop="code" label="部门编码" width="150" />
          <el-table-column prop="parent" label="上级部门" width="150">
            <template slot-scope="{ row }">
              <span v-if="row.parent">{{ row.parent.name }}</span>
              <span v-else class="text-muted">顶级部门</span>
            </template>
          </el-table-column>
          <el-table-column prop="manager" label="部门负责人" width="120">
            <template slot-scope="{ row }">
              <el-tag v-if="row.manager" type="success" size="small">
                {{ row.manager.name }}
              </el-tag>
              <span v-else class="text-muted">未设置</span>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="部门描述" show-overflow-tooltip />
          <el-table-column prop="sort" label="排序" width="80" />
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
              <el-button size="small" type="primary" @click="handleAddChild(row)">
                添加子部门
              </el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <DepartmentDialog
      :visible.sync="showCreateDialog"
      :department="editDepartment"
      :parent-department="parentDepartment"
      @success="handleDialogSuccess"
    />

    <!-- 组织管理模态框 -->
    <OrganizationManagementModal
      :visible.sync="showOrganizationModal"
      @refresh="handleOrganizationRefresh"
    />
  </div>
</template>

<script>
import { departmentApi } from '@/services/departmentApi'
import DepartmentDialog from '@/components/organization/DepartmentDialog.vue'
import OrganizationManagementModal from '@/components/organization/OrganizationManagementModal.vue'
import { formatDate } from '@/utils/date'

export default {
  name: 'DepartmentListView',
  components: {
    DepartmentDialog,
    OrganizationManagementModal
  },
  data() {
    return {
      loading: false,
      treeLoading: false,
      showStatistics: false,
      showTreeView: true,
      showCreateDialog: false,
      showOrganizationModal: false,
      departments: [],
      departmentTree: [],
      statistics: null,
      editDepartment: null,
      parentDepartment: null,
      selectedDepartments: [],
      searchForm: {
        keyword: '',
        parent_id: undefined,
        status: ''
      },
      // 级联选择器配置
      cascaderProps: {
        value: 'id',
        label: 'name',
        children: 'children',
        emitPath: false,
        checkStrictly: true
      },
      // 树形控件配置
      treeProps: {
        children: 'children',
        label: 'name'
      }
    }
  },
  computed: {
    totalEmployees() {
      if (!this.statistics || !this.statistics.employee_counts) {
        return 0
      }
      return Object.values(this.statistics.employee_counts).reduce((sum, count) => sum + count, 0)
    }
  },
  mounted() {
    this.fetchDepartments()
    this.fetchDepartmentTree()
    this.fetchStatistics()
  },
  methods: {
    formatDate,
    async fetchDepartments() {
      this.loading = true
      try {
        const response = await departmentApi.getDepartments()
        // 处理API响应数据结构
        const responseData = response.data || response
        this.departments = responseData.data || responseData || []
      } catch (error) {
        console.error('获取部门列表失败:', error)
        this.$message.error('获取部门列表失败')
      } finally {
        this.loading = false
      }
    },
    async fetchDepartmentTree() {
      this.treeLoading = true
      try {
        const response = await departmentApi.getDepartmentTree()
        // 处理API响应数据结构
        const responseData = response.data || response
        this.departmentTree = responseData.data || responseData || []
      } catch (error) {
        console.error('获取部门树失败:', error)
        // 如果获取部门树失败，尝试使用普通部门列表构建树形结构
        try {
          const flatResponse = await departmentApi.getDepartments()
          const flatData = flatResponse.data || flatResponse
          const flatDepartments = flatData.data || flatData || []
          this.departmentTree = this.buildTreeFromFlat(flatDepartments)
        } catch (flatError) {
          console.error('获取部门列表也失败:', flatError)
          this.$message.error('获取部门数据失败')
          this.departmentTree = []
        }
      } finally {
        this.treeLoading = false
      }
    },
    async fetchStatistics() {
      try {
        const response = await departmentApi.getDepartmentStatistics()
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
      // 实现搜索逻辑
      this.fetchDepartments()
    },
    handleReset() {
      Object.assign(this.searchForm, {
        keyword: '',
        parent_id: undefined,
        status: ''
      })
      this.fetchDepartments()
    },
    handleNodeClick(data) {
      // 树节点点击事件
      console.log('Node clicked:', data)
    },
    handleEdit(department) {
      this.editDepartment = department
      this.parentDepartment = null
      this.showCreateDialog = true
    },
    handleAddChild(department) {
      this.editDepartment = null
      this.parentDepartment = department
      this.showCreateDialog = true
    },
    async handleDelete(department) {
      try {
        await this.$confirm(
          `确定要删除部门"${department.name}"吗？删除后其子部门也将被删除。`,
          '确认删除',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        await departmentApi.deleteDepartment(department.id)
        this.$message.success('删除成功')
        this.fetchDepartments()
        this.fetchDepartmentTree()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除失败')
        }
      }
    },
    handleSelectionChange(selected) {
      this.selectedDepartments = selected
    },
    handleDialogSuccess() {
      this.showCreateDialog = false
      this.editDepartment = null
      this.parentDepartment = null
      this.fetchDepartments()
      this.fetchDepartmentTree()
      if (this.showStatistics) {
        this.fetchStatistics()
      }
    },
    handleOrganizationRefresh() {
      this.fetchDepartments()
      this.fetchDepartmentTree()
      if (this.showStatistics) {
        this.fetchStatistics()
      }
      this.$message.success('组织架构数据已刷新')
    },
    // 从扁平化数据构建树形结构
    buildTreeFromFlat(flatDepartments) {
      if (!flatDepartments || flatDepartments.length === 0) {
        return []
      }

      const map = {}
      const roots = []

      // 创建映射
      flatDepartments.forEach(dept => {
        map[dept.id] = { ...dept, children: [] }
      })

      // 构建树形结构
      flatDepartments.forEach(dept => {
        if (dept.parent_id && map[dept.parent_id]) {
          map[dept.parent_id].children.push(map[dept.id])
        } else {
          roots.push(map[dept.id])
        }
      })

      return roots
    },
    // 计算基础统计信息
    calculateBasicStatistics() {
      const total = this.departments.length
      const withManager = this.departments.filter(dept => dept.manager).length
      const byLevel = {}
      
      this.departments.forEach(dept => {
        const level = dept.level || 1
        byLevel[level] = (byLevel[level] || 0) + 1
      })

      return {
        total,
        with_manager: withManager,
        by_level: byLevel,
        employee_counts: {}
      }
    }
  }
}
</script>

<style scoped>
.department-list {
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

.text-muted {
  color: #999;
  font-style: italic;
}

/* 树形视图样式 */
.tree-node {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  padding-right: 10px;
}

.node-content {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.node-label {
  font-weight: 500;
}

.node-actions {
  display: flex;
  gap: 5px;
}

:deep(.el-tree-node__content) {
  height: 40px;
}

:deep(.el-tree-node__expand-icon) {
  padding: 6px;
}
</style>