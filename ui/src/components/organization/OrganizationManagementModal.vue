<template>
  <el-dialog
    :visible.sync="dialogVisible"
    title="组织管理"
    width="95%"
    height="90%"
    class="organization-modal"
    :before-close="handleClose"
    :close-on-click-modal="false"
  >
    <!-- 头部工具栏 -->
    <div class="modal-header">
      <div class="header-actions">
        <div class="header-title">
          <h3>组织架构图</h3>
        </div>
        
        <div class="header-controls">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索组织单元"
            size="small"
            style="width: 250px"
            prefix-icon="el-icon-search"
            @input="handleSearch"
            clearable
          />
          <el-button 
            type="primary" 
            size="small" 
            @click="showCreateUnitDialog = true"
            icon="el-icon-plus"
          >
            新建组织单元
          </el-button>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="modal-content">
      <!-- 组织架构图 -->
      <div class="chart-panel">
        <div class="panel-header">
          <div class="org-stats" v-if="organizationStats">
            <span class="stat-item">
              <i class="el-icon-office-building"></i>
              {{ organizationStats.totalDepartments }}个单元
            </span>
            <span class="stat-item">
              <i class="el-icon-user"></i>
              {{ organizationStats.totalEmployees }}名员工
            </span>
            <span class="stat-item">
              <i class="el-icon-s-grid"></i>
              {{ organizationStats.maxLevel }}级层次
            </span>
          </div>
        </div>
        
        <div class="organization-chart-container">
          <WorkdayOrgChart 
            :organization-data="organizationTree"
            :selected-unit="selectedUnit"
            @unit-selected="handleChartNodeSelected"
            @create-department="handleCreateDepartment"
            @create-employee="handleCreateEmployee"
            class="main-org-chart"
          />
        </div>
      </div>
    </div>

    <!-- 创建组织单元对话框 -->
    <OrganizationUnitDialog
      :visible.sync="showCreateUnitDialog"
      :parent-unit="selectedUnit"
      :unit-types="unitTypes"
      @success="handleCreateSuccess"
    />

    <!-- 编辑组织单元对话框 -->
    <OrganizationUnitDialog
      :visible.sync="showEditUnitDialog"
      :unit="editingUnit"
      :unit-types="unitTypes"
      :is-edit="true"
      @success="handleEditSuccess"
    />

    <!-- 员工分配对话框 -->
    <EmployeeAssignmentDialog
      :visible.sync="showAssignmentDialog"
      :organization-unit="selectedUnit"
      :assignment="editingAssignment"
      @success="handleAssignmentSuccess"
    />

  </el-dialog>
</template>

<script>
import OrganizationUnitDialog from './OrganizationUnitDialog.vue'
import EmployeeAssignmentDialog from './EmployeeAssignmentDialog.vue'
import WorkdayOrgChart from './WorkdayOrgChart.vue'
import { organizationApiService } from '@/services/organizationApi'
import { departmentApi } from '@/services/departmentApi'
import { employeeApiService } from '@/services/employeeApi'

