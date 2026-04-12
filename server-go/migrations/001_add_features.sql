-- 数据库迁移脚本
-- 为6个新功能创建所需的表

-- 1. 保存的筛选器表
CREATE TABLE IF NOT EXISTS saved_filters (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  filters TEXT,
  is_default INTEGER DEFAULT 0,
  create_time TEXT
);

-- 2. 工作流规则表
CREATE TABLE IF NOT EXISTS workflow_rules (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  trigger_type TEXT NOT NULL,
  trigger_value TEXT NOT NULL,
  action_type TEXT NOT NULL,
  action_config TEXT,
  enabled INTEGER DEFAULT 1,
  priority INTEGER DEFAULT 1,
  create_time TEXT
);

-- 插入示例工作流规则（可选）
-- 当状态变为"洽谈中"时，自动添加"待跟进"标签
-- INSERT INTO workflow_rules (name, trigger_type, trigger_value, action_type, action_config, enabled, priority, create_time)
-- VALUES ('洽谈中自动添加标签', 'status_change', '洽谈中', 'tag_add', '{"tag_ids": [1]}', 1, 1, datetime('now'));
