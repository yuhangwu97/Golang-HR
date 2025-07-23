<template>
  <el-dialog
    :visible.sync="dialogVisible"
    :title="isEdit ? '编辑职级' : '新增职级'"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="职级名称" prop="name">
        <el-input
          v-model="form.name"
          placeholder="请输入职级名称"
          maxlength="50"
        />
      </el-form-item>
      
      <el-form-item label="职级编码" prop="code">
        <el-input
          v-model="form.code"
          placeholder="请输入职级编码"
          maxlength="20"
        />
      </el-form-item>
      
      <el-form-item label="职级等级" prop="level">
        <el-input-number
          v-model="form.level"
          :min="1"
          :max="20"
          placeholder="请输入职级等级"
          style="width: 100%"
        />
        <div class="form-tip">等级数字越大，职级越高</div>
      </el-form-item>
      
      <el-form-item label="薪资范围" required>
        <div class="salary-range">
          <el-input-number
            v-model="form.min_salary"
            :min="0"
            :max="1000000"
            :step="1000"
            placeholder="最低薪资"
            style="width: 200px; margin-right: 10px;"
          />
          <span style="margin: 0 10px; line-height: 32px;">至</span>
          <el-input-number
            v-model="form.max_salary"
            :min="0"
            :max="1000000"
            :step="1000"
            placeholder="最高薪资"
            style="width: 200px"
          />
        </div>
      </el-form-item>
      
      <el-form-item label="职级描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="4"
          placeholder="请输入职级描述"
          maxlength="500"
        />
      </el-form-item>
      
      <el-form-item label="状态" prop="status">
        <el-radio-group v-model="form.status">
          <el-radio label="active">活跃</el-radio>
          <el-radio label="inactive">非活跃</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>

    <span slot="footer" class="dialog-footer">
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button
          type="primary"
          :loading="loading"
          @click="handleSubmit"
        >
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </div>
    </span>
  </el-dialog>
</template>

<script>
import { jobLevelApi } from '@/services/jobLevelApi'

export default {
  name: 'JobLevelDialog',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    jobLevel: {
      type: Object,
      default: null
    }
  },
  data() {
    return {
      loading: false,
      form: {
        name: '',
        code: '',
        level: 1,
        min_salary: 0,
        max_salary: 0,
        description: '',
        status: 'active'
      }
    }
  },
  computed: {
    dialogVisible: {
      get() {
        return this.visible
      },
      set(value) {
        this.$emit('update:visible', value)
      }
    },
    isEdit() {
      return !!this.jobLevel
    },
    rules() {
      return {
        name: [
          { required: true, message: '请输入职级名称', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        code: [
          { required: true, message: '请输入职级编码', trigger: 'blur' },
          { min: 1, max: 20, message: '长度在 1 到 20 个字符', trigger: 'blur' },
          { pattern: /^[A-Za-z0-9_-]+$/, message: '编码只能包含字母、数字、下划线和连字符', trigger: 'blur' }
        ],
        level: [
          { required: true, message: '请输入职级等级', trigger: 'blur' },
          { type: 'number', min: 1, max: 20, message: '等级必须在 1 到 20 之间', trigger: 'blur' }
        ],
        min_salary: [
          { required: true, message: '请输入最低薪资', trigger: 'blur' },
          { type: 'number', min: 0, message: '薪资不能为负数', trigger: 'blur' },
          { validator: this.validateSalaryRange, trigger: 'blur' }
        ],
        max_salary: [
          { required: true, message: '请输入最高薪资', trigger: 'blur' },
          { type: 'number', min: 0, message: '薪资不能为负数', trigger: 'blur' },
          { validator: this.validateSalaryRange, trigger: 'blur' }
        ],
        description: [
          { required: true, message: '请输入职级描述', trigger: 'blur' },
          { max: 500, message: '长度不能超过 500 个字符', trigger: 'blur' }
        ],
        status: [
          { required: true, message: '请选择状态', trigger: 'change' }
        ]
      }
    }
  },
  watch: {
    visible(val) {
      if (val) {
        this.initForm()
      }
    }
  },
  methods: {
    validateSalaryRange(rule, value, callback) {
      if (this.form.min_salary >= this.form.max_salary) {
        callback(new Error('最低薪资必须小于最高薪资'))
      } else {
        callback()
      }
    },
    initForm() {
      if (this.jobLevel) {
        Object.assign(this.form, {
          name: this.jobLevel.name,
          code: this.jobLevel.code,
          level: this.jobLevel.level,
          min_salary: this.jobLevel.min_salary,
          max_salary: this.jobLevel.max_salary,
          description: this.jobLevel.description,
          status: this.jobLevel.status
        })
      } else {
        Object.assign(this.form, {
          name: '',
          code: '',
          level: 1,
          min_salary: 0,
          max_salary: 0,
          description: '',
          status: 'active'
        })
      }
    },
    async handleSubmit() {
      if (!this.$refs.formRef) return
      
      try {
        await this.$refs.formRef.validate()
        this.loading = true
        
        // 调用实际API
        if (this.isEdit && this.jobLevel) {
          await jobLevelApi.updateJobLevel(this.jobLevel.id, this.form)
          this.$message.success('更新成功')
        } else {
          await jobLevelApi.createJobLevel(this.form)
          this.$message.success('创建成功')
        }
        
        this.$emit('success')
        this.handleClose()
      } catch (error) {
        console.error('操作失败:', error)
        if (error.response?.data?.message) {
          this.$message.error(error.response.data.message)
        } else {
          this.$message.error(this.isEdit ? '更新失败' : '创建失败')
        }
      } finally {
        this.loading = false
      }
    },
    handleClose() {
      this.$refs.formRef.clearValidate()
      this.dialogVisible = false
    }
  }
}
</script>

<style scoped>
.dialog-footer {
  text-align: right;
}

.salary-range {
  display: flex;
  align-items: center;
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}
</style>