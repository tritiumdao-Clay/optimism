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

  L1StandardBridgeProxy address: 0x687E5591Cd5f1F48bE01886Ff9dE5012aD38C759

  TestERC20:
  ```

</details>


- <details><summary>L2(opEOS)</summary>

  ```js
  chainId: 42069

  L2StandardBridgeProxy address: 0x4200000000000000000000000000000000000010

  L2TestERC20 address:
  ```
  </details>

</details>

## 方法调用

当前接口API未包括erc20 transfer, approve等接口

<details><summary>前端接口</summary>

- <details><summary class="green">depositETH: 充值ETH从L1到L2(caller: L1StandardBridgeProxy)</summary>

  ```javascript
  function depositETH(
            uint32 amount,              // 充值金额, 单位 wei
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

- <details><summary class="green">finalize: 用户再L1链Claim提现交易</summary>

  //todo
  ```javascript
  function finalize(
  );
  ```

  </details>
</details>

