<template>
  <div class="floating-todo" :class="{ 'is-expanded': isExpanded }">
    <!-- 悬浮按钮 -->
    <div class="float-trigger" @click="toggleExpanded">
      <i class="el-icon-notebook-2"></i>
      <span class="badge" v-if="todoCount > 0">{{ todoCount }}</span>
    </div>
    
    <!-- 待办事项面板 -->
    <transition name="slide-fade">
      <div v-if="isExpanded" class="todo-panel">
        <div class="todo-header">
          <h3>待办事项</h3>
          <div class="header-actions">
            <el-button 
              type="text" 
              size="small" 
              @click="showAddDialog = true"
              title="添加待办"
            >
              <i class="el-icon-plus"></i>
            </el-button>
            <el-button 
              type="text" 
              size="small" 
              @click="isExpanded = false"
              title="收起"
            >
              <i class="el-icon-close"></i>
            </el-button>
          </div>
        </div>
        
        <div class="todo-content">
          <div class="todo-filter">
            <el-radio-group v-model="filterType" size="mini">
              <el-radio-button label="all">全部</el-radio-button>
              <el-radio-button label="pending">待完成</el-radio-button>
              <el-radio-button label="completed">已完成</el-radio-button>
            </el-radio-group>
          </div>
          
          <div class="todo-list" v-if="filteredTodos.length">
            <div 
              v-for="(todo, index) in filteredTodos" 
              :key="todo.id"
              class="todo-item"
              :class="{ completed: todo.completed }"
            >
              <el-checkbox 
                v-model="todo.completed" 
                @change="updateTodo(todo)"
                :disabled="todo.completed"
              >
                {{ todo.text }}
              </el-checkbox>
              <div class="todo-meta">
                <span class="priority" :class="todo.priority">
                  {{ getPriorityText(todo.priority) }}
                </span>
                <span class="due-date" v-if="todo.dueDate">
                  {{ formatDate(todo.dueDate) }}
                </span>
              </div>
              <div class="todo-actions">
                <el-button 
                  type="text" 
                  size="mini" 
                  @click="editTodo(todo)"
                  title="编辑"
                >
                  <i class="el-icon-edit"></i>
                </el-button>
                <el-button 
                  type="text" 
                  size="mini" 
                  @click="deleteTodo(todo.id)"
                  title="删除"
                >
                  <i class="el-icon-delete"></i>
                </el-button>
              </div>
            </div>
          </div>
          
          <div class="todo-empty" v-else>
            <i class="el-icon-document-checked"></i>
            <p>{{ filterType === 'all' ? '暂无待办事项' : '暂无相关待办事项' }}</p>
          </div>
        </div>
      </div>
    </transition>
    
    <!-- 添加待办对话框 -->
    <el-dialog
      title="添加待办事项"
      :visible.sync="showAddDialog"
      width="400px"
      class="todo-dialog"
    >
      <el-form :model="newTodo" :rules="todoRules" ref="todoForm">
        <el-form-item label="内容" prop="text">
          <el-input
            v-model="newTodo.text"
            type="textarea"
            :rows="3"
            placeholder="请输入待办内容"
          />
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-select v-model="newTodo.priority" style="width: 100%">
            <el-option label="高优先级" value="high"></el-option>
            <el-option label="中优先级" value="medium"></el-option>
            <el-option label="低优先级" value="low"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="截止日期" prop="dueDate">
          <el-date-picker
            v-model="newTodo.dueDate"
            type="date"
            placeholder="选择截止日期"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="addTodo">添加</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import dayjs from 'dayjs'

