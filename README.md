# Atom Swap Demo
Atom swap demo between a Besu private network and a Fabric network. In this demo atempt an atom swap transfer between a ERC20 token and a Fabric UTXO Token. 

Imagine that Alice has 100 GLDR token in a Besu private network and Bob has 100 GoTokens in a Fabric network. Alice wants to send 50 GLDR to Bob and Bob wants to send 50 GoTokens to Alice. They can do this by using an atom swap.

For that, a Hashtime Locked Contract (HTLC) is deployed in both networks. Alice locks 50 GLDR in the HTLC contract in the Besu network and Bob locks 50 GoTokens in the HTLC contract in the Fabric network. Both contracts have the same hashlock and timelock. Alice reveals the hashlock secret to Bob when she claims her GoTokens in the Fabric network. Bob then takes the revealed secret and claims the GLDR tokens locked in the HTLC contract in Besu.

# Besu
In the Besu section of this demo we have an ERC20 contract an a HTLC contract.

## Besu prequisites

- Besu binary
- Truffle
    - That means you need NodeJS and NPM
- Docker and Docker Compose

## Besu setup
To setup the Besu network and deploy the contracts simply run the following commands from the repository's root directory:
```bash
cd besu
./startDev.sh
```
Your network will be up and running in a few seconds. You can check the status of the network by running:
```bash
docker logs besu-node-0 -f
```

# Fabric
In the Fabric section of this demo we have a GoToken chaincode with HTLC functionality. 

NOTE: This chaincode implements a *partial* ERC20 interface for demo purposes. Also, for simplicity, we are not signing requests made by a particular account. This means that anyone can mint, burn and transfer tokens for any account. This is not a problem for this demo, but it is not recommended for production.

## Fabric prequisites
- Fabric binaries
- Go 1.13 minimum
- Docker and Docker Compose

## Fabric setup
To setup the Fabric network and deploy the chaincode simply run the following commands from the repository's root directory:
```bash
cd fabric
./startDev.sh
```
Your network will be up and running in a few seconds. You can check the status of the network by running:
```bash
docker logs peer0.org1.example.com -f
```
For easier interaction with the Fabric network, we will deploy GoInitus, a web application that allows you to interact with an API. To do that, run the following commands from the `fabric` directory:
```bash
./run-cc-web.sh
```

# HTLC Demo

## Setup Besu accounts for the demo
To setup the accounts we will use Truffle to transfer GLDR from the GoLedgerToken contract to Alice's account. To do that, run the following commands from the `besu/contracts` directory:
```bash
truffle console --config truffle-config.js
```
Once in the truffle console, run the following commands:
```javascript
let token_instance = await GoLedgerToken.deployed()
let result = await token_instance.transfer("0xfe3b557e8fb62b89f4916b721be55ceb828dbd73", 100)
```
This will transfer 100 GLDR tokens to Alice's account. You can check the balance of Alice's account by running:
```javascript
let balance = await token_instance.balanceOf("0xfe3b557e8fb62b89f4916b721be55ceb828dbd73") 
balance.toNumber()
```

## Setup Fabric accounts for the demo

In the `fabric/keys` directory there are public and private keys for Alice and Bob. 
```
keys
├── alice.key
├── alice.pub
├── bob.key
└── bob.pub
```

Now, we will mint tokens for Bob. To do that, you can use GoInitus by accessing `http://localhost:8080` in your browser or you can run the following `curl` command:
```bash
 curl -X \
 POST "http://localhost:80/api/invoke/mint" -H 'Content-Type: application/json' \
 -H 'cache-control: no-cache' \
  -d '{"amount":"100","to":"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUFySHR5bW9jZnY4QmNmUWgwZ2xacwpNdnF3bzNWUUVEZzI5bmhuRjJJTitFV2RTWE5oblRsd3h0MFZTSExMV2c5OCtFbUtaTUFFZ0tHRExUd2FaMXh5Cmh1QTIzTFF3QlgvQ08rWjVyNDBuQ2U3QTk1NTFrNTJGL3ZKdVVrUFJHMUQwNEhYN0hMYkhYaVBoRlM5eDk4c3MKZHljY0hnejBqbU1DMENWN0tUR3N5L3VKZlFPV1diYUM2MXZIK0RKOURqeGMxMVdKY3pFaDFNTlYzVFFZS1pCaQpPWlk1Q1NaRmtWOCsxdzFIYkN2RjA4UW9VMVFZSm1XdkFueVAzVFhzOUVhdER0Y0ZNL3RrbHRKaTdPaFlpVlB1Cko0QmVkZEo5TTEvNFJBVFdIcTBrb3VvdjNLdGVVeTFjUi85akJpYmo1UmZLUno4YmhDVUgxdTh3NXpSS3ZYNWwKN3dJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="}'
```

## The Demo itself

