# Atom Swap Demo
Atom swap demo between a Besu private network and a Fabric network. In this demo we atempt an atom swap transfer between a ERC20 token and a Fabric Token. 

Imagine that Alice has 100 GLDR token in a Besu private network and Bob has 100 GoTokens in a Fabric network. Alice wants to send 70 GLDR to Bob and Bob wants to send 70 GoTokens to Alice. They can do this by using an atom swap.

For that, a Hashtime Locked Contract (HTLC) is deployed in both networks. Alice locks 70 GLDR in the HTLC contract in the Besu network and Bob locks 70 GoTokens in the HTLC contract in the Fabric network. Both contracts have the same hashlock and timelock. Alice reveals the hashlock secret to Bob when she claims her GoTokens in the Fabric network. Bob then takes the revealed secret and claims the GLDR tokens locked in the HTLC contract in Besu.

# Besu
In the Besu section of this demo we have an ERC20 contract and a HTLC contract.

## Besu prequisites

- Besu binary (https://besu.hyperledger.org/private-networks/get-started/install/binary-distribution)
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

NOTE: This chaincode implements a *partial* ERC20 interface for demo purposes. Also, for simplicity, we are not signing requests made by a particular account. This means that anyone can mint, burn and transfer tokens for any account. This is not a problem for this demo, but should not be used in production.

## Fabric prequisites
- Fabric binaries (already available in the repository for the ones using Linux)
  - If using MacOS, you need to install the Fabric binaries yourself (https://hyperledger-fabric.readthedocs.io/en/release-2.5/install.html)
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

For the representation of the addresses of Bob and Alice we will use the same ones as in the Besu network, for simplicity. 
```
Alice: 0xfe3b557e8fb62b89f4916b721be55ceb828dbd73
Bob: 0x627306090abaB3A6e1400e9345bC60c78a8BEf57
```

We will create the wallets for Alice and Bob. To do that, you can use GoInitus by accessing `http://localhost:8080` in your browser or you can run the following `curl` command:
```bash
 curl -X \
 POST "http://localhost:80/api/invoke/createWallet" -H 'Content-Type: application/json' \
 -H 'cache-control: no-cache' \
  -d '{"address":"0xfe3b557e8fb62b89f4916b721be55ceb828dbd73","label":"Alice"}'

 curl -X \
 POST "http://localhost:80/api/invoke/createWallet" -H 'Content-Type: application/json' \
 -H 'cache-control: no-cache' \
  -d '{"address":"0x627306090abaB3A6e1400e9345bC60c78a8BEf57","label":"Bob"}'
```
Now, we will mint 100 tokens for Bob:
```bash
 curl -X \
 POST "http://localhost:80/api/invoke/mint" -H 'Content-Type: application/json' \
 -H 'cache-control: no-cache' \
  -d '{"amount":"100","to":{"@assetType":"wallet","@key":"wallet:20e5d9cf-f138-5e56-b921-ed05ce10c7ba"}}'
```

## The Demo itself

We've setup Alice and Bob's accounts in both networks. In the Besu network, Alice has 100 GLDR and in the Fabric network, Bob has 100 GoTokens. Now, we will do an atom swap between them.

### Step 1: Alice locks 70 GLDR in the HTLC contract in the Besu network

To lock 70 GLDR in the HTLC contract, we will use Truffle to call the `approve` function of ERC20 contract and then the `newSwap` function of the HTLC contract with the hash of the secret (in this example the secret is `mysecret`). To do that, run the following commands from the `besu/contracts` directory, this time using the `truffle-alice.js` configuration file:
```bash
truffle console --config truffle-alice.js
```

Once in the truffle console, run the following commands:
```javascript
let token_instance = await GoLedgerToken.deployed()
result = await token_instance.approve("0xF12b5dd4EAD5F743C6BaA640B0216200e89B60Da", 70)

let htlc_instance = await HTLCTokenSwap.deployed()
let result = await htlc_instance.newSwap("0xf17f52151EbEF6C7334FAD080c5704D77216b732", "0x41133a9ee618c42f1d3b40a69a750cb9e6df6b1801de6743012f829855ec7df0", 1000000000, 70)
```

### Step 2: Bob locks 70 GoTokens in the HTLC contract in the Fabric network with the same hashlock

To lock 70 GoTokens in the HTLC contract, we will use GoInitus to call the `newSwap` function of the HTLC chaincode with the hash of the secret. To do that, access `http://localhost:8080` in your browser and click on the `New Swap` transaction. Or run the following `curl` command:
```bash
 curl -X \
 POST "http://localhost:80/api/invoke/newSwap" -H 'Content-Type: application/json' \
 -H 'cache-control: no-cache' \
  -d '{"id":"1","fromWallet":{"@assetType":"wallet","@key":"wallet:20e5d9cf-f138-5e56-b921-ed05ce10c7ba"},"toWallet":{"@assetType":"wallet","@key":"wallet:6dc547ce-4fde-502a-96b6-4e2dbb6ac434"},"amount":"70","hashlock":"41133a9ee618c42f1d3b40a69a750cb9e6df6b1801de6743012f829855ec7df0","timelock":"2023-09-30T03:00:00.000Z"}'
``` 

### Step 3: Alice reveals the secret to Bob by finishing the swap in the Fabric network

Alice will now claim the 70 GoTokens by revealing the secret to Bob in the Fabric network.

Run the following `curl` command:
```bash
 curl -X POST "http://localhost:80/api/invoke/finishSwap" -H 'Content-Type: application/json' -H 'cache-control: no-cache'  -d '{"swap":{"@assetType":"hashTimeLock","@key":"hashTimeLock:996ad860-2a9a-504f-8861-aeafd0b2ae29"},"secret":"mysecret"}'
```
Check the balances of Alice and Bob and notice that Alice now has 70 GoTokens and Bob has 30 GoTokens.

### Step 4: Bob claims the 70 GLDR locked in the HTLC contract in the Besu network with the secret revealed

Bob will now claim the 70 GLDR by using the revealed secret in the Besu network. Run the following commands from the `besu/contracts` directory, this time using the `truffle-bob.js` configuration file:

```bash
truffle console --config truffle-bob.js
```

Once in the truffle console, run the following commands:
```javascript
let htlc_instance = await HTLCTokenSwap.deployed()
let result = await htlc_instance.finalizeSwap(1, "mysecret")
```

Check the balances of Alice and Bob and notice that Alice now has 30 GLDR and Bob has 70 GLDR. You can do that by calling the `balanceOf` function of the ERC20 contract in the Besu network:
```javascript
let token_instance = await GoLedgerToken.deployed()
let balance = await token_instance.balanceOf("0xfe3b557e8fb62b89f4916b721be55ceb828dbd73") // Alice
balance.toNumber()

balance = await token_instance.balanceOf("0xf17f52151EbEF6C7334FAD080c5704D77216b732") // Bob
balance.toNumber()
```

Our atom swap is complete!

# Challenges

* How to automate such a process?
* How to make it scalable?
* Where can it be used other than crypto exchange? 
  * Car ownership transfer, where a Fabric network holds car documents and a Besu network holds the fiat currency

# Improvements to be made

1. Implement a routine to check timelock expiration on Fabric's side