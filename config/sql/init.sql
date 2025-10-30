CREATE DATABASE IF NOT EXISTS rpc CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE rpc;

CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_name VARCHAR(16) NOT NULL COMMENT '用户名',
  password VARCHAR(255) NOT NULL COMMENT '密码',
  avatar VARCHAR(255) DEFAULT NULL COMMENT '头像URL',
  totp VARCHAR(255) DEFAULT NULL COMMENT 'TOTP密钥',
  created_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  updated_at DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  deleted_at DATETIME(6) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_users_user_name (user_name),
  KEY idx_deleted_at (deleted_at)
) DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

CREATE TABLE IF NOT EXISTS videos (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '发布者ID',
    video_url VARCHAR(255) NOT NULL COMMENT '视频URL',
    cover_url VARCHAR(255) DEFAULT NULL COMMENT '封面URL',
    title VARCHAR(255) NOT NULL COMMENT '标题',
    description TEXT COMMENT '描述',
    visit_count BIGINT UNSIGNED DEFAULT 0 COMMENT '访问量',
    like_count BIGINT UNSIGNED DEFAULT 0 COMMENT '点赞数',
    comment_count BIGINT UNSIGNED DEFAULT 0 COMMENT '评论数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME NULL COMMENT '删除时间',
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at),
    INDEX idx_visit_count (visit_count),
    INDEX idx_like_count (like_count),
    INDEX idx_deleted_at (deleted_at)
) DEFAULT CHARSET=utf8mb4 COMMENT='视频表';

CREATE TABLE IF NOT EXISTS follows (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '被关注者ID',
    follower_id BIGINT UNSIGNED NOT NULL COMMENT '关注者ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME NULL COMMENT '删除时间',
    UNIQUE KEY uk_user_follower (user_id, follower_id),
    INDEX idx_follower_user_created (follower_id, user_id, created_at),
    INDEX idx_user_follower_created (user_id, follower_id, created_at),
    INDEX idx_deleted_at (deleted_at)
) DEFAULT CHARSET=utf8mb4 COMMENT='关注表';

CREATE TABLE IF NOT EXISTS likes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '点赞用户ID',
    target_id BIGINT UNSIGNED NOT NULL COMMENT '目标ID（视频或评论）',
    type TINYINT UNSIGNED NOT NULL COMMENT '类型：1-视频，2-评论',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME NULL COMMENT '删除时间',
    UNIQUE KEY uk_user_target_type (user_id, target_id, type),
    INDEX idx_target_type (target_id, type),
    INDEX idx_user_type_created (user_id, type, created_at),
    INDEX idx_deleted_at (deleted_at)
) DEFAULT CHARSET=utf8mb4 COMMENT='点赞表';

CREATE TABLE IF NOT EXISTS comments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '评论用户ID',
    video_id BIGINT UNSIGNED NOT NULL COMMENT '视频ID',
    parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父评论ID',
    like_count BIGINT UNSIGNED DEFAULT 0 COMMENT '点赞数',
    child_count BIGINT UNSIGNED DEFAULT 0 COMMENT '子评论数',
    content TEXT NOT NULL COMMENT '评论内容',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME NULL COMMENT '删除时间',
    INDEX idx_video_created (video_id, parent_id, created_at),
    INDEX idx_parent_created (parent_id, created_at),
    INDEX idx_user_id (user_id),
    INDEX idx_deleted_at (deleted_at)
) DEFAULT CHARSET=utf8mb4 COMMENT='评论表';

CREATE TABLE IF NOT EXISTS last_logout_times (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    logout_time DATETIME NOT NULL COMMENT '登出时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_user_id (user_id)
) DEFAULT CHARSET=utf8mb4 COMMENT='最后登出时间表';

CREATE TABLE IF NOT EXISTS private_messages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    from_user_id BIGINT UNSIGNED NOT NULL COMMENT '发送者ID',
    to_user_id BIGINT UNSIGNED NOT NULL COMMENT '接收者ID',
    content TEXT NOT NULL COMMENT '消息内容',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME NULL COMMENT '删除时间',
    INDEX idx_from_to_created (from_user_id, to_user_id, created_at),
    INDEX idx_to_from_created (to_user_id, from_user_id, created_at),
    INDEX idx_deleted_at (deleted_at)
) DEFAULT CHARSET=utf8mb4 COMMENT='私聊消息表';

CREATE TABLE IF NOT EXISTS group_messages (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    from_user_id BIGINT UNSIGNED NOT NULL COMMENT '发送者ID',
    group_id BIGINT UNSIGNED NOT NULL COMMENT '群组ID',
    content TEXT NOT NULL COMMENT '消息内容',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME NULL COMMENT '删除时间',
    INDEX idx_group_created (group_id, created_at),
    INDEX idx_deleted_at (deleted_at)
) DEFAULT CHARSET=utf8mb4 COMMENT='群聊消息表';

CREATE TABLE IF NOT EXISTS group_members (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    group_id BIGINT UNSIGNED NOT NULL COMMENT '群组ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME NULL COMMENT '删除时间',
    UNIQUE KEY uk_group_user (group_id, user_id),
    INDEX idx_group_id (group_id),
    INDEX idx_user_id (user_id),
    INDEX idx_deleted_at (deleted_at)
) DEFAULT CHARSET=utf8mb4 COMMENT='群成员表';