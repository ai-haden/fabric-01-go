## Having the chain be accessible from a DNS

Since GoDaddy is a PoS, will use self-signed certificates, since supported by Fabric.

Update DNS in `compose/compose-test-net.yaml`
```
orderer.example.com:
  environment:
    - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
    - ORDERER_GENERAL_TLS_ENABLED=true
  ports:
    - "7050:7050"
peer0.org1.example.com:
  environment:
    - CORE_PEER_ADDRESS=slowclock.com:7051
    - CORE_PEER_GOSSIP_EXTERNALENDPOINT=slowclock.com:7051
    - CORE_PEER_TLS_ENABLED=true
  ports:
    - "7051:7051"
peer0.org2.example.com:
  environment:
    - CORE_PEER_ADDRESS=slowclock.com:9051
    - CORE_PEER_GOSSIP_EXTERNALENDPOINT=slowclock.com:9051
    - CORE_PEER_TLS_ENABLED=true
  ports:
    - "9051:9051"
```

And `compose/docker/docker-compose-test-net.yaml`

```
peer0.org1.example.com:
  environment:
    - CORE_PEER_ADDRESS=slowclock.com:7051
    - CORE_PEER_GOSSIP_EXTERNALENDPOINT=slowclock.com:7051
    - CORE_PEER_TLS_ENABLED=true
  ports:
    - "7051:7051"
peer0.org2.example.com:
  environment:
    - CORE_PEER_ADDRESS=slowclock.com:9051
    - CORE_PEER_GOSSIP_EXTERNALENDPOINT=slowclock.com:9051
    - CORE_PEER_TLS_ENABLED=true
  ports:
    - "9051:9051"
```
And `connection-org1.json`

```
"orderers": {
  "orderer.example.com": {
    "url": "grpcs://slowclock.com:7050",
    "tlsCACerts": { "path": "/path/to/tlsca.example.com-cert.pem" }
  }
},
"peers": {
  "peer0.org1.example.com": {
    "url": "grpcs://slowclock.com:7051",
    "tlsCACerts": { "path": "/path/to/org1-tls-ca.crt" }
  }
}
```