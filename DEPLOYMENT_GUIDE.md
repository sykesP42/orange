# 博主管理系统 - 完整部署指南

## 项目结构

```
/var/www/
├── blogger/              # 前端项目
│   ├── dist/            # 构建后的静态文件
│   ├── src/
│   ├── package.json
│   └── vite.config.js
│
└── blogger-server/       # 后端项目
    ├── server/
    │   └── index.js     # 后端入口文件
    ├── data/            # 数据存储目录
    ├── package.json
    └── node_modules/
```

## 1. 后端部署

### 1.1 进入后端项目目录
```bash
cd /var/www/blogger-server
```

### 1.2 安装依赖
```bash
npm install --production
```

### 1.3 创建数据目录
```bash
mkdir -p data
```

### 1.4 检查后端配置
编辑 `server/index.js`，确保监听地址正确：
```javascript
app.listen(3002, '0.0.0.0', () => {
  console.log('服务器运行在 http://0.0.0.0:3002')
})
```

### 1.5 使用PM2启动后端
```bash
# 启动后端服务
pm2 start server/index.js --name "blogger-server"

# 设置开机自启
pm2 startup
pm2 save

# 查看服务状态
pm2 status

# 查看日志
pm2 logs blogger-server
```

### 1.6 测试后端API
```bash
# 测试博主列表API
curl http://127.0.0.1:3002/api/blogger/list

# 测试登录API
curl -X POST http://127.0.0.1:3002/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
```

## 2. 前端部署

### 2.1 进入前端项目目录
```bash
cd /var/www/blogger
```

### 2.2 安装依赖
```bash
npm install
```

### 2.3 构建前端
```bash
npm run build
```

### 2.4 验证构建结果
```bash
ls -la dist/
```

应该看到 `index.html` 和其他静态文件。

## 3. Nginx配置

### 3.1 创建Nginx配置文件
```bash
nano /www/server/panel/vhost/nginx/test.conf
```

### 3.2 添加以下配置
```nginx
server {
    listen 80;
    server_name 123.56.99.154;  # 替换为您的域名或IP

    # 前端静态文件
    location / {
        root /var/www/blogger/dist;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    # 后端API代理
    location /api {
        proxy_pass http://127.0.0.1:3002;  # 重要：端口必须是3002
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # Gzip压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/javascript application/json;

    # 静态资源缓存
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

### 3.3 测试并重启Nginx
```bash
# 测试配置
nginx -t

# 重启Nginx
systemctl reload nginx

# 查看Nginx状态
systemctl status nginx
```

## 4. 防火墙配置

### 4.1 Ubuntu (UFW)
```bash
# 开放端口
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw allow 3002/tcp  # 后端API（可选）

# 启用防火墙
sudo ufw enable

# 查看状态
sudo ufw status
```

### 4.2 CentOS (firewalld)
```bash
# 开放端口
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --permanent --add-port=3002/tcp

# 重载防火墙
sudo firewall-cmd --reload

# 查看状态
sudo firewall-cmd --list-all
```

## 5. 云服务器安全组配置

登录您的云服务商控制台（阿里云/腾讯云等），确保以下端口已开放：
- **80** (HTTP) - 必须
- **443** (HTTPS) - 可选
- **22** (SSH) - 必须
- **3002** (后端API) - 可选，如果通过Nginx代理则不需要

## 6. 快速部署脚本

### 6.1 后端部署脚本
```bash
#!/bin/bash

echo "开始部署后端..."

cd /var/www/blogger-server

# 安装依赖
npm install --production

# 创建数据目录
mkdir -p data

# 启动服务
pm2 start server/index.js --name "blogger-server"
pm2 save

echo "后端部署完成！"
```

### 6.2 前端部署脚本
```bash
#!/bin/bash

echo "开始部署前端..."

cd /var/www/blogger

# 安装依赖
npm install

# 构建前端
npm run build

echo "前端部署完成！"
```

## 7. 常用管理命令

### 7.1 后端管理
```bash
# 查看后端状态
pm2 status

# 查看后端日志
pm2 logs blogger-server

# 重启后端
pm2 restart blogger-server

# 停止后端
pm2 stop blogger-server

# 删除后端
pm2 delete blogger-server
```

### 7.2 Nginx管理
```bash
# 重启Nginx
systemctl restart nginx

# 重新加载配置（不中断服务）
systemctl reload nginx

# 查看Nginx状态
systemctl status nginx

# 查看错误日志
tail -f /var/log/nginx/error.log

# 查看访问日志
tail -f /var/log/nginx/access.log
```

### 7.3 系统管理
```bash
# 查看端口监听
netstat -tlnp | grep 3002

# 查看进程
ps aux | grep node

# 查看磁盘空间
df -h

