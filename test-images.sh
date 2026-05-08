#!/bin/bash
# 图片显示问题诊断脚本
# 用于检查上传的图片是否能正确访问

echo "=========================================="
echo "🔍 Orange 系统图片显示诊断工具"
echo "=========================================="
echo ""

# 配置信息
SERVER_IP="123.56.99.154"
BASE_URL="http://${SERVER_IP}:8080"
API_URL="${BASE_URL}/api"

echo "📡 测试服务器连接..."
if curl -s --connect-timeout 5 "${BASE_URL}" > /dev/null 2>&1; then
    echo "✅ 前端服务正常 (${BASE_URL})"
else
    echo "❌ 前端服务无法访问 (${BASE_URL})"
fi

if curl -s --connect-timeout 5 "${API_URL}/health" > /dev/null 2>&1; then
    echo "✅ 后端API正常 (${API_URL})"
else
    echo "❌ 后端API无法访问 (${API_URL})"
fi

echo ""
echo "📁 检查 /uploads/ 路径是否可访问..."
UPLOAD_TEST=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/uploads/default-avatar.png" 2>/dev/null)
if [ "$UPLOAD_TEST" = "200" ]; then
    echo "✅ /uploads/ 路径可正常访问 (HTTP ${UPLOAD_TEST})"
elif [ "$UPLOAD_TEST" = "404" ]; then
    echo "⚠️  /uploads/ 路径返回404 (可能是默认头像不存在)"
    echo "   但路径本身应该是通的"
else
    echo "❌ /uploads/ 路径异常 (HTTP ${UPLOAD_TEST})"
fi

echo ""
echo "🗄️  从数据库获取实际图片URL样本..."

