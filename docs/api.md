# **opEOS - 联调对接文档**

<details><summary>联调日志</summary>

| 修改类型 | 说明 | 日期    |
|---|---|-------|
| | |  |

</details>

<details><summary>合约地址</summary>

- <details><summary>L1(eosevmtest)</summary>

  ```js
  chainId: 15557
  L1_RPC: https://api.testnet.evm.eosnetwork.com

  L1StandardBridgeProxy address: 0x44332FC3Bf2a38e742DeF9d6cf20bd2ca4a3a9A8

  OptimismPortalProxy address : 0xdd52D429c7c85d2122EbEB3C5808fbf73caBe927

  TestERC20: 0x2475009a64a6846d02661611dc34C20A95eaBdFC
  ```

</details>


- <details><summary>L2(opEOS)</summary>

  ```js
  chainId: 4556

  L2StandardBridgeProxy address: 0x4200000000000000000000000000000000000010

  L2TestERC20 address: 0xb60bfc844edb68251e5ff04d5d1081ed294a0a6b
  ```
  </details>

</details>

## 方法调用

当前接口API未包括erc20 transfer, approve等接口

<details><summary>前端接口</summary>

- <details><summary class="green">depositETH: 充值ETH从L1到L2(caller: L1StandardBridgeProxy)</summary>

  ```javascript
  function depositETH(
            uint32 minGasLimit,              // 默认参数1000000
            bytes data                  //默认参数填写0x0
  )
  ```
  </details>

- <details><summary class="green">depositERC20: 充值ERC20从L1到L2(caller: L1StandardBridgeProxy)</summary>

  ```javascript
  function depositERC20(
    address l1TokenAddr,    //L1链的erc20地址
    address l2TokenAddr,    //L2链的erc20地址
    uint256 amount,         //充值金额
    uint32 minGasLimit,     //最低gas, 当前填写值是1000000, 后续需要准确再修改
    bytes data              //默认参数填写0x0
  )
  ```

  </details>

- <details><summary class="green">withdraw: 提现从L2到L1(caller: L2StandardBridgeProxy)</summary>

  ```javascript
  function withdraw(
    address addr,           //L2链的erc20地址, 如果提现是ETH, 填写地址为"0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000"
    uint256 amount,         //提现金额
    uint32 minGasLimit,     //最低gas, 当前填写值是1000000, 后续需要准确再修改
    bytes data              //默认参数填写0x0
  )
  ```

  </details>

- <details><summary class="green">prove: L1链Claim提现交易(caller: OptimismPortalProxy)</summary>

  ```javascript
  struct WithdrawalTransaction {
        uint256 nonce;
        address sender;
        address target;
        uint256 value;
        uint256 gasLimit;
        bytes data;
  }

  struct OutputRootProof {
        bytes32 version;
        bytes32 stateRoot;
        bytes32 messagePasserStorageRoot;
        bytes32 latestBlockhash;
  }

  function proveWithdrawalTransaction(
      Types.WithdrawalTransaction memory _tx,
      uint256 _l2OutputIndex,
      Types.OutputRootProof calldata _outputRootProof,
      bytes[] calldata _withdrawalProof
  )
  ```

  </details>

- <details><summary class="green">finalize: L1链Claim提现交易(caller: OptimismPortalProxy)</summary>

  ```javascript
  struct WithdrawalTransaction {
        uint256 nonce;
        address sender;
        address target;
        uint256 value;
        uint256 gasLimit;
        bytes data;
  }

  function finalizeWithdrawalTransaction(
    Types.WithdrawalTransaction memory _tx
  )
  ```

  </details>
</details>
