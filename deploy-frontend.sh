#!/bin/bash
# 部署前端文件到服务器的脚本

SERVER="root@123.56.99.154"
CONTAINER="orange-frontend"
HTML_DIR="/usr/share/nginx/html"

echo "🚀 开始部署前端文件..."

# 上传根目录文件
for file in index.html index.html.gz manifest.json manifest.json.gz offline.html offline.html.gz sw.js sw.js.gz; do
    if [ -f "dist/$file" ]; then
        echo "📤 上传: $file"
        scp "dist/$file" "${SERVER}:/tmp/"
        ssh ${SERVER} "docker cp /tmp/${file} ${CONTAINER}:${HTML_DIR}/"
    fi
done

# 上传assets目录
echo "📦 上传 assets 目录..."
ssh ${SERVER} "docker exec ${CONTAINER} sh -c 'mkdir -p ${HTML_DIR}/assets'"

for file in dist/assets/*; do
    filename=$(basename "$file")
    echo "  📄 $filename"
    scp "$file" "${SERVER}:/tmp/"
    ssh ${SERVER} "docker cp /tmp/${filename} ${CONTAINER}:${HTML_DIR}/assets/"
done

# 上传pwa-icons目录
echo "🎨 上传 pwa-icons 目录..."
ssh ${SERVER} "docker exec ${CONTAINER} sh -c 'mkdir -p ${HTML_DIR}/pwa-icons'"

for file in dist/pwa-icons/*; do
    filename=$(basename "$file")
    echo "  🖼️ $filename"
    scp "$file" "${SERVER}:/tmp/"
    ssh ${SERVER} "docker cp /tmp/${filename} ${CONTAINER}:${HTML_DIR}/pwa-icons/"
done

# 重载Nginx
echo "♻️ 重载 Nginx..."
ssh ${SERVER} "docker exec ${CONTAINER} nginx -s reload"

echo "✅ 部署完成！"
