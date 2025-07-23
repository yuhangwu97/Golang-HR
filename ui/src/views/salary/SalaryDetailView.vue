<template>
  <div class="salary-detail">
    <div class="header">
      <h1>薪资详情</h1>
      <div class="actions">
        <button @click="goBack" class="btn btn-secondary">
          返回
        </button>
        <button 
          v-if="salary && salary.status === 'approved'"
          @click="showPayrollRecords = true" 
          class="btn btn-primary"
        >
          发放记录
        </button>
      </div>
    </div>

    <div v-if="loading" class="loading">
      加载中...
    </div>

    <div v-else-if="salary" class="salary-content">
      <!-- 基本信息 -->
      <div class="info-section">
        <h3>基本信息</h3>
        <div class="info-grid">
          <div class="info-item">
            <label>薪资ID:</label>
            <span>{{ salary.id }}</span>
          </div>
          <div class="info-item">
            <label>员工姓名:</label>
            <span>{{ salary.employee?.name || '-' }}</span>
          </div>
          <div class="info-item">
            <label>员工工号:</label>
            <span>{{ salary.employee?.employee_id || '-' }}</span>
          </div>
          <div class="info-item">
            <label>部门:</label>
            <span>{{ salary.employee?.department?.name || '-' }}</span>
          </div>
          <div class="info-item">
            <label>薪资月份:</label>
            <span>{{ salary.month }}</span>
          </div>
          <div class="info-item">
            <label>状态:</label>
            <span :class="getStatusClass(salary.status)">
              {{ getStatusText(salary.status) }}
            </span>
          </div>
        </div>
      </div>

      <!-- 薪资详情 -->
      <div class="info-section">
        <h3>薪资详情</h3>
        <div class="salary-breakdown">
          <div class="breakdown-group">
            <h4>收入项目</h4>
            <div class="breakdown-item">
              <label>基本薪资:</label>
              <span class="amount">{{ formatCurrency(salary.base_salary) }}</span>
            </div>
            <div class="breakdown-item">
              <label>奖金:</label>
              <span class="amount">{{ formatCurrency(salary.bonus) }}</span>
            </div>
            <div class="breakdown-item">
              <label>津贴:</label>
              <span class="amount">{{ formatCurrency(salary.allowance) }}</span>
            </div>
            <div class="breakdown-item total">
              <label>应发薪资:</label>
              <span class="amount">{{ formatCurrency(salary.gross_salary) }}</span>
            </div>
          </div>

          <div class="breakdown-group">
            <h4>扣除项目</h4>
            <div class="breakdown-item">
              <label>个人所得税:</label>
              <span class="amount deduction">-{{ formatCurrency(salary.tax) }}</span>
            </div>
            <div class="breakdown-item">
              <label>社会保险:</label>
              <span class="amount deduction">-{{ formatCurrency(salary.social_security) }}</span>
            </div>
            <div class="breakdown-item">
              <label>住房公积金:</label>
              <span class="amount deduction">-{{ formatCurrency(salary.housing_fund) }}</span>
            </div>
            <div class="breakdown-item">
              <label>其他扣款:</label>
              <span class="amount deduction">-{{ formatCurrency(salary.deduction) }}</span>
            </div>
          </div>

          <div class="breakdown-group">
            <div class="breakdown-item final-amount">
              <label>实发薪资:</label>
              <span class="amount">{{ formatCurrency(salary.net_salary) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 备注 -->
      <div v-if="salary.remark" class="info-section">
        <h3>备注</h3>
        <p>{{ salary.remark }}</p>
      </div>

      <!-- 时间信息 -->
      <div class="info-section">
        <h3>时间信息</h3>
        <div class="info-grid">
          <div class="info-item">
            <label>创建时间:</label>
            <span>{{ formatDateTime(salary.created_at) }}</span>
          </div>
          <div class="info-item">
            <label>更新时间:</label>
            <span>{{ formatDateTime(salary.updated_at) }}</span>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="actions-section">
        <button 
          v-if="salary.status === 'calculated'"
          @click="approveSalary('approved')"
          class="btn btn-success"
        >
          审批通过
        </button>
        <button 
          v-if="salary.status === 'calculated'"
          @click="showRejectModal = true"
          class="btn btn-danger"
        >
          拒绝
        </button>
        <button 
          v-if="salary.status === 'approved'"
          @click="showPayrollModal = true"
          class="btn btn-primary"
        >
          发放薪资
        </button>
      </div>
    </div>

    <!-- 发放薪资模态框 -->
    <div v-if="showPayrollModal" class="modal">
      <div class="modal-content">
        <h3>发放薪资</h3>
        <div class="salary-info">
          <p>员工: {{ salary?.employee?.name }}</p>
          <p>月份: {{ salary?.month }}</p>
          <p>实发薪资: {{ salary ? formatCurrency(salary.net_salary) : '' }}</p>
        </div>
        <form @submit.prevent="processPayroll">
          <div class="form-group">
            <label>发放方式:</label>
            <select v-model="payrollForm.payment_method" required>
              <option value="bank">银行转账</option>
              <option value="cash">现金</option>
              <option value="check">支票</option>
            </select>
          </div>
          <div class="form-group">
            <label>银行账户:</label>
            <input type="text" v-model="payrollForm.bank_account" required />
          </div>
          <div class="form-actions">
            <button type="submit" class="btn btn-primary">确认发放</button>
            <button type="button" @click="closePayrollModal" class="btn">
              取消
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 拒绝模态框 -->
    <div v-if="showRejectModal" class="modal">
      <div class="modal-content">
        <h3>拒绝薪资</h3>
        <form @submit.prevent="approveSalary('rejected')">
          <div class="form-group">
            <label>拒绝原因:</label>
            <textarea v-model="rejectReason" rows="4" required></textarea>
          </div>
          <div class="form-actions">
            <button type="submit" class="btn btn-danger">确认拒绝</button>
            <button type="button" @click="showRejectModal = false" class="btn">
              取消
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 发放记录模态框 -->
    <div v-if="showPayrollRecords" class="modal large">
      <div class="modal-content">
        <h3>发放记录</h3>
        <div class="payroll-records">
          <table class="records-table">
            <thead>
              <tr>
                <th>发放日期</th>
                <th>发放方式</th>
                <th>银行账户</th>
                <th>发放金额</th>
                <th>状态</th>
                <th>处理人</th>
                <th>备注</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="record in payrollRecords" :key="record.id">
                <td>{{ formatDateTime(record.payment_date) }}</td>
                <td>{{ record.payment_method }}</td>
                <td>{{ record.bank_account }}</td>
                <td>{{ formatCurrency(record.payment_amount) }}</td>
                <td>{{ record.status }}</td>
                <td>{{ record.processor?.name || '-' }}</td>
                <td>{{ record.remark || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="form-actions">
          <button @click="showPayrollRecords = false" class="btn">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script >
import { ref, reactive, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { mapActions, mapGetters } from 'vuex'

const route = useRoute();
const router = useRouter();
const salaryStore = useSalaryStore();

const salary = ref(null);
const payrollRecords = ref([]);
const loading = ref(true);
const showPayrollModal = ref(false);
const showRejectModal = ref(false);
const showPayrollRecords = ref(false);
const rejectReason = ref('');

const payrollForm = reactive({
  payment_method: 'bank',
  bank_account: ''
});

onMounted(async () => {
  const id = Number(route.params.id);
  if (id) {
    await loadSalaryDetail(id);
  }
});

const loadSalaryDetail = async (id) => {
  try {
    salary.value = await salaryStore.fetchSalaryDetail(id);
  } catch (error) {
    console.error('Failed to load salary detail:', error);
  } finally {
    loading.value = false;
  }
};

const loadPayrollRecords = async () => {
  if (!salary.value) return;
  try {
    payrollRecords.value = await salaryStore.fetchPayrollRecords(salary.value.id);
  } catch (error) {
    console.error('Failed to load payroll records:', error);
  }
};

const goBack = () => {
  router.push('/salary');
};

const approveSalary = async (status) => {
  if (!salary.value) return;
  
  try {
    const remark = status === 'rejected' ? rejectReason.value : '';
    await salaryStore.approveSalary(salary.value.id, status, remark);
    alert(status === 'approved' ? '审批通过' : '已拒绝');
    await loadSalaryDetail(salary.value.id);
    showRejectModal.value = false;
    rejectReason.value = '';
  } catch (error) {
    alert('操作失败');
  }
};

const processPayroll = async () => {
  if (!salary.value) return;
  
  try {
    await salaryStore.processPayroll(salary.value.id, payrollForm);
    alert('薪资发放处理成功');
    closePayrollModal();
    await loadSalaryDetail(salary.value.id);
  } catch (error) {
    alert('薪资发放处理失败');
  }
};

const closePayrollModal = () => {
  showPayrollModal.value = false;
  payrollForm.payment_method = 'bank';
  payrollForm.bank_account = '';
};

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('zh-CN', {
    style: 'currency',
    currency: 'CNY'
  }).format(amount);
};

const formatDateTime = (dateTime) => {
  if (!dateTime) return '-';
  return new Date(dateTime).toLocaleString('zh-CN');
};

const getStatusClass = (status) => {
  const statusClasses = {
    draft: 'status-draft',
    calculated: 'status-calculated',
    approved: 'status-approved',
    paid: 'status-paid',
    rejected: 'status-rejected'
  };
  return statusClasses[status] || '';
};

const getStatusText = (status) => {
  const statusTexts = {
    draft: '草稿',
    calculated: '已计算',
    approved: '已审批',
    paid: '已发放',
    rejected: '已拒绝'
  };
  return statusTexts[status] || status;
};

// 当显示发放记录时加载数据
const handleShowPayrollRecords = async () => {
  showPayrollRecords.value = true;
  await loadPayrollRecords();
};
</script>

<style scoped>
.salary-detail {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.loading {
  text-align: center;
  padding: 50px;
  font-size: 18px;
}

.salary-content {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.info-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.info-section h3 {
  margin: 0 0 20px 0;
  color: #333;
  border-bottom: 2px solid #007bff;
  padding-bottom: 10px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 15px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid #eee;
}

.info-item label {
  font-weight: bold;
  color: #555;
}

.salary-breakdown {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 30px;
}

.breakdown-group {
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 20px;
}

.breakdown-group h4 {
  margin: 0 0 15px 0;
  color: #333;
  text-align: center;
  padding-bottom: 10px;
  border-bottom: 1px solid #ddd;
}

.breakdown-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #f5f5f5;
}

.breakdown-item.total {
  border-top: 2px solid #007bff;
  margin-top: 10px;
  padding-top: 15px;
  font-weight: bold;
}

.breakdown-item.final-amount {
  font-size: 18px;
  font-weight: bold;
  color: #007bff;
  border: 2px solid #007bff;
  border-radius: 8px;
  padding: 15px;
  margin-top: 20px;
}

.amount {
  font-weight: bold;
  color: #28a745;
}

.amount.deduction {
  color: #dc3545;
}

.actions-section {
  display: flex;
  gap: 10px;
  justify-content: center;
  padding: 20px;
}

.status-draft { color: #6c757d; }
.status-calculated { color: #17a2b8; }
.status-approved { color: #28a745; }
.status-paid { color: #007bff; }
.status-rejected { color: #dc3545; }

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal.large .modal-content {
  min-width: 80%;
  max-width: 90%;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
  min-width: 400px;
  max-height: 80vh;
  overflow-y: auto;
}

.salary-info {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 5px;
  margin-bottom: 15px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.form-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 20px;
}

.payroll-records {
  margin-bottom: 20px;
}

.records-table {
  width: 100%;
  border-collapse: collapse;
}

.records-table th,
.records-table td {
  padding: 10px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

.records-table th {
  background: #f8f9fa;
  font-weight: bold;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-primary {
  background: #007bff;
  color: white;
}

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-success {
  background: #28a745;
  color: white;
}

.btn-danger {
  background: #dc3545;
  color: white;
}

.btn:hover {
  opacity: 0.9;
}
</style>