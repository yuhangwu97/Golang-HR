-- HR系统数据库初始化脚本 - 重新设计版本
-- 创建日期: 2025-07-23
-- 确保每个部门都有负责人，所有密码统一设置

-- 创建数据库
CREATE DATABASE IF NOT EXISTS gin_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE gin_db;

-- 删除现有表（如果存在）以重新创建
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS `organization_change_logs`;
DROP TABLE IF EXISTS `organization_histories`;
DROP TABLE IF EXISTS `organization_changes`;
DROP TABLE IF EXISTS `employee_assignments`;
DROP TABLE IF EXISTS `organization_units`;
DROP TABLE IF EXISTS `department_assignments`;
DROP TABLE IF EXISTS `performances`;
DROP TABLE IF EXISTS `candidates`;
DROP TABLE IF EXISTS `recruitments`;
DROP TABLE IF EXISTS `payroll_records`;
DROP TABLE IF EXISTS `salaries`;
DROP TABLE IF EXISTS `leaves`;
DROP TABLE IF EXISTS `attendances`;
DROP TABLE IF EXISTS `work_experiences`;
DROP TABLE IF EXISTS `employees`;
DROP TABLE IF EXISTS `positions`;
DROP TABLE IF EXISTS `job_levels`;
DROP TABLE IF EXISTS `departments`;
DROP TABLE IF EXISTS `role_permissions`;
DROP TABLE IF EXISTS `user_roles`;
DROP TABLE IF EXISTS `permissions`;
DROP TABLE IF EXISTS `roles`;
DROP TABLE IF EXISTS `users`;
SET FOREIGN_KEY_CHECKS = 1;

-- ==============================================
-- 1. 系统基础表
-- ==============================================

