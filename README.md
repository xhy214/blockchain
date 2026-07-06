# 音乐数字版权保护平台

基于 Hyperledger Fabric 2.5 + Go + Gin 的区块链版权存证与授权系统。

## 技术栈

| 层 | 技术 |
|---|---|
| 区块链 | Hyperledger Fabric 2.5 + CouchDB |
| 链码 | Go + fabric-contract-api-go |
| 后端 | Go 1.22+ + Gin |
| 数据库 | MySQL 8 |
| 认证 | JWT (HS256) |

## 环境要求

- Docker + Docker Compose
- Go 1.22+
- MySQL 8
- Fabric 二进制工具：`cryptogen`、`configtxgen`（[安装指引](https://hyperledger-fabric.readthedocs.io/en/release-2.5/prereqs.html)）

## 启动步骤

### 1. 启动 Fabric 网络

```bash
cd network/scripts
bash bootstrap.sh   # 生成证书 → 启动容器 → 创建通道
```

### 2. 部署链码

```bash
bash deploy.sh      # 打包 → 安装 → 审批 → 提交链码
```

### 3. 初始化 MySQL

```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS copyright_db DEFAULT CHARSET utf8mb4;"
# 后端启动时会自动建表
```

修改 `backend/config.yaml` 中的 MySQL DSN 以匹配你的数据库密码。

### 4. 启动后端 API

```bash
cd backend
go mod tidy
go run main.go
```

服务监听 `http://localhost:8080`。

## API 接口一览

### 用户模块
| Method | Path | 说明 |
|---|---|---|
| POST | `/api/v1/user/register` | 注册 |
| POST | `/api/v1/user/login` | 登录，返回 JWT |
| GET | `/api/v1/user/profile` | 当前用户信息 |

### 版权模块
| Method | Path | 说明 |
|---|---|---|
| POST | `/api/v1/copyright/register` | 版权存证（上传音频文件） |
| GET | `/api/v1/copyright/:workID` | 作品详情 |
| GET | `/api/v1/copyright/my/list` | 我的作品列表 |
| GET | `/api/v1/copyright/search` | 搜索 `?keyword=&page=&size=` |
| GET | `/api/v1/copyright/:workID/history` | 链上历史 |
| POST | `/api/v1/copyright/verify-hash` | 文件哈希验真 |
| POST | `/api/v1/copyright/transfer` | 版权转让 |

### 授权模块
| Method | Path | 说明 |
|---|---|---|
| POST | `/api/v1/license/grant` | 发放授权 |
| GET | `/api/v1/license/verify` | 核验授权 `?workID=&licenseeID=` |
| GET | `/api/v1/license/my` | 我的授权列表 |
| POST | `/api/v1/license/revoke` | 撤销授权 |
| POST | `/api/v1/license/record-usage` | 记录使用 |

### 其他
| Method | Path | 说明 |
|---|---|---|
| GET | `/api/v1/copyright/:workID/certificate` | 下载存证证书 PDF |
| POST | `/api/v1/dispute/file` | 提交版权争议 |
| GET | `/api/v1/dispute/:workID` | 查询争议 |

## 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": { }
}
```

### 错误码

| code | 含义 |
|---|---|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 未登录 / Token 无效 |
| 1003 | 权限不足 |
| 2001 | 作品不存在 |
| 2003 | 文件哈希不匹配 |
| 3001 | 授权不存在 |
| 3002 | 授权已过期或无效 |
| 5001 | Fabric 网络错误 |

## 停止 / 清理

```bash
cd network/scripts
bash teardown.sh
```