export default {
  name: 'OrganizationManagementModal',
  components: {
    OrganizationUnitDialog,
    EmployeeAssignmentDialog,
    WorkdayOrgChart
  },
  props: {
    visible: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      // 基础状态
      searchKeyword: '',
      selectedUnit: null,
      
      // 数据
      organizationTree: [],
      allEmployees: [], // 存储所有员工数据
      organizationStats: null, // 组织统计信息
      
      // 对话框状态
      showCreateUnitDialog: false,
      showEditUnitDialog: false,
      showAssignmentDialog: false,
      
      // 编辑状态
      editingUnit: null,
      editingAssignment: null,
      
      
      // 组织单元类型
      unitTypes: [
        { value: 'company', label: '公司/法人实体', icon: 'el-icon-office-building', color: '' },
        { value: 'business_unit', label: '业务单元', icon: 'el-icon-s-cooperation', color: 'primary' },
        { value: 'department', label: '部门', icon: 'el-icon-user', color: 'success' },
        { value: 'team', label: '团队', icon: 'el-icon-s-custom', color: 'warning' },
        { value: 'cost_center', label: '成本中心', icon: 'el-icon-s-finance', color: 'info' },
        { value: 'location', label: '地理位置', icon: 'el-icon-location', color: 'danger' },
        { value: 'project', label: '项目组', icon: 'el-icon-s-flag', color: 'warning' }
      ]
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
    visible(val) {
      if (val) {
        this.initData()
      }
    },
    selectedUnit() {
      // 选中单元变化时可以执行相关操作
      if (this.selectedUnit) {
        console.log('选中组织单元:', this.selectedUnit.name)
      }
    }
  },
  methods: {
    // 初始化数据
    async initData() {
      try {
        // 只在初始化时加载组织树
        await this.loadOrganizationTree()
        // 其他数据将在选中单元时加载
      } catch (error) {
        console.error('加载组织数据失败:', error)
        this.$message.error('加载组织数据失败')
      }
    },
    
    // 加载组织架构树
    async loadOrganizationTree() {
      try {
        const response = await departmentApi.getDepartmentTree()
        const responseData = response.data || response
        const treeData = responseData.data || []
        
        console.log('原始部门树数据:', treeData)
        
        // 转换数据格式并计算员工和子单元数量
        this.organizationTree = await this.processTreeData(treeData)
        
        console.log('处理后的部门树数据:', this.organizationTree)
        
        // 计算组织统计信息
        this.calculateOrganizationStats()
      } catch (error) {
        console.error('获取组织架构树失败:', error)
        this.$message.error('获取组织架构失败')
        this.organizationTree = []
      }
    },
    
    // 处理树形数据
    async processTreeData(treeData) {
      // 先获取所有员工数据用于计算
      await this.loadAllEmployeesData()
      
      const processNode = (node, parentPath = '') => {
        // 构建层级路径
        const hierarchyPath = parentPath ? `${parentPath} > ${node.name}` : node.name
        
        // 处理子节点
        const children = []
        if (node.children && node.children.length > 0) {
          for (const child of node.children) {
            children.push(processNode(child, hierarchyPath))
          }
        }
        
        // 计算真实的直接员工数量
        const directEmployeeCount = this.getDirectEmployeeCount(node.id)
        
        // 计算子单元数量（真实的子部门数量）
        const subunitCount = children.length
        
        // 计算层级员工数量（包括所有子部门的员工）
        const hierarchicalEmployeeCount = this.calculateHierarchicalEmployeeCount(node, children)
        
        return {
          id: node.id,
          name: node.name || node.short_name || `部门${node.id}`,
          code: node.code,
          short_name: node.short_name,
          type: node.type || 'department',
          level: node.level || 1,
          description: node.description,
          manager_id: node.manager_id,
          manager: node.manager,
          employeeCount: directEmployeeCount,
          hierarchicalEmployeeCount: hierarchicalEmployeeCount,
          subunitCount: subunitCount,
          hierarchyPath: hierarchyPath,
          children: children,
          // 保留原始数据
          originalData: node
        }
      }
      
      const result = []
      for (const node of treeData) {
        result.push(processNode(node))
      }
      
      return result
    },
    
    // 计算层级员工数量（包括所有子部门的员工）
    calculateHierarchicalEmployeeCount(node, processedChildren = []) {
      // 获取当前部门的直接员工数量
      let totalCount = this.getDirectEmployeeCount(node.id)
      
      // 如果有已处理的子节点，使用其数据
      if (processedChildren && processedChildren.length > 0) {
        for (const child of processedChildren) {
          totalCount += child.hierarchicalEmployeeCount || 0
        }
      } else if (node.children && node.children.length > 0) {
        // 如果是原始数据，递归计算
        for (const child of node.children) {
          totalCount += this.calculateHierarchicalEmployeeCount(child)
        }
      }
      
      return totalCount
    },
    
    // 加载所有员工数据
    async loadAllEmployeesData() {
      try {
        const response = await employeeApiService.getEmployees({
          page: 1,
          pageSize: 10000, // 获取所有员工
          status: 'active'
        })
        const responseData = response.data || response
        this.allEmployees = responseData.data || responseData.employees || []
        
        console.log('已加载员工数据:', this.allEmployees.length, '条记录')
      } catch (error) {
        console.error('获取员工数据失败:', error)
        this.allEmployees = []
      }
    },
    
    // 获取部门的直接员工数量
    getDirectEmployeeCount(departmentId) {
      if (!this.allEmployees || !departmentId) return 0
      
      return this.allEmployees.filter(emp => {
        const empDeptId = emp.department_id || (emp.department && emp.department.id)
        return empDeptId === departmentId
      }).length
    },
    
    // 获取部门的直接员工列表
    getDirectEmployees(departmentId) {
      if (!this.allEmployees || !departmentId) return []
      
      return this.allEmployees.filter(emp => {
        const empDeptId = emp.department_id || (emp.department && emp.department.id)
        return empDeptId === departmentId
      })
    },
    
    // 获取部门完整的员工列表（包含部门负责人和所有下属，含子部门员工）
    getDepartmentCompleteEmployeeList(departmentId) {
      if (!this.allEmployees || !departmentId) return []
      
      const result = []
      const departmentNode = this.findDepartmentNode(departmentId)
      
      if (!departmentNode) return []
      
      // 收集所有相关员工
      const collectEmployees = (node) => {
        // 获取当前部门的所有员工（包括部门负责人）
        const deptEmployees = this.allEmployees.filter(emp => {
          const empDeptId = emp.department_id || (emp.department && emp.department.id)
          return empDeptId === node.id
        })
        
        // 添加部门负责人标识
        deptEmployees.forEach(emp => {
          if (emp.id === node.manager_id) {
            emp.isDepartmentManager = true
            emp.managerDepartmentName = node.name
          }
        })
        
        result.push(...deptEmployees)
        
        // 递归处理子部门
        if (node.children && node.children.length > 0) {
          node.children.forEach(child => {
            collectEmployees(child)
          })
        }
      }
      
      collectEmployees(departmentNode)
      
      // 按角色排序：部门负责人优先，然后按员工ID排序
      return result.sort((a, b) => {
        if (a.isDepartmentManager && !b.isDepartmentManager) return -1
        if (!a.isDepartmentManager && b.isDepartmentManager) return 1
        return a.id - b.id
      })
    },
    
    // 查找部门节点（递归搜索）
    findDepartmentNode(departmentId) {
      const findInTree = (nodes) => {
        for (const node of nodes) {
          if (node.id === departmentId) {
            return node
          }
          if (node.children && node.children.length > 0) {
            const found = findInTree(node.children)
            if (found) return found
          }
        }
        return null
      }
      
      return findInTree(this.organizationTree)
    },
    
    // 计算组织统计信息
    calculateOrganizationStats() {
      if (!this.organizationTree || this.organizationTree.length === 0) {
        this.organizationStats = null
        return
      }
      
      let totalDepartments = 0
      let totalEmployees = 0
      let maxLevel = 0
      
      const countNodes = (nodes) => {
        for (const node of nodes) {
          totalDepartments++
          totalEmployees += node.employeeCount || 0
          maxLevel = Math.max(maxLevel, node.level || 1)
          
          if (node.children && node.children.length > 0) {
            countNodes(node.children)
          }
        }
      }
      
      countNodes(this.organizationTree)
      
      this.organizationStats = {
        totalDepartments,
        totalEmployees,
        maxLevel
      }
    },
    
    
    // 搜索处理
    handleSearch() {
      // 实现搜索逻辑
    },
    
    
    // 图表节点选择处理
    handleChartNodeSelected(nodeData) {
      this.selectedUnit = nodeData
      console.log('选中组织单元:', nodeData.name)
    },
    
    // 处理创建子部门
    handleCreateDepartment(parentUnit) {
      this.selectedUnit = parentUnit
      this.showCreateUnitDialog = true
      console.log('创建子部门，父级单元:', parentUnit.name)
    },
    
    // 处理创建员工
    handleCreateEmployee(departmentUnit) {
      this.selectedUnit = departmentUnit
      this.showAssignmentDialog = true
      console.log('创建员工，所属部门:', departmentUnit.name)
    },
    
    
    // 节点操作
    handleNodeAction(command) {
      const { action, data } = command
      
      switch (action) {
        case 'view':
          this.handleViewUnit(data)
          break
        case 'edit':
          this.handleEditUnit(data)
          break
        case 'addChild':
          this.handleAddChild(data)
          break
        case 'assign':
          this.handleAssignEmployee(data)
          break
        case 'move':
          this.handleMoveUnit(data)
          break
        case 'delete':
          this.handleDeleteUnit(data)
          break
      }
    },
    
    // 查看单元详情
    handleViewUnit(unit) {
      this.selectedUnit = unit
      this.currentView = 'hierarchy'
    },
    
    // 编辑单元
    handleEditUnit(unit) {
      this.editingUnit = unit
      this.showEditUnitDialog = true
    },
    
    // 添加子单元
    handleAddChild(unit) {
      this.selectedUnit = unit
      this.showCreateUnitDialog = true
    },
    
    // 分配员工
    handleAssignEmployee(unit) {
      this.selectedUnit = unit
      this.editingAssignment = null
      this.showAssignmentDialog = true
    },
    
    // 移动单元
    handleMoveUnit(unit) {
      this.$message.info('移动功能开发中')
    },
    
    // 删除单元
    async handleDeleteUnit(unit) {
      try {
        await this.$confirm(`确定要删除组织单元"${unit.name}"吗？`, '确认删除', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        // 调用API删除
        await departmentApi.deleteDepartment(unit.id)
        this.$message.success('删除成功')
        this.loadOrganizationTree()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除失败:', error)
          this.$message.error('删除失败: ' + (error.response?.data?.message || error.message || '未知错误'))
        }
      }
    },
    
    
    // 获取单元类型图标
    getUnitTypeIcon(type) {
      const unitType = this.unitTypes.find(t => t.value === type)
      return unitType ? unitType.icon : 'el-icon-office-building'
    },
    
    
    // 获取单元类型颜色
    getUnitTypeColor(type) {
      const unitType = this.unitTypes.find(t => t.value === type)
      return unitType ? unitType.color : ''
    },
    
    // 获取单元类型名称
    getUnitTypeName(type) {
      const unitType = this.unitTypes.find(t => t.value === type)
      return unitType ? unitType.label : type
    },
    
    // 单元选择
    handleUnitSelected(unit) {
      this.selectedUnit = unit
    },
    
    
    // 刷新数据
    handleRefresh() {
      this.initData()
    },
    
    // 创建成功
    handleCreateSuccess() {
      this.showCreateUnitDialog = false
      this.loadOrganizationTree()
      this.$message.success('创建成功')
    },
    
    // 编辑成功
    handleEditSuccess() {
      this.showEditUnitDialog = false
      this.editingUnit = null
      this.loadOrganizationTree()
      this.$message.success('编辑成功')
    },
    
    // 分配成功
    handleAssignmentSuccess() {
      this.showAssignmentDialog = false
      this.editingAssignment = null
      this.loadOrganizationTree()
      this.$message.success('分配成功')
    },
    
    
    // 关闭对话框
    handleClose() {
      this.dialogVisible = false
      this.selectedUnit = null
    }
  }
}
</script>

