# 音乐数字版权保护平台

基于 Hyperledger Fabric 2.5 + Go + Gin 的音乐版权存证与授权系统。支持版权存证、授权发放与核验、版权转让、争议存证、PDF 存证证书下载。

## 技术栈

| 层 | 技术 |
|---|---|
| 区块链 | Hyperledger Fabric 2.5 + CouchDB |
| 链码 | Go + fabric-contract-api-go |
| 后端 | Go 1.22+ + Gin |
| 数据库 | MySQL 8 |
| 前端 | Vue 3 + Vite + Element Plus |
| 认证 | JWT (HS256) |

## 部署环境

本项目在 **WSL 2（Ubuntu 24.04）** 中运行，Fabric 网络跑在 Docker 容器内。已配置国内网络（Docker 镜像加速、Go 模块代理、npm 镜像），无需翻墙。

WSL 内已装好：Docker、Go 1.22、Node 18、MySQL 8、Fabric 工具链（`cryptogen` / `configtxgen` v2.5.16）。项目代码位于 Windows D 盘 `D:\blockchainlesson\Github\blockchain`，WSL 内通过 `/mnt/d/blockchainlesson/Github/blockchain` 访问。

## 启动步骤（在 WSL 内执行）

### 一键启动（推荐）

**终端 1 —— 后端（自动包含 Fabric 网络启动与链码部署）：**

```bash
bash /mnt/d/blockchainlesson/Github/blockchain/start_backend.sh
```

脚本自动完成：前置检查（docker / go / config.yaml）→ `bootstrap.sh`（生成证书 → 启动容器 → 创建通道）→ 修正证书目录权限 → `deploy.sh`（部署链码）→ 启动后端 `:8080`。中途需要输入一次 sudo 密码。

启动成功日志依次出现：`MySQL connected` → `Fabric Gateway connected` → `Server starting on :8080`。

**终端 2 —— 前端：**

```bash
bash /mnt/d/blockchainlesson/Github/blockchain/start_frontend.sh
```

自动检查 node/npm，`node_modules` 缺失时自动 `npm install`，然后启动 Vite 开发服务器。

浏览器访问 **http://localhost:5173**（vite 将 `/api` 代理到后端 8080）。

> Ctrl+C 停止前后端进程；Fabric 容器不受影响，彻底关闭见下方「停止 / 清理」。

### 手动分步执行（备查，与脚本等价）

#### 1. 进入 WSL，确认 Docker 已启动

```bash
wsl
sudo systemctl start docker    # 未自动启动时执行
```

#### 2. 启动 Fabric 网络

```bash
cd /mnt/d/blockchainlesson/Github/blockchain/network/scripts
sudo bash bootstrap.sh    # 生成证书 → 启动容器 → 创建通道
```

> **注意**：bootstrap 会重新生成 crypto-config（root 所有）。后端启动前需授权：

```bash
sudo chown -R xhy:xhy /mnt/d/blockchainlesson/Github/blockchain/network/crypto-config
```

#### 3. 部署链码

```bash
cd /mnt/d/blockchainlesson/Github/blockchain/network/scripts
sudo bash deploy.sh       # 打包 → 安装 → 审批 → 提交链码
```

#### 4. 启动后端 API

```bash
cd /mnt/d/blockchainlesson/Github/blockchain/backend
GOPROXY=https://goproxy.cn,direct GOSUMDB=off go run main.go
```

#### 5. 启动前端

```bash
cd /mnt/d/blockchainlesson/Github/blockchain/frontend
npm install     # 首次需要
npm run dev
```

## 停止 / 清理

```bash
cd /mnt/d/blockchainlesson/Github/blockchain/network/scripts
sudo bash teardown.sh    # 停止容器、删除链码镜像、清空证书与通道工件
```

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

## 常见问题与修复记录

- **链码容器崩溃 `too_many_pings`**：Fabric 2.5 + Go 1.20+ 的 gRPC keepalive 兼容性问题。已修改链码 vendored shim（`fabric-chaincode-go/.../internal/config.go`）为 `PermitWithoutStream=false`、`Time=24h` 修复。**勿重跑 `go mod vendor` 覆盖**。
- **后端查询超时 `DeadlineExceeded`**：`backend/service/fabric.go` 中 Gateway 四个超时已从 `0` 改为 `60 * time.Second`（`0` 会被当作立即超时）。
- **空响应 `unexpected end of JSON input`**：链码查询函数已改为返回显式空数组 `[]`（非 `nil`）。
- **`chaincode already installed`**：deploy.sh 的 install 步骤已加 `|| true` 幂等处理；全新网络 sequence 从 1 开始。
- **首次查询较慢**：链码容器空闲约 5 分钟后被回收，首次查询需冷启动约 30 秒，属正常现象。
- **国内网络**：Docker 镜像源配置在 `/etc/docker/daemon.json`；Go 用 `GOPROXY=https://goproxy.cn,direct GOSUMDB=off`；npm 用 `registry.npmmirror.com`。
