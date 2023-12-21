export const SAFE_VOTING_PERIOD_TIME = 2100;
export const SAFE_IBC_TRANSFER = 4100;

export const ibcTransfer = () => new Promise((resolve) => setTimeout(() => resolve(SAFE_IBC_TRANSFER), SAFE_IBC_TRANSFER));
export const votingPeriod = () => new Promise((resolve) => setTimeout(() => resolve(SAFE_VOTING_PERIOD_TIME), SAFE_VOTING_PERIOD_TIME));