# 查看内存使用
free -h
```

## 8. 故障排查

### 8.1 登录失败，提示"请检查网络"

#### 检查后端服务
```bash
pm2 status
```

如果后端未运行：
```bash
cd /var/www/blogger-server
pm2 start server/index.js --name "blogger-server"
```

#### 检查后端监听端口
```bash
netstat -tlnp | grep 3002
```

应该显示 `0.0.0.0:3002` 或 `127.0.0.1:3002`

#### 测试后端API
```bash
curl http://127.0.0.1:3002/api/blogger/list
```

#### 检查Nginx配置
```bash
cat /www/server/panel/vhost/nginx/test.conf | grep proxy_pass
```

应该显示 `proxy_pass http://127.0.0.1:3002;`

#### 重启Nginx
```bash
nginx -t
systemctl reload nginx
```

### 8.2 前端页面无法访问

#### 检查前端文件
```bash
ls -la /var/www/blogger/dist
```

应该看到 `index.html` 文件。

#### 检查Nginx配置
```bash
cat /www/server/panel/vhost/nginx/test.conf | grep root
```

应该显示 `root /var/www/blogger/dist;`

#### 检查Nginx权限
```bash
ls -la /var/www/blogger/dist
```

确保Nginx用户有读取权限。

### 8.3 后端API返回500错误

#### 查看后端日志
```bash
pm2 logs blogger-server
```

#### 检查数据目录权限
```bash
ls -la /var/www/blogger-server/data
```

确保数据目录存在且有写权限。

## 9. 数据备份

### 9.1 创建备份脚本
```bash
nano /var/www/blogger-server/backup.sh
```

添加以下内容：
```bash
#!/bin/bash
BACKUP_DIR="/var/backups/blogger-manager"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

# 备份数据文件
tar -czf $BACKUP_DIR/data_$DATE.tar.gz -C /var/www/blogger-server data/

# 删除30天前的备份
find $BACKUP_DIR -name "data_*.tar.gz" -mtime +30 -delete

echo "Backup completed: data_$DATE.tar.gz"
```

### 9.2 设置定时备份
```bash
chmod +x /var/www/blogger-server/backup.sh

# 添加到crontab（每天凌晨2点备份）
crontab -e
# 添加以下行
0 2 * * * /var/www/blogger-server/backup.sh
```

## 10. 部署检查清单

- [ ] 后端项目已上传到 `/var/www/blogger-server`
- [ ] 前端项目已上传到 `/var/www/blogger`
- [ ] 后端依赖已安装
- [ ] 前端已构建
- [ ] 后端服务已通过PM2启动
- [ ] 后端监听3002端口
- [ ] 后端API可以正常访问
- [ ] Nginx配置已完成
- [ ] Nginx已重启
- [ ] 防火墙规则已配置
- [ ] 云服务器安全组已开放80端口
- [ ] 前端页面可以正常访问
- [ ] 登录功能正常
- [ ] 备份脚本已设置

## 11. 更新部署

### 11.1 更新后端
```bash
cd /var/www/blogger-server
git pull  # 或上传新文件
npm install --production
pm2 restart blogger-server
```

### 11.2 更新前端
```bash
cd /var/www/blogger
git pull  # 或上传新文件
npm install
npm run build
```

## 12. 性能优化

### 12.1 PM2集群模式
```bash
# 根据CPU核心数启动多个实例
pm2 start server/index.js --name "blogger-server" -i max
```

### 12.2 Nginx缓存配置
在Nginx配置中添加：
```nginx
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=my_cache:10m max_size=1g inactive=60m;

location /api {
    proxy_cache my_cache;
    proxy_cache_valid 200 5m;
    proxy_pass http://127.0.0.1:3002;
    # ... 其他配置
}
```

## 13. 安全建议

1. **定期更新系统**
   ```bash
   sudo apt update && sudo apt upgrade -y
   ```

2. **配置SSH密钥认证**
   ```bash
   # 禁用密码登录
   sudo nano /etc/ssh/sshd_config
   # 设置 PasswordAuthentication no
   sudo systemctl restart sshd
   ```

3. **配置防火墙**
   ```bash
   sudo ufw enable
   ```

4. **定期备份数据**
   ```bash
   # 确保备份脚本正常运行
   crontab -l
   ```

5. **监控服务状态**
   ```bash
   # 设置监控脚本
   pm2 install pm2-logrotate
   ```

## 14. 联系支持

如遇到问题，请提供以下信息：
1. PM2服务状态：`pm2 status`
2. 后端日志：`pm2 logs blogger-server --lines 50`
3. Nginx日志：`sudo tail -50 /var/log/nginx/error.log`
4. 网络端口检查：`netstat -tlnp | grep 3002`
5. Nginx配置：`cat /www/server/panel/vhost/nginx/test.conf | grep proxy_pass`
