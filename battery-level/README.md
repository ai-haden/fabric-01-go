## Getting the code built and on the chain


1. Package the chaincode, as in the fabric-friendly folder, inside the recommended go folder:

`~/go/src/github.com/cartheur/fabric-samples/chaincode/battery_level_chaincode`

Run the command. Grok has a version that doesn't work in Fabric 2.5.11:

`peer lifecycle chaincode package batterylevelcc.tar.gz -p . --label batterylevelcc_1.0`

2. Install the chaincode

`peer lifecycle chaincode install batterylevelcc.tar.gz`

As an error occurs "Error: failed to retrieve endorser client for install: endorser client failed to connect to localhost:7051: failed to create new connection: context deadline exceeded":


