const HTLC_Contract = artifacts.require("HTLCTokenSwap.sol");

module.exports = function(deployer) {
  deployer.deploy(HTLC_Contract, "0x0E6dD3FaCDB7F50484e50EFb83A12291589c6073", {gas: 5000000});
};