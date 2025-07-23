<template>
  <div class="dashboard">
    <!-- 欢迎横幅 -->
    <div class="welcome-banner slide-up">
      <div class="banner-content">
        <div class="welcome-info">
          <div class="welcome-avatar">
            <img v-if="user?.avatar" :src="user.avatar" :alt="user.name">
            <div v-else class="avatar-placeholder">{{ (user?.name || 'U').charAt(0).toUpperCase() }}</div>
          </div>
          <div class="welcome-text">
            <h1>欢迎回来，{{ (user && user.name) || '用户' }}！</h1>
            <p>{{ getGreeting() }}，让我们开始今天的工作吧</p>
          </div>
        </div>
        <div class="banner-actions">
          <div class="weather-widget">
            <i class="el-icon-sunny"></i>
            <span>{{ weather.temperature }}°C</span>
            <span class="weather-desc">{{ weather.description }}</span>
          </div>
          <div class="date-widget">
            <div class="date-main">{{ currentDate.date }}</div>
            <div class="date-sub">{{ currentDate.dayOfWeek }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 关键指标卡片 -->
    <div class="metrics-grid fade-in">
      <div class="metric-card" v-for="(metric, index) in keyMetrics" :key="index">
        <div class="metric-icon" :style="{ background: metric.color }">
          <i :class="metric.icon"></i>
        </div>
        <div class="metric-content">
          <div class="metric-label">{{ metric.label }}</div>
          <div class="metric-value">
            <span class="value">{{ metric.value }}</span>
            <span class="unit">{{ metric.unit }}</span>
          </div>
          <div class="metric-trend" :class="metric.trend > 0 ? 'positive' : 'negative'">
            <i :class="metric.trend > 0 ? 'el-icon-arrow-up' : 'el-icon-arrow-down'"></i>
            {{ Math.abs(metric.trend) }}% 
            <span class="trend-text">{{ metric.trend > 0 ? '比上月增长' : '比上月下降' }}</span>
          </div>
        </div>
        <div class="metric-chart">
          <canvas :ref="`chart-${index}`" width="80" height="40"></canvas>
        </div>
      </div>
    </div>
    
    <!-- 主要图表区域 -->
    <div class="charts-section zoom-in">
      <el-row :gutter="24">
        <!-- 部门分布图表 -->
        <el-col :xs="24" :lg="16" :xl="16">
          <div class="chart-card">
            <div class="chart-header">
              <div class="chart-title">
                <h3>部门人员分布</h3>
                <p>各部门员工数量统计</p>
              </div>
              <div class="chart-controls">
                <el-select v-model="chartTimeRange" size="small" style="width: 120px">
                  <el-option label="本月" value="month"></el-option>
                  <el-option label="本季度" value="quarter"></el-option>
                  <el-option label="本年" value="year"></el-option>
                </el-select>
              </div>
            </div>
            <div class="chart-content" v-loading="loading">
              <canvas ref="departmentChart" width="600" height="300"></canvas>
            </div>
          </div>
        </el-col>

        <!-- 员工状态饼图 -->
        <el-col :xs="24" :lg="8" :xl="8">
          <div class="chart-card">
            <div class="chart-header">
              <div class="chart-title">
                <h3>员工状态分布</h3>
                <p>在职/离职比例</p>
              </div>
            </div>
            <div class="chart-content">
              <canvas ref="statusChart" width="300" height="300"></canvas>
            </div>
          </div>
        </el-col>
      </el-row>

      <el-row :gutter="24" style="margin-top: 24px">
        <!-- 入职趋势图 -->
        <el-col :xs="24" :lg="12">
          <div class="chart-card">
            <div class="chart-header">
              <div class="chart-title">
                <h3>入职趋势</h3>
                <p>近12个月入职人数变化</p>
              </div>
            </div>
            <div class="chart-content">
              <canvas ref="hiresTrendChart" width="400" height="250"></canvas>
            </div>
          </div>
        </el-col>

        <!-- 年龄分布图 -->
        <el-col :xs="24" :lg="12">
          <div class="chart-card">
            <div class="chart-header">
              <div class="chart-title">
                <h3>年龄分布</h3>
                <p>员工年龄段统计</p>
              </div>
            </div>
            <div class="chart-content">
              <canvas ref="ageDistributionChart" width="400" height="250"></canvas>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 活动和通知区域 -->
    <div class="activity-section fade-in">
      <el-row :gutter="24">
        <!-- 最近活动 -->
        <el-col :xs="24" :lg="12">
          <div class="activity-card">
            <div class="activity-header">
              <h3>最近活动</h3>
              <el-button type="text" size="small">查看全部</el-button>
            </div>
            <div class="activity-timeline">
              <div 
                v-for="(activity, index) in recentActivities" 
                :key="index"
                class="timeline-item"
              >
                <div class="timeline-dot" :style="{ background: activity.color }"></div>
                <div class="timeline-content">
                  <div class="activity-title">{{ activity.title }}</div>
                  <div class="activity-desc">{{ activity.description }}</div>
                  <div class="activity-time">{{ activity.time }}</div>
                </div>
              </div>
            </div>
          </div>
        </el-col>

        <!-- 待办事项 -->
        <el-col :xs="24" :lg="12">
          <div class="todo-card">
            <div class="todo-header">
              <h3>待办事项</h3>
              <el-button type="text" size="small" @click="showAddTodo = true">
                <i class="el-icon-plus"></i>
                添加
              </el-button>
            </div>
            <div class="todo-list">
              <div 
                v-for="(todo, index) in todoList" 
                :key="index"
                class="todo-item"
                :class="{ completed: todo.completed }"
              >
                <el-checkbox 
                  v-model="todo.completed" 
                  @change="updateTodo(index)"
                >
                  {{ todo.text }}
                </el-checkbox>
                <div class="todo-priority" :class="todo.priority">
                  {{ todo.priority === 'high' ? '高' : todo.priority === 'medium' ? '中' : '低' }}
                </div>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 快捷操作区域 -->
    <div class="quick-actions-section zoom-in">
      <div class="actions-header">
        <h3>快捷操作</h3>
        <p>常用功能快速入口</p>
      </div>
      <div class="actions-grid">
        <div 
          v-for="(action, index) in quickActions" 
          :key="index"
          class="action-item"
          @click="handleQuickAction(action)"
        >
          <div class="action-icon" :style="{ background: action.color }">
            <i :class="action.icon"></i>
          </div>
          <div class="action-content">
            <div class="action-title">{{ action.title }}</div>
            <div class="action-desc">{{ action.description }}</div>
          </div>
          <div class="action-arrow">
            <i class="el-icon-arrow-right"></i>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加待办对话框 -->
    <el-dialog
      title="添加待办事项"
      :visible.sync="showAddTodo"
      width="400px"
    >
      <el-form>
        <el-form-item label="内容">
          <el-input v-model="newTodo.text" placeholder="请输入待办事项"></el-input>
        </el-form-item>
        <el-form-item label="优先级">
          <el-select v-model="newTodo.priority" style="width: 100%">
            <el-option label="高" value="high"></el-option>
            <el-option label="中" value="medium"></el-option>
            <el-option label="低" value="low"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="showAddTodo = false">取消</el-button>
        <el-button type="primary" @click="addTodo">确定</el-button>
      </div>
    </el-dialog>
    <!-- 悬浮待办事项 -->
    <FloatingTodo />
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { employeeService } from '@/services/employee'
import { departmentApi } from '@/services/departmentApi'
import { salaryService } from '@/services/salary'
import { dashboardService } from '@/services/dashboard'
import { formatDate } from '@/utils/date'
import Chart from 'chart.js/auto'
import FloatingTodo from '@/components/FloatingTodo.vue'
import dayjs from 'dayjs'

export default {
  name: 'DashboardView',
  components: {
    FloatingTodo
  },
  data() {
    return {
      loading: false,
      chartTimeRange: 'month',
      showAddTodo: false,
      refreshTimer: null,
      
      // 统计数据
      statistics: {
        total: 0,
        byDepartment: {},
        byStatus: {},
        recentHires: [],
        monthlyTrend: [],
        ageDistribution: {}
      },
      
      // 天气数据
      weather: {
        temperature: 24,
        description: '晴朗'
      },
      
      // 日期数据
      currentDate: {
        date: '',
        dayOfWeek: ''
      },
      
      // 图表实例
      chartInstances: {
        department: null,
        status: null,
        hiresTrend: null,
        ageDistribution: null
      },
      
      // 关键指标
      keyMetrics: [
        {
          label: '总员工数',
          value: 0,
          unit: '人',
          icon: 'el-icon-user',
          color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          trend: 8.5
        },
        {
          label: '本月入职',
          value: 0,
          unit: '人',
          icon: 'el-icon-plus',
          color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
          trend: 12.3
        },
        {
          label: '平均薪资',
          value: 0,
          unit: 'K',
          icon: 'el-icon-coin',
          color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
          trend: 5.2
        },
        {
          label: '出勤率',
          value: 0,
          unit: '%',
          icon: 'el-icon-time',
          color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
          trend: -2.1
        }
      ],
      
      // 最近活动
      recentActivities: [],
      
      // 待办事项
      todoList: [
        {
          text: '审核张三的入职申请',
          completed: false,
          priority: 'high'
        },
        {
          text: '准备月度绩效报告',
          completed: false,
          priority: 'medium'
        },
        {
          text: '更新员工培训计划',
          completed: true,
          priority: 'low'
        },
        {
          text: '安排新员工入职培训',
          completed: false,
          priority: 'high'
        }
      ],
      
      // 新待办事项
      newTodo: {
        text: '',
        priority: 'medium'
      },
      
      // 快捷操作
      quickActions: [
        {
          title: '添加员工',
          description: '录入新员工信息',
          icon: 'el-icon-user-solid',
          color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          action: 'add-employee'
        },
        {
          title: '薪资管理',
          description: '查看和管理薪资',
          icon: 'el-icon-coin',
          color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
          action: 'salary-management'
        },
        {
          title: '考勤统计',
          description: '查看考勤报表',
          icon: 'el-icon-time',
          color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
          action: 'attendance-stats'
        },
        {
          title: '绩效考核',
          description: '员工绩效评估',
          icon: 'el-icon-trophy',
          color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
          action: 'performance'
        },
        {
          title: '培训计划',
          description: '安排员工培训',
          icon: 'el-icon-reading',
          color: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
          action: 'training'
        },
        {
          title: '报表中心',
          description: '各类数据报表',
          icon: 'el-icon-data-analysis',
          color: 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)',
          action: 'reports'
        }
      ]
    }
  },
  computed: {
    ...mapGetters('auth', ['user'])
  },
  mounted() {
    this.initDashboard()
    this.startAutoRefresh()
  },
  beforeDestroy() {
    this.stopAutoRefresh()
  },
  methods: {
    formatDate,
    
    // 初始化仪表盘
    async initDashboard() {
      this.updateDateTime()
      this.fetchStatistics()
      this.fetchWeather()
      this.initCharts()
    },
    
    // 更新日期时间
    updateDateTime() {
      const now = new Date()
      const days = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
      this.currentDate = {
        date: now.getDate().toString().padStart(2, '0'),
        dayOfWeek: days[now.getDay()]
      }
    },
    
    // 获取问候语
    getGreeting() {
      const hour = new Date().getHours()
      if (hour < 12) return '早上好'
      if (hour < 18) return '下午好'
      return '晚上好'
    },
    
    // 获取统计数据
    async fetchStatistics() {
      this.loading = true
      try {
        // 使用 dashboardService 获取综合统计数据
        const dashboardResponse = await dashboardService.getDashboardData()
        
        if (dashboardResponse.success) {
          const data = dashboardResponse.data
          
          // 更新关键指标
          this.updateKeyMetrics(data)
          
          // 更新统计数据
          this.statistics = {
            total: data.employee.total || 0,
            byDepartment: data.department.employee_counts || {},
            byStatus: data.employee.by_status || {},
            recentHires: data.employee.recent_hires || [],
            monthlyTrend: data.employee.monthly_trend || [],
            ageDistribution: data.employee.age_distribution || {},
            departments: data.departments || []
          }
          
          // 更新最近活动
          this.updateRecentActivities(data)
        } else {
          console.error('获取仪表盘数据失败:', dashboardResponse.error)
          this.setDefaultData()
        }
        
        
        // 更新图表
        this.$nextTick(() => {
          this.updateCharts()
        })
        
      } catch (error) {
        console.error('获取统计数据失败:', error)
        this.$message.error('获取统计数据失败')
        this.setDefaultData()
      } finally {
        this.loading = false
      }
    },
    
    // 更新关键指标
    updateKeyMetrics(data) {
      const employeeData = data.employee || {}
      const salaryData = data.salary || {}
      const attendanceData = data.attendance || {}
      
      // 总员工数
      this.keyMetrics[0].value = employeeData.total || 0
      this.keyMetrics[0].trend = this.calculateTrend(employeeData.total, employeeData.last_month_total)
      
      // 本月入职
      this.keyMetrics[1].value = employeeData.current_month_hires || 0
      this.keyMetrics[1].trend = this.calculateTrend(employeeData.current_month_hires, employeeData.last_month_hires)
      
      // 平均薪资
      const avgSalary = salaryData.average_salary || 0
      this.keyMetrics[2].value = Math.round(avgSalary / 1000) // 转换为K
      this.keyMetrics[2].trend = this.calculateTrend(salaryData.average_salary, salaryData.last_month_average_salary)
      
      // 出勤率
      this.keyMetrics[3].value = Math.round(attendanceData.attendance_rate || 95)
      this.keyMetrics[3].trend = this.calculateTrend(attendanceData.attendance_rate, attendanceData.last_month_attendance_rate)
    },
    
    // 计算趋势百分比
    calculateTrend(current, previous) {
      if (!previous || previous === 0) return 0
      return Math.round(((current - previous) / previous) * 100 * 10) / 10
    },
    
    // 更新最近活动
    updateRecentActivities(data) {
      const recentHires = data.employee?.recent_hires || []
      const activities = []
      
      // 添加最近入职活动
      recentHires.slice(0, 5).forEach(hire => {
        activities.push({
          title: '新员工入职',
          description: `${hire.name} 加入了 ${hire.department?.name || '公司'}`,
          time: this.formatRelativeTime(hire.hire_date),
          color: '#67c23a'
        })
      })
      
      // 添加一些系统活动（模拟）
      activities.push(
        {
          title: '薪资发放',
          description: '本月薪资已发放完成',
          time: '3小时前',
          color: '#409eff'
        },
        {
          title: '考勤统计',
          description: '生成本周考勤报表',
          time: '1天前',
          color: '#e6a23c'
        }
      )
      
      this.recentActivities = activities.slice(0, 6)
    },
    
    // 设置默认数据
    setDefaultData() {
      this.statistics = {
        total: 0,
        byDepartment: {},
        byStatus: { active: 0, inactive: 0 },
        recentHires: [],
        monthlyTrend: [],
        ageDistribution: {}
      }
      
      this.keyMetrics.forEach(metric => {
        metric.value = 0
        metric.trend = 0
      })
      
      this.recentActivities = []
    },
    
    // 格式化相对时间
    formatRelativeTime(dateString) {
      const date = new Date(dateString)
      const now = new Date()
      const diffInHours = Math.floor((now - date) / (1000 * 60 * 60))
      
      if (diffInHours < 1) return '刚刚'
      if (diffInHours < 24) return `${diffInHours}小时前`
      
      const diffInDays = Math.floor(diffInHours / 24)
      if (diffInDays < 30) return `${diffInDays}天前`
      
      const diffInMonths = Math.floor(diffInDays / 30)
      return `${diffInMonths}个月前`
    },
    
    // 更新图表
    updateCharts() {
      this.renderDepartmentChart()
      this.renderStatusChart()
      this.renderHiresTrendChart()
      this.renderAgeDistributionChart()
    },
    
    // 获取薪资统计
    async fetchSalaryStatistics() {
      try {
        const response = await salaryService.getSalaryStatistics()
        return response.data || response
      } catch (error) {
        console.error('获取薪资统计失败:', error)
        // 返回默认数据
        return {
          average: 0,
          median: 0,
          range: {
            min: 0,
            max: 0
          }
        }
      }
    },
    
    // 更新关键指标
    updateKeyMetrics(basicStats, salaryStats) {
      this.keyMetrics[0].value = basicStats.total || 0
      this.keyMetrics[1].value = this.getMonthlyHires()
      this.keyMetrics[2].value = salaryStats?.average || 0
      this.keyMetrics[3].value = this.calculateAttendanceRate()
    },
    
    // 更新最近活动
    updateRecentActivities(basicStats) {
      const activities = []
      const recentHires = basicStats.recent_hires || basicStats.recentHires || []
      
      // 从最近入职员工生成活动
      recentHires.slice(0, 3).forEach(employee => {
        const departmentName = employee.department?.name || '未分配部门'
        const hireDateObj = new Date(employee.hire_date)
        const timeDiff = Date.now() - hireDateObj.getTime()
        const daysAgo = Math.floor(timeDiff / (1000 * 60 * 60 * 24))
        
        let timeText = ''
        if (daysAgo === 0) timeText = '今天'
        else if (daysAgo === 1) timeText = '昨天'
        else timeText = `${daysAgo}天前`
        
        activities.push({
          title: `${employee.name}入职${departmentName}`,
          description: `新员工${employee.name}已完成入职手续，分配至${departmentName}`,
          time: timeText,
          color: '#52c41a'
        })
      })
      
      // 添加一些系统活动
      if (activities.length < 4) {
        activities.push({
          title: '系统统计更新',
          description: '员工统计数据已更新，当前总员工数：' + (basicStats.total || 0),
          time: '刚刚',
          color: '#1890ff'
        })
      }
      
      this.recentActivities = activities
    },
    
    // 获取本月入职人数
    getMonthlyHires() {
      const currentMonth = new Date().getMonth()
      const currentYear = new Date().getFullYear()
      
      // 从最近入职员工中统计本月入职人数
      const recentHires = this.statistics.recent_hires || this.statistics.recentHires || []
      const monthlyHires = recentHires.filter(employee => {
        if (!employee.hire_date) return false
        const hireDate = new Date(employee.hire_date)
        return hireDate.getMonth() === currentMonth && hireDate.getFullYear() === currentYear
      })
      
      return monthlyHires.length
    },
    
    // 计算出勤率
    calculateAttendanceRate() {
      // 模拟计算逻辑
      return 94.7
    },
    
    // 获取天气数据
    async fetchWeather() {
      // 模拟天气API调用
      try {
        // 实际项目中可以调用天气API
        this.weather = {
          temperature: 24,
          description: '晴朗'
        }
      } catch (error) {
        console.error('获取天气数据失败:', error)
      }
    },
    
    // 初始化图表
    initCharts() {
      this.$nextTick(() => {
        this.renderDepartmentChart()
        this.renderStatusChart()
        this.renderHiresTrendChart()
        this.renderAgeDistributionChart()
        this.renderMetricCharts()
      })
    },
    
    // 渲染部门分布图表
    renderDepartmentChart() {
      const canvas = this.$refs.departmentChart
      if (!canvas) return
      
      // 销毁现有图表
      if (this.chartInstances.department) {
        this.chartInstances.department.destroy()
      }
      
      const ctx = canvas.getContext('2d')
      const departmentData = this.statistics.byDepartment || {}
      
      // 处理部门数据，确保显示部门名称而不是ID
      const departments = this.statistics.departments || []
      const labels = []
      const values = []
      
      Object.entries(departmentData).slice(0, 8).forEach(([deptId, count]) => {
        const dept = departments.find(d => d.id == deptId)
        labels.push(dept ? dept.name : `部门${deptId}`)
        values.push(count)
      })
      
      this.chartInstances.department = new Chart(ctx, {
        type: 'bar',
        data: {
          labels: labels.length > 0 ? labels : ['技术部', '产品部', '市场部', '人事部', '财务部'],
          datasets: [{
            label: '员工数量',
            data: values.length > 0 ? values : [25, 18, 15, 12, 8],
            backgroundColor: [
              'rgba(102, 126, 234, 0.8)',
              'rgba(118, 75, 162, 0.8)',
              'rgba(240, 147, 251, 0.8)',
              'rgba(245, 87, 108, 0.8)',
              'rgba(79, 172, 254, 0.8)',
              'rgba(45, 206, 137, 0.8)',
              'rgba(254, 176, 25, 0.8)',
              'rgba(255, 99, 132, 0.8)'
            ],
            borderColor: [
              'rgba(102, 126, 234, 1)',
              'rgba(118, 75, 162, 1)',
              'rgba(240, 147, 251, 1)',
              'rgba(245, 87, 108, 1)',
              'rgba(79, 172, 254, 1)',
              'rgba(45, 206, 137, 1)',
              'rgba(254, 176, 25, 1)',
              'rgba(255, 99, 132, 1)'
            ],
            borderWidth: 2,
            borderRadius: 8,
            borderSkipped: false
          }]
        },
        options: {
          responsive: true,
          plugins: {
            legend: {
              display: false
            },
            tooltip: {
              backgroundColor: 'rgba(0, 0, 0, 0.8)',
              titleColor: 'white',
              bodyColor: 'white',
              borderColor: 'rgba(255, 255, 255, 0.1)',
              borderWidth: 1,
              cornerRadius: 8
            }
          },
          scales: {
            y: {
              beginAtZero: true,
              grid: {
                color: 'rgba(0, 0, 0, 0.1)',
                drawBorder: false
              },
              ticks: {
                color: '#666'
              }
            },
            x: {
              grid: {
                display: false
              },
              ticks: {
                color: '#666'
              }
            }
          }
        }
      })
    },
    
    // 渲染状态饼图
    renderStatusChart() {
      const canvas = this.$refs.statusChart
      if (!canvas) return
      
      // 销毁现有图表
      if (this.chartInstances.status) {
        this.chartInstances.status.destroy()
      }
      
      const ctx = canvas.getContext('2d')
      const statusData = this.statistics.byStatus || {}
      
      const activeCount = statusData.active || 0
      const inactiveCount = statusData.inactive || 0
      
      this.chartInstances.status = new Chart(ctx, {
        type: 'doughnut',
        data: {
          labels: ['在职', '离职'],
          datasets: [{
            data: [activeCount || 180, inactiveCount || 20],
            backgroundColor: [
              '#67c23a',
              '#f56c6c'
            ],
            borderColor: [
              '#67c23a',
              '#f56c6c'
            ],
            borderWidth: 2,
            cutout: '60%'
          }]
        },
        options: {
          responsive: true,
          plugins: {
            legend: {
              position: 'bottom',
              labels: {
                usePointStyle: true,
                padding: 20,
                color: '#666'
              }
            },
            tooltip: {
              backgroundColor: 'rgba(0, 0, 0, 0.8)',
              titleColor: 'white',
              bodyColor: 'white',
              borderColor: 'rgba(255, 255, 255, 0.1)',
              borderWidth: 1,
              cornerRadius: 8,
              callbacks: {
                label: function(context) {
                  const total = context.dataset.data.reduce((a, b) => a + b, 0)
                  const percentage = ((context.parsed / total) * 100).toFixed(1)
                  return `${context.label}: ${context.parsed}人 (${percentage}%)`
                }
              }
            }
          },
          maintainAspectRatio: false
        }
      })
    },
    
    // 渲染入职趋势图
    renderHiresTrendChart() {
      const canvas = this.$refs.hiresTrendChart
      if (!canvas) return
      
      // 销毁现有图表
      if (this.chartInstances.hiresTrend) {
        this.chartInstances.hiresTrend.destroy()
      }
      
      const ctx = canvas.getContext('2d')
      const trendData = this.statistics.monthlyTrend || []
      
      const labels = trendData.map(item => item.month)
      const hiresData = trendData.map(item => item.hires)
      const leavesData = trendData.map(item => item.leaves)
      
      this.chartInstances.hiresTrend = new Chart(ctx, {
        type: 'line',
        data: {
          labels: labels.length > 0 ? labels : ['1月', '2月', '3月', '4月', '5月', '6月'],
          datasets: [{
            label: '入职人数',
            data: hiresData.length > 0 ? hiresData : [8, 12, 15, 10, 18, 22],
            borderColor: '#409eff',
            backgroundColor: 'rgba(64, 158, 255, 0.1)',
            borderWidth: 3,
            fill: true,
            tension: 0.4,
            pointBackgroundColor: '#409eff',
            pointBorderColor: '#fff',
            pointBorderWidth: 2,
            pointRadius: 6,
            pointHoverRadius: 8
          }, {
            label: '离职人数',
            data: leavesData.length > 0 ? leavesData : [2, 3, 1, 4, 2, 5],
            borderColor: '#f56c6c',
            backgroundColor: 'rgba(245, 108, 108, 0.1)',
            borderWidth: 3,
            fill: true,
            tension: 0.4,
            pointBackgroundColor: '#f56c6c',
            pointBorderColor: '#fff',
            pointBorderWidth: 2,
            pointRadius: 6,
            pointHoverRadius: 8
          }]
        },
        options: {
          responsive: true,
          plugins: {
            legend: {
              position: 'top',
              labels: {
                usePointStyle: true,
                padding: 20,
                color: '#666'
              }
            },
            tooltip: {
              backgroundColor: 'rgba(0, 0, 0, 0.8)',
              titleColor: 'white',
              bodyColor: 'white',
              borderColor: 'rgba(255, 255, 255, 0.1)',
              borderWidth: 1,
              cornerRadius: 8,
              intersect: false,
              mode: 'index'
            }
          },
          scales: {
            y: {
              beginAtZero: true,
              grid: {
                color: 'rgba(0, 0, 0, 0.1)',
                drawBorder: false
              },
              ticks: {
                color: '#666'
              }
            },
            x: {
              grid: {
                display: false
              },
              ticks: {
                color: '#666'
              }
            }
          },
          maintainAspectRatio: false
        }
      })
    },
    
    // 渲染年龄分布图
    renderAgeDistributionChart() {
      const canvas = this.$refs.ageDistributionChart
      if (!canvas) return
      
      // 销毁现有图表
      if (this.chartInstances.ageDistribution) {
        this.chartInstances.ageDistribution.destroy()
      }
      
      const ctx = canvas.getContext('2d')
      const ageData = this.statistics.ageDistribution || {}
      
      const labels = Object.keys(ageData)
      const values = Object.values(ageData)
      
      this.chartInstances.ageDistribution = new Chart(ctx, {
        type: 'bar',
        data: {
          labels: labels.length > 0 ? labels : ['20-25', '26-30', '31-35', '36-40', '41-45', '46-50', '50+'],
          datasets: [{
            label: '员工数量',
            data: values.length > 0 ? values : [45, 78, 65, 42, 28, 15, 8],
            backgroundColor: 'rgba(79, 172, 254, 0.8)',
            borderColor: 'rgba(79, 172, 254, 1)',
            borderWidth: 2,
            borderRadius: 8,
            borderSkipped: false
          }]
        },
        options: {
          responsive: true,
          plugins: {
            legend: {
              display: false
            },
            tooltip: {
              backgroundColor: 'rgba(0, 0, 0, 0.8)',
              titleColor: 'white',
              bodyColor: 'white',
              borderColor: 'rgba(255, 255, 255, 0.1)',
              borderWidth: 1,
              cornerRadius: 8,
              callbacks: {
                label: function(context) {
                  return `${context.label}岁: ${context.parsed.y}人`
                }
              }
            }
          },
          scales: {
            y: {
              beginAtZero: true,
              grid: {
                color: 'rgba(0, 0, 0, 0.1)',
                drawBorder: false
              },
              ticks: {
                color: '#666'
              }
            },
            x: {
              grid: {
                display: false
              },
              ticks: {
                color: '#666'
              }
            }
          },
          maintainAspectRatio: false
        }
      })
    },
    
    // 渲染指标卡片中的小图表
    renderMetricCharts() {
      this.keyMetrics.forEach((metric, index) => {
        const canvas = this.$refs[`chart-${index}`]?.[0]
        if (canvas) {
          const ctx = canvas.getContext('2d')
          this.drawSparkline(ctx, canvas.width, canvas.height, metric.trend)
        }
      })
    },
    
    // 绘制条形图
    drawBarChart(ctx, width, height) {
      // 使用真实的部门数据
      const departmentData = this.statistics.by_department || this.statistics.byDepartment || {}
      const data = Object.entries(departmentData).map(([name, count]) => ({
        label: name || '未分配',
        value: count
      })).slice(0, 5) // 只显示前5个部门
      
      // 如果没有数据，使用默认数据
      if (data.length === 0) {
        data.push(
          { label: '技术部', value: 0 },
          { label: '产品部', value: 0 },
          { label: '市场部', value: 0 },
          { label: '人事部', value: 0 },
          { label: '财务部', value: 0 }
        )
      }
      
      ctx.clearRect(0, 0, width, height)
      
      const maxValue = Math.max(...data.map(d => d.value))
      const barWidth = width / data.length - 20
      const barMaxHeight = height - 60
      
      data.forEach((item, index) => {
        const barHeight = (item.value / maxValue) * barMaxHeight
        const x = index * (barWidth + 20) + 10
        const y = height - barHeight - 30
        
        // 绘制渐变条形
        const gradient = ctx.createLinearGradient(0, y, 0, y + barHeight)
        gradient.addColorStop(0, '#667eea')
        gradient.addColorStop(1, '#764ba2')
        
        ctx.fillStyle = gradient
        ctx.fillRect(x, y, barWidth, barHeight)
        
        // 绘制标签
        ctx.fillStyle = '#666'
        ctx.font = '12px Arial'
        ctx.textAlign = 'center'
        ctx.fillText(item.label, x + barWidth / 2, height - 10)
        ctx.fillText(item.value, x + barWidth / 2, y - 5)
      })
    },
    
    // 绘制饼图
    drawPieChart(ctx, width, height) {
      // 使用真实的状态数据
      const statusData = this.statistics.by_status || this.statistics.byStatus || {}
      const data = [
        { label: '在职', value: statusData.active || statusData['在职'] || 0, color: '#52c41a' },
        { label: '离职', value: statusData.inactive || statusData['离职'] || 0, color: '#ff4d4f' }
      ]
      
      ctx.clearRect(0, 0, width, height)
      
      const centerX = width / 2
      const centerY = height / 2
      const radius = Math.min(width, height) / 2 - 40
      
      let currentAngle = -Math.PI / 2
      const total = data.reduce((sum, item) => sum + item.value, 0)
      
      data.forEach(item => {
        const sliceAngle = (item.value / total) * 2 * Math.PI
        
        ctx.beginPath()
        ctx.moveTo(centerX, centerY)
        ctx.arc(centerX, centerY, radius, currentAngle, currentAngle + sliceAngle)
        ctx.closePath()
        ctx.fillStyle = item.color
        ctx.fill()
        
        // 绘制标签
        const labelAngle = currentAngle + sliceAngle / 2
        const labelX = centerX + Math.cos(labelAngle) * (radius + 20)
        const labelY = centerY + Math.sin(labelAngle) * (radius + 20)
        
        ctx.fillStyle = '#666'
        ctx.font = '12px Arial'
        ctx.textAlign = 'center'
        ctx.fillText(`${item.label}: ${item.value}`, labelX, labelY)
        
        currentAngle += sliceAngle
      })
    },
    
    // 绘制折线图
    drawLineChart(ctx, width, height) {
      const data = this.statistics.monthlyTrend.slice(-6) // 最近6个月
      
      ctx.clearRect(0, 0, width, height)
      
      if (data.length === 0) return
      
      const maxValue = Math.max(...data.map(d => d.hires))
      const stepX = width / (data.length - 1)
      const stepY = (height - 40) / maxValue
      
      // 绘制线条
      ctx.strokeStyle = '#1890ff'
      ctx.lineWidth = 3
      ctx.beginPath()
      
      data.forEach((item, index) => {
        const x = index * stepX
        const y = height - 20 - (item.hires * stepY)
        
        if (index === 0) {
          ctx.moveTo(x, y)
        } else {
          ctx.lineTo(x, y)
        }
        
        // 绘制点
        ctx.fillStyle = '#1890ff'
        ctx.beginPath()
        ctx.arc(x, y, 4, 0, 2 * Math.PI)
        ctx.fill()
      })
      
      ctx.stroke()
    },
    
    // 绘制年龄分布图
    drawAgeChart(ctx, width, height) {
      const data = Object.entries(this.statistics.ageDistribution || {})
      
      ctx.clearRect(0, 0, width, height)
      
      if (data.length === 0) return
      
      const maxValue = Math.max(...data.map(([, value]) => value))
      const barWidth = width / data.length - 10
      const barMaxHeight = height - 40
      
      data.forEach(([age, count], index) => {
        const barHeight = (count / maxValue) * barMaxHeight
        const x = index * (barWidth + 10) + 5
        const y = height - barHeight - 20
        
        ctx.fillStyle = '#722ed1'
        ctx.fillRect(x, y, barWidth, barHeight)
        
        // 标签
        ctx.fillStyle = '#666'
        ctx.font = '10px Arial'
        ctx.textAlign = 'center'
        ctx.fillText(age, x + barWidth / 2, height - 5)
        ctx.fillText(count, x + barWidth / 2, y - 5)
      })
    },
    
    // 绘制迷你图表
    drawSparkline(ctx, width, height, trend) {
      ctx.clearRect(0, 0, width, height)
      
      // 生成一些随机数据点
      const points = []
      for (let i = 0; i < 10; i++) {
        points.push(Math.random() * height)
      }
      
      const stepX = width / (points.length - 1)
      
      ctx.strokeStyle = trend > 0 ? '#52c41a' : '#ff4d4f'
      ctx.lineWidth = 2
      ctx.beginPath()
      
      points.forEach((point, index) => {
        const x = index * stepX
        const y = height - point
        
        if (index === 0) {
          ctx.moveTo(x, y)
        } else {
          ctx.lineTo(x, y)
        }
      })
      
      ctx.stroke()
    },
    
    // 开始自动刷新
    startAutoRefresh() {
      this.refreshTimer = setInterval(() => {
        this.updateDateTime()
        // 每5分钟刷新一次数据
        if (new Date().getMinutes() % 5 === 0) {
          this.fetchStatistics()
        }
      }, 60000) // 每分钟更新一次时间
    },
    
    // 停止自动刷新
    stopAutoRefresh() {
      if (this.refreshTimer) {
        clearInterval(this.refreshTimer)
        this.refreshTimer = null
      }
    },
    
    // 处理快捷操作
    handleQuickAction(action) {
      switch (action.action) {
        case 'add-employee':
          this.$router.push('/employees/create')
          break
        case 'salary-management':
          this.$router.push('/salary')
          break
        case 'attendance-stats':
          this.$router.push('/attendance/stats')
          break
        case 'performance':
          this.$router.push('/performance')
          break
        case 'training':
          this.$router.push('/training')
          break
        case 'reports':
          this.$router.push('/reports')
          break
        default:
          this.$message.info('功能开发中...')
      }
    },
    
    // 更新待办事项
    updateTodo(index) {
      // 这里可以调用API更新数据库
      console.log('更新待办事项:', this.todoList[index])
    },
    
    // 添加待办事项
    addTodo() {
      if (!this.newTodo.text.trim()) {
        this.$message.warning('请输入待办事项内容')
        return
      }
      
      this.todoList.unshift({
        text: this.newTodo.text,
        completed: false,
        priority: this.newTodo.priority
      })
      
      this.newTodo = { text: '', priority: 'medium' }
      this.showAddTodo = false
      this.$message.success('添加成功')
    },

    // 导出数据
    async exportData() {
      try {
        this.$message({
          message: '正在导出数据...',
          type: 'info'
        })
        const blob = await employeeService.exportEmployees('excel')
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `employees_${this.formatDate(new Date())}.xlsx`
        a.click()
        window.URL.revokeObjectURL(url)
        this.$message.success('导出成功')
      } catch (error) {
        this.$message.error('导出失败')
      }
    }
  }
}
</script>

