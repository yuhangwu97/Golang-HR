-- 修复departments表缺少status字段的问题
USE gin_db;

-- 添加status字段到departments表
ALTER TABLE departments ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态';

-- 查看表结构确认字段已添加
DESCRIBE departments;