<template>
  <el-dialog
    :visible.sync="dialogVisible"
    :title="isEdit ? '编辑角色' : '新增角色'"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="角色名称" prop="name">
        <el-input
          v-model="form.name"
          placeholder="请输入角色名称"
          maxlength="50"
        />
      </el-form-item>
      
      <el-form-item label="角色编码" prop="code">
        <el-input
          v-model="form.code"
          placeholder="请输入角色编码"
          maxlength="50"
        />
      </el-form-item>
      
      <el-form-item label="角色描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="4"
          placeholder="请输入角色描述"
          maxlength="500"
        />
      </el-form-item>
      
      <el-form-item label="状态" prop="status">
        <el-radio-group v-model="form.status">
          <el-radio value="active">活跃</el-radio>
          <el-radio value="inactive">非活跃</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>

    <span slot="footer" class="dialog-footer">
      <el-button @click="handleClose">取消</el-button>
      <el-button
        type="primary"
        :loading="loading"
        @click="handleSubmit"
      >
        {{ isEdit ? '更新' : '创建' }}
      </el-button>
    </span>
  </el-dialog>
</template>

<script>
import { roleApi } from '@/services/systemApi'

export default {
  name: 'RoleDialog',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    role: {
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
        description: '',
        status: 'active'
      },
      rules: {
        name: [
          { required: true, message: '请输入角色名称', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        code: [
          { required: true, message: '请输入角色编码', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' },
          { pattern: /^[A-Za-z0-9_-]+$/, message: '编码只能包含字母、数字、下划线和连字符', trigger: 'blur' }
        ],
        description: [
          { required: true, message: '请输入角色描述', trigger: 'blur' },
          { max: 500, message: '长度不能超过 500 个字符', trigger: 'blur' }
        ],
        status: [
          { required: true, message: '请选择状态', trigger: 'change' }
        ]
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
      return !!this.role
    }
  },
  watch: {
    visible: {
      handler(visible) {
        if (visible) {
          this.initForm()
        }
      },
      immediate: true
    }
  },
  methods: {
    initForm() {
      if (this.role) {
        this.form = {
          name: this.role.name,
          code: this.role.code,
          description: this.role.description,
          status: this.role.status
        }
      } else {
        this.form = {
          name: '',
          code: '',
          description: '',
          status: 'active'
        }
      }
    },
    async handleSubmit() {
      if (!this.$refs.formRef) return
      
      try {
        await this.$refs.formRef.validate()
        this.loading = true
        
        if (this.isEdit && this.role) {
          await roleApi.updateRole(this.role.id, this.form)
          this.$message.success('更新成功')
        } else {
          await roleApi.createRole(this.form)
          this.$message.success('创建成功')
        }
        
        this.$emit('success')
      } catch (error) {
        if (error.response && error.response.data && error.response.data.error) {
          this.$message.error(error.response.data.error)
        } else {
          this.$message.error(this.isEdit ? '更新失败' : '创建失败')
        }
      } finally {
        this.loading = false
      }
    },
    handleClose() {
      if (this.$refs.formRef) {
        this.$refs.formRef.clearValidate()
      }
      this.dialogVisible = false
    }
  }
}
</script>

<style scoped>
.dialog-footer {
  text-align: right;
}
</style>