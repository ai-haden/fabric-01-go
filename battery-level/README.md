## Getting the code built and on the chain


Why no `peer` at the terminal for deployment, although installed. Try:

```
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_MSPCONFIGPATH=$PWD/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_ADDRESS=localhost:7051
```

```
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/bootstrap.sh
chmod +x bootstrap.sh
./bootstrap.sh 2.5.0 1.4.9
```