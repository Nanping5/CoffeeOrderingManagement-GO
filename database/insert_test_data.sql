-- ============================================
-- 插入测试数据 - 用于论文查询演示
-- ============================================

USE coffee_ordering;

-- 1. 插入测试用户（除了 admin 之外的普通用户）
-- 注意：触发器会自动为这些用户创建积分账户
INSERT INTO users (username, email, password, phone, role, is_active, is_verified) VALUES
('alice', 'alice@example.com', '$2a$10$nqeLcqA2u0be4ZO5TL2xvujnOTAAeZISDzVewv7UAiJN6pSuLB8xm', '13800000001', 'user', TRUE, TRUE),
('bob', 'bob@example.com', '$2a$10$nqeLcqA2u0be4ZO5TL2xvujnOTAAeZISDzVewv7UAiJN6pSuLB8xm', '13800000002', 'user', TRUE, TRUE),
('charlie', 'charlie@example.com', '$2a$10$nqeLcqA2u0be4ZO5TL2xvujnOTAAeZISDzVewv7UAiJN6pSuLB8xm', '13800000003', 'user', TRUE, TRUE),
('diana', 'diana@example.com', '$2a$10$nqeLcqA2u0be4ZO5TL2xvujnOTAAeZISDzVewv7UAiJN6pSuLB8xm', '13800000004', 'user', TRUE, TRUE),
('eve', 'eve@example.com', '$2a$10$nqeLcqA2u0be4ZO5TL2xvujnOTAAeZISDzVewv7UAiJN6pSuLB8xm', '13800000005', 'user', TRUE, TRUE);

-- 2. 更新用户积分账户（使用 UPDATE 而不是 INSERT，因为触发器已自动创建）
UPDATE user_points SET total_points = 500, lifetime_points = 500, member_level = 'bronze' WHERE user_id = 2;
UPDATE user_points SET total_points = 1500, lifetime_points = 1500, member_level = 'silver' WHERE user_id = 3;
UPDATE user_points SET total_points = 6000, lifetime_points = 6000, member_level = 'gold' WHERE user_id = 4;
UPDATE user_points SET total_points = 25000, lifetime_points = 25000, member_level = 'platinum' WHERE user_id = 5;
UPDATE user_points SET total_points = 800, lifetime_points = 800, member_level = 'bronze' WHERE user_id = 6;

-- 3. 插入测试订单
INSERT INTO orders (user_id, order_number, pickup_code, status, customer_name, customer_phone,
                   customer_points_used, points_deduction_amount, points_earned, member_level_at_time) VALUES
(2, 'ORD20250129001', 'A001', 'completed', '张三', '13800000001', 0, 0.00, 36, 'bronze'),
(2, 'ORD20250129002', 'A002', 'completed', '张三', '13800000001', 100, 5.00, 44, 'bronze'),
(3, 'ORD20250129003', 'A003', 'completed', '李四', '13800000002', 200, 10.00, 52, 'silver'),
(3, 'ORD20250129004', 'A004', 'pending', '李四', '13800000002', 0, 0.00, 0, 'silver'),
(4, 'ORD20250129005', 'A005', 'ready', '王五', '13800000003', 500, 25.00, 78, 'gold'),
(5, 'ORD20250129006', 'A006', 'preparing', '赵六', '13800000004', 1000, 75.00, 130, 'platinum'),
(NULL, 'ORD20250129007', 'A007', 'completed', '游客订单', '13900000000', 0, 0.00, 18, NULL),
(6, 'ORD20250129008', 'A008', 'pending', '小红', '13800000005', 0, 0.00, 0, 'bronze');

-- 4. 插入订单明细
-- 订单1: 美式咖啡 x2
INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price) VALUES
(1, 1, 2, 18.00);

-- 订单2: 拿铁咖啡 x1 + 卡布奇诺 x1 + 摩卡咖啡 x1
INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price) VALUES
(2, 2, 1, 22.00),
(2, 3, 1, 24.00),
(2, 4, 1, 26.00);

-- 订单3: 焦糖玛奇朵 x2 + 抹茶拿铁 x1
INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price) VALUES
(3, 5, 2, 28.00),
(3, 8, 1, 25.00);

-- 订单4: 冰咖啡 x3
INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price) VALUES
(4, 7, 3, 20.00);

-- 订单5: 拿铁咖啡 x2 + 蓝莓松饼 x2 + 巧克力蛋糕 x1
INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price) VALUES
(5, 2, 2, 22.00),
(5, 12, 2, 15.00),
(5, 13, 1, 20.00);

-- 订单6: 浓缩咖啡 x2 + 摩卡咖啡 x2 + 芝士蛋糕 x2 + 三明治 x1
INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price) VALUES
(6, 6, 2, 15.00),
(6, 4, 2, 26.00),
(6, 14, 2, 22.00),
(6, 16, 1, 32.00);

-- 订单7: 美式咖啡 x1 + 柠檬茶 x1
INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price) VALUES
(7, 1, 1, 18.00),
(7, 9, 1, 16.00);

-- 订单8: 红茶 x2 + 热巧克力 x1
INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price) VALUES
(8, 10, 2, 14.00),
(8, 11, 1, 18.00);

-- 5. 插入积分交易记录
INSERT INTO point_transactions (user_id, order_id, transaction_type, points_change, points_balance, description) VALUES
-- Alice 的积分记录
(2, 1, 'earned', 36, 36, '订单完成获得积分'),
(2, 2, 'used', -100, 0, '订单积分抵扣'),
(2, 2, 'earned', 44, 44, '订单完成获得积分'),

-- Bob 的积分记录
(3, 3, 'used', -200, 0, '订单积分抵扣'),
(3, 3, 'earned', 52, 52, '订单完成获得积分'),

-- Charlie 的积分记录
(4, 5, 'used', -500, 0, '订单积分抵扣'),
(4, 5, 'earned', 78, 78, '订单完成获得积分'),

-- Diana 的积分记录
(5, 6, 'used', -1000, 0, '订单积分抵扣'),
(5, 6, 'earned', 130, 130, '订单完成获得积分'),

-- 游客订单（无积分）
-- (NULL, 7, ...) - 游客没有积分记录

-- Eve 的积分记录（订单还未完成，所以没有积分获得）
-- (6, 8, ...) - 订单 pending，暂无积分

-- 6. 更新订单状态（模拟订单完成）
-- 注意：实际场景中这些应该通过业务逻辑自动触发
UPDATE orders SET status = 'completed' WHERE id IN (1, 2, 3, 5, 7);

SELECT '测试数据插入完成！' AS message;
SELECT '用户数量:' AS info, COUNT(*) AS count FROM users;
SELECT '订单数量:' AS info, COUNT(*) AS count FROM orders;
SELECT '订单明细数量:' AS info, COUNT(*) AS count FROM order_items;
SELECT '积分交易数量:' AS info, COUNT(*) AS count FROM point_transactions;
