import { requestEnvelope } from './http'
import type { RuntimeConfigPayload, RuntimeConfigStatus, Template } from './types'

export async function getRuntimeConfigStatus(): Promise<RuntimeConfigStatus> {
  return requestEnvelope<RuntimeConfigStatus>('/config/runtime', {
    method: 'GET',
  })
}

export async function saveRuntimeConfig(payload: RuntimeConfigPayload): Promise<RuntimeConfigStatus> {
  return requestEnvelope<RuntimeConfigStatus>('/config/runtime', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  })
}

export async function getTemplatesConfig(): Promise<string> {
  return requestEnvelope<string>('/config/templates', {
    method: 'GET',
  })
}

export async function saveTemplatesConfig(payload: Template[]): Promise<string> {
  return requestEnvelope<string>('/config/templates', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  })
}

export async function getSandboxTemplateConfig(): Promise<string> {
  return requestEnvelope<string>('/config/sandbox-template', {
    method: 'GET',
  })
}

export async function saveSandboxTemplateConfig(payload: string): Promise<string> {
  return requestEnvelope<string>('/config/sandbox-template', {
    method: 'POST',
    headers: {
      'Content-Type': 'text/plain; charset=utf-8',
    },
    body: payload,
  })
}
