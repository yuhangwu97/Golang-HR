<template>
  <div class="organization-chart-view">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-title">
        <div class="title-icon">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <div class="title-content">
          <h1>组织架构</h1>
          <p>可视化组织结构管理</p>
        </div>
      </div>
      
      <div class="page-actions">
        <el-button 
          @click="refreshChart"
          icon="el-icon-refresh"
          class="action-btn"
        >
          刷新
        </el-button>
        <el-button 
          @click="exportChart"
          icon="el-icon-download"
          class="action-btn"
        >
          导出图表
        </el-button>
      </div>
    </div>


    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- 图表区域（占满全屏） -->
      <div class="chart-container">
        <div class="chart-wrapper">
          <WorkdayOrgChart 
            :organization-data="organizationTree"
            :selected-unit="selectedUnit"
            :view-type="viewType"
            @unit-selected="handleUnitSelected"
            @create-department="handleCreateDepartment"
            @create-employee="handleCreateEmployee"
            @view-type-change="handleViewTypeChange"
            @view-employee="handleViewEmployee"
            @edit-employee="handleEditEmployee"
            ref="orgChart"
            class="main-org-chart"
          />
        </div>
      </div>
      
      <!-- 悬浮的选中单元信息面板 -->
      <div v-if="selectedUnit" class="floating-unit-info">
        <div class="floating-header">
          <h3>选中单元信息</h3>
          <el-button 
            type="text" 
            icon="el-icon-close" 
            @click="selectedUnit = null"
            class="close-btn"
          />
        </div>
        <div class="selected-unit-info">
          <div class="unit-name">{{ selectedUnit.name }}</div>
          <div class="unit-details">
            <div class="detail-item">
              <span class="label">类型：</span>
              <span class="value">{{ getUnitTypeName(selectedUnit.type) }}</span>
            </div>
            <div class="detail-item" v-if="selectedUnit.type !== 'employee'">
              <span class="label">负责人：</span>
              <span class="value">{{ selectedUnit.manager ? selectedUnit.manager.name : '无' }}</span>
            </div>
            <div class="detail-item" v-if="selectedUnit.type !== 'employee'">
              <span class="label">直接员工：</span>
              <span class="value">{{ selectedUnit.employeeCount || 0 }}人</span>
            </div>
            <div class="detail-item" v-if="selectedUnit.type !== 'employee'">
              <span class="label">子单元：</span>
              <span class="value">{{ selectedUnit.subunitCount || 0 }}个</span>
            </div>
            <!-- 员工信息显示 -->
            <div v-if="selectedUnit.type === 'employee'" class="employee-details">
              <div class="detail-item">
                <span class="label">员工工号：</span>
                <span class="value">{{ selectedUnit.employeeData?.employee_id || '-' }}</span>
              </div>
              <div class="detail-item">
                <span class="label">职位：</span>
                <span class="value">{{ selectedUnit.employeeData?.position?.name || '未设置' }}</span>
              </div>
              <div class="detail-item">
                <span class="label">状态：</span>
                <span class="value">{{ selectedUnit.employeeData?.status === 'active' ? '在职' : '离职' }}</span>
              </div>
            </div>
          </div>
          
          <!-- 部门可用职位列表 -->
          <div v-if="selectedUnit.type !== 'employee' && departmentPositions.length > 0" class="department-positions">
            <h4>部门职位</h4>
            <div class="positions-list">
              <el-tag 
                v-for="position in departmentPositions" 
                :key="position.id" 
                size="mini" 
                class="position-tag"
              >
                {{ position.name }}
              </el-tag>
            </div>
          </div>
          
          <!-- 操作按钮 -->
          <div class="unit-actions">
            <!-- 部门操作按钮 -->
            <template v-if="selectedUnit.type !== 'employee'">
              <el-button size="mini" type="primary" @click="createChildUnit">
                <i class="el-icon-plus"></i> 新增子部门
              </el-button>
              <el-button size="mini" type="success" @click="createEmployee">
                <i class="el-icon-user"></i> 新增员工
              </el-button>
            </template>
            
            <!-- 员工操作按钮 -->
            <template v-else>
              <el-button size="mini" type="primary" @click="viewEmployee">
                <i class="el-icon-view"></i> 查看员工
              </el-button>
              <el-button size="mini" type="warning" @click="editEmployee">
                <i class="el-icon-edit"></i> 修改员工
              </el-button>
            </template>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建部门对话框 -->
    <DepartmentDialog
      :visible.sync="showCreateUnitDialog"
      :parent-department="selectedUnit"
      @success="handleCreateSuccess"
    />

    <!-- 编辑部门对话框 -->
    <DepartmentDialog
      :visible.sync="showEditUnitDialog"
      :department="editingUnit"
      @success="handleEditSuccess"
    />

  </div>
