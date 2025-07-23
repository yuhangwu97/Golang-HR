<template>
  <div class="department-tree-select">
    <el-select
      :value="selectedLabel"
      :placeholder="placeholder"
      :disabled="disabled"
      :clearable="clearable"
      :filterable="filterable"
      :filter-method="filterMethod"
      @clear="handleClear"
      @focus="handleFocus"
      @blur="handleBlur"
      ref="select"
      popper-class="department-tree-select-dropdown"
    >
      <el-option :value="selectedValue" :label="selectedLabel" style="display: none;" />
      <div class="tree-select-dropdown">
        <div v-if="showSearch" class="tree-search">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索部门..."
            prefix-icon="el-icon-search"
            size="small"
            @input="handleSearch"
          />
        </div>
        <div class="tree-container">
          <el-tree
            ref="tree"
            :data="treeData"
            :props="treeProps"
            :node-key="nodeKey"
            :expand-on-click-node="false"
            :check-on-click-node="true"
            :filter-node-method="filterNode"
            :default-expand-all="defaultExpandAll"
            :highlight-current="true"
            @node-click="handleNodeClick"
            @current-change="handleCurrentChange"
            class="department-tree"
          >
            <div class="tree-node" slot-scope="{ node, data }">
              <div class="node-content">
                <div class="node-info">
                  <i :class="getNodeIcon(data)" class="node-icon"></i>
                  <span class="node-name">{{ data.name }}</span>
                  <span v-if="showCode && data.code" class="node-code">{{ data.code }}</span>
                </div>
                <div class="node-meta">
                  <el-tag v-if="showType" :type="getTypeColor(data.type)" size="mini">
                    {{ getTypeName(data.type) }}
                  </el-tag>
                  <span v-if="showEmployeeCount" class="employee-count" :title="`直接员工: ${data.employeeCount || 0}人, 含子部门: ${data.hierarchicalEmployeeCount || data.employeeCount || 0}人`">
                    {{ data.hierarchicalEmployeeCount || data.employeeCount || 0 }}人
                  </span>
                </div>
              </div>
            </div>
          </el-tree>
        </div>
        <div v-if="!treeData.length" class="empty-state">
          <i class="el-icon-folder-opened"></i>
          <p>暂无部门数据</p>
        </div>
      </div>
    </el-select>
  </div>
</template>

