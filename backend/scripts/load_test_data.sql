-- 导入测试数据脚本
-- 使用方法：mysql -u root -p gin_db < scripts/load_test_data.sql

USE gin_db;

-- 清空现有数据（注意：这会删除所有现有数据）
SET FOREIGN_KEY_CHECKS = 0;

-- 清空表数据（按依赖关系顺序）
TRUNCATE TABLE candidates;
TRUNCATE TABLE recruitments;
TRUNCATE TABLE payroll_records;
TRUNCATE TABLE salaries;
TRUNCATE TABLE leaves;
TRUNCATE TABLE attendances;
TRUNCATE TABLE role_permissions;
TRUNCATE TABLE user_roles;
TRUNCATE TABLE work_experiences;
TRUNCATE TABLE users;
TRUNCATE TABLE employees;
TRUNCATE TABLE positions;
TRUNCATE TABLE departments;
TRUNCATE TABLE job_levels;
TRUNCATE TABLE permissions;
TRUNCATE TABLE roles;

-- 重置自增ID
ALTER TABLE candidates AUTO_INCREMENT = 1;
ALTER TABLE recruitments AUTO_INCREMENT = 1;
ALTER TABLE payroll_records AUTO_INCREMENT = 1;
ALTER TABLE salaries AUTO_INCREMENT = 1;
ALTER TABLE leaves AUTO_INCREMENT = 1;
ALTER TABLE attendances AUTO_INCREMENT = 1;
ALTER TABLE role_permissions AUTO_INCREMENT = 1;
ALTER TABLE user_roles AUTO_INCREMENT = 1;
ALTER TABLE work_experiences AUTO_INCREMENT = 1;
ALTER TABLE users AUTO_INCREMENT = 1;
ALTER TABLE employees AUTO_INCREMENT = 1;
ALTER TABLE positions AUTO_INCREMENT = 1;
ALTER TABLE departments AUTO_INCREMENT = 1;
ALTER TABLE job_levels AUTO_INCREMENT = 1;
ALTER TABLE permissions AUTO_INCREMENT = 1;
ALTER TABLE roles AUTO_INCREMENT = 1;

-- 导入基础数据

-- 1. 角色数据
INSERT INTO `roles` (`name`, `code`, `description`, `status`, `created_at`, `updated_at`) VALUES
('系统管理员', 'admin', '系统管理员，拥有所有权限', 'active', NOW(), NOW()),
('人事管理员', 'hr', '人事管理员，负责人事相关业务', 'active', NOW(), NOW()),
('普通员工', 'employee', '普通员工用户', 'active', NOW(), NOW());

-- 2. 权限数据
INSERT INTO `permissions` (`name`, `code`, `resource`, `action`, `description`, `status`, `created_at`, `updated_at`) VALUES
('用户查看', 'user:view', 'user', 'view', '查看用户信息', 'active', NOW(), NOW()),
('用户创建', 'user:create', 'user', 'create', '创建用户', 'active', NOW(), NOW()),
('员工查看', 'employee:view', 'employee', 'view', '查看员工信息', 'active', NOW(), NOW()),
('员工创建', 'employee:create', 'employee', 'create', '创建员工', 'active', NOW(), NOW()),
('部门管理', 'department:manage', 'department', 'manage', '部门管理权限', 'active', NOW(), NOW());

-- 3. 部门数据
INSERT INTO `departments` (`name`, `code`, `parent_id`, `description`, `status`, `sort`, `created_at`, `updated_at`) VALUES
('总经办', 'CEO', NULL, '公司最高管理层', 'active', 100, NOW(), NOW()),
('技术部', 'TECH', NULL, '负责技术研发', 'active', 90, NOW(), NOW()),
('人事部', 'HR', NULL, '负责人力资源管理', 'active', 80, NOW(), NOW()),
('财务部', 'FINANCE', NULL, '负责财务管理', 'active', 70, NOW(), NOW()),
('销售部', 'SALES', NULL, '负责销售业务', 'active', 60, NOW(), NOW());

-- 4. 职级数据
INSERT INTO `job_levels` (`name`, `code`, `level`, `min_salary`, `max_salary`, `description`, `status`, `created_at`, `updated_at`) VALUES
('初级', 'JUNIOR', 1, 5000.00, 8000.00, '初级员工', 'active', NOW(), NOW()),
('中级', 'MIDDLE', 2, 8000.00, 15000.00, '中级员工', 'active', NOW(), NOW()),
('高级', 'SENIOR', 3, 15000.00, 25000.00, '高级员工', 'active', NOW(), NOW()),
('专家', 'EXPERT', 4, 25000.00, 40000.00, '专家级员工', 'active', NOW(), NOW()),
('总监', 'DIRECTOR', 5, 40000.00, 80000.00, '部门总监', 'active', NOW(), NOW());