</template>

<script>
import WorkdayOrgChart from '@/components/organization/WorkdayOrgChart.vue'
import DepartmentDialog from '@/components/organization/DepartmentDialog.vue'
import { departmentApi } from '@/services/departmentApi'
import { employeeApiService } from '@/services/employeeApi'
import { positionApi } from '@/services/positionApi'

export default {
  name: 'OrganizationChartView',
  components: {
    WorkdayOrgChart,
    DepartmentDialog
  },
  data() {
    return {
      loading: false,
      viewType: 'tree',
      selectedUnit: null,
      
      // 数据
      organizationTree: [],
      allEmployees: [],
      departmentPositions: [], // 当前选中部门的职位列表
      selectedDepartmentInfo: null, // 选中部门信息，用于创建员工时预填充
      
      // 对话框状态
      showCreateUnitDialog: false,
      showEditUnitDialog: false,
      
      // 编辑状态
      editingUnit: null
    }
  },
  async mounted() {
    await this.initData()
  },
  methods: {
    // 初始化数据
    async initData() {
      this.loading = true
      try {
        await Promise.all([
          this.loadOrganizationTree(),
          this.loadAllEmployeesData()
        ])
      } catch (error) {
        console.error('加载组织数据失败:', error)
        this.$message.error('加载组织数据失败')
      } finally {
        this.loading = false
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
      } catch (error) {
        console.error('获取组织架构树失败:', error)
        this.$message.error('获取组织架构失败')
        this.organizationTree = []
      }
    },

    // 处理树形数据
    async processTreeData(treeData) {
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
        
        // 获取直接员工列表
        const directEmployees = this.getDirectEmployees(node.id)
        
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
          directEmployees: directEmployees, // 添加直接员工列表
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
        if (this.allEmployees.length > 0) {
          console.log('员工数据示例:', this.allEmployees.slice(0, 3))
        }
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
      
      const employees = this.allEmployees.filter(emp => {
        const empDeptId = emp.department_id || (emp.department && emp.department.id)
        return empDeptId === departmentId
      })
      
      if (employees.length > 0) {
        console.log(`部门 ${departmentId} 的员工:`, employees.map(emp => ({ id: emp.id, name: emp.name, employee_id: emp.employee_id })))
      }
      
      return employees
    },

    // 加载部门职位信息
    async loadDepartmentPositions(departmentId) {
      try {
        const response = await positionApi.getPositionsByDepartment(departmentId)
        this.departmentPositions = response.data || []
      } catch (error) {
        console.error('获取部门职位失败:', error)
        this.departmentPositions = []
      }
    },

    // 获取单元类型名称
    getUnitTypeName(type) {
      const typeMap = {
        'company': '公司/法人实体',
        'business_unit': '业务单元',
        'department': '部门',
        'team': '团队',
        'cost_center': '成本中心',
        'location': '地理位置',
        'project': '项目组'
      }
      return typeMap[type] || type
    },

    // 处理单元选择
    handleUnitSelected(unitData) {
      this.selectedUnit = unitData
      console.log('选中组织单元:', unitData.name)
      
      // 如果选中的是部门，加载部门职位信息
      if (unitData.type !== 'employee') {
        this.loadDepartmentPositions(unitData.id)
      } else {
        this.departmentPositions = []
      }
    },

    // 处理创建子部门
    handleCreateDepartment(parentUnit) {
      this.selectedUnit = parentUnit
      this.showCreateUnitDialog = true
    },

    // 处理创建员工
    handleCreateEmployee(departmentUnit) {
      this.selectedUnit = departmentUnit
      
      // 首先加载该部门的职位信息
      this.loadDepartmentPositions(departmentUnit.id).then(() => {
        // 构建员工创建页面需要的部门和管理者信息
        const departmentInfo = {
          id: departmentUnit.id,
          name: departmentUnit.name,
          code: departmentUnit.code,
          manager_id: departmentUnit.manager_id,
          manager: departmentUnit.manager,
          positions: this.departmentPositions
        }
        
        // 通过路由参数传递部门信息到员工创建页面
        this.$router.push({
          path: '/employees/create',
          query: {
            department_id: departmentInfo.id,
            department_name: departmentInfo.name,
            department_code: departmentInfo.code,
            direct_manager_id: departmentInfo.manager_id,
            direct_manager_name: departmentInfo.manager?.name,
            // 将职位信息作为JSON字符串传递
            department_positions: JSON.stringify(departmentInfo.positions || [])
          }
        })
      })
    },

    // 创建子单元
    createChildUnit() {
      if (!this.selectedUnit) return
      this.showCreateUnitDialog = true
    },
    
    // 创建员工
    createEmployee() {
      if (!this.selectedUnit) return
      
      // 构建员工创建页面需要的部门和管理者信息
      const departmentInfo = {
        id: this.selectedUnit.id,
        name: this.selectedUnit.name,
        code: this.selectedUnit.code,
        manager_id: this.selectedUnit.manager_id,
        manager: this.selectedUnit.manager,
        positions: this.departmentPositions // 传递当前部门的职位列表
      }
      
      // 通过路由参数传递部门信息到员工创建页面
      this.$router.push({
        path: '/employees/create',
        query: {
          department_id: departmentInfo.id,
          department_name: departmentInfo.name,
          department_code: departmentInfo.code,
          direct_manager_id: departmentInfo.manager_id,
          direct_manager_name: departmentInfo.manager?.name,
          // 将职位信息作为JSON字符串传递
          department_positions: JSON.stringify(departmentInfo.positions || [])
        }
      })
    },


    // 刷新图表
    async refreshChart() {
      await this.initData()
      this.$message.success('图表已刷新')
    },

    // 展开全部
    expandAll() {
      if (this.$refs.orgChart) {
        this.$refs.orgChart.expandAll()
      }
    },

    // 收起全部
    collapseAll() {
      if (this.$refs.orgChart) {
        this.$refs.orgChart.collapseAll()
      }
    },

    // 导出图表
    exportChart() {
      if (this.$refs.orgChart) {
        this.$refs.orgChart.exportChart()
      }
    },

    // 创建成功回调
    handleCreateSuccess() {
      this.showCreateUnitDialog = false
      this.loadOrganizationTree()
      this.$message.success('创建成功')
    },

    // 编辑成功回调
    handleEditSuccess() {
      this.showEditUnitDialog = false
      this.editingUnit = null
      this.loadOrganizationTree()
      this.$message.success('编辑成功')
    },


    // 处理视图类型变化
    handleViewTypeChange(newViewType) {
      this.viewType = newViewType
    },
    
    // 查看员工详情
    handleViewEmployee(employeeData) {
      this.$router.push(`/employees/${employeeData.id}`)
    },
    
    // 编辑员工
    handleEditEmployee(employeeData) {
      this.$router.push(`/employees/${employeeData.id}/edit`)
    },

    // 查看员工（侧边栏按钮）
    async viewEmployee() {
      if (this.selectedUnit?.type === 'employee' && this.selectedUnit.employeeData) {
        try {
          // 先验证员工是否存在
          await employeeApiService.getEmployee(this.selectedUnit.employeeData.id)
          this.$router.push(`/employees/${this.selectedUnit.employeeData.id}`)
        } catch (error) {
          console.error('员工不存在:', error)
          this.$message.error('员工信息不存在或已被删除')
        }
      }
    },

    // 修改员工（侧边栏按钮）
    async editEmployee() {
      if (this.selectedUnit?.type === 'employee' && this.selectedUnit.employeeData) {
        try {
          // 先验证员工是否存在
          await employeeApiService.getEmployee(this.selectedUnit.employeeData.id)
          this.$router.push(`/employees/${this.selectedUnit.employeeData.id}/edit`)
        } catch (error) {
          console.error('员工不存在:', error)
          this.$message.error('员工信息不存在或已被删除，请刷新页面重新加载数据')
          // 刷新组织架构数据
          this.initData()
        }
      }
    }
  }
}
</script>

