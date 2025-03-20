## Having the chain be accessible from a DNS

Since GoDaddy is a PoS, will use self-signed certificates, since supported by Fabric.

Update the docker compose files

First `compose/compose-test-net.yaml`
```
cd /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network
nano compose/compose-test-net.yaml

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
nano compose/docker/docker-compose-test-net.yaml

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
nano /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json

{
  "name": "test-network-org1",
  "version": "1.0.0",
  "client": {
    "organization": "Org1"
  },
  "organizations": {
    "Org1": {
      "mspid": "Org1MSP",
      "peers": ["peer0.org1.example.com"]
    }
  },
  "orderers": {
    "orderer.example.com": {
      "url": "grpcs://slowclock.com:7050",
      "tlsCACerts": {
        "path": "/home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
      }
    }
  },
  "peers": {
    "peer0.org1.example.com": {
      "url": "grpcs://slowclock.com:7051",
      "tlsCACerts": {
        "path": "/home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
      }
    }
  }
}
```

Copy to the machine controlling the robot

```
cp /home/cartheur/go/src/github.com/cartheur/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json /home/user/robot-client-dotnet/
```

Restart the network
```
./network.sh down
./network.sh up
./network.sh createChannel -c haden  # If data cleared
```

Open firewall ports

```
sudo ufw allow 7050
sudo ufw allow 7051
sudo ufw allow 9051
sudo ufw status
```

### Modify the c# code

```
using Hyperledger.Fabric.SDK;
using Hyperledger.Fabric.SDK.Channels;
using Hyperledger.Fabric.SDK.Identity;
using Hyperledger.Fabric.SDK.Security;
using System;
using System.IO;
using System.Net.Http;
using System.Threading.Tasks;

namespace BatteryReporter
{
    class Program
    {
        static async Task Main(string[] args)
        {
            if (args.Length != 3)
            {
                Console.WriteLine("Usage: dotnet run -- <robotID> <batteryLevel> <timestamp>");
                return;
            }
            await ReportBattery(args[0], args[1], args[2]);
        }

        static async Task ReportBattery(string robotID, string batteryLevel, string timestamp)
        {
            try
            {
                // Trust self-signed certs
                var handler = new HttpClientHandler();
                handler.ServerCertificateCustomValidationCallback = (message, cert, chain, errors) => true;
                HFClient client = HFClient.Create(handler);

                // Crypto suite
                ICryptoSuite cryptoSuite = CryptoSuite.Factory.GetCryptoSuite();
                client.CryptoSuite = cryptoSuite;

                // Load user
                string certPath = Path.Combine(Directory.GetCurrentDirectory(), "wallet", "admin-cert.pem");
                string keyPath = Path.Combine(Directory.GetCurrentDirectory(), "wallet", "admin-key.pem");
                string cert = File.ReadAllText(certPath);
                string key = File.ReadAllText(keyPath);
                IEnrollment enrollment = new X509Enrollment(key, cert);
                User user = new SampleUser("admin", "Org1MSP", enrollment);
                client.UserContext = user;

                // Load connection profile
                string connectionProfilePath = Path.Combine(Directory.GetCurrentDirectory(), "connection-org1.json");
                string connectionProfileJson = File.ReadAllText(connectionProfilePath);

                // Connect to channel
                Channel channel = client.NewChannel("haden");
                channel.Initialize(connectionProfileJson);

                // Submit transaction
                ChaincodeID chaincodeID = ChaincodeID.NewBuilder().SetName("batterylevelcc").Build();
                TransactionProposalRequest request = client.NewTransactionProposalRequest();
                request.SetChaincodeID(chaincodeID);
                request.SetFcn("reportBattery");
                request.SetArgs(robotID, batteryLevel, timestamp);

                var responses = await channel.SendTransactionProposalAsync(request);
                var transaction = await channel.SendTransactionAsync(responses);
                Console.WriteLine($"Battery reported: {robotID}, {transaction.TransactionID}");
            }
            catch (Exception e)
            {
                Console.WriteLine($"Failed: {e.Message}");
            }
        }
    }
}
```

Run the caller
```
static void ReportBatteryLevel()
{
    string robotID = "Robot1";
    string batteryLevel = GetBatteryLevel().ToString();
    string timestamp = DateTime.UtcNow.ToString("yyyy-MM-ddTHH:mm:ssZ");
    string[] args = { robotID, batteryLevel, timestamp };
    Task.Run(() => Program.Main(args)).Wait();
}
```

Verify (on host)

```
set_org1
peer chaincode query -C haden -n batterylevelcc -c '{"Args":["queryBattery","Robot1"]}'
```