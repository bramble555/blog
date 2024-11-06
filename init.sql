
CREATE DATABASE IF NOT EXISTS gvb;

USE gvb;


CREATE TABLE IF NOT EXISTS banner_models (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    hash VARCHAR(32) COMMENT 'Hash 值',
    name VARCHAR(100) COMMENT '文件名',
    image_type TINYINT(1) DEFAULT 1 COMMENT '图片类型,默认是本地',
    CHECK (image_type IN (1, 2))
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;

-- DROP TABLE banner_models;

CREATE TABLE IF NOT EXISTS advert_models (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    title VARCHAR(50) NOT NULL UNIQUE COMMENT '标题',
    href VARCHAR(200) NOT NULL COMMENT '链接',
    images VARCHAR(200) NOT NULL COMMENT '图片',
    is_show TINYINT(1) DEFAULT 1 COMMENT '是否显示,1为是,0为否',
    UNIQUE KEY `idx_title` ( `title` ) USING BTREE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;

-- ALTER TABLE advert_models ADD CONSTRAINT idx_title UNIQUE (title);
CREATE TABLE IF NOT EXISTS menu_models (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    title VARCHAR(32) UNIQUE NOT NULL COMMENT '菜单标题',
    path VARCHAR(32) UNIQUE NOT NULL COMMENT '菜单路径',
    slogan VARCHAR(64) NULL COMMENT '菜单口号或标语',
    abstract VARCHAR(100) NULL COMMENT'菜单简介',
    abstract_time SMALLINT(2) DEFAULT 0 COMMENT '简介的切换时间，单位为秒',
    banner_id BIGINT COMMENT '关联的横幅ID',
    sort TINYINT(5) UNIQUE NOT NULL COMMENT '菜单的顺序,0表示最高优先级',
    UNIQUE KEY `idx_title` (`title`) USING BTREE,
    UNIQUE KEY `idx_path` (`path`) USING BTREE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;
-- DROP TABLE  menu_models;

CREATE TABLE IF NOT EXISTS user_models (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    username VARCHAR(36) NOT NULL UNIQUE COMMENT '用户名，唯一标识',
    password VARCHAR(72) NOT NULL COMMENT '用户密码',
    avatar VARCHAR(256) COMMENT '用户头像URL',
    email VARCHAR(128) NULL COMMENT '用户邮箱',
    phone VARCHAR(18) NULL COMMENT '用户电话', 
    addr VARCHAR(64) NULL COMMENT '用户地址',
    token VARCHAR(64) NULL COMMENT '用户身份令牌',
    ip VARCHAR(20) DEFAULT '127.0.0.1' COMMENT '用户最后登录IP',
    role SMALLINT(1) DEFAULT 1 COMMENT '用户角色,默认值为1',
    sign_status SMALLINT(1) COMMENT '用户签到状态',
    article_id BIGINT COMMENT '文章ID',
    collect_id BIGINT COMMENT '收藏文章ID',
    UNIQUE KEY `idx_username` (`username`) USING BTREE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS article_models (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '文章ID',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    title VARCHAR(255) NOT NULL UNIQUE COMMENT '文章标题',
    abstract TEXT NOT NULL COMMENT '文章简介',
    content TEXT NOT NULL COMMENT '文章内容',
    look_count INT DEFAULT 0 COMMENT '浏览量',
    comment_count INT DEFAULT 0 COMMENT '评论量',
    digg_count INT DEFAULT 0 COMMENT '点赞量',
    collects_count INT DEFAULT 0 COMMENT '收藏量',
    category VARCHAR(20) NOT NULL COMMENT '文章分类',
    source VARCHAR(255) COMMENT '文章来源',
    link VARCHAR(255) COMMENT '原文链接',
    tags TEXT COMMENT '文章标签（以逗号分隔）',
    banner_id BIGINT NOT NULL COMMENT '文章封面ID',
    banner_url VARCHAR(255) COMMENT '封面图片链接',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    username VARCHAR(255) COMMENT '用户名',
    user_avatar VARCHAR(255) COMMENT '用户头像',
    UNIQUE KEY `idx_title` ( `title` ) USING BTREE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS comment_models (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '评论ID',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    content VARCHAR(256) NOT NULL COMMENT '评论内容',
    digg_count INT DEFAULT 0 COMMENT '点赞数',
    comment_count INT DEFAULT 0 COMMENT '子评论数量',
    parent_comment_id BIGINT COMMENT '父级评论ID', -- 用于层级评论
    article_id BIGINT COMMENT '文章ID',
    user_id BIGINT COMMENT '评论的用户ID'
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;
CREATE TABLE IF NOT EXISTS message_models (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '消息 ID',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    send_user_id BIGINT NOT NULL COMMENT '发送人 ID',
    send_username VARCHAR(100) COMMENT '发送人用户名',
    send_user_avatar VARCHAR(255) COMMENT '发送人头像',
    rev_user_id BIGINT NOT NULL COMMENT '接收人 ID',
    rev_username VARCHAR(100) COMMENT '接收人用户名',
    rev_user_avatar VARCHAR(255) COMMENT '接收人头像',
    is_read BOOLEAN DEFAULT FALSE COMMENT '接收方是否查看',
    content TEXT COMMENT '消息内容'
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS tag_models (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '消息 ID',
    title VARCHAR(100) NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_title (title) -- 为 title 字段创建索引
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS article_tag_models (
    article_id BIGINT NOT NULL COMMENT '文章 ID',
    tag_id BIGINT NOT NULL COMMENT '标签 ID',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (article_id, tag_id) COMMENT '复合主键，确保每一对 (article_id, tag_id) 是唯一的'
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



