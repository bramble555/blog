# 🚀 GVB Blog - 现代化全栈博客系统

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.0+-4FC08D?style=flat-square&logo=vue.js)](https://vuejs.org/)
[![Elasticsearch](https://img.shields.io/badge/Elasticsearch-8.19+-005571?style=flat-square&logo=elasticsearch)](https://www.elastic.co/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)

GVB Blog 是一款专为开发者设计的高性能、现代化全栈博客系统。后端基于 Go 语言进行极致性能优化，前端采用 Vue 3 打造流畅的交互体验。

---

## ✨ 核心特性

### 📝 文章管理
- **Markdown 支持**: Markdown 写作与安全渲染。
- **自动摘要**: 基于内容自动生成文章简介。
- **动态封面**: 从横幅库随机选择文章封面。
- **标签与日历**: 标签关联与文章日历统计。

### 🔍 搜索能力
- **Elasticsearch 集成**: 标题/内容关键词检索与高亮展示。
- **MySQL 自动降级**: ES 不可用时自动切换到 `LIKE` 查询。

### 👥 用户与社交

- **权限安全**: JWT 鉴权与角色管理（Admin/User）。
- **账号能力**: 注册/登录、密码修改、邮箱验证码绑定。
- **互动体验**: 文章点赞、收藏、评论与消息系统。
- **即时通信**: WebSocket 群聊与聊天记录回放。

### 📊 系统运维
- **运行指标**: 通过 `/api/debug/vars` 暴露运行指标与连接池状态。
- **分布式 ID**: 雪花算法（Snowflake）生成全局唯一 ID。
- **流量控制**: 基于 IP 的接口限流机制（Redis/Memory），防止恶意刷刷接口。
- **配置管理**: 基于 Viper 实现多环境配置读取与热加载。
- **日志系统**: Logrus 结合日志切割与归档，记录关键操作流水。
- **自动同步**: MySQL → Elasticsearch 定时同步任务（Cron）。
- **敏感词过滤**: 自己写了前缀树实现评论内容过滤。
- **平滑关闭**: 优雅退出，确保请求处理完成。

### 🛡 安全机制
- **令牌管理**: JWT 鉴权机制.
- **密码安全**: BCrypt 密码加密存储。
- **内容安全**: Markdown 渲染 XSS 过滤。
- **跨域处理**: 严格的 CORS 策略配置。

---

## 🛠 技术栈

### 后端 (Backend)
- **核心框架**: Go (Gin)
- **数据库**: MySQL / Redis
- **全文搜索**: Elasticsearch 8.19.x
- **持久层**: GORM
- **鉴权与安全**: JWT / bcrypt / Session
- **实时通信**: WebSocket
- **任务调度**: Cron
- **内容处理**: Markdown 渲染与 XSS 过滤
- **工具集**: Logrus, Viper, Snowflake

### 前端 (Frontend)
- **框架**: Vue 3 (Composition API)
- **构建工具**: Vite
- **状态管理**: Pinia
- **UI 组件**: Element Plus
- **样式**: Tailwind CSS 4
- **图表**: Echarts

---

## 🚀 快速开始

### 1. 环境依赖
- **Go** 1.23+
- **Node.js** 18+
- **MySQL** 8.0+
- **Elasticsearch** 8.x（可选）
- **Redis**

### 2. 后端配置
1. 克隆代码库。
2. 配置 `setting.yaml` 中的数据库、Redis、Elasticsearch 连接信息。
3. 首次启动时会自动执行 `init.sql` 初始化表结构。
4. 如不使用 ES，可将 `elasticsearch.enable` 设为 `false`。
5. 运行服务:
   ```bash
   go run main.go
   ```
   *推荐开发中使用 `air` 实现热重载。*

### 3. 前端部署
1. 进入 `frontend` 目录。
2. 安装依赖:
   ```bash
   npm install
   ```
3. 启动开发服务器:
   ```bash
   npm run dev
   ```

### 4. 命令行工具 (CLI)
直接通过终端创建管理员账户:
```bash
go run main.go -n=admin -p=password123 -r=1
```

## ✅ 测试
- **功能范围**: 覆盖敏感词过滤、评论内容处理（敏感词替换 + 特殊字符转义）、边界输入与并发安全场景。
- **使用方法**:
  ```bash
  go test ./test -run TestProcessCommentContent
  go test ./test
  ```
- **注意事项**:
  - 依赖项目根目录下的 `filter.txt` 词库文件。
  - 建议在项目根目录执行测试，已包含自动路径回退以适配 `test` 包运行目录。

## 📄 开源协议
本项目采用 [MIT License](LICENSE) 开源协议。
由 **Go** & **Vue 3** 驱动。Made with ❤️ by Bramble555.
