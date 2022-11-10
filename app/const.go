package app

const (
	// AccountAddressPrefix is the prefix of bech32 encoded address
	AccountAddressPrefix = "terra"

	// AppName is the application name
	AppName = "terra"

	// CoinType is the LUNA coin type as defined in SLIP44 (https://github.com/satoshilabs/slips/blob/master/slip-0044.md)
	CoinType = 330

	// BondDenom staking denom
	BondDenom = "uluna"

	authzMsgExec                        = "/cosmos.authz.v1beta1.MsgExec"
	authzMsgGrant                       = "/cosmos.authz.v1beta1.MsgGrant"
	authzMsgRevoke                      = "/cosmos.authz.v1beta1.MsgRevoke"
	bankMsgSend                         = "/cosmos.bank.v1beta1.MsgSend"
	bankMsgMultiSend                    = "/cosmos.bank.v1beta1.MsgMultiSend"
	distrMsgSetWithdrawAddr             = "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress"
	distrMsgWithdrawValidatorCommission = "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission"
	distrMsgFundCommunityPool           = "/cosmos.distribution.v1beta1.MsgFundCommunityPool"
	distrMsgWithdrawDelegatorReward     = "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"
	feegrantMsgGrantAllowance           = "/cosmos.feegrant.v1beta1.MsgGrantAllowance"
	feegrantMsgRevokeAllowance          = "/cosmos.feegrant.v1beta1.MsgRevokeAllowance"
	govMsgVoteWeighted                  = "/cosmos.gov.v1beta1.MsgVoteWeighted"
	govMsgSubmitProposal                = "/cosmos.gov.v1beta1.MsgSubmitProposal"
	govMsgDeposit                       = "/cosmos.gov.v1beta1.MsgDeposit"
	govMsgVote                          = "/cosmos.gov.v1beta1.MsgVote"
	stakingMsgEditValidator             = "/cosmos.staking.v1beta1.MsgEditValidator"
	stakingMsgDelegate                  = "/cosmos.staking.v1beta1.MsgDelegate"
	stakingMsgUndelegate                = "/cosmos.staking.v1beta1.MsgUndelegate"
	stakingMsgBeginRedelegate           = "/cosmos.staking.v1beta1.MsgBeginRedelegate"
	stakingMsgCreateValidator           = "/cosmos.staking.v1beta1.MsgCreateValidator"
	vestingMsgCreateVestingAccount      = "/cosmos.vesting.v1beta1.MsgCreateVestingAccount"
	transferMsgTransfer                 = "/ibc.applications.transfer.v1.MsgTransfer"
	wasmMsgStoreCode                    = "/cosmwasm.wasm.v1.MsgStoreCode"
	wasmMsgExecuteContract              = "/cosmwasm.wasm.v1.MsgExecuteContract"
	wasmMsgInstantiateContract          = "/cosmwasm.wasm.v1.MsgInstantiateContract"
	wasmMsgMigrateContract              = "/cosmwasm.wasm.v1.MsgMigrateContract"

	// UpgradeName gov proposal name
	UpgradeName = "2.2.0"
)
