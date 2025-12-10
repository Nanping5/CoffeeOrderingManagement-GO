-- ============================================
-- 咖啡点餐管理系统 - 一键部署数据库脚本
-- 执行方式: mysql -u root -p < setup.sql
-- ============================================

-- 创建数据库
DROP DATABASE IF EXISTS coffee_ordering;
CREATE DATABASE coffee_ordering CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE coffee_ordering;

-- ============================================
-- 1. 用户表
-- ============================================
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(20) DEFAULT 'user',
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    avatar_url VARCHAR(255),
    gender ENUM('male', 'female', 'other'),
    birth_date DATE,
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_email (email),
    INDEX idx_username (username),
    INDEX idx_role (role),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 2. 菜单表
-- ============================================
CREATE TABLE menu_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    category VARCHAR(50) NOT NULL DEFAULT 'coffee',
    image_url VARCHAR(255),
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category (category),
    INDEX idx_available (is_available)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 3. 订单表
-- ============================================
CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NULL,
    pickup_code VARCHAR(10) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    status ENUM('pending', 'preparing', 'ready', 'completed', 'cancelled') DEFAULT 'pending',
    customer_name VARCHAR(100),
    customer_phone VARCHAR(20),
    notes TEXT,
    customer_points_used INT DEFAULT 0,
    points_deduction_amount DECIMAL(10,2) DEFAULT 0.00,
    points_earned INT DEFAULT 0,
    original_total_price DECIMAL(10,2),
    final_payment_amount DECIMAL(10,2),
    member_level_at_time ENUM('bronze','silver','gold','platinum'),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_pickup_code (pickup_code),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 4. 订单明细表
