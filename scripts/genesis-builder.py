import sys
import json
import argparse


def init_default_argument_parser(prog_desc, default_chain_id,
                                 default_genesis_time, default_pretty):
    # Load genesis template via argument
    parser = argparse.ArgumentParser(description=prog_desc)

    parser.add_argument(
        'genesis-template',
        help='template genesis.json file',
        type=argparse.FileType('r'), default=sys.stdin,
    )
    parser.add_argument('--chain-id', type=str,
                        help='chain id for new genesis', default=default_chain_id)
    parser.add_argument('--genesis-time', type=str,
                        help='genesis time for new genesis', default=default_genesis_time)
    parser.add_argument('--pretty', type=bool,
                        default=default_pretty)
    return parser


def main(argument_parser, process_genesis_func):
    args = argument_parser.parse_args()
    if args.chain_id.strip() == '':
        sys.exit('chain-id required')

    genesis = json.loads(args.exported_genesis.read())
    genesis = process_genesis_func(genesis=genesis, parsed_args=args,)

    if args.pretty:
        raw_genesis = json.dumps(genesis, indent=4, sort_keys=True)
    else:
        raw_genesis = json.dumps(genesis, indent=None,
                                 sort_keys=False, separators=(',', ':'))

    print(raw_genesis)


def process_raw_genesis(genesis, parsed_args):
    # Consensus Params: Block
    genesis['consensus_params']['block'] = {
        'max_bytes': '1000000',
        'max_gas': '100000000',
        'time_iota_ms': '1000',
    }

    # Mint: target inflation 7%
    genesis['app_state']['mint'] = {
        'minter': {
            'inflation': '0.070000000000000000',
            'annual_provisions': '0.000000000000000000'
        },
        'params': {
            'mint_denom': 'uluna',
            'inflation_rate_change': '0.000000000000000000',
            'inflation_max': '0.070000000000000000',
            'inflation_min': '0.070000000000000000',
            'goal_bonded': '0.670000000000000000',
            'blocks_per_year': '6311520'
        }
    },

    # Staking: change bond_denom to uluna
    genesis['app_state']['staking']['params']['bond_denom'] = 'uluna'

    # Crisis: change constant fee to 512 LUNA
    genesis['app_state']['crisis']['params']['constant_fee'] = {
        'denom': 'uluna',
        'amount': '512000000',
    }

    # Gov: change min deposit to 512 LUNA
    genesis['app_state']['gov']['deposit_params']['min_deposit'] = {
        'denom': 'uluna',
        'amount': '512000000',
    }

    # Account Registration


def add_normal_account(genesis, address, amount):
    genesis['app_state']['auth']['accounts'].append({
        '@type': '/cosmos.auth.v1beta1.BaseAccount',
        'address': address,
        'pub_key': None,
        'account_number': '0',
        'sequence': '0'
    })

    genesis['app_state']['bank']['balances'].append({
        'address': address,
        'coins': [
            {
                'denom': 'uluna',
                'amount': amount
            }
        ]
    })


def add_continuous_vesting_account(genesis, address, total_amount, vesting_amount, start_time, end_time):
    genesis['app_state']['auth']['accounts'].append({
        '@type': '/cosmos.vesting.v1beta1.ContinuousVestingAccount',
        'base_vesting_account': {
            'base_account': {
                'address': address,
                'pub_key': None,
                'account_number': '0',
                'sequence': '0'
            },
            'original_vesting': [
                {
                    'denom': 'uluna',
                    'amount': vesting_amount
                }
            ],
            'delegated_free': [],
            'delegated_vesting': [],
            'end_time': end_time
        },
        'start_time': start_time
    })

    genesis['app_state']['bank']['balances'].append({
        'address': address,
        'coins': [
            {
                'denom': 'uluna',
                'amount': total_amount
            }
        ]
    })


def add_delayed_vesting_account(genesis, address, total_amount, vesting_amount, end_time):
    genesis['app_state']['auth']['accounts'].append({
        '@type': '/cosmos.vesting.v1beta1.DelayedVestingAccount',
        'base_vesting_account': {
            'base_account': {
                'address': address,
                'pub_key': None,
                'account_number': '0',
                'sequence': '0'
            },
            'original_vesting': [
                {
                    'denom': 'uluna',
                    'amount': vesting_amount
                }
            ],
            'delegated_free': [],
            'delegated_vesting': [],
            'end_time': end_time
        }
    })

    genesis['app_state']['bank']['balances'].append({
        'address': address,
        'coins': [
            {
                'denom': 'uluna',
                'amount': total_amount
            }
        ]
    })


if __name__ == '__main__':
    parser = init_default_argument_parser(
        prog_desc='Genesis Builder for Terra Revival',
        default_chain_id='columbus-6',
        default_genesis_time='2022-05-30T00:00:00Z',
        default_pretty=False
    )
    main(parser, process_raw_genesis)

# 02-820-4494 + 4492
