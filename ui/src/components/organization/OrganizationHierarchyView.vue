<template>
  <div class="hierarchy-view">
    <div class="view-header">
      <h3>组织层级详情</h3>
      <div class="header-actions">
        <el-button size="small" @click="refreshData">
          <i class="el-icon-refresh"></i>
          刷新
        </el-button>
        <el-button size="small" type="primary" @click="exportHierarchy">
          <i class="el-icon-download"></i>
          导出
        </el-button>
      </div>
    </div>

    <div v-if="!selectedUnit" class="empty-state">
      <div class="empty-content">
        <i class="el-icon-office-building"></i>
        <p>请从左侧选择一个组织单元查看详情</p>
      </div>
    </div>

    <div v-else class="unit-details">
      <!-- 基本信息卡片 -->
      <el-card class="info-card">
        <div slot="header" class="card-header">
          <div class="unit-title">
            <i :class="getUnitTypeIcon(selectedUnit.type)"></i>
            <span class="unit-name">{{ selectedUnit.name }}</span>
            <el-tag 
              :type="getUnitTypeColor(selectedUnit.type)"
              size="small"
            >
              {{ getUnitTypeName(selectedUnit.type) }}
            </el-tag>
          </div>
          <div class="unit-actions">
            <el-button size="mini" @click="$emit('edit-unit', selectedUnit)">
              <i class="el-icon-edit"></i>
              编辑
            </el-button>
            <el-button size="mini" type="danger" @click="$emit('delete-unit', selectedUnit)">
              <i class="el-icon-delete"></i>
              删除
            </el-button>
          </div>
        </div>

        <div class="unit-info">
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="info-item">
                <label>组织编码：</label>
                <span>{{ selectedUnit.code || '-' }}</span>
              </div>
              <div class="info-item">
                <label>组织层级：</label>
                <span>{{ selectedUnit.level || '-' }}</span>
              </div>
              <div class="info-item">
                <label>状态：</label>
                <el-tag 
                  :type="(selectedUnit.is_active || selectedUnit.status === 'active') ? 'success' : 'danger'"
                  size="mini"
                >
                  {{ (selectedUnit.is_active || selectedUnit.status === 'active') ? '激活' : '停用' }}
                </el-tag>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="info-item">
                <label>负责人：</label>
                <span>{{ (selectedUnit.manager && selectedUnit.manager.name) || selectedUnit.managerName || '-' }}</span>
              </div>
              <div class="info-item">
                <label>创建时间：</label>
                <span>{{ formatDate(selectedUnit.created_at || selectedUnit.createdAt) }}</span>
              </div>
              <div class="info-item">
                <label>更新时间：</label>
                <span>{{ formatDate(selectedUnit.updated_at || selectedUnit.updatedAt) }}</span>
              </div>
            </el-col>
          </el-row>
        </div>
      </el-card>

      <!-- 统计信息卡片 -->
      <el-card class="stats-card">
        <div slot="header">
          <span>统计信息</span>
        </div>
        <div class="stats-row">
          <div class="stat-item">
            <div class="stat-icon"><i class="el-icon-user"></i></div>
            <div class="stat-content">
              <div class="stat-value">{{ selectedUnit.employeeCount || 0 }}</div>
              <div class="stat-label">直接员工</div>
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-icon"><i class="el-icon-s-custom"></i></div>
            <div class="stat-content">
              <div class="stat-value">{{ selectedUnit.hierarchicalEmployeeCount || 0 }}</div>
              <div class="stat-label">总员工数</div>
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-icon"><i class="el-icon-office-building"></i></div>
            <div class="stat-content">
              <div class="stat-value">{{ selectedUnit.subunitCount || 0 }}</div>
              <div class="stat-label">子单元数</div>
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-icon"><i class="el-icon-s-grid"></i></div>
            <div class="stat-content">
              <div class="stat-value">{{ getMaxLevel(selectedUnit) }}</div>
              <div class="stat-label">最大层级</div>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 层级关系图 -->
      <el-card class="hierarchy-card">
        <div slot="header" class="card-header">
          <span>层级关系</span>
          <div class="header-actions">
            <el-button-group>
              <el-button 
                size="mini" 
                :type="viewMode === 'tree' ? 'primary' : ''"
                @click="viewMode = 'tree'"
              >
                <i class="el-icon-s-data"></i>
                树形视图
              </el-button>
              <el-button 
                size="mini" 
                :type="viewMode === 'chart' ? 'primary' : ''"
                @click="viewMode = 'chart'"
              >
                <i class="el-icon-s-grid"></i>
                图形视图
              </el-button>
            </el-button-group>
          </div>
        </div>
        
        <!-- 层级路径和关系图 -->
        <div class="hierarchy-section">
          <div class="hierarchy-path">
            <div v-if="loading.path" class="path-loading">
              <i class="el-icon-loading"></i>
              <span>加载路径中...</span>
            </div>
            <div v-else class="breadcrumb-container">
              <div class="path-title">层级路径：</div>
              <el-breadcrumb separator=">" class="hierarchy-breadcrumb">
                <el-breadcrumb-item 
                  v-for="(item, index) in hierarchyPath"
                  :key="index"
                  @click.native="$emit('unit-selected', item)"
                  :class="{'current-item': item.id === selectedUnit.id}"
                >
                  <div class="breadcrumb-item-content">
                    <i :class="getUnitTypeIcon(item.type)" :style="getIconColorStyle(item.type)"></i>
                    <span>{{ item.name }}</span>
                    <span class="level-indicator">L{{ item.level || (index + 1) }}</span>
                  </div>
                </el-breadcrumb-item>
              </el-breadcrumb>
            </div>
          </div>
          
          <!-- 组织关系图 -->
          <div class="org-relationship-chart">
            <div class="relationship-title">组织关系：</div>
            <div class="relationship-flow">
              <div 
                v-for="(unit, index) in hierarchyPath" 
                :key="unit.id"
                class="relationship-node"
                :class="{'current-node': unit.id === selectedUnit.id}"
                @click="$emit('unit-selected', unit)"
              >
                <div class="node-card">
                  <div class="node-header">
                    <i :class="getUnitTypeIcon(unit.type)" :style="getIconColorStyle(unit.type)"></i>
                    <span class="node-title">{{ unit.name }}</span>
                  </div>
                  <div class="node-info">
                    <div class="node-level">L{{ unit.level || (index + 1) }}</div>
                    <div v-if="unit.manager && unit.manager.name" class="node-manager">
                      负责人：{{ unit.manager.name }}
                    </div>
                  </div>
                </div>
                <div v-if="index < hierarchyPath.length - 1" class="relationship-arrow">
                  <i class="el-icon-right"></i>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Workday风格组织图形视图 -->
        <div v-if="viewMode === 'chart'" class="hierarchy-chart">
          <WorkdayOrgChart 
            :organization-data="hierarchyTree"
            :selected-unit="selectedUnit"
            @unit-selected="handleChartNodeClick"
          />
        </div>
        
        <!-- 树形视图 -->
        <div v-if="viewMode === 'tree'" class="hierarchy-tree">
          <div v-if="loading.tree" class="tree-loading">
            <i class="el-icon-loading"></i>
            <span>加载组织架构中...</span>
          </div>
          <el-tree
            v-else
            :data="hierarchyTree"
            :props="treeProps"
            node-key="id"
            :expand-on-click-node="false"
            :default-expanded-keys="expandedKeys"
            @node-click="handleTreeNodeClick"
          >
            <div class="tree-node" slot-scope="{ node, data }">
              <div class="node-content">
                <div class="node-info">
                  <i :class="getUnitTypeIcon(data.type)"></i>
                  <span class="node-name">{{ data.name }}</span>
                  <el-tag 
                    v-if="data.type !== 'company'" 
                    size="mini" 
                    :type="getUnitTypeColor(data.type)"
                  >
                    {{ getUnitTypeName(data.type) }}
                  </el-tag>
                </div>
                <div class="node-stats">
                  <span class="employee-count">
                    <i class="el-icon-user"></i>
                    {{ data.employeeCount || 0 }}
                  </span>
                </div>
              </div>
            </div>
          </el-tree>
        </div>
      </el-card>

      <!-- 子单元列表 -->
      <el-card class="subunits-card">
        <div slot="header" class="card-header">
          <span>子单元 ({{ subunits.length }})</span>
          <el-button 
            size="mini" 
            type="primary" 
            @click="showAddSubunitDialog = true"
          >
            <i class="el-icon-plus"></i>
            添加子单元
          </el-button>
        </div>
        
        <div v-if="subunits.length === 0" class="empty-subunits">
          <p>暂无子单元</p>
        </div>
        
        <div v-else class="subunits-grid">
          <div 
            v-for="subunit in subunits"
            :key="subunit.id"
            class="subunit-card"
            @click="$emit('unit-selected', subunit)"
          >
            <div class="subunit-header">
              <div class="subunit-title">
                <i :class="getUnitTypeIcon(subunit.type)" class="subunit-icon"></i>
                <span class="subunit-name">{{ subunit.name }}</span>
              </div>
              <div class="subunit-manager">
                <span v-if="subunit.manager && subunit.manager.name" class="manager-text">
                  <i class="el-icon-user-solid"></i>
                  {{ subunit.manager.name }}
                </span>
                <span v-else class="no-manager-text">无负责人</span>
              </div>
            </div>
            <div class="subunit-stats">
              <div class="stat-mini">
                <i class="el-icon-user"></i>
                <span>{{ subunit.employeeCount || 0 }}</span>
              </div>
              <div class="stat-mini">
                <i class="el-icon-office-building"></i>
                <span>{{ subunit.subunitCount || 0 }}</span>
              </div>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 直接员工列表 -->
      <el-card class="employees-card">
        <div slot="header" class="card-header">
          <span>直接员工 ({{ directEmployees.length }})</span>
          <el-button 
            size="mini" 
            type="primary" 
            @click="showAssignEmployeeDialog = true"
          >
            <i class="el-icon-plus"></i>
            分配员工
          </el-button>
        </div>
        
        <div v-if="directEmployees.length === 0" class="empty-employees">
          <p>暂无直接员工</p>
        </div>
        
        <div v-else class="employees-grid">
          <div 
            v-for="employee in directEmployees"
            :key="employee.id"
            class="employee-card"
          >
            <div class="employee-info">
              <div class="employee-header">
                <span class="employee-name">{{ employee.name }}</span>
                <span class="employee-id">{{ employee.employeeId }}</span>
              </div>
              <div class="employee-details">
                <div class="employee-position">
                  <i class="el-icon-s-custom"></i>
                  {{ employee.position || '未设置职位' }}
                </div>
                <div v-if="employee.directManager" class="employee-manager">
                  <i class="el-icon-user"></i>
                  直接上级：{{ employee.directManager }}
                </div>
              </div>
            </div>
            <div class="employee-actions">
              <el-button 
                size="mini" 
                type="text" 
                icon="el-icon-edit"
                @click="handleEditAssignment(employee)"
                title="编辑"
              ></el-button>
            </div>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script>
