# 咖啡点餐管理系统

> 现代化的自助点餐系统，使用 Go + Gin + Vue 3 + Element Plus 构建

## 🚀 快速开始

### 1. 部署数据库

```bash
mysql -u root -p < database/setup.sql
```

### 2. 启动后端

```bash
cd backend-go
cp .env.example .env  # 配置数据库连接
go mod download
go run main.go
```
服务运行在 `http://localhost:8081`

### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```
服务运行在 `http://localhost:3001`

### 4. 访问系统

| 入口 | 地址 | 账号 |
|------|------|------|
| 顾客端 | http://localhost:3001 | 无需登录 |
| 管理后台 | http://localhost:3001/admin/login | admin / admin123 |

## 📋 功能特性

### 顾客端
- 浏览菜单（分类筛选、排序）
- 购物车管理
- 下单获取取餐码
- 会员注册/登录
- 积分累积与抵扣
- 订单历史查询

### 管理后台
- 仪表板数据统计
- 菜单管理（CRUD、上下架）
- 订单管理（状态更新）
- 会员数据统计
- 系统设置


## 🏗️ 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.21+, Gin, GORM, JWT |
| 前端 | Vue 3, Vite, Element Plus, Pinia |
| 数据库 | MySQL 8.0+ |
| 样式 | SCSS, CSS Variables |

## 📁 项目结构

```
CoffeeOrderingManagement/
├── backend-go/           # Go 后端
│   ├── config/          # 配置
│   ├── database/        # 数据库连接
│   ├── handlers/        # 请求处理
│   ├── middleware/      # 中间件
│   ├── models/          # 数据模型
│   ├── routes/          # 路由
│   ├── services/        # 业务逻辑
│   └── main.go
├── frontend/            # Vue 前端
│   ├── src/
│   │   ├── api/        # API 接口
│   │   ├── components/ # 组件
│   │   ├── views/      # 页面
│   │   ├── store/      # 状态管理
│   │   ├── router/     # 路由
│   │   └── styles/     # 样式
│   └── package.json
└── database/
    └── setup.sql       # 一键部署脚本
```

## 🗄️ 数据库配置

编辑 `backend-go/.env`：

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=coffee_ordering
PORT=8081
CORS_ORIGINS=http://localhost:3001
JWT_SECRET=your-secret-key
```

## 🔧 常用命令

```bash
# 后端
cd backend-go
go run main.go          # 开发运行
go build -o app         # 编译
make run                # Makefile 运行

# 前端
cd frontend
npm run dev             # 开发服务器
npm run build           # 生产构建
npm run preview         # 预览构建
```

## 🐳 Docker 部署

```bash
cd backend-go
docker build -t coffee-backend .
docker run -p 8081:8081 --env-file .env coffee-backend
```

## 📚 文档

- [API 接口文档](API_REFERENCE.md)

## 📄 许可证

MIT License
