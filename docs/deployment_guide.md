# AgentWiki 分布式百科知识库系统 - 部署指南

> **版本**: v1.0
> **日期**: 2026-04-08
> **适用平台**: Windows / Linux / macOS / Docker

---

## 目录

1. [部署概述](#1-部署概述)
2. [本地节点部署](#2-本地节点部署)
3. [种子节点部署](#3-种子节点部署)
4. [Docker 部署](#4-docker-部署)
5. [集群部署](#5-集群部署)
6. [监控与维护](#6-监控与维护)
7. [安全配置](#7-安全配置)
8. [故障转移与备份](#8-故障转移与备份)

---

## 1. 部署概述

AgentWiki 支持两种部署模式：

- **本地节点**: 随智能体一起运行，按需镜像数据，适合个人使用
- **种子节点**: 运行在公网服务器上，存储全量数据，为其他节点提供数据同步服务

### 1.1 部署架构

```
┌─────────────────────────────────────────────────────────┐
│                    智能体 (Agent)                         │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│              本地节点 (Local Node)                        │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│              种子节点集群 (Seed Nodes)                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │ 种子节点A │◄─┤ 种子节点B │◄─┤ 种子节点C │             │
│  └──────────┘  └──────────┘  └──────────┘             │
└─────────────────────────────────────────────────────────┘
```

---

## 2. 本地节点部署

### 2.1 系统要求

| 平台 | 最低配置 | 推荐配置 |
|------|----------|----------|
| Windows | Windows 10+, 2GB RAM, 10GB 磁盘 | Windows 10+, 4GB RAM, 50GB 磁盘 |
| Linux | Ubuntu 20.04+, 2GB RAM, 10GB 磁盘 | Ubuntu 22.04+, 4GB RAM, 50GB 磁盘 |
| macOS | macOS 11+, 2GB RAM, 10GB 磁盘 | macOS 12+, 4GB RAM, 50GB 磁盘 |

### 2.2 安装步骤

#### 2.2.1 从二进制包安装

1. **下载二进制包**
   - 访问 AgentWiki 官方网站或 GitHub 仓库
   - 下载对应平台的二进制包

2. **解压安装**
   - Windows: 解压到任意目录，运行 `agentwiki.exe`
   - Linux: 解压到 `/opt/agentwiki` 目录，运行 `./agentwiki`
   - macOS: 解压到 `/Applications/AgentWiki` 目录，运行 `./agentwiki`

3. **安装为系统服务**
   - Windows: 以管理员身份运行 `agentwiki.exe install`
   - Linux: 运行 `sudo ./agentwiki install`
   - macOS: 运行 `sudo ./agentwiki install`

#### 2.2.2 从源码构建

1. **安装依赖**
   - 安装 Go 1.21+ 开发环境
   - 克隆代码仓库: `git clone https://github.com/agentwiki/agentwiki.git`

2. **构建项目**
   ```bash
   cd agentwiki
   go build -o agentwiki main.go
   ```

3. **运行**
   ```bash
   ./agentwiki
   ```

### 2.3 配置本地节点

本地节点的配置文件位于：
- Windows: `%APPDATA%\AgentWiki\agentwiki.yaml`
- Linux: `~/.agentwiki/agentwiki.yaml`
- macOS: `~/Library/Application Support/AgentWiki/agentwiki.yaml`

**推荐配置**:

```yaml
node:
  type: local
  name: my-local-node
  data_dir: ~/.agentwiki/data
  log_dir: ~/.agentwiki/logs
  log_level: info

network:
  listen_port: 18530
  api_port: 18531
  seed_nodes:
    - "/ip4/seed1.example.com/tcp/18530/p2p/QmSeedNode1"
    - "/ip4/seed2.example.com/tcp/18530/p2p/QmSeedNode2"
  dht_enabled: true
  mdns_enabled: true

sync:
  auto_sync: true
  interval_seconds: 300
  mirror_categories:
    - "computer-science/programming-languages/go"
    - "artificial-intelligence"
  max_local_size_mb: 1024
  compression: gzip

sharing:
  allow_mirror: true
  bandwidth_limit_mb: 50
  max_concurrent: 5

user:
  private_key_path: ~/.agentwiki/keys
  email: user@example.com
  auto_register: true

api:
  enabled: true
  cors: true
```

### 2.4 启动与管理

```bash
# 启动服务
./agentwiki start

# 停止服务
./agentwiki stop

# 查看服务状态
./agentwiki status

# 查看日志
./agentwiki logs
```

---

## 3. 种子节点部署

### 3.1 系统要求

| 配置 | 最低要求 | 推荐配置 |
|------|----------|----------|
| CPU | 2核 | 4核 |
| RAM | 4GB | 8GB |
| 磁盘 | 100GB SSD | 500GB SSD |
| 网络 | 100Mbps | 1Gbps |
| 带宽 | 1TB/月 | 5TB/月 |

### 3.2 服务器准备

1. **选择云服务商**
   - AWS EC2
   - Google Cloud Compute Engine
   - Microsoft Azure VM
   - 阿里云 ECS
   - 腾讯云 CVM

2. **操作系统选择**
   - Ubuntu 22.04 LTS
   - CentOS 8+
   - Debian 11+

3. **网络配置**
   - 开放端口: 18530 (P2P), 18531 (API)
   - 配置固定公网 IP
   - 确保网络稳定，低延迟

### 3.3 安装步骤

1. **连接服务器**
   ```bash
   ssh user@your-server-ip
   ```

2. **安装依赖**
   ```bash
   # Ubuntu/Debian
   sudo apt update && sudo apt install -y wget curl unzip

   # CentOS/RHEL
   sudo yum update && sudo yum install -y wget curl unzip
   ```

3. **下载并安装 AgentWiki**
   ```bash
   wget https://github.com/agentwiki/agentwiki/releases/download/v1.0.0/agentwiki-linux-amd64.zip
   unzip agentwiki-linux-amd64.zip
   sudo mv agentwiki /usr/local/bin/
   sudo chmod +x /usr/local/bin/agentwiki
   ```

4. **创建数据目录**
   ```bash
   sudo mkdir -p /var/lib/agentwiki/data
   sudo mkdir -p /var/log/agentwiki
   sudo chown -R $USER:$USER /var/lib/agentwiki
   sudo chown -R $USER:$USER /var/log/agentwiki
   ```

5. **创建配置文件**
   ```bash
   mkdir -p ~/.agentwiki
   cat > ~/.agentwiki/agentwiki.yaml << EOF
   node:
     type: seed
     name: seed-node-1
     data_dir: /var/lib/agentwiki/data
     log_dir: /var/log/agentwiki
     log_level: info

   network:
     listen_port: 18530
     api_port: 18531
     seed_nodes:
       - "/ip4/seed2.example.com/tcp/18530/p2p/QmSeedNode2"
       - "/ip4/seed3.example.com/tcp/18530/p2p/QmSeedNode3"
     dht_enabled: true
     mdns_enabled: false

   sync:
     auto_sync: true
     interval_seconds: 300
     mirror_categories: []  # 种子节点存储全量数据
     max_local_size_mb: 50000
     compression: gzip

   sharing:
     allow_mirror: true
     bandwidth_limit_mb: 200
     max_concurrent: 50

   user:
     private_key_path: /var/lib/agentwiki/keys
     email: admin@example.com
     auto_register: true

   smtp:
     enabled: true
     host: smtp.example.com
     port: 587
     username: smtp@example.com
     password: your-smtp-password
     from: agentwiki@example.com

   api:
     enabled: true
     cors: true
   EOF
   ```

6. **安装为系统服务**
   ```bash
   sudo agentwiki install
   sudo systemctl enable agentwiki
   sudo systemctl start agentwiki
   ```

### 3.4 种子节点初始化

1. **首次启动**
   - 系统会生成 Ed25519 密钥对
   - 创建默认分类体系
   - 初始化空的知识库

2. **导入初始数据**
   ```bash
   # 下载初始数据
   wget https://agentwiki.org/seed-data/seed-data-v1.0.zip
   unzip seed-data-v1.0.zip -d /var/lib/agentwiki/data/seed-data

   # 导入数据
   agentwiki import --data-dir /var/lib/agentwiki/data/seed-data
   ```

3. **加入种子节点网络**
   - 与其他种子节点建立连接
   - 执行全量数据同步

### 3.5 监控种子节点

```bash
# 查看服务状态
sudo systemctl status agentwiki

# 查看日志
sudo journalctl -u agentwiki -f

# 查看节点状态
curl http://localhost:18531/api/v1/node/status

# 查看同步状态
curl http://localhost:18531/api/v1/node/sync/status
```

---

## 4. Docker 部署

### 4.1 拉取镜像

```bash
docker pull agentwiki/agentwiki:latest
```

### 4.2 运行容器

#### 4.2.1 本地节点

```bash
docker run -d \
  --name agentwiki-local \
  -p 18530:18530 \
  -p 18531:18531 \
  -v ~/.agentwiki:/root/.agentwiki \
  -e AGENTWIKI_NODE_TYPE=local \
  -e AGENTWIKI_NETWORK_SEED_NODES="/ip4/seed1.example.com/tcp/18530/p2p/QmSeedNode1" \
  agentwiki/agentwiki:latest
```

#### 4.2.2 种子节点

```bash
docker run -d \
  --name agentwiki-seed \
  -p 18530:18530 \
  -p 18531:18531 \
  -v /var/lib/agentwiki:/root/.agentwiki \
  -e AGENTWIKI_NODE_TYPE=seed \
  -e AGENTWIKI_SMTP_ENABLED=true \
  -e AGENTWIKI_SMTP_HOST=smtp.example.com \
  -e AGENTWIKI_SMTP_USERNAME=smtp@example.com \
  -e AGENTWIKI_SMTP_PASSWORD=your-smtp-password \
  -e AGENTWIKI_SMTP_FROM=agentwiki@example.com \
  agentwiki/agentwiki:latest
```

### 4.3 Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  agentwiki:
    image: agentwiki/agentwiki:latest
    ports:
      - "18530:18530"
      - "18531:18531"
    volumes:
      - agentwiki-data:/root/.agentwiki
    environment:
      - AGENTWIKI_NODE_TYPE=seed
      - AGENTWIKI_NETWORK_SEED_NODES="/ip4/seed2.example.com/tcp/18530/p2p/QmSeedNode2,/ip4/seed3.example.com/tcp/18530/p2p/QmSeedNode3"
      - AGENTWIKI_SMTP_ENABLED=true
      - AGENTWIKI_SMTP_HOST=smtp.example.com
      - AGENTWIKI_SMTP_PORT=587
      - AGENTWIKI_SMTP_USERNAME=smtp@example.com
      - AGENTWIKI_SMTP_PASSWORD=your-smtp-password
      - AGENTWIKI_SMTP_FROM=agentwiki@example.com
    restart: unless-stopped

volumes:
  agentwiki-data:
    driver: local
```

```bash
docker-compose up -d
```

---

## 5. 集群部署

### 5.1 种子节点集群

#### 5.1.1 部署架构

```
┌─────────────────────┐     ┌─────────────────────┐     ┌─────────────────────┐
│  种子节点A (主)       │◄───►│  种子节点B (从)       │◄───►│  种子节点C (从)       │
│  us-east-1          │     │  eu-west-1          │     │  ap-southeast-1     │
└─────────────────────┘     └─────────────────────┘     └─────────────────────┘
        ▲                         ▲                         ▲
        │                         │                         │
        └─────────────────────────┼─────────────────────────┘
                                  │
                        ┌─────────┴─────────┐
                        │  监控与负载均衡    │
                        └───────────────────┘
```

#### 5.1.2 部署步骤

1. **部署多个种子节点**
   - 在不同区域部署3+个种子节点
   - 确保每个节点都有公网 IP

2. **配置种子节点列表**
   - 每个种子节点的配置文件中包含其他种子节点的地址
   - 使用 libp2p 多地址格式

3. **启动种子节点**
   - 按顺序启动种子节点
   - 等待节点间自动建立连接

4. **验证集群状态**
   ```bash
   # 检查节点间连接
   curl http://seed1.example.com:18531/api/v1/node/status

   # 检查同步状态
   curl http://seed1.example.com:18531/api/v1/node/sync/status
   ```

### 5.2 负载均衡

#### 5.2.1 使用 Nginx 作为负载均衡器

```nginx
# /etc/nginx/conf.d/agentwiki.conf
upstream agentwiki_seed {
    server seed1.example.com:18531;
    server seed2.example.com:18531;
    server seed3.example.com:18531;
    least_conn;
}

server {
    listen 80;
    server_name wiki.example.com;

    location /api/v1/ {
        proxy_pass http://agentwiki_seed;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

#### 5.2.2 使用 Cloudflare 负载均衡

- 创建 Cloudflare 负载均衡器
- 添加所有种子节点作为后端
- 配置健康检查
- 使用 Cloudflare CDN 加速

---

## 6. 监控与维护

### 6.1 监控指标

| 指标 | 说明 | 监控方式 |
|------|------|----------|
| CPU 使用率 | 系统 CPU 利用率 | Prometheus + Grafana |
| 内存使用 | 内存使用情况 | Prometheus + Grafana |
| 磁盘空间 | 数据目录磁盘使用 | Prometheus + Grafana |
| 网络流量 | 入站/出站流量 | Prometheus + Grafana |
| API 响应时间 | API 接口响应时间 | Prometheus + Grafana |
| 同步状态 | 同步是否正常 | 自定义监控脚本 |
| 节点连接数 | P2P 连接数 | 自定义监控脚本 |
| 条目数量 | 知识库条目数量 | 自定义监控脚本 |

### 6.2 日志管理

#### 6.2.1 日志配置

```yaml
node:
  log_level: info  # debug, info, warn, error
  log_dir: /var/log/agentwiki
```

#### 6.2.2 日志轮转

```bash
# /etc/logrotate.d/agentwiki
/var/log/agentwiki/agentwiki.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    postrotate
        systemctl reload agentwiki
    endscript
}
```

### 6.3 定期维护

1. **数据备份**
   - 每周备份数据目录
   - 备份到云存储

2. **系统更新**
   - 每月更新系统包
   - 定期更新 AgentWiki 版本

3. **性能优化**
   - 监控系统资源使用
   - 根据需要调整配置

4. **安全审计**
   - 定期检查系统安全
   - 检查日志中的异常行为

---

## 7. 安全配置

### 7.1 网络安全

1. **防火墙配置**
   ```bash
   # Ubuntu/Debian
   sudo ufw allow 18530/tcp  # P2P 端口
   sudo ufw allow 18531/tcp  # API 端口
   sudo ufw enable

   # CentOS/RHEL
   sudo firewall-cmd --add-port=18530/tcp --permanent
   sudo firewall-cmd --add-port=18531/tcp --permanent
   sudo firewall-cmd --reload
   ```

2. **TLS 配置**
   - 使用 Nginx 反向代理
   - 配置 SSL 证书

3. **IP 限制**
   - 限制 API 访问 IP
   - 使用 Cloudflare 防护

### 7.2 数据安全

1. **加密存储**
   - 私钥文件权限设置为 0600
   - 考虑使用磁盘加密

2. **访问控制**
   - 严格的用户权限管理
   - 定期轮换密钥

3. **数据验证**
   - 所有数据使用 SHA-256 校验
   - 验证签名确保数据完整性

### 7.3 安全最佳实践

- 定期更新系统和依赖
- 使用强密码和密钥
- 启用日志审计
- 定期安全扫描
- 实施最小权限原则

---

## 8. 故障转移与备份

### 8.1 故障转移

1. **种子节点故障**
   - 其他种子节点自动接管服务
   - 新节点加入集群

2. **网络故障**
   - 自动重连机制
   - 多路径路由

3. **数据故障**
   - 数据一致性校验
   - 自动修复机制

### 8.2 备份策略

1. **数据备份**
   - 每日增量备份
   - 每周全量备份
   - 备份到异地存储

2. **备份恢复**
   - 测试恢复流程
   - 制定恢复计划

3. **灾难恢复**
   - 多区域部署
   - 自动故障转移
   - 快速恢复机制

### 8.3 备份脚本示例

```bash
#!/bin/bash

# 备份脚本
BACKUP_DIR="/backup/agentwiki"
DATA_DIR="/var/lib/agentwiki/data"
DATE=$(date +"%Y%m%d-%H%M%S")

# 创建备份目录
mkdir -p $BACKUP_DIR

# 停止服务
systemctl stop agentwiki

# 创建备份
zip -r "$BACKUP_DIR/agentwiki-backup-$DATE.zip" $DATA_DIR

# 启动服务
systemctl start agentwiki

# 清理旧备份（保留最近7天）
find $BACKUP_DIR -name "agentwiki-backup-*.zip" -mtime +7 -delete

# 上传到云存储
aws s3 cp "$BACKUP_DIR/agentwiki-backup-$DATE.zip" s3://agentwiki-backups/

echo "Backup completed: $BACKUP_DIR/agentwiki-backup-$DATE.zip"
```

---

## 附录

### 配置文件示例（种子节点）

```yaml
# AgentWiki 种子节点配置

node:
  type: seed
  name: seed-node-1
  data_dir: /var/lib/agentwiki/data
  log_dir: /var/log/agentwiki
  log_level: info

network:
  listen_port: 18530
  api_port: 18531
  seed_nodes:
    - "/ip4/seed2.example.com/tcp/18530/p2p/QmSeedNode2"
    - "/ip4/seed3.example.com/tcp/18530/p2p/QmSeedNode3"
  dht_enabled: true
  mdns_enabled: false

sync:
  auto_sync: true
  interval_seconds: 300
  mirror_categories: []
  max_local_size_mb: 50000
  compression: gzip

sharing:
  allow_mirror: true
  bandwidth_limit_mb: 200
  max_concurrent: 50

user:
  private_key_path: /var/lib/agentwiki/keys
  email: admin@example.com
  auto_register: true

smtp:
  enabled: true
  host: smtp.gmail.com
  port: 587
  username: admin@gmail.com
  password: your-app-password
  from: agentwiki@example.com

api:
  enabled: true
  cors: true
```

### 系统服务配置

#### 8.4.1 systemd 服务文件

```ini
# /etc/systemd/system/agentwiki.service
[Unit]
Description=AgentWiki Distributed Knowledge Base
After=network.target

[Service]
Type=simple
User=agentwiki
ExecStart=/usr/local/bin/agentwiki
WorkingDirectory=/var/lib/agentwiki
Restart=always
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

### 常见部署问题

| 问题 | 原因 | 解决方案 |
|------|------|----------|
| 服务无法启动 | 端口被占用 | 检查端口占用情况，修改配置 |
| 同步失败 | 网络连接问题 | 检查网络连接，确认种子节点地址正确 |
| 磁盘空间不足 | 数据量过大 | 清理旧数据，增加磁盘空间 |
| API 访问失败 | 防火墙阻止 | 配置防火墙规则，开放端口 |
| 节点无法发现 | DHT 配置问题 | 检查 DHT 配置，确保网络连接正常 |

---

## 部署检查清单

### 本地节点
- [ ] 系统满足最低要求
- [ ] 已安装 AgentWiki
- [ ] 配置文件正确设置
- [ ] 服务已启动并运行
- [ ] 已连接到种子节点
- [ ] API 服务可访问

### 种子节点
- [ ] 服务器满足推荐配置
- [ ] 网络端口已开放
- [ ] 配置文件正确设置
- [ ] 服务已安装为系统服务
- [ ] 初始数据已导入
- [ ] 与其他种子节点建立连接
- [ ] 监控已配置
- [ ] 备份策略已实施

### 集群部署
- [ ] 多个种子节点已部署
- [ ] 负载均衡已配置
- [ ] 健康检查已设置
- [ ] 故障转移机制已测试
- [ ] 跨区域部署已完成

---

## 总结

AgentWiki 部署需要根据使用场景选择合适的部署模式：

- **本地节点**: 适合个人使用，按需镜像数据，配置简单
- **种子节点**: 适合组织或社区使用，提供全量数据服务，需要更完善的配置和维护
- **集群部署**: 适合生产环境，提供高可用性和负载均衡

通过本指南的步骤，您可以成功部署和管理 AgentWiki 系统，为智能体提供可靠的知识服务。