<style scoped>
/* 整体布局 */
.dashboard {
  padding: 0;
  background: var(--background-color);
  min-height: 100vh;
}

/* 欢迎横幅 */
.welcome-banner {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  padding: 32px;
  margin-bottom: 32px;
  color: white;
  position: relative;
  overflow: hidden;
}

.welcome-banner::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, transparent 70%);
  animation: float 6s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(180deg); }
}

.banner-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: relative;
  z-index: 2;
}

.welcome-info {
  display: flex;
  align-items: center;
  gap: 20px;
}

.welcome-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  overflow: hidden;
  border: 3px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.1);
}

.welcome-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 600;
  color: white;
  background: rgba(255, 255, 255, 0.2);
}

.welcome-text h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 600;
  margin-bottom: 8px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.welcome-text p {
  margin: 0;
  font-size: 16px;
  opacity: 0.9;
}

.banner-actions {
  display: flex;
  gap: 24px;
  align-items: center;
}

.weather-widget {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  backdrop-filter: blur(10px);
}

.weather-widget i {
  font-size: 20px;
  color: #ffd700;
}

.weather-desc {
  font-size: 14px;
  opacity: 0.8;
}

.date-widget {
  text-align: center;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  backdrop-filter: blur(10px);
}

.date-main {
  font-size: 24px;
  font-weight: 600;
  line-height: 1;
}