-- 5. 职位数据
INSERT INTO `positions` (`name`, `code`, `department_id`, `description`, `requirements`, `status`, `created_at`, `updated_at`) VALUES
('总经理', 'CEO', 1, '公司总经理', '具有丰富的管理经验', 'active', NOW(), NOW()),
('技术总监', 'CTO', 2, '技术部门负责人', '具有丰富的技术管理经验', 'active', NOW(), NOW()),
('人事经理', 'HR_MANAGER', 3, '人事部门负责人', '具有人力资源管理经验', 'active', NOW(), NOW()),
('财务经理', 'FINANCE_MANAGER', 4, '财务部门负责人', '具有财务管理经验', 'active', NOW(), NOW()),
('销售经理', 'SALES_MANAGER', 5, '销售部门负责人', '具有销售管理经验', 'active', NOW(), NOW());

-- 6. 员工数据 
INSERT INTO `employees` (`employee_id`, `name`, `email`, `phone`, `gender`, `status`, `department_id`, `position_id`, `job_level_id`, `hire_date`, `base_salary`, `created_at`, `updated_at`) VALUES
('EMP001', '张三', 'zhangsan@company.com', '13800138001', 'male', 'active', 1, 1, 5, '2020-01-01', 50000.00, NOW(), NOW()),
('EMP002', '李四', 'lisi@company.com', '13800138002', 'male', 'active', 2, 2, 4, '2020-02-01', 30000.00, NOW(), NOW()),
('EMP003', '王五', 'wangwu@company.com', '13800138003', 'female', 'active', 3, 3, 3, '2020-03-01', 20000.00, NOW(), NOW()),
('EMP004', '赵六', 'zhaoliu@company.com', '13800138004', 'male', 'active', 4, 4, 3, '2020-04-01', 22000.00, NOW(), NOW()),
('EMP005', '孙七', 'sunqi@company.com', '13800138005', 'female', 'active', 5, 5, 2, '2020-05-01', 15000.00, NOW(), NOW());

-- 7. 用户数据（密码都是admin123的bcrypt加密）
INSERT INTO `users` (`username`, `email`, `password`, `employee_id`, `status`, `created_at`, `updated_at`) VALUES
('zhangsan', 'zhangsan@company.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2uheWG/igi.', 1, 'active', NOW(), NOW()),
('lisi', 'lisi@company.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2uheWG/igi.', 2, 'active', NOW(), NOW()),
('wangwu', 'wangwu@company.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2uheWG/igi.', 3, 'active', NOW(), NOW()),
('zhaoliu', 'zhaoliu@company.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2uheWG/igi.', 4, 'active', NOW(), NOW()),
('sunqi', 'sunqi@company.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2uheWG/igi.', 5, 'active', NOW(), NOW());

-- 8. 用户角色关联
INSERT INTO `user_roles` (`user_id`, `role_id`, `created_at`) VALUES
(1, 1, NOW()),  -- 张三是管理员
(2, 1, NOW()),  -- 李四是管理员  
(3, 2, NOW()),  -- 王五是人事
(4, 3, NOW()),  -- 赵六是员工
(5, 3, NOW());  -- 孙七是员工

-- 9. 角色权限关联
INSERT INTO `role_permissions` (`role_id`, `permission_id`, `created_at`) VALUES
(1, 1, NOW()), -- 管理员拥有所有权限
(1, 2, NOW()),
(1, 3, NOW()),
(1, 4, NOW()),
(1, 5, NOW()),
(2, 3, NOW()), -- 人事拥有员工查看权限
(2, 4, NOW()), -- 人事拥有员工创建权限
(3, 3, NOW()); -- 员工拥有员工查看权限

-- 更新部门管理者
UPDATE `departments` SET `manager_id` = 1 WHERE `id` = 1;
UPDATE `departments` SET `manager_id` = 2 WHERE `id` = 2;
UPDATE `departments` SET `manager_id` = 3 WHERE `id` = 3;
UPDATE `departments` SET `manager_id` = 4 WHERE `id` = 4;
UPDATE `departments` SET `manager_id` = 5 WHERE `id` = 5;

SET FOREIGN_KEY_CHECKS = 1;

SELECT '测试数据导入完成！' as message;