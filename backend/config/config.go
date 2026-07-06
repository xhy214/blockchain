package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	JWT    JWTConfig    `yaml:"jwt"`
	Fabric FabricConfig `yaml:"fabric"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type MySQLConfig struct {
	DSN string `yaml:"dsn"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type FabricConfig struct {
	PeerEndpoint  string `yaml:"peer_endpoint"`
	GatewayPeer   string `yaml:"gateway_peer"`
	Channel       string `yaml:"channel"`
	Chaincode     string `yaml:"chaincode"`
	MSPID         string `yaml:"msp_id"`
	CertPath      string `yaml:"cert_path"`
	KeyPath       string `yaml:"key_path"`
	TLSCert       string `yaml:"tls_cert"`
}

var DefaultConfig = Config{
	Server: ServerConfig{Port: "8080"},
	MySQL: MySQLConfig{
		DSN: "root:password@tcp(127.0.0.1:3306)/copyright_db?charset=utf8mb4&parseTime=True&loc=Local",
	},
	JWT: JWTConfig{
		Secret:      "replace_with_strong_secret_in_production",
		ExpireHours: 24,
	},
	Fabric: FabricConfig{
		PeerEndpoint:  "localhost:7051",
		GatewayPeer:   "peer0.org1.example.com",
		Channel:       "mychannel",
		Chaincode:     "copyright",
		MSPID:         "Org1MSP",
		CertPath:      "../network/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem",
		KeyPath:       "../network/crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/",
		TLSCert:       "../network/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt",
	},
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := DefaultConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
