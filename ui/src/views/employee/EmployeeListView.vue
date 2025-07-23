<template>
  <div class="employee-list">
    <!-- 简化的页面头部 -->
    <div class="page-header-simple">
      <div class="page-title">
        <div class="title-icon">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 12C14.21 12 16 10.21 16 8C16 5.79 14.21 4 12 4C9.79 4 8 5.79 8 8C8 10.21 9.79 12 12 12ZM12 14C9.33 14 4 15.34 4 18V20H20V18C20 15.34 14.67 14 12 14Z" fill="currentColor"/>
          </svg>
        </div>
        <div class="title-content">
          <h1>员工管理</h1>
        </div>
      </div>
    </div>

    <!-- 浮动操作按钮 -->
    <div class="floating-actions">
      <el-dropdown trigger="click" @command="handleFloatingAction">
        <el-button type="primary" class="floating-trigger" circle>
          <i class="el-icon-plus"></i>
        </el-button>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item command="add">
            <i class="el-icon-plus"></i>
            添加员工
          </el-dropdown-item>
          <el-dropdown-item command="import" divided>
            <i class="el-icon-upload2"></i>
            导入员工
          </el-dropdown-item>
          <el-dropdown-item command="export">
            <i class="el-icon-download"></i>
            导出员工
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>

    <!-- 搜索筛选区域 -->
    <el-card class="search-card fade-in">
      <div class="search-header">
        <h3>筛选条件</h3>
        <el-button type="text" @click="toggleSearchExpand" class="expand-btn">
          <i :class="searchExpanded ? 'el-icon-arrow-up' : 'el-icon-arrow-down'"></i>
          {{ searchExpanded ? '收起' : '展开' }}
        </el-button>
      </div>
      <transition name="slide-fade">
        <div v-show="searchExpanded" class="search-content">
          <el-form :model="searchForm" inline @submit.native.prevent="handleSearch">
            <div class="search-row">
              <el-form-item label="关键词" class="search-item">
                <el-input
                  v-model="searchForm.keyword"
                  placeholder="姓名、邮箱、工号"
                  style="width: 200px"
                  clearable
                  prefix-icon="el-icon-search"
                />
              </el-form-item>
              
              <el-form-item label="部门" class="search-item">
                <DepartmentTreeSelect
                  v-model="searchForm.department_id"
                  :tree-data="departmentTree"
                  placeholder="请选择部门"
                  style="width: 200px"
                  @change="handleDepartmentFilterChange"
                  :show-code="true"
                  :show-employee-count="true"
                  :default-expand-all="false"
                  :clearable="true"
                />
                <!-- Debug info -->
                <div v-if="!isDepartmentTreeLoaded" style="font-size: 12px; color: #999; margin-top: 4px;">
                  部门数据加载中...
                </div>
                <div v-else style="font-size: 12px; color: #67c23a; margin-top: 4px;">
                  已加载 {{ departmentTree.length }} 个部门
                </div>
              </el-form-item>
              
              <el-form-item label="状态" class="search-item">
                <el-select
                  v-model="searchForm.status"
                  placeholder="选择状态"
                  style="width: 120px"
                  clearable
                >
                  <el-option label="全部" value="" />
                  <el-option label="在职" value="active" />
                  <el-option label="离职" value="inactive" />
                </el-select>
              </el-form-item>
            </div>
            
            <div class="search-actions">
              <el-button type="primary" @click="handleSearch" class="search-btn">
                <i class="el-icon-search"></i>
                搜索
              </el-button>
              <el-button @click="resetSearch" class="reset-btn">
                <i class="el-icon-refresh"></i>
                重置
              </el-button>
            </div>
          </el-form>
        </div>
      </transition>
    </el-card>

    <!-- 员工列表 -->
    <el-card>
      <div slot="header" class="table-header">
        <span>员工列表</span>
        <div class="table-actions">
          <span class="count-info">
            共 {{ pagination.total }} 人
          </span>
          <el-button-group>
            <el-button 
              size="small" 
              :type="displayMode === 'table' ? 'primary' : 'default'"
              @click="displayMode = 'table'"
            >
              <i class="el-icon-menu"></i>
              表格
            </el-button>
            <el-button 
              size="small" 
              :type="displayMode === 'cards' ? 'primary' : 'default'"
              @click="displayMode = 'cards'"
            >
              <i class="el-icon-postcard"></i>
              卡片
            </el-button>
          </el-button-group>
          <el-button
            size="small"
            :disabled="!selectedRowKeys.length"
            @click="handleBulkDelete"
          >
            批量删除
          </el-button>
        </div>
      </div>

      <!-- 表格视图 -->
      <el-table
        v-if="displayMode === 'table'"
        :data="employees"
        v-loading="loading"
        @selection-change="handleSelectionChange"
        row-key="id"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column label="头像" width="60">
          <template slot-scope="scope">
            <el-avatar>
              {{ (scope.row.name && scope.row.name.charAt(0).toUpperCase()) || '' }}
            </el-avatar>
          </template>
        </el-table-column>
        
        <el-table-column label="姓名" width="150">
          <template slot-scope="scope">
            <div class="name-cell">
              <div class="name">{{ scope.row.name }}</div>
              <div class="employee-id">{{ scope.row.employee_id }}</div>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="phone" label="手机号" />
        
        <el-table-column label="部门">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.department" type="primary" size="small">
              {{ getDepartmentHierarchyPath(scope.row.department) }}
            </el-tag>
            <span v-else class="text-gray">未分配</span>
          </template>
        </el-table-column>
        
        <el-table-column label="职位">
          <template slot-scope="scope">
            <span v-if="scope.row.position">{{ scope.row.position.name }}</span>
            <span v-else class="text-gray">未设置</span>
          </template>
        </el-table-column>
        
        <el-table-column label="直属领导">
          <template slot-scope="scope">
            <span v-if="scope.row.manager">{{ scope.row.manager.name }}</span>
            <span v-else class="text-gray">无</span>
          </template>
        </el-table-column>
        
        <el-table-column label="状态" width="80">
          <template slot-scope="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ scope.row.status === 'active' ? '在职' : '离职' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="入职日期" width="120">
          <template slot-scope="scope">
            <span v-if="scope.row.hire_date">
              {{ formatDate(scope.row.hire_date) }}
            </span>
            <span v-else class="text-gray">-</span>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="180">
          <template slot-scope="scope">
            <el-button
              type="text"
              size="small"
              @click="viewEmployeeDetails(scope.row)"
              icon="el-icon-view"
            >
              查看
            </el-button>
            <el-button
              type="text"
              size="small"
              @click="editEmployee(scope.row)"
              icon="el-icon-edit"
            >
              编辑
            </el-button>
            <el-popconfirm
              title="确定要删除这个员工吗？"
              @onConfirm="deleteEmployee(scope.row)"
            >
              <el-button 
                slot="reference" 
                type="text" 
                size="small" 
                style="color: #f56c6c;"
                icon="el-icon-delete"
              >
                删除
              </el-button>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 卡片视图 -->
      <div v-if="displayMode === 'cards'" class="cards-container" v-loading="loading">
        <div class="employee-cards">
          <div 
            v-for="employee in employees" 
            :key="employee.id" 
            class="employee-card"
            @click="viewEmployeeDetails(employee)"
          >
            <div class="card-header">
              <div class="employee-avatar">
                <el-avatar :size="50" :src="employee.avatar">
                  {{ (employee.name && employee.name.charAt(0).toUpperCase()) || '' }}
                </el-avatar>
                <div class="status-badge" :class="employee.status">
                  <i class="el-icon-success" v-if="employee.status === 'active'"></i>
                  <i class="el-icon-error" v-else></i>
                </div>
              </div>
              <div class="card-actions">
                <el-dropdown trigger="click" @command="handleCardCommand">
                  <span class="el-dropdown-link">
                    <i class="el-icon-more"></i>
                  </span>
                  <el-dropdown-menu slot="dropdown">
                    <el-dropdown-item :command="{action: 'view', employee: employee}">
                      <i class="el-icon-view"></i> 查看详情
                    </el-dropdown-item>
                    <el-dropdown-item :command="{action: 'edit', employee: employee}">
                      <i class="el-icon-edit"></i> 编辑
                    </el-dropdown-item>
                    <el-dropdown-item :command="{action: 'delete', employee: employee}" divided>
                      <i class="el-icon-delete"></i> 删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </el-dropdown>
              </div>
            </div>
            
            <div class="card-body">
              <h3 class="employee-name">{{ employee.name }}</h3>
              <p class="employee-id">工号: {{ employee.employee_id }}</p>
              
              <div class="employee-info">
                <div class="info-item">
                  <i class="el-icon-office-building"></i>
                  <span>{{ employee.department ? getDepartmentHierarchyPath(employee.department) : '未分配部门' }}</span>
                </div>
                <div class="info-item">
                  <i class="el-icon-suitcase"></i>
                  <span>{{ employee.position ? employee.position.name : '未设置职位' }}</span>
                </div>
                <div class="info-item">
                  <i class="el-icon-message"></i>
                  <span>{{ employee.email }}</span>
                </div>
                <div class="info-item" v-if="employee.phone">
                  <i class="el-icon-phone"></i>
                  <span>{{ employee.phone }}</span>
                </div>
              </div>
            </div>
            
            <div class="card-footer">
              <div class="hire-date" v-if="employee.hire_date">
                <i class="el-icon-date"></i>
                <span>入职: {{ formatDate(employee.hire_date) }}</span>
              </div>
              <el-tag 
                :type="employee.status === 'active' ? 'success' : 'danger'" 
                size="mini"
                class="status-tag"
              >
                {{ employee.status === 'active' ? '在职' : '离职' }}
              </el-tag>
            </div>
          </div>
        </div>
      </div>


      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          :current-page.sync="pagination.current"
          :page-size.sync="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 导入对话框 -->
    <el-dialog
      title="导入员工数据"
      :visible.sync="showImport"
      @close="showImport = false"
    >
      <el-upload
        :file-list="fileList"
        :before-upload="beforeUpload"
        accept=".xlsx,.xls,.csv"
        action=""
      >
        <el-button>
          <i class="el-icon-upload2"></i>
          选择文件
        </el-button>
      </el-upload>
      <div class="upload-tips">
        <p>支持文件格式：Excel (.xlsx, .xls) 或 CSV (.csv)</p>
        <p>请确保文件包含必要的列：姓名、邮箱等</p>
      </div>
      <div slot="footer" class="dialog-footer">
        <el-button @click="showImport = false">取消</el-button>
        <el-button type="primary" @click="handleImport">确定</el-button>
      </div>
    </el-dialog>

    <!-- 员工详情弹窗 -->
    <el-dialog
      title="员工详情"
      :visible.sync="showEmployeeDetails"
      width="800px"
      @close="showEmployeeDetails = false"
    >
      <div v-if="selectedEmployee" class="employee-detail-dialog">
        <el-row :gutter="24">
          <el-col :span="8">
            <div class="employee-avatar-section">
              <el-avatar :size="80" :src="selectedEmployee.avatar">
                {{ (selectedEmployee.name && selectedEmployee.name.charAt(0).toUpperCase()) || '' }}
              </el-avatar>
              <h3>{{ selectedEmployee.name }}</h3>
              <p>{{ selectedEmployee.position ? selectedEmployee.position.name : '未设置职位' }}</p>
              <el-tag :type="selectedEmployee.status === 'active' ? 'success' : 'danger'">
                {{ selectedEmployee.status === 'active' ? '在职' : '离职' }}
              </el-tag>
            </div>
          </el-col>
          <el-col :span="16">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="工号">{{ selectedEmployee.employee_id || '-' }}</el-descriptions-item>
              <el-descriptions-item label="邮箱">{{ selectedEmployee.email || '-' }}</el-descriptions-item>
              <el-descriptions-item label="手机号">{{ selectedEmployee.phone || '-' }}</el-descriptions-item>
              <el-descriptions-item label="部门">
                {{ selectedEmployee.department ? getDepartmentHierarchyPath(selectedEmployee.department) : '未分配' }}
              </el-descriptions-item>
              <el-descriptions-item label="职位">
                {{ selectedEmployee.position ? selectedEmployee.position.name : '未设置' }}
              </el-descriptions-item>
              <el-descriptions-item label="入职日期">
                {{ selectedEmployee.hire_date ? formatDate(selectedEmployee.hire_date) : '-' }}
              </el-descriptions-item>
            </el-descriptions>
          </el-col>
        </el-row>
      </div>
      <div slot="footer" class="dialog-footer">
        <el-button @click="showEmployeeDetails = false">关闭</el-button>
        <el-button type="primary" @click="editEmployee(selectedEmployee)" v-if="selectedEmployee">
          <i class="el-icon-edit"></i>
          编辑
        </el-button>
        <el-button @click="viewEmployee(selectedEmployee)" v-if="selectedEmployee">
          <i class="el-icon-view"></i>
          查看详情
        </el-button>
      </div>
    </el-dialog>

  </div>
