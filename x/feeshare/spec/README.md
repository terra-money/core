<!--
order: 0
title: "FeeShare Overview"
parent:
  title: "feeshare"
-->

# `feeshare`

## Abstract

This document specifies the internal `x/feeshare` module, which was originally developed by the Juno Network. This documentation is a fork of the [original documentation](https://github.com/CosmosContracts/juno/tree/main/x/feeshare/spec). 

This custom implementation of the `x/feeshare` module enables the splitting of revenue from transaction fees between validators and registered smart contracts. Developers can register their smart contracts and every time someone interacts with a registered smart contract, the contract deployer or their assigned withdrawal account receives a part of the transaction fees. If multiple contracts are involved in a transaction, the FeeShare revenue is split evenly between all registered contracts that participated in the transaction. 

## Contents

1. **[Concepts](01_concepts.md)**
2. **[State](02_state.md)**
3. **[State Transitions](03_state_transitions.md)**
4. **[Transactions](04_transactions.md)**
5. **[Post](05_post.md)**
6. **[Events](06_events.md)**
7. **[Parameters](07_parameters.md)**
8. **[Clients](08_clients.md)**


