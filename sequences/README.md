## Summary sequence of this repo

If a previous install is in-place, do a clean reset from the `../fabric-samples/test-network` folder:

`rm -rf peer0_org1_data/* peer0_org2_data/* orderer_data/*`

But better and more professional to uptick the label version number.

```
peer lifecycle chaincode package batterycc.tar.gz -p . --label batterylevelcc_1.0
```
peer lifecycle chaincode package batterycc.tar.gz -p . --label batterylevelcc_1.1
```
# Package (if not done)
cd /home/cartheur/go/src/github.com/cartheur/fabric-samples/chaincode/battery_level_chaincode/
peer lifecycle chaincode package batterycc.tar.gz -p . --label batterylevelcc_1.0

# Install (if not done)
set_org1
peer lifecycle chaincode install batterycc.tar.gz
set_org2
peer lifecycle chaincode install batterycc.tar.gz

# Get Package ID
set_org1
peer lifecycle chaincode queryinstalled  # Note the Package ID

# Approve
set_org1
peer lifecycle chaincode approveformyorg -o localhost:7050 --channelID haden --name batterylevelcc --version 1.0 --package-id batterylevelcc_1.0:<some-hash> --sequence 1 --tls --cafile /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
set_org2
peer lifecycle chaincode approveformyorg -o localhost:7050 --channelID haden --name batterylevelcc --version 1.0 --package-id batterylevelcc_1.0:<some-hash> --sequence 1 --tls --cafile /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# Commit
set_org1
peer lifecycle chaincode commit -o localhost:7050 --channelID haden --name batterylevelcc --version 1.0 --sequence 1 --tls --cafile /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses localhost:7051 --peerAddresses localhost:9051 --tlsRootCertFiles /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --tlsRootCertFiles /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

# Test
peer chaincode invoke -o localhost:7050 --tls --cafile /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C haden -n batterylevelcc -c '{"Args":["reportBattery","Robot1","8000","2025-03-20T10:03:00Z"]}'
cd /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network
./network.sh down
./network.sh up
set_org1
peer chaincode query -C haden -n batterylevelcc -c '{"Args":["queryBattery","Robot1"]}'
```