.date-sub {
  font-size: 12px;
  opacity: 0.8;
  margin-top: 4px;
}

/* 关键指标网格 */
.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.metric-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.05);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.metric-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: var(--primary-color);
}

.metric-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.12);
}

.metric-card {
  display: flex;
  align-items: center;
  gap: 16px;
}

.metric-icon {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 24px;
  flex-shrink: 0;
}

.metric-content {
  flex: 1;
}

.metric-label {
  font-size: 14px;
  color: #8c8c8c;
  margin-bottom: 4px;
}

.metric-value {
  display: flex;
  align-items: baseline;
  gap: 4px;
  margin-bottom: 8px;
}

.metric-value .value {
  font-size: 32px;
  font-weight: 700;
  color: #262626;
}

.metric-value .unit {
  font-size: 16px;
  color: #8c8c8c;
}

.metric-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 500;
}

.metric-trend.positive {
  color: #52c41a;
}

.metric-trend.negative {
  color: #ff4d4f;
}

.metric-trend .trend-text {
  margin-left: 4px;
  color: #8c8c8c;
}

.metric-chart {
  width: 80px;
  height: 40px;
  flex-shrink: 0;
}

/* 图表区域 */
.charts-section {
  margin-bottom: 32px;
}

.chart-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.05);
  height: 100%;
  transition: all 0.3s;
}

