## Multiplexed `export`

### Do the deployment

```
cd /home/cartheur/go/src/github.com/cartheur/fabric-samples/chaincode/battery_level_chaincode/
peer lifecycle chaincode package batterylevelcc.tar.gz -p . --label batterylevelcc_1.0
set_org1
peer lifecycle chaincode install batterylevelcc.tar.gz
set_org2
peer lifecycle chaincode install batterylevelcc.tar.gz
```

### Verify the install
```
cd /home/cartheur/go/src/github.com/cartheur/fabric-samples/chaincode/battery_level_chaincode/
set_org1
peer lifecycle chaincode queryinstalled
set_org2
peer lifecycle chaincode queryinstalled
```
### Approve by both orgs
```
set_org1
peer lifecycle chaincode approveformyorg -o localhost:7050 --channelID haden --name batterylevelcc --version 1.1 --package-id batterylevelcc_1.1:df995bb9611f838b8143eaa4d62bd1779eee51d1e1412572e379c2902f344eb0 --sequence 1 --tls --cafile /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

set_org2
peer lifecycle chaincode approveformyorg -o localhost:7050 --channelID haden --name batterylevelcc --version 1.1 --package-id batterylevelcc_1.1:df995bb9611f838b8143eaa4d62bd1779eee51d1e1412572e379c2902f344eb0 --sequence 1 --tls --cafile /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```
### Commit
```
set_org1
peer lifecycle chaincode commit -o localhost:7050 --channelID haden --name batterylevelcc --version 1.1 --sequence 1 --tls --cafile /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses localhost:7051 --peerAddresses localhost:9051 --tlsRootCertFiles /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --tlsRootCertFiles /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
```
Expected
```
Committed chaincode definition for chaincode 'batterylevelcc' on channel 'haden':
Version: 1.0, Sequence: 1, Endorsement Plugin: escc, Validation Plugin: vscc
```

### Test invocation (on org1)
```
peer chaincode invoke -o localhost:7050 --tls --cafile /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C haden -n batterylevelcc -c '{"Args":["reportBattery","Robot1","8000","2025-03-20T10:03:00Z"]}'
```

### Query (for the queryBattery addition to the code)

`peer chaincode query -C haden -n batterylevelcc -c '{"Args":["queryBattery","Robot1"]}'`

### Test persistence

```
cd /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network

./network.sh down
./network.sh up

set_org1
peer chaincode query -C haden -n batterylevelcc -c '{"Args":["queryBattery","Robot1"]}'
```

### Query again

```
set_org1
peer chaincode query -C haden -n batterylevelcc -c '{"Args":["queryBattery","Robot1"]}'
```