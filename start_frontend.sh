#!/bin/bash
# ==============================================================================
# 🚀 一行启动前端（直接复制粘贴执行即可，无需提前 cd）：
#
#     bash /home/iloveruci/dyl/blockchain/start_frontend.sh
#
# 功能：自动定位 frontend 目录 → 检查 node_modules → npm run dev 启动 Vite 开发服务器
# ==============================================================================
set -e

# —— 定位仓库根 ——
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_DIR="$REPO_ROOT/frontend"

# —— 颜色 ——
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'
log_info()  { echo -e "${GREEN}▸ $*${NC}"; }
log_warn()  { echo -e "${YELLOW}▸ $*${NC}"; }
log_error() { echo -e "${RED}✗ $*${NC}" 1>&2; }

# ==============================================================================
# 0. 前置检查
# ==============================================================================
log_info "[0/2] 前置检查：node / npm"

if ! command -v node >/dev/null 2>&1; then
    log_error "未找到 node。请先安装 Node.js 18+："
    log_warn  "  Ubuntu:  curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash - && sudo apt install -y nodejs"
    log_warn  "  WSL:    同上，或 Windows 安装 Node.js 后 WSL 自动继承"
    exit 1
fi

# npm 可能通过 corepack 管理（Debian/Ubuntu 新版 nodejs 包），shim 不在默认 PATH 里
# 把 corepack shim 目录加到 PATH，这样 command -v npm 能命中
if ! command -v npm >/dev/null 2>&1; then
    if [ -x /usr/share/nodejs/corepack/shims/npm ]; then
        export PATH="/usr/share/nodejs/corepack/shims:$PATH"
    fi
fi

if ! command -v npm >/dev/null 2>&1; then
    log_error "未找到 npm。尝试运行：sudo corepack enable"
    log_warn  "  或者重新安装 Node.js（含 npm）：curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash - && sudo apt install -y nodejs"
    exit 1
fi

log_info "node $(node -v)  npm $(npm -v)"

# ==============================================================================
# 1. 依赖检查（首次 clone 后 node_modules 不存在 → 自动 npm install）
# ==============================================================================
if [ ! -d "$FRONTEND_DIR/node_modules" ]; then
    log_info "[1/2] 未检测到 node_modules，执行 npm install（首次启动较慢，请耐心等待）..."
    cd "$FRONTEND_DIR"
    npm install
else
    log_info "[1/2] node_modules 已存在，跳过安装（如需重装：cd frontend && rm -rf node_modules && npm install）"
fi

# ==============================================================================
# 2. 启动 Vite 开发服务器
# ==============================================================================
log_info "[2/2] 启动 Vite 开发服务器：npm run dev"
log_info "  Ctrl+C 停止。默认地址 http://localhost:5173 （Vite 启动后会打印实际端口）"
log_info "  请确保后端已启动（bash start_backend.sh），否则登录 / 接口会报错。"

cd "$FRONTEND_DIR"
exec npm run dev
