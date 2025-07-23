// 员工相关常量和辅助函数

// 性别选项
export const GENDER_OPTIONS = [
  { label: '男', value: 'male' },
  { label: '女', value: 'female' }
]

// 合同类型选项
export const CONTRACT_TYPE_OPTIONS = [
  { label: '全职', value: 'full_time' },
  { label: '兼职', value: 'part_time' },
  { label: '实习', value: 'intern' },
  { label: '合同工', value: 'contract' }
]

// 学历选项
export const EDUCATION_OPTIONS = [
  { label: '高中', value: 'high_school' },
  { label: '大专', value: 'college' },
  { label: '本科', value: 'bachelor' },
  { label: '硕士', value: 'master' },
  { label: '博士', value: 'doctor' }
]

// 创建空的员工表单数据
export function createEmptyEmployeeForm() {
  return {
    name: '',
    email: '',
    phone: '',
    gender: '',
    birthday: null,
    id_card: '',
    department_id: null,
    position_id: null,
    job_level_id: null,
    manager_id: null,
    hire_date: null,
    probation_end_date: null,
    contract_type: '',
    base_salary: null,
    address: '',
    emergency_contact: '',
    emergency_phone: '',
    education: '',
    school: '',
    major: ''
  }
}

// 验证规则
export const EMPLOYEE_RULES = {
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { min: 2, message: '姓名至少2个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }
  ],
  department_id: [
    { required: true, message: '请选择部门', trigger: 'change' }
  ],
  hire_date: [
    { required: true, message: '请选择入职日期', trigger: 'change' }
  ]
}