<style scoped>
.organization-modal {
  min-height: 600px;
}

.organization-modal :deep(.el-dialog) {
  margin-top: 5vh !important;
  margin-bottom: 5vh !important;
  height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.organization-modal :deep(.el-dialog__body) {
  padding: 0;
  height: calc(100% - 120px);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.modal-header {
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
  background: #f8f9fa;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-controls {
  display: flex;
  gap: 12px;
  align-items: center;
}

.modal-content {
  display: flex;
  height: 100%;
  overflow: hidden;
}

.left-panel {
  width: 400px; /* 增加左侧面板宽度 */
  min-width: 350px; /* 设置最小宽度 */
  border-right: 1px solid #e4e7ed;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  resize: horizontal; /* 允许水平拖拽调整宽度 */
}

.panel-header {
  padding: 16px;
  border-bottom: 1px solid #e4e7ed;
  background: #fafafa;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.header-title {
  flex: 1;
}

.panel-header h4 {
  margin: 0 0 8px 0;
  color: #303133;
  font-size: 16px;
  font-weight: 600;
}

.org-stats {
  display: flex;
  gap: 16px;
  flex-wrap: nowrap;
  align-items: center;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #606266;
  background: #f5f7fa;
  padding: 4px 8px;
  border-radius: 12px;
}

.stat-item i {
  color: #409eff;
}

.tree-controls {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
  align-self: flex-start;
}

.organization-chart-container {
  flex: 1;
  overflow: hidden;
  background: white;
  border-radius: 6px;
  min-height: 0; /* 允许flex子项收缩 */
}

.left-org-chart {
  height: 100%;
  background: white;
}


.right-panel {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.hierarchy-view,
.chart-view,
.assignments-view, 
.changes-view {
  flex: 1;
  overflow: auto;
  padding: 20px;
  height: 100%;
  min-height: 0; /* 允许flex子项收缩 */
}


/* 响应式适配 */
@media (max-width: 1200px) {
  .left-panel {
    width: 350px; /* 在中等屏幕保持合适宽度 */
    min-width: 300px;
  }
  
  .node-name {
    max-width: 150px; /* 中等屏幕下减小最大宽度 */
  }
}

@media (max-width: 768px) {
  .modal-content {
    flex-direction: column;
  }
  
  .left-panel {
    width: 100%;
    max-height: 300px;
    resize: none; /* 小屏幕下禁用拖拽调整 */
  }
  
  .tree-node {
    min-height: 50px; /* 小屏幕下减小最小高度 */
    padding: 8px 10px;
  }
  
  .node-name {
    max-width: 120px; /* 小屏幕下进一步减小宽度 */
    font-size: 13px;
  }
  
  .node-stats {
    min-width: 60px;
  }
  
  .header-actions {
    flex-direction: column;
    gap: 12px;
  }
  
  .header-controls {
    width: 100%;
    justify-content: space-between;
  }
}

/* 添加tooltip效果 */
.node-name:hover {
  position: relative;
  z-index: 1000;
}

/* 为长文本添加更好的视觉反馈 */
.tree-node:hover .node-name {
  color: #409eff;
  transition: color 0.3s;
}

</style>