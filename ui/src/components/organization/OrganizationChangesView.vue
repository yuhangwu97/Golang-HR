<template>
  <div class="organization-changes-view">
    <!-- 头部工具栏 -->
    <div class="changes-header">
      <div class="header-info">
        <h3 v-if="selectedUnit">{{ selectedUnit.name }} - 变更记录</h3>
        <h3 v-else>组织架构变更记录</h3>
        <p class="header-description">查看和管理组织架构的历史变更</p>
      </div>
      
      <div class="header-actions">
        <el-button-group>
          <el-button 
            :type="viewMode === 'changes' ? 'primary' : ''"
            @click="viewMode = 'changes'"
            icon="el-icon-document"
            size="small"
          >
            变更记录
          </el-button>
          <el-button 
            :type="viewMode === 'history' ? 'primary' : ''"
            @click="viewMode = 'history'"
            icon="el-icon-time"
            size="small"
          >
            历史快照
          </el-button>
          <el-button 
            :type="viewMode === 'timeline' ? 'primary' : ''"
            @click="viewMode = 'timeline'"
            icon="el-icon-s-grid"
            size="small"
          >
            时间线
          </el-button>
        </el-button-group>
        
        <div class="filter-controls">
          <el-select
            v-model="filterStatus"
            placeholder="状态筛选"
            size="small"
            style="width: 120px"
            @change="loadChanges"
            clearable
          >
            <el-option label="待审批" value="pending"></el-option>
            <el-option label="已审批" value="approved"></el-option>
            <el-option label="已拒绝" value="rejected"></el-option>
            <el-option label="已执行" value="implemented"></el-option>
            <el-option label="已取消" value="cancelled"></el-option>
          </el-select>
          
          <el-date-picker
            v-model="dateRange"
            type="datetimerange"
            size="small"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="yyyy-MM-dd HH:mm"
            value-format="yyyy-MM-dd HH:mm:ss"
            @change="loadChanges"
            style="width: 300px"
          />
          
          <el-button 
            type="primary" 
            size="small" 
            @click="showCreateChangeDialog = true"
            icon="el-icon-plus"
          >
            新建变更
          </el-button>
        </div>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="changes-content" v-loading="loading">
      <!-- 变更记录视图 -->
      <div v-if="viewMode === 'changes'" class="changes-list-view">
        <el-table 
          :data="changes" 
          stripe 
          style="width: 100%"
          @row-click="handleRowClick"
          row-key="id"
        >
          <el-table-column type="expand">
            <template slot-scope="{row}">
              <div class="change-details">
                <div class="detail-row">
                  <span class="detail-label">变更描述：</span>
                  <span class="detail-value">{{ row.change_description || '无' }}</span>
                </div>
                <div class="detail-row">
                  <span class="detail-label">变更原因：</span>
                  <span class="detail-value">{{ row.change_reason || '无' }}</span>
                </div>
                <div class="detail-row">
                  <span class="detail-label">业务理由：</span>
                  <span class="detail-value">{{ row.business_justification || '无' }}</span>
                </div>
                <div class="detail-row" v-if="row.impact_description">
                  <span class="detail-label">影响分析：</span>
                  <span class="detail-value">{{ row.impact_description }}</span>
                </div>
                <div class="detail-row" v-if="row.approval_note">
                  <span class="detail-label">审批意见：</span>
                  <span class="detail-value">{{ row.approval_note }}</span>
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column prop="change_type" label="变更类型" width="100">
            <template slot-scope="{row}">
              <el-tag :type="getChangeTypeColor(row.change_type)" size="small">
                {{ getChangeTypeName(row.change_type) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column prop="entity_name" label="变更对象" width="180"/>
          
          <el-table-column label="变更内容" width="200">
            <template slot-scope="{row}">
              <div v-if="row.field_name" class="change-content">
                <div class="field-name">{{ getFieldDisplayName(row.field_name) }}</div>
                <div class="value-change">
                  <span class="old-value">{{ row.old_value || '空' }}</span>
                  <i class="el-icon-right"></i>
                  <span class="new-value">{{ row.new_value || '空' }}</span>
                </div>
              </div>
              <div v-else>{{ row.change_description }}</div>
            </template>
          </el-table-column>
          
          <el-table-column prop="status" label="状态" width="100">
            <template slot-scope="{row}">
              <el-tag :type="getStatusColor(row.status)" size="small">
                {{ getStatusName(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column label="发起人" width="120">
            <template slot-scope="{row}">
              {{ row.initiator ? row.initiator.name : '系统' }}
            </template>
          </el-table-column>
          
          <el-table-column label="审批人" width="120">
            <template slot-scope="{row}">
              {{ row.approver ? row.approver.name : '-' }}
            </template>
          </el-table-column>
          
          <el-table-column prop="created_at" label="创建时间" width="150">
            <template slot-scope="{row}">
              {{ formatDateTime(row.created_at) }}
            </template>
          </el-table-column>
          
          <el-table-column prop="effective_date" label="生效时间" width="150">
            <template slot-scope="{row}">
              {{ row.effective_date ? formatDateTime(row.effective_date) : '-' }}
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="160" fixed="right">
            <template slot-scope="{row}">
              <el-button-group>
                <el-button 
                  v-if="row.status === 'pending'" 
                  type="success" 
                  size="mini"
                  @click.stop="handleApprove(row)"
                >
                  审批
                </el-button>
                <el-button 
                  v-if="row.status === 'pending'" 
                  type="danger" 
                  size="mini"
                  @click.stop="handleReject(row)"
                >
                  拒绝
                </el-button>
                <el-button 
                  type="primary" 
                  size="mini"
                  @click.stop="handleViewDetail(row)"
                >
                  详情
                </el-button>
              </el-button-group>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination
            background
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="pagination.page"
            :page-size="pagination.size"
            :page-sizes="[10, 20, 50, 100]"
            :total="pagination.total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>

      <!-- 历史快照视图 -->
      <div v-if="viewMode === 'history'" class="history-view">
        <div class="history-controls">
          <el-button @click="createSnapshot" type="primary" size="small">
            创建快照
          </el-button>
          <el-button @click="compareVersions" :disabled="selectedHistories.length !== 2" size="small">
            比较版本
          </el-button>
        </div>
        
        <el-table 
          :data="histories" 
          @selection-change="handleHistorySelection"
          row-key="id"
        >
          <el-table-column type="selection" :selectable="row => true" width="55"/>
          
          <el-table-column prop="snapshot_date" label="快照时间" width="150">
            <template slot-scope="{row}">
              {{ formatDateTime(row.snapshot_date) }}
            </template>
          </el-table-column>
          
          <el-table-column prop="snapshot_reason" label="快照原因" width="200"/>
          
          <el-table-column prop="unit_name" label="组织单元" width="180"/>
          
          <el-table-column prop="change_type" label="变更类型" width="100">
            <template slot-scope="{row}">
              <el-tag :type="getChangeTypeColor(row.change_type)" size="small">
                {{ getChangeTypeName(row.change_type) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column prop="changed_by_name" label="变更人" width="120"/>
          
          <el-table-column prop="employee_count" label="员工数" width="80"/>
          
          <el-table-column label="操作" width="160">
            <template slot-scope="{row}">
              <el-button-group>
                <el-button type="primary" size="mini" @click="viewHistory(row)">
                  查看
                </el-button>
                <el-button type="warning" size="mini" @click="rollbackToHistory(row)">
                  回滚
                </el-button>
              </el-button-group>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 时间线视图 -->
      <div v-if="viewMode === 'timeline'" class="timeline-view">
        <el-timeline>
          <el-timeline-item
            v-for="(item, index) in timeline"
            :key="index"
            :timestamp="formatDateTime(item.timestamp)"
            :type="getTimelineType(item.type)"
            :icon="getTimelineIcon(item.type)"
          >
            <el-card>
              <div class="timeline-content">
                <div class="timeline-header">
                  <span class="timeline-title">{{ item.title }}</span>
                  <el-tag :type="getChangeTypeColor(item.change_type)" size="small">
                    {{ getChangeTypeName(item.change_type) }}
                  </el-tag>
                </div>
                <div class="timeline-description">{{ item.description }}</div>
                <div class="timeline-meta">
                  <span>操作人：{{ item.operator_name }}</span>
                  <span v-if="item.impacted_employees">影响员工：{{ item.impacted_employees }}人</span>
                </div>
              </div>
            </el-card>
          </el-timeline-item>
        </el-timeline>
      </div>
    </div>

    <!-- 变更详情对话框 -->
    <el-dialog
      :visible.sync="showDetailDialog"
      title="变更详情"
      width="800px"
    >
      <div v-if="selectedChange" class="change-detail-dialog">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="变更类型">
            <el-tag :type="getChangeTypeColor(selectedChange.change_type)">
              {{ getChangeTypeName(selectedChange.change_type) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="变更对象">
            {{ selectedChange.entity_name }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusColor(selectedChange.status)">
              {{ getStatusName(selectedChange.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="发起人">
            {{ selectedChange.initiator ? selectedChange.initiator.name : '系统' }}
          </el-descriptions-item>
          <el-descriptions-item label="审批人">
            {{ selectedChange.approver ? selectedChange.approver.name : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="生效时间">
            {{ selectedChange.effective_date ? formatDateTime(selectedChange.effective_date) : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="变更原因" :span="2">
            {{ selectedChange.change_reason || '无' }}
          </el-descriptions-item>
          <el-descriptions-item label="变更描述" :span="2">
            {{ selectedChange.change_description || '无' }}
          </el-descriptions-item>
          <el-descriptions-item label="业务理由" :span="2">
            {{ selectedChange.business_justification || '无' }}
          </el-descriptions-item>
          <el-descriptions-item label="影响分析" :span="2">
            {{ selectedChange.impact_description || '无' }}
          </el-descriptions-item>
          <el-descriptions-item label="审批意见" :span="2">
            {{ selectedChange.approval_note || '无' }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    <!-- 审批对话框 -->
    <el-dialog
      :visible.sync="showApprovalDialog"
      title="审批变更"
      width="600px"
    >
      <el-form :model="approvalForm" :rules="approvalRules" ref="approvalForm">
        <el-form-item label="审批结果" prop="result">
          <el-radio-group v-model="approvalForm.result">
            <el-radio label="approved">通过</el-radio>
            <el-radio label="rejected">拒绝</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="审批意见" prop="note">
          <el-input
            v-model="approvalForm.note"
            type="textarea"
            :rows="4"
            placeholder="请输入审批意见"
          />
        </el-form-item>
        <el-form-item label="生效时间" v-if="approvalForm.result === 'approved'">
          <el-date-picker
            v-model="approvalForm.effective_date"
            type="datetime"
            placeholder="选择生效时间"
            format="yyyy-MM-dd HH:mm:ss"
            value-format="yyyy-MM-dd HH:mm:ss"
          />
        </el-form-item>
      </el-form>
      
      <div slot="footer">
        <el-button @click="showApprovalDialog = false">取消</el-button>
        <el-button type="primary" @click="submitApproval">确定</el-button>
      </div>
    </el-dialog>

    <!-- 新建变更对话框 -->
    <el-dialog
      :visible.sync="showCreateChangeDialog"
      title="新建变更"
      width="800px"
    >
      <el-form :model="changeForm" :rules="changeRules" ref="changeForm" label-width="120px">
        <el-form-item label="变更类型" prop="change_type">
          <el-select v-model="changeForm.change_type" placeholder="选择变更类型">
            <el-option label="创建" value="create"></el-option>
            <el-option label="更新" value="update"></el-option>
            <el-option label="删除" value="delete"></el-option>
            <el-option label="移动" value="move"></el-option>
            <el-option label="激活" value="activate"></el-option>
            <el-option label="停用" value="deactivate"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="变更对象" prop="entity_name">
          <el-input v-model="changeForm.entity_name" placeholder="输入变更对象名称"/>
        </el-form-item>
        <el-form-item label="变更原因" prop="change_reason">
          <el-input v-model="changeForm.change_reason" type="textarea" :rows="3"/>
        </el-form-item>
        <el-form-item label="变更描述" prop="change_description">
          <el-input v-model="changeForm.change_description" type="textarea" :rows="3"/>
        </el-form-item>
        <el-form-item label="业务理由" prop="business_justification">
          <el-input v-model="changeForm.business_justification" type="textarea" :rows="3"/>
        </el-form-item>
        <el-form-item label="生效时间" prop="effective_date">
          <el-date-picker
            v-model="changeForm.effective_date"
            type="datetime"
            placeholder="选择生效时间"
            format="yyyy-MM-dd HH:mm:ss"
            value-format="yyyy-MM-dd HH:mm:ss"
          />
        </el-form-item>
      </el-form>
      
      <div slot="footer">
        <el-button @click="showCreateChangeDialog = false">取消</el-button>
        <el-button type="primary" @click="submitCreateChange">提交</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import organizationApi from '@/services/organizationApi'

export default {
  name: 'OrganizationChangesView',
  props: {
    selectedUnit: {
      type: Object,
      default: null
    },
    changes: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      loading: false,
      viewMode: 'changes',
      
      // 筛选
      filterStatus: '',
      dateRange: [],
      
      // 数据
      histories: [],
      timeline: [],
      selectedHistories: [],
      
      // 分页
      pagination: {
        page: 1,
        size: 20,
        total: 0
      },
      
      // 对话框
      showDetailDialog: false,
      showApprovalDialog: false,
      showCreateChangeDialog: false,
      selectedChange: null,
      
      // 表单
      approvalForm: {
        result: 'approved',
        note: '',
        effective_date: ''
      },
      changeForm: {
        change_type: '',
        entity_name: '',
        change_reason: '',
        change_description: '',
        business_justification: '',
        effective_date: ''
      },
      
      // 验证规则
      approvalRules: {
        result: [{ required: true, message: '请选择审批结果', trigger: 'change' }],
        note: [{ required: true, message: '请输入审批意见', trigger: 'blur' }]
      },
      changeRules: {
        change_type: [{ required: true, message: '请选择变更类型', trigger: 'change' }],
        entity_name: [{ required: true, message: '请输入变更对象', trigger: 'blur' }],
        change_reason: [{ required: true, message: '请输入变更原因', trigger: 'blur' }]
      }
    }
  },
  watch: {
    selectedUnit: {
      handler() {
        this.loadData()
      },
      immediate: true
    },
    viewMode() {
      this.loadData()
    }
  },
  methods: {
    // 加载数据
    async loadData() {
      switch (this.viewMode) {
        case 'changes':
          await this.loadChanges()
          break
        case 'history':
          await this.loadHistories()
          break
        case 'timeline':
          await this.loadTimeline()
          break
      }
    },
    
    // 加载变更记录
    async loadChanges() {
      try {
        this.loading = true
        const params = {
          page: this.pagination.page,
          size: this.pagination.size
        }
        
        if (this.selectedUnit) {
          params.unit_id = this.selectedUnit.id
        }
        if (this.filterStatus) {
          params.status = this.filterStatus
        }
        if (this.dateRange && this.dateRange.length === 2) {
          params.start_date = this.dateRange[0]
          params.end_date = this.dateRange[1]
        }
        
        const response = await organizationApi.getOrganizationChanges(params)
        this.pagination.total = response.data.total
        this.$emit('update:changes', response.data.items)
      } catch (error) {
        console.error('Failed to load changes:', error)
        this.$message.error('加载变更记录失败')
      } finally {
        this.loading = false
      }
    },
    
    // 加载历史记录
    async loadHistories() {
      try {
        this.loading = true
        const params = {}
        if (this.selectedUnit) {
          params.unit_id = this.selectedUnit.id
        }
        
        const response = await organizationApi.getOrganizationSnapshots(params)
        this.histories = response.data.items || []
      } catch (error) {
        console.error('Failed to load histories:', error)
        this.$message.error('加载历史记录失败')
      } finally {
        this.loading = false
      }
    },
    
    // 加载时间线
    async loadTimeline() {
      try {
        this.loading = true
        const params = {}
        if (this.selectedUnit) {
          params.unit_id = this.selectedUnit.id
        }
        
        const response = await organizationApi.getOrganizationTimeline(params)
        this.timeline = response.data.items || []
      } catch (error) {
        console.error('Failed to load timeline:', error)
        this.$message.error('加载时间线失败')
      } finally {
        this.loading = false
      }
    },
    
    // 表格行点击
    handleRowClick(row) {
      this.selectedChange = row
      this.showDetailDialog = true
    },
    
    // 审批
    handleApprove(change) {
      this.selectedChange = change
      this.approvalForm = {
        result: 'approved',
        note: '',
        effective_date: ''
      }
      this.showApprovalDialog = true
    },
    
    // 拒绝
    handleReject(change) {
      this.selectedChange = change
      this.approvalForm = {
        result: 'rejected',
        note: '',
        effective_date: ''
      }
      this.showApprovalDialog = true
    },
    
    // 查看详情
    handleViewDetail(change) {
      this.selectedChange = change
      this.showDetailDialog = true
    },
    
    // 提交审批
    async submitApproval() {
      try {
        await this.$refs.approvalForm.validate()
        
        const data = {
          approval_note: this.approvalForm.note,
          effective_date: this.approvalForm.effective_date
        }
        
        if (this.approvalForm.result === 'approved') {
          await organizationApi.approveChange(this.selectedChange.id, data)
          this.$message.success('审批通过')
        } else {
          await organizationApi.rejectChange(this.selectedChange.id, data)
          this.$message.success('已拒绝变更')
        }
        
        this.showApprovalDialog = false
        this.loadChanges()
      } catch (error) {
        if (error !== 'validation failed') {
          console.error('Approval failed:', error)
          this.$message.error('审批失败')
        }
      }
    },
    
    // 提交新建变更
    async submitCreateChange() {
      try {
        await this.$refs.changeForm.validate()
        
        const data = {
          ...this.changeForm,
          entity_type: 'organization_unit',
          entity_id: this.selectedUnit?.id
        }
        
        await organizationApi.createOrganizationChange(data)
        this.$message.success('变更申请提交成功')
        this.showCreateChangeDialog = false
        this.changeForm = {
          change_type: '',
          entity_name: '',
          change_reason: '',
          change_description: '',
          business_justification: '',
          effective_date: ''
        }
        this.loadChanges()
      } catch (error) {
        if (error !== 'validation failed') {
          console.error('Create change failed:', error)
          this.$message.error('提交变更失败')
        }
      }
    },
    
    // 历史选择
    handleHistorySelection(selection) {
      this.selectedHistories = selection
    },
    
    // 创建快照
    async createSnapshot() {
      try {
        const data = {
          snapshot_reason: '手动创建快照',
          unit_id: this.selectedUnit?.id
        }
        
        await organizationApi.createOrganizationSnapshot(data)
        this.$message.success('快照创建成功')
        this.loadHistories()
      } catch (error) {
        console.error('Create snapshot failed:', error)
        this.$message.error('创建快照失败')
      }
    },
    
    // 比较版本
    compareVersions() {
      if (this.selectedHistories.length !== 2) {
        this.$message.warning('请选择两个版本进行比较')
        return
      }
      
      // TODO: 实现版本比较
      this.$message.info('版本比较功能开发中')
    },
    
    // 查看历史
    viewHistory(history) {
      // TODO: 实现历史查看
      this.$message.info('历史查看功能开发中')
    },
    
    // 回滚到历史版本
    async rollbackToHistory(history) {
      try {
        await this.$confirm(`确定要回滚到 ${this.formatDateTime(history.snapshot_date)} 的版本吗？`, '确认回滚', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        await organizationApi.rollbackToHistory(history.unit_id, history.id)
        this.$message.success('回滚成功')
        this.$emit('refresh')
      } catch (error) {
        if (error !== 'cancel') {
          console.error('Rollback failed:', error)
          this.$message.error('回滚失败')
        }
      }
    },
    
    // 分页
    handleSizeChange(size) {
      this.pagination.size = size
      this.loadChanges()
    },
    
    handleCurrentChange(page) {
      this.pagination.page = page
      this.loadChanges()
    },
    
    // 获取变更类型颜色
    getChangeTypeColor(type) {
      const colors = {
        create: 'success',
        update: 'primary',
        delete: 'danger',
        move: 'warning',
        activate: 'success',
        deactivate: 'info',
        assign: 'primary',
        unassign: 'warning',
        reassign: 'primary',
        merge: 'warning',
        split: 'warning'
      }
      return colors[type] || ''
    },
    
    // 获取变更类型名称
    getChangeTypeName(type) {
      const names = {
        create: '创建',
        update: '更新',
        delete: '删除',
        move: '移动',
        activate: '激活',
        deactivate: '停用',
        assign: '分配员工',
        unassign: '取消分配',
        reassign: '重新分配',
        merge: '合并',
        split: '拆分'
      }
      return names[type] || type
    },
    
    // 获取状态颜色
    getStatusColor(status) {
      const colors = {
        pending: 'warning',
        approved: 'success',
        rejected: 'danger',
        implemented: 'primary',
        cancelled: 'info'
      }
      return colors[status] || ''
    },
    
    // 获取状态名称
    getStatusName(status) {
      const names = {
        pending: '待审批',
        approved: '已审批',
        rejected: '已拒绝',
        implemented: '已执行',
        cancelled: '已取消'
      }
      return names[status] || status
    },
    
    // 获取字段显示名称
    getFieldDisplayName(fieldName) {
      const names = {
        name: '名称',
        parent_id: '上级单元',
        manager_id: '负责人',
        status: '状态',
        description: '描述'
      }
      return names[fieldName] || fieldName
    },
    
    // 获取时间线类型
    getTimelineType(type) {
      const types = {
        create: 'success',
        update: 'primary',
        delete: 'danger',
        move: 'warning'
      }
      return types[type] || 'primary'
    },
    
    // 获取时间线图标
    getTimelineIcon(type) {
      const icons = {
        create: 'el-icon-plus',
        update: 'el-icon-edit',
        delete: 'el-icon-delete',
        move: 'el-icon-sort'
      }
      return icons[type] || 'el-icon-document'
    },
    
    // 格式化时间
    formatDateTime(dateTime) {
      if (!dateTime) return '-'
      return new Date(dateTime).toLocaleString('zh-CN')
    }
  }
}
</script>

<style scoped>
.organization-changes-view {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.changes-header {
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
  background: #fafafa;
}

.header-info h3 {
  margin: 0 0 8px 0;
  color: #303133;
  font-size: 18px;
  font-weight: 600;
}

.header-description {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.header-actions {
  margin-top: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-controls {
  display: flex;
  gap: 12px;
  align-items: center;
}

.changes-content {
  flex: 1;
  padding: 20px;
  overflow: auto;
}

.changes-list-view,
.history-view,
.timeline-view {
  height: 100%;
}

.history-controls {
  margin-bottom: 16px;
  display: flex;
  gap: 12px;
}

.change-details {
  padding: 16px;
  background: #f8f9fa;
  border-radius: 6px;
}

.detail-row {
  margin-bottom: 12px;
  display: flex;
  align-items: flex-start;
}

.detail-row:last-child {
  margin-bottom: 0;
}

.detail-label {
  min-width: 80px;
  color: #606266;
  font-weight: 500;
}

.detail-value {
  flex: 1;
  color: #303133;
}

.change-content .field-name {
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.value-change {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.old-value {
  color: #f56c6c;
  text-decoration: line-through;
}

.new-value {
  color: #67c23a;
  font-weight: 500;
}

.timeline-content {
  padding: 12px;
}

.timeline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.timeline-title {
  font-weight: 500;
  color: #303133;
  font-size: 16px;
}

.timeline-description {
  color: #606266;
  margin-bottom: 8px;
  line-height: 1.5;
}

.timeline-meta {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: #909399;
}

.pagination-wrapper {
  margin-top: 20px;
  text-align: right;
}

.change-detail-dialog {
  padding: 16px;
}

/* 响应式适配 */
@media (max-width: 768px) {
  .header-actions {
    flex-direction: column;
    gap: 16px;
  }
  
  .filter-controls {
    width: 100%;
    justify-content: space-between;
    flex-wrap: wrap;
  }
  
  .value-change {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
}
</style>