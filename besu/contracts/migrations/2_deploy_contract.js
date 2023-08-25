const HTLC_Contract = artifacts.require("HTLCTokenSwap.sol");

module.exports = function(deployer) {
  deployer.deploy(HTLC_Contract, "0x8CdaF0CD259887258Bc13a92C0a6dA92698644C0", {gas: 5000000});
};