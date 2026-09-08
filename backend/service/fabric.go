package service

import (
	"crypto/tls"
        "crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"blockchain/backend/config"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type FabricClient struct {
	contract *client.Contract
	gw       *client.Gateway
	conn     *grpc.ClientConn
}

func NewFabricClient(cfg *config.FabricConfig) (*FabricClient, error) {
	conn, err := newGrpcConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("grpc connection: %w", err)
	}

	id, err := newIdentity(cfg)
	if err != nil {
		return nil, fmt.Errorf("identity: %w", err)
	}

	sign, err := newSign(cfg)
	if err != nil {
		return nil, fmt.Errorf("sign: %w", err)
	}

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(conn),
		client.WithEvaluateTimeout(60 * time.Second),
		client.WithEndorseTimeout(60 * time.Second),
		client.WithSubmitTimeout(60 * time.Second),
		client.WithCommitStatusTimeout(60 * time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("gateway connect: %w", err)
	}

	network := gw.GetNetwork(cfg.Channel)
	contract := network.GetContract(cfg.Chaincode)

	return &FabricClient{
		contract: contract,
		gw:       gw,
		conn:     conn,
	}, nil
}

func (c *FabricClient) Submit(fn string, args ...string) ([]byte, error) {
	data, err := c.contract.SubmitTransaction(fn, args...)
	if err != nil {
		return data, unwrapChaincodeErr(err)
	}
	return data, nil
}

func (c *FabricClient) Evaluate(fn string, args ...string) ([]byte, error) {
	data, err := c.contract.EvaluateTransaction(fn, args...)
	if err != nil {
		return data, unwrapChaincodeErr(err)
	}
	return data, nil
}

// Fabric 网关把链码返回的原始错误包在 gRPC status details 里，
// 不解开的话前端只能看到 "failed to endorse transaction" 这类无用信息
func unwrapChaincodeErr(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	msg := ""
	for _, d := range st.Details() {
		if det, ok := d.(*gateway.ErrorDetail); ok {
			if msg != "" {
				msg += "; "
			}
			msg += det.GetMessage()
		}
	}
	if msg == "" {
		return err
	}
	return fmt.Errorf("%s", msg)
}

func (c *FabricClient) Close() {
	if c.gw != nil {
		c.gw.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func newGrpcConnection(cfg *config.FabricConfig) (*grpc.ClientConn, error) {
	certPool := x509.NewCertPool()
	pem, err := os.ReadFile(cfg.TLSCert)
	if err != nil {
		return nil, fmt.Errorf("read TLS cert: %w", err)
	}
	if !certPool.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("failed to parse TLS cert")
	}

	tlsConfig := &tls.Config{
		ServerName: cfg.GatewayPeer,
		RootCAs:    certPool,
	}
	return grpc.Dial(cfg.PeerEndpoint, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
}

func newIdentity(cfg *config.FabricConfig) (*identity.X509Identity, error) {
	certPEM, err := os.ReadFile(cfg.CertPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(certPEM)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
     		return nil, err
	}
 	return identity.NewX509Identity(cfg.MSPID, cert)
}

func newSign(cfg *config.FabricConfig) (identity.Sign, error) {
        files, err := os.ReadDir(cfg.KeyPath)
        if err != nil {
                return nil, err
        }
        var keyPath string
        for _, f := range files {
                if !f.IsDir() && len(f.Name()) >= 3 && f.Name()[len(f.Name())-3:] == "_sk" {
                        keyPath = filepath.Join(cfg.KeyPath, f.Name())
                        break
                }
        }
        if keyPath == "" {
                return nil, fmt.Errorf("no private key found in %s", cfg.KeyPath)
        }

        keyPEM, err := os.ReadFile(keyPath)
        if err != nil {
                return nil, err
        }

        block, _ := pem.Decode(keyPEM)
        if block == nil {
                return nil, fmt.Errorf("failed to decode PEM block")
        }

        key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
        if err != nil {
                // Try EC parse as fallback
                key, err = x509.ParseECPrivateKey(block.Bytes)
                if err != nil {
                        return nil, fmt.Errorf("failed to parse private key: %w", err)
                }
        }

        return identity.NewPrivateKeySign(key)
  }
