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

### Test persistence

```
cd /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network
./network.sh down
./network.sh up
set_org1
peer chaincode query -C haden -n batterylevelcc -c '{"Args":["queryBattery","Robot1"]}'
```

### Approve by both orgs



### Committ