.chart-card:hover {
  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.12);
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.chart-title h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 4px;
}

.chart-title p {
  margin: 0;
  font-size: 14px;
  color: #8c8c8c;
}

.chart-content {
  position: relative;
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-content canvas {
  max-width: 100%;
  max-height: 100%;
}

/* 活动和待办区域 */
.activity-section {
  margin-bottom: 32px;
}

.activity-card, .todo-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.05);
  height: 100%;
}

.activity-header, .todo-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.activity-header h3, .todo-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #262626;
}

.activity-timeline {
  position: relative;
}

.activity-timeline::before {
  content: '';
  position: absolute;
  left: 8px;
  top: 0;
  bottom: 0;
  width: 2px;
  background: #f0f0f0;
}

.timeline-item {
  position: relative;
  padding-left: 32px;
  margin-bottom: 24px;
}

.timeline-item:last-child {
  margin-bottom: 0;
}

.timeline-dot {
  position: absolute;
  left: 0;
  top: 4px;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 3px solid white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.timeline-content {
  background: #fafafa;
  border-radius: 12px;
  padding: 16px;
  transition: all 0.3s;
}

.timeline-content:hover {
  background: #f0f0f0;
}

.activity-title {
  font-size: 14px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 4px;
}

.activity-desc {
  font-size: 13px;
  color: #595959;
  margin-bottom: 8px;
  line-height: 1.4;
}

.activity-time {
  font-size: 12px;
  color: #8c8c8c;
}

/* 待办事项 */
.todo-list {
  max-height: 400px;
  overflow-y: auto;
}

.todo-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
  transition: all 0.3s;
}

