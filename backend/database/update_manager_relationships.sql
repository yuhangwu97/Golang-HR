-- ============================================
-- 更新管理关系脚本
-- 确保每个部门都有负责人，每个员工都有直属领导
-- ============================================

-- 1. 首先更新员工的manager_id，建立直属领导关系
UPDATE `employees` SET `manager_id` = 1 WHERE `id` IN (2, 3); -- CTO、CFO 直属 CEO
UPDATE `employees` SET `manager_id` = 2 WHERE `id` IN (4, 5); -- 技术相关人员直属 CTO
UPDATE `employees` SET `manager_id` = 1 WHERE `id` IN (6, 7); -- 市场总监、HR总监 直属 CEO
UPDATE `employees` SET `manager_id` = 5 WHERE `id` IN (8, 9); -- 开发经理、产品经理 直属产品总监
UPDATE `employees` SET `manager_id` = 8 WHERE `id` IN (10, 11, 12, 13, 14); -- 开发人员直属开发经理
UPDATE `employees` SET `manager_id` = 4 WHERE `id` IN (15, 18); -- 测试经理、运维经理 直属技术总监
UPDATE `employees` SET `manager_id` = 15 WHERE `id` IN (16, 17); -- 测试人员直属测试经理
UPDATE `employees` SET `manager_id` = 6 WHERE `id` = 19; -- 销售经理直属市场总监
UPDATE `employees` SET `manager_id` = 3 WHERE `id` = 20; -- 财务专员直属 CFO

-- 2. 更新部门的manager_id，指定每个部门的负责人
UPDATE `departments` SET `manager_id` = 1 WHERE `id` = 1; -- 公司CEO
UPDATE `departments` SET `manager_id` = 2 WHERE `id` = 2; -- 技术研发事业部-CTO
UPDATE `departments` SET `manager_id` = 6 WHERE `id` = 3; -- 市场运营事业部-市场总监
UPDATE `departments` SET `manager_id` = 7 WHERE `id` = 4; -- 职能支持中心-HR总监
UPDATE `departments` SET `manager_id` = 5 WHERE `id` = 5; -- 产品开发部-产品总监
UPDATE `departments` SET `manager_id` = 4 WHERE `id` = 6; -- 平台架构部-技术总监
UPDATE `departments` SET `manager_id` = 15 WHERE `id` = 7; -- 质量保证部-测试经理
UPDATE `departments` SET `manager_id` = 6 WHERE `id` = 8; -- 市场推广部-市场总监
UPDATE `departments` SET `manager_id` = 19 WHERE `id` = 9; -- 销售部-销售经理
UPDATE `departments` SET `manager_id` = 6 WHERE `id` = 10; -- 客户服务部-市场总监(临时)
UPDATE `departments` SET `manager_id` = 7 WHERE `id` = 11; -- 人力资源部-HR总监
UPDATE `departments` SET `manager_id` = 3 WHERE `id` = 12; -- 财务部-CFO
UPDATE `departments` SET `manager_id` = 1 WHERE `id` = 13; -- 法务部-CEO(临时)
UPDATE `departments` SET `manager_id` = 8 WHERE `id` = 14; -- 前端开发团队-开发经理
UPDATE `departments` SET `manager_id` = 8 WHERE `id` = 15; -- 后端开发团队-开发经理
UPDATE `departments` SET `manager_id` = 8 WHERE `id` = 16; -- 移动开发团队-开发经理
UPDATE `departments` SET `manager_id` = 4 WHERE `id` = 17; -- 数据团队-技术总监
UPDATE `departments` SET `manager_id` = 18 WHERE `id` = 18; -- 运维团队-运维经理
UPDATE `departments` SET `manager_id` = 15 WHERE `id` = 19; -- 测试团队-测试经理
UPDATE `departments` SET `manager_id` = 1 WHERE `id` = 20; -- 北京分公司-CEO

-- 3. 验证查询：检查部门负责人关系
SELECT 
    d.id as dept_id,
    d.name as dept_name,
    d.manager_id,
    e.name as manager_name,
    e.employee_id as manager_employee_id
FROM departments d
LEFT JOIN employees e ON d.manager_id = e.id
WHERE d.is_active = 1
ORDER BY d.level, d.sort;

-- 4. 验证查询：检查员工直属领导关系
SELECT 
    emp.id,
    emp.employee_id,
    emp.name as employee_name,
    emp.manager_id,
    mgr.name as manager_name,
    mgr.employee_id as manager_employee_id,
    dept.name as department_name
FROM employees emp
LEFT JOIN employees mgr ON emp.manager_id = mgr.id
LEFT JOIN departments dept ON emp.department_id = dept.id
WHERE emp.status = 'active'
ORDER BY emp.department_id, emp.id;

-- 5. 检查没有直属领导的员工（除了CEO）
SELECT 
    emp.id,
    emp.employee_id,
    emp.name as employee_name,
    pos.name as position_name
FROM employees emp
LEFT JOIN positions pos ON emp.position_id = pos.id
WHERE emp.manager_id IS NULL 
    AND emp.id != 1 -- CEO除외
    AND emp.status = 'active';

-- 6. 检查没有负责人的部门
SELECT 
    d.id as dept_id,
    d.name as dept_name,
    d.type as dept_type
FROM departments d
WHERE d.manager_id IS NULL 
    AND d.is_active = 1;