package config

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

const TemplatesConfigMapKey = "config-templates"

// SandboxTemplateConfigMapKey sandbox template is k8s Resource Definition(ReplicasSet) for sandbox
const SandboxTemplateConfigMapKey = "config-sandbox-template"

const RuntimeConfigMapKey = "config-runtime"

func WatchConfigMap() func(configMap *corev1.ConfigMap) {
	var lastTemplatesContent string
	var lastSandboxTemplateContent string
	var lastRuntimeConfigContent string

	return func(configMap *corev1.ConfigMap) {
		templatesContent := configMap.Data[TemplatesConfigMapKey]
		if templatesContent != "" && templatesContent != lastTemplatesContent {
			klog.Info("watching ConfigMap changed, templates content updated")
			Cfg.ShouldLoadTemplates(templatesContent)
			lastTemplatesContent = templatesContent
		}

		sandboxTemplateContent := configMap.Data[SandboxTemplateConfigMapKey]
		if sandboxTemplateContent != "" && sandboxTemplateContent != lastSandboxTemplateContent {
			klog.Info("watching ConfigMap changed, sandbox template content updated")
			SandboxDeployTemplate = sandboxTemplateContent
			lastSandboxTemplateContent = sandboxTemplateContent
		}

		runtimeConfigContent, ok := configMap.Data[RuntimeConfigMapKey]
		if ok && runtimeConfigContent != "" && runtimeConfigContent != lastRuntimeConfigContent {
			klog.Info("watching ConfigMap changed, runtime config updated")
			Cfg.ApplyRuntimeConfigContent(runtimeConfigContent)
			lastRuntimeConfigContent = runtimeConfigContent
		}
	}
}