# 尝试登录获取token
LOGIN_RESPONSE=$(curl -s -X POST "${API_URL}/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin123"}' 2>/dev/null)

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$TOKEN" ] && [ "$TOKEN" != "" ] && [ "$TOKEN" != "null" ]; then
    echo "✅ 登录成功"

    # 获取用户列表中的头像
    echo ""
    echo "👤 检查用户头像..."
    USERS_RESPONSE=$(curl -s "${API_URL}/users/list" \
        -H "Authorization: Bearer ${TOKEN}" 2>/dev/null)

    # 提取前3个用户的avatar字段
    AVATARS=$(echo "$USERS_RESPONSE" | grep -o '"avatar":"[^"]*"' | head -3)
    if [ -n "$AVATARS" ]; then
        echo "$AVATARS" | while read line; do
            AVATAR_URL=$(echo "$line" | cut -d'"' -f4)
            if [ -n "$AVATAR_URL" ] && [ "$AVATAR_URL" != "" ] && [ "$AVATAR_URL" != "null" ]; then
                # 构建完整URL
                if [[ "$AVATAR_URL" == /* ]]; then
                    FULL_URL="${BASE_URL}${AVATAR_URL}"
                else
                    FULL_URL="${AVATAR_URL}"
                fi

                # 测试访问
                HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${FULL_URL}" 2>/dev/null)
                if [ "$HTTP_CODE" = "200" ]; then
                    echo "✅ 用户头像可访问: ${AVATAR_URL} (HTTP ${HTTP_CODE})"
                else
                    echo "❌ 用户头像不可访问: ${AVATAR_URL} (HTTP ${HTTP_CODE})"
                fi
            fi
        done
    else
        echo "⚠️  未找到用户头像数据"
    fi

    # 获取团队列表中的logo和背景图
    echo ""
    echo "🏢 检查团队Logo和背景图..."
    TEAMS_RESPONSE=$(curl -s "${API_URL}/teams" \
        -H "Authorization: Bearer ${TOKEN}" 2>/dev/null)

    LOGOS=$(echo "$TEAMS_RESPONSE" | grep -o '"logo":"[^"]*"' | head -3)
    BG_IMAGES=$(echo "$TEAMS_RESPONSE" | grep -o '"bg_image":"[^"]*"' | head -3)

    if [ -n "$LOGOS" ]; then
        echo ""
        echo "📷 团队Logo:"
        echo "$LOGOS" | while read line; do
            LOGO_URL=$(echo "$line" | cut -d'"' -f4)
            if [ -n "$LOGO_URL" ] && [ "$LOGO_URL" != "" ] && [ "$LOGO_URL" != "null" ]; then
                if [[ "$LOGO_URL" == /* ]]; then
                    FULL_URL="${BASE_URL}${LOGO_URL}"
                else
                    FULL_URL="${LOGO_URL}"
                fi

                HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${FULL_URL}" 2>/dev/null)
                if [ "$HTTP_CODE" = "200" ]; then
                    echo "✅ 团队Logo可访问: ${LOGO_URL} (HTTP ${HTTP_CODE})"
                else
                    echo "❌ 团队Logo不可访问: ${LOGO_URL} (HTTP ${HTTP_CODE})"
                fi
            fi
        done
    else
        echo "⚠️  未找到团队Logo数据"
    fi

    if [ -n "$BG_IMAGES" ]; then
        echo ""
        echo "🖼️  团队背景图:"
        echo "$BG_IMAGES" | while read line; do
            BG_URL=$(echo "$line" | cut -d'"' -f4)
            if [ -n "$BG_URL" ] && [ "$BG_URL" != "" ] && [ "$BG_URL" != "null" ]; then
                if [[ "$BG_URL" == /* ]]; then
                    FULL_URL="${BASE_URL}${BG_URL}"
                else
                    FULL_URL="${BG_URL}"
                fi

                HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${FULL_URL}" 2>/dev/null)
                if [ "$HTTP_CODE" = "200" ]; then
                    echo "✅ 团队背景图可访问: ${BG_URL} (HTTP ${HTTP_CODE})"
                else
                    echo "❌ 团队背景图不可访问: ${BG_URL} (HTTP ${HTTP_CODE})"
                fi
            fi
        done
    else
        echo "⚠️  未找到团队背景图数据"
    fi

    # 获取博主头像
    echo ""
    echo "📝 检查博主头像..."
    BLOGGERS_RESPONSE=$(curl -s "${API_URL}/blogger/list?page=1&pageSize=5" \
        -H "Authorization: Bearer ${TOKEN}" 2>/dev/null)

    BLOGGER_AVATARS=$(echo "$BLOGGERS_RESPONSE" | grep -o '"avatar":"[^"]*"' | head -5)
    if [ -n "$BLOGGER_AVATARS" ]; then
        echo "$BLOGGER_AVATARS" | while read line; do
            BLOGGER_AVATAR=$(echo "$line" | cut -d'"' -f4)
            if [ -n "$BLOGGER_AVATAR" ] && [ "$BLOGGER_AVATAR" != "" ] && [ "$BLOGGER_AVATAR" != "null" ]; then
                if [[ "$BLOGGER_AVATAR" == /* ]]; then
                    FULL_URL="${BASE_URL}${BLOGGER_AVATAR}"
                else
                    FULL_URL="${BLOGGER_AVATAR}"
                fi

                HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${FULL_URL}" 2>/dev/null)
                if [ "$HTTP_CODE" = "200" ]; then
                    echo "✅ 博主头像可访问: ${BLOGGER_AVATAR} (HTTP ${HTTP_CODE})"
                else
                    echo "❌ 博主头像不可访问: ${BLOGGER_AVATAR} (HTTP ${HTTP_CODE})"
                fi
            fi
        done
    else
        echo "⚠️  未找到博主头像数据"
    fi

else
    echo "❌ 登录失败，无法获取详细数据"
    echo "   响应: ${LOGIN_RESPONSE}"
fi

echo ""
echo "=========================================="
echo "🔧 常见问题及解决方案"
echo "=========================================="
echo ""
echo "1️⃣  如果图片返回 404:"
echo "   → 文件未正确保存到 uploads 目录"
echo "   → 检查 Docker 卷挂载: ./uploads:/app/uploads"
echo "   → 运行: docker exec orange-backend ls -la /app/uploads/"
echo ""
echo "2️⃣  如果图片返回 403:"
echo "   → 权限问题或认证中间件拦截"
echo "   → 检查 Nginx /uploads/ 配置"
echo "   → 检查后端 Static 路由配置"
echo ""
echo "3️⃣  如果图片路径为空或 null:"
echo "   → 数据库中未存储图片路径"
echo "   → 上传功能可能失败但未报错"
echo "   → 需要重新上传图片"
echo ""
echo "4️⃣  如果图片路径格式异常:"
echo "   → 可能包含完整URL或相对路径错误"
echo "   → 应该以 /uploads/ 开头"
echo ""
echo "5️⃣  Content-Security-Policy 问题:"
echo "   → Nginx CSP 设置了 img-src 'self'"
echo "   → 外部图片会被阻止"
echo "   → 所有图片必须来自同源 (/uploads/)"
echo ""
echo "=========================================="
