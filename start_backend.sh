#!/bin/bash
# ==============================================================================
# 🚀 一行启动（直接复制粘贴执行即可，无需提前 cd）：
#
#     bash /home/iloveruci/dyl/blockchain/start_backend.sh
#
# 功能：把手动执行的 5 步打包成一步：
#   1. sudo bash network/scripts/bootstrap.sh   (清理 + 生成证书 + 启动容器 + 创建通道)
#   2. sudo chown -R 当前用户:当前用户 network/crypto-config  （解决 root 生成后 go 读不到证书的问题）
#   3. sudo bash network/scripts/deploy.sh      （打包 / 安装 / 审批 / 提交链码）
#   4. cd backend && GOPROXY=... GOSUMDB=off go run main.go   （启动后端服务）
# ==============================================================================
set -e

# —— 定位仓库根：无论你从哪个目录执行 bash start_backend.sh，都能找到 bootstrap.sh / deploy.sh / backend ——
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$REPO_ROOT/network/scripts"
BACKEND_DIR="$REPO_ROOT/backend"

# —— 当前用户（即便是 sudo 执行本脚本，SUDO_USER 也会是真正登录用户） ——
#    用法：sudo bash start_backend.sh   →  ORIGINAL_USER=iloveruci
#          bash start_backend.sh       →  ORIGINAL_USER=$(whoami)
if [ -n "$SUDO_USER" ]; then
    ORIGINAL_USER="$SUDO_USER"
else
    ORIGINAL_USER="$(whoami)"
fi

# —— 颜色输出，可读性好 ——
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'
log_info()  { echo -e "${GREEN}▸ $*${NC}"; }
log_warn()  { echo -e "${YELLOW}▸ $*${NC}"; }
log_error() { echo -e "${RED}✗ $*${NC}" 1>&2; }

cd "$REPO_ROOT"

# ==============================================================================
# 0. 前置检查
# ==============================================================================
log_info "[0/4] 前置检查：docker / docker compose / go / sudo 权限"

if ! command -v docker >/dev/null 2>&1; then
    log_error "未找到 docker。请先安装 Docker Desktop (WSL) 或 docker.io，并把当前用户加入 docker 组：sudo usermod -aG docker $ORIGINAL_USER"
    exit 1
fi

if ! command -v go >/dev/null 2>&1; then
    log_error "未找到 go。请先安装 Go 1.21+：https://go.dev/doc/install"
    exit 1
fi

# 检查 backend/config.yaml 是否存在（否则后端启动会 1045 数据库密码错）
if [ ! -f "$BACKEND_DIR/config.yaml" ]; then
    log_error "后端配置文件不存在：$BACKEND_DIR/config.yaml"
    log_warn  "  请先复制模板再修改密码："
    log_warn  "     cp $BACKEND_DIR/config.yaml.example $BACKEND_DIR/config.yaml"
    log_warn  "     vim $BACKEND_DIR/config.yaml   # 把 mysql.dsn 里的 '你的MySQL密码' 改成真实密码"
    log_warn  "  config.yaml 已在 .gitignore 中，不用担心密码被提交。"
    exit 1
fi

# 询问 sudo（保证后面 bootstrap/deploy/chown 都有 sudo 权限，不需要中途再输一次）
log_info "需要 sudo 权限运行 bootstrap.sh / deploy.sh / chown，请输入密码（若已免密则直接跳过）："
sudo -v || { log_error "sudo 认证失败"; exit 1; }

# ==============================================================================
# 1. bootstrap.sh
# ==============================================================================
log_info "[1/4] 运行 bootstrap.sh：清理旧数据 → 生成证书 → 启动容器 → 创建通道..."
cd "$SCRIPTS_DIR"
sudo bash "$SCRIPTS_DIR/bootstrap.sh"

# ==============================================================================
# 2. chown 证书目录（root 生成证书 → 给普通用户读，否则 go run main.go 会报 cert 文件权限错）
# ==============================================================================
log_info "[2/4] 修正证书目录权限：network/crypto-config → chown $ORIGINAL_USER:$ORIGINAL_USER ..."
cd "$REPO_ROOT"
sudo chown -R "${ORIGINAL_USER}:${ORIGINAL_USER}" network/crypto-config network/channel-artifacts 2>/dev/null || true

# ==============================================================================
# 3. deploy.sh 部署链码（打包 → 两 org 安装 → 审批 → 提交）
# ==============================================================================
log_info "[3/4] 运行 deploy.sh：打包并部署版权链码 copyright_1.0 到 mychannel..."
cd "$SCRIPTS_DIR"
sudo bash "$SCRIPTS_DIR/deploy.sh"

# ==============================================================================
# 4. 启动后端服务（go run main.go，用户身份运行，非 root）
# ==============================================================================
log_info "[4/4] 启动后端服务：cd backend && GOPROXY=https://goproxy.cn,direct GOSUMDB=off go run main.go"
log_info "  Ctrl+C 停止后端；链码容器 / Fabric 网络不会自动停止。"
log_info "  如需彻底关闭 Fabric 网络：cd network/scripts && sudo bash teardown.sh"

cd "$BACKEND_DIR"

# —— 以「原始登录用户」身份执行 go run（不能是 sudo root，否则家目录/Go 缓存全在 /root 下） ——
if [ "$EUID" = "0" ]; then
    # 如果用户是用 sudo bash start_backend.sh 运行的，这里 su 回到普通用户
    exec sudo -u "$ORIGINAL_USER" -H env GOPROXY="https://goproxy.cn,direct" GOSUMDB=off go run main.go
else
    exec env GOPROXY="https://goproxy.cn,direct" GOSUMDB=off go run main.go
fi
