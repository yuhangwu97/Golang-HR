# Golang-HR 人力资源管理系统

> 基于 Vue 3 + Golang Gin + MySQL + Redis 的现代化人力资源管理系统

## 📋 项目简介

Golang-HR 是一个功能完整的人力资源管理系统，模仿北森云系统的核心功能，提供员工管理、组织架构、考勤、薪资、招聘、绩效等全方位的HR管理解决方案。

### 🌟 核心特性

- **员工管理**: 员工档案、入职离职、组织关系管理
- **组织架构**: 部门管理、职位管理、层级关系
- **考勤管理**: 签到签退、请假审批、考勤统计
- **薪资管理**: 薪资计算、发放记录、薪资报表
- **招聘管理**: 职位发布、简历管理、面试流程
- **绩效管理**: 目标设定、考核评估、结果分析
- **系统管理**: 用户权限、角色管理、系统配置

## 🏗️ 技术架构

### 后端技术栈
- **框架**: Golang + Gin
- **数据库**: MySQL 8.0
- **缓存**: Redis
- **ORM**: GORM
- **认证**: JWT
- **架构模式**: DDD + 依赖注入 + DAO模式

### 前端技术栈
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI库**: Ant Design Vue
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP客户端**: Axios

### 高级DAO架构特性
- **读写分离**: 支持主从数据库配置
- **智能缓存**: 多种缓存策略(WriteThrough/WriteBack/WriteAround/CacheAside)
- **连接池管理**: 高效的数据库连接管理
- **事务管理**: 本地和分布式事务支持
- **查询构建器**: 流式API构建复杂SQL查询
- **负载均衡**: 多种负载均衡策略

## 📁 项目结构

```
Golang-HR/
├── backend/                 # 后端服务
│   ├── controllers/        # 控制器层
│   ├── services/          # 业务逻辑层
│   ├── dao/              # 数据访问层
│   │   ├── interfaces/   # DAO接口定义
│   │   ├── mysql/       # MySQL实现
│   │   ├── redis/       # Redis实现
│   │   ├── cache/       # 缓存策略
│   │   ├── pool/        # 连接池管理
│   │   ├── query/       # 查询构建器
│   │   └── transaction/ # 事务管理
│   ├── models/           # 数据模型
│   ├── middleware/       # 中间件
│   ├── routes/          # 路由定义
│   ├── config/          # 配置管理
│   └── utils/           # 工具函数
└── ui/                   # 前端应用
    ├── src/
    │   ├── components/   # 公共组件
    │   ├── views/       # 页面组件
    │   ├── stores/      # 状态管理
    │   ├── services/    # API服务
    │   ├── types/       # 类型定义
    │   └── utils/       # 工具函数
    ├── public/          # 静态资源
    └── package.json     # 依赖配置
```

## 🚀 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+

### 后端启动

```bash
cd backend

# 安装依赖
go mod tidy

# 配置环境变量
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=hr_system
export REDIS_HOST=localhost
export REDIS_PORT=6379

# 运行服务
go run main.go
```

### 前端启动

```bash
cd ui

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

## 📊 数据库设计

### 核心数据表

- **users**: 用户账户表
- **employees**: 员工信息表
- **departments**: 部门信息表
- **positions**: 职位信息表
- **job_levels**: 职级信息表
- **attendances**: 考勤记录表
- **leaves**: 请假记录表
- **salaries**: 薪资记录表
- **recruitments**: 招聘职位表
- **candidates**: 候选人信息表
- **performances**: 绩效考核表

## 🔧 API接口

### 员工管理
```
GET    /api/v1/employees          # 获取员工列表
POST   /api/v1/employees          # 创建员工
GET    /api/v1/employees/:id      # 获取员工详情
PUT    /api/v1/employees/:id      # 更新员工信息
DELETE /api/v1/employees/:id      # 删除员工
GET    /api/v1/employees/search   # 搜索员工
GET    /api/v1/employees/export   # 导出员工数据
POST   /api/v1/employees/import   # 导入员工数据
```

### 部门管理
```
GET    /api/v1/departments        # 获取部门列表
POST   /api/v1/departments        # 创建部门
GET    /api/v1/departments/:id    # 获取部门详情
PUT    /api/v1/departments/:id    # 更新部门信息
DELETE /api/v1/departments/:id    # 删除部门
GET    /api/v1/departments/tree   # 获取部门树形结构
```

## 🎯 功能特色

### 1. 高级DAO架构
- **接口抽象**: 通过接口层实现数据库无关性
- **读写分离**: 支持主从数据库配置，提升性能
- **智能缓存**: 自动缓存失效，支持多级缓存
- **连接池**: 高效的数据库连接管理
- **事务支持**: 完整的事务管理机制

### 2. 企业级特性
- **权限管理**: 基于RBAC的细粒度权限控制
- **数据导入导出**: 支持Excel/CSV格式
- **批量操作**: 支持批量更新和删除
- **审计日志**: 完整的操作记录
- **多租户**: 支持多租户架构

### 3. 现代化前端
- **响应式设计**: 适配各种设备尺寸
- **实时更新**: WebSocket实时通信
- **组件化开发**: 高度可复用的组件
- **类型安全**: 完整的TypeScript支持

## 🔐 安全特性

- JWT身份认证
- RBAC权限控制
- SQL注入防护
- XSS攻击防护
- CSRF令牌验证
- 敏感数据加密

## 📈 性能优化

- 数据库索引优化
- Redis缓存策略
- 连接池管理
- 分页查询
- 懒加载
- CDN加速

## 🤝 贡献指南

欢迎提交Issue和Pull Request来帮助改进项目。

## 📄 许可证

本项目采用 MIT 许可证。详情请见 [LICENSE](LICENSE) 文件。

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

---

**注意**: 这是一个演示项目，模仿北森云系统的功能特性，仅供学习和参考使用。