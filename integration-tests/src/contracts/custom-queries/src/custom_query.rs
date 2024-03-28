use cosmwasm_schema::{cw_serde, QueryResponses};
use cosmwasm_std::CustomQuery;

#[cw_serde]
#[derive(QueryResponses)]
pub enum CustomQueries {
    #[returns(AllianceResponse)]
    Alliance { denom: String },
    #[returns(DelegationResponse)]
    Delegation { denom: String, validator: String, delegator: String },
    #[returns(DelegationRewardsResponse)]
    DelegationRewards { denom: String, validator: String, delegator: String },

    #[returns(TokenFactoryResponses)]
    Token(TokenQuery)
}

#[cw_serde]
pub enum TokenQuery {
    FullDenom { creator_addr: String, subdenom: String },
    Admin { denom: String },
    Metadata { denom: String },
    DenomsByCreator { creator: String },
    Params {},
}

impl CustomQuery for CustomQueries {}

#[cw_serde]
pub struct RewardWeightRange {
    pub min: String,
    pub max: String,
}

#[cw_serde]
pub struct AllianceResponse {
    pub denom: String,
    pub reward_weight: String,
    pub take_rate: String,
    pub total_tokens: String,
    pub total_validator_shares: String,
    pub reward_start_time: u64,
    pub reward_change_rate: String,
    pub last_reward_change_time: u64,
    pub reward_weight_range: RewardWeightRange,
    pub is_initialized: bool,

}

#[cw_serde]
pub struct DelegationResponse {
    pub delegator: String,
    pub validator: String,
    pub denom: String,
    pub amount: String,
}

#[cw_serde]
pub struct DelegationRewardsResponse {
    pub rewards: Vec<Coin>,
}

#[cw_serde]
pub enum TokenFactoryResponses {
   FullDenomResponse(FullDenomResponse),
   AdminResponse(AdminResponse),
   DenumUnit(DenomUnit),
   MetadataReponse(MetadataResponse),
    DenomsByCreatorResponse(DenomsByCreatorResponse),
    ParamsResponse(ParamsResponse),
}

#[cw_serde]
 pub struct FullDenomResponse {
    pub denom: String,
}

#[cw_serde]
pub struct AdminResponse {
    pub admin: String,
}

#[cw_serde]
pub struct DenomUnit {
    pub denom: String,
    pub exponent: u32,
    pub aliases: Vec<String>,
}

#[cw_serde]
pub struct Metadata {
    pub description: String,
    pub denom_units: Vec<DenomUnit>,
    pub base: String,
    pub display: String,
    pub name: String,
    pub symbol: String,
}

#[cw_serde]
pub struct MetadataResponse {
    pub metadata: Option<Metadata>,
}

#[cw_serde]
pub struct DenomsByCreatorResponse {
    pub denoms: Vec<String>,
}

#[cw_serde]
pub struct Coin {
    pub denom: String,
    pub amount: String,
}

#[cw_serde]
pub struct Params {
    pub denom_creation_fee: Vec<Coin>,
}

#[cw_serde]
pub struct ParamsResponse {
    pub params: Params,
}