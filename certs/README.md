## Running with CA Certs

`./network up -ca`

Will get three containers for each of the three: two peers and one orderer.

###  Update Wallet with Correct Certs

If Using CA Mode: Enroll an admin to get matching certs

```
export PATH=${PWD}/../bin:$PATH
export FABRIC_CA_CLIENT_HOME=${PWD}
fabric-ca-client enroll -u http://admin:adminpw@localhost:7054 --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
```
With the error:

```
Error: Failed to parse response: Client sent an HTTP request to an HTTPS server.
: invalid character 'C' looking for beginning of value
```

Copy certs to wallet-ca

```
cp ${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/cert.pem wallet-ca/admin-cert.pem
cp ${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystores/*_sk wallet-ca/admin-key.pem
cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt wallet-ca/tls-ca.crt
```
