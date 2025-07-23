-- 创建薪资周期表
CREATE TABLE IF NOT EXISTS `payroll_periods` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '周期名称',
    `period_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '周期类型',
    `year` int NOT NULL COMMENT '年份',
    `month` int DEFAULT NULL COMMENT '月份',
    `quarter` int DEFAULT NULL COMMENT '季度',
    `start_date` date NOT NULL COMMENT '开始日期',
    `end_date` date NOT NULL COMMENT '结束日期',
    `pay_date` date DEFAULT NULL COMMENT '发薪日期',
    `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'draft' COMMENT '状态',
    `is_locked` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否锁定',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_period_type` (`period_type`),
    KEY `idx_year` (`year`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入示例数据
INSERT INTO `payroll_periods` (`name`, `period_type`, `year`, `month`, `start_date`, `end_date`, `status`) VALUES
('2025年1月', 'monthly', 2025, 1, '2025-01-01', '2025-01-31', 'draft'),
('2025年2月', 'monthly', 2025, 2, '2025-02-01', '2025-02-28', 'draft'),
('2025年3月', 'monthly', 2025, 3, '2025-03-01', '2025-03-31', 'draft'),
('2025年第一季度', 'quarterly', 2025, NULL, '2025-01-01', '2025-03-31', 'draft'),
('2025年度', 'yearly', 2025, NULL, '2025-01-01', '2025-12-31', 'draft');