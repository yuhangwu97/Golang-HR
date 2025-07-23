<template>
  <div class="job-level-list">
    <div class="page-header">
      <h2>职级管理</h2>
      <div class="actions">
        <el-button type="primary" @click="showCreateDialog = true">
          <i class="el-icon-plus"></i>
          新增职级
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
              <div class="stat-label">总职级数</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ activeJobLevels }}</div>
              <div class="stat-label">活跃职级</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ inactiveJobLevels }}</div>
              <div class="stat-label">非活跃职级</div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            <div class="stat-card">
              <div class="stat-number">{{ totalEmployees }}</div>
              <div class="stat-label">员工总数</div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 搜索过滤 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="职级名称">
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入职级名称或编码"
            clearable
            @keyup.enter.native="handleSearch"
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

    <!-- 职级表格 -->
    <el-card>
      <el-table
        v-loading="loading"
        :data="jobLevels"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="职级名称" width="200" />
        <el-table-column prop="code" label="职级编码" width="150" />
        <el-table-column prop="level" label="职级等级" width="120">
          <template slot-scope="{ row }">
            <el-tag type="info">{{ row.level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="min_salary" label="最低薪资" width="120">
          <template slot-scope="{ row }">
            ¥{{ row.min_salary ? row.min_salary.toLocaleString() : '0' }}
          </template>
        </el-table-column>
        <el-table-column prop="max_salary" label="最高薪资" width="120">
          <template slot-scope="{ row }">
            ¥{{ row.max_salary ? row.max_salary.toLocaleString() : '0' }}
          </template>
        </el-table-column>
        <el-table-column prop="description" label="职级描述" show-overflow-tooltip />
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
    </el-card>

    <!-- 创建/编辑对话框 -->
    <JobLevelDialog
      :visible.sync="showCreateDialog"
      :job-level="editJobLevel"
      @success="handleDialogSuccess"
    />
  </div>
</template>

<script>
import { jobLevelApi } from '@/services/jobLevelApi'
import JobLevelDialog from '@/components/organization/JobLevelDialog.vue'
import { formatDate } from '@/utils/date'

export default {
  name: 'JobLevelListView',
  components: {
    JobLevelDialog
  },
  data() {
    return {
      loading: false,
      showStatistics: false,
      showCreateDialog: false,
      jobLevels: [],
      statistics: null,
      editJobLevel: null,
      selectedJobLevels: [],
      searchForm: {
        keyword: '',
        status: '',
        page: 1,
        page_size: 20
      },
      pagination: {
        page: 1,
        page_size: 20,
        total_items: 0,
        total_pages: 0
      }
    }
  },
  computed: {
    activeJobLevels() {
      return (this.statistics && this.statistics.by_status && this.statistics.by_status.active) || 0
    },
    inactiveJobLevels() {
      return (this.statistics && this.statistics.by_status && this.statistics.by_status.inactive) || 0
    },
    totalEmployees() {
      if (!this.statistics || !this.statistics.employee_count) {
        return 0
      }
      return Object.values(this.statistics.employee_count).reduce((sum, count) => sum + count, 0)
    }
  },
  mounted() {
    this.fetchJobLevels()
    this.fetchStatistics()
  },
  methods: {
    formatDate,
    async fetchJobLevels() {
      this.loading = true
      try {
        const response = await jobLevelApi.getJobLevels({
          ...this.searchForm,
          page: this.pagination.page,
          page_size: this.pagination.page_size
        })
        
        // 处理API响应数据结构
        const responseData = response.data || response
        this.jobLevels = responseData.data || []
        this.pagination.total_items = responseData.total_items || 0
        this.pagination.total_pages = responseData.total_pages || 0
      } catch (error) {
        console.error('获取职级列表失败:', error)
        this.$message.error('获取职级列表失败')
      } finally {
        this.loading = false
      }
    },
    async fetchStatistics() {
      try {
        const response = await jobLevelApi.getJobLevelStatistics()
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
      this.fetchJobLevels()
    },
    handleReset() {
      Object.assign(this.searchForm, {
        keyword: '',
        status: '',
        page: 1,
        page_size: 20
      })
      this.pagination.page = 1
      this.fetchJobLevels()
    },
    handleEdit(jobLevel) {
      this.editJobLevel = jobLevel
      this.showCreateDialog = true
    },
    async handleDelete(jobLevel) {
      try {
        await this.$confirm(
          `确定要删除职级"${jobLevel.name}"吗？`,
          '确认删除',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        await jobLevelApi.deleteJobLevel(jobLevel.id)
        this.$message.success('删除成功')
        this.fetchJobLevels()
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error('删除失败')
        }
      }
    },
    handleSelectionChange(selected) {
      this.selectedJobLevels = selected
    },
    handleSizeChange(size) {
      this.pagination.page_size = size
      this.pagination.page = 1
      this.fetchJobLevels()
    },
    handleCurrentChange(page) {
      this.pagination.page = page
      this.fetchJobLevels()
    },
    handleDialogSuccess() {
      this.showCreateDialog = false
      this.editJobLevel = null
      this.fetchJobLevels()
      if (this.showStatistics) {
        this.fetchStatistics()
      }
    },
    // 计算基础统计信息
    calculateBasicStatistics() {
      const total = this.jobLevels.length
      const activeJobLevels = this.jobLevels.filter(level => level.status === 'active').length
      const inactiveJobLevels = total - activeJobLevels
      const employeeCount = {}
      
      this.jobLevels.forEach(level => {
        // 假设职级有员工数量字段，如果没有则默认为0
        employeeCount[level.name] = level.employee_count || 0
      })

      return {
        total,
        active_job_levels: activeJobLevels,
        inactive_job_levels: inactiveJobLevels,
        employee_count: employeeCount
      }
    }
  }
}
</script>

<style scoped>
.job-level-list {
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
</style>