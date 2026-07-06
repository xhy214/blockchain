#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NETWORK_DIR="$(dirname "$SCRIPT_DIR")"
ORDERER_CA="/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

echo "=== 1. 打包链码 ==="
docker exec cli peer lifecycle chaincode package copyright.tar.gz \
    --path /opt/gopath/src/github.com/hyperledger/fabric/peer/chaincode/copyright \
    --lang golang \
    --label copyright_1.0

echo "=== 2. Org1 安装链码 ==="
docker exec cli peer lifecycle chaincode install copyright.tar.gz

echo "=== 3. Org2 安装链码 ==="
docker exec \
    -e CORE_PEER_ADDRESS=peer0.org2.example.com:9051 \
    -e CORE_PEER_LOCALMSPID=Org2MSP \
    -e CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/server.crt \
    -e CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/server.key \
    -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
    -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp \
    cli peer lifecycle chaincode install copyright.tar.gz

echo "=== 4. 获取 Package ID ==="
PACKAGE_ID=$(docker exec cli peer lifecycle chaincode queryinstalled 2>&1 | grep "copyright_1.0" | awk -F'Package ID: ' '{print $2}' | awk -F', Label' '{print $1}')
echo "Package ID: $PACKAGE_ID"

echo "=== 5. Org1 审批链码 ==="
docker exec cli peer lifecycle chaincode approveformyorg \
    -o orderer.example.com:7050 \
    --channelID mychannel \
    --name copyright \
    --version 1.0 \
    --package-id "$PACKAGE_ID" \
    --sequence 1 \
    --tls --cafile "$ORDERER_CA"

echo "=== 6. Org2 审批链码 ==="
docker exec \
    -e CORE_PEER_ADDRESS=peer0.org2.example.com:9051 \
    -e CORE_PEER_LOCALMSPID=Org2MSP \
    -e CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/server.crt \
    -e CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/server.key \
    -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
    -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp \
    cli peer lifecycle chaincode approveformyorg \
    -o orderer.example.com:7050 \
    --channelID mychannel \
    --name copyright \
    --version 1.0 \
    --package-id "$PACKAGE_ID" \
    --sequence 1 \
    --tls --cafile "$ORDERER_CA"

echo "=== 7. 检查提交就绪状态 ==="
docker exec cli peer lifecycle chaincode checkcommitreadiness \
    --channelID mychannel \
    --name copyright \
    --version 1.0 \
    --sequence 1 \
    --output json

echo "=== 8. 提交链码 ==="
docker exec cli peer lifecycle chaincode commit \
    -o orderer.example.com:7050 \
    --channelID mychannel \
    --name copyright \
    --version 1.0 \
    --sequence 1 \
    --peerAddresses peer0.org1.example.com:7051 \
    --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
    --peerAddresses peer0.org2.example.com:9051 \
    --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
    --tls --cafile "$ORDERER_CA"

echo "=== 9. 验证链码已提交 ==="
docker exec cli peer lifecycle chaincode querycommitted \
    --channelID mychannel --name copyright

echo "=== 链码部署完成 ==="
