export const SAFE_BLOCK_INCLUSION_TIME = 1100;
export const blockInclusion = () => new Promise((resolve) => setTimeout(()=>resolve(SAFE_BLOCK_INCLUSION_TIME), SAFE_BLOCK_INCLUSION_TIME));
