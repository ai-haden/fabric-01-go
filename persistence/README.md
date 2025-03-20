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

_Check Persistence_

`peer chaincode invoke -o localhost:7050 --tls --cafile /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C haden -n batterylevelcc -c '{"Args":["reportBattery","Robot1","8000","2025-03-20T10:03:00Z"]}'`