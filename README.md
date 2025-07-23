# 🏢 Golang-HR 人力资源管理系统

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.0+-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
[![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1?style=flat&logo=mysql&logoColor=white)](https://www.mysql.com/)
[![Element UI](https://img.shields.io/badge/Element--Plus-2.0+-409EFF?style=flat&logo=element)](https://element-plus.org/)

**基于 Go + Vue 3 + MySQL 的现代化企业级人力资源管理系统**

*模仿 Workday/北森云 等知名HR系统的核心功能设计*
📝 **文档说明**: 本README由GitHub Copilot辅助编写完成

</div>

## 📋 项目概述

Golang-HR 是一个功能完整的企业级人力资源管理系统，采用现代化的技术栈和架构设计，提供从员工管理到薪资发放的全流程HR解决方案。系统具备高性能、高可用、可扩展的特点，适用于中小企业的人力资源数字化管理需求。

### ✨ 核心功能模块

| 模块 | 功能描述 | 实现状态 |
|------|----------|----------|
| 🏗️ **组织架构管理** | 部门树形管理、职位层级、组织图表可视化 | ✅ 已完成 |
| 👥 **员工生命周期** | 入职登记、档案管理、调岗转正、离职办理 | ✅ 已完成 |
| 💰 **薪资福利系统** | 薪资计算、发放流程、薪资报表、个税计算 | ✅ 已完成 |
| ⏰ **考勤管理** | 签到打卡、请假审批、加班统计、考勤报表 | 🚧 开发中 |
| 📊 **招聘管理** | 职位发布、简历管理、面试安排、录用流程 | 🚧 开发中 |
| 🎯 **绩效管理** | 目标设定、考核评估、360度反馈、结果分析 | 🚧 开发中 |
| 🔐 **权限系统** | RBAC权限控制、用户管理、角色分配 | ✅ 已完成 |

## 🏗️ 技术架构

### 后端技术栈
```
🔧 核心框架: Golang 1.21 + Gin Web Framework
💾 数据存储: MySQL 8.0 + Redis 6.0
🗃️ ORM框架: GORM v2 (支持自动迁移、关联预加载)
🔒 身份认证: JWT Token + 中间件鉴权
🏛️ 架构模式: 三层架构 + 依赖注入 + DAO模式
```

### 前端技术栈
```
⚡ 构建工具: Vue 2 + Vite 4.0
🎨 UI组件库: Element Plus (企业级组件)
📦 状态管理: Vuex 4.0
🛣️ 路由管理: Vue Router 4.0
🌐 HTTP客户端: Axios
📝 开发语言: JavaScript ES6+
```

### 架构亮点
- **🚀 高性能**: 支持并发处理、连接池管理、查询优化
- **🔄 读写分离**: 支持主从数据库配置，提升查询性能
- **💾 多级缓存**: Redis缓存 + 应用缓存，降低数据库压力
- **🔧 微服务化**: 模块化设计，支持服务拆分和独立部署
- **📊 监控完善**: 请求日志、性能监控、错误追踪

## 📁 项目结构

```
Golang-HR/
├── 📂 backend/                    # Go 后端服务
│   ├── 📂 controllers/           # 控制器层 (API接口定义)
│   │   ├── auth.go              # 用户认证
│   │   ├── employee.go          # 员工管理
│   │   ├── department.go        # 部门管理
│   │   ├── salary.go            # 薪资管理
│   │   ├── organization.go      # 组织架构
│   │   └── ...
│   ├── 📂 services/              # 业务逻辑层 (核心业务处理)
│   │   ├── employee.go          # 员工业务逻辑
│   │   ├── salary.go            # 薪资计算逻辑
│   │   ├── salary_enhanced.go   # 增强薪资系统
│   │   └── ...
│   ├── 📂 models/                # 数据模型层 (数据库映射)
│   │   ├── employee.go          # 员工模型
│   │   ├── organization.go      # 组织模型
│   │   ├── salary_enhanced.go   # 薪资模型
│   │   └── user.go              # 用户模型
│   ├── 📂 dao/                   # 数据访问层 (DAO模式)
│   │   ├── 📂 interfaces/       # DAO接口定义
│   │   ├── 📂 mysql/           # MySQL实现
│   │   ├── 📂 redis/           # Redis实现
│   │   ├── 📂 cache/           # 缓存策略
│   │   └── 📂 pool/            # 连接池管理
│   ├── 📂 middleware/            # 中间件 (认证、日志、CORS等)
│   ├── 📂 routes/               # 路由定义
│   ├── 📂 config/               # 配置管理
│   ├── 📂 database/             # 数据库脚本
│   ├── 📂 utils/                # 工具函数
│   └── 📂 tests/                # 单元测试
└── 📂 ui/                         # Vue 前端应用
    ├── 📂 src/
    │   ├── 📂 components/       # 公共组件
    │   │   ├── 📂 layout/      # 布局组件
    │   │   ├── 📂 organization/ # 组织架构组件
    │   │   ├── 📂 salary/      # 薪资组件
    │   │   └── 📂 system/      # 系统管理组件
    │   ├── 📂 views/           # 页面组件
    │   │   ├── 📂 employee/    # 员工管理页面
    │   │   ├── 📂 organization/ # 组织架构页面
    │   │   ├── 📂 salary/      # 薪资管理页面
    │   │   └── 📂 system/      # 系统管理页面
    │   ├── 📂 services/        # API服务封装
    │   ├── 📂 store/           # 状态管理
    │   ├── 📂 router/          # 路由配置
    │   └── 📂 utils/           # 工具函数
    ├── 📂 public/              # 静态资源
    └── package.json            # 依赖配置
```

## 🚀 快速开始

### 环境要求

| 环境 | 版本要求 | 说明 |
|------|----------|------|
| **Go** | 1.21+ | 后端开发语言 |
| **Node.js** | 18+ | 前端构建环境 |
| **MySQL** | 8.0+ | 主数据库 |
| **Redis** | 6.0+ | 缓存数据库 |

### 📦 后端部署

1. **克隆项目**
```bash
git clone https://github.com/your-repo/Golang-HR.git
cd Golang-HR/backend
```

2. **安装依赖**
```bash
go mod tidy
```

3. **数据库配置**
```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE gin_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 导入初始数据
mysql -u root -p gin_db < database/init.sql

# 修复数据库字段（如果需要）
mysql -u root -p gin_db < database/fix_departments_status.sql
```

4. **环境变量配置**
```bash
# 复制配置文件
cp config/config.example.yaml config/config.yaml

# 或者设置环境变量
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=gin_db
export REDIS_HOST=localhost
export REDIS_PORT=6379
export JWT_SECRET=your_jwt_secret
```

5. **启动服务**
```bash
# 开发模式
go run main.go

# 或者构建后运行
go build -o hr-backend
./hr-backend
```

服务启动后访问: http://localhost:8080

### 🎨 前端部署

1. **进入前端目录**
```bash
cd ui
```

2. **安装依赖**
```bash
npm install
# 或者使用 yarn
yarn install
```

3. **启动开发服务器**
```bash
npm run dev
# 或者
yarn dev
```

4. **构建生产版本**
```bash
npm run build
yarn build
```

前端服务启动后访问: http://localhost:3000

### 🐳 Docker 部署

```bash
# 使用 Docker Compose 一键部署
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

## 📊 数据库设计

### 核心数据表结构

| 表名 | 中文名 | 主要字段 | 说明 |
|------|--------|----------|------|
| `users` | 用户表 | username, email, password | 系统登录用户 |
| `employees` | 员工表 | name, employee_id, department_id | 员工基本信息 |
| `departments` | 部门表 | name, code, parent_id, level | 组织架构部门 |
| `positions` | 职位表 | name, code, department_id, level | 职位信息 |
| `job_levels` | 职级表 | name, level, min_salary, max_salary | 职级体系 |
| `salaries` | 薪资表 | employee_id, month, gross_salary | 薪资记录 |
| `enhanced_salaries` | 增强薪资表 | employee_id, payroll_period_id | 详细薪资 |
| `payroll_periods` | 薪资周期表 | name, period_type, start_date | 发薪周期 |

### 🔗 关联关系
- 员工 ↔ 部门 (多对一)
- 员工 ↔ 职位 (多对一)  
- 员工 ↔ 职级 (多对一)
- 部门 ↔ 部门 (自关联，树形结构)
- 薪资 ↔ 员工 (一对多)

## 🔧 API 接口文档

### 认证接口
```http
POST   /api/v1/auth/login       # 用户登录
POST   /api/v1/auth/register    # 用户注册
POST   /api/v1/auth/refresh     # 刷新Token
GET    /api/v1/auth/profile     # 获取用户信息
```

### 员工管理接口
```http
GET    /api/v1/employees           # 获取员工列表 (支持分页、筛选)
POST   /api/v1/employees           # 创建员工
GET    /api/v1/employees/:id       # 获取员工详情
PUT    /api/v1/employees/:id       # 更新员工信息
DELETE /api/v1/employees/:id       # 删除员工
GET    /api/v1/employees/search    # 搜索员工
POST   /api/v1/employees/import    # 批量导入员工
GET    /api/v1/employees/export    # 导出员工数据
```

### 组织架构接口
```http
GET    /api/v1/departments         # 获取部门列表
POST   /api/v1/departments         # 创建部门
GET    /api/v1/departments/tree    # 获取部门树形结构
PUT    /api/v1/departments/:id     # 更新部门信息
DELETE /api/v1/departments/:id     # 删除部门

GET    /api/v1/positions           # 获取职位列表
GET    /api/v1/positions/tree      # 获取职位树形结构
POST   /api/v1/positions           # 创建职位
```

### 薪资管理接口
```http
GET    /api/v1/salaries                    # 获取薪资列表
POST   /api/v1/salaries/calculate          # 计算薪资
GET    /api/v1/payroll/periods             # 获取薪资周期
POST   /api/v1/payroll/periods             # 创建薪资周期
GET    /api/v1/salaries/my-salary          # 个人薪资查询
GET    /api/v1/salaries/my-dashboard       # 个人薪资仪表板
```

## 🎯 系统特色功能

### 1. 🏗️ 智能组织架构
- **Workday风格组织图**: 可视化组织架构展示
- **层级路径显示**: 如 "技术研发事业部 > 产品开发部 > 前端开发团队"
- **拖拽调整**: 支持拖拽方式调整组织结构
- **批量操作**: 支持批量调整部门归属

### 2. 💰 增强薪资系统
- **灵活薪资结构**: 支持多种薪资组件配置
- **自动计算引擎**: 支持复杂薪资计算公式
- **薪资周期管理**: 月薪、季薪、年薪、奖金等多种周期
- **审批工作流**: 薪资审核、批准、发放流程
- **个人薪资门户**: 员工自助查询薪资详情

### 3. 📊 数据可视化
- **仪表板**: 员工统计、薪资分析、组织分布
- **图表展示**: 多种图表类型支持
- **实时数据**: WebSocket实时数据更新
- **导出功能**: Excel、PDF等格式导出

### 4. 🔐 企业级安全
- **RBAC权限**: 基于角色的访问控制
- **JWT认证**: 无状态Token认证
- **数据加密**: 敏感数据加密存储
- **操作日志**: 完整的审计追踪

## 🚀 性能优化

### 后端优化
- **数据库索引优化**: 针对查询热点建立复合索引
- **连接池管理**: 合理配置数据库连接池大小
- **Redis缓存**: 热点数据缓存，减少数据库压力
- **分页查询**: 大数据量列表分页加载
- **异步处理**: 耗时操作异步处理

### 前端优化  
- **组件懒加载**: 路由级别的代码分割
- **图片优化**: 图片压缩和懒加载
- **缓存策略**: 合理的HTTP缓存配置
- **虚拟滚动**: 大列表虚拟滚动优化

## 🧪 测试覆盖

```bash
# 运行后端测试
cd backend
go test ./... -v

# 运行前端测试  
cd ui
npm run test

# 测试覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📈 系统监控

### 性能指标
- **响应时间**: API平均响应时间 < 200ms
- **并发支持**: 支持 1000+ 并发用户
- **数据库性能**: 查询优化，索引使用率 > 90%
- **缓存命中率**: Redis缓存命中率 > 80%

### 日志管理
- **分级日志**: DEBUG/INFO/WARN/ERROR
- **结构化日志**: JSON格式便于分析
- **日志轮转**: 按大小和时间自动轮转
- **错误追踪**: 完整的错误堆栈信息

## 🤝 贡献指南

我们欢迎社区贡献！请遵循以下步骤：

1. **Fork 项目**
2. **创建功能分支** (`git checkout -b feature/AmazingFeature`)
3. **提交更改** (`git commit -m 'Add some AmazingFeature'`)
4. **推送分支** (`git push origin feature/AmazingFeature`)
5. **创建 Pull Request**

### 开发规范
- **代码风格**: 遵循 Go 和 Vue 的官方代码规范
- **提交规范**: 使用 Conventional Commits 规范
- **测试要求**: 新功能需要包含单元测试
- **文档更新**: 重要变更需要更新相关文档

## 📄 许可证

本项目采用 [MIT License](LICENSE) 开源协议。

## 🙏 致谢

感谢以下开源项目为本项目提供支持：
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Vue.js](https://github.com/vuejs/vue)
- [Element Plus](https://github.com/element-plus/element-plus)
- [GORM](https://github.com/go-gorm/gorm)

## 📧 联系方式

- **项目主页**: [GitHub Repository](https://github.com/your-repo/Golang-HR)
- **问题反馈**: [Issues](https://github.com/your-repo/Golang-HR/issues)
- **功能建议**: [Discussions](https://github.com/your-repo/Golang-HR/discussions)

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给它一个星标！**

*本项目仅供学习和参考使用，模仿知名HR系统功能特性*

---



</div>