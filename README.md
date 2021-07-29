hyperledger fabric 区块链浏览器，与官方的区块链浏览器功能大致相同，使用go 和vue开发

演示

![](C:\Users\iamlz\Desktop\运维\20210726101944.png)

![20210726102029](C:\Users\iamlz\Desktop\运维\20210726102029.png)

![20210726102058](C:\Users\iamlz\Desktop\运维\20210726102058.png)

![20210726102120](C:\Users\iamlz\Desktop\运维\20210726102120.png)

![20210726102146](C:\Users\iamlz\Desktop\运维\20210726102146.png)

![20210726102205](C:\Users\iamlz\Desktop\运维\20210726102205.png)

### 克隆仓库

```bash
# git clone https://github.com/iamlzw/blockchain-explorer.git
```

### 创建数据库

```bash
# cd blockchain-explorer/service
# psql createDB.sql
```

#### 修改postgresql 连接配置

修改service.go中的连接配置

```go
func SqlOpen() {
	var err error
	db, err = sql.Open("postgres", "port=5432 user=hppoc password=password dbname=fabricexplorer sslmode=disable")
	common.CheckErr(err)
}
```

### 修改配置文件

#### 修改config_e2e文件

```yaml
version: 1.0.0

client:
  organization: Org1
  logging:
    level: info
  # Root of the MSP directories with keys and certs.
  cryptoconfig:
    path: E:\workspace\go\src\github.com\lifegoeson\blockchain-explorer\crypto-config
  credentialStore:
    path: "/tmp/state-store"
    cryptoStore:
      path: /tmp/msp

channels:
  mychannel:
    peers:
      peer0.org1.example.com:
        endorsingPeer: true

        chaincodeQuery: true

        ledgerQuery: true

        eventSource: true

organizations:
  Org1:
    mspid: Org1MSP
    cryptoPath:  peerOrganizations\org1.example.com\users\{username}@org1.example.com\msp
    peers:
      - peer0.org1.example.com
peers:
  peer0.org1.example.com:
    url: grpcs://192.168.126.128:7051

    grpcOptions:
      ssl-target-name-override: peer0.org1.example.com
      keep-alive-time: 0s
      keep-alive-timeout: 20s
      keep-alive-permit: false
      fail-fast: false
      allow-insecure: false
    tlsCACerts:
      path: E:\workspace\go\src\github.com\lifegoeson\blockchain-explorer\crypto-config\peerOrganizations\org1.example.com\tlsca\tlsca.org1.example.com-cert.pem

entityMatchers:
  peer:
    - pattern: (\w+).org1.example.(\w+)
      urlSubstitutionExp: grpcs://192.168.126.128:7051
      sslTargetOverrideUrlSubstitutionExp: peer0.org1.example.com
      mappedHost: peer0.org1.example.com
```

#### 修改默认配置文件default_config.yaml

```bash
default_org_user: "User1"
default_org_name: "Org1"
default_peer_url: "grpcs://192.168.126.128:7051"
default_server_name: "peer0.org1.example.com"
default_channel_name: "mychannel"
default_org_admin: "Admin"
```

#### 修改服务发现配置文件discovery_config.yaml

```yaml
version: 0
tlsconfig:
  certpath: ""
  keypath: ""
  peercacertpath: E:\workspace\go\src\github.com\lifegoeson\blockchain-explorer\crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/tls/ca.crt
  timeout: 0s
signerconfig:
  mspid: Org1MSP
  identitypath: E:\workspace\go\src\github.com\lifegoeson\blockchain-explorer\crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem
  keypath: E:\workspace\go\src\github.com\lifegoeson\blockchain-explorer\crypto-config/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/c5dad79b0eb8ca81ce0078d204d3cc6872e5d64d64789c097dd2e30b2231ca6a_sk
```

### 启动服务

```bash
# cd blockchain-explorer/
# go run main.go
```

### 启动前端服务

#### 克隆仓库

```bash
# git clone https://github.com/iamlzw/blockchain-explorer-app.git
```

#### 启动服务

```bash
# cd blockchain-explorer-app
# npm install
# npm run dev
```



