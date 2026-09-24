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
  confirmer VARCHAR(80) NOT NULL DEFAULT '' COMMENT '确认人，确认前为空',
  provider VARCHAR(80) NOT NULL COMMENT '提供方',
  learner VARCHAR(80) NOT NULL COMMENT '学习方',
  pair_name VARCHAR(160) NOT NULL,
  offer_skill VARCHAR(120) NOT NULL,
  wanted_skill VARCHAR(120) NOT NULL,
  slot VARCHAR(80) NOT NULL COMMENT '共同时段',
  meeting_mode VARCHAR(20) NOT NULL COMMENT 'online 线上会议 / offline 线下面对面',
  location VARCHAR(255) NOT NULL COMMENT '线上会议链接或线下地点',
  agenda TEXT,
  status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT 'pending 待确认 / confirmed 待完成 / completed 已完成 / cancelled 已取消',
  cancelled_by VARCHAR(80) NOT NULL DEFAULT '',
  created_at VARCHAR(32) NOT NULL,
  confirmed_at VARCHAR(32) NOT NULL DEFAULT '',
  completed_at VARCHAR(32) NOT NULL DEFAULT '',
  cancelled_at VARCHAR(32) NOT NULL DEFAULT '',
  INDEX idx_appointment_status (status),
  INDEX idx_appointment_slot (slot)
) COMMENT='交换预约：pending/confirmed 占用时段，取消后释放';

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