export default {
  name: 'FloatingTodo',
  data() {
    return {
      isExpanded: false,
      showAddDialog: false,
      filterType: 'all',
      todos: [
        {
          id: 1,
          text: '审核新员工入职申请',
          priority: 'high',
          completed: false,
          dueDate: '2025-07-25',
          createdAt: '2025-07-22'
        },
        {
          id: 2,
          text: '准备月度员工绩效报告',
          priority: 'medium',
          completed: false,
          dueDate: '2025-07-30',
          createdAt: '2025-07-21'
        },
        {
          id: 3,
          text: '更新员工手册',
          priority: 'low',
          completed: true,
          dueDate: '2025-07-20',
          createdAt: '2025-07-18'
        }
      ],
      newTodo: {
        text: '',
        priority: 'medium',
        dueDate: null
      },
      todoRules: {
        text: [
          { required: true, message: '请输入待办内容', trigger: 'blur' },
          { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
        ],
        priority: [
          { required: true, message: '请选择优先级', trigger: 'change' }
        ]
      }
    }
  },
  computed: {
    todoCount() {
      return this.todos.filter(todo => !todo.completed).length
    },
    filteredTodos() {
      switch (this.filterType) {
        case 'pending':
          return this.todos.filter(todo => !todo.completed)
        case 'completed':
          return this.todos.filter(todo => todo.completed)
        default:
          return this.todos
      }
    }
  },
  methods: {
    toggleExpanded() {
      this.isExpanded = !this.isExpanded
    },
    updateTodo(todo) {
      if (todo.completed) {
        this.$message.success('任务已完成')
      }
      this.saveTodos()
    },
    addTodo() {
      this.$refs.todoForm.validate((valid) => {
        if (valid) {
          const newTodo = {
            id: Date.now(),
            text: this.newTodo.text,
            priority: this.newTodo.priority,
            completed: false,
            dueDate: this.newTodo.dueDate ? dayjs(this.newTodo.dueDate).format('YYYY-MM-DD') : null,
            createdAt: dayjs().format('YYYY-MM-DD')
          }
          this.todos.unshift(newTodo)
          this.saveTodos()
          this.showAddDialog = false
          this.resetForm()
          this.$message.success('待办事项添加成功')
        }
      })
    },
    editTodo(todo) {
      this.$message.info('编辑功能开发中')
    },
    deleteTodo(id) {
      this.$confirm('确定要删除这个待办事项吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.todos = this.todos.filter(todo => todo.id !== id)
        this.saveTodos()
        this.$message.success('删除成功')
      })
    },
    resetForm() {
      this.newTodo = {
        text: '',
        priority: 'medium',
        dueDate: null
      }
      this.$refs.todoForm && this.$refs.todoForm.resetFields()
    },
    saveTodos() {
      localStorage.setItem('hr-todos', JSON.stringify(this.todos))
    },
    loadTodos() {
      const saved = localStorage.getItem('hr-todos')
      if (saved) {
        this.todos = JSON.parse(saved)
      }
    },
    getPriorityText(priority) {
      const map = {
        high: '高',
        medium: '中',
        low: '低'
      }
      return map[priority] || '中'
    },
    formatDate(date) {
      return dayjs(date).format('MM-DD')
    }
  },
  mounted() {
    this.loadTodos()
    // 点击外部关闭
    document.addEventListener('click', (e) => {
      if (!this.$el.contains(e.target)) {
        this.isExpanded = false
      }
    })
  }
}
</script>

<style scoped>
.floating-todo {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
}

.float-trigger {
  width: 56px;
  height: 56px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  cursor: pointer;
  box-shadow: 0 4px 20px rgba(102, 126, 234, 0.3);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
}

.float-trigger:hover {
  transform: scale(1.05);
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.4);
}

.float-trigger i {
  font-size: 24px;
}

.badge {
  position: absolute;
  top: -5px;
  right: -5px;
  min-width: 20px;
  height: 20px;
  background: #f56c6c;
  color: white;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  border: 2px solid white;
}

.todo-panel {
  position: absolute;
  top: 70px;
  right: 0;
  width: 350px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
}

.todo-header {
  padding: 16px 20px;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-bottom: 1px solid #e9ecef;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.todo-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.todo-content {
  padding: 16px 20px;
  max-height: 400px;
  overflow-y: auto;
}

.todo-filter {
  margin-bottom: 16px;
}

.todo-filter :deep(.el-radio-group) {
  display: flex;
  width: 100%;
}

.todo-filter :deep(.el-radio-button) {
  flex: 1;
}

.todo-filter :deep(.el-radio-button__inner) {
  width: 100%;
}

.todo-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.todo-item {
  padding: 12px;
  background: #fafafa;
  border-radius: 8px;
  border-left: 4px solid #409eff;
  transition: all 0.3s;
}

.todo-item:hover {
  background: #f0f9ff;
  transform: translateX(2px);
}

.todo-item.completed {
  border-left-color: #67c23a;
  opacity: 0.6;
}

.todo-item.completed .todo-text {
  text-decoration: line-through;
}

.todo-meta {
  display: flex;
  gap: 8px;
  margin-top: 8px;
  font-size: 12px;
}

.priority {
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
  color: white;
}

.priority.high {
  background: #f56c6c;
}

.priority.medium {
  background: #e6a23c;
}

.priority.low {
  background: #909399;
}

.due-date {
  color: #909399;
}

.todo-actions {
  display: flex;
  gap: 4px;
  margin-top: 8px;
  opacity: 0;
  transition: opacity 0.3s;
}

.todo-item:hover .todo-actions {
  opacity: 1;
}

.todo-empty {
  text-align: center;
  padding: 40px 20px;
  color: #909399;
}

.todo-empty i {
  font-size: 48px;
  margin-bottom: 16px;
  display: block;
}

.todo-empty p {
  margin: 0;
  font-size: 14px;
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
  transform: translateX(20px);
  opacity: 0;
}

/* 对话框样式 */
.todo-dialog :deep(.el-dialog) {
  border-radius: 12px;
}

.todo-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-bottom: 1px solid #e9ecef;
}

/* 滚动条样式 */
.todo-content::-webkit-scrollbar {
  width: 6px;
}

.todo-content::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.todo-content::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.todo-content::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .floating-todo {
    top: 15px;
    right: 15px;
  }
  
  .todo-panel {
    width: 300px;
  }
  
  .float-trigger {
    width: 48px;
    height: 48px;
  }
  
  .float-trigger i {
    font-size: 20px;
  }
}
</style>