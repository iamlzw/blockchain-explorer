module github.com/lifegoeson/blockchain-explorer

go 1.15

require (
	github.com/Shopify/sarama v1.29.1 // indirect
	github.com/fsouza/go-dockerclient v1.7.3 // indirect
	github.com/gin-gonic/gin v1.7.2
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/hashicorp/go-version v1.3.0 // indirect
	github.com/hyperledger/fabric v1.4.3
	github.com/hyperledger/fabric-amcl v0.0.0-20210603140002-2670f91851c8 // indirect
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/lib/pq v1.10.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/viper v1.3.2
	github.com/syndtr/goleveldb v1.0.0 // indirect
	go.uber.org/zap v1.18.1 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/hyperledger/fabric-sdk-go v1.0.0 => ../fabric-sdk-go
