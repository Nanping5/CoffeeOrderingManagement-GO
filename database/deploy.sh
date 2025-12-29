#!/bin/bash
# ============================================
# 咖啡点餐系统 - 数据库部署脚本
# ============================================

echo "========================================="
echo "咖啡点餐系统 - 数据库部署"
echo "========================================="
echo ""

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SQL_FILE="$SCRIPT_DIR/setup.sql"

# 检查SQL文件是否存在
if [ ! -f "$SQL_FILE" ]; then
    echo "错误: 找不到 setup.sql 文件"
    echo "路径: $SQL_FILE"
    exit 1
fi

echo "找到SQL文件: $SQL_FILE"
echo ""

# 提示用户输入密码
echo "请输入MySQL root密码:"
echo "（如果直接按回车，将尝试无密码连接）"
read -s MYSQL_PASSWORD

echo ""
echo "正在部署数据库..."
echo "----------------------------------------"

# 执行SQL文件
if [ -z "$MYSQL_PASSWORD" ]; then
    # 无密码连接
    mysql -u root < "$SQL_FILE" 2>&1
else
    # 有密码连接
    mysql -u root -p"$MYSQL_PASSWORD" < "$SQL_FILE" 2>&1
fi

# 检查执行结果
if [ $? -eq 0 ]; then
    echo ""
    echo "========================================="
    echo "数据库部署成功！"
    echo "========================================="
    echo ""
    echo "管理员账户: admin / admin123"
    echo ""
else
    echo ""
    echo "========================================="
    echo "数据库部署失败！"
    echo "========================================="
    echo ""
    echo "请检查:"
    echo "1. MySQL服务是否已启动"
    echo "2. root密码是否正确"
    echo "3. setup.sql文件是否存在"
    echo ""
    exit 1
fi
