#!/bin/bash
# ==============================================================================
# 🚀 日常启动（保留链上数据，直接复制粘贴执行即可，无需提前 cd）：
#
#     bash /mnt/d/blockchainlesson/Github/blockchain/start_backend.sh
#
# 🧨 首次初始化 / 彻底重置（清空全部链上数据，MySQL 用户保留）：
#
#     bash /mnt/d/blockchainlesson/Github/blockchain/start_backend.sh --init
#
# 日常启动只做 3 步：
#   1. cd network && docker compose up -d   （容器停了就拉起，数据卷原样保留，账本不丢）
#   2. 必要时 chown 证书目录                 （证书已归当前用户则跳过，全程无需输密码）
#   3. cd backend && go run main.go
#
# 只有 --init 才会执行 bootstrap.sh（清理 + 生成证书 + 建通道）和 deploy.sh（部署链码），
# 这两步会删掉全部账本数据卷，执行前有二次确认。
# ==============================================================================
set -e

# —— 定位仓库根：无论你从哪个目录执行 bash start_backend.sh，都能找到 bootstrap.sh / deploy.sh / backend ——
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$REPO_ROOT/network/scripts"
BACKEND_DIR="$REPO_ROOT/backend"

# —— 当前用户（即便是 sudo 执行本脚本，SUDO_USER 也会是真正登录用户） ——
#    用法：sudo bash start_backend.sh   →  ORIGINAL_USER=登录用户
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

# —— 模式解析：run（默认，保留数据） / --init（重置环境） ——
MODE="${1:-run}"
case "$MODE" in
    ""|run)
        MODE="run" ;;
    init|--init)
        MODE="init" ;;
    *)
        log_error "未知参数: $MODE （用法：bash start_backend.sh [--init]）"
        exit 1 ;;
esac

# ==============================================================================
# 0. 前置检查
# ==============================================================================
log_info "[0] 前置检查：docker / docker compose / go"

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

if [ "$MODE" = "init" ]; then
    echo ""
    log_warn "⚠⚠⚠  你正在执行 --init（重置环境）：bootstrap.sh 会 docker compose down -v 删除全部数据卷，"
    log_warn "        链上的存证 / 授权 / 争议记录将全部丢失（MySQL 里的用户账号保留）。"
    log_warn "        如果只是想重启服务，请直接运行: bash start_backend.sh（不带参数）"
    echo ""
    read -r -p "确认要重置 Fabric 网络吗？输入 y 继续，其他任意键取消: " CONFIRM
    if [ "$CONFIRM" != "y" ]; then
        log_error "已取消。"
        exit 1
    fi

    log_info "需要 sudo 权限运行 bootstrap.sh / deploy.sh / chown，请输入密码（若已免密则直接跳过）："
    sudo -v || { log_error "sudo 认证失败"; exit 1; }

    # ==============================================================================
    # init-1. bootstrap.sh（清理 + 生成证书 + 启动容器 + 创建通道）
    # ==============================================================================
    log_info "[init 1/3] 运行 bootstrap.sh：清理旧数据 → 生成证书 → 启动容器 → 创建通道..."
    cd "$SCRIPTS_DIR"
    sudo bash "$SCRIPTS_DIR/bootstrap.sh"

    # ==============================================================================
    # init-2. chown 证书目录（root 生成证书 → 给普通用户读，否则 go run main.go 会报 cert 文件权限错）
    # ==============================================================================
    log_info "[init 2/3] 修正证书目录权限：network/crypto-config → chown $ORIGINAL_USER:$ORIGINAL_USER ..."
    cd "$REPO_ROOT"
    sudo chown -R "${ORIGINAL_USER}:${ORIGINAL_USER}" network/crypto-config network/channel-artifacts 2>/dev/null || true

    # ==============================================================================
    # init-3. deploy.sh 部署链码（打包 → 两 org 安装 → 审批 → 提交）
    # ==============================================================================
    log_info "[init 3/3] 运行 deploy.sh：打包并部署版权链码 copyright_1.0 到 mychannel..."
    cd "$SCRIPTS_DIR"
    sudo bash "$SCRIPTS_DIR/deploy.sh"

else
    # ==============================================================================
    # run-1. 证书目录若归 root（上次 init 后没 chown 过）才需要 sudo 修一次，否则全程免密码
    # ==============================================================================
    if [ ! -d "$REPO_ROOT/network/crypto-config" ]; then
        log_error "network/crypto-config 不存在，Fabric 网络尚未初始化。"
        log_warn  "  首次使用请执行：bash $REPO_ROOT/start_backend.sh --init"
        exit 1
    fi

    if [ ! -w "$REPO_ROOT/network/crypto-config" ]; then
        log_info "证书目录当前用户不可写（可能归 root），需要 sudo chown 一次，请输入密码："
        sudo -v || { log_error "sudo 认证失败"; exit 1; }
        sudo chown -R "${ORIGINAL_USER}:${ORIGINAL_USER}" "$REPO_ROOT/network/crypto-config" "$REPO_ROOT/network/channel-artifacts" 2>/dev/null || true
    fi

    # ==============================================================================
    # run-2. 拉起网络容器（不动数据卷，账本 / 已部署链码原样保留）
    # ==============================================================================
    log_info "[run 1/2] 启动 Fabric 网络容器：cd network && docker compose up -d（保留原有数据卷，账本数据不丢）..."
    cd "$REPO_ROOT/network"
    docker compose up -d
    log_info "等待容器就绪 (8s)..."
    sleep 8

    # ==============================================================================
    # run-3. 启动后端服务
    # ==============================================================================
    log_info "[run 2/2] 启动后端服务：cd backend && GOPROXY=https://goproxy.cn,direct GOSUMDB=off go run main.go"
fi

log_info "  Ctrl+C 停止后端；Fabric 网络容器会继续在后台运行，下次仍可直接 bash start_backend.sh 拉起后端。"
log_info "  如需彻底关闭并清空 Fabric 网络：cd network/scripts && sudo bash teardown.sh（会删除全部链上数据，不可逆）"

cd "$BACKEND_DIR"

# —— 以「原始登录用户」身份执行 go run（不能是 sudo root，否则家目录/Go 缓存全在 /root 下） ——
if [ "$EUID" = "0" ]; then
    # 如果用户是用 sudo bash start_backend.sh 运行的，这里 su 回到普通用户
    exec sudo -u "$ORIGINAL_USER" -H env GOPROXY="https://goproxy.cn,direct" GOSUMDB=off go run main.go
else
    exec env GOPROXY="https://goproxy.cn,direct" GOSUMDB=off go run main.go
fi