import dayjs from 'dayjs'
import { organizationApiService } from '@/services/organizationApi'
import WorkdayOrgChart from './WorkdayOrgChart.vue'

export default {
  name: 'OrganizationHierarchyView',
  components: {
    WorkdayOrgChart
  },
  props: {
    selectedUnit: {
      type: Object,
      default: null
    },
    organizationData: {
      type: Object,
      default: () => ({})
    }
  },
  data() {
    return {
      showAddSubunitDialog: false,
      showAssignEmployeeDialog: false,
      hierarchyPath: [],
      subunits: [],
      directEmployees: [],
      viewMode: 'tree',
      hierarchyTree: [],
      expandedKeys: [],
      loading: {
        details: false,
        path: false,
        tree: false
      },
      treeProps: {
        children: 'children',
        label: 'name'
      },
      
      // 单元类型配置
      unitTypes: [
        { value: 'company', label: '公司', icon: 'el-icon-office-building', color: '' },
        { value: 'business_unit', label: '业务单元', icon: 'el-icon-s-cooperation', color: 'primary' },
        { value: 'department', label: '部门', icon: 'el-icon-user', color: 'success' },
        { value: 'team', label: '团队', icon: 'el-icon-s-custom', color: 'warning' },
        { value: 'cost_center', label: '成本中心', icon: 'el-icon-s-finance', color: 'info' },
        { value: 'location', label: '地理位置', icon: 'el-icon-location', color: 'danger' },
        { value: 'project', label: '项目组', icon: 'el-icon-s-flag', color: 'warning' }
      ],
      
      // 分配类型配置
      assignmentTypes: [
        { value: 'primary', label: '主要', color: 'success' },
        { value: 'additional', label: '额外', color: 'warning' },
        { value: 'temporary', label: '临时', color: 'info' },
        { value: 'project', label: '项目', color: 'primary' }
      ]
    }
  },
  watch: {
    selectedUnit: {
      handler(newVal) {
        if (newVal) {
          this.loadUnitDetails()
        }
      },
      immediate: true
    }
  },
  methods: {
    // 加载单元详情
    async loadUnitDetails() {
      if (!this.selectedUnit) return
      
      this.loading.details = true
      try {
        // 加载层级路径
        await this.buildHierarchyPath()
        
        // 加载子单元
        await this.loadSubunits()
        
        // 加载直接员工
        await this.loadDirectEmployees()
        
      } catch (error) {
        console.error('Failed to load unit details:', error)
        this.$message.error('加载组织单元详情失败')
      } finally {
        this.loading.details = false
      }
    },
    
    // 构建层级路径
    async buildHierarchyPath() {
      if (!this.selectedUnit) {
        this.hierarchyPath = []
        return
      }

      this.loading.path = true
      try {
        // 使用部门树API获取层级路径
        const { departmentApi } = await import('@/services/departmentApi')
        try {
          const response = await departmentApi.getDepartmentPath(this.selectedUnit.id)
          const responseData = response.data || response
          this.hierarchyPath = responseData.data || responseData || []
        } catch (error) {
          console.warn('Department path API not available, building from hierarchyPath')
          // 如果API不可用，使用组织单元的 hierarchyPath 构建路径
          if (this.selectedUnit.hierarchyPath) {
            const pathNames = this.selectedUnit.hierarchyPath.split(' > ')
            this.hierarchyPath = pathNames.map((name, index) => ({
              id: `path_${index}`,
              name: name,
              type: index === 0 ? 'company' : 'department'
            }))
          } else {
            this.hierarchyPath = [this.selectedUnit]
          }
        }
      } catch (error) {
        console.error('Failed to load hierarchy path:', error)
        // 发生错误时，使用当前选中单位作为路径
        this.hierarchyPath = [this.selectedUnit]
      } finally {
        this.loading.path = false
      }
      
      // 构建层级树
      await this.buildHierarchyTree()
    },
    
    // 构建层级树
    async buildHierarchyTree() {
      this.loading.tree = true
      try {
        // 调用部门树API获取完整的组织架构树
        const { departmentApi } = await import('@/services/departmentApi')
        const response = await departmentApi.getDepartmentTree()
        const responseData = response.data || response
        this.hierarchyTree = responseData.data || responseData || []
        
        // 设置展开的节点（展开到当前选中单位的路径上的所有节点）
        this.expandedKeys = this.getExpandedKeys()
      } catch (error) {
        console.error('Failed to load hierarchy tree:', error)
        // 发生错误时使用空树
        this.hierarchyTree = []
        this.expandedKeys = []
      } finally {
        this.loading.tree = false
      }
    },

    // 获取需要展开的节点ID列表
    getExpandedKeys() {
      if (!this.hierarchyPath || this.hierarchyPath.length === 0) {
        return []
      }
      
      // 展开层级路径中所有节点的ID
      return this.hierarchyPath.map(item => item.id)
    },
    
    // 处理树节点点击
    handleTreeNodeClick(data) {
      this.$emit('unit-selected', data)
    },
    
    // 处理图表节点点击
    handleChartNodeClick(nodeData) {
      this.$emit('unit-selected', nodeData)
    },
    
    // 加载子单元
    async loadSubunits() {
      if (!this.selectedUnit) {
        this.subunits = []
        return
      }
      
      try {
        // 使用selectedUnit的children数据，如果存在的话
        if (this.selectedUnit.children && this.selectedUnit.children.length > 0) {
          this.subunits = this.selectedUnit.children.map(child => ({
            id: child.id,
            name: child.name,
            code: child.code,
            type: child.type,
            level: child.level,
            description: child.description,
            manager_id: child.manager_id,
            manager: child.manager,
            employeeCount: child.employeeCount || 0,
            subunitCount: child.subunitCount || 0,
            hierarchicalEmployeeCount: child.hierarchicalEmployeeCount || 0
          }))
        } else {
          // 如果没有children数据，从部门API获取
          const { departmentApi } = await import('@/services/departmentApi')
          try {
            const response = await departmentApi.getDepartmentsByParent(this.selectedUnit.id)
            const responseData = response.data || response
            this.subunits = responseData.data || responseData || []
          } catch (error) {
            console.warn('Department children API not available')
            this.subunits = []
          }
        }
      } catch (error) {
        console.error('Failed to load subunits:', error)
        this.subunits = []
      }
    },
    
    // 加载直接员工
    async loadDirectEmployees() {
      if (!this.selectedUnit) {
        this.directEmployees = []
        return
      }
      
      try {
        // 调用API获取该部门的直接员工
        const { employeeApiService } = await import('@/services/employeeApi')
        const response = await employeeApiService.getEmployees({
          department_id: this.selectedUnit.id,
          page: 1,
          pageSize: 1000
        })
        
        const responseData = response.data || response
        const employees = responseData.data || responseData.employees || []
        
        // 转换为直接员工格式
        this.directEmployees = employees.map(emp => ({
          id: emp.id,
          name: emp.name,
          employeeId: emp.employee_id,
          email: emp.email,
          phone: emp.phone,
          position: emp.position ? emp.position.name : '未设置',
          positionId: emp.position ? emp.position.id : null,
          assignmentType: emp.assignment_type || 'primary',
          managementType: emp.management_type || 'line',
          workPercentage: emp.work_percentage || 100,
          directManager: emp.manager ? emp.manager.name : '无',
          directManagerId: emp.manager_id,
          effectiveDate: emp.hire_date,
          expirationDate: emp.contract_end_date,
          status: emp.status,
          // 保留完整的员工信息
          employeeData: emp
        }))
      } catch (error) {
        console.error('Failed to load direct employees:', error)
        this.directEmployees = []
      }
    },
    
    // 刷新数据
    refreshData() {
      this.loadUnitDetails()
    },
    
    // 导出层级
    exportHierarchy() {
      this.$message.info('导出功能开发中')
    },
    
    // 编辑分配
    handleEditAssignment(assignment) {
      this.$emit('edit-assignment', assignment)
    },
    
    // 移除分配
    handleRemoveAssignment(assignment) {
      this.$emit('remove-assignment', assignment)
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
    
    // 获取图标颜色样式
    getIconColorStyle(type) {
      const typeColors = {
        company: '#595959',
        department: '#52c41a', 
        team: '#fa8c16',
        business_unit: '#409eff',
        location: '#ff4d4f',
        cost_center: '#722ed1',
        project: '#eb2f96'
      }
      return {
        color: typeColors[type] || '#409eff'
      }
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
    
    // 格式化日期
    formatDate(date) {
      return date ? dayjs(date).format('YYYY-MM-DD') : '-'
    },
    
    // 计算最大层级
    getMaxLevel(unit) {
      if (!unit || !unit.children || unit.children.length === 0) {
        return unit ? unit.level || 1 : 0
      }
      
      let maxLevel = unit.level || 1
      const getChildMaxLevel = (nodes) => {
        for (const child of nodes) {
          if (child.level > maxLevel) {
            maxLevel = child.level
          }
          if (child.children && child.children.length > 0) {
            getChildMaxLevel(child.children)
          }
        }
      }
      
      getChildMaxLevel(unit.children)
      return maxLevel
    }
  }
}
</script>

<style scoped>
.hierarchy-view {
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

.unit-details {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-card,
.stats-card,
.hierarchy-card,
.subunits-card,
.employees-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.unit-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.unit-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.unit-actions {
  display: flex;
  gap: 8px;
}

.unit-info {
  padding: 16px 0;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.info-item label {
  min-width: 80px;
  color: #606266;
  font-size: 14px;
}

.info-item span {
  color: #303133;
  font-size: 14px;
}

.stats-row {
  display: flex;
  gap: 24px;
  justify-content: space-around;
  flex-wrap: wrap;
}

.stats-card .stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  flex: 1;
  min-width: 120px;
}

.stat-icon {
  width: 36px;
  height: 36px;
  background: #409eff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
  line-height: 1;
}

.stat-label {
  font-size: 12px;
  color: #909399;
  line-height: 1;
}

.hierarchy-section {
  background: #fafafa;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 16px;
}

.hierarchy-path {
  margin-bottom: 20px;
}

.breadcrumb-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.path-title {
  font-weight: 500;
  color: #303133;
  font-size: 14px;
  flex-shrink: 0;
}

.hierarchy-breadcrumb :deep(.el-breadcrumb__inner) {
  cursor: pointer;
  padding: 0;
}

.hierarchy-breadcrumb :deep(.el-breadcrumb__inner:hover) {
  color: #409eff;
}

.breadcrumb-item-content {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.3s;
}

.breadcrumb-item-content:hover {
  background: #e6f3ff;
}

.current-item .breadcrumb-item-content {
  background: #409eff;
  color: white;
}

.current-item .breadcrumb-item-content i {
  color: white !important;
}

.level-indicator {
  background: rgba(255, 255, 255, 0.3);
  padding: 2px 4px;
  border-radius: 8px;
  font-size: 10px;
  font-weight: bold;
}

.org-relationship-chart {
  border-top: 1px solid #e4e7ed;
  padding-top: 16px;
}

.relationship-title {
  font-weight: 500;
  color: #303133;
  font-size: 14px;
  margin-bottom: 12px;
}

.relationship-flow {
  display: flex;
  align-items: center;
  gap: 8px;
  overflow-x: auto;
  padding: 8px 0;
}

.relationship-node {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.node-card {
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  padding: 12px;
  background: white;
  cursor: pointer;
  transition: all 0.3s;
  min-width: 140px;
}

.node-card:hover {
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
}

.current-node .node-card {
  border-color: #409eff;
  background: #e6f3ff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2);
}

.node-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
}

.node-title {
  font-weight: 500;
  font-size: 13px;
  color: #303133;
}

.node-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.node-level {
  background: #f0f9ff;
  color: #409eff;
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 8px;
  font-weight: bold;
  align-self: flex-start;
}

.node-manager {
  font-size: 11px;
  color: #67c23a;
}

.relationship-arrow {
  color: #409eff;
  font-size: 16px;
  font-weight: bold;
  opacity: 0.7;
}

.path-loading,
.tree-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 20px;
  color: #909399;
  font-size: 14px;
}

.path-loading i,
.tree-loading i {
  font-size: 16px;
}

.hierarchy-chart {
  margin-top: 20px;
}

.chart-container {
  width: 100%;
  height: 400px;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  overflow: hidden;
}

.org-chart {
  width: 100%;
  height: 100%;
}

.chart-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #909399;
  background: #fafafa;
}

.chart-placeholder i {
  font-size: 48px;
  margin-bottom: 16px;
  color: #c0c4cc;
}

.chart-placeholder p {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
}

.chart-placeholder .tip {
  margin-top: 8px;
  font-size: 12px;
  color: #c0c4cc;
}

.hierarchy-tree {
  margin-top: 20px;
  max-height: 400px;
  overflow: auto;
}

.hierarchy-tree .tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 8px 12px;
  border-radius: 6px;
  transition: background-color 0.3s;
}

.hierarchy-tree .tree-node:hover {
  background-color: #f0f9ff;
}

.hierarchy-tree .node-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.hierarchy-tree .node-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.hierarchy-tree .node-name {
  font-weight: 500;
  color: #303133;
  font-size: 14px;
}

.hierarchy-tree .node-stats {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #909399;
}

.hierarchy-tree .employee-count {
  display: flex;
  align-items: center;
  gap: 4px;
}

.empty-subunits,
.empty-employees {
  text-align: center;
  color: #909399;
  padding: 40px 0;
}

.subunits-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.subunit-card {
  padding: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  background: white;
}

.subunit-card:hover {
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
  transform: translateY(-2px);
}

.subunit-header {
  margin-bottom: 12px;
}

.subunit-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.subunit-icon {
  color: #409eff;
  font-size: 16px;
}

.subunit-name {
  font-weight: 500;
  color: #303133;
  font-size: 14px;
}

.subunit-manager {
  font-size: 12px;
}

.manager-text {
  color: #67c23a;
  display: flex;
  align-items: center;
  gap: 4px;
}

.no-manager-text {
  color: #909399;
  font-style: italic;
}

.subunit-stats {
  display: flex;
  gap: 16px;
  justify-content: flex-end;
}

.stat-mini {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #909399;
  background: #f5f7fa;
  padding: 4px 8px;
  border-radius: 12px;
}

.stat-mini i {
  color: #409eff;
}

.employees-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 12px;
  margin-top: 16px;
  max-height: 400px;
  overflow-y: auto;
}

.employee-card {
  padding: 12px;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  background: white;
  transition: all 0.3s;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.employee-card:hover {
  border-color: #409eff;
  box-shadow: 0 2px 4px rgba(64, 158, 255, 0.1);
}

.employee-info {
  flex: 1;
}

.employee-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.employee-name {
  font-weight: 500;
  color: #303133;
  font-size: 14px;
}

.employee-id {
  font-size: 12px;
  color: #909399;
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 10px;
}

.employee-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.employee-position,
.employee-manager {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #606266;
}

.employee-position i,
.employee-manager i {
  color: #409eff;
}

.employee-actions {
  opacity: 0;
  transition: opacity 0.3s;
}

.employee-card:hover .employee-actions {
  opacity: 1;
}
</style>