<script>
export default {
  name: 'DepartmentTreeSelect',
  props: {
    // 选中的值
    value: {
      type: [String, Number],
      default: null
    },
    // 占位符
    placeholder: {
      type: String,
      default: '请选择部门'
    },
    // 是否禁用
    disabled: {
      type: Boolean,
      default: false
    },
    // 是否可清空
    clearable: {
      type: Boolean,
      default: true
    },
    // 是否可搜索
    filterable: {
      type: Boolean,
      default: true
    },
    // 树形数据
    treeData: {
      type: Array,
      default: () => []
    },
    // 是否显示搜索框
    showSearch: {
      type: Boolean,
      default: true
    },
    // 是否显示部门编码
    showCode: {
      type: Boolean,
      default: true
    },
    // 是否显示部门类型
    showType: {
      type: Boolean,
      default: false
    },
    // 是否显示员工数量
    showEmployeeCount: {
      type: Boolean,
      default: false
    },
    // 是否默认展开所有节点
    defaultExpandAll: {
      type: Boolean,
      default: false
    },
    // 节点键值
    nodeKey: {
      type: String,
      default: 'id'
    },
    // 只能选择叶子节点
    leafOnly: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      searchKeyword: '',
      selectedValue: this.value,
      selectedLabel: '',
      selectedNode: null,
      treeProps: {
        children: 'children',
        label: 'name'
      }
    }
  },
  computed: {
    // 扁平化的部门数据（用于搜索）
    flattenData() {
      const result = []
      const flatten = (nodes, level = 0) => {
        nodes.forEach(node => {
          result.push({
            ...node,
            level
          })
          if (node.children && node.children.length) {
            flatten(node.children, level + 1)
          }
        })
      }
      flatten(this.treeData)
      return result
    }
  },
  watch: {
    value(newVal) {
      this.selectedValue = newVal
      this.updateSelectedLabel()
    },
    treeData: {
      handler() {
        this.updateSelectedLabel()
      },
      immediate: true
    }
  },
  methods: {
    // 获取节点图标
    getNodeIcon(data) {
      const iconMap = {
        'company': 'el-icon-office-building',
        'business_unit': 'el-icon-s-cooperation',
        'department': 'el-icon-user',
        'team': 'el-icon-s-custom',
        'cost_center': 'el-icon-s-finance',
        'location': 'el-icon-location',
        'project': 'el-icon-s-flag'
      }
      return iconMap[data.type] || 'el-icon-office-building'
    },
    
    // 获取类型颜色
    getTypeColor(type) {
      const colorMap = {
        'company': '',
        'business_unit': 'primary',
        'department': 'success',
        'team': 'warning',
        'cost_center': 'info',
        'location': 'danger',
        'project': 'warning'
      }
      return colorMap[type] || ''
    },
    
    // 获取类型名称
    getTypeName(type) {
      const nameMap = {
        'company': '公司',
        'business_unit': '事业部',
        'department': '部门',
        'team': '团队',
        'cost_center': '成本中心',
        'location': '地点',
        'project': '项目'
      }
      return nameMap[type] || type
    },
    
    // 树节点过滤
    filterNode(value, data) {
      if (!value) return true
      return data.name.includes(value) || 
             (data.code && data.code.includes(value))
    },
    
    // 搜索处理
    handleSearch() {
      this.$refs.tree.filter(this.searchKeyword)
    },
    
    // 过滤方法
    filterMethod(value) {
      this.searchKeyword = value
      this.handleSearch()
    },
    
    // 节点点击
    handleNodeClick(data, node) {
      // 如果设置了只能选择叶子节点
      if (this.leafOnly && node.childNodes.length > 0) {
        return
      }
      
      this.selectedValue = data[this.nodeKey]
      this.selectedLabel = data.name
      this.selectedNode = data
      
      // 关闭下拉框
      this.$refs.select.blur()
      
      // 触发change事件
      this.$emit('input', this.selectedValue)
      this.$emit('change', this.selectedValue, data)
    },
    
    // 当前节点变化
    handleCurrentChange(data) {
      if (data) {
        this.selectedNode = data
      }
    },
    
    // 清空
    handleClear() {
      this.selectedValue = null
      this.selectedLabel = ''
      this.selectedNode = null
      this.$refs.tree.setCurrentKey(null)
      this.$emit('input', null)
      this.$emit('change', null, null)
    },
    
    // 获取焦点
    handleFocus() {
      this.$emit('focus')
    },
    
    // 失去焦点
    handleBlur() {
      this.$emit('blur')
    },
    
    // 更新选中标签
    updateSelectedLabel() {
      if (this.selectedValue && this.flattenData.length) {
        const node = this.flattenData.find(item => item[this.nodeKey] === this.selectedValue)
        if (node) {
          this.selectedLabel = node.name
          this.selectedNode = node
          this.$nextTick(() => {
            this.$refs.tree.setCurrentKey(this.selectedValue)
          })
        }
      } else {
        this.selectedLabel = ''
        this.selectedNode = null
      }
    },
    
    // 获取选中的节点数据
    getSelectedNode() {
      return this.selectedNode
    },
    
    // 设置选中的节点
    setSelectedNode(key) {
      this.selectedValue = key
      this.updateSelectedLabel()
    }
  }
}
</script>

<style scoped>
.department-tree-select {
  width: 100%;
}

.tree-select-dropdown {
  padding: 0;
  max-height: 400px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.tree-search {
  padding: 8px 12px;
  border-bottom: 1px solid #e4e7ed;
  background: #fafafa;
}

.tree-container {
  flex: 1;
  overflow: auto;
  padding: 8px 0;
}

.department-tree {
  background: transparent;
  border: none;
}

.department-tree :deep(.el-tree-node__content) {
  padding: 8px 12px;
  height: auto;
  min-height: 36px;
}

.department-tree :deep(.el-tree-node__content:hover) {
  background-color: #f0f9ff;
}

.department-tree :deep(.el-tree-node.is-current > .el-tree-node__content) {
  background-color: #e1f3ff;
  color: #1890ff;
}

.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 4px 0;
}

.node-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.node-info {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.node-icon {
  font-size: 14px;
  color: #666;
}

.node-name {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
}

.node-code {
  font-size: 12px;
  color: #909399;
  margin-left: 8px;
}

.node-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.employee-count {
  font-size: 12px;
  color: #909399;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: #909399;
}

.empty-state i {
  font-size: 48px;
  margin-bottom: 16px;
  display: block;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}
</style>

<style>
.department-tree-select-dropdown {
  padding: 0 !important;
}

.department-tree-select-dropdown .el-select-dropdown__item {
  display: none;
}
</style>