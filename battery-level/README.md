## Getting the code built and on the chain


1. Package the chaincode, as in the fabric-friendly folder, inside the recommended go folder:

`~/go/src/github.com/cartheur/fabric-samples/chaincode/battery_level_chaincode`

Run the command. Grok has a version that doesn't work in Fabric 2.5.11:

`peer lifecycle chaincode package batterylevelcc.tar.gz -p . --label batterylevelcc_1.0`

2. Install the chaincode

`peer lifecycle chaincode install batterylevelcc.tar.gz`

As an error occurs "Error: failed to retrieve endorser client for install: endorser client failed to connect to localhost:7051: failed to create new connection: context deadline exceeded":

```
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH=/home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export FABRIC_CFG_PATH=/home/cartheur/go/src/github.com/cartheur/fabric-samples/config/
```
Despite this:

`Error: failed to retrieve endorser client for install: endorser client failed to connect to localhost:7051: failed to create new connection: context deadline exceeded`

Check with docker logs on the container:

`docker logs peer0.org1.example.com`

_TLS is the issue_

```
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=/home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
```

With TLS being solved, the next error is:

`Error: chaincode install failed with status: 500 - failed to invoke backing implementation of 'InstallChaincode': could not build chaincode: docker build failed: docker image build failed: docker build failed: Error returned from build: 1 "go: downloading go1.24.1 (linux/amd64)
`
_Try without TLS_

```
export CORE_PEER_TLS_ENABLED=false
peer lifecycle chaincode install batterycc.tar.gz
```

A different error but still not effective. Issue is other networks in the instance, for myself the experiments with Prometheus and Kind, so

```
./network.sh down
docker rm -f $(docker ps -aq)
docker network prune
./network.sh up
```

However, we now get a docker build error:

`battery_level_chaincode.go:6:5: missing go.sum entry for module providing package github.com/hyperledger/fabric-chaincode-go/shim (imported by battery_level_chaincode); to add:
        go get battery_level_chaincode`

Most likely a Go compatibility issue, since (currently) Fabric uses 1.18, the fix is to modify go.mod to:

```
module battery_level_chaincode

go 1.18

require (
    github.com/hyperledger/fabric-chaincode-go v0.0.0-20220920211402-79e4c7985b55
    github.com/hyperledger/fabric-protos-go v0.3.0
)
```

Then:

```
go get github.com/hyperledger/fabric-chaincode-go/shim@latest
go get github.com/hyperledger/fabric-protos-go/peer@latest
go mod tidy
go mod vendor
```
But:

`go: github.com/hyperledger/fabric-chaincode-go@v0.0.0-20220920211402-79e4c7985b55: invalid version: unknown revision 79e4c7985b55`
 
But:

`go.mod:6:5: usage: require module/path v1.2.3`

Replace with:

```
github.com/hyperledger/fabric-chaincode-go v0.0.0-20220920211402-6ab0e950e7dd
github.com/hyperledger/fabric-protos-go v0.3.0
```

And:

`peer lifecycle chaincode package batterycc.tar.gz -p . --label batterycc_1.0`

