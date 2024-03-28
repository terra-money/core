use cosmwasm_schema::write_api;
use smart_accounts_packages::{
    instantiate_models::InstantiateMsg, query_models::QueryMsg, sudo_models::SudoMsg,
};

fn main() {
    write_api! {
        instantiate: InstantiateMsg,
        sudo: SudoMsg,
        query: QueryMsg,
    }
}
