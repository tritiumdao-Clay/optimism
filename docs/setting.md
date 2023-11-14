

### Layer2设置


// deploy-config/getting-started.json
```json
{
    "baseFeeVaultRecipient": "string", //baseFee收取地址
    "l1FeeVaultRecipient": "string", //l1 calldata费用收取地址
    "sequencerFeeVaultRecipient": "string" //tipFee收取地址

    "baseFeeVaultMinimumWithdrawalAmount": uint256 //最小提取金额 
    "l1FeeVaultMinimumWithdrawalAmount": uint256 //最小提取金额
    "sequencerFeeVaultWithdrawalAmount": uint256 //最小提取金额

    "gasPriceOracleOverhead": uint256 //每笔交易固定的L1费用
    "gasPriceOracleScalar": uint256 //每笔交易动态的L1费用

    "eip1559Denominator": uint256 //
    "eip1559Elasticity": uint256 //

    "l2GenesisBlockGasLimit": "string" //
    "l2GenesisBlockBaseFeePerGas": "string" //
}


```
