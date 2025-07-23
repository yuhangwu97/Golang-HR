<template>
  <el-dialog
    :visible.sync="dialogVisible"
    :title="isEdit ? '编辑权限' : '新增权限'"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="权限名称" prop="name">
        <el-input
          v-model="form.name"
          placeholder="请输入权限名称"
          maxlength="100"
        />
      </el-form-item>
      
      <el-form-item label="权限编码" prop="code">
        <el-input
          v-model="form.code"
          placeholder="请输入权限编码"
          maxlength="100"
        />
      </el-form-item>
      
      <el-form-item label="资源类型" prop="resource">
        <el-select
          v-model="form.resource"
          placeholder="请选择资源类型"
          style="width: 100%"
        >
          <el-option
            v-for="(label, value) in resourceLabels"
            :key="value"
            :label="label"
            :value="value"
          />
        </el-select>
      </el-form-item>
      
      <el-form-item label="操作类型" prop="action">
        <el-select
          v-model="form.action"
          placeholder="请选择操作类型"
          style="width: 100%"
        >
          <el-option
            v-for="(label, value) in actionLabels"
            :key="value"
            :label="label"
            :value="value"
          />
        </el-select>
      </el-form-item>
      
      <el-form-item label="权限描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="4"
          placeholder="请输入权限描述"
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
import { permissionApi } from '@/services/systemApi'

export default {
  name: 'PermissionDialog',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    permission: {
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
        resource: '',
        action: '',
        description: '',
        status: 'active'
      },
      // 资源和操作的选项
      resourceLabels: {
        employee: '员工管理',
        department: '部门管理',
        position: '职位管理',
        job_level: '职级管理',
        salary: '薪资管理',
        attendance: '考勤管理',
        system: '系统管理',
        user: '用户管理',
        role: '角色管理',
        permission: '权限管理'
      },
      actionLabels: {
        create: '创建',
        read: '查看',
        update: '更新',
        delete: '删除',
        manage: '管理',
        assign: '分配',
        approve: '审批'
      },
      rules: {
        name: [
          { required: true, message: '请输入权限名称', trigger: 'blur' },
          { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
        ],
        code: [
          { required: true, message: '请输入权限编码', trigger: 'blur' },
          { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' },
          { pattern: /^[A-Za-z0-9_.-]+$/, message: '编码只能包含字母、数字、下划线、点和连字符', trigger: 'blur' }
        ],
        resource: [
          { required: true, message: '请选择资源类型', trigger: 'change' }
        ],
        action: [
          { required: true, message: '请选择操作类型', trigger: 'change' }
        ],
        description: [
          { required: true, message: '请输入权限描述', trigger: 'blur' },
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
      return !!this.permission
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
      if (this.permission) {
        this.form = {
          name: this.permission.name,
          code: this.permission.code,
          resource: this.permission.resource,
          action: this.permission.action,
          description: this.permission.description,
          status: this.permission.status
        }
      } else {
        this.form = {
          name: '',
          code: '',
          resource: '',
          action: '',
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
        
        if (this.isEdit && this.permission) {
          await permissionApi.updatePermission(this.permission.id, this.form)
          this.$message.success('更新成功')
        } else {
          await permissionApi.createPermission(this.form)
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