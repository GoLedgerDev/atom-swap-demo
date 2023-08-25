# Atom Swap Demo
Atom swap demo between a Besu private network and a Fabric network. In this demo atempt an atom swap transfer between a ERC20 token and a Fabric UTXO Token. 

Imagine that Alice has 100 GLDR token in a Besu private network and Bob has 100 GoTokens in a Fabric network. Alice wants to send 50 GLDR to Bob and Bob wants to send 50 GoTokens to Alice. They can do this by using an atom swap.

For that, a Hashtime Locked Contract (HTLC) is deployed in both networks. Alice locks 50 GLDR in the HTLC contract in the Besu network and Bob locks 50 GoTokens in the HTLC contract in the Fabric network. Both contracts have the same hashlock and timelock. Alice reveals the hashlock secret to Bob when she claims her GoTokens in the Fabric network. Bob then takes the revealed secret and claims the GLDR tokens locked in the HTLC contract in Besu.

# Besu
In the Besu section of this demo we have an ERC20 contract an a HTLC contract.

## Besu prequisites

- Besu binary
- Truffle
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

# HTLC Demo

## Setup accounts for the demo
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
