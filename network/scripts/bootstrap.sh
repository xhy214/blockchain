#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NETWORK_DIR="$(dirname "$SCRIPT_DIR")"

# ------------------------------------------------------------------
# 自动定位 Fabric 工具链（cryptogen / configtxgen）
# 兼容：普通用户 PATH、sudo 重置 PATH、fabric-samples 默认安装位置
# ------------------------------------------------------------------
find_fabric_bin() {
    # —— WSL 兼容性：获取"非 sudo 时"真正登录用户的家目录
    #    因为 sudo -E 下 $HOME 可能还是 root/的 /root，但 fabric-samples 往往装在普通用户 HOME 下
    local real_home="$HOME"
    if [ -n "$SUDO_USER" ] && [ "$SUDO_USER" != "root" ]; then
        # 有 sudo 但 SUDO_USER 是普通用户 → 用 getent passwd 拿到他的 $HOME（ubuntu 兼容）
        local uhome
        uhome="$(getent passwd "$SUDO_USER" 2>/dev/null | cut -d: -f6 || true)"
        if [ -d "$uhome" ]; then
            real_home="$uhome"
        else
            # 兜底：/home/$SUDO_USER（WSL / Ubuntu 默认结构）
            real_home="/home/$SUDO_USER"
        fi
    fi

    # 候选列表：既包含当前用户家目录，也包含 WSL 下把 fabric-samples 放在 Windows /mnt/c/Users 上的常见路径
    local candidates=(
        "$real_home/hyperledger/fabric-samples/bin"
        "$HOME/hyperledger/fabric-samples/bin"
        # WSL 常见：同事可能把 fabric-samples 下载在 Windows 家目录，直接被 /mnt/c/Users 挂载
        "/mnt/c/Users/$USER/hyperledger/fabric-samples/bin"
        "/mnt/c/Users/$USER/fabric-samples/bin"
        "/mnt/c/hyperledger/fabric-samples/bin"
        # Linux 系统级目录
        "/usr/local/bin"
        "/usr/bin"
        "/opt/fabric/bin"
        "/opt/hyperledger-fabric/bin"
    )
    for d in "${candidates[@]}"; do
        if [ -x "$d/cryptogen" ] && [ -x "$d/configtxgen" ]; then
            echo "$d"
            return 0
        fi
    done
    # 最后兜底：用 which 查当前执行用户（即使用户已 export PATH）
    local cg_path
    cg_path="$(command -v cryptogen 2>/dev/null || true)"
    if [ -n "$cg_path" ]; then
        dirname "$cg_path"
        return 0
    fi
    return 1
}

FABRIC_BIN="$(find_fabric_bin || true)"
if [ -z "$FABRIC_BIN" ]; then
    # —— 给 WSL 的更有用的错误提示（同事如果路径没放对，能直接看到他机器上当前有哪些 fabric-samples 候选） ——
    echo "ERROR: 找不到 cryptogen / configtxgen。"
    echo ""
    echo "  Ubuntu / WSL 建议安装方式："
    echo "    curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh"
    echo "    bash install-fabric.sh binary"
    echo "  这会把二进制装到: ./fabric-samples/bin"
    echo ""
    echo "  或者把已经下载的 fabric-samples/bin 目录放到以下任一位置（都识别得到）："
    echo "    - ~/hyperledger/fabric-samples/bin          （推荐，当前用户家目录）"
    echo "    - /mnt/c/Users/你的Windows用户名/hyperledger/fabric-samples/bin  （WSL 放在 Windows 盘上）"
    echo "    - /usr/local/bin                           （系统级）"
    echo ""
    echo "  当前检测到的执行用户信息："
    echo "    USER=$USER   SUDO_USER=${SUDO_USER:-<空>}   HOME=$HOME"
    echo "    PATH=$PATH"
    exit 1
fi
export PATH="$FABRIC_BIN:$PATH"
echo "使用 Fabric 工具链: $FABRIC_BIN (cryptogen $(cryptogen version 2>/dev/null | head -1))"

cd "$NETWORK_DIR"

echo "=== 清理旧的证书和通道工件 ==="
docker compose down -v --remove-orphans 2>/dev/null || true
rm -rf crypto-config channel-artifacts

echo "=== 1. 生成证书 ==="
cryptogen generate --config=crypto-config.yaml --output=crypto-config

echo "=== 2. 生成创世块和通道配置 ==="
mkdir -p channel-artifacts
export FABRIC_CFG_PATH="$NETWORK_DIR"

configtxgen -profile TwoOrgsOrdererGenesis -channelID system-channel \
    -outputBlock ./channel-artifacts/genesis.block

configtxgen -profile TwoOrgsChannel \
    -outputCreateChannelTx ./channel-artifacts/mychannel.tx \
    -channelID mychannel

configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate \
    ./channel-artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP

configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate \
    ./channel-artifacts/Org2MSPanchors.tx -channelID mychannel -asOrg Org2MSP

echo "=== 3. 启动网络容器 ==="
docker compose up -d
echo "等待容器就绪..."
sleep 8

ORDERER_CA="$NETWORK_DIR/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
ORG1_PEER_TLS="$NETWORK_DIR/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"

echo "=== 4. 创建通道 ==="
docker exec cli peer channel create \
    -o orderer.example.com:7050 \
    -c mychannel \
    -f /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/mychannel.tx \
    --tls \
    --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

echo "=== 5. Org1 peer0 加入通道 ==="
docker exec cli peer channel join -b mychannel.block

echo "=== 6. Org2 peer0 加入通道 ==="
docker exec \
    -e CORE_PEER_ADDRESS=peer0.org2.example.com:9051 \
    -e CORE_PEER_LOCALMSPID=Org2MSP \
    -e CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/server.crt \
    -e CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/server.key \
    -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
    -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp \
    cli peer channel join -b mychannel.block

echo "=== 7. 更新锚节点 ==="
docker exec cli peer channel update \
    -o orderer.example.com:7050 \
    -c mychannel \
    -f /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/Org1MSPanchors.tx \
    --tls \
    --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

docker exec \
    -e CORE_PEER_ADDRESS=peer0.org2.example.com:9051 \
    -e CORE_PEER_LOCALMSPID=Org2MSP \
    -e CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/server.crt \
    -e CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/server.key \
    -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
    -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp \
    cli peer channel update \
    -o orderer.example.com:7050 \
    -c mychannel \
    -f /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/Org2MSPanchors.tx \
    --tls \
    --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

echo "=== 网络启动完成，运行 deploy.sh 部署链码 ==="
