<template>
  <el-dialog
    :visible.sync="dialogVisible"
    :title="isEdit ? '编辑职位' : '新增职位'"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="职位名称" prop="name">
        <el-input
          v-model="form.name"
          placeholder="请输入职位名称"
          maxlength="50"
        />
      </el-form-item>
      
      <el-form-item label="职位编码" prop="code">
        <el-input
          v-model="form.code"
          placeholder="请输入职位编码"
          maxlength="20"
        />
      </el-form-item>
      
      <el-form-item label="所属部门" prop="department_id">
        <DepartmentTreeSelect
          v-model="form.department_id"
          :tree-data="departmentTree"
          placeholder="请选择所属部门"
          @change="handleDepartmentChange"
          :show-code="true"
          :show-employee-count="true"
          :default-expand-all="false"
        />
      </el-form-item>
      
      <el-form-item label="上级职位" prop="parent_id">
        <el-select
          v-model="form.parent_id"
          placeholder="请选择上级职位（可选）"
          style="width: 100%"
          clearable
        >
          <el-option
            v-for="position in availableParentPositions"
            :key="position.id"
            :label="`${position.name} (${position.code})`"
            :value="position.id"
          />
        </el-select>
      </el-form-item>
      
      <el-form-item label="职位层级" prop="level">
        <el-input-number
          v-model="form.level"
          :min="1"
          :max="10"
          placeholder="请输入职位层级"
          style="width: 100%"
        />
      </el-form-item>
      
      <el-form-item label="排序" prop="sort">
        <el-input-number
          v-model="form.sort"
          :min="0"
          placeholder="请输入排序值"
          style="width: 100%"
        />
      </el-form-item>
      
      <el-form-item label="职位描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="请输入职位描述"
          maxlength="500"
        />
      </el-form-item>
      
      <el-form-item label="任职要求" prop="requirements">
        <el-input
          v-model="form.requirements"
          type="textarea"
          :rows="4"
          placeholder="请输入任职要求"
          maxlength="1000"
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
import { positionApi } from '@/services/positionApi'
import { departmentApi } from '@/services/departmentApi'
import DepartmentTreeSelect from '@/components/common/DepartmentTreeSelect.vue'

export default {
  name: 'PositionDialog',
  components: {
    DepartmentTreeSelect
  },
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    position: {
      type: Object,
      default: null
    },
    departments: {
      type: Array,
      default: () => []
    },
    parentPosition: {
      type: Object,
      default: null
    }
  },
  data() {
    return {
      loading: false,
      availableParentPositions: [],
      departmentTree: [],
      form: {
        name: '',
        code: '',
        department_id: undefined,
        parent_id: undefined,
        level: 1,
        sort: 0,
        description: '',
        requirements: '',
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
      return !!this.position
    },
    rules() {
      return {
        name: [
          { required: true, message: '请输入职位名称', trigger: 'blur' },
          { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
        ],
        code: [
          { required: true, message: '请输入职位编码', trigger: 'blur' },
          { min: 1, max: 20, message: '长度在 1 到 20 个字符', trigger: 'blur' },
          { pattern: /^[A-Za-z0-9_-]+$/, message: '编码只能包含字母、数字、下划线和连字符', trigger: 'blur' }
        ],
        department_id: [
          { required: true, message: '请选择所属部门', trigger: 'change' }
        ],
        level: [
          { required: true, message: '请输入职位层级', trigger: 'blur' }
        ],
        description: [
          { required: true, message: '请输入职位描述', trigger: 'blur' },
          { max: 500, message: '长度不能超过 500 个字符', trigger: 'blur' }
        ],
        requirements: [
          { required: true, message: '请输入任职要求', trigger: 'blur' },
          { max: 1000, message: '长度不能超过 1000 个字符', trigger: 'blur' }
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
        this.fetchAvailableParentPositions()
        this.loadDepartmentTree()
      }
    }
  },
  methods: {
    initForm() {
      if (this.position) {
        Object.assign(this.form, {
          name: this.position.name,
          code: this.position.code,
          department_id: this.position.department_id,
          parent_id: this.position.parent_id,
          level: this.position.level || 1,
          sort: this.position.sort || 0,
          description: this.position.description,
          requirements: this.position.requirements,
          status: this.position.status
        })
      } else {
        Object.assign(this.form, {
          name: '',
          code: '',
          department_id: undefined,
          parent_id: this.parentPosition ? this.parentPosition.id : undefined,
          level: this.parentPosition ? this.parentPosition.level + 1 : 1,
          sort: 0,
          description: '',
          requirements: '',
          status: 'active'
        })
      }
    },
    async fetchAvailableParentPositions() {
      try {
        const response = await positionApi.getAllPositions()
        const allPositions = response.data || []
        // 排除当前编辑的职位作为父级选项
        this.availableParentPositions = allPositions.filter(pos => {
          return !this.position || pos.id !== this.position.id
        })
      } catch (error) {
        this.$message.error('获取可选父级职位失败')
      }
    },
    async handleSubmit() {
      if (!this.$refs.formRef) return
      
      try {
        await this.$refs.formRef.validate()
        this.loading = true
        
        // 调用实际API
        if (this.isEdit && this.position) {
          await positionApi.updatePosition(this.position.id, this.form)
          this.$message.success('更新成功')
        } else {
          await positionApi.createPosition(this.form)
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

    // 加载部门树数据
    async loadDepartmentTree() {
      try {
        const response = await departmentApi.getDepartmentTree()
        this.departmentTree = this.processDepartmentTree(response.data?.data || response.data || [])
      } catch (error) {
        console.error('加载部门树失败:', error)
        this.$message.error('加载部门树失败')
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
          manager: node.manager,
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

    // 计算层级员工数量
    calculateHierarchicalEmployeeCount(node) {
      let totalCount = node.employee_count || 0
      
      if (node.children && node.children.length > 0) {
        for (const child of node.children) {
          totalCount += this.calculateHierarchicalEmployeeCount(child)
        }
      }
      
      return totalCount
    },

    // 处理部门变化
    async handleDepartmentChange(departmentId, departmentData) {
      if (departmentId) {
        // 根据部门自动推荐上级职位
        await this.suggestParentPosition(departmentId, departmentData)
      }
    },

    // 根据部门推荐上级职位
    async suggestParentPosition(departmentId, departmentData) {
      try {
        // 获取该部门及父部门的职位
        const response = await positionApi.getPositionsByDepartment(departmentId)
        const departmentPositions = response.data || []
        
        // 如果部门有经理职位，推荐为上级
        const managerPositions = departmentPositions.filter(pos => 
          pos.name.includes('经理') || pos.name.includes('总监') || pos.name.includes('主管')
        )
        
        if (managerPositions.length > 0) {
          // 选择层级最高的管理职位
          const topManagerPosition = managerPositions.reduce((prev, current) => 
            (prev.level < current.level) ? prev : current
          )
          this.form.parent_id = topManagerPosition.id
          this.form.level = topManagerPosition.level + 1
          this.$message.success(`已自动推荐上级职位：${topManagerPosition.name}`)
        } else if (departmentData && departmentData.manager) {
          // 如果部门有直属领导，查找该领导的职位
          const leaderPositions = this.availableParentPositions.filter(pos => 
            pos.department_id === departmentId && pos.name.includes('经理')
          )
          if (leaderPositions.length > 0) {
            this.form.parent_id = leaderPositions[0].id
            this.form.level = leaderPositions[0].level + 1
            this.$message.success(`已自动推荐上级职位：${leaderPositions[0].name}`)
          }
        }
      } catch (error) {
        console.error('推荐上级职位失败:', error)
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
</style>