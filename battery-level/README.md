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

-----

Still getting the error:

`Error building image: docker build failed: Error returned from build: 1 "go: downloading go1.24.1 (linux/amd64)
battery_level_chaincode.go:6:5: missing go.sum entry for module providing package github.com/hyperledger/fabric-chaincode-go/shim (imported by battery_level_chaincode); to add:
        go get battery_level_chaincode`

-----

Using the proper fucking code, running:

```
go mod vendor
peer lifecycle chaincode package batterylevelcc.tar.gz -p . --label batterylevelcc_1.0
peer lifecycle chaincode install batterylevelcc.tar.gz
```

Gets:

```
2025-03-19 14:51:32.512 CET 0001 INFO [cli.lifecycle.chaincode] submitInstallProposal -> Installed remotely: response:<status:200 payload:"\nSbatterylevelcc_1.0:1833bd409463ffd8dbad6cc34ff5620e155901b7fd40ebdb6486a76e54485078\022\022batterylevelcc_1.0" > 
2025-03-19 14:51:32.513 CET 0002 INFO [cli.lifecycle.chaincode] submitInstallProposal -> Chaincode code package identifier: batterylevelcc_1.0:1833bd409463ffd8dbad6cc34ff5620e155901b7fd40ebdb6486a76e54485078
```

Check the results with:

`peer lifecycle chaincode queryinstalled`

Received:

```
Installed chaincodes on peer:
Package ID: batterylevelcc_1.0:1833bd409463ffd8dbad6cc34ff5620e155901b7fd40ebdb6486a76e54485078, Label: batterylevelcc_1.0
```

Approve:

`peer lifecycle chaincode approveformyorg -o localhost:7050 --channelID haden --name batterylevelcc --version 1.0 --package-id batterylevelcc_1.0:1833bd409463ffd8dbad6cc34ff5620e155901b7fd40ebdb6486a76e54485078 --sequence 1 --tls --cafile /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem`

If an error:

`peer channel list`

If there are no channels, will need to create one since doing a `docker-compose` on this one:

`./network.sh createChannel -c haden`

Go back and see if you can `Approve`.

Yes! Yields:

`2025-03-19 15:07:17.627 CET 0001 INFO [chaincodeCmd] ClientWait -> txid [738ec33189b475d25836dd4a4bd37ef55f9e4f1a649a955585b8fc1614e482d0] committed with status (VALID) at localhost:7051`

Commit:

`peer lifecycle chaincode commit -o localhost:7050 --channelID haden --name batterylevelcc --version 1.0 --sequence 1 --tls --cafile /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem`

Gets:

```
2025-03-19 15:08:02.023 CET 0001 INFO [chaincodeCmd] ClientWait -> txid [44a03e42b8bc59176bc2d3418f4219f0009af87b71339b78e8112317bc6d4223] committed with status (ENDORSEMENT_POLICY_FAILURE) at localhost:7051
Error: transaction invalidated with status (ENDORSEMENT_POLICY_FAILURE)
```
Because need approval from _both_ orgs, since Org1 is only approved from above. Now need to switch orgs in `export`:

```
export CORE_PEER_LOCALMSPID=Org2MSP
export CORE_PEER_MSPCONFIGPATH=/home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
export CORE_PEER_TLS_ROOTCERT_FILE=/home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
```
And the second approval:

`peer lifecycle chaincode approveformyorg -o localhost:7050 --channelID haden --name batterylevelcc --version 1.0 --package-id batterylevelcc_1.0:1833bd409463ffd8dbad6cc34ff5620e155901b7fd40ebdb6486a76e54485078 --sequence 1 --tls --cafile /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem`

Yes! Yields:

`2025-03-19 15:15:04.452 CET 0001 INFO [chaincodeCmd] ClientWait -> txid [765d5bad299d45b03df8f2bcba8a69cf32edeecdab7ca31e091ac0fa8ae1e1fc] committed with status (VALID) at localhost:9051`

Verify:

`peer chaincode invoke -o localhost:7050 --tls --cafile /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C haden -n batterylevelcc -c '{"Args":["reportBattery","Robot1","8000","2025-03-19T10:03:00Z"]}'`

Errors with:

`Error: endorsement failure during invoke. response: status:500 message:"make sure the chaincode batterylevelcc has been successfully defined on channel haden and try again: chaincode batterylevelcc not found" `

-----

Getting lost now. Check with:

`peer lifecycle chaincode querycommitted --channelID haden --name batterylevelcc`

Gives the error:

`Error: query failed with status: 404 - namespace batterylevelcc is not defined`