-- 用户表
CREATE TABLE `users` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
    `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `employee_id` bigint unsigned DEFAULT NULL,
    `status` enum('active','inactive','locked') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active',
    `last_login_at` datetime DEFAULT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`),
    UNIQUE KEY `idx_email` (`email`),
    KEY `idx_employee_id` (`employee_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 角色表
CREATE TABLE `roles` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
    `description` text COLLATE utf8mb4_unicode_ci,
    `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name` (`name`),
    UNIQUE KEY `idx_code` (`code`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 权限表
CREATE TABLE `permissions` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `code` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `resource` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `action` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
    `description` text COLLATE utf8mb4_unicode_ci,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_code` (`code`),
    KEY `idx_resource` (`resource`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 用户角色关联表
CREATE TABLE `user_roles` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id` bigint unsigned NOT NULL,
    `role_id` bigint unsigned NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_role` (`user_id`, `role_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 角色权限关联表
CREATE TABLE `role_permissions` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `role_id` bigint unsigned NOT NULL,
    `permission_id` bigint unsigned NOT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_role_permission` (`role_id`, `permission_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==============================================
-- 2. 职级表（需要在员工表之前创建）
-- ==============================================

CREATE TABLE `job_levels` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `code` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `level` int NOT NULL,
    `min_salary` decimal(10,2) DEFAULT NULL,
    `max_salary` decimal(10,2) DEFAULT NULL,
    `description` text COLLATE utf8mb4_unicode_ci,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_code` (`code`),
    UNIQUE KEY `idx_level` (`level`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==============================================
-- 3. 员工表（需要在部门表之前创建以避免循环依赖）
-- ==============================================

CREATE TABLE `employees` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `employee_id` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `avatar` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `gender` enum('male','female','other') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `birthday` date DEFAULT NULL,
    `id_card` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `department_id` bigint unsigned DEFAULT NULL,
    `position_id` bigint unsigned DEFAULT NULL,
    `job_level_id` bigint unsigned DEFAULT NULL,
    `manager_id` bigint unsigned DEFAULT NULL,
    `functional_manager_id` bigint unsigned DEFAULT NULL,
    `secondary_department_id` bigint unsigned DEFAULT NULL,
    `work_percentage` decimal(5,2) NOT NULL DEFAULT '100.00',
    `assignment_type` enum('primary','additional','temporary','project') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'primary',
    `management_type` enum('line','matrix','functional') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'line',
    `hire_date` date NOT NULL,
    `probation_end_date` date DEFAULT NULL,
    `contract_start_date` date DEFAULT NULL,
    `contract_end_date` date DEFAULT NULL,
    `contract_type` enum('full_time','part_time','contract','intern','consultant') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'full_time',
    `base_salary` decimal(10,2) DEFAULT NULL,
    `address` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `emergency_contact` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `emergency_phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `education` enum('high_school','associate','bachelor','master','phd') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `school` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `major` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `status` enum('active','inactive','on_leave','terminated') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_employee_id` (`employee_id`),
    UNIQUE KEY `idx_email` (`email`),
    KEY `idx_department_id` (`department_id`),
    KEY `idx_position_id` (`position_id`),
    KEY `idx_job_level_id` (`job_level_id`),
    KEY `idx_manager_id` (`manager_id`),
    KEY `idx_functional_manager_id` (`functional_manager_id`),
    KEY `idx_secondary_department_id` (`secondary_department_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==============================================
-- 4. 部门表
-- ==============================================

CREATE TABLE `departments` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
    `code` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `short_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `description` text COLLATE utf8mb4_unicode_ci,
    `parent_id` bigint unsigned DEFAULT NULL,
    `level` int NOT NULL DEFAULT '1',
    `sort` int NOT NULL DEFAULT '0',
    `type` enum('company','business_unit','department','team','cost_center','location','project') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'department',
    `manager_id` bigint unsigned NOT NULL, -- 修改为NOT NULL，确保每个部门都有负责人
    `functional_manager_id` bigint unsigned DEFAULT NULL,
    `address` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `website` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `country_code` char(2) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `currency_code` char(3) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `time_zone` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `cost_center` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `external_id` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `effective_date` date DEFAULT NULL,
    `expiration_date` date DEFAULT NULL,
    `is_active` tinyint(1) NOT NULL DEFAULT '1',
    `is_headquarters` tinyint(1) NOT NULL DEFAULT '0',
    `allow_subunits` tinyint(1) NOT NULL DEFAULT '1',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_code` (`code`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_manager_id` (`manager_id`),
    KEY `idx_functional_manager_id` (`functional_manager_id`),
    KEY `idx_type` (`type`),
    KEY `idx_is_active` (`is_active`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==============================================
-- 5. 职位表
-- ==============================================

CREATE TABLE `positions` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
    `code` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `description` text COLLATE utf8mb4_unicode_ci,
    `requirements` text COLLATE utf8mb4_unicode_ci,
    `parent_id` bigint unsigned DEFAULT NULL,
    `level` int NOT NULL DEFAULT '1',
    `sort` int NOT NULL DEFAULT '0',
    `department_id` bigint unsigned DEFAULT NULL,
    `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_code` (`code`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_department_id` (`department_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==============================================
-- 6. 薪资管理表
-- ==============================================

-- 薪资记录表
CREATE TABLE IF NOT EXISTS `salaries` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `employee_id` bigint unsigned NOT NULL,
    `month` date NOT NULL,
    `base_salary` decimal(10,2) NOT NULL,
    `overtime_pay` decimal(10,2) DEFAULT '0.00',
    `bonus` decimal(10,2) DEFAULT '0.00',
    `allowance` decimal(10,2) DEFAULT '0.00',
    `commission` decimal(10,2) DEFAULT '0.00',
    `deduction` decimal(10,2) DEFAULT '0.00',
    `tax` decimal(10,2) DEFAULT '0.00',
    `social_insurance` decimal(10,2) DEFAULT '0.00',
    `housing_fund` decimal(10,2) DEFAULT '0.00',
    `gross_salary` decimal(10,2) NOT NULL,
    `net_salary` decimal(10,2) NOT NULL,
    `payment_date` date DEFAULT NULL,
    `status` enum('draft','confirmed','paid','cancelled') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'draft',
    `notes` text COLLATE utf8mb4_unicode_ci,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_employee_month` (`employee_id`, `month`),
    KEY `idx_month` (`month`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`),
    CONSTRAINT `salaries_employee_id_foreign` FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 薪资发放记录表
CREATE TABLE IF NOT EXISTS `payroll_records` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `salary_id` bigint unsigned NOT NULL,
    `payment_method` enum('bank_transfer','cash','check','digital_wallet') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'bank_transfer',
    `bank_account` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `bank_name` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `transaction_id` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `amount` decimal(10,2) NOT NULL,
    `processor_id` bigint unsigned NOT NULL,
    `processed_at` datetime NOT NULL,
    `status` enum('pending','processed','failed','refunded') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending',
    `notes` text COLLATE utf8mb4_unicode_ci,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_salary_id` (`salary_id`),
    KEY `idx_processor_id` (`processor_id`),
    KEY `idx_status` (`status`),
    KEY `idx_processed_at` (`processed_at`),
    KEY `idx_deleted_at` (`deleted_at`),
    CONSTRAINT `payroll_records_salary_id_foreign` FOREIGN KEY (`salary_id`) REFERENCES `salaries` (`id`) ON DELETE CASCADE,
    CONSTRAINT `payroll_records_processor_id_foreign` FOREIGN KEY (`processor_id`) REFERENCES `employees` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==============================================
-- 7. 插入初始化数据
-- ==============================================

-- 插入职级数据
INSERT INTO `job_levels` (`id`, `name`, `code`, `level`, `min_salary`, `max_salary`, `description`) VALUES
(1, '实习生', 'INTERN', 1, 3000.00, 5000.00, '实习阶段员工'),
(2, '初级专员', 'JUNIOR', 2, 5000.00, 8000.00, '初级专业技术人员'),
(3, '专员', 'SPECIALIST', 3, 8000.00, 12000.00, '专业技术人员'),
(4, '高级专员', 'SENIOR_SPEC', 4, 12000.00, 18000.00, '高级专业技术人员'),
(5, '主管', 'SUPERVISOR', 5, 15000.00, 25000.00, '基层管理人员'),
(6, '经理', 'MANAGER', 6, 20000.00, 35000.00, '中层管理人员'),
(7, '高级经理', 'SENIOR_MGR', 7, 30000.00, 50000.00, '高级管理人员'),
(8, '总监', 'DIRECTOR', 8, 45000.00, 80000.00, '部门总监级别'),
(9, '副总裁', 'VP', 9, 70000.00, 120000.00, '副总裁级别'),
(10, 'C级高管', 'C_LEVEL', 10, 100000.00, 300000.00, 'CEO/CTO/CFO等高管');

-- 插入员工数据（确保每个部门都有专属管理者）
INSERT INTO `employees` (`id`, `employee_id`, `name`, `email`, `phone`, `gender`, `birthday`, `job_level_id`, `hire_date`, `contract_type`, `base_salary`, `education`, `status`) VALUES
-- 高层管理（C级）
(1, 'CEO001', '张伟', 'zhangwei@company.com', '13800138001', 'male', '1980-05-15', 10, '2020-01-01', 'full_time', 200000.00, 'master', 'active'),
(2, 'CTO001', '李娜', 'lina@company.com', '13800138002', 'female', '1985-08-22', 10, '2020-02-01', 'full_time', 180000.00, 'master', 'active'),
(3, 'CFO001', '王强', 'wangqiang@company.com', '13800138003', 'male', '1982-03-10', 10, '2020-03-01', 'full_time', 160000.00, 'master', 'active'),
(4, 'CHR001', '刘芳', 'liufang@company.com', '13800138004', 'female', '1987-11-28', 8, '2020-04-01', 'full_time', 80000.00, 'bachelor', 'active'),

-- 事业部总监级别
(5, 'DIR001', '陈明', 'chenming@company.com', '13800138005', 'male', '1988-07-14', 8, '2020-05-01', 'full_time', 75000.00, 'bachelor', 'active'),
(6, 'DIR002', '周丽', 'zhouli@company.com', '13800138006', 'female', '1990-12-03', 8, '2020-06-01', 'full_time', 70000.00, 'bachelor', 'active'),
(7, 'DIR003', '黄华', 'huanghua@company.com', '13800138007', 'male', '1989-09-17', 8, '2020-07-01', 'full_time', 65000.00, 'bachelor', 'active'),

-- 部门经理级别（为避免重复，新增专门的部门经理）
(8, 'MGR001', '赵敏', 'zhaomin@company.com', '13800138008', 'female', '1991-04-25', 6, '2020-08-01', 'full_time', 35000.00, 'bachelor', 'active'),
(9, 'MGR002', '孙杰', 'sunjie@company.com', '13800138009', 'male', '1992-01-12', 6, '2020-09-01', 'full_time', 32000.00, 'bachelor', 'active'),
(10, 'MGR003', '钱琳', 'qianlin@company.com', '13800138010', 'female', '1993-06-08', 6, '2020-10-01', 'full_time', 30000.00, 'bachelor', 'active'),
(11, 'MGR004', '郭飞', 'guofei@company.com', '13800138011', 'male', '1996-11-07', 6, '2021-05-01', 'full_time', 30000.00, 'master', 'active'),
(12, 'MGR005', '韩丽丽', 'hanlili@company.com', '13800138012', 'female', '1990-12-16', 6, '2021-08-01', 'full_time', 28000.00, 'bachelor', 'active'),
(13, 'MGR006', '冯宁', 'fengning@company.com', '13800138013', 'male', '1997-04-03', 6, '2021-09-01', 'full_time', 25000.00, 'bachelor', 'active'),

-- 新增专门的法务部和财务部经理
(14, 'LAW001', '吴静', 'wujing@company.com', '13800138014', 'female', '1988-03-15', 6, '2021-01-01', 'full_time', 32000.00, 'master', 'active'),
(15, 'FIN001', '马强', 'maqiang@company.com', '13800138015', 'male', '1985-11-20', 6, '2021-02-01', 'full_time', 35000.00, 'master', 'active'),

-- 团队负责人级别
(16, 'LEAD001', '李雪', 'lixue@company.com', '13800138016', 'female', '1994-02-18', 5, '2021-03-01', 'full_time', 22000.00, 'bachelor', 'active'),
(17, 'LEAD002', '张鹏', 'zhangpeng@company.com', '13800138017', 'male', '1993-08-14', 5, '2021-04-01', 'full_time', 23000.00, 'bachelor', 'active'),
(18, 'LEAD003', '王丽', 'wangli@company.com', '13800138018', 'female', '1992-05-20', 5, '2021-05-01', 'full_time', 24000.00, 'bachelor', 'active'),
(19, 'LEAD004', '刘波', 'liubo@company.com', '13800138019', 'male', '1995-10-30', 5, '2021-06-01', 'full_time', 21000.00, 'bachelor', 'active'),
(20, 'LEAD005', '陈霞', 'chenxia@company.com', '13800138020', 'female', '1994-07-12', 5, '2021-07-01', 'full_time', 20000.00, 'bachelor', 'active'),
(21, 'LEAD006', '杨明', 'yangming@company.com', '13800138021', 'male', '1996-09-08', 5, '2021-08-01', 'full_time', 22000.00, 'bachelor', 'active'),

-- 新增北京分公司总经理
(22, 'BJ001', '徐建国', 'xujianguo@company.com', '13800138022', 'male', '1983-06-25', 8, '2021-09-01', 'full_time', 60000.00, 'master', 'active'),

-- 普通员工
(23, 'ENG001', '吴涛', 'wutao@company.com', '13800138023', 'male', '1994-02-18', 4, '2021-01-01', 'full_time', 18000.00, 'bachelor', 'active'),
(24, 'ENG002', '胡燕', 'huyan@company.com', '13800138024', 'female', '1995-10-30', 4, '2021-02-01', 'full_time', 18000.00, 'bachelor', 'active'),
(25, 'ENG003', '林宇', 'linyu@company.com', '13800138025', 'male', '1993-08-14', 4, '2021-03-01', 'full_time', 18000.00, 'bachelor', 'active'),
(26, 'ENG004', '何雪', 'hexue@company.com', '13800138026', 'female', '1992-05-20', 3, '2021-04-01', 'full_time', 12000.00, 'bachelor', 'active'),
(27, 'ENG005', '邓梅', 'dengmei@company.com', '13800138027', 'female', '1994-03-22', 4, '2021-06-01', 'full_time', 16000.00, 'bachelor', 'active'),
(28, 'ENG006', '许刚', 'xugang@company.com', '13800138028', 'male', '1995-07-09', 3, '2021-07-01', 'full_time', 12000.00, 'bachelor', 'active'),
(29, 'HR001', '蔡红', 'caihong@company.com', '13800138029', 'female', '1991-08-11', 3, '2021-10-01', 'full_time', 12000.00, 'bachelor', 'active'),
(30, 'MKT001', '宋江', 'songjiang@company.com', '13800138030', 'male', '1992-12-05', 4, '2021-11-01', 'full_time', 15000.00, 'bachelor', 'active');

-- 插入部门数据（修正每个部门的负责人层级匹配）
INSERT INTO `departments` (`id`, `name`, `code`, `short_name`, `description`, `parent_id`, `level`, `sort`, `type`, `manager_id`, `country_code`, `currency_code`, `time_zone`, `is_active`, `is_headquarters`, `allow_subunits`) VALUES
-- 公司层级（第1层）：CEO管理
(1, '科技创新有限公司', 'COMP001', '科技创新', '一家专注于技术创新的现代化企业', NULL, 1, 1, 'company', 1, 'CN', 'CNY', 'Asia/Shanghai', 1, 1, 1),

-- 事业部层级（第2层）：C级高管管理
(2, '技术研发事业部', 'BU001', '技术研发', '负责公司核心技术研发工作', 1, 2, 1, 'business_unit', 2, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(3, '市场运营事业部', 'BU002', '市场运营', '负责市场推广和运营工作', 1, 2, 2, 'business_unit', 6, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(4, '职能支持中心', 'BU003', '职能支持', '提供人力、财务、法务等职能支持', 1, 2, 3, 'business_unit', 4, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(20, '北京分公司', 'LOC001', '北京分公司', '北京地区分公司', 1, 2, 4, 'location', 22, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),

-- 部门层级（第3层）：总监级管理
(5, '产品开发部', 'DEPT001', '产品开发', '负责产品设计和开发', 2, 3, 1, 'department', 5, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(6, '平台架构部', 'DEPT002', '平台架构', '负责技术平台和系统架构', 2, 3, 2, 'department', 7, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(7, '质量保证部', 'DEPT003', '质量保证', '负责产品质量测试和保证', 2, 3, 3, 'department', 8, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(8, '市场推广部', 'DEPT004', '市场推广', '负责产品市场推广和品牌建设', 3, 3, 1, 'department', 9, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(9, '销售部', 'DEPT005', '销售部', '负责产品销售和客户关系管理', 3, 3, 2, 'department', 10, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(10, '客户服务部', 'DEPT006', '客户服务', '负责客户服务和支持', 3, 3, 3, 'department', 11, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(11, '人力资源部', 'DEPT007', '人力资源', '负责人力资源管理和发展', 4, 3, 1, 'department', 12, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(12, '财务部', 'DEPT008', '财务部', '负责财务管理和会计核算', 4, 3, 2, 'department', 15, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(13, '法务部', 'DEPT009', '法务部', '负责法律事务和合规管理', 4, 3, 3, 'department', 14, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),

-- 团队层级（第4层）：团队负责人管理
(14, '前端开发团队', 'TEAM001', '前端团队', '负责前端界面开发', 5, 4, 1, 'team', 16, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 0),
(15, '后端开发团队', 'TEAM002', '后端团队', '负责后端服务开发', 5, 4, 2, 'team', 17, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 0),
(16, '移动开发团队', 'TEAM003', '移动团队', '负责移动应用开发', 5, 4, 3, 'team', 18, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 0),
(17, '数据团队', 'TEAM004', '数据团队', '负责数据分析和挖掘', 6, 4, 1, 'team', 19, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 0),
(18, '运维团队', 'TEAM005', '运维团队', '负责系统运维和部署', 6, 4, 2, 'team', 20, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 0),
(19, '测试团队', 'TEAM006', '测试团队', '负责功能和性能测试', 7, 4, 1, 'team', 21, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 0);

-- 更新员工的部门关联和管理关系（修正层级匹配）
-- 第1层：公司层级
UPDATE `employees` SET `department_id` = 1, `manager_id` = NULL WHERE `id` = 1; -- 张伟(CEO)：管理总公司，无上级

-- 第2层：事业部层级（C级高管直属CEO）
UPDATE `employees` SET `department_id` = 2, `manager_id` = 1 WHERE `id` = 2; -- 李娜(CTO)：管理技术研发事业部
UPDATE `employees` SET `department_id` = 4, `manager_id` = 1 WHERE `id` = 3; -- 王强(CFO)：在职能支持中心工作
UPDATE `employees` SET `department_id` = 4, `manager_id` = 1 WHERE `id` = 4; -- 刘芳(CHR)：管理职能支持中心
UPDATE `employees` SET `department_id` = 3, `manager_id` = 1 WHERE `id` = 6; -- 周丽(市场总监)：管理市场运营事业部
UPDATE `employees` SET `department_id` = 20, `manager_id` = 1 WHERE `id` = 22; -- 徐建国：管理北京分公司

-- 第3层：部门层级（总监级管理部门）
UPDATE `employees` SET `department_id` = 5, `manager_id` = 2 WHERE `id` = 5; -- 陈明(产品总监)：管理产品开发部
UPDATE `employees` SET `department_id` = 6, `manager_id` = 2 WHERE `id` = 7; -- 黄华(架构总监)：管理平台架构部
UPDATE `employees` SET `department_id` = 7, `manager_id` = 2 WHERE `id` = 8; -- 赵敏(QA经理)：管理质量保证部
UPDATE `employees` SET `department_id` = 8, `manager_id` = 6 WHERE `id` = 9; -- 孙杰(市场经理)：管理市场推广部
UPDATE `employees` SET `department_id` = 9, `manager_id` = 6 WHERE `id` = 10; -- 钱琳(销售经理)：管理销售部
UPDATE `employees` SET `department_id` = 10, `manager_id` = 6 WHERE `id` = 11; -- 郭飞(客服经理)：管理客户服务部
UPDATE `employees` SET `department_id` = 11, `manager_id` = 4 WHERE `id` = 12; -- 韩丽丽(HR经理)：管理人力资源部
UPDATE `employees` SET `department_id` = 13, `manager_id` = 4 WHERE `id` = 14; -- 吴静(法务经理)：管理法务部
UPDATE `employees` SET `department_id` = 12, `manager_id` = 4 WHERE `id` = 15; -- 马强(财务经理)：管理财务部

-- 第4层：团队层级（团队负责人管理团队）
UPDATE `employees` SET `department_id` = 14, `manager_id` = 5 WHERE `id` = 16; -- 李雪(前端负责人)：管理前端开发团队
UPDATE `employees` SET `department_id` = 15, `manager_id` = 5 WHERE `id` = 17; -- 张鹏(后端负责人)：管理后端开发团队
UPDATE `employees` SET `department_id` = 16, `manager_id` = 5 WHERE `id` = 18; -- 王丽(移动负责人)：管理移动开发团队
UPDATE `employees` SET `department_id` = 17, `manager_id` = 7 WHERE `id` = 19; -- 刘波(数据负责人)：管理数据团队
UPDATE `employees` SET `department_id` = 18, `manager_id` = 7 WHERE `id` = 20; -- 陈霞(运维负责人)：管理运维团队
UPDATE `employees` SET `department_id` = 19, `manager_id` = 8 WHERE `id` = 21; -- 杨明(测试负责人)：管理测试团队

-- 第5层：普通员工（归属各团队，直属团队负责人）
UPDATE `employees` SET `department_id` = 14, `manager_id` = 16 WHERE `id` = 23; -- 吴涛(前端工程师)：归属前端团队
UPDATE `employees` SET `department_id` = 15, `manager_id` = 17 WHERE `id` = 24; -- 胡燕(后端工程师)：归属后端团队
UPDATE `employees` SET `department_id` = 16, `manager_id` = 18 WHERE `id` = 25; -- 林宇(移动工程师)：归属移动团队
UPDATE `employees` SET `department_id` = 5, `manager_id` = 5 WHERE `id` = 26; -- 何雪(产品专员)：归属产品开发部
UPDATE `employees` SET `department_id` = 19, `manager_id` = 21 WHERE `id` = 27; -- 邓梅(测试工程师)：归属测试团队
UPDATE `employees` SET `department_id` = 18, `manager_id` = 20 WHERE `id` = 28; -- 许刚(运维工程师)：归属运维团队
UPDATE `employees` SET `department_id` = 11, `manager_id` = 12 WHERE `id` = 29; -- 蔡红(HR专员)：归属人力资源部
UPDATE `employees` SET `department_id` = 8, `manager_id` = 9 WHERE `id` = 30; -- 宋江(市场专员)：归属市场推广部

-- 修正一个特殊情况：冯宁作为副总监
UPDATE `employees` SET `department_id` = 3, `manager_id` = 6 WHERE `id` = 13; -- 冯宁(副总监)：协助市场总监管理市场运营事业部

-- 修正关键的部门负责人问题
-- 1. 科技创新有限公司的负责人应该是张伟(ID=1)
UPDATE `departments` SET `manager_id` = 1 WHERE `id` = 1;
-- 2. 前端开发团队的负责人应该是李雪(ID=16) 
UPDATE `departments` SET `manager_id` = 16 WHERE `id` = 14;
-- 3. 北京分公司的负责人应该是徐建国(ID=22)
UPDATE `departments` SET `manager_id` = 22 WHERE `id` = 20;

-- 确保员工部门分配正确
UPDATE `employees` SET `department_id` = 1 WHERE `id` = 1;   -- 张伟在总公司
UPDATE `employees` SET `department_id` = 14 WHERE `id` = 16; -- 李雪在前端开发团队
UPDATE `employees` SET `department_id` = 20 WHERE `id` = 22; -- 徐建国在北京分公司

-- 插入职位数据
INSERT INTO `positions` (`id`, `name`, `code`, `description`, `requirements`, `parent_id`, `level`, `sort`, `department_id`, `status`) VALUES
(1, '首席执行官', 'CEO', '公司最高管理者，负责公司整体战略和运营', '具有丰富的企业管理经验，优秀的领导能力', NULL, 1, 1, 1, 'active'),
(2, '首席技术官', 'CTO', '负责公司技术战略和技术团队管理', '具有深厚的技术背景和管理经验', 1, 2, 1, 2, 'active'),
(3, '首席财务官', 'CFO', '负责公司财务战略和财务管理', '具有财务管理和投资经验', 1, 2, 2, 12, 'active'),
(4, '人力资源总监', 'CHR', '负责人力资源战略和HR团队管理', '具有人力资源管理经验', 1, 2, 3, 11, 'active'),
(5, '产品总监', 'PROD_DIR', '负责产品战略和产品团队管理', '具有产品管理经验和市场洞察力', 2, 3, 1, 5, 'active'),
(6, '市场总监', 'MKT_DIR', '负责市场战略和营销团队管理', '具有市场营销经验和品牌管理能力', 1, 3, 2, 3, 'active'),
(7, '技术架构总监', 'ARCH_DIR', '负责技术架构和基础设施管理', '具有系统架构和技术管理经验', 2, 3, 3, 6, 'active'),
(8, '开发经理', 'DEV_MGR', '负责开发团队管理和项目执行', '具有软件开发和团队管理经验', 5, 4, 1, 14, 'active'),
(9, '产品经理', 'PROD_MGR', '负责产品规划和需求管理', '具有产品管理经验和用户洞察', 5, 4, 2, 15, 'active'),
(10, '技术经理', 'TECH_MGR', '负责技术团队管理', '具有技术管理经验', 5, 4, 3, 16, 'active'),
(11, '高级开发工程师', 'SR_DEV', '负责核心功能开发和技术攻关', '具有丰富的开发经验和技术能力', 8, 5, 1, 14, 'active'),
(12, '开发工程师', 'DEV_ENG', '负责功能开发和代码实现', '具有软件开发能力和学习能力', 11, 6, 1, 15, 'active'),
(13, '移动开发工程师', 'MOB_DEV', '负责移动应用开发', '具有移动开发经验', 10, 6, 1, 16, 'active'),
(14, '产品专员', 'PROD_SPEC', '协助产品规划和需求分析', '产品相关专业，有一定产品经验', 9, 6, 1, 5, 'active'),
(15, '测试经理', 'TEST_MGR', '负责测试团队管理和质量保证', '具有测试管理经验和质量意识', 7, 4, 4, 19, 'active'),
(16, '测试工程师', 'TEST_ENG', '负责功能测试和缺陷发现', '具有测试技能和细心负责', 15, 5, 1, 19, 'active'),
(17, '运维经理', 'OPS_MGR', '负责系统运维和基础设施管理', '具有运维管理经验和技术能力', 7, 4, 5, 18, 'active'),
(18, '市场经理', 'MKT_MGR', '负责市场推广和品牌建设', '具有市场营销经验', 6, 4, 1, 8, 'active'),
(19, '人力资源专员', 'HR_SPEC', '负责招聘、培训等HR事务', '人力资源相关专业，沟通能力强', 4, 5, 1, 11, 'active');

-- 更新员工的职位关联
UPDATE `employees` SET `position_id` = 1 WHERE `id` = 1; -- CEO
UPDATE `employees` SET `position_id` = 2 WHERE `id` = 2; -- CTO
UPDATE `employees` SET `position_id` = 3 WHERE `id` = 3; -- CFO
UPDATE `employees` SET `position_id` = 4 WHERE `id` = 4; -- CHR
UPDATE `employees` SET `position_id` = 5 WHERE `id` = 5; -- 产品总监
UPDATE `employees` SET `position_id` = 6 WHERE `id` = 6; -- 市场总监
UPDATE `employees` SET `position_id` = 7 WHERE `id` = 7; -- 架构总监
UPDATE `employees` SET `position_id` = 8 WHERE `id` = 8; -- 开发经理
UPDATE `employees` SET `position_id` = 9 WHERE `id` = 9; -- 产品经理
UPDATE `employees` SET `position_id` = 10 WHERE `id` = 10; -- 技术经理
UPDATE `employees` SET `position_id` = 11 WHERE `id` = 11; -- 高级开发工程师
UPDATE `employees` SET `position_id` = 12 WHERE `id` = 12; -- 开发工程师
UPDATE `employees` SET `position_id` = 13 WHERE `id` = 13; -- 移动开发工程师
UPDATE `employees` SET `position_id` = 14 WHERE `id` = 14; -- 产品专员
UPDATE `employees` SET `position_id` = 15 WHERE `id` = 15; -- 测试经理
UPDATE `employees` SET `position_id` = 16 WHERE `id` = 16; -- 测试工程师
UPDATE `employees` SET `position_id` = 16 WHERE `id` = 17; -- 测试工程师
UPDATE `employees` SET `position_id` = 17 WHERE `id` = 18; -- 运维经理
UPDATE `employees` SET `position_id` = 18 WHERE `id` = 19; -- 市场经理
UPDATE `employees` SET `position_id` = 19 WHERE `id` = 20; -- HR专员

-- 插入角色数据
INSERT INTO `roles` (`id`, `name`, `code`, `description`, `status`) VALUES
(1, '超级管理员', 'super_admin', '系统超级管理员，拥有所有权限', 'active'),
(2, '系统管理员', 'admin', '系统管理员，负责系统配置和用户管理', 'active'),
(3, 'HR管理员', 'hr_admin', 'HR管理员，负责人力资源管理', 'active'),
(4, 'HR专员', 'hr_specialist', 'HR专员，处理日常人力资源事务', 'active'),
(5, '部门经理', 'dept_manager', '部门经理，管理部门员工和事务', 'active'),
(6, '项目经理', 'project_manager', '项目经理，负责项目管理', 'active'),
(7, '团队主管', 'team_lead', '团队主管，负责团队管理', 'active'),
(8, '高级员工', 'senior_employee', '高级员工，有一定管理权限', 'active'),
(9, '普通员工', 'employee', '普通员工，基本权限', 'active');

-- 插入权限数据  
INSERT INTO `permissions` (`id`, `name`, `code`, `resource`, `action`, `description`) VALUES
(1, '查看用户', 'user.view', 'user', 'view', '查看用户信息'),
(2, '创建用户', 'user.create', 'user', 'create', '创建新用户'),
(3, '编辑用户', 'user.edit', 'user', 'edit', '编辑用户信息'),
(4, '删除用户', 'user.delete', 'user', 'delete', '删除用户'),
(5, '查看员工', 'employee.view', 'employee', 'view', '查看员工信息'),
(6, '创建员工', 'employee.create', 'employee', 'create', '创建员工档案'),
(7, '编辑员工', 'employee.edit', 'employee', 'edit', '编辑员工信息'),
(8, '删除员工', 'employee.delete', 'employee', 'delete', '删除员工档案'),
(9, '查看部门', 'department.view', 'department', 'view', '查看部门信息'),
(10, '创建部门', 'department.create', 'department', 'create', '创建部门'),
(11, '编辑部门', 'department.edit', 'department', 'edit', '编辑部门信息'),
(12, '删除部门', 'department.delete', 'department', 'delete', '删除部门');

-- 插入用户数据（使用指定的密码哈希）
INSERT INTO `users` (`id`, `username`, `email`, `password`, `employee_id`, `status`) VALUES
(1, 'zhangwei', 'zhangwei@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 1, 'active'),
(2, 'lina', 'lina@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 2, 'active'),
(3, 'wangqiang', 'wangqiang@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 3, 'active'),
(4, 'liufang', 'liufang@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 4, 'active'),
(5, 'chenming', 'chenming@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 5, 'active'),
(6, 'zhouli', 'zhouli@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 6, 'active'),
(7, 'huanghua', 'huanghua@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 7, 'active'),
(8, 'zhaomin', 'zhaomin@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 8, 'active'),
(9, 'sunjie', 'sunjie@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 9, 'active'),
(10, 'qianlin', 'qianlin@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 10, 'active'),
(11, 'guofei', 'guofei@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 11, 'active'),
(12, 'hanlili', 'hanlili@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 12, 'active'),
(13, 'fengning', 'fengning@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 13, 'active'),
(14, 'wujing', 'wujing@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 14, 'active'),
(15, 'maqiang', 'maqiang@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 15, 'active'),
(16, 'lixue', 'lixue@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 16, 'active'),
(17, 'zhangpeng', 'zhangpeng@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 17, 'active'),
(18, 'wangli', 'wangli@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 18, 'active'),
(19, 'liubo', 'liubo@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 19, 'active'),
(20, 'chenxia', 'chenxia@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 20, 'active'),
(21, 'yangming', 'yangming@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 21, 'active'),
(22, 'xujianguo', 'xujianguo@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 22, 'active'),
(23, 'wutao', 'wutao@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 23, 'active'),
(24, 'huyan', 'huyan@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 24, 'active'),
(25, 'linyu', 'linyu@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 25, 'active'),
(26, 'hexue', 'hexue@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 26, 'active'),
(27, 'dengmei', 'dengmei@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 27, 'active'),
(28, 'xugang', 'xugang@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 28, 'active'),
(29, 'caihong', 'caihong@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 29, 'active'),
(30, 'songjiang', 'songjiang@company.com', '$2a$10$34E5nzYEgcRxXbxASDmHderhljgdm0Cr4RwvPEoYL8M9BdUXI.Jai', 30, 'active');

-- 插入用户角色关联
INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES
(1, 1), -- CEO - 超级管理员
(2, 2), -- CTO - 系统管理员  
(3, 2), -- CFO - 系统管理员
(4, 3), -- CHR - HR管理员
(5, 5), -- 产品总监 - 部门经理
(6, 5), -- 市场总监 - 部门经理
(7, 5), -- 架构总监 - 部门经理
(8, 5), -- 客户服务部经理 - 部门经理
(9, 5), -- 人力资源部经理 - 部门经理
(10, 5), -- 市场运营副经理 - 部门经理
(11, 5), -- 质量保证部经理 - 部门经理
(12, 5), -- 市场推广部经理 - 部门经理
(13, 5), -- 销售部经理 - 部门经理
(14, 5), -- 法务部经理 - 部门经理
(15, 5), -- 财务部经理 - 部门经理
(16, 7), -- 前端团队负责人 - 团队主管
(17, 7), -- 后端团队负责人 - 团队主管
(18, 7), -- 移动团队负责人 - 团队主管
(19, 7), -- 数据团队负责人 - 团队主管
(20, 7), -- 运维团队负责人 - 团队主管
(21, 7), -- 测试团队负责人 - 团队主管
(22, 5), -- 北京分公司总经理 - 部门经理
(23, 9), -- 前端工程师 - 普通员工
(24, 9), -- 后端工程师 - 普通员工
(25, 9), -- 移动工程师 - 普通员工
(26, 9), -- 产品专员 - 普通员工
(27, 9), -- 测试工程师 - 普通员工
(28, 9), -- 运维工程师 - 普通员工
(29, 4), -- HR专员 - HR专员
(30, 9); -- 市场专员 - 普通员工

-- 插入角色权限关联
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
-- 超级管理员拥有所有权限
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8), (1, 9), (1, 10), (1, 11), (1, 12),
-- 系统管理员权限
(2, 1), (2, 2), (2, 3), (2, 5), (2, 6), (2, 7), (2, 9), (2, 10), (2, 11),
-- HR管理员权限
(3, 5), (3, 6), (3, 7), (3, 9),
-- HR专员权限
(4, 5), (4, 9),
-- 部门经理权限
(5, 5), (5, 7), (5, 9),
-- 项目经理权限
(6, 5), (6, 9),
-- 团队主管权限
(7, 5), (7, 9),
-- 高级员工权限
(8, 5), (8, 9),
-- 普通员工权限
(9, 5), (9, 9);

-- ==============================================
-- 7. 添加外键约束
-- ==============================================

ALTER TABLE `users` ADD CONSTRAINT `fk_users_employee` FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;
ALTER TABLE `user_roles` ADD CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
ALTER TABLE `user_roles` ADD CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE;
ALTER TABLE `role_permissions` ADD CONSTRAINT `fk_role_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE;
ALTER TABLE `role_permissions` ADD CONSTRAINT `fk_role_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE;

ALTER TABLE `departments` ADD CONSTRAINT `fk_departments_parent` FOREIGN KEY (`parent_id`) REFERENCES `departments` (`id`) ON DELETE SET NULL;
ALTER TABLE `departments` ADD CONSTRAINT `fk_departments_manager` FOREIGN KEY (`manager_id`) REFERENCES `employees` (`id`) ON DELETE RESTRICT;
ALTER TABLE `departments` ADD CONSTRAINT `fk_departments_functional_manager` FOREIGN KEY (`functional_manager_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;

ALTER TABLE `positions` ADD CONSTRAINT `fk_positions_parent` FOREIGN KEY (`parent_id`) REFERENCES `positions` (`id`) ON DELETE SET NULL;
ALTER TABLE `positions` ADD CONSTRAINT `fk_positions_department` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`) ON DELETE SET NULL;

ALTER TABLE `employees` ADD CONSTRAINT `fk_employees_department` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`) ON DELETE SET NULL;
ALTER TABLE `employees` ADD CONSTRAINT `fk_employees_position` FOREIGN KEY (`position_id`) REFERENCES `positions` (`id`) ON DELETE SET NULL;
ALTER TABLE `employees` ADD CONSTRAINT `fk_employees_job_level` FOREIGN KEY (`job_level_id`) REFERENCES `job_levels` (`id`) ON DELETE SET NULL;
ALTER TABLE `employees` ADD CONSTRAINT `fk_employees_manager` FOREIGN KEY (`manager_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;
ALTER TABLE `employees` ADD CONSTRAINT `fk_employees_functional_manager` FOREIGN KEY (`functional_manager_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;
ALTER TABLE `employees` ADD CONSTRAINT `fk_employees_secondary_department` FOREIGN KEY (`secondary_department_id`) REFERENCES `departments` (`id`) ON DELETE SET NULL;

-- 插入示例薪资数据
INSERT INTO `salaries` (`employee_id`, `month`, `base_salary`, `overtime_pay`, `bonus`, `allowance`, `commission`, `deduction`, `tax`, `social_insurance`, `housing_fund`, `gross_salary`, `net_salary`, `status`) VALUES
-- 2024年12月薪资数据
(1, '2024-12-01', 200000.00, 0.00, 50000.00, 5000.00, 0.00, 2000.00, 35000.00, 8000.00, 4000.00, 255000.00, 206000.00, 'paid'),
(2, '2024-12-01', 180000.00, 2000.00, 30000.00, 4000.00, 0.00, 1500.00, 28000.00, 7200.00, 3600.00, 216000.00, 175700.00, 'paid'),
(3, '2024-12-01', 160000.00, 0.00, 25000.00, 3500.00, 0.00, 1000.00, 24000.00, 6400.00, 3200.00, 188500.00, 153900.00, 'paid'),
(4, '2024-12-01', 80000.00, 1000.00, 15000.00, 2000.00, 0.00, 500.00, 12000.00, 3200.00, 1600.00, 98000.00, 80700.00, 'paid'),
(5, '2024-12-01', 75000.00, 1500.00, 12000.00, 2500.00, 0.00, 800.00, 11000.00, 3000.00, 1500.00, 91000.00, 74700.00, 'paid'),
(6, '2024-12-01', 70000.00, 800.00, 10000.00, 2000.00, 0.00, 600.00, 10000.00, 2800.00, 1400.00, 82800.00, 68000.00, 'paid'),
(7, '2024-12-01', 65000.00, 1200.00, 8000.00, 1800.00, 0.00, 400.00, 9000.00, 2600.00, 1300.00, 76000.00, 62700.00, 'paid'),
(8, '2024-12-01', 35000.00, 2000.00, 5000.00, 1500.00, 0.00, 300.00, 5500.00, 1400.00, 700.00, 43500.00, 35600.00, 'paid'),
(9, '2024-12-01', 32000.00, 1800.00, 4500.00, 1200.00, 0.00, 250.00, 5000.00, 1280.00, 640.00, 39500.00, 32330.00, 'paid'),
(10, '2024-12-01', 30000.00, 1500.00, 4000.00, 1000.00, 0.00, 200.00, 4500.00, 1200.00, 600.00, 36500.00, 30000.00, 'paid'),
-- 2025年1月薪资数据（当前月）
(1, '2025-01-01', 200000.00, 0.00, 40000.00, 5000.00, 0.00, 2000.00, 34000.00, 8000.00, 4000.00, 245000.00, 197000.00, 'confirmed'),
(2, '2025-01-01', 180000.00, 1500.00, 25000.00, 4000.00, 0.00, 1500.00, 27000.00, 7200.00, 3600.00, 210500.00, 171200.00, 'confirmed'),
(3, '2025-01-01', 160000.00, 0.00, 20000.00, 3500.00, 0.00, 1000.00, 23000.00, 6400.00, 3200.00, 183500.00, 149900.00, 'confirmed'),
(4, '2025-01-01', 80000.00, 800.00, 12000.00, 2000.00, 0.00, 500.00, 11500.00, 3200.00, 1600.00, 94800.00, 77700.00, 'confirmed'),
(5, '2025-01-01', 75000.00, 1200.00, 10000.00, 2500.00, 0.00, 800.00, 10500.00, 3000.00, 1500.00, 88700.00, 72900.00, 'confirmed'),
(6, '2025-01-01', 70000.00, 600.00, 8000.00, 2000.00, 0.00, 600.00, 9500.00, 2800.00, 1400.00, 80600.00, 66300.00, 'confirmed'),
(7, '2025-01-01', 65000.00, 1000.00, 6000.00, 1800.00, 0.00, 400.00, 8500.00, 2600.00, 1300.00, 73800.00, 61000.00, 'confirmed'),
(8, '2025-01-01', 35000.00, 1800.00, 4000.00, 1500.00, 0.00, 300.00, 5200.00, 1400.00, 700.00, 42300.00, 34700.00, 'confirmed'),
(9, '2025-01-01', 32000.00, 1600.00, 3500.00, 1200.00, 0.00, 250.00, 4800.00, 1280.00, 640.00, 38300.00, 31330.00, 'confirmed'),
(10, '2025-01-01', 30000.00, 1400.00, 3000.00, 1000.00, 0.00, 200.00, 4300.00, 1200.00, 600.00, 35400.00, 29100.00, 'confirmed');

-- ==============================================
-- 8. 创建索引以提升查询性能
-- ==============================================

CREATE INDEX `idx_employees_dept_status` ON `employees` (`department_id`, `status`);
CREATE INDEX `idx_employees_manager_status` ON `employees` (`manager_id`, `status`);
CREATE INDEX `idx_departments_manager_active` ON `departments` (`manager_id`, `is_active`);
CREATE INDEX `idx_departments_parent_level` ON `departments` (`parent_id`, `level`);

-- ==============================================
-- 9. 验证数据一致性
-- ==============================================

-- 验证每个部门都有负责人
SELECT 
    d.id,
    d.name,
    d.manager_id,
    e.name as manager_name,
    e.employee_id as manager_employee_id
FROM departments d
LEFT JOIN employees e ON d.manager_id = e.id
WHERE d.is_active = 1
ORDER BY d.level, d.sort;

-- 验证员工管理关系
SELECT 
    emp.id,
    emp.employee_id,
    emp.name as employee_name,
    emp.manager_id,
    mgr.name as manager_name,
    dept.name as department_name
FROM employees emp
LEFT JOIN employees mgr ON emp.manager_id = mgr.id  
LEFT JOIN departments dept ON emp.department_id = dept.id
WHERE emp.status = 'active'
ORDER BY emp.department_id, emp.id;

-- 统计信息
SELECT 
    '总部门数' as item,
    COUNT(*) as count
FROM departments 
WHERE is_active = 1
UNION ALL
SELECT 
    '有负责人的部门数' as item,
    COUNT(*) as count  
FROM departments 
WHERE is_active = 1 AND manager_id IS NOT NULL
UNION ALL
SELECT 
    '总员工数' as item,
    COUNT(*) as count
FROM employees 
WHERE status = 'active'
UNION ALL
SELECT 
    '有管理者的员工数' as item,
    COUNT(*) as count
FROM employees 
WHERE status = 'active' AND manager_id IS NOT NULL;

-- 初始化完成提示
SELECT '=== 数据库初始化完成 ===' as message;
SELECT '所有部门都已分配负责人' as status;
SELECT '所有用户密码已统一设置' as password_status;