<style scoped>
.organization-chart-view {
  min-height: 100vh;
  background: #f8f9fa;
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.2);
}

.page-title {
  display: flex;
  align-items: center;
  gap: 16px;
}

.title-icon {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.title-content h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 600;
  color: white;
}

.title-content p {
  margin: 8px 0 0;
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}

.page-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  padding: 12px 24px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.1);
  color: white;
  backdrop-filter: blur(10px);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.primary-btn {
  background: rgba(255, 255, 255, 0.95);
  color: #667eea;
  border: none;
}

.primary-btn:hover {
  background: white;
  color: #667eea;
}


/* 主要内容区域 */
.main-content {
  position: relative;
  height: calc(100vh - 200px);
}


/* 悬浮的选中单元信息面板 */
.floating-unit-info {
  position: absolute;
  bottom: 20px;
  right: 20px;
  width: 320px;
  max-height: 500px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  overflow: hidden;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  animation: slideInFromRight 0.3s ease-out;
}

@keyframes slideInFromRight {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

.floating-unit-info:hover {
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.2);
  transform: translateY(-2px);
}

.floating-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.floating-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: white;
}

.close-btn {
  color: white !important;
  font-size: 16px;
  padding: 0 !important;
  min-width: auto !important;
}

.close-btn:hover {
  color: rgba(255, 255, 255, 0.8) !important;
}

