var GoLedgerToken = artifacts.require("./GoLedgerToken.sol");

module.exports = function(deployer) {
  deployer.deploy(GoLedgerToken, 1000000000, {gas: 5000000});
};
