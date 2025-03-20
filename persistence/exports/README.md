## Making these `export`s persistent and consistent

Edit `~/.bashrc`

`nano ~/.bashrc`

Add this to the end:

```
# Fabric Configuration
export FABRIC_CFG_PATH=/home/cartheur/go/src/github.com/cartheur/fabric-samples/config/

# Go Path (adjust if Go is installed elsewhere)
export PATH=$PATH:/usr/local/go/bin

# Peer Settings for Org1 (optional, see below)
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH=/home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=/home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
```
Reload

`source ~/.bashrc`

### Handling the Org1 and Org2 switching behavior

Create a script

`nano ~/fabric-env.sh`

Add the following

```
#!/bin/bash

# Base settings
export FABRIC_CFG_PATH=/home/cartheur/go/src/github.com/cartheur/fabric-samples/config/
export PATH=$PATH:/usr/local/go/bin

# Function to switch orgs
set_org1() {
    export CORE_PEER_LOCALMSPID=Org1MSP
    export CORE_PEER_MSPCONFIGPATH=/home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_TLS_ROOTCERT_FILE=/home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    echo "Switched to Org1"
}

set_org2() {
    export CORE_PEER_LOCALMSPID=Org2MSP
    export CORE_PEER_MSPCONFIGPATH=/home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    export CORE_PEER_ADDRESS=localhost:9051
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_TLS_ROOTCERT_FILE=/home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
    echo "Switched to Org2"
}

# Default to Org1
set_org1
```
Make it executable

`chmod +x ~/fabric-env.sh`

Add to `~/.bashrc`

```
echo "source ~/fabric-env.sh" >> ~/.bashrc
source ~/.bashrc
```

Use:

```
set_org1  # Switch to Org1
set_org2  # Switch to Org2
```
