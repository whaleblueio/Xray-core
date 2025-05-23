module github.com/whaleblueio/Xray-core

go 1.17

replace github.com/whaleblueio/Xray-core => ./

require (
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.5
	github.com/gorilla/websocket v1.4.2
	github.com/juju/ratelimit v1.0.1
	github.com/lucas-clemente/quic-go v0.20.0
	//github.com/marten-seemann/qtls-go1-17 v0.1.0 // indirect
	github.com/miekg/dns v1.1.41
	github.com/pelletier/go-toml v1.8.1
	github.com/pires/go-proxyproto v0.5.0
	github.com/refraction-networking/utls v0.0.0-20201210053706-2179f286686b
	github.com/seiflotfy/cuckoofilter v0.0.0-20201222105146-bc6005554a0c
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/xtls/go v0.0.0-20201118062508-3632bf3b7499
	go.starlark.net v0.0.0-20210312235212-74c10e2c17dc
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007
	google.golang.org/grpc v1.36.1
	google.golang.org/protobuf v1.26.0
	h12.io/socks v1.0.2

)
