-- HR系统数据库初始化脚本
-- 创建日期: 2024-07-22
-- 包含所有表的创建语句和初始化数据

-- 创建数据库
CREATE DATABASE IF NOT EXISTS gin_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE gin_db;

-- ==============================================
-- 1. 系统基础表
-- ==============================================

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
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
CREATE TABLE IF NOT EXISTS `roles` (
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
CREATE TABLE IF NOT EXISTS `permissions` (
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
CREATE TABLE IF NOT EXISTS `user_roles` (
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
CREATE TABLE IF NOT EXISTS `role_permissions` (
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
-- 2. 组织架构表
-- ==============================================

-- 部门表
CREATE TABLE IF NOT EXISTS `departments` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
    `code` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `short_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `description` text COLLATE utf8mb4_unicode_ci,
    `parent_id` bigint unsigned DEFAULT NULL,
    `level` int NOT NULL DEFAULT '1',
    `sort` int NOT NULL DEFAULT '0',
    `type` enum('company','business_unit','department','team','cost_center','location','project') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'department',
    `manager_id` bigint unsigned DEFAULT NULL,
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

-- 职位表
CREATE TABLE IF NOT EXISTS `positions` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
    `code` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
    `description` text COLLATE utf8mb4_unicode_ci,
    `requirements` text COLLATE utf8mb4_unicode_ci,
    `parent_id` bigint unsigned DEFAULT NULL,
    `level` int NOT NULL DEFAULT '1',
    `sort` int NOT NULL DEFAULT '0',
    `department_id` bigint unsigned DEFAULT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_code` (`code`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_department_id` (`department_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 职级表
CREATE TABLE IF NOT EXISTS `job_levels` (
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

-- 员工表
CREATE TABLE IF NOT EXISTS `employees` (
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

-- 部门分配表
CREATE TABLE IF NOT EXISTS `department_assignments` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `employee_id` bigint unsigned NOT NULL,
    `department_id` bigint unsigned NOT NULL,
    `assignment_type` enum('primary','additional','temporary','project') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'primary',
    `management_type` enum('line','matrix','functional') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'line',
    `work_percentage` decimal(5,2) NOT NULL DEFAULT '100.00',
    `is_primary` tinyint(1) NOT NULL DEFAULT '0',
    `position_id` bigint unsigned DEFAULT NULL,
    `job_level_id` bigint unsigned DEFAULT NULL,
    `direct_manager_id` bigint unsigned DEFAULT NULL,
    `functional_manager_id` bigint unsigned DEFAULT NULL,
    `effective_date` date NOT NULL,
    `expiration_date` date DEFAULT NULL,
    `status` enum('active','inactive','expired') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_employee_id` (`employee_id`),
    KEY `idx_department_id` (`department_id`),
    KEY `idx_position_id` (`position_id`),
    KEY `idx_job_level_id` (`job_level_id`),
    KEY `idx_direct_manager_id` (`direct_manager_id`),
    KEY `idx_functional_manager_id` (`functional_manager_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==============================================
-- 3. 组织架构扩展表
-- ==============================================

-- 组织单元表
CREATE TABLE IF NOT EXISTS `organization_units` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
    `code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
    `description` text COLLATE utf8mb4_unicode_ci,
    `parent_id` bigint unsigned DEFAULT NULL,
    `level` int NOT NULL DEFAULT '1',
    `sort` int NOT NULL DEFAULT '0',
    `type` enum('company','business_unit','department','team','cost_center','location','project') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'department',
    `manager_id` bigint unsigned DEFAULT NULL,
    `functional_manager_id` bigint unsigned DEFAULT NULL,
    `cost_center` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `effective_date` date DEFAULT NULL,
    `expiration_date` date DEFAULT NULL,
    `is_active` tinyint(1) NOT NULL DEFAULT '1',
    `status` enum('active','inactive','planning','dissolved') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active',
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

-- 员工分配表
CREATE TABLE IF NOT EXISTS `employee_assignments` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `employee_id` bigint unsigned NOT NULL,
    `organization_unit_id` bigint unsigned NOT NULL,
    `assignment_type` enum('primary','additional','temporary','project') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'primary',
    `management_type` enum('line','matrix','functional') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'line',
    `work_percentage` decimal(5,2) NOT NULL DEFAULT '100.00',
    `is_primary` tinyint(1) NOT NULL DEFAULT '0',
    `position_id` bigint unsigned DEFAULT NULL,
    `job_level_id` bigint unsigned DEFAULT NULL,
    `direct_manager_id` bigint unsigned DEFAULT NULL,
    `functional_manager_id` bigint unsigned DEFAULT NULL,
    `effective_date` date NOT NULL,
    `expiration_date` date DEFAULT NULL,
    `status` enum('active','inactive','expired') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_employee_id` (`employee_id`),
    KEY `idx_organization_unit_id` (`organization_unit_id`),
    KEY `idx_position_id` (`position_id`),
    KEY `idx_job_level_id` (`job_level_id`),
    KEY `idx_direct_manager_id` (`direct_manager_id`),
    KEY `idx_functional_manager_id` (`functional_manager_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==============================================
-- 4. 组织变更管理表
-- ==============================================

-- 组织变更表
CREATE TABLE IF NOT EXISTS `organization_changes` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `entity_type` enum('department','position','employee','organization_unit') COLLATE utf8mb4_unicode_ci NOT NULL,
    `entity_id` bigint unsigned NOT NULL,
    `change_type` enum('create','update','delete','move','assign','merge','split') COLLATE utf8mb4_unicode_ci NOT NULL,
    `field_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `old_value` text COLLATE utf8mb4_unicode_ci,
    `new_value` text COLLATE utf8mb4_unicode_ci,
    `description` text COLLATE utf8mb4_unicode_ci,
    `reason` text COLLATE utf8mb4_unicode_ci,
    `impact_analysis` text COLLATE utf8mb4_unicode_ci,
    `initiator_id` bigint unsigned NOT NULL,
    `approver_id` bigint unsigned DEFAULT NULL,
    `approval_note` text COLLATE utf8mb4_unicode_ci,
    `effective_date` datetime DEFAULT NULL,
    `status` enum('pending','approved','rejected','implemented','cancelled') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_entity` (`entity_type`, `entity_id`),
    KEY `idx_change_type` (`change_type`),
    KEY `idx_initiator_id` (`initiator_id`),
    KEY `idx_approver_id` (`approver_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 组织历史表
CREATE TABLE IF NOT EXISTS `organization_histories` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `snapshot_date` datetime NOT NULL,
    `snapshot_reason` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `unit_id` bigint unsigned NOT NULL,
    `unit_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
    `unit_code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
    `unit_type` enum('company','business_unit','department','team','cost_center','location','project') COLLATE utf8mb4_unicode_ci NOT NULL,
    `parent_id` bigint unsigned DEFAULT NULL,
    `parent_name` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `level` int NOT NULL,
    `manager_id` bigint unsigned DEFAULT NULL,
    `manager_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `functional_manager_id` bigint unsigned DEFAULT NULL,
    `functional_manager_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `employee_count` int NOT NULL DEFAULT '0',
    `direct_reports` int NOT NULL DEFAULT '0',
    `subunit_count` int NOT NULL DEFAULT '0',
    `status` enum('active','inactive','planning','dissolved') COLLATE utf8mb4_unicode_ci NOT NULL,
    `is_active` tinyint(1) NOT NULL,
    `change_type` enum('create','update','delete','move','assign','merge','split') COLLATE utf8mb4_unicode_ci NOT NULL,
    `changed_by` bigint unsigned NOT NULL,
    `changed_by_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `hierarchy_path` varchar(1000) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_unit_id` (`unit_id`),
    KEY `idx_snapshot_date` (`snapshot_date`),
    KEY `idx_change_type` (`change_type`),
    KEY `idx_changed_by` (`changed_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 组织变更日志表
CREATE TABLE IF NOT EXISTS `organization_change_logs` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `entity_type` enum('department','position','employee','organization_unit') COLLATE utf8mb4_unicode_ci NOT NULL,
    `entity_id` bigint unsigned NOT NULL,
    `action` enum('create','update','delete','move','assign','merge','split') COLLATE utf8mb4_unicode_ci NOT NULL,
    `old_value` text COLLATE utf8mb4_unicode_ci,
    `new_value` text COLLATE utf8mb4_unicode_ci,
    `field_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `operator_id` bigint unsigned NOT NULL,
    `operator_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `reason` text COLLATE utf8mb4_unicode_ci,
    `ip_address` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `user_agent` text COLLATE utf8mb4_unicode_ci,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_entity` (`entity_type`, `entity_id`),
    KEY `idx_action` (`action`),
    KEY `idx_operator_id` (`operator_id`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==============================================
-- 5. 业务管理表
-- ==============================================

-- 工作经历表
CREATE TABLE IF NOT EXISTS `work_experiences` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `employee_id` bigint unsigned NOT NULL,
    `company` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
    `position` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
    `start_date` date NOT NULL,
    `end_date` date DEFAULT NULL,
    `description` text COLLATE utf8mb4_unicode_ci,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_employee_id` (`employee_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 考勤表
CREATE TABLE IF NOT EXISTS `attendances` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `employee_id` bigint unsigned NOT NULL,
    `date` date NOT NULL,
    `check_in_time` time DEFAULT NULL,
    `check_out_time` time DEFAULT NULL,
    `work_hours` decimal(4,2) DEFAULT NULL,
    `break_hours` decimal(4,2) DEFAULT NULL,
    `overtime_hours` decimal(4,2) DEFAULT NULL,
    `status` enum('present','absent','late','early_leave','half_day','holiday','sick_leave') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'present',
    `notes` text COLLATE utf8mb4_unicode_ci,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_employee_date` (`employee_id`, `date`),
    KEY `idx_date` (`date`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 请假表
CREATE TABLE IF NOT EXISTS `leaves` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `employee_id` bigint unsigned NOT NULL,
    `type` enum('annual','sick','personal','maternity','paternity','bereavement','emergency','unpaid') COLLATE utf8mb4_unicode_ci NOT NULL,
    `start_date` date NOT NULL,
    `end_date` date NOT NULL,
    `days` decimal(4,2) NOT NULL,
    `reason` text COLLATE utf8mb4_unicode_ci,
    `approver_id` bigint unsigned DEFAULT NULL,
    `approval_note` text COLLATE utf8mb4_unicode_ci,
    `status` enum('pending','approved','rejected','cancelled') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_employee_id` (`employee_id`),
    KEY `idx_approver_id` (`approver_id`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`),
    KEY `idx_start_date` (`start_date`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 薪资表
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
    KEY `idx_deleted_at` (`deleted_at`)
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
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==============================================
-- 6. 招聘管理表
-- ==============================================

-- 招聘表
CREATE TABLE IF NOT EXISTS `recruitments` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
    `department_id` bigint unsigned NOT NULL,
    `position_id` bigint unsigned NOT NULL,
    `job_level_id` bigint unsigned DEFAULT NULL,
    `recruiter_id` bigint unsigned NOT NULL,
    `urgency` enum('low','medium','high','urgent') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'medium',
    `headcount` int NOT NULL DEFAULT '1',
    `salary_min` decimal(10,2) DEFAULT NULL,
    `salary_max` decimal(10,2) DEFAULT NULL,
    `requirements` text COLLATE utf8mb4_unicode_ci,
    `description` text COLLATE utf8mb4_unicode_ci,
    `benefits` text COLLATE utf8mb4_unicode_ci,
    `work_location` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `work_type` enum('full_time','part_time','contract','intern','remote') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'full_time',
    `start_date` date NOT NULL,
    `end_date` date DEFAULT NULL,
    `status` enum('active','paused','closed','cancelled') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_department_id` (`department_id`),
    KEY `idx_position_id` (`position_id`),
    KEY `idx_recruiter_id` (`recruiter_id`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 候选人表
CREATE TABLE IF NOT EXISTS `candidates` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `recruitment_id` bigint unsigned NOT NULL,
    `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `resume_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `cover_letter` text COLLATE utf8mb4_unicode_ci,
    `experience_years` int DEFAULT NULL,
    `current_company` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `current_position` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `current_salary` decimal(10,2) DEFAULT NULL,
    `expected_salary` decimal(10,2) DEFAULT NULL,
    `source` enum('internal','referral','job_board','social_media','recruitment_agency','career_fair','other') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `referrer_id` bigint unsigned DEFAULT NULL,
    `stage` enum('applied','screening','interview','assessment','offer','hired','rejected') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'applied',
    `status` enum('active','withdrawn','rejected','hired','on_hold') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active',
    `notes` text COLLATE utf8mb4_unicode_ci,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_recruitment_id` (`recruitment_id`),
    KEY `idx_referrer_id` (`referrer_id`),
    KEY `idx_stage` (`stage`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 绩效评估表
CREATE TABLE IF NOT EXISTS `performances` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `employee_id` bigint unsigned NOT NULL,
    `evaluator_id` bigint unsigned NOT NULL,
    `evaluation_period` enum('monthly','quarterly','semi_annual','annual') COLLATE utf8mb4_unicode_ci NOT NULL,
    `period_start` date NOT NULL,
    `period_end` date NOT NULL,
    `overall_score` decimal(4,2) DEFAULT NULL,
    `performance_rating` enum('outstanding','exceeds','meets','below','unsatisfactory') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `goals_achievement` text COLLATE utf8mb4_unicode_ci,
    `strengths` text COLLATE utf8mb4_unicode_ci,
    `areas_for_improvement` text COLLATE utf8mb4_unicode_ci,
    `development_plan` text COLLATE utf8mb4_unicode_ci,
    `evaluator_comments` text COLLATE utf8mb4_unicode_ci,
    `employee_comments` text COLLATE utf8mb4_unicode_ci,
    `status` enum('draft','submitted','reviewed','approved','rejected') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'draft',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_employee_id` (`employee_id`),
    KEY `idx_evaluator_id` (`evaluator_id`),
    KEY `idx_period` (`period_start`, `period_end`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==============================================
-- 7. 插入初始化数据
-- ==============================================

-- 角色数据
INSERT INTO `roles` (`id`, `name`, `code`, `description`, `status`) VALUES
(1, '超级管理员', 'super_admin', '系统超级管理员，拥有所有权限', 'active'),
(2, '系统管理员', 'admin', '系统管理员，负责系统配置和用户管理', 'active'),
(3, 'HR管理员', 'hr_admin', 'HR管理员，负责人力资源管理', 'active'),
(4, 'HR专员', 'hr_specialist', 'HR专员，处理日常人力资源事务', 'active'),
(5, '部门经理', 'dept_manager', '部门经理，管理部门员工和事务', 'active'),
(6, '项目经理', 'project_manager', '项目经理，负责项目管理', 'active'),
(7, '团队主管', 'team_lead', '团队主管，负责团队管理', 'active'),
(8, '高级员工', 'senior_employee', '高级员工，有一定管理权限', 'active'),
(9, '普通员工', 'employee', '普通员工，基本权限', 'active'),
(10, '实习生', 'intern', '实习生，受限权限', 'active'),
(11, '顾问', 'consultant', '外部顾问，特定项目权限', 'active'),
(12, '财务管理员', 'finance_admin', '财务管理员，负责财务相关功能', 'active'),
(13, '财务专员', 'finance_specialist', '财务专员，处理财务事务', 'active'),
(14, 'IT管理员', 'it_admin', 'IT管理员，负责技术支持', 'active'),
(15, '审计员', 'auditor', '内部审计员，负责审计工作', 'active'),
(16, '法务专员', 'legal_specialist', '法务专员，处理法律事务', 'active'),
(17, '安全管理员', 'security_admin', '安全管理员，负责信息安全', 'active'),
(18, '访客', 'guest', '临时访客权限', 'active'),
(19, '只读用户', 'readonly_user', '只读权限用户', 'active'),
(20, '离职员工', 'former_employee', '离职员工，保留历史数据访问', 'inactive');

-- 权限数据
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
(12, '删除部门', 'department.delete', 'department', 'delete', '删除部门'),
(13, '查看薪资', 'salary.view', 'salary', 'view', '查看薪资信息'),
(14, '编辑薪资', 'salary.edit', 'salary', 'edit', '编辑薪资信息'),
(15, '查看考勤', 'attendance.view', 'attendance', 'view', '查看考勤记录'),
(16, '编辑考勤', 'attendance.edit', 'attendance', 'edit', '编辑考勤记录'),
(17, '查看请假', 'leave.view', 'leave', 'view', '查看请假记录'),
(18, '审批请假', 'leave.approve', 'leave', 'approve', '审批请假申请'),
(19, '查看招聘', 'recruitment.view', 'recruitment', 'view', '查看招聘信息'),
(20, '管理招聘', 'recruitment.manage', 'recruitment', 'manage', '管理招聘流程');

-- 职级数据
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
(10, '总裁', 'PRESIDENT', 10, 100000.00, 200000.00, '总裁级别'),
(11, '首席执行官', 'CEO', 11, 150000.00, 300000.00, '首席执行官'),
(12, '首席技术官', 'CTO', 11, 150000.00, 280000.00, '首席技术官'),
(13, '首席财务官', 'CFO', 11, 150000.00, 280000.00, '首席财务官'),
(14, '首席运营官', 'COO', 11, 150000.00, 280000.00, '首席运营官'),
(15, '首席人力资源官', 'CHRO', 11, 120000.00, 250000.00, '首席人力资源官'),
(16, '资深顾问', 'SR_CONSULTANT', 8, 50000.00, 100000.00, '资深外部顾问'),
(17, '首席顾问', 'CHIEF_CONSULTANT', 9, 80000.00, 150000.00, '首席外部顾问'),
(18, '合伙人', 'PARTNER', 10, 120000.00, 250000.00, '合伙人级别'),
(19, '临时工', 'TEMP', 1, 2000.00, 4000.00, '临时工作人员'),
(20, '兼职', 'PART_TIME', 2, 3000.00, 8000.00, '兼职工作人员');

-- 部门数据 (按层级顺序插入)
INSERT INTO `departments` (`id`, `name`, `code`, `short_name`, `description`, `parent_id`, `level`, `sort`, `type`, `manager_id`, `country_code`, `currency_code`, `time_zone`, `is_active`, `is_headquarters`, `allow_subunits`) VALUES
(1, '科技创新有限公司', 'COMP001', '科技创新', '一家专注于技术创新的现代化企业', NULL, 1, 1, 'company', 1, 'CN', 'CNY', 'Asia/Shanghai', 1, 1, 1),
(2, '技术研发事业部', 'BU001', '技术研发', '负责公司核心技术研发工作', 1, 2, 1, 'business_unit', 2, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(3, '市场运营事业部', 'BU002', '市场运营', '负责市场推广和运营工作', 1, 2, 2, 'business_unit', 6, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(4, '职能支持中心', 'BU003', '职能支持', '提供人力、财务、法务等职能支持', 1, 2, 3, 'business_unit', 7, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(5, '产品开发部', 'DEPT001', '产品开发', '负责产品设计和开发', 2, 3, 1, 'department', 5, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(6, '平台架构部', 'DEPT002', '平台架构', '负责技术平台和系统架构', 2, 3, 2, 'department', 4, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(7, '质量保证部', 'DEPT003', '质量保证', '负责产品质量测试和保证', 2, 3, 3, 'department', 15, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(8, '市场推广部', 'DEPT004', '市场推广', '负责产品市场推广和品牌建设', 3, 3, 1, 'department', 6, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(9, '销售部', 'DEPT005', '销售部', '负责产品销售和客户关系管理', 3, 3, 2, 'department', 19, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(10, '客户服务部', 'DEPT006', '客户服务', '负责客户服务和支持', 3, 3, 3, 'department', 6, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(11, '人力资源部', 'DEPT007', '人力资源', '负责人力资源管理和发展', 4, 3, 1, 'department', 7, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(12, '财务部', 'DEPT008', '财务部', '负责财务管理和会计核算', 4, 3, 2, 'department', 3, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(13, '法务部', 'DEPT009', '法务部', '负责法律事务和合规管理', 4, 3, 3, 'department', 1, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1),
(14, '前端开发团队', 'TEAM001', '前端团队', '负责前端界面开发', 5, 4, 1, 'team', 8, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 0),
(15, '后端开发团队', 'TEAM002', '后端团队', '负责后端服务开发', 5, 4, 2, 'team', 8, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 0),
(16, '移动开发团队', 'TEAM003', '移动团队', '负责移动应用开发', 5, 4, 3, 'team', 8, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 0),
(17, '数据团队', 'TEAM004', '数据团队', '负责数据分析和挖掘', 6, 4, 1, 'team', 4, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 0),
(18, '运维团队', 'TEAM005', '运维团队', '负责系统运维和部署', 6, 4, 2, 'team', 18, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 0),
(19, '测试团队', 'TEAM006', '测试团队', '负责功能和性能测试', 7, 4, 1, 'team', 15, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 0),
(20, '北京分公司', 'LOC001', '北京分公司', '北京地区分公司', 1, 2, 4, 'location', 1, 'CN', 'CNY', 'Asia/Shanghai', 1, 0, 1);

-- 职位数据
INSERT INTO `positions` (`id`, `name`, `code`, `description`, `requirements`, `parent_id`, `level`, `sort`, `department_id`) VALUES
(1, '首席执行官', 'CEO', '公司最高管理者，负责公司整体战略和运营', '具有丰富的企业管理经验，优秀的领导能力', NULL, 1, 1, 1),
(2, '首席技术官', 'CTO', '负责公司技术战略和技术团队管理', '具有深厚的技术背景和管理经验', 1, 2, 1, 2),
(3, '首席财务官', 'CFO', '负责公司财务战略和财务管理', '具有财务管理和投资经验', 1, 2, 2, 12),
(4, '技术总监', 'TECH_DIR', '负责技术部门管理和技术决策', '资深技术专家，具有团队管理能力', 2, 3, 1, 2),
(5, '产品总监', 'PROD_DIR', '负责产品战略和产品团队管理', '具有产品管理经验和市场洞察力', 2, 3, 2, 5),
(6, '市场总监', 'MKT_DIR', '负责市场战略和营销团队管理', '具有市场营销经验和品牌管理能力', NULL, 3, 1, 8),
(7, '人力资源总监', 'HR_DIR', '负责人力资源战略和HR团队管理', '具有人力资源管理经验', NULL, 3, 2, 11),
(8, '开发经理', 'DEV_MGR', '负责开发团队管理和项目执行', '具有软件开发和团队管理经验', 4, 4, 1, 5),
(9, '产品经理', 'PROD_MGR', '负责产品规划和需求管理', '具有产品管理经验和用户洞察', 5, 4, 1, 5),
(10, '高级开发工程师', 'SR_DEV', '负责核心功能开发和技术攻关', '具有丰富的开发经验和技术能力', 8, 5, 1, 5),
(11, '开发工程师', 'DEV_ENG', '负责功能开发和代码实现', '具有软件开发能力和学习能力', 10, 6, 1, 5),
(12, '初级开发工程师', 'JR_DEV', '协助功能开发和代码维护', '计算机相关专业，有一定编程基础', 11, 7, 1, 5),
(13, '测试经理', 'TEST_MGR', '负责测试团队管理和质量保证', '具有测试管理经验和质量意识', 4, 4, 2, 7),
(14, '高级测试工程师', 'SR_TEST', '负责测试方案设计和执行', '具有丰富的测试经验', 13, 5, 1, 7),
(15, '测试工程师', 'TEST_ENG', '负责功能测试和缺陷发现', '具有测试技能和细心负责', 14, 6, 1, 7),
(16, '运维经理', 'OPS_MGR', '负责系统运维和基础设施管理', '具有运维管理经验和技术能力', 4, 4, 3, 6),
(17, '高级运维工程师', 'SR_OPS', '负责系统部署和运维优化', '具有丰富的运维经验', 16, 5, 1, 6),
(18, '销售经理', 'SALES_MGR', '负责销售团队管理和业绩达成', '具有销售管理经验和客户资源', 6, 4, 1, 9),
(19, '人力资源专员', 'HR_SPEC', '负责招聘、培训等HR事务', '人力资源相关专业，沟通能力强', 7, 5, 1, 11),
(20, '财务专员', 'FIN_SPEC', '负责财务核算和报表制作', '财务会计相关专业，细心负责', NULL, 5, 1, 12);

-- 员工数据
INSERT INTO `employees` (`id`, `employee_id`, `name`, `email`, `phone`, `gender`, `birthday`, `department_id`, `position_id`, `job_level_id`, `manager_id`, `hire_date`, `contract_type`, `base_salary`, `education`, `status`) VALUES
(1, 'EMP001', '张伟', 'zhangwei@company.com', '13800138001', 'male', '1980-05-15', 1, 1, 11, NULL, '2020-01-01', 'full_time', 200000.00, 'master', 'active'),
(2, 'EMP002', '李娜', 'lina@company.com', '13800138002', 'female', '1985-08-22', 2, 2, 12, 1, '2020-02-01', 'full_time', 180000.00, 'master', 'active'),
(3, 'EMP003', '王强', 'wangqiang@company.com', '13800138003', 'male', '1982-03-10', 12, 3, 13, 1, '2020-03-01', 'full_time', 160000.00, 'master', 'active'),
(4, 'EMP004', '刘芳', 'liufang@company.com', '13800138004', 'female', '1987-11-28', 2, 4, 8, 2, '2020-04-01', 'full_time', 80000.00, 'bachelor', 'active'),
(5, 'EMP005', '陈明', 'chenming@company.com', '13800138005', 'male', '1988-07-14', 5, 5, 8, 2, '2020-05-01', 'full_time', 75000.00, 'bachelor', 'active'),
(6, 'EMP006', '周丽', 'zhouli@company.com', '13800138006', 'female', '1990-12-03', 8, 6, 8, 1, '2020-06-01', 'full_time', 70000.00, 'bachelor', 'active'),
(7, 'EMP007', '黄华', 'huanghua@company.com', '13800138007', 'male', '1989-09-17', 11, 7, 8, 1, '2020-07-01', 'full_time', 65000.00, 'bachelor', 'active'),
(8, 'EMP008', '赵敏', 'zhaomin@company.com', '13800138008', 'female', '1991-04-25', 5, 8, 6, 5, '2020-08-01', 'full_time', 35000.00, 'bachelor', 'active'),
(9, 'EMP009', '孙杰', 'sunjie@company.com', '13800138009', 'male', '1992-01-12', 5, 9, 6, 5, '2020-09-01', 'full_time', 32000.00, 'bachelor', 'active'),
(10, 'EMP010', '钱琳', 'qianlin@company.com', '13800138010', 'female', '1993-06-08', 5, 10, 4, 8, '2020-10-01', 'full_time', 18000.00, 'bachelor', 'active'),
(11, 'EMP011', '吴涛', 'wutao@company.com', '13800138011', 'male', '1994-02-18', 14, 11, 3, 8, '2021-01-01', 'full_time', 12000.00, 'bachelor', 'active'),
(12, 'EMP012', '胡燕', 'huyan@company.com', '13800138012', 'female', '1995-10-30', 15, 11, 3, 8, '2021-02-01', 'full_time', 12000.00, 'bachelor', 'active'),
(13, 'EMP013', '林宇', 'linyu@company.com', '13800138013', 'male', '1993-08-14', 16, 11, 3, 8, '2021-03-01', 'full_time', 12000.00, 'bachelor', 'active'),
(14, 'EMP014', '何雪', 'hexue@company.com', '13800138014', 'female', '1992-05-20', 5, 12, 2, 10, '2021-04-01', 'full_time', 8000.00, 'bachelor', 'active'),
(15, 'EMP015', '郭飞', 'guofei@company.com', '13800138015', 'male', '1996-11-07', 7, 13, 6, 4, '2021-05-01', 'full_time', 30000.00, 'master', 'active'),
(16, 'EMP016', '邓梅', 'dengmei@company.com', '13800138016', 'female', '1994-03-22', 19, 14, 4, 15, '2021-06-01', 'full_time', 16000.00, 'bachelor', 'active'),
(17, 'EMP017', '许刚', 'xugang@company.com', '13800138017', 'male', '1995-07-09', 19, 15, 3, 16, '2021-07-01', 'full_time', 10000.00, 'bachelor', 'active'),
(18, 'EMP018', '韩丽丽', 'hanlili@company.com', '13800138018', 'female', '1990-12-16', 18, 16, 6, 4, '2021-08-01', 'full_time', 28000.00, 'bachelor', 'active'),
(19, 'EMP019', '冯宁', 'fengning@company.com', '13800138019', 'male', '1997-04-03', 9, 18, 6, 6, '2021-09-01', 'full_time', 25000.00, 'bachelor', 'active'),
(20, 'EMP020', '蔡红', 'caihong@company.com', '13800138020', 'female', '1991-08-11', 11, 19, 3, 7, '2021-10-01', 'full_time', 9000.00, 'bachelor', 'active');

-- 用户数据 (基于员工创建对应用户)
INSERT INTO `users` (`id`, `username`, `email`, `password`, `employee_id`, `status`) VALUES
(1, 'zhangwei', 'zhangwei@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 1, 'active'),
(2, 'lina', 'lina@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 2, 'active'),
(3, 'wangqiang', 'wangqiang@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 3, 'active'),
(4, 'liufang', 'liufang@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 4, 'active'),
(5, 'chenming', 'chenming@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 5, 'active'),
(6, 'zhouli', 'zhouli@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 6, 'active'),
(7, 'huanghua', 'huanghua@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 7, 'active'),
(8, 'zhaomin', 'zhaomin@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 8, 'active'),
(9, 'sunjie', 'sunjie@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 9, 'active'),
(10, 'qianlin', 'qianlin@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 10, 'active'),
(11, 'wutao', 'wutao@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 11, 'active'),
(12, 'huyan', 'huyan@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 12, 'active'),
(13, 'linyu', 'linyu@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 13, 'active'),
(14, 'hexue', 'hexue@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 14, 'active'),
(15, 'guofei', 'guofei@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 15, 'active'),
(16, 'dengmei', 'dengmei@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 16, 'active'),
(17, 'xugang', 'xugang@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 17, 'active'),
(18, 'hanlili', 'hanlili@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 18, 'active'),
(19, 'fengning', 'fengning@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 19, 'active'),
(20, 'caihong', 'caihong@company.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKyZFa.LNK6mCdwKewCdpT2RLH4m', 20, 'active');

-- 用户角色关联
INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES
(1, 1), (2, 2), (3, 12), (4, 2), (5, 2),
(6, 5), (7, 3), (8, 5), (9, 6), (10, 8),
(11, 9), (12, 9), (13, 9), (14, 9), (15, 5),
(16, 8), (17, 9), (18, 8), (19, 8), (20, 4);

-- 角色权限关联 (超级管理员拥有所有权限)
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8), (1, 9), (1, 10),
(1, 11), (1, 12), (1, 13), (1, 14), (1, 15), (1, 16), (1, 17), (1, 18), (1, 19), (1, 20),
-- 系统管理员权限
(2, 1), (2, 2), (2, 3), (2, 5), (2, 6), (2, 7), (2, 9), (2, 10), (2, 11),
-- HR管理员权限
(3, 5), (3, 6), (3, 7), (3, 9), (3, 13), (3, 15), (3, 17), (3, 18), (3, 19), (3, 20),
-- 普通员工权限
(9, 5), (9, 9), (9, 15), (9, 17);

-- 考勤数据
INSERT INTO `attendances` (`employee_id`, `date`, `check_in_time`, `check_out_time`, `work_hours`, `status`) VALUES
(1, '2024-07-01', '09:00:00', '18:00:00', 8.00, 'present'),
(2, '2024-07-01', '09:15:00', '18:30:00', 8.25, 'late'),
(3, '2024-07-01', '08:45:00', '17:45:00', 8.00, 'present'),
(4, '2024-07-01', '09:00:00', '18:00:00', 8.00, 'present'),
(5, '2024-07-01', '09:30:00', '18:00:00', 7.50, 'late'),
(6, '2024-07-01', '09:00:00', '17:30:00', 7.50, 'early_leave'),
(7, '2024-07-01', '09:00:00', '18:00:00', 8.00, 'present'),
(8, '2024-07-01', '09:00:00', '18:00:00', 8.00, 'present'),
(9, '2024-07-01', NULL, NULL, 0.00, 'absent'),
(10, '2024-07-01', '09:00:00', '18:00:00', 8.00, 'present'),
(11, '2024-07-01', '09:00:00', '18:00:00', 8.00, 'present'),
(12, '2024-07-01', '09:05:00', '18:05:00', 8.00, 'present'),
(13, '2024-07-01', '09:00:00', '18:00:00', 8.00, 'present'),
(14, '2024-07-01', '09:00:00', '18:00:00', 8.00, 'present'),
(15, '2024-07-01', '08:50:00', '18:10:00', 8.33, 'present'),
(16, '2024-07-01', '09:00:00', '18:00:00', 8.00, 'present'),
(17, '2024-07-01', '09:00:00', '18:00:00', 8.00, 'present'),
(18, '2024-07-01', '09:00:00', '18:00:00', 8.00, 'present'),
(19, '2024-07-01', '09:20:00', '18:00:00', 7.67, 'late'),
(20, '2024-07-01', '09:00:00', '18:00:00', 8.00, 'present');

-- 薪资数据
INSERT INTO `salaries` (`employee_id`, `month`, `base_salary`, `bonus`, `allowance`, `deduction`, `tax`, `social_insurance`, `housing_fund`, `gross_salary`, `net_salary`, `status`) VALUES
(1, '2024-07-01', 200000.00, 20000.00, 5000.00, 0.00, 45000.00, 8000.00, 6000.00, 225000.00, 166000.00, 'confirmed'),
(2, '2024-07-01', 180000.00, 15000.00, 4000.00, 0.00, 39800.00, 7200.00, 5400.00, 199000.00, 146600.00, 'confirmed'),
(3, '2024-07-01', 160000.00, 12000.00, 3000.00, 0.00, 35000.00, 6400.00, 4800.00, 175000.00, 128800.00, 'confirmed'),
(4, '2024-07-01', 80000.00, 5000.00, 2000.00, 0.00, 17400.00, 3200.00, 2400.00, 87000.00, 64000.00, 'confirmed'),
(5, '2024-07-01', 75000.00, 4000.00, 1500.00, 0.00, 16100.00, 3000.00, 2250.00, 80500.00, 59150.00, 'confirmed'),
(6, '2024-07-01', 70000.00, 3500.00, 1000.00, 0.00, 14900.00, 2800.00, 2100.00, 74500.00, 54700.00, 'confirmed'),
(7, '2024-07-01', 65000.00, 3000.00, 1000.00, 0.00, 13800.00, 2600.00, 1950.00, 69000.00, 50650.00, 'confirmed'),
(8, '2024-07-01', 35000.00, 2000.00, 500.00, 0.00, 7500.00, 1400.00, 1050.00, 37500.00, 27550.00, 'confirmed'),
(9, '2024-07-01', 32000.00, 1500.00, 500.00, 0.00, 6800.00, 1280.00, 960.00, 34000.00, 24960.00, 'confirmed'),
(10, '2024-07-01', 18000.00, 1000.00, 300.00, 0.00, 3828.00, 720.00, 540.00, 19300.00, 14212.00, 'confirmed'),
(11, '2024-07-01', 12000.00, 500.00, 200.00, 0.00, 2540.00, 480.00, 360.00, 12700.00, 9320.00, 'confirmed'),
(12, '2024-07-01', 12000.00, 500.00, 200.00, 0.00, 2540.00, 480.00, 360.00, 12700.00, 9320.00, 'confirmed'),
(13, '2024-07-01', 12000.00, 800.00, 200.00, 0.00, 2600.00, 480.00, 360.00, 13000.00, 9560.00, 'confirmed'),
(14, '2024-07-01', 8000.00, 300.00, 100.00, 0.00, 1660.00, 320.00, 240.00, 8400.00, 6180.00, 'confirmed'),
(15, '2024-07-01', 30000.00, 2000.00, 500.00, 0.00, 6500.00, 1200.00, 900.00, 32500.00, 23900.00, 'confirmed'),
(16, '2024-07-01', 16000.00, 800.00, 200.00, 0.00, 3360.00, 640.00, 480.00, 17000.00, 12520.00, 'confirmed'),
(17, '2024-07-01', 10000.00, 500.00, 100.00, 0.00, 2120.00, 400.00, 300.00, 10600.00, 7780.00, 'confirmed'),
(18, '2024-07-01', 28000.00, 1500.00, 400.00, 0.00, 5956.00, 1120.00, 840.00, 29900.00, 21984.00, 'confirmed'),
(19, '2024-07-01', 25000.00, 1200.00, 300.00, 0.00, 5304.00, 1000.00, 750.00, 26500.00, 19446.00, 'confirmed'),
(20, '2024-07-01', 9000.00, 400.00, 100.00, 0.00, 1900.00, 360.00, 270.00, 9500.00, 6970.00, 'confirmed');

-- 请假数据
INSERT INTO `leaves` (`employee_id`, `type`, `start_date`, `end_date`, `days`, `reason`, `approver_id`, `status`) VALUES
(10, 'annual', '2024-07-15', '2024-07-19', 5.00, '年假休息', 8, 'approved'),
(14, 'sick', '2024-07-08', '2024-07-09', 2.00, '感冒发烧', 8, 'approved'),
(17, 'personal', '2024-07-22', '2024-07-22', 1.00, '个人事务', 15, 'pending'),
(12, 'annual', '2024-08-01', '2024-08-07', 7.00, '暑假旅游', 8, 'approved'),
(11, 'sick', '2024-07-05', '2024-07-05', 1.00, '医院检查', 8, 'approved'),
(20, 'annual', '2024-07-25', '2024-07-26', 2.00, '家庭聚会', 7, 'approved'),
(16, 'personal', '2024-07-30', '2024-07-30', 1.00, '银行办事', 15, 'approved'),
(13, 'annual', '2024-08-10', '2024-08-16', 7.00, '年假出游', 8, 'pending'),
(19, 'sick', '2024-07-12', '2024-07-13', 2.00, '肠胃不适', 6, 'approved'),
(15, 'annual', '2024-07-28', '2024-08-02', 6.00, '年假休息', 4, 'approved'),
(9, 'personal', '2024-07-20', '2024-07-20', 1.00, '驾照考试', 8, 'approved'),
(18, 'annual', '2024-08-05', '2024-08-09', 5.00, '回老家探亲', 4, 'pending'),
(6, 'sick', '2024-07-18', '2024-07-18', 1.00, '牙科治疗', 2, 'approved'),
(7, 'annual', '2024-08-12', '2024-08-18', 7.00, '海外旅游', 1, 'pending'),
(4, 'personal', '2024-07-24', '2024-07-24', 1.00, '孩子家长会', 2, 'approved'),
(5, 'annual', '2024-08-15', '2024-08-21', 7.00, '年假休息', 2, 'pending'),
(8, 'sick', '2024-07-16', '2024-07-17', 2.00, '感冒休息', 4, 'approved'),
(1, 'annual', '2024-09-01', '2024-09-07', 7.00, 'CEO年假', NULL, 'approved'),
(2, 'personal', '2024-07-29', '2024-07-29', 1.00, '技术会议', 1, 'approved'),
(3, 'annual', '2024-08-20', '2024-08-26', 7.00, 'CFO年假', 1, 'pending');

-- 工作经历数据
INSERT INTO `work_experiences` (`employee_id`, `company`, `position`, `start_date`, `end_date`, `description`) VALUES
(1, '阿里巴巴集团', '技术总监', '2015-01-01', '2019-12-31', '负责电商平台技术架构设计'),
(2, '腾讯科技', '高级架构师', '2017-03-01', '2020-01-31', '负责微信支付系统架构'),
(3, '华为技术', '财务经理', '2016-06-01', '2020-02-29', '负责企业级财务管理'),
(4, '百度', '技术主管', '2018-01-01', '2020-03-31', '负责搜索引擎技术优化'),
(5, '字节跳动', '产品经理', '2019-01-01', '2020-04-30', '负责短视频产品规划'),
(6, '美团', '市场经理', '2017-05-01', '2020-05-31', '负责本地生活服务推广'),
(7, '滴滴出行', 'HR经理', '2018-03-01', '2020-06-30', '负责人才招聘和培养'),
(8, '京东', '开发经理', '2019-02-01', '2020-07-31', '负责电商系统开发'),
(9, '小米科技', '产品经理', '2018-06-01', '2020-08-31', '负责智能硬件产品'),
(10, '网易', '高级工程师', '2019-03-01', '2020-09-30', '负责游戏后端开发'),
(11, '搜狐', '前端工程师', '2020-01-01', '2020-12-31', '负责门户网站前端开发'),
(12, '新浪', '后端工程师', '2020-02-01', '2021-01-31', '负责微博系统开发'),
(13, '快手', '移动开发工程师', '2020-03-01', '2021-02-28', '负责短视频APP开发'),
(14, '拼多多', '初级工程师', '2020-06-01', '2021-03-31', '负责电商平台功能开发'),
(15, '携程', '测试经理', '2019-05-01', '2021-04-30', '负责旅游平台质量保证'),
(16, '去哪儿', '高级测试工程师', '2020-01-01', '2021-05-31', '负责在线旅游测试'),
(17, '蚂蚁金服', '测试工程师', '2020-07-01', '2021-06-30', '负责金融系统测试'),
(18, '顺丰科技', '运维工程师', '2019-08-01', '2021-07-31', '负责物流系统运维'),
(19, '苏宁易购', '销售专员', '2020-01-01', '2021-08-31', '负责电商平台销售'),
(20, '58同城', 'HR专员', '2020-10-01', '2021-09-30', '负责招聘和员工关系管理');

-- 组织单元数据 (Organization Units)
INSERT INTO `organization_units` (`id`, `name`, `code`, `description`, `parent_id`, `level`, `sort`, `type`, `manager_id`, `functional_manager_id`, `is_active`, `status`) VALUES
(1, '科技创新集团', 'ORG001', '科技创新集团总部', NULL, 1, 1, 'company', 1, NULL, 1, 'active'),
(2, '研发中心', 'ORG002', '技术研发中心', 1, 2, 1, 'business_unit', 2, NULL, 1, 'active'),
(3, '运营中心', 'ORG003', '业务运营中心', 1, 2, 2, 'business_unit', 6, NULL, 1, 'active'),
(4, '支持中心', 'ORG004', '职能支持中心', 1, 2, 3, 'business_unit', 7, NULL, 1, 'active'),
(5, '产品研发部', 'ORG005', '产品研发部门', 2, 3, 1, 'department', 5, NULL, 1, 'active'),
(6, '技术架构部', 'ORG006', '技术架构部门', 2, 3, 2, 'department', 4, NULL, 1, 'active'),
(7, '质量管理部', 'ORG007', '质量管理部门', 2, 3, 3, 'department', 15, NULL, 1, 'active'),
(8, '市场营销部', 'ORG008', '市场营销部门', 3, 3, 1, 'department', 6, NULL, 1, 'active'),
(9, '销售管理部', 'ORG009', '销售管理部门', 3, 3, 2, 'department', 19, NULL, 1, 'active'),
(10, '客户成功部', 'ORG010', '客户成功部门', 3, 3, 3, 'department', NULL, NULL, 1, 'active'),
(11, '人力资源部', 'ORG011', '人力资源部门', 4, 3, 1, 'department', 7, NULL, 1, 'active'),
(12, '财务管理部', 'ORG012', '财务管理部门', 4, 3, 2, 'department', 3, NULL, 1, 'active'),
(13, '法务合规部', 'ORG013', '法务合规部门', 4, 3, 3, 'department', NULL, NULL, 1, 'active'),
(14, '前端技术团队', 'ORG014', '前端开发团队', 5, 4, 1, 'team', 11, NULL, 1, 'active'),
(15, '后端技术团队', 'ORG015', '后端开发团队', 5, 4, 2, 'team', 12, NULL, 1, 'active'),
(16, '移动端团队', 'ORG016', '移动应用开发团队', 5, 4, 3, 'team', 13, NULL, 1, 'active'),
(17, '数据技术团队', 'ORG017', '大数据和AI团队', 6, 4, 1, 'team', NULL, NULL, 1, 'active'),
(18, '基础架构团队', 'ORG018', '基础设施团队', 6, 4, 2, 'team', 18, NULL, 1, 'active'),
(19, '测试技术团队', 'ORG019', '测试工程团队', 7, 4, 1, 'team', 16, NULL, 1, 'active'),
(20, '上海研发中心', 'ORG020', '上海地区研发中心', 2, 3, 4, 'location', NULL, NULL, 1, 'active');

-- 员工分配数据 (Employee Assignments)
INSERT INTO `employee_assignments` (`employee_id`, `organization_unit_id`, `assignment_type`, `management_type`, `work_percentage`, `is_primary`, `position_id`, `job_level_id`, `direct_manager_id`, `effective_date`, `status`) VALUES
(1, 1, 'primary', 'line', 100.00, 1, 1, 11, NULL, '2020-01-01', 'active'),
(2, 2, 'primary', 'line', 100.00, 1, 2, 12, 1, '2020-02-01', 'active'),
(3, 12, 'primary', 'line', 100.00, 1, 3, 13, 1, '2020-03-01', 'active'),
(4, 6, 'primary', 'line', 100.00, 1, 4, 8, 2, '2020-04-01', 'active'),
(5, 5, 'primary', 'line', 100.00, 1, 5, 8, 2, '2020-05-01', 'active'),
(6, 8, 'primary', 'line', 100.00, 1, 6, 8, 1, '2020-06-01', 'active'),
(7, 11, 'primary', 'line', 100.00, 1, 7, 8, 1, '2020-07-01', 'active'),
(8, 5, 'primary', 'line', 100.00, 1, 8, 6, 5, '2020-08-01', 'active'),
(9, 5, 'primary', 'line', 100.00, 1, 9, 6, 5, '2020-09-01', 'active'),
(10, 5, 'primary', 'line', 100.00, 1, 10, 4, 8, '2020-10-01', 'active'),
(11, 14, 'primary', 'line', 100.00, 1, 11, 3, 8, '2021-01-01', 'active'),
(12, 15, 'primary', 'line', 100.00, 1, 11, 3, 8, '2021-02-01', 'active'),
(13, 16, 'primary', 'line', 100.00, 1, 11, 3, 8, '2021-03-01', 'active'),
(14, 5, 'primary', 'line', 100.00, 1, 12, 2, 10, '2021-04-01', 'active'),
(15, 7, 'primary', 'line', 100.00, 1, 13, 6, 4, '2021-05-01', 'active'),
(16, 19, 'primary', 'line', 100.00, 1, 14, 4, 15, '2021-06-01', 'active'),
(17, 19, 'primary', 'line', 100.00, 1, 15, 3, 16, '2021-07-01', 'active'),
(18, 18, 'primary', 'line', 100.00, 1, 16, 6, 4, '2021-08-01', 'active'),
(19, 9, 'primary', 'line', 100.00, 1, 18, 6, 6, '2021-09-01', 'active'),
(20, 11, 'primary', 'line', 100.00, 1, 19, 3, 7, '2021-10-01', 'active');

-- 招聘数据
INSERT INTO `recruitments` (`title`, `department_id`, `position_id`, `job_level_id`, `recruiter_id`, `urgency`, `headcount`, `salary_min`, `salary_max`, `requirements`, `description`, `start_date`, `status`) VALUES
('高级Java开发工程师', 5, 10, 4, 20, 'high', 2, 15000.00, 25000.00, '5年以上Java开发经验，熟悉Spring框架', '负责核心业务系统开发', '2024-07-01', 'active'),
('前端架构师', 5, 10, 7, 20, 'medium', 1, 30000.00, 45000.00, '8年以上前端经验，精通React/Vue', '负责前端技术架构设计', '2024-07-05', 'active'),
('产品经理', 5, 9, 6, 20, 'medium', 1, 20000.00, 35000.00, '3年以上产品经验，有B端产品经验', '负责产品规划和设计', '2024-07-10', 'active'),
('数据分析师', 6, 10, 4, 20, 'low', 1, 12000.00, 20000.00, '熟悉SQL和Python，有数据分析经验', '负责业务数据分析', '2024-07-12', 'active'),
('运维工程师', 6, 16, 5, 20, 'medium', 1, 15000.00, 25000.00, '熟悉Linux和Docker，有云平台经验', '负责系统运维和部署', '2024-07-15', 'active'),
('测试工程师', 7, 15, 3, 20, 'medium', 2, 8000.00, 15000.00, '2年以上测试经验，熟悉自动化测试', '负责功能和性能测试', '2024-07-18', 'active'),
('UI设计师', 5, 10, 4, 20, 'low', 1, 10000.00, 18000.00, '3年以上UI设计经验，熟悉设计工具', '负责产品界面设计', '2024-07-20', 'active'),
('市场专员', 8, 18, 3, 20, 'medium', 1, 6000.00, 12000.00, '市场营销相关专业，有推广经验', '负责市场推广活动', '2024-07-22', 'active'),
('销售代表', 9, 18, 3, 20, 'high', 3, 8000.00, 15000.00, '有销售经验，沟通能力强', '负责客户开发和维护', '2024-07-25', 'active'),
('财务分析师', 12, 20, 4, 20, 'medium', 1, 12000.00, 20000.00, '财务专业，有分析经验', '负责财务数据分析', '2024-07-28', 'active'),
('法务专员', 13, 20, 3, 20, 'low', 1, 10000.00, 18000.00, '法律专业，通过司法考试优先', '负责合同审核和法律咨询', '2024-07-30', 'active'),
('人力资源专员', 11, 19, 3, 20, 'medium', 1, 7000.00, 12000.00, '人力资源专业，有招聘经验', '负责招聘和员工关系', '2024-08-01', 'active'),
('实习生-开发', 5, 12, 1, 20, 'low', 5, 3000.00, 5000.00, '计算机相关专业在校生', '参与项目开发实习', '2024-08-05', 'active'),
('客服专员', 10, 20, 2, 20, 'medium', 2, 5000.00, 8000.00, '有客服经验，服务意识强', '负责客户咨询和投诉处理', '2024-08-08', 'active'),
('业务分析师', 3, 9, 5, 20, 'medium', 1, 18000.00, 28000.00, '有业务分析经验，逻辑思维强', '负责业务流程分析优化', '2024-08-10', 'active'),
('技术文档工程师', 2, 10, 3, 20, 'low', 1, 8000.00, 15000.00, '技术写作能力强，有文档经验', '负责技术文档编写', '2024-08-12', 'active'),
('DevOps工程师', 6, 16, 6, 20, 'high', 1, 20000.00, 35000.00, '熟悉CI/CD，有DevOps经验', '负责持续集成和部署', '2024-08-15', 'active'),
('安全工程师', 6, 16, 5, 20, 'medium', 1, 18000.00, 30000.00, '信息安全专业，有安全测试经验', '负责系统安全评估', '2024-08-18', 'active'),
('算法工程师', 6, 10, 7, 20, 'high', 1, 25000.00, 40000.00, '机器学习背景，有算法经验', '负责AI算法研发', '2024-08-20', 'active'),
('项目经理', 2, 6, 6, 20, 'medium', 1, 20000.00, 35000.00, '有项目管理经验，PMP认证优先', '负责项目规划和执行', '2024-08-22', 'active');

-- 候选人数据
INSERT INTO `candidates` (`recruitment_id`, `name`, `email`, `phone`, `experience_years`, `current_company`, `current_position`, `current_salary`, `expected_salary`, `source`, `stage`, `status`) VALUES
(1, '张三', 'zhangsan@email.com', '13901234567', 6, '阿里巴巴', 'Java开发工程师', 18000.00, 22000.00, 'job_board', 'interview', 'active'),
(1, '李四', 'lisi@email.com', '13901234568', 5, '腾讯', '后端开发工程师', 16000.00, 20000.00, 'referral', 'assessment', 'active'),
(2, '王五', 'wangwu@email.com', '13901234569', 8, '字节跳动', '前端架构师', 35000.00, 40000.00, 'referral', 'offer', 'active'),
(3, '赵六', 'zhaoliu@email.com', '13901234570', 4, '美团', '产品经理', 18000.00, 25000.00, 'job_board', 'screening', 'active'),
(4, '钱七', 'qianqi@email.com', '13901234571', 3, '京东', '数据分析师', 10000.00, 15000.00, 'social_media', 'applied', 'active'),
(5, '孙八', 'sunba@email.com', '13901234572', 4, '华为', '运维工程师', 14000.00, 18000.00, 'job_board', 'interview', 'active'),
(6, '周九', 'zhoujiu@email.com', '13901234573', 3, '百度', '测试工程师', 9000.00, 12000.00, 'referral', 'assessment', 'active'),
(7, '吴十', 'wushi@email.com', '13901234574', 4, '网易', 'UI设计师', 12000.00, 15000.00, 'job_board', 'interview', 'active'),
(8, '郑一', 'zhengyi@email.com', '13901234575', 2, '小米', '市场专员', 5000.00, 8000.00, 'social_media', 'screening', 'active'),
(9, '王二', 'wanger@email.com', '13901234576', 3, '滴滴', '销售代表', 7000.00, 10000.00, 'referral', 'applied', 'active'),
(10, '李三', 'lisan@email.com', '13901234577', 5, '蚂蚁金服', '财务分析师', 15000.00, 18000.00, 'job_board', 'offer', 'active'),
(11, '张四', 'zhangsi@email.com', '13901234578', 2, '顺丰', '法务专员', 8000.00, 12000.00, 'job_board', 'interview', 'active'),
(12, '赵五', 'zhaowu@email.com', '13901234579', 1, '58同城', 'HR专员', 6000.00, 8000.00, 'social_media', 'screening', 'active'),
(13, '钱六', 'qianliu@email.com', '13901234580', 0, '北京大学', '计算机学生', 0.00, 4000.00, 'career_fair', 'applied', 'active'),
(14, '孙七', 'sunqi@email.com', '13901234581', 2, '去哪儿', '客服专员', 4000.00, 6000.00, 'job_board', 'interview', 'active'),
(15, '周八', 'zhouba@email.com', '13901234582', 6, '携程', '业务分析师', 20000.00, 25000.00, 'referral', 'assessment', 'active'),
(16, '吴九', 'wujiu@email.com', '13901234583', 3, '新浪', '技术文档工程师', 7000.00, 10000.00, 'job_board', 'screening', 'active'),
(17, '郑十', 'zhengshi@email.com', '13901234584', 5, '快手', 'DevOps工程师', 18000.00, 25000.00, 'referral', 'offer', 'active'),
(18, '王一二', 'wangyier@email.com', '13901234585', 4, '拼多多', '安全工程师', 16000.00, 22000.00, 'job_board', 'interview', 'active'),
(19, '李三四', 'lisansi@email.com', '13901234586', 7, '商汤科技', '算法工程师', 30000.00, 35000.00, 'referral', 'assessment', 'active');

-- 绩效评估数据
INSERT INTO `performances` (`employee_id`, `evaluator_id`, `evaluation_period`, `period_start`, `period_end`, `overall_score`, `performance_rating`, `goals_achievement`, `strengths`, `areas_for_improvement`, `status`) VALUES
(10, 8, 'quarterly', '2024-04-01', '2024-06-30', 4.2, 'exceeds', '超额完成季度开发任务', '技术能力强，学习能力突出', '需要提升团队协作能力', 'approved'),
(11, 8, 'quarterly', '2024-04-01', '2024-06-30', 3.8, 'meets', '按时完成前端开发任务', '代码质量好，响应速度快', '需要学习更多新技术', 'approved'),
(12, 8, 'quarterly', '2024-04-01', '2024-06-30', 4.0, 'meets', '完成后端接口开发', '逻辑思维清晰，bug率低', '需要提升沟通表达能力', 'approved'),
(13, 8, 'quarterly', '2024-04-01', '2024-06-30', 3.9, 'meets', '移动端功能开发顺利', '对新技术敏感，学习速度快', '需要更多项目经验', 'approved'),
(14, 10, 'quarterly', '2024-04-01', '2024-06-30', 3.5, 'meets', '基本完成分配的开发任务', '工作认真负责，态度积极', '技术深度需要加强', 'approved'),
(16, 15, 'quarterly', '2024-04-01', '2024-06-30', 4.1, 'exceeds', '测试效率和质量都很高', '细心负责，测试覆盖率高', '可以学习自动化测试技能', 'approved'),
(17, 16, 'quarterly', '2024-04-01', '2024-06-30', 3.7, 'meets', '完成功能测试工作', '执行力强，发现问题及时', '需要提升测试设计能力', 'approved'),
(19, 6, 'quarterly', '2024-04-01', '2024-06-30', 4.3, 'exceeds', '销售业绩超过目标20%', '客户沟通能力强，成单率高', '需要拓展更多客户渠道', 'approved'),
(20, 7, 'quarterly', '2024-04-01', '2024-06-30', 3.8, 'meets', '招聘任务完成率90%', '面试评估准确，候选人质量高', '需要拓展招聘渠道', 'approved'),
(8, 5, 'quarterly', '2024-04-01', '2024-06-30', 4.0, 'meets', '团队管理和项目进度良好', '管理能力不断提升', '需要更多技术前瞻性', 'approved'),
(9, 5, 'quarterly', '2024-04-01', '2024-06-30', 3.9, 'meets', '产品需求分析准确', '需求理解能力强', '需要提升数据分析能力', 'approved'),
(15, 4, 'quarterly', '2024-04-01', '2024-06-30', 4.2, 'exceeds', '质量管理体系建设成效显著', '质量意识强，流程规范', '需要关注新技术质量标准', 'approved'),
(18, 4, 'quarterly', '2024-04-01', '2024-06-30', 3.8, 'meets', '系统稳定性保持在99.5%以上', '运维技能扎实，响应及时', '需要学习云原生技术', 'approved'),
(5, 2, 'quarterly', '2024-04-01', '2024-06-30', 4.1, 'exceeds', '产品规划和执行都很出色', '产品思维敏锐，执行力强', '需要更多行业洞察', 'approved'),
(6, 1, 'quarterly', '2024-04-01', '2024-06-30', 4.0, 'meets', '市场推广效果良好', '市场敏感度高，策划能力强', '需要提升数字化营销能力', 'approved'),
(7, 1, 'quarterly', '2024-04-01', '2024-06-30', 3.9, 'meets', 'HR体系建设稳步推进', '人才识别能力强', '需要更多数字化HR经验', 'approved'),
(4, 2, 'quarterly', '2024-04-01', '2024-06-30', 4.0, 'meets', '技术架构优化成果明显', '技术视野开阔，决策准确', '需要更多团队激励技巧', 'approved'),
(3, 1, 'quarterly', '2024-04-01', '2024-06-30', 4.2, 'exceeds', '财务管理规范，成本控制有效', '财务专业能力强', '需要更多业务理解', 'approved'),
(2, 1, 'quarterly', '2024-04-01', '2024-06-30', 4.1, 'exceeds', '技术战略执行良好', '技术领导力突出', '需要更多商业思维', 'approved'),
(1, NULL, 'quarterly', '2024-04-01', '2024-06-30', 4.5, 'outstanding', '公司整体发展超预期', '战略眼光和执行力都很强', '需要更多国际化视野', 'approved');

-- 组织变更数据
INSERT INTO `organization_changes` (`entity_type`, `entity_id`, `change_type`, `field_name`, `old_value`, `new_value`, `description`, `reason`, `initiator_id`, `status`) VALUES
('department', 14, 'create', NULL, NULL, '前端开发团队', '新建前端开发团队', '业务发展需要独立的前端团队', 5, 'implemented'),
('department', 15, 'create', NULL, NULL, '后端开发团队', '新建后端开发团队', '业务发展需要独立的后端团队', 5, 'implemented'),
('department', 16, 'create', NULL, NULL, '移动开发团队', '新建移动开发团队', '移动业务扩展需要专门团队', 5, 'implemented'),
('employee', 11, 'assign', 'department_id', '5', '14', '员工转移到前端团队', '专业化分工需要', 8, 'implemented'),
('employee', 12, 'assign', 'department_id', '5', '15', '员工转移到后端团队', '专业化分工需要', 8, 'implemented'),
('employee', 13, 'assign', 'department_id', '5', '16', '员工转移到移动团队', '专业化分工需要', 8, 'implemented'),
('department', 20, 'create', NULL, NULL, '北京分公司', '设立北京分公司', '业务扩展到北京地区', 1, 'approved'),
('organization_unit', 20, 'create', NULL, NULL, '上海研发中心', '设立上海研发中心', '就近服务华东地区客户', 2, 'approved'),
('employee', 14, 'update', 'job_level_id', '1', '2', '实习生转正为初级专员', '实习期满，表现优秀', 8, 'implemented'),
('department', 7, 'update', 'manager_id', NULL, '15', '任命质量保证部经理', '部门管理需要', 4, 'implemented'),
('employee', 20, 'update', 'base_salary', '8000', '9000', '薪资调整', '年度调薪', 7, 'implemented'),
('employee', 17, 'update', 'base_salary', '9000', '10000', '薪资调整', '年度调薪', 15, 'implemented'),
('department', 13, 'update', 'description', '法务部门', '法务合规部门', '部门职能扩展', '合规要求提升', 1, 'implemented'),
('organization_unit', 13, 'update', 'name', '法务部', '法务合规部', '组织单元名称更新', '职能范围扩大', 1, 'implemented'),
('employee', 16, 'update', 'position_id', '15', '14', '职位晋升', '工作表现优秀，晋升为高级测试工程师', 15, 'approved'),
('department', 18, 'update', 'manager_id', NULL, '18', '任命运维团队负责人', '团队管理需要', 4, 'pending'),
('employee', 19, 'update', 'department_id', '10', '9', '部门调整', '销售业务重组', 6, 'pending'),
('organization_unit', 17, 'update', 'status', 'planning', 'active', '数据团队正式启动', '技术架构完成', 2, 'pending'),
('employee', 10, 'update', 'job_level_id', '4', '5', '职级晋升', '技术能力突出，晋升为主管级别', 8, 'pending'),
('department', 6, 'update', 'name', '平台架构部', '技术架构部', '部门名称调整', '更好反映部门职能', 2, 'pending');

-- ==============================================
-- 8. 外键约束设置
-- ==============================================

-- 添加外键约束 (在数据插入完成后添加，避免插入时的循环依赖问题)
ALTER TABLE `users` ADD CONSTRAINT `fk_users_employee` FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;
ALTER TABLE `user_roles` ADD CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;
ALTER TABLE `user_roles` ADD CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE;
ALTER TABLE `role_permissions` ADD CONSTRAINT `fk_role_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE;
ALTER TABLE `role_permissions` ADD CONSTRAINT `fk_role_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE;

ALTER TABLE `departments` ADD CONSTRAINT `fk_departments_parent` FOREIGN KEY (`parent_id`) REFERENCES `departments` (`id`) ON DELETE SET NULL;
ALTER TABLE `departments` ADD CONSTRAINT `fk_departments_manager` FOREIGN KEY (`manager_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;
ALTER TABLE `departments` ADD CONSTRAINT `fk_departments_functional_manager` FOREIGN KEY (`functional_manager_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;

ALTER TABLE `positions` ADD CONSTRAINT `fk_positions_parent` FOREIGN KEY (`parent_id`) REFERENCES `positions` (`id`) ON DELETE SET NULL;
ALTER TABLE `positions` ADD CONSTRAINT `fk_positions_department` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`) ON DELETE SET NULL;

ALTER TABLE `employees` ADD CONSTRAINT `fk_employees_department` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`) ON DELETE SET NULL;
ALTER TABLE `employees` ADD CONSTRAINT `fk_employees_position` FOREIGN KEY (`position_id`) REFERENCES `positions` (`id`) ON DELETE SET NULL;
ALTER TABLE `employees` ADD CONSTRAINT `fk_employees_job_level` FOREIGN KEY (`job_level_id`) REFERENCES `job_levels` (`id`) ON DELETE SET NULL;
ALTER TABLE `employees` ADD CONSTRAINT `fk_employees_manager` FOREIGN KEY (`manager_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;
ALTER TABLE `employees` ADD CONSTRAINT `fk_employees_functional_manager` FOREIGN KEY (`functional_manager_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;
ALTER TABLE `employees` ADD CONSTRAINT `fk_employees_secondary_department` FOREIGN KEY (`secondary_department_id`) REFERENCES `departments` (`id`) ON DELETE SET NULL;

ALTER TABLE `department_assignments` ADD CONSTRAINT `fk_dept_assignments_employee` FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`) ON DELETE CASCADE;
ALTER TABLE `department_assignments` ADD CONSTRAINT `fk_dept_assignments_department` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`) ON DELETE CASCADE;
ALTER TABLE `department_assignments` ADD CONSTRAINT `fk_dept_assignments_position` FOREIGN KEY (`position_id`) REFERENCES `positions` (`id`) ON DELETE SET NULL;
ALTER TABLE `department_assignments` ADD CONSTRAINT `fk_dept_assignments_job_level` FOREIGN KEY (`job_level_id`) REFERENCES `job_levels` (`id`) ON DELETE SET NULL;
ALTER TABLE `department_assignments` ADD CONSTRAINT `fk_dept_assignments_direct_manager` FOREIGN KEY (`direct_manager_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;
ALTER TABLE `department_assignments` ADD CONSTRAINT `fk_dept_assignments_functional_manager` FOREIGN KEY (`functional_manager_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;

ALTER TABLE `organization_units` ADD CONSTRAINT `fk_org_units_parent` FOREIGN KEY (`parent_id`) REFERENCES `organization_units` (`id`) ON DELETE SET NULL;
ALTER TABLE `organization_units` ADD CONSTRAINT `fk_org_units_manager` FOREIGN KEY (`manager_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;
ALTER TABLE `organization_units` ADD CONSTRAINT `fk_org_units_functional_manager` FOREIGN KEY (`functional_manager_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;

ALTER TABLE `employee_assignments` ADD CONSTRAINT `fk_emp_assignments_employee` FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`) ON DELETE CASCADE;
ALTER TABLE `employee_assignments` ADD CONSTRAINT `fk_emp_assignments_org_unit` FOREIGN KEY (`organization_unit_id`) REFERENCES `organization_units` (`id`) ON DELETE CASCADE;
ALTER TABLE `employee_assignments` ADD CONSTRAINT `fk_emp_assignments_position` FOREIGN KEY (`position_id`) REFERENCES `positions` (`id`) ON DELETE SET NULL;
ALTER TABLE `employee_assignments` ADD CONSTRAINT `fk_emp_assignments_job_level` FOREIGN KEY (`job_level_id`) REFERENCES `job_levels` (`id`) ON DELETE SET NULL;
ALTER TABLE `employee_assignments` ADD CONSTRAINT `fk_emp_assignments_direct_manager` FOREIGN KEY (`direct_manager_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;
ALTER TABLE `employee_assignments` ADD CONSTRAINT `fk_emp_assignments_functional_manager` FOREIGN KEY (`functional_manager_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;

ALTER TABLE `organization_changes` ADD CONSTRAINT `fk_org_changes_initiator` FOREIGN KEY (`initiator_id`) REFERENCES `employees` (`id`) ON DELETE RESTRICT;
ALTER TABLE `organization_changes` ADD CONSTRAINT `fk_org_changes_approver` FOREIGN KEY (`approver_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;

ALTER TABLE `organization_histories` ADD CONSTRAINT `fk_org_histories_unit` FOREIGN KEY (`unit_id`) REFERENCES `organization_units` (`id`) ON DELETE CASCADE;
ALTER TABLE `organization_histories` ADD CONSTRAINT `fk_org_histories_changed_by` FOREIGN KEY (`changed_by`) REFERENCES `employees` (`id`) ON DELETE RESTRICT;

ALTER TABLE `organization_change_logs` ADD CONSTRAINT `fk_org_change_logs_operator` FOREIGN KEY (`operator_id`) REFERENCES `employees` (`id`) ON DELETE RESTRICT;

ALTER TABLE `work_experiences` ADD CONSTRAINT `fk_work_experiences_employee` FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`) ON DELETE CASCADE;
ALTER TABLE `attendances` ADD CONSTRAINT `fk_attendances_employee` FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`) ON DELETE CASCADE;
ALTER TABLE `leaves` ADD CONSTRAINT `fk_leaves_employee` FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`) ON DELETE CASCADE;
ALTER TABLE `leaves` ADD CONSTRAINT `fk_leaves_approver` FOREIGN KEY (`approver_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;
ALTER TABLE `salaries` ADD CONSTRAINT `fk_salaries_employee` FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`) ON DELETE CASCADE;
ALTER TABLE `payroll_records` ADD CONSTRAINT `fk_payroll_records_salary` FOREIGN KEY (`salary_id`) REFERENCES `salaries` (`id`) ON DELETE CASCADE;
ALTER TABLE `payroll_records` ADD CONSTRAINT `fk_payroll_records_processor` FOREIGN KEY (`processor_id`) REFERENCES `employees` (`id`) ON DELETE RESTRICT;

ALTER TABLE `recruitments` ADD CONSTRAINT `fk_recruitments_department` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`) ON DELETE RESTRICT;
ALTER TABLE `recruitments` ADD CONSTRAINT `fk_recruitments_position` FOREIGN KEY (`position_id`) REFERENCES `positions` (`id`) ON DELETE RESTRICT;
ALTER TABLE `recruitments` ADD CONSTRAINT `fk_recruitments_job_level` FOREIGN KEY (`job_level_id`) REFERENCES `job_levels` (`id`) ON DELETE SET NULL;
ALTER TABLE `recruitments` ADD CONSTRAINT `fk_recruitments_recruiter` FOREIGN KEY (`recruiter_id`) REFERENCES `employees` (`id`) ON DELETE RESTRICT;

ALTER TABLE `candidates` ADD CONSTRAINT `fk_candidates_recruitment` FOREIGN KEY (`recruitment_id`) REFERENCES `recruitments` (`id`) ON DELETE CASCADE;
ALTER TABLE `candidates` ADD CONSTRAINT `fk_candidates_referrer` FOREIGN KEY (`referrer_id`) REFERENCES `employees` (`id`) ON DELETE SET NULL;

ALTER TABLE `performances` ADD CONSTRAINT `fk_performances_employee` FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`) ON DELETE CASCADE;
ALTER TABLE `performances` ADD CONSTRAINT `fk_performances_evaluator` FOREIGN KEY (`evaluator_id`) REFERENCES `employees` (`id`) ON DELETE RESTRICT;

-- ==============================================
-- 9. 索引优化
-- ==============================================

-- 创建复合索引以提升查询性能
CREATE INDEX `idx_employees_dept_status` ON `employees` (`department_id`, `status`);
CREATE INDEX `idx_employees_manager_status` ON `employees` (`manager_id`, `status`);
CREATE INDEX `idx_dept_assignments_emp_status` ON `department_assignments` (`employee_id`, `status`);
CREATE INDEX `idx_emp_assignments_org_status` ON `employee_assignments` (`organization_unit_id`, `status`);
CREATE INDEX `idx_attendances_emp_date` ON `attendances` (`employee_id`, `date`);
CREATE INDEX `idx_salaries_emp_month` ON `salaries` (`employee_id`, `month`);
CREATE INDEX `idx_leaves_emp_status` ON `leaves` (`employee_id`, `status`);
CREATE INDEX `idx_org_changes_entity_status` ON `organization_changes` (`entity_type`, `entity_id`, `status`);
CREATE INDEX `idx_org_histories_unit_date` ON `organization_histories` (`unit_id`, `snapshot_date`);
CREATE INDEX `idx_recruitments_dept_status` ON `recruitments` (`department_id`, `status`);
CREATE INDEX `idx_candidates_recruitment_status` ON `candidates` (`recruitment_id`, `status`);
CREATE INDEX `idx_performances_emp_period` ON `performances` (`employee_id`, `period_start`, `period_end`);

-- ==============================================
-- 初始化脚本完成
-- ==============================================