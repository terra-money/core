export const SAFE_BLOCK_INCLUSION_TIME = 1100;
export const blockInclusion = () => new Promise(res => setTimeout(res, SAFE_BLOCK_INCLUSION_TIME));
