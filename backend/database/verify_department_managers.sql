-- ============================================
-- 验证部门负责人配置
-- 确保每个层级的部门都有负责人
-- ============================================

USE gin_db;

-- 1. 查看所有部门的负责人配置（按层级排序）
SELECT 
    d.id as dept_id,
    d.name as dept_name,
    d.level,
    d.type,
    d.manager_id,
    e.name as manager_name,
    e.employee_id as manager_employee_id,
    pos.name as manager_position,
    CASE 
        WHEN d.manager_id IS NULL THEN '❌ 缺少负责人'
        ELSE '✅ 有负责人'
    END as manager_status
FROM departments d
LEFT JOIN employees e ON d.manager_id = e.id
LEFT JOIN positions pos ON e.position_id = pos.id
WHERE d.is_active = 1
ORDER BY d.level, d.sort;

-- 2. 按层级统计部门负责人配置
SELECT 
    d.level as department_level,
    COUNT(*) as total_departments,
    COUNT(d.manager_id) as departments_with_manager,
    COUNT(*) - COUNT(d.manager_id) as departments_without_manager,
    ROUND(COUNT(d.manager_id) * 100.0 / COUNT(*), 2) as manager_coverage_percentage
FROM departments d
WHERE d.is_active = 1
GROUP BY d.level
ORDER BY d.level;

-- 3. 查看哪些员工同时担任多个部门的负责人
SELECT 
    e.id as employee_id,
    e.name as employee_name,
    e.employee_id,
    pos.name as position,
    COUNT(d.id) as departments_managed,
    GROUP_CONCAT(CONCAT(d.name, '(', d.code, ')') ORDER BY d.level, d.name SEPARATOR ', ') as managed_departments
FROM employees e
INNER JOIN departments d ON e.id = d.manager_id
LEFT JOIN positions pos ON e.position_id = pos.id
WHERE d.is_active = 1 AND e.status = 'active'
GROUP BY e.id, e.name, e.employee_id, pos.name
HAVING COUNT(d.id) > 1
ORDER BY COUNT(d.id) DESC, e.name;

-- 4. 查看部门层级结构和负责人链
WITH RECURSIVE dept_hierarchy AS (
    -- 根部门
    SELECT 
        d.id,
        d.name,
        d.code,
        d.level,
        d.parent_id,
        d.manager_id,
        e.name as manager_name,
        CAST(d.name AS CHAR(1000)) as hierarchy_path,
        CAST(CONCAT(e.name, '(', e.employee_id, ')') AS CHAR(1000)) as manager_path
    FROM departments d
    LEFT JOIN employees e ON d.manager_id = e.id
    WHERE d.parent_id IS NULL AND d.is_active = 1
    
    UNION ALL
    
    -- 子部门
    SELECT 
        d.id,
        d.name,
        d.code,
        d.level,
        d.parent_id,
        d.manager_id,
        e.name as manager_name,
        CAST(CONCAT(dh.hierarchy_path, ' > ', d.name) AS CHAR(1000)) as hierarchy_path,
        CAST(CONCAT(dh.manager_path, ' > ', e.name, '(', e.employee_id, ')') AS CHAR(1000)) as manager_path
    FROM departments d
    INNER JOIN dept_hierarchy dh ON d.parent_id = dh.id
    LEFT JOIN employees e ON d.manager_id = e.id
    WHERE d.is_active = 1
)
SELECT 
    level,
    hierarchy_path as department_hierarchy,
    manager_path as manager_hierarchy
FROM dept_hierarchy
ORDER BY level, hierarchy_path;

-- 5. 验证员工和部门的关系一致性
SELECT 
    '员工部门关系验证' as check_type,
    e.id as employee_id,
    e.name as employee_name,
    e.department_id as assigned_dept_id,
    dept.name as assigned_dept_name,
    GROUP_CONCAT(DISTINCT managed_dept.name ORDER BY managed_dept.name SEPARATOR ', ') as managed_departments
FROM employees e
INNER JOIN departments dept ON e.department_id = dept.id
LEFT JOIN departments managed_dept ON e.id = managed_dept.manager_id AND managed_dept.is_active = 1
WHERE e.status = 'active'
GROUP BY e.id, e.name, e.department_id, dept.name
HAVING managed_departments IS NOT NULL
ORDER BY e.name;

-- 6. 检查是否有循环管理关系
SELECT 
    '循环管理关系检查' as check_type,
    emp1.id as employee_id,
    emp1.name as employee_name,
    emp1.manager_id,
    mgr.name as direct_manager,
    CASE 
        WHEN emp1.id IN (
            SELECT DISTINCT d.manager_id 
            FROM departments d 
            WHERE d.manager_id IN (
                SELECT e2.id 
                FROM employees e2 
                WHERE e2.manager_id = emp1.id
            )
        ) THEN '⚠️ 可能存在循环管理'
        ELSE '✅ 正常'
    END as relationship_status
FROM employees emp1
LEFT JOIN employees mgr ON emp1.manager_id = mgr.id
WHERE emp1.status = 'active'
ORDER BY emp1.name;

-- 7. 总结报告
SELECT '=== 部门负责人配置总结 ===' as summary;
SELECT 
    CONCAT('总部门数: ', COUNT(*)) as total_info
FROM departments WHERE is_active = 1
UNION ALL
SELECT 
    CONCAT('有负责人部门数: ', COUNT(*)) as manager_info
FROM departments WHERE is_active = 1 AND manager_id IS NOT NULL
UNION ALL
SELECT 
    CONCAT('负责人覆盖率: ', ROUND(
        (SELECT COUNT(*) FROM departments WHERE is_active = 1 AND manager_id IS NOT NULL) * 100.0 / 
        (SELECT COUNT(*) FROM departments WHERE is_active = 1), 2
    ), '%') as coverage_info;