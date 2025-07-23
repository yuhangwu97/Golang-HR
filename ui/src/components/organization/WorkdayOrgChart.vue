<template>
  <div class="workday-org-chart">
    <div class="chart-header">
      <div class="chart-title">
        <h3>组织架构图</h3>
        <span class="subtitle">{{ selectedUnit ? selectedUnit.name : '全公司' }} 组织结构</span>
      </div>
      <div class="chart-controls">
        <el-button-group>
          <el-button 
            size="small" 
            :type="viewType === 'tree' ? 'primary' : ''"
            @click="$emit('view-type-change', 'tree')"
            icon="el-icon-s-data"
          >
            树形图
          </el-button>
          <el-button 
            size="small" 
            :type="viewType === 'radial' ? 'primary' : ''"
            @click="$emit('view-type-change', 'radial')"
            icon="el-icon-pie-chart"
          >
            径向图
          </el-button>
          <el-button 
            size="small" 
            :type="viewType === 'horizontal' ? 'primary' : ''"
            @click="$emit('view-type-change', 'horizontal')"
            icon="el-icon-s-grid"
          >
            横向图
          </el-button>
        </el-button-group>
        
        <el-button-group class="zoom-controls">
          <el-button size="small" @click="zoomIn" icon="el-icon-zoom-in" title="放大"></el-button>
          <el-button size="small" @click="zoomOut" icon="el-icon-zoom-out" title="缩小"></el-button>
          <el-button size="small" @click="resetZoom" icon="el-icon-refresh" title="重置"></el-button>
        </el-button-group>
        
        <el-button-group>
          <el-button size="small" @click="expandAll" icon="el-icon-folder-opened">展开全部</el-button>
          <el-button size="small" @click="collapseAll" icon="el-icon-folder">收起全部</el-button>
        </el-button-group>
        
        <el-button size="small" @click="exportChart" icon="el-icon-download">导出</el-button>
      </div>
    </div>
    
    <div class="chart-container" ref="chartContainer">
      <svg ref="orgChartSvg" class="org-chart-svg">
        <defs>
          <!-- 定义渐变和阴影 -->
          <linearGradient id="nodeGradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" style="stop-color:#409eff;stop-opacity:0.1" />
            <stop offset="100%" style="stop-color:#409eff;stop-opacity:0.3" />
          </linearGradient>
          <filter id="shadow" x="-50%" y="-50%" width="200%" height="200%">
            <feDropShadow dx="2" dy="2" stdDeviation="3" flood-color="#00000020"/>
          </filter>
        </defs>
        <g class="chart-content"></g>
      </svg>
      
      <!-- 节点详情弹窗 -->
      <div 
        v-if="hoveredNode" 
        class="node-tooltip" 
        :style="tooltipStyle"
        @mouseenter="keepTooltip = true"
        @mouseleave="hideTooltip"
      >
        <div class="tooltip-header">
          <div class="node-avatar">
            <img v-if="hoveredNode.data.manager && hoveredNode.data.manager.avatar" 
                 :src="hoveredNode.data.manager.avatar" 
                 :alt="hoveredNode.data.manager.name">
            <i v-else :class="getUnitTypeIcon(hoveredNode.data.type)" class="default-avatar"></i>
          </div>
          <div class="node-info">
            <div class="node-name">{{ hoveredNode.data.name }}</div>
            <div class="node-manager">
              {{ hoveredNode.data.manager ? hoveredNode.data.manager.name : '无负责人' }}
            </div>
          </div>
        </div>
        <div class="tooltip-content">
          <div v-if="hoveredNode.data.type === 'employee'" class="employee-tooltip">
            <div class="info-row">
              <span class="label">员工工号：</span>
              <span class="value">{{ hoveredNode.data.employeeData.employee_id || '-' }}</span>
            </div>
            <div class="info-row">
              <span class="label">部门：</span>
              <span class="value">{{ hoveredNode.data.employeeData.department ? hoveredNode.data.employeeData.department.name : '-' }}</span>
            </div>
            <div class="info-row">
              <span class="label">职位：</span>
              <span class="value">{{ hoveredNode.data.employeeData.position ? hoveredNode.data.employeeData.position.name : '未设置' }}</span>
            </div>
            <div class="info-row">
              <span class="label">状态：</span>
              <span class="value">{{ hoveredNode.data.employeeData.status === 'active' ? '在职' : '离职' }}</span>
            </div>
          </div>
          <div v-else class="department-tooltip">
            <div class="info-row">
              <span class="label">部门类型：</span>
              <span class="value">{{ getUnitTypeName(hoveredNode.data.type) }}</span>
            </div>
            <div class="info-row">
              <span class="label">直接员工：</span>
              <span class="value">{{ hoveredNode.data.employeeCount || 0 }}人</span>
            </div>
            <div class="info-row">
              <span class="label">子部门：</span>
              <span class="value">{{ hoveredNode.data.subunitCount || 0 }}个</span>
            </div>
            <div class="info-row">
              <span class="label">组织层级：</span>
              <span class="value">L{{ hoveredNode.data.level || 1 }}</span>
            </div>
          </div>
        </div>
        <div class="tooltip-actions">
          <el-button-group v-if="hoveredNode.data.type !== 'employee'">
            <el-button size="mini" type="primary" @click="selectNode(hoveredNode.data)">
              <i class="el-icon-view"></i> 查看详情
            </el-button>
            <el-dropdown trigger="click" size="mini" @command="handleNodeAction">
              <el-button size="mini" type="success">
                <i class="el-icon-plus"></i> 创建
                <i class="el-icon-arrow-down el-icon--right"></i>
              </el-button>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item :command="{action: 'createDepartment', data: hoveredNode.data}">
                  <i class="el-icon-office-building"></i> 创建子部门
                </el-dropdown-item>
                <el-dropdown-item :command="{action: 'createEmployee', data: hoveredNode.data}">
                  <i class="el-icon-user"></i> 创建员工
                </el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
          </el-button-group>
          <el-button-group v-else>
            <el-button size="mini" type="primary" @click="viewEmployeeDetail(hoveredNode.data.employeeData)">
              <i class="el-icon-view"></i> 查看员工
            </el-button>
            <el-button size="mini" type="warning" @click="editEmployee(hoveredNode.data.employeeData)">
              <i class="el-icon-edit"></i> 编辑
            </el-button>
          </el-button-group>
        </div>
      </div>
    </div>
    
    <!-- 图例 -->
    <div class="chart-legend">
      <div class="legend-title">图例</div>
      <div class="legend-items">
        <div class="legend-item" v-for="unitType in legendItems" :key="unitType.value">
          <i :class="unitType.icon" :style="{ color: unitType.color }"></i>
          <span>{{ unitType.label }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import * as d3 from 'd3'

export default {
  name: 'WorkdayOrgChart',
  props: {
    organizationData: {
      type: Array,
      default: () => []
    },
    selectedUnit: {
      type: Object,
      default: null
    },
    viewType: {
      type: String,
      default: 'tree'
    }
  },
  data() {
    return {
      svg: null,
      g: null,
      zoom: null,
      root: null,
      hoveredNode: null,
      keepTooltip: false,
      tooltipStyle: {
        left: '0px',
        top: '0px',
        display: 'none'
      },
      expandedNodes: new Set(), // 存储展开的节点ID
      nodeData: new Map(), // 存储节点数据缓存
      
      // 单元类型配置
      unitTypes: [
        { value: 'company', label: '公司', icon: 'el-icon-office-building', color: '#595959' },
        { value: 'business_unit', label: '业务单元', icon: 'el-icon-s-cooperation', color: '#409eff' },
        { value: 'department', label: '部门', icon: 'el-icon-user', color: '#52c41a' },
        { value: 'team', label: '团队', icon: 'el-icon-s-custom', color: '#fa8c16' },
        { value: 'cost_center', label: '成本中心', icon: 'el-icon-s-finance', color: '#722ed1' },
        { value: 'location', label: '地理位置', icon: 'el-icon-location', color: '#ff4d4f' },
        { value: 'project', label: '项目组', icon: 'el-icon-s-flag', color: '#eb2f96' }
      ],
      
      currentTransform: d3.zoomIdentity
    }
  },
  emits: ['unit-selected', 'create-department', 'create-employee', 'view-type-change', 'view-employee', 'edit-employee'],
  computed: {
    legendItems() {
      // 只显示当前数据中存在的单元类型
      const existingTypes = new Set()
      const collectTypes = (nodes) => {
        nodes.forEach(node => {
          existingTypes.add(node.type)
          if (node.children && node.children.length > 0) {
            collectTypes(node.children)
          }
        })
      }
      collectTypes(this.organizationData)
      
      return this.unitTypes.filter(type => existingTypes.has(type.value))
    }
  },
  watch: {
    organizationData: {
      handler(newData) {
        if (newData && newData.length > 0) {
          this.initChart()
        }
      },
      immediate: true,
      deep: true
    },
    selectedUnit(newUnit) {
      if (newUnit && this.root) {
        this.highlightSelectedNode(newUnit.id)
      }
    },
    viewType() {
      this.initChart()
    }
  },
  mounted() {
    this.initSvg()
    if (this.organizationData && this.organizationData.length > 0) {
      this.initChart()
    }
  },
  methods: {
    initSvg() {
      const container = this.$refs.chartContainer
      const svg = d3.select(this.$refs.orgChartSvg)
      
      // 设置SVG尺寸
      const containerRect = container.getBoundingClientRect()
      svg.attr('width', containerRect.width)
         .attr('height', containerRect.height)
      
      // 初始化缩放
      this.zoom = d3.zoom()
        .scaleExtent([0.1, 3])
        .on('zoom', (event) => {
          this.currentTransform = event.transform
          svg.select('.chart-content').attr('transform', event.transform)
        })
      
      svg.call(this.zoom)
      
      this.svg = svg
      this.g = svg.select('.chart-content')
    },
    
    initChart() {
      if (!this.organizationData || this.organizationData.length === 0) return
      
      // 清空现有内容
      this.g.selectAll('*').remove()
      
      // 缓存所有节点数据
      this.cacheNodeData(this.organizationData[0])
      
      // 初始化时展开前两层
      if (this.expandedNodes.size === 0) {
        this.initializeDefaultExpansion()
      }
      
      // 根据视图类型渲染
      switch (this.viewType) {
        case 'tree':
          this.renderTreeChart()
          break
        case 'radial':
          this.renderRadialChart()
          break
        case 'horizontal':
          this.renderHorizontalChart()
          break
      }
    },
    
    renderTreeChart() {
      const width = this.svg.attr('width')
      const height = this.svg.attr('height')
      
      // 重新构建展开状态的层次结构
      this.root = this.buildExpandedHierarchy()
      
      const treeLayout = d3.tree()
        .size([width - 200, height - 200])
        .separation((a, b) => a.parent === b.parent ? 1 : 1.2)
      
      treeLayout(this.root)
      
      // 调整根节点位置
      this.root.descendants().forEach(d => {
        d.y += 100
        d.x += width / 2 - 100
      })
      
      this.drawNodes()
      this.drawLinks()
    },
    
    buildExpandedHierarchy() {
      const buildNode = (nodeData) => {
        const node = {
          ...nodeData,
          children: []
        }
        
        if (this.isExpanded(nodeData.id)) {
          const cachedData = this.nodeData.get(nodeData.id)
          const allChildren = []
          
          // 添加子部门
          if (cachedData && cachedData.children) {
            const departmentChildren = cachedData.children
              .filter(child => this.isExpanded(child.id))
              .map(child => buildNode(child))
            allChildren.push(...departmentChildren)
          }
          
          // 添加直接员工（作为叶子节点）
          if (cachedData && cachedData.directEmployees && cachedData.directEmployees.length > 0 && this.isExpanded(`${nodeData.id}_employees`)) {
            const employeeNodes = cachedData.directEmployees.map(emp => ({
              id: `emp_${emp.id}`,
              name: emp.name,
              type: 'employee',
              employeeData: emp,
              level: (nodeData.level || 1) + 1,
              children: [] // 员工是叶子节点
            }))
            allChildren.push(...employeeNodes)
          }
          
          node.children = allChildren
        }
        
        return node
      }
      
      const hierarchyData = buildNode(this.organizationData[0])
      return d3.hierarchy(hierarchyData, d => d.children)
    },
    
    renderRadialChart() {
      const width = this.svg.attr('width')
      const height = this.svg.attr('height')
      const radius = Math.min(width, height) / 2 - 100
      
      // 重新构建展开状态的层次结构
      this.root = this.buildExpandedHierarchy()
      
      const treeLayout = d3.tree()
        .size([2 * Math.PI, radius])
        .separation((a, b) => (a.parent === b.parent ? 1 : 2) / a.depth)
      
      treeLayout(this.root)
      
      // 转换为笛卡尔坐标
      this.root.descendants().forEach(d => {
        const angle = d.x
        const r = d.y
        d.x = r * Math.cos(angle - Math.PI / 2) + width / 2
        d.y = r * Math.sin(angle - Math.PI / 2) + height / 2
      })
      
      this.drawNodes()
      this.drawRadialLinks()
    },
    
    renderHorizontalChart() {
      const width = this.svg.attr('width')
      const height = this.svg.attr('height')
      
      // 重新构建展开状态的层次结构
      this.root = this.buildExpandedHierarchy()
      
      const treeLayout = d3.tree()
        .size([height - 200, width - 400])
        .separation((a, b) => a.parent === b.parent ? 1 : 1.2)
      
      treeLayout(this.root)
      
      // 调整为水平布局
      this.root.descendants().forEach(d => {
        const temp = d.x
        d.x = d.y + 150
        d.y = temp + 100
      })
      
      this.drawNodes()
      this.drawLinks()
    },
    
    drawNodes() {
      const nodes = this.g.selectAll('.node')
        .data(this.root.descendants())
        .enter().append('g')
        .attr('class', 'node')
        .attr('transform', d => `translate(${d.x},${d.y})`)
        .style('cursor', 'pointer')
        .on('click', (_, d) => this.handleNodeClick(d))
        .on('mouseenter', (event, d) => this.showTooltip(event, d))
        .on('mouseleave', () => this.hideTooltip())
      
      // 绘制节点背景
      nodes.append('rect')
        .attr('class', 'node-rect')
        .attr('x', -75)
        .attr('y', -25)
        .attr('width', 150)
        .attr('height', 50)
        .attr('rx', 8)
        .attr('ry', 8)
        .style('fill', d => d.data.type === 'employee' ? this.getEmployeeNodeColor() : this.getNodeColor(d.data.type))
        .style('stroke', d => d.data.id === this.selectedUnit?.id ? '#409eff' : '#e4e7ed')
        .style('stroke-width', d => d.data.id === this.selectedUnit?.id ? 3 : 1)
        .style('filter', 'url(#shadow)')
      
      // 绘制图标
      nodes.append('foreignObject')
        .attr('x', -70)
        .attr('y', -20)
        .attr('width', 20)
        .attr('height', 20)
        .append('xhtml:div')
        .style('font-size', '16px')
        .style('color', d => d.data.type === 'employee' ? '#409eff' : this.getIconColor(d.data.type))
        .html(d => `<i class="${d.data.type === 'employee' ? 'el-icon-user' : this.getUnitTypeIcon(d.data.type)}"></i>`)
      
      // 绘制主要文本
      nodes.append('text')
        .attr('class', 'node-name')
        .attr('x', -45)
        .attr('y', -5)
        .style('font-size', '12px')
        .style('font-weight', 'bold')
        .style('fill', '#303133')
        .text(d => this.truncateText(d.data.name, 12))
      
      // 绘制负责人信息（包含职位）
      nodes.append('text')
        .attr('class', 'node-manager')
        .attr('x', -45)
        .attr('y', 8)
        .style('font-size', '10px')
        .style('fill', '#67c23a')
        .text(d => {
          if (d.data.type === 'employee') {
            // 员工节点显示职位
            const emp = d.data.employeeData
            return emp && emp.position ? this.truncateText(emp.position.name, 10) : '未设置职位'
          } else {
            // 部门节点显示负责人和职位
            const manager = d.data.manager
            if (!manager) return '无负责人'
            
            const managerName = this.truncateText(manager.name, 8)
            const position = manager.position ? manager.position.name : ''
            return position ? `${managerName} | ${this.truncateText(position, 8)}` : managerName
          }
        })
      
      // 绘制展开/收起按钮（对于有子节点的节点，但不包括员工节点）
      const expandButton = nodes.filter(d => d.data.type !== 'employee' && this.hasChildren(d.data))
        .append('g')
        .attr('class', 'expand-button')
        .style('cursor', 'pointer')
        .on('click', (event, d) => {
          event.stopPropagation()
          this.toggleNode(d)
        })
      
      expandButton.append('circle')
        .attr('cx', 65)
        .attr('cy', 0)
        .attr('r', 10)
        .style('fill', d => this.isExpanded(d.data.id) ? '#f56c6c' : '#67c23a')
        .style('stroke', '#fff')
        .style('stroke-width', 2)
      
      expandButton.append('text')
        .attr('x', 65)
        .attr('y', 4)
        .style('font-size', '12px')
        .style('font-weight', 'bold')
        .style('fill', '#fff')
        .style('text-anchor', 'middle')
        .text(d => this.isExpanded(d.data.id) ? '-' : '+')
      
      // 绘制员工数量标识（只对部门节点显示）
      const departmentNodes = nodes.filter(d => d.data.type !== 'employee')
      
      departmentNodes.append('circle')
        .attr('class', 'employee-badge')
        .attr('cx', 65)
        .attr('cy', -15)
        .attr('r', 8)
        .style('fill', '#409eff')
        .style('stroke', '#fff')
        .style('stroke-width', 1)
      
      departmentNodes.append('text')
        .attr('class', 'employee-count')
        .attr('x', 65)
        .attr('y', -11)
        .style('font-size', '8px')
        .style('font-weight', 'bold')
        .style('fill', '#fff')
        .style('text-anchor', 'middle')
        .text(d => d.data.employeeCount || 0)
    },
    
    drawLinks() {
      this.g.selectAll('.link')
        .data(this.root.links())
        .enter().insert('path', '.node')
        .attr('class', 'link')
        .attr('d', d => {
          // 动态计算连接线路径，考虑展开按钮位置
          const sourceY = d.source.y + (this.hasChildren(d.source.data) ? 10 : 25)
          const targetY = d.target.y - 25
          return `M${d.source.x},${sourceY}
                  C${d.source.x},${(sourceY + targetY) / 2}
                   ${d.target.x},${(sourceY + targetY) / 2}
                   ${d.target.x},${targetY}`
        })
        .style('fill', 'none')
        .style('stroke', '#409eff')
        .style('stroke-width', 2)
        .style('opacity', 0.7)
    },
    
    drawRadialLinks() {
      this.g.selectAll('.link')
        .data(this.root.links())
        .enter().insert('path', '.node')
        .attr('class', 'link')
        .attr('d', d3.linkRadial()
          .angle(d => d.x)
          .radius(d => d.y))
        .style('fill', 'none')
        .style('stroke', '#409eff')
        .style('stroke-width', 2)
        .style('opacity', 0.7)
    },
    
    showTooltip(event, d) {
      this.hoveredNode = d
      this.keepTooltip = false
      
      const containerRect = this.$refs.chartContainer.getBoundingClientRect()
      const mouseX = event.clientX - containerRect.left
      const mouseY = event.clientY - containerRect.top
      
      this.tooltipStyle = {
        left: `${mouseX + 20}px`,
        top: `${mouseY - 10}px`,
        display: 'block'
      }
    },
    
    hideTooltip() {
      setTimeout(() => {
        if (!this.keepTooltip) {
          this.hoveredNode = null
          this.tooltipStyle.display = 'none'
        }
      }, 100)
    },
    
    selectNode(nodeData) {
      this.hideTooltip()
      this.$emit('unit-selected', nodeData)
    },
    
    selectNodeAndExpand(nodeData) {
      // 选中节点并智能展开
      this.selectNode(nodeData)
      
      // 获取节点的父级路径
      const parentPath = this.getNodePath(nodeData.id)
      
      // 收起所有不在路径上的兄弟节点
      this.collapseUnrelatedNodes(parentPath, nodeData.id)
      
      // 展开选中节点的子节点
      this.expandNodeChildren(nodeData.id)
      
      // 重新渲染
      this.updateChart()
    },
    
    highlightSelectedNode(nodeId) {
      this.g.selectAll('.node-rect')
        .style('stroke', d => d.data.id === nodeId ? '#409eff' : '#e4e7ed')
        .style('stroke-width', d => d.data.id === nodeId ? 3 : 1)
    },
    
    
    zoomIn() {
      this.svg.transition().duration(300).call(
        this.zoom.scaleBy, 1.5
      )
    },
    
    zoomOut() {
      this.svg.transition().duration(300).call(
        this.zoom.scaleBy, 1 / 1.5
      )
    },
    
    resetZoom() {
      this.svg.transition().duration(500).call(
        this.zoom.transform,
        d3.zoomIdentity
      )
    },
    
    exportChart() {
      // 导出SVG为图片
      const svgElement = this.$refs.orgChartSvg
      const serializer = new XMLSerializer()
      const svgString = serializer.serializeToString(svgElement)
      const canvas = document.createElement('canvas')
      const ctx = canvas.getContext('2d')
      const img = new Image()
      
      img.onload = () => {
        canvas.width = img.width
        canvas.height = img.height
        ctx.drawImage(img, 0, 0)
        
        // 下载图片
        const link = document.createElement('a')
        link.download = `组织架构图-${new Date().toISOString().slice(0, 10)}.png`
        link.href = canvas.toDataURL()
        link.click()
      }
      
      img.src = 'data:image/svg+xml;base64,' + btoa(decodeURIComponent(encodeURIComponent(svgString)))
    },
    
    // 工具方法
    getUnitTypeIcon(type) {
      const unitType = this.unitTypes.find(t => t.value === type)
      return unitType ? unitType.icon : 'el-icon-office-building'
    },
    
    getUnitTypeName(type) {
      const unitType = this.unitTypes.find(t => t.value === type)
      return unitType ? unitType.label : type
    },
    
    getNodeColor(type) {
      const baseColors = {
        company: '#f0f2f5',
        business_unit: '#e6f3ff',
        department: '#f6ffed',
        team: '#fff7e6',
        cost_center: '#f9f0ff',
        location: '#fff1f0',
        project: '#fff0f6'
      }
      return baseColors[type] || '#f5f7fa'
    },
    
    getEmployeeNodeColor() {
      return '#e8f4fd' // 员工节点使用淡蓝色背景
    },
    
    getIconColor(type) {
      const unitType = this.unitTypes.find(t => t.value === type)
      return unitType ? unitType.color : '#409eff'
    },
    
    truncateText(text, maxLength) {
      if (!text) return ''
      return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
    },
    
    // 初始化默认展开状态（展开前两层）
    initializeDefaultExpansion() {
      this.expandedNodes.clear()
      
      const expandToLevel = (nodeData, currentLevel = 1, maxLevel = 2) => {
        if (currentLevel <= maxLevel) {
          this.expandedNodes.add(nodeData.id)
          // 也为该节点添加员工显示标记
          this.expandedNodes.add(`${nodeData.id}_employees`)
          
          if (nodeData.children && nodeData.children.length > 0) {
            nodeData.children.forEach(child => {
              expandToLevel(child, currentLevel + 1, maxLevel)
            })
          }
        }
      }
      
      expandToLevel(this.organizationData[0])
    },
    
    cacheNodeData(nodeData) {
      // 递归缓存所有节点数据
      this.nodeData.set(nodeData.id, nodeData)
      if (nodeData.children && nodeData.children.length > 0) {
        nodeData.children.forEach(child => {
          this.cacheNodeData(child)
        })
      }
    },
    
    hasChildren(nodeData) {
      const cachedData = this.nodeData.get(nodeData.id)
      return cachedData && (
        (cachedData.children && cachedData.children.length > 0) ||
        (cachedData.directEmployees && cachedData.directEmployees.length > 0)
      )
    },
    
    isExpanded(nodeId) {
      return this.expandedNodes.has(nodeId)
    },
    
    handleNodeClick(d) {
      // 单击节点选择并展开其子节点
      this.selectNodeAndExpand(d.data)
    },
    
    toggleNode(d) {
      const nodeId = d.data.id
      
      if (this.isExpanded(nodeId)) {
        // 收起节点
        this.expandedNodes.delete(nodeId)
        // 移除子节点
        if (d.children) {
          d._children = d.children
          d.children = null
        }
      } else {
        // 展开节点
        this.expandedNodes.add(nodeId)
        // 加载子节点
        this.loadChildNodes(d)
      }
      
      this.updateChart()
    },
    
    loadChildNodes(d) {
      const cachedData = this.nodeData.get(d.data.id)
      if (cachedData && cachedData.children) {
        // 从缓存中加载子节点
        if (d._children) {
          d.children = d._children
          d._children = null
        } else {
          // 创建新的子节点
          d.children = cachedData.children.map(childData => {
            const childNode = d3.hierarchy(childData, d => {
              // 子节点初始不展开
              return this.isExpanded(d.id) ? d.children : []
            })
            return childNode
          })
        }
      }
    },
    
    updateChart() {
      // 重新计算布局和重绘
      this.initChart()
    },
    
    expandAll() {
      // 展开所有节点（包括员工）
      this.expandedNodes.clear()
      this.nodeData.forEach((nodeData, nodeId) => {
        this.expandedNodes.add(nodeId)
        // 也展开该节点的员工
        if (nodeData.directEmployees && nodeData.directEmployees.length > 0) {
          this.expandedNodes.add(`${nodeId}_employees`)
        }
      })
      this.updateChart()
    },
    
    collapseAll() {
      // 收起所有节点，只保留根节点
      this.expandedNodes.clear()
      if (this.organizationData && this.organizationData.length > 0) {
        this.expandedNodes.add(this.organizationData[0].id)
        // 保留根节点的员工显示
        this.expandedNodes.add(`${this.organizationData[0].id}_employees`)
      }
      this.updateChart()
    },
    
    handleNodeAction(command) {
      const { action, data } = command
      this.hideTooltip()
      
      switch (action) {
        case 'createDepartment':
          this.$emit('create-department', data)
          break
        case 'createEmployee':
          this.$emit('create-employee', data)
          break
      }
    },
    
    async viewEmployeeDetail(employeeData) {
      this.hideTooltip()
      // 验证员工数据是否完整
      if (!employeeData || !employeeData.id) {
        this.$message.error('员工信息不完整')
        return
      }
      this.$emit('view-employee', employeeData)
    },
    
    async editEmployee(employeeData) {
      this.hideTooltip()
      // 验证员工数据是否完整
      if (!employeeData || !employeeData.id) {
        this.$message.error('员工信息不完整')
        return
      }
      this.$emit('edit-employee', employeeData)
    },
    
    // 获取节点的路径（从根到该节点）
    getNodePath(nodeId) {
      const path = []
      
      const findPath = (nodeData, targetId, currentPath = []) => {
        currentPath.push(nodeData.id)
        
        if (nodeData.id === targetId) {
          path.splice(0, path.length, ...currentPath)
          return true
        }
        
        if (nodeData.children) {
          for (const child of nodeData.children) {
            if (findPath(child, targetId, [...currentPath])) {
              return true
            }
          }
        }
        
        return false
      }
      
      findPath(this.organizationData[0], nodeId)
      return path
    },
    
    // 收起与选中节点无关的节点
    collapseUnrelatedNodes(nodePath, selectedNodeId) {
      const pathSet = new Set(nodePath)
      
      // 保留路径上的节点和选中节点的直接子节点
      const nodesToKeep = new Set(pathSet)
      
      // 为路径上的每个节点添加员工显示标记
      pathSet.forEach(nodeId => {
        nodesToKeep.add(`${nodeId}_employees`)
      })
      
      // 添加选中节点的直接子节点
      const selectedNodeData = this.nodeData.get(selectedNodeId)
      if (selectedNodeData && selectedNodeData.children) {
        selectedNodeData.children.forEach(child => {
          nodesToKeep.add(child.id)
          nodesToKeep.add(`${child.id}_employees`)
        })
      }
      
      // 清理expandedNodes，只保留需要的节点
      const newExpandedNodes = new Set()
      this.expandedNodes.forEach(nodeId => {
        if (nodesToKeep.has(nodeId)) {
          newExpandedNodes.add(nodeId)
        }
      })
      
      this.expandedNodes = newExpandedNodes
    },
    
    // 展开节点的子节点
    expandNodeChildren(nodeId) {
      const nodeData = this.nodeData.get(nodeId)
      if (nodeData) {
        this.expandedNodes.add(nodeId)
        // 显示该节点的员工
        this.expandedNodes.add(`${nodeId}_employees`)
        
        // 展开子部门
        if (nodeData.children) {
          nodeData.children.forEach(child => {
            this.expandedNodes.add(child.id)
            // 也显示子部门的员工
            this.expandedNodes.add(`${child.id}_employees`)
          })
        }
      }
    }
  }
}
</script>

<style scoped>
.workday-org-chart {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #fafafa;
  border-radius: 8px;
  overflow: hidden;
}

.chart-header {
  padding: 16px 20px;
  background: white;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-title h3 {
  margin: 0 0 4px 0;
  color: #303133;
  font-size: 16px;
  font-weight: 600;
}

.subtitle {
  color: #909399;
  font-size: 12px;
}

.chart-controls {
  display: flex;
  gap: 12px;
  align-items: center;
}

.zoom-controls {
  margin-left: 8px;
}

.chart-container {
  flex: 1;
  position: relative;
  overflow: hidden;
  background: white;
}

.org-chart-svg {
  width: 100%;
  height: 100%;
  cursor: grab;
}

.org-chart-svg:active {
  cursor: grabbing;
}

/* 节点样式 */
.org-chart-svg :deep(.node) {
  transition: all 0.3s ease;
}

.org-chart-svg :deep(.node:hover .node-rect) {
  transform: scale(1.05);
  filter: url(#shadow) brightness(1.1);
}

.org-chart-svg :deep(.link) {
  transition: all 0.3s ease;
}

.org-chart-svg :deep(.link:hover) {
  stroke-width: 3;
  opacity: 1;
}

/* 工具提示样式 */
.node-tooltip {
  position: absolute;
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  padding: 0;
  z-index: 1000;
  max-width: 280px;
  overflow: hidden;
}

.tooltip-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #f8f9fa;
  border-bottom: 1px solid #e4e7ed;
}

.node-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  background: #409eff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.node-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-avatar {
  color: white;
  font-size: 20px;
}

.node-info .node-name {
  font-weight: 600;
  color: #303133;
  font-size: 14px;
}

.node-info .node-manager {
  color: #67c23a;
  font-size: 12px;
  margin-top: 2px;
}

.tooltip-content {
  padding: 12px 16px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 12px;
}

.info-row:last-child {
  margin-bottom: 0;
}

.label {
  color: #909399;
}

.value {
  color: #303133;
  font-weight: 500;
}

.tooltip-actions {
  padding: 8px 16px 12px;
  border-top: 1px solid #f0f0f0;
}

/* 图例样式 */
.chart-legend {
  padding: 12px 20px;
  background: white;
  border-top: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  gap: 20px;
}

.legend-title {
  font-weight: 600;
  color: #303133;
  font-size: 12px;
}

.legend-items {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #606266;
}

.legend-item i {
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .chart-header {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }
  
  .chart-controls {
    width: 100%;
    justify-content: space-between;
  }
  
  .legend-items {
    gap: 8px;
  }
  
  .node-tooltip {
    max-width: 250px;
  }
}
</style>