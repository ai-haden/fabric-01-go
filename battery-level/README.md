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
export CORE_PEER_MSPCONFIGPATH=/home/user/fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export FABRIC_CFG_PATH=/home/user/fabric/fabric-samples/config/
```
