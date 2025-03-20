## Making what is ephemeral persist

After some intensive playing-around, add the lines to `compose/compose-test-net.yaml` _and_ `compose/docker/docker-compose-test-net.yaml`:

```
peer0.org1.example.com:
  # Base config...
  volumes:
    - ./peer0_org1_data:/var/hyperledger/production
peer0.org2.example.com:
  # Base config...
  volumes:
    - ./peer0_org2_data:/var/hyperledger/production
orderer.example.com:
  # Existing config...
  volumes:
    - ./orderer_data:/var/hyperledger/production
```
Create the `haden` channel:

`./network.sh createChannel -c haden`

`grep` to see how it went:

```
docker inspect peer0.org1.example.com | grep -A 5 "Mounts"
docker inspect peer0.org2.example.com | grep -A 5 "Mounts"
docker inspect orderer.example.com | grep -A 5 "Mounts"
```

Running the command a few times will result in the correct pointing to the new folders created in the `test-network` directory.

[Redeploy](/battery-level/README.md#package-the-chaincode-as-in-the-fabric-friendly-folder-inside-the-recommended-go-folder) the `haden` batterylevel package.

However, getting the error:

`InitCmd -> Fatal error when initializing core config : error when reading core config file: Config File "core" Not Found in "[/home/user/fabric-samples/chaincode/battery_level_chaincode]"`

The fabric cannot find the config path, so set it:

`export FABRIC_CFG_PATH=/home/cartheur/go/src/github.com/cartheur/fabric-samples/config/`

Then the next error:

`Error: failed to normalize chaincode path: failed to determine module root: exec: "go": executable file not found in $PATH`

Stupid go in the path, fixed.

Since not yet persistent, run the series of `export`s in the terminal:

```
cd /home/cartheur/go/src/github.com/cartheur/fabric-samples/chaincode/battery_level_chaincode/
export FABRIC_CFG_PATH=/home/cartheur/go/src/github.com/cartheur/fabric-samples/config/
peer lifecycle chaincode package batterycc.tar.gz -p . --label batterylevelcc_1.0
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH=/home/cartheur/go/src/github.com/cartheur/fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=/home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
peer lifecycle chaincode install batterycc.tar.gz
```

_Check Persistence_

`peer chaincode invoke -o localhost:7050 --tls --cafile /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C haden -n batterylevelcc -c '{"Args":["reportBattery","Robot1","8000","2025-03-20T10:03:00Z"]}'`