CREATE TABLE IF NOT EXISTS users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(80) NOT NULL,
  major VARCHAR(120) NOT NULL,
  credit_score INT NOT NULL DEFAULT 80,
  credit_level VARCHAR(40) NOT NULL
);

CREATE TABLE IF NOT EXISTS skills (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  title VARCHAR(120) NOT NULL,
  category VARCHAR(40) NOT NULL,
  level_score INT NOT NULL,
  campus VARCHAR(40) NOT NULL,
  description TEXT NOT NULL,
  portfolio VARCHAR(160) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS needs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  title VARCHAR(120) NOT NULL,
  category VARCHAR(40) NOT NULL,
  campus VARCHAR(40) NOT NULL,
  expect_time VARCHAR(80) NOT NULL,
  budget_type VARCHAR(40) NOT NULL,
  description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS appointments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  match_id BIGINT NOT NULL,
  initiator VARCHAR(80) NOT NULL COMMENT '发起人',
  confirmer VARCHAR(80) NOT NULL COMMENT '待确认/确认人',
  provider VARCHAR(80) NOT NULL COMMENT '技能提供方',
  learner VARCHAR(80) NOT NULL COMMENT '技能学习方',
  offer_skill VARCHAR(120) NOT NULL,
  wanted_skill VARCHAR(120) NOT NULL,
  pair_name VARCHAR(120) NOT NULL COMMENT '双方展示名，如 林澈 ↔ 孟野',
  exchange_slot VARCHAR(80) NOT NULL COMMENT '共同可用时段',
  location_type VARCHAR(16) NOT NULL COMMENT 'ONLINE 线上会议 / OFFLINE 线下地点',
  place VARCHAR(200) NOT NULL COMMENT '会议链接或线下地点',
  agenda VARCHAR(200) NOT NULL DEFAULT '' COMMENT '协商议程',
  status VARCHAR(16) NOT NULL COMMENT 'PENDING 待确认 / CONFIRMED 待完成 / COMPLETED 已完成 / CANCELLED 已取消',
  created_at VARCHAR(16) NOT NULL,
  updated_at VARCHAR(16) NOT NULL DEFAULT '',
  INDEX idx_appointment_slot (exchange_slot, status),
  INDEX idx_appointment_participant (initiator, confirmer)
);

CREATE TABLE IF NOT EXISTS reviews (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  from_user VARCHAR(80) NOT NULL,
  to_user VARCHAR(80) NOT NULL,
  rating INT NOT NULL,
  content TEXT NOT NULL
);

INSERT INTO users(name, major, credit_score, credit_level) VALUES
('林澈', '新闻传播 2023', 91, '黄金导师'),
('孟野', '音乐表演 2022', 88, '白银协作者'),
('周芮', '统计学 2021', 93, '黄金导师');

INSERT INTO skills(user_id, title, category, level_score, campus, description, portfolio) VALUES
(1, '毕业照人像摄影', '摄影', 92, '东校区', '提供构图、修图和毕业季跟拍，可交换吉他入门课。', '12组校园人像作品'),
(2, '民谣吉他陪练', '乐器', 81, '西校区', '节奏型、弹唱和舞台经验分享，想找人拍宣传照。', '校园音乐节演出视频'),
(3, 'Python 数据分析', '编程', 88, '中心校区', 'pandas、可视化、论文数据清洗辅导。', '3份课程项目证书');

-- 预约状态机：PENDING 待对方确认 -> CONFIRMED 待完成 -> COMPLETED 已完成；完成前任一方可 CANCELLED。
-- PENDING / CONFIRMED 占用时段，同一参与者在同一 exchange_slot 存在这两种状态时不允许新建预约；
-- COMPLETED / CANCELLED 释放时段，可再次发起预约。
INSERT INTO appointments(
  match_id, initiator, confirmer, provider, learner, offer_skill, wanted_skill,
  pair_name, exchange_slot, location_type, place, agenda, status, created_at, updated_at
) VALUES
(1, '孟野', '林澈', '林澈', '孟野', '毕业照人像摄影', '民谣吉他陪练',
 '林澈 ↔ 孟野', '周六上午', 'OFFLINE', '东校区湖边', '先拍宣传照，再约 2 次吉他课',
 'CONFIRMED', '2026-09-20 09:30', '2026-09-21 12:10'),
(2, '许安', '周芮', '周芮', '许安', 'Python 数据分析', '论文数据清洗',
 '周芮 ↔ 许安', '周二晚', 'ONLINE', 'https://meeting.example.com/room/python-2048', '导入问卷 CSV 并完成基础可视化',
 'PENDING', '2026-09-22 20:15', '2026-09-22 20:15');
