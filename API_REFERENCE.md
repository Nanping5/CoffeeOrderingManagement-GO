# API 接口文档

**Base URL**: `http://localhost:8081/api`

## 认证方式

用户接口需要 JWT Token：
```
Authorization: Bearer <token>
```

管理员接口支持 JWT 或简化 Token：
```
Authorization: Bearer fake-admin-token-12345
```

---

## 公开接口

### 健康检查
```
GET /health
```

### 菜单

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /menu | 获取菜单列表 |
| GET | /menu/categories | 获取分类列表 |
| GET | /menu/:id | 获取商品详情 |

**GET /menu 参数**：
- `page` - 页码（默认 1）
- `per_page` - 每页数量（默认 20）
- `category` - 分类筛选
- `keyword` - 关键词搜索

### 订单

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /orders | 创建订单 |
| GET | /orders/:id | 获取订单详情 |
| GET | /orders/pickup/:code | 取餐码查询 |

**POST /orders 请求体**：
```json
{
  "items": [{"menu_id": 1, "quantity": 2, "unit_price": 18.00}],
  "total_price": 36.00,
  "notes": "少糖",
  "user_id": 1,
  "points_used": 100
}
```

---

## 用户认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /auth/register | 用户注册 |
| POST | /auth/login | 用户登录 |
| POST | /auth/admin/login | 管理员登录 |

**POST /auth/register**：
```json
{
  "username": "user1",
  "email": "user@example.com",
  "password": "password123",
  "phone": "13800138000"
}
```

**POST /auth/login**：
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```


---

## 用户接口（需要登录）

### 个人信息

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /user/profile | 获取个人信息 |
| PUT | /user/profile | 更新个人信息 |
| PUT | /user/password | 修改密码 |
| GET | /user/orders | 获取我的订单 |

### 积分系统

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /user/points | 获取积分信息 |
| GET | /user/points/transactions | 获取积分记录 |
| GET | /user/points/level | 获取会员等级 |

**GET /user/points 响应**：
```json
{
  "success": true,
  "data": {
    "total_points": 1500,
    "available_points": 1200,
    "frozen_points": 300,
    "member_level": "silver",
    "level_display_name": "银牌会员"
  }
}
```

---

## 管理员接口

### 菜单管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /admin/menu | 获取所有菜单 |
| POST | /admin/menu | 创建菜单项 |
| PUT | /admin/menu/:id | 更新菜单项 |
| DELETE | /admin/menu/:id | 删除菜单项 |
| PATCH | /admin/menu/:id/toggle | 切换上下架 |

### 订单管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /admin/orders | 获取所有订单 |
| GET | /admin/orders/:id | 获取订单详情 |
| PUT | /admin/orders/:id/status | 更新订单状态 |
| DELETE | /admin/orders/:id | 删除订单 |
| GET | /admin/orders/statistics | 订单统计 |

**订单状态**：`pending` → `preparing` → `ready` → `completed` / `cancelled`

**GET /admin/orders/statistics 响应**：
```json
{
  "success": true,
  "data": {
    "total_orders": 150,
    "total_revenue": 4500.00,
    "today_orders": 25,
    "today_revenue": 750.00,
    "pending_count": 10,
    "preparing_count": 5,
    "ready_count": 3,
    "completed_count": 130,
    "cancelled_count": 2,
    "total_users": 50,
    "member_levels": {
      "bronze": 30,
      "silver": 15,
      "gold": 4,
      "platinum": 1
    },
    "top_products": [
      {"menu_name": "拿铁咖啡", "quantity": 120, "revenue": 2640.00}
    ]
  }
}
```

---

## 响应格式

**成功**：
```json
{
  "success": true,
  "data": {...},
  "message": "操作成功"
}
```

**失败**：
```json
{
  "success": false,
  "message": "错误信息"
}
```

## HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 404 | 资源不存在 |
| 500 | 服务器错误 |

---

**版本**: 3.0.0