-- ============================================
CREATE TABLE order_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,
    menu_item_id INT NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items(id) ON DELETE RESTRICT,
    INDEX idx_order_id (order_id),
    INDEX idx_menu_item_id (menu_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 5. 会员积分表
-- ============================================
CREATE TABLE user_points (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,
    total_points INT DEFAULT 0,
    available_points INT DEFAULT 0,
    frozen_points INT DEFAULT 0,
    lifetime_points INT DEFAULT 0,
    member_level ENUM('bronze', 'silver', 'gold', 'platinum') DEFAULT 'bronze',
    level_upgrade_date TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_member_level (member_level)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 6. 积分变动记录表
-- ============================================
CREATE TABLE point_transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    order_id INT NULL,
    transaction_type ENUM('earned', 'used', 'expired', 'refunded', 'signup_bonus', 'birthday_bonus', 'referral_bonus') NOT NULL,
    points_change INT NOT NULL,
    points_balance INT NOT NULL,
    description VARCHAR(255) NOT NULL,
    reference_id VARCHAR(100),
    metadata JSON,
    expires_at TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_type (transaction_type),
    INDEX idx_order_id (order_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 7. 会员等级配置表
-- ============================================
CREATE TABLE member_levels (
    id INT AUTO_INCREMENT PRIMARY KEY,
    level_name ENUM('bronze', 'silver', 'gold', 'platinum') NOT NULL UNIQUE,
    level_display_name VARCHAR(50) NOT NULL,
    min_points INT NOT NULL,
    points_earning_rate DECIMAL(5,4) NOT NULL DEFAULT 1.0000,
    points_discount_rate DECIMAL(5,4) NOT NULL DEFAULT 0.0000,
    max_discount_percentage DECIMAL(5,2) NOT NULL DEFAULT 20.00,
    benefits JSON,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 8. 插入默认会员等级配置
-- ============================================
INSERT INTO member_levels (level_name, level_display_name, min_points, points_earning_rate, points_discount_rate, max_discount_percentage, benefits, sort_order) VALUES
('bronze', '铜牌会员', 0, 1.0000, 0.0000, 20.00, '{"description": "基础会员权益", "features": ["积分累积", "生日礼品"]}', 1),
('silver', '银牌会员', 1000, 1.2000, 0.0500, 30.00, '{"description": "银牌会员专享权益", "features": ["积分累积", "生日双倍积分", "专属优惠券"]}', 2),
('gold', '金牌会员', 5000, 1.5000, 0.1000, 40.00, '{"description": "金牌尊贵权益", "features": ["积分累积", "生日三倍积分", "专属优惠券", "新品优先体验"]}', 3),
('platinum', '白金会员', 20000, 2.0000, 0.1500, 50.00, '{"description": "白金顶级权益", "features": ["积分累积", "生日五倍积分", "专属客服", "免费配送"]}', 4);

-- ============================================
-- 9. 插入管理员账户 (密码: admin123)
-- ============================================
INSERT INTO users (username, email, password, role, is_active, is_verified) VALUES
('admin', 'admin@coffee.com', '$2a$10$nqeLcqA2u0be4ZO5TL2xvujnOTAAeZISDzVewv7UAiJN6pSuLB8xm', 'admin', TRUE, TRUE);

-- ============================================
-- 10. 插入示例菜单数据
-- ============================================
INSERT INTO menu_items (name, description, price, category, image_url, is_available) VALUES
('美式咖啡', '经典美式咖啡，口感浓郁醇厚', 18.00, '咖啡', 'https://images.unsplash.com/photo-1514432324607-a09d9b4aefdd?w=400', TRUE),
('拿铁咖啡', '意式浓缩咖啡配蒸汽牛奶，口感顺滑', 22.00, '咖啡', 'https://images.unsplash.com/photo-1461023058943-07fcbe16d735?w=400', TRUE),
('卡布奇诺', '浓缩咖啡、蒸汽牛奶和奶泡的完美组合', 24.00, '咖啡', 'https://images.unsplash.com/photo-1572442388796-11668a67e53d?w=400', TRUE),
('摩卡咖啡', '巧克力与咖啡的完美融合', 26.00, '咖啡', 'https://images.unsplash.com/photo-1578314675249-a6910f80cc4e?w=400', TRUE),
('焦糖玛奇朵', '香草糖浆和焦糖酱的甜蜜点缀', 28.00, '咖啡', 'https://images.unsplash.com/photo-1485808191679-5f86510681a2?w=400', TRUE),
('浓缩咖啡', '纯正意式浓缩咖啡，强烈浓郁', 15.00, '咖啡', 'https://images.unsplash.com/photo-1510707577719-ae7c14805e3a?w=400', TRUE),
('冰咖啡', '清爽冰镇咖啡，夏日首选', 20.00, '咖啡', 'https://images.unsplash.com/photo-1517701550927-30cf4ba1dba5?w=400', TRUE),
('抹茶拿铁', '日式抹茶与牛奶的清新搭配', 25.00, '茶饮', 'https://images.unsplash.com/photo-1515823064-d6e0c04616a7?w=400', TRUE),
('柠檬茶', '清新柠檬茶，解渴提神', 16.00, '茶饮', 'https://images.unsplash.com/photo-1556679343-c7306c1976bc?w=400', TRUE),
('红茶', '醇厚红茶，经典英式风味', 14.00, '茶饮', 'https://images.unsplash.com/photo-1594631252845-29fc4cc8cde9?w=400', TRUE),
('热巧克力', '浓郁热巧克力，温暖甜蜜', 18.00, '其他', 'https://images.unsplash.com/photo-1542990253-0d0f5be5f0ed?w=400', TRUE),
('蓝莓松饼', '新鲜蓝莓制作的松软松饼', 15.00, '甜点', 'https://images.unsplash.com/photo-1607958996333-41aef7caefaa?w=400', TRUE),
('巧克力蛋糕', '浓郁巧克力蛋糕，甜而不腻', 20.00, '甜点', 'https://images.unsplash.com/photo-1578985545062-69928b1d9587?w=400', TRUE),
('芝士蛋糕', '纽约风味芝士蛋糕，口感丰富', 22.00, '甜点', 'https://images.unsplash.com/photo-1533134242443-d4fd215305ad?w=400', TRUE),
('羊角面包', '法式经典羊角面包，酥脆可口', 12.00, '甜点', 'https://images.unsplash.com/photo-1555507036-ab1f4038808a?w=400', TRUE),
('三明治', '新鲜蔬菜和优质肉类制作', 32.00, '其他', 'https://images.unsplash.com/photo-1528735602780-2552fd46c7af?w=400', TRUE),
('橙汁', '鲜榨橙汁，维C丰富', 12.00, '其他', 'https://images.unsplash.com/photo-1621506289937-a8e4df240d0b?w=400', TRUE),
('矿泉水', '天然矿泉水，纯净清爽', 5.00, '其他', 'https://images.unsplash.com/photo-1548839140-29a749e1cf4d?w=400', TRUE);

-- ============================================
-- 11. 创建触发器（用户注册后自动创建积分账户）
-- ============================================
DROP TRIGGER IF EXISTS after_user_insert;
DELIMITER //
CREATE TRIGGER after_user_insert
AFTER INSERT ON users
FOR EACH ROW
BEGIN
    INSERT INTO user_points (user_id, total_points, available_points, frozen_points, lifetime_points, member_level)
    VALUES (NEW.id, 0, 0, 0, 0, 'bronze');
END //
DELIMITER ;

-- ============================================
-- 完成提示
-- ============================================
SELECT '========================================' as '';
SELECT '咖啡点餐系统数据库部署完成！' as message;
SELECT '========================================' as '';
SELECT '管理员账户: admin / admin123' as credentials;
SELECT CONCAT('菜单项数量: ', COUNT(*)) as menu_count FROM menu_items;
SELECT CONCAT('用户数量: ', COUNT(*)) as user_count FROM users;