/* 选中单元信息 */
.selected-unit-info {
  padding: 16px 20px;
  background: transparent;
  max-height: 400px;
  overflow-y: auto;
}

/* 自定义滚动条样式 */
.selected-unit-info::-webkit-scrollbar {
  width: 6px;
}

.selected-unit-info::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.05);
  border-radius: 3px;
}

.selected-unit-info::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
}

.selected-unit-info::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.3);
}

.unit-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
}

.unit-details {
  margin-bottom: 16px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 12px;
}

.detail-item .label {
  color: #909399;
}

.detail-item .value {
  color: #303133;
  font-weight: 500;
}

.unit-actions {
  display: flex;
  gap: 8px;
}

.unit-actions .el-button {
  flex: 1;
}

/* 部门职位列表样式 */
.department-positions {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e4e7ed;
}

.department-positions h4 {
  margin: 0 0 8px 0;
  font-size: 12px;
  font-weight: 600;
  color: #303133;
}

.positions-list {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.position-tag {
  font-size: 11px;
  padding: 2px 6px;
}

/* 员工详情样式 */
.employee-details {
  background: #f8f9fa;
  padding: 8px;
  border-radius: 4px;
  margin: 8px 0;
}

/* 图表区域（占满全屏） */
.chart-container {
  width: 100%;
  height: 100%;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.chart-wrapper {
  height: 100%;
  width: 100%;
}

.main-org-chart {
  height: 100%;
  width: 100%;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .floating-unit-info {
    width: 280px;
    max-height: 400px;
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
  
  .page-actions {
    width: 100%;
    justify-content: flex-start;
  }
  
  .floating-unit-info {
    width: 260px;
    max-height: 350px;
    bottom: 10px;
    right: 10px;
  }
  
  .chart-container {
    height: 500px;
  }
}
</style>