export const SAFE_VOTING_PERIOD_TIME = 4100;
export const SAFE_BLOCK_INCLUSION_TIME = 1100;
export const blockInclusion = () => new Promise((resolve) => setTimeout(()=>resolve(SAFE_BLOCK_INCLUSION_TIME), SAFE_BLOCK_INCLUSION_TIME));
export const votingPeriod = () => new Promise((resolve) => setTimeout(()=>resolve(SAFE_VOTING_PERIOD_TIME), SAFE_VOTING_PERIOD_TIME));
