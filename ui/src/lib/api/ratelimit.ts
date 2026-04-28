import { requestEnvelope } from './http'
import type { RateLimitStatusData } from './types'

export type NormalizedRateLimitStatusData = Omit<RateLimitStatusData, 'default_config'> & {
  default_config: RateLimitStatusData['default_config'] & { enabled: boolean }
}

export async function getRateLimitStatus(): Promise<NormalizedRateLimitStatusData> {
  const data = await requestEnvelope<RateLimitStatusData | undefined>('/ratelimit', {
    method: 'GET',
  })

  return {
    default_config: {
      enabled: data?.default_config?.enabled ?? true,
      max_concurrency: data?.default_config?.max_concurrency ?? 0,
      max_sandbox: data?.default_config?.max_sandbox ?? 0,
    },
    users: Array.isArray(data?.users) ? data.users : [],
  }
}
