<template>
  <div class="test-create">
    <el-card>
      <div slot="header">
        <span>测试添加员工功能</span>
      </div>
      
      <el-button @click="testAPIs" type="primary">测试所有API</el-button>
      
      <div v-if="testResults.length > 0" style="margin-top: 20px;">
        <h3>测试结果：</h3>
        <el-alert
          v-for="(result, index) in testResults"
          :key="index"
          :title="result.name"
          :type="result.success ? 'success' : 'error'"
          :description="result.message"
          show-icon
          style="margin-bottom: 10px;"
        />
      </div>
      
      <div v-if="apiData" style="margin-top: 20px;">
        <h3>API数据：</h3>
        <el-collapse>
          <el-collapse-item title="部门数据" name="departments">
            <pre>{{ JSON.stringify(apiData.departments, null, 2) }}</pre>
          </el-collapse-item>
          <el-collapse-item title="职级数据" name="jobLevels">
            <pre>{{ JSON.stringify(apiData.jobLevels, null, 2) }}</pre>
          </el-collapse-item>
          <el-collapse-item title="经理数据" name="managers">
            <pre>{{ JSON.stringify(apiData.managers, null, 2) }}</pre>
          </el-collapse-item>
        </el-collapse>
      </div>
      
      <div style="margin-top: 20px;">
        <el-button @click="goToCreate" type="success">前往添加员工页面</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import { employeeService } from '@/services/employee'

export default {
  name: 'TestCreateView',
  data() {
    return {
      testResults: [],
      apiData: null
    }
  },
  methods: {
    async testAPIs() {
      this.testResults = []
      this.apiData = {
        departments: [],
        jobLevels: [],
        managers: []
      }
      
      // 测试部门API
      try {
        const deptResult = await employeeService.getDepartments()
        this.testResults.push({
          name: '部门列表API',
          success: deptResult.success,
          message: deptResult.success ? `成功获取 ${deptResult.data.length} 个部门` : deptResult.error
        })
        this.apiData.departments = deptResult.data
      } catch (error) {
        this.testResults.push({
          name: '部门列表API',
          success: false,
          message: error.message
        })
      }
      
      // 测试职级API
      try {
        const jobLevelResult = await employeeService.getJobLevels()
        this.testResults.push({
          name: '职级列表API',
          success: jobLevelResult.success,
          message: jobLevelResult.success ? `成功获取 ${jobLevelResult.data.length} 个职级` : jobLevelResult.error
        })
        this.apiData.jobLevels = jobLevelResult.data
      } catch (error) {
        this.testResults.push({
          name: '职级列表API',
          success: false,
          message: error.message
        })
      }
      
      // 测试经理API
      try {
        const managerResult = await employeeService.getManagers()
        this.testResults.push({
          name: '经理列表API',
          success: managerResult.success,
          message: managerResult.success ? `成功获取 ${managerResult.data.length} 个经理` : managerResult.error
        })
        this.apiData.managers = managerResult.data
      } catch (error) {
        this.testResults.push({
          name: '经理列表API',
          success: false,
          message: error.message
        })
      }
    },
    
    goToCreate() {
      this.$router.push('/employees/create')
    }
  }
}
</script>

<style scoped>
.test-create {
  padding: 20px;
}

pre {
  background: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
  overflow-x: auto;
  max-height: 300px;
}
</style>