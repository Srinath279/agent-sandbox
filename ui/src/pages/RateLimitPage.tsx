import { ChangeEvent, useCallback, useEffect, useRef, useState } from 'react'

import { getRateLimitStatus, type NormalizedRateLimitStatusData } from '../lib/api/ratelimit'
import type { RateLimitUserStatus } from '../lib/api/types'

const emptyStatus: NormalizedRateLimitStatusData = {
  default_config: {
    enabled: false,
    max_concurrency: 0,
    max_sandbox: 0,
  },
  users: [],
}

const refreshIntervalOptions = [2000, 5000, 10000, 30000]

function formatLimit(value: number): string {
  return value > 0 ? String(value) : 'Unlimited'
}

function formatUsagePercent(user: RateLimitUserStatus): string {
  if (user.sandbox_max <= 0) {
    return '-'
  }
  return `${user.sandbox_usage_percent}%`
}

function progressClass(percent: number): string {
  if (percent >= 90) {
    return 'progress-error'
  }
  if (percent >= 70) {
    return 'progress-warning'
  }
  return 'progress-info'
}

export default function RateLimitPage() {
  const [status, setStatus] = useState<NormalizedRateLimitStatusData>(emptyStatus)
  const [isLoading, setIsLoading] = useState(false)
  const [loadError, setLoadError] = useState('')
  const [isAutoRefresh, setIsAutoRefresh] = useState(true)
  const [refreshIntervalMs, setRefreshIntervalMs] = useState(5000)
  const requestInFlightRef = useRef(false)

  const loadStatus = useCallback(async () => {
    if (requestInFlightRef.current) {
      return
    }

    requestInFlightRef.current = true
    setIsLoading(true)
    setLoadError('')

    try {
      const data = await getRateLimitStatus()
      setStatus(data)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to load rate limit status'
      setLoadError(message)
      setStatus(emptyStatus)
    } finally {
      requestInFlightRef.current = false
      setIsLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadStatus()
  }, [loadStatus])

  useEffect(() => {
    if (!isAutoRefresh) {
      return
    }

    const timer = window.setInterval(() => {
      void loadStatus()
    }, refreshIntervalMs)

    return () => {
      window.clearInterval(timer)
    }
  }, [isAutoRefresh, refreshIntervalMs, loadStatus])

  const handleRefreshIntervalChange = (event: ChangeEvent<HTMLSelectElement>) => {
    const parsed = Number.parseInt(event.target.value, 10)
    if (!Number.isNaN(parsed) && parsed > 0) {
      setRefreshIntervalMs(parsed)
    }
  }

  return (
    <>
      <header className="card border border-base-300 bg-base-100 shadow-sm">
        <div className="card-body gap-3">
          <div>
            <h2 className="text-2xl font-semibold">Capacity</h2>
            <p className="text-sm text-base-content/70">View sandbox capacity and creation concurrency limits by user.</p>
          </div>
        </div>
      </header>

      {loadError && (
        <section>
          <div className="alert alert-error">
            <span>{loadError}</span>
          </div>
        </section>
      )}

      <section className="grid gap-3 md:grid-cols-3">
        <div className="card border border-base-300 bg-base-100 shadow-sm">
          <div className="card-body gap-2 py-4">
            <h3 className="card-title text-lg">Rate Limit Status</h3>
            <div className="stat p-0">
              <div className={`stat-value text-3xl ${status.default_config.enabled ? 'text-success' : 'text-warning'}`}>
                {status.default_config.enabled ? 'Enabled' : 'Disabled'}
              </div>
              <div className="stat-desc">Sandbox creation rate limiting</div>
            </div>
          </div>
        </div>
        <div className="card border border-base-300 bg-base-100 shadow-sm">
          <div className="card-body gap-2 py-4">
            <h3 className="card-title text-lg">Default Sandbox Capacity</h3>
            <div className="stat p-0">
              <div className="stat-value text-3xl">{formatLimit(status.default_config.max_sandbox)}</div>
              <div className="stat-desc">Maximum existing sandboxes per user</div>
            </div>
          </div>
        </div>
        <div className="card border border-base-300 bg-base-100 shadow-sm">
          <div className="card-body gap-2 py-4">
            <h3 className="card-title text-lg">Default Creation Concurrency</h3>
            <div className="stat p-0">
              <div className="stat-value text-3xl">{formatLimit(status.default_config.max_concurrency)}</div>
              <div className="stat-desc">Maximum in-flight creation requests per user</div>
            </div>
          </div>
        </div>
      </section>

      <section>
        <div className="card border border-base-300 bg-base-100 shadow-sm">
          <div className="card-body gap-4">
            <div className="flex items-center justify-between gap-2">
              <h3 className="card-title text-lg">User Limits</h3>
              <div className="flex flex-wrap items-center gap-2">
                <div className="badge badge-outline">{status.users.length} users</div>
                <label className="label cursor-pointer gap-2 py-0">
                  <span className="label-text text-sm">Auto Refresh</span>
                  <input
                    className="toggle toggle-sm"
                    type="checkbox"
                    checked={isAutoRefresh}
                    onChange={() => {
                      setIsAutoRefresh((prev) => !prev)
                    }}
                  />
                </label>
                <label className="flex items-center gap-2">
                  <span className="text-sm">Interval</span>
                  <select className="select select-sm select-bordered" value={String(refreshIntervalMs)} onChange={handleRefreshIntervalChange} disabled={!isAutoRefresh}>
                    {refreshIntervalOptions.map((interval) => (
                      <option key={interval} value={interval}>
                        {interval / 1000}s
                      </option>
                    ))}
                  </select>
                </label>
                <button
                  className={`btn btn-sm btn-outline ${isLoading ? 'btn-disabled' : ''}`}
                  type="button"
                  onClick={() => {
                    void loadStatus()
                  }}
                  disabled={isLoading}
                >
                  {isLoading ? 'Refreshing...' : 'Refresh'}
                </button>
              </div>
            </div>

            <div className="h-[calc(100vh-31rem)] overflow-auto rounded-box border border-base-300">
              <table className="table table-pin-rows table-zebra">
                <thead>
                  <tr>
                    <th>#</th>
                    <th>User</th>
                    <th className="text-center">Creation Concurrency</th>
                    <th className="text-center">Sandbox Capacity</th>
                    <th className="min-w-48">Capacity Usage</th>
                  </tr>
                </thead>
                <tbody>
                  {status.users.length === 0 ? (
                    <tr>
                      <td className="text-center text-base-content/70" colSpan={5}>
                        {isLoading ? 'Loading rate limit status...' : 'No user rate limit data found.'}
                      </td>
                    </tr>
                  ) : (
                    status.users.map((user, index) => {
                      const usagePercent = Math.max(0, Math.min(user.sandbox_usage_percent, 100))
                      return (
                        <tr key={user.user || `user-${index}`}>
                          <td>{index + 1}</td>
                          <td className="font-medium">{user.user || '-'}</td>
                          <td className="text-center">
                            <div className="badge badge-secondary badge-sm">
                              {user.concurrency_active} / {formatLimit(user.concurrency_max)}
                            </div>
                          </td>
                          <td className="text-center">
                            <div className="badge badge-primary badge-sm">
                              {user.sandbox_current} / {formatLimit(user.sandbox_max)}
                            </div>
                          </td>
                          <td>
                            <div className="flex items-center gap-3">
                              <progress className={`progress ${progressClass(usagePercent)} w-32`} value={usagePercent} max="100" />
                              <span className="text-sm tabular-nums">{formatUsagePercent(user)}</span>
                            </div>
                          </td>
                        </tr>
                      )
                    })
                  )}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </section>
    </>
  )
}