.todo-item:last-child {
  border-bottom: none;
}

.todo-item:hover {
  background: #fafafa;
  margin: 0 -12px;
  padding: 12px;
  border-radius: 8px;
}

.todo-item.completed {
  opacity: 0.5;
}

.todo-item.completed :deep(.el-checkbox__label) {
  text-decoration: line-through;
}

.todo-priority {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.todo-priority.high {
  background: #fff2f0;
  color: #ff4d4f;
}

.todo-priority.medium {
  background: #fff7e6;
  color: #fa8c16;
}

.todo-priority.low {
  background: #f6ffed;
  color: #52c41a;
}

/* 快捷操作区域 */
.quick-actions-section {
  margin-bottom: 32px;
}

.actions-header {
  text-align: center;
  margin-bottom: 24px;
}

.actions-header h3 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 8px;
}

.actions-header p {
  margin: 0;
  font-size: 16px;
  color: #8c8c8c;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.action-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.05);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.action-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.12);
}

.action-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 24px;
  flex-shrink: 0;
}

.action-content {
  flex: 1;
}

.action-title {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 4px;
}

.action-desc {
  font-size: 14px;
  color: #8c8c8c;
}

.action-arrow {
  color: #bfbfbf;
  font-size: 16px;
  transition: all 0.3s;
}

.action-item:hover .action-arrow {
  color: #1890ff;
  transform: translateX(4px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .dashboard {
    padding: 16px;
  }

  .welcome-banner {
    padding: 24px;
  }

  .banner-content {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }

  .banner-actions {
    flex-direction: column;
    gap: 16px;
  }

  .metrics-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .metric-card {
    flex-direction: column;
    text-align: center;
  }

  .actions-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .chart-content {
    height: 200px;
  }
}

@media (max-width: 480px) {
  .welcome-avatar {
    width: 48px;
    height: 48px;
  }

  .welcome-text h1 {
    font-size: 20px;
  }

  .welcome-text p {
    font-size: 14px;
  }

  .metric-value .value {
    font-size: 24px;
  }

  .action-item {
    padding: 16px;
  }
}

/* 加载状态 */
.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f0f0f0;
  border-top: 4px solid #1890ff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 滚动条美化 */
.todo-list::-webkit-scrollbar {
  width: 6px;
}

.todo-list::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.todo-list::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.todo-list::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* 对话框样式 */
:deep(.el-dialog) {
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
}

:deep(.el-dialog__header) {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-bottom: 1px solid #e9ecef;
  border-radius: 16px 16px 0 0;
}

:deep(.el-dialog__title) {
  font-weight: 600;
  color: #262626;
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-dialog__footer) {
  padding: 16px 24px;
  border-top: 1px solid #f0f0f0;
}
</style>