</template>

<script>
import dayjs from 'dayjs'
import { employeeApiService } from '@/services/employeeApi'
import { departmentApi } from '@/services/departmentApi'
import DepartmentTreeSelect from '@/components/common/DepartmentTreeSelect.vue'

export default {
  name: 'EmployeeListView',
  components: {
    DepartmentTreeSelect
  },
  data() {
    return {
      loading: false,
      employees: [],
      departments: [],
      departmentTree: [],
      selectedRowKeys: [],
      showImport: false,
      fileList: [],
      showEmployeeDetails: false,
      selectedEmployee: null,
      displayMode: 'table', // 'table' | 'cards'
      searchExpanded: true,
      // 搜索表单
      searchForm: {
        keyword: '',
        department: '',
        department_id: null,
        status: '',
      },
      // 分页配置
      pagination: {
        current: 1,
        pageSize: 10,
        total: 0,
      }
    }
  },
  computed: {
    // 部门树是否已加载
    isDepartmentTreeLoaded() {
      return this.departmentTree && this.departmentTree.length > 0
    }
  },
  async mounted() {
    await Promise.all([
      this.fetchEmployees(),
      this.fetchDepartments()
    ])
  },
  methods: {
    // 获取部门层级路径（跳过第一层）
    getDepartmentHierarchyPath(department) {
      if (!department) return ''
      
      // 构建部门层级路径
      const buildPath = (dept) => {
        const path = []
        let current = dept
        
        // 向上遍历所有父级部门
        while (current) {
          if (current.name) {
            path.unshift(current.name)
          }
          current = current.parent
        }
        
        // 跳过第一层（根部门），如果有多于一层的话
        if (path.length > 1) {
          path.shift()
        }
        
        // 如果跳过第一层后还有内容，用 > 连接；否则显示原始名称
        return path.length > 0 ? path.join(' > ') : department.name
      }
      
      return buildPath(department)
    },
    
    // 处理浮动操作按钮
    handleFloatingAction(command) {
      switch (command) {
        case 'add':
          this.$router.push('/employees/create')
          break
        case 'import':
          this.showImport = true
          break
        case 'export':
          this.exportData()
          break
      }
    },
    // 格式化日期
    formatDate(date) {
      if (!date) return ''
      return dayjs(date).format('YYYY-MM-DD')
    },
    // 获取员工列表（分页）
    async fetchEmployees() {
      this.loading = true
      try {
        const params = {
          page: this.pagination.current,
          pageSize: this.pagination.pageSize,
          ...this.searchForm,
        }
        const response = await employeeApiService.getEmployees(params)
        
        // 处理API响应的数据结构
        const responseData = response.data || response
        this.employees = responseData.data || responseData.employees || []
        this.pagination.total = responseData.pagination?.total || responseData.total || 0
      } catch (error) {
        console.error('获取员工列表失败:', error)
        this.$message.error('获取员工列表失败: ' + (error.response?.data?.message || error.message || '未知错误'))
        this.employees = []
        this.pagination.total = 0
      } finally {
        this.loading = false
      }
    },
    // 搜索
    async handleSearch() {
      this.pagination.current = 1
      await this.fetchEmployees()
    },
    // 重置搜索
    async resetSearch() {
      Object.assign(this.searchForm, {
        keyword: '',
        department: '',
        department_id: null,
        status: '',
      })
      this.pagination.current = 1
      await this.fetchEmployees()
    },
    // 分页大小改变
    async handleSizeChange(size) {
      this.pagination.pageSize = size
      this.pagination.current = 1
      await this.fetchEmployees()
    },
    // 当前页改变
    async handleCurrentChange(page) {
      this.pagination.current = page
      await this.fetchEmployees()
    },
    // 选择变化
    handleSelectionChange(selection) {
      this.selectedRowKeys = selection.map(item => item.id)
    },
    // 查看员工
    viewEmployee(employee) {
      this.$router.push(`/employees/${employee.id}`)
    },
    // 查看员工详情弹窗
    viewEmployeeDetails(employee) {
      this.selectedEmployee = employee
      this.showEmployeeDetails = true
    },
    // 处理卡片操作命令
    handleCardCommand(command) {
      const { action, employee } = command
      switch (action) {
        case 'view':
          this.viewEmployeeDetails(employee)
          break
        case 'edit':
          this.editEmployee(employee)
          break
        case 'delete':
          this.$confirm('确定要删除这个员工吗？', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            this.deleteEmployee(employee)
          })
          break
      }
    },
    // 发送消息
    sendMessage() {
      this.$message.info('消息功能待开发')
    },
    // 获取性别文本
    getGenderText(gender) {
      const genderMap = {
        'male': '男',
        'female': '女'
      }
      return genderMap[gender] || '未知'
    },
    // 获取学历文本
    getEducationText(education) {
      const educationMap = {
        'primary': '小学',
        'junior': '初中',
        'senior': '高中',
        'associate': '大专',
        'bachelor': '本科',
        'master': '硕士',
        'doctorate': '博士'
      }
      return educationMap[education] || education
    },
    // 编辑员工
    editEmployee(employee) {
      this.$router.push(`/employees/${employee.id}/edit`)
    },
    // 删除员工
    async deleteEmployee(employee) {
      try {
        await employeeApiService.deleteEmployee(employee.id)
        this.$message.success('删除成功')
        await this.fetchEmployees()
      } catch (error) {
        console.error('删除员工失败:', error)
        this.$message.error('删除失败: ' + (error.response?.data?.message || error.message || '未知错误'))
      }
    },
    // 批量删除
    async handleBulkDelete() {
      if (!this.selectedRowKeys.length) {
        this.$message.warning('请选择要删除的员工')
        return
      }

      try {
        // 批量删除每个选中的员工
        await Promise.all(
          this.selectedRowKeys.map(id => employeeApiService.deleteEmployee(id))
        )
        this.$message.success('批量删除成功')
        this.selectedRowKeys = []
        await this.fetchEmployees()
      } catch (error) {
        console.error('批量删除失败:', error)
        this.$message.error('批量删除失败: ' + (error.response?.data?.message || error.message || '未知错误'))
      }
    },
    // 导出数据
    async exportData() {
      try {
        this.$message({
          message: '正在导出数据...',
          type: 'info'
        })
        // 使用employee.js中的导出方法
        const { employeeService } = await import('@/services/employee')
        const response = await employeeService.exportEmployees('excel')
        
        // 创建下载链接
        const url = window.URL.createObjectURL(new Blob([response.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', `employees_${new Date().toISOString().split('T')[0]}.xlsx`)
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
        
        this.$message.success('导出成功')
      } catch (error) {
        console.error('导出失败:', error)
        this.$message.error('导出失败: ' + (error.response?.data?.message || error.message || '未知错误'))
      }
    },
    // 导入数据
    async handleImport() {
      if (!this.fileList.length) {
        this.$message.warning('请选择要导入的文件')
        return
      }

      try {
        const file = this.fileList[0].raw || this.fileList[0]
        // 使用employee.js中的导入方法
        const { employeeService } = await import('@/services/employee')
        await employeeService.importEmployees(file)
        
        this.$message.success('导入成功')
        this.showImport = false
        this.fileList = []
        await this.fetchEmployees()
      } catch (error) {
        console.error('导入失败:', error)
        this.$message.error('导入失败: ' + (error.response?.data?.message || error.message || '未知错误'))
      }
    },
    // 文件上传前检查
    beforeUpload(file) {
      const isValidType = file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' ||
        file.type === 'application/vnd.ms-excel' ||
        file.type === 'text/csv'
      
      if (!isValidType) {
        this.$message.error('文件格式不正确')
        return false
      }
      
      const isLt10M = file.size / 1024 / 1024 < 10
      if (!isLt10M) {
        this.$message.error('文件大小不能超过 10MB')
        return false
      }
      
      return false // 阻止自动上传
    },
    // 获取部门列表
    async fetchDepartments() {
      try {
        const [departmentsRes, departmentTreeRes] = await Promise.all([
          departmentApi.getDepartments(),
          departmentApi.getDepartmentTree()
        ])
        
        // 处理部门列表响应
        const departmentsData = departmentsRes.data || departmentsRes
        this.departments = departmentsData.data || departmentsData.departments || []
        
        // 处理部门树响应
        const treeData = departmentTreeRes.data?.data || departmentTreeRes.data || []
        this.departmentTree = this.processDepartmentTree(treeData)
        
        console.log('部门树数据加载完成:', this.departmentTree)
      } catch (error) {
        console.error('获取部门列表失败:', error)
        // 如果API调用失败，使用备用数据以保证页面功能
        this.departments = [
          { id: 1, name: '技术部', parent_id: null },
          { id: 2, name: '产品部', parent_id: null },
          { id: 3, name: '市场部', parent_id: null },
          { id: 4, name: '前端组', parent_id: 1 },
          { id: 5, name: '后端组', parent_id: 1 }
        ]
        this.departmentTree = []
      }
    },
    
    // 处理部门树数据
    processDepartmentTree(treeData) {
      const processNode = (node) => {
        const children = []
        if (node.children && node.children.length > 0) {
          for (const child of node.children) {
            children.push(processNode(child))
          }
        }
        
        // 使用后端提供的真实数据
        const directEmployeeCount = node.employee_count || 0
        const hierarchicalEmployeeCount = this.calculateHierarchicalEmployeeCount(node)
        
        return {
          id: node.id,
          name: node.name,
          code: node.code,
          type: node.type || 'department',
          level: node.level || 1,
          description: node.description,
          manager_id: node.manager_id,
          employeeCount: directEmployeeCount,
          hierarchicalEmployeeCount: hierarchicalEmployeeCount,
          children: children
        }
      }
      
      const result = []
      for (const node of treeData) {
        result.push(processNode(node))
      }
      
      return result
    },
    
    // 计算层级员工数量（包括所有子部门的员工）
    calculateHierarchicalEmployeeCount(node) {
      let totalCount = node.employee_count || 0
      
      if (node.children && node.children.length > 0) {
        for (const child of node.children) {
          totalCount += this.calculateHierarchicalEmployeeCount(child)
        }
      }
      
      return totalCount
    },
    
    // 处理部门筛选变化
    handleDepartmentFilterChange(departmentId, departmentData) {
      // 当选择部门时，更新搜索表单并触发搜索
      this.searchForm.department_id = departmentId
      this.searchForm.department = departmentData ? departmentData.name : ''
      this.handleSearch()
    },
    
    // 获取总员工数
    getTotalEmployees() {
      return this.employees.length
    },
    // 切换搜索展开状态
    toggleSearchExpand() {
      this.searchExpanded = !this.searchExpanded
    }
  }
}
</script>

<style scoped>
/* 整体容器 */
.employee-list {
  padding: 0;
  position: relative;
}

/* 简化的页面头部 */
.page-header-simple {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
  padding: 16px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.2);
  position: relative;
  overflow: hidden;
}

/* 原始页面头部样式保留以防其他地方使用 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  padding: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.2);
  position: relative;
  overflow: hidden;
}

.page-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="%23ffffff" fill-opacity="0.1"/><circle cx="75" cy="75" r="1" fill="%23ffffff" fill-opacity="0.1"/><circle cx="50" cy="10" r="1" fill="%23ffffff" fill-opacity="0.1"/><circle cx="20" cy="80" r="1" fill="%23ffffff" fill-opacity="0.1"/></pattern></defs><rect width="100%" height="100%" fill="url(%23grain)"/></svg>');
  opacity: 0.3;
  pointer-events: none;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 16px;
  position: relative;
  z-index: 1;
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
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.title-content p {
  margin: 8px 0 0;
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}

/* 简化头部的标题样式 */
.page-header-simple .title-content h1 {
  font-size: 24px;
  margin: 0;
}

.page-actions {
  display: flex;
  gap: 12px;
  position: relative;
  z-index: 1;
}

.action-btn {
  padding: 12px 24px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.1);
  color: white;
  backdrop-filter: blur(10px);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
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

/* 搜索卡片 */
.search-card {
  margin-bottom: 24px;
  border-radius: 12px;
  border: none;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.search-card :deep(.el-card__body) {
  padding: 0;
}

.search-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-bottom: 1px solid #e9ecef;
}

.search-header h3 {
  margin: 0;
  color: #495057;
  font-size: 16px;
  font-weight: 600;
}

.expand-btn {
  color: #6c757d;
  padding: 4px 8px;
  transition: all 0.3s;
}

.expand-btn:hover {
  color: #007bff;
  background: rgba(0, 123, 255, 0.1);
}

.search-content {
  padding: 24px;
}

.search-row {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 16px;
}

.search-item {
  margin-bottom: 0;
}

.search-item :deep(.el-form-item__label) {
  font-weight: 500;
  color: #495057;
}

.search-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.search-btn {
  background: linear-gradient(135deg, #007bff, #0056b3);
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  transition: all 0.3s;
}

.search-btn:hover {
  background: linear-gradient(135deg, #0056b3, #007bff);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 123, 255, 0.3);
}

.reset-btn {
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  color: #6c757d;
  padding: 10px 20px;
  border-radius: 6px;
  transition: all 0.3s;
}

.reset-btn:hover {
  background: #e9ecef;
  border-color: #adb5bd;
}

/* 表格头部 */
.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding: 0 8px;
}

.table-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.count-info {
  color: #6c757d;
  font-size: 14px;
  padding: 8px 16px;
  background: #f8f9fa;
  border-radius: 20px;
  border: 1px solid #e9ecef;
}

/* 员工信息卡片 */
.name-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.name-cell .name {
  font-weight: 600;
  color: #212529;
  font-size: 14px;
}

.name-cell .employee-id {
  font-size: 12px;
  color: #6c757d;
  background: #f8f9fa;
  padding: 2px 8px;
  border-radius: 12px;
  display: inline-block;
  width: fit-content;
}

.text-gray {
  color: #6c757d;
  font-style: italic;
}

/* 分页样式 */
.pagination {
  margin-top: 24px;
  text-align: right;
  padding: 16px 0;
}

.pagination :deep(.el-pagination) {
  display: inline-flex;
  align-items: center;
  background: white;
  border-radius: 8px;
  padding: 8px 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}


/* 导入对话框样式 */
.upload-tips {
  margin-top: 16px;
  color: #6c757d;
  font-size: 12px;
  background: #f8f9fa;
  padding: 12px;
  border-radius: 6px;
  border-left: 4px solid #007bff;
}

.upload-tips p {
  margin: 4px 0;
}

/* 员工详情弹窗样式 */
.employee-detail-dialog {
  padding: 20px;
}

.employee-avatar-section {
  text-align: center;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  margin-bottom: 20px;
}

.employee-avatar-section h3 {
  margin: 16px 0 8px;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.employee-avatar-section p {
  margin: 8px 0 16px;
  color: #606266;
  font-size: 14px;
}

.employee-detail-dialog :deep(.el-descriptions) {
  margin-top: 20px;
}

.employee-detail-dialog :deep(.el-descriptions-item__label) {
  font-weight: 600;
  color: #303133;
  width: 120px;
}

.employee-detail-dialog :deep(.el-descriptions-item__content) {
  color: #606266;
}

/* 动画效果 */
.slide-fade-enter-active {
  transition: all 0.3s ease;
}

.slide-fade-leave-active {
  transition: all 0.3s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateY(-10px);
  opacity: 0;
}

/* 表格样式增强 */
.employee-list :deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.employee-list :deep(.el-table__header) {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
}

.employee-list :deep(.el-table th) {
  background: transparent;
  color: #495057;
  font-weight: 600;
  border-bottom: 2px solid #dee2e6;
}

.employee-list :deep(.el-table tr) {
  transition: background-color 0.3s;
}

.employee-list :deep(.el-table tr:hover) {
  background: rgba(0, 123, 255, 0.03);
}

.employee-list :deep(.el-table td) {
  border-bottom: 1px solid #f8f9fa;
  padding: 12px;
}

/* 卡片样式增强 */
.employee-list :deep(.el-card) {
  border: none;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s;
}

.employee-list :deep(.el-card:hover) {
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.employee-list :deep(.el-card__header) {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-bottom: 1px solid #e9ecef;
}

/* 按钮组样式 */
.employee-list :deep(.el-button-group) {
  border-radius: 6px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.employee-list :deep(.el-button-group .el-button) {
  border-radius: 0;
  border-right: 1px solid rgba(0, 0, 0, 0.1);
  transition: all 0.3s;
}

.employee-list :deep(.el-button-group .el-button:last-child) {
  border-right: none;
}

.employee-list :deep(.el-button-group .el-button:hover) {
  transform: none;
  z-index: 1;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .page-header-simple {
    padding: 12px 16px;
  }
  
  .page-header-simple .title-content h1 {
    font-size: 20px;
  }
  
  .page-actions {
    justify-content: center;
  }
  
  .search-row {
    flex-direction: column;
  }
  
  .search-actions {
    justify-content: center;
  }
  
  .table-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .table-actions {
    justify-content: center;
  }
}

/* 加载动画 */
.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
  background: #f8f9fa;
  border-radius: 8px;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e9ecef;
  border-top: 4px solid #007bff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 浮动操作按钮 */
.floating-actions {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 1000;
}

.floating-trigger {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.3);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 20px;
  border: none;
}

.floating-trigger:hover {
  transform: translateY(-2px) scale(1.05);
  box-shadow: 0 8px 24px rgba(64, 158, 255, 0.4);
}

.floating-trigger:active {
  transform: translateY(0) scale(0.95);
}

/* 下拉菜单样式增强 */
.floating-actions :deep(.el-dropdown-menu) {
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  border: none;
  padding: 8px 0;
  min-width: 160px;
}

.floating-actions :deep(.el-dropdown-menu__item) {
  padding: 12px 20px;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.floating-actions :deep(.el-dropdown-menu__item:hover) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.floating-actions :deep(.el-dropdown-menu__item i) {
  font-size: 16px;
  width: 16px;
  text-align: center;
}

/* 响应式适配 */
@media (max-width: 768px) {
  .floating-actions {
    bottom: 16px;
    right: 16px;
  }
  
  .floating-trigger {
    width: 48px;
    height: 48px;
    font-size: 18px;
  }
}
</style>