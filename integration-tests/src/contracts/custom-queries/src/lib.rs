pub mod custom_query;

use cosmwasm_std::{entry_point, Deps, DepsMut, Env, MessageInfo, Response, StdError, Binary, QueryRequest, to_json_binary, StdResult};
use cosmwasm_schema::{cw_serde};
use cw2::set_contract_version;
use crate::custom_query::{CustomQueries, AllianceResponse, DelegationResponse, DelegationRewardsResponse, TokenQuery, FullDenomResponse, AdminResponse, MetadataResponse, DenomsByCreatorResponse, ParamsResponse};

const CONTRACT_NAME: &str = "crates.io:smart-auth";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[cw_serde]
pub struct InstantiateMsg {}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, StdError> {
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;
    Ok(Response::new().add_attribute("contract", format!("{} {}", CONTRACT_NAME, CONTRACT_VERSION)))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps<CustomQueries>, _env: Env, msg: CustomQueries) -> StdResult<Binary> {
    match msg.clone() {
        CustomQueries::Alliance {..} => {
            let res: AllianceResponse = deps.querier.query(&QueryRequest::Custom(msg))?;
            to_json_binary(&res)
        },
        CustomQueries::Delegation {..} => {
            let res: DelegationResponse = deps.querier.query(&QueryRequest::Custom(msg))?;
            to_json_binary(&res)
        }
        CustomQueries::DelegationRewards {..}=> {
            let res: DelegationRewardsResponse = deps.querier.query(&QueryRequest::Custom(msg))?;
            to_json_binary(&res)
        },
        CustomQueries::Token(TokenQuery) => {
            match TokenQuery {
                TokenQuery::FullDenom { .. } => {
                    let res: FullDenomResponse = deps.querier.query(&QueryRequest::Custom(msg))?;
                    to_json_binary(&res)
                },
                TokenQuery::Admin { .. } => {
                    let res: AdminResponse= deps.querier.query(&QueryRequest::Custom(msg))?;
                    to_json_binary(&res)
                },
                TokenQuery::Metadata { .. } => {
                    let res: MetadataResponse= deps.querier.query(&QueryRequest::Custom(msg))?;
                    to_json_binary(&res)
                }
                TokenQuery::DenomsByCreator { .. } => {
                    let res: DenomsByCreatorResponse = deps.querier.query(&QueryRequest::Custom(msg))?;
                    to_json_binary(&res)
                }
                TokenQuery::Params { .. } => {
                    let res: ParamsResponse = deps.querier.query(&QueryRequest::Custom(msg))?;
                    to_json_binary(&res)
                }
            }
        }
    }
}