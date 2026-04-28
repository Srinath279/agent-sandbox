/*
 * Copyright 2025 The https://github.com/agent-sandbox/agent-sandbox Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/kelseyhightower/envconfig"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

const Version = "0.4.3"

type Resources struct {
	CPU         string `json:"cpu"`
	Memory      string `json:"memory"`
	CPULimit    string `json:"cpuLimit"`
	MemoryLimit string `json:"memoryLimit"`
}

type TemplatePool struct {
	Size       int       `json:"size" required:"false"`
	ReadySize  int       `json:"readySize" required:"false"`
	ProbePort  int       `json:"probePort" required:"false"`
	WarmupCmd  string    `json:"warmupCmd" required:"false"`
	StartupCmd string    `json:"startupCmd" required:"false"`
	Resources  Resources `json:"resources" required:"false"`
}

type Template struct {
	Name           string            `json:"name" required:"false"`
	Pattern        string            `json:"pattern" required:"false"`
	Image          string            `json:"image" required:"false"`
	Port           int               `json:"port" required:"false"`
	Type           string            `json:"type" required:"false" description:"dynamic or static, default is static, dynamic means template is dynamic by regexp"`
	Metadata       map[string]string `json:"metadata" required:"false"`
	NoStartupProbe bool              `json:"noStartupProbe" required:"false"`
	Args           []string          `json:"args" required:"false"`
	Resources      Resources         `json:"resources"  required:"false"`
	Pool           TemplatePool      `json:"pool" required:"false"`
	Description    string            `json:"description" required:"false"`
}

type RateLimitConfig struct {
	Enabled        bool `json:"enabled" split_words:"true" default:"true"`
	MaxConcurrency int  `json:"max_concurrency" split_words:"true" default:"10"`
	MaxSandbox     int  `json:"max_sandbox" split_words:"true" default:"100"`
}

func (c *RateLimitConfig) UnmarshalJSON(data []byte) error {
	type rateLimitConfig RateLimitConfig
	var raw struct {
		rateLimitConfig
		EnabledLegacy        *bool `json:"Enabled"`
		MaxConcurrencyLegacy *int  `json:"MaxConcurrency"`
		MaxSandboxLegacy     *int  `json:"MaxSandbox"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*c = RateLimitConfig(raw.rateLimitConfig)
	if raw.EnabledLegacy != nil {
		c.Enabled = *raw.EnabledLegacy
	}
	if raw.MaxConcurrencyLegacy != nil {
		c.MaxConcurrency = *raw.MaxConcurrencyLegacy
	}
	if raw.MaxSandboxLegacy != nil {
		c.MaxSandbox = *raw.MaxSandboxLegacy
	}
	return nil
}

func (c *RateLimitConfig) Validate() {
	if c.MaxConcurrency < 0 {
		klog.Warningf("MaxConcurrency invalid (%d), using default 10", c.MaxConcurrency)
		c.MaxConcurrency = 10
	}
	if c.MaxSandbox < 0 {
		klog.Warningf("MaxSandbox invalid (%d), using default 100", c.MaxSandbox)
		c.MaxSandbox = 100
	}
}

type UserRateLimitConfig struct {
	User           string `json:"user"`
	MaxConcurrency int    `json:"max_concurrency"`
	MaxSandbox     int    `json:"max_sandbox"`
}

type RuntimeConfig struct {
	SystemToken            *string          `json:"system_token,omitempty"`
	APITokensRaw           *string          `json:"api_tokens_raw,omitempty"`
	RateLimit              *RateLimitConfig `json:"rate_limit,omitempty"`
	RateLimitUsersRaw      *string          `json:"rate_limit_users_raw,omitempty"`
	SandboxDefaultImage    *string          `json:"sandbox_default_image,omitempty"`
	SandboxDefaultTemplate *string          `json:"sandbox_default_template,omitempty"`
}

func (c *UserRateLimitConfig) Validate() bool {
	if c.User == "" {
		return false
	}
	if c.MaxConcurrency < 0 {
		klog.Warningf("UserRateLimitConfig.MaxConcurrency negative for user=%s, will use default", c.User)
		c.MaxConcurrency = 0
	}
	if c.MaxSandbox < 0 {
		klog.Warningf("UserRateLimitConfig.MaxSandbox negative for user=%s, will use default", c.User)
		c.MaxSandbox = 0
	}
	return true
}

var Cfg *Config
var Templates []*Template
var SandboxDeployTemplate string

type Config struct {
	KubeClient kubernetes.Interface `ignored:"true"`

	APIVersion  string   `split_words:"true" default:"v1" required:"false"`
	APIBaseURL  string   `split_words:"true" default:"" required:"false"`
	ServerAddr  string   `split_words:"true" default:"0.0.0.0:10000" required:"false"`
	SystemToken string   `split_words:"true" default:"sys-2492a85b10ed4cb083b2c76b181eac96" required:"false"`
	APITokens   []string `ignored:"true"`

	RateLimit      RateLimitConfig       `split_words:"true"`
	RateLimitUsers []UserRateLimitConfig `json:"rateLimitUsers"`

	// witch Kubernetes namespace to create sandboxes Replicaset&Pod in
	SandboxNamespace string `split_words:"true" default:"default" required:"false"`

	SandboxTemplateFile string `split_words:"true" default:"config/sandbox.yaml" required:"false"`

	SandboxTemplatesConfigFile string `split_words:"true" default:"config/templates.json" required:"false"`
	SandboxDefaultImage        string `split_words:"true" default:"ghcr.io/agent-infra/sandbox:latest" required:"false"`
	SandboxDefaultTemplate     string `split_words:"true" default:"aio" required:"false"`

	ConfigmapName string `split_words:"true" default:"agent-sandbox" required:"false"`

	EnvName string `split_words:"true" default:"dev" required:"false"`

	APITokensRaw      string `split_words:"true" default:"testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f" required:"false"`
	RateLimitUsersRaw string `split_words:"true" required:"false"`

	// Runtime config snapshots are read concurrently by request handlers.
	runtimeMu           sync.RWMutex `ignored:"true"`
	apiTokensValue      atomic.Value `ignored:"true"`
	rateLimitValue      atomic.Value `ignored:"true"`
	rateLimitUsersValue atomic.Value `ignored:"true"`
}

func InitConfig() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		klog.Fatal("Failed to process config: ", err)
	}

	cfg.APIBaseURL = "/api/" + cfg.APIVersion

	Cfg = &cfg

	Cfg.ApplyRuntimeConfig(RuntimeConfig{})

	return Cfg
}

func (c *Config) ApplyRuntimeConfig(runtimeConfig RuntimeConfig) {
	if c == nil {
		return
	}

	c.runtimeMu.Lock()
	defer c.runtimeMu.Unlock()

	if runtimeConfig.SystemToken != nil {
		c.SystemToken = strings.TrimSpace(*runtimeConfig.SystemToken)
	}
	if runtimeConfig.APITokensRaw != nil {
		c.APITokensRaw = strings.TrimSpace(*runtimeConfig.APITokensRaw)
	}
	if runtimeConfig.RateLimit != nil {
		c.RateLimit = *runtimeConfig.RateLimit
	}
	if runtimeConfig.RateLimitUsersRaw != nil {
		c.RateLimitUsersRaw = strings.TrimSpace(*runtimeConfig.RateLimitUsersRaw)
	}
	if runtimeConfig.SandboxDefaultImage != nil {
		c.SandboxDefaultImage = strings.TrimSpace(*runtimeConfig.SandboxDefaultImage)
	}
	if runtimeConfig.SandboxDefaultTemplate != nil {
		c.SandboxDefaultTemplate = strings.TrimSpace(*runtimeConfig.SandboxDefaultTemplate)
	}

	c.RateLimit.Validate()
	c.APITokens = ParseAPITokens(c.SystemToken, c.APITokensRaw)
	c.RateLimitUsers = parseRateLimitUsers(c.RateLimitUsersRaw)
	c.apiTokensValue.Store(c.APITokens)
	c.rateLimitValue.Store(c.RateLimit)
	c.rateLimitUsersValue.Store(c.RateLimitUsers)
}

func (c *Config) ApplyRuntimeConfigContent(content string) {
	if c == nil || strings.TrimSpace(content) == "" {
		return
	}

	var runtimeConfig RuntimeConfig
	if err := json.Unmarshal([]byte(content), &runtimeConfig); err != nil {
		klog.Warningf("Failed to parse runtime config, ignoring value: %v", err)
		return
	}

	c.ApplyRuntimeConfig(runtimeConfig)
}

func (c *Config) MergedRuntimeConfig(runtimeConfig RuntimeConfig) RuntimeConfig {
	if c == nil {
		return runtimeConfig
	}

	c.runtimeMu.RLock()
	defer c.runtimeMu.RUnlock()

	systemToken := c.SystemToken
	apiTokensRaw := c.APITokensRaw
	rateLimit := c.RateLimit
	rateLimitUsersRaw := c.RateLimitUsersRaw
	sandboxDefaultImage := c.SandboxDefaultImage
	sandboxDefaultTemplate := c.SandboxDefaultTemplate

	if runtimeConfig.SystemToken != nil {
		systemToken = strings.TrimSpace(*runtimeConfig.SystemToken)
	}
	if runtimeConfig.APITokensRaw != nil {
		apiTokensRaw = strings.TrimSpace(*runtimeConfig.APITokensRaw)
	}
	if runtimeConfig.RateLimit != nil {
		rateLimit = *runtimeConfig.RateLimit
	}
	if runtimeConfig.RateLimitUsersRaw != nil {
		rateLimitUsersRaw = strings.TrimSpace(*runtimeConfig.RateLimitUsersRaw)
	}
	if runtimeConfig.SandboxDefaultImage != nil {
		sandboxDefaultImage = strings.TrimSpace(*runtimeConfig.SandboxDefaultImage)
	}
	if runtimeConfig.SandboxDefaultTemplate != nil {
		sandboxDefaultTemplate = strings.TrimSpace(*runtimeConfig.SandboxDefaultTemplate)
	}

	rateLimit.Validate()

	return RuntimeConfig{
		SystemToken:            &systemToken,
		APITokensRaw:           &apiTokensRaw,
		RateLimit:              &rateLimit,
		RateLimitUsersRaw:      &rateLimitUsersRaw,
		SandboxDefaultImage:    &sandboxDefaultImage,
		SandboxDefaultTemplate: &sandboxDefaultTemplate,
	}
}

func RuntimeConfigContent(runtimeConfig RuntimeConfig) (string, error) {
	content, err := json.MarshalIndent(runtimeConfig, "", "  ")
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (c *Config) RuntimeConfigContent() (string, error) {
	if c == nil {
		return "", fmt.Errorf("config is nil")
	}

	return RuntimeConfigContent(c.MergedRuntimeConfig(RuntimeConfig{}))
}

func (c *Config) RuntimeConfigSnapshot() RuntimeConfig {
	return c.MergedRuntimeConfig(RuntimeConfig{})
}

func (c *Config) GetAPITokens() []string {
	if c == nil {
		return nil
	}
	if value := c.apiTokensValue.Load(); value != nil {
		return value.([]string)
	}
	return c.APITokens
}

func (c *Config) GetRateLimit() RateLimitConfig {
	if c == nil {
		return RateLimitConfig{}
	}
	if value := c.rateLimitValue.Load(); value != nil {
		return value.(RateLimitConfig)
	}
	return c.RateLimit
}

func (c *Config) GetRateLimitUsers() []UserRateLimitConfig {
	if c == nil {
		return nil
	}
	if value := c.rateLimitUsersValue.Load(); value != nil {
		return value.([]UserRateLimitConfig)
	}
	return c.RateLimitUsers
}

func (c *Config) GetSystemToken() string {
	if c == nil {
		return ""
	}
	c.runtimeMu.RLock()
	defer c.runtimeMu.RUnlock()
	return c.SystemToken
}

func (c *Config) GetSandboxDefaultImage() string {
	if c == nil {
		return ""
	}
	c.runtimeMu.RLock()
	defer c.runtimeMu.RUnlock()
	return c.SandboxDefaultImage
}

func (c *Config) GetSandboxDefaultTemplate() string {
	if c == nil {
		return ""
	}
	c.runtimeMu.RLock()
	defer c.runtimeMu.RUnlock()
	return c.SandboxDefaultTemplate
}

func ParseAPITokens(systemToken, apiTokensRaw string) []string {
	raw := strings.TrimSpace(systemToken)
	if strings.TrimSpace(apiTokensRaw) != "" {
		raw = raw + "," + apiTokensRaw
	}

	tokens := strings.Split(raw, ",")
	validTokens := make([]string, 0, len(tokens))
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" || len(token) < 5 {
			continue
		}
		validTokens = append(validTokens, token)
	}
	return validTokens
}

func ParseRateLimitUsers(rateLimitUsersRaw string) ([]UserRateLimitConfig, error) {
	if strings.TrimSpace(rateLimitUsersRaw) == "" {
		return nil, nil
	}

	var rateLimitUsers []UserRateLimitConfig
	if err := json.Unmarshal([]byte(rateLimitUsersRaw), &rateLimitUsers); err != nil {
		return nil, err
	}

	validUsers := make([]UserRateLimitConfig, 0, len(rateLimitUsers))
	for i := range rateLimitUsers {
		uc := &rateLimitUsers[i]
		if uc.Validate() {
			validUsers = append(validUsers, *uc)
		}
	}
	return validUsers, nil
}

func parseRateLimitUsers(rateLimitUsersRaw string) []UserRateLimitConfig {
	rateLimitUsers, err := ParseRateLimitUsers(rateLimitUsersRaw)
	if err != nil {
		klog.Warningf("Failed to parse RATE_LIMIT_USERS, ignoring value: %v", err)
		return nil
	}
	return rateLimitUsers
}

func (c *Config) CheckAndSaveConfigToConfigmap() {
	templatesContent, err := c.ReadTemplatesFromCM()
	if templatesContent == "" {
		klog.Info("templates config is empty, will load from local file system")

		fileName := c.SandboxTemplatesConfigFile
		content, readErr := os.ReadFile(fileName)
		if readErr != nil {
			klog.Fatalf("Failed to read Template config file %s error: %v", fileName, readErr)
		}

		templatesContent = string(content)
		klog.Infof("Loaded Template config from file %s", fileName)

		if err = c.SaveTemplatesToCM(templatesContent); err != nil {
			klog.Fatalf("Failed to save Templates config to configmap, error: %v", err)
		}
		klog.Info("Templates config saved to configmap successfully")
	} else {
		klog.Info("templates config already exists in configmap")
	}

	sandboxTemplateContent, err := c.ReadSandboxTemplateFromCM()
	if err != nil {
		klog.Fatalf("Failed to read sandbox template from configmap: %v", err)
	}
	if sandboxTemplateContent == "" {
		klog.Info("sandbox template config is empty, will load from local file system")

		fileName := c.SandboxTemplateFile
		content, readErr := os.ReadFile(fileName)
		if readErr != nil {
			klog.Fatalf("Failed to read sandbox template file %s error: %v", fileName, readErr)
		}

		sandboxTemplateContent = string(content)
		klog.Infof("Loaded sandbox template config from file %s", fileName)

		if err = c.SaveSandboxTemplateToCM(sandboxTemplateContent); err != nil {
			klog.Fatalf("Failed to save sandbox template config to configmap, error: %v", err)
		}
		klog.Info("Sandbox template config saved to configmap successfully")
	} else {
		klog.Info("sandbox template config already exists in configmap")
	}

	runtimeConfigContent, err := c.ReadRuntimeConfigFromCM()
	if err != nil {
		klog.Fatalf("Failed to read runtime config from configmap: %v", err)
	}
	if runtimeConfigContent == "" {
		klog.Info("runtime config is empty, will save current environment config")
		runtimeConfigContent, err = c.RuntimeConfigContent()
		if err != nil {
			klog.Fatalf("Failed to marshal runtime config: %v", err)
		}
		if err = c.SaveRuntimeConfigToCM(runtimeConfigContent); err != nil {
			klog.Fatalf("Failed to save runtime config to configmap, error: %v", err)
		}
		klog.Info("Runtime config saved to configmap successfully")
	} else {
		klog.Info("runtime config already exists in configmap")
		c.ApplyRuntimeConfigContent(runtimeConfigContent)
	}

}

// ShouldLoadTemplates load templates config
func (c *Config) ShouldLoadTemplates(templatesContent string) {
	var tpls []*Template
	err := json.Unmarshal([]byte(templatesContent), &tpls)
	if err != nil {
		klog.Errorf("Failed to unmarshal Template config templatesContent %s error: %v", templatesContent, err)
	}

	//check envs not empty
	if len(tpls) == 0 {
		klog.Errorf("No Templates  found in config content %s", templatesContent)
	}

	//varify each env has name  image and description
	for _, env := range tpls {
		if env.Name == "" || env.Image == "" || env.Description == "" {
			klog.Errorf("Invalid Template config in templatesContent %s: %+v, name image and desc must not dempty", templatesContent, env)
		}
	}

	Templates = tpls

	//log loaded envs
	for _, env := range Templates {
		klog.Infof("Loaded Template object: %+v", env)
	}
}

func (c *Config) SaveTemplatesToCM(templatesContent string) error {
	cmClient := c.KubeClient.CoreV1().ConfigMaps(c.SandboxNamespace)

	existCm, err := cmClient.Get(context.TODO(), Cfg.ConfigmapName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      Cfg.ConfigmapName,
					Namespace: c.SandboxNamespace,
				},
				Data: map[string]string{
					TemplatesConfigMapKey: templatesContent,
				},
			}
			_, createErr := cmClient.Create(context.TODO(), cm, metav1.CreateOptions{})
			return createErr
		}

		return err
	}

	if existCm.Data == nil {
		existCm.Data = map[string]string{}
	}
	existCm.Data[TemplatesConfigMapKey] = templatesContent
	_, err = cmClient.Update(context.TODO(), existCm, metav1.UpdateOptions{})
	return err
}

func (c *Config) ReadTemplatesFromCM() (content string, err error) {
	templatesContent := ""

	existCm, err := c.KubeClient.CoreV1().ConfigMaps(c.SandboxNamespace).Get(context.TODO(), Cfg.ConfigmapName, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		klog.Info("templates configmap not found, return empty content")
		return templatesContent, nil
	} else if err != nil {
		klog.Errorf("Failed to get ConfigMap for Templates config: %v", err)
		return templatesContent, err
	}

	return existCm.Data[TemplatesConfigMapKey], nil
}

func (c *Config) SaveSandboxTemplateToCM(content string) error {
	cmClient := c.KubeClient.CoreV1().ConfigMaps(c.SandboxNamespace)

	existCm, err := cmClient.Get(context.TODO(), Cfg.ConfigmapName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      Cfg.ConfigmapName,
					Namespace: c.SandboxNamespace,
				},
				Data: map[string]string{
					SandboxTemplateConfigMapKey: content,
				},
			}
			_, createErr := cmClient.Create(context.TODO(), cm, metav1.CreateOptions{})
			return createErr
		}

		return err
	}

	if existCm.Data == nil {
		existCm.Data = map[string]string{}
	}
	existCm.Data[SandboxTemplateConfigMapKey] = content
	_, err = cmClient.Update(context.TODO(), existCm, metav1.UpdateOptions{})
	return err
}

func (c *Config) ReadSandboxTemplateFromCM() (content string, err error) {
	sandboxTemplateContent := ""

	existCm, err := c.KubeClient.CoreV1().ConfigMaps(c.SandboxNamespace).Get(context.TODO(), Cfg.ConfigmapName, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		klog.Info("sandbox template configmap not found, return empty content")
		return sandboxTemplateContent, nil
	} else if err != nil {
		klog.Errorf("Failed to get ConfigMap for sandbox template config: %v", err)
		return sandboxTemplateContent, err
	}

	return existCm.Data[SandboxTemplateConfigMapKey], nil
}

func (c *Config) ReadRuntimeConfigFromCM() (content string, err error) {
	runtimeConfigContent := ""

	existCm, err := c.KubeClient.CoreV1().ConfigMaps(c.SandboxNamespace).Get(context.TODO(), Cfg.ConfigmapName, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		klog.Info("runtime configmap not found, return empty content")
		return runtimeConfigContent, nil
	} else if err != nil {
		klog.Errorf("Failed to get ConfigMap for runtime config: %v", err)
		return runtimeConfigContent, err
	}

	return existCm.Data[RuntimeConfigMapKey], nil
}

func (c *Config) SaveRuntimeConfigToCM(content string) error {
	cmClient := c.KubeClient.CoreV1().ConfigMaps(c.SandboxNamespace)

	existCm, err := cmClient.Get(context.TODO(), Cfg.ConfigmapName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      Cfg.ConfigmapName,
					Namespace: c.SandboxNamespace,
				},
				Data: map[string]string{
					RuntimeConfigMapKey: content,
				},
			}
			_, createErr := cmClient.Create(context.TODO(), cm, metav1.CreateOptions{})
			return createErr
		}

		return err
	}

	if existCm.Data == nil {
		existCm.Data = map[string]string{}
	}
	existCm.Data[RuntimeConfigMapKey] = content
	_, err = cmClient.Update(context.TODO(), existCm, metav1.UpdateOptions{})
	return err
}

func GetTemplateByName(name string) (*Template, error) {
	for _, t := range Templates {
		if t.Name == name {
			return t, nil
		}
	}

	for _, t := range Templates {
		if t.Type == "dynamic" {
			image := t.Image
			//match by regexp
			re := regexp.MustCompile(t.Pattern)
			match := re.FindStringSubmatch(name)
			if len(match) == 0 {
				continue
			}

			if len(match) > 0 {
				versionIndex := re.SubexpIndex("version")
				nameIndex := re.SubexpIndex("name")
				if nameIndex == -1 || versionIndex == -1 {
					continue
				}

				tversion := match[versionIndex]
				tname := match[nameIndex]
				image = strings.ReplaceAll(image, "<name>", tname)
				image = strings.ReplaceAll(image, "<version>", tversion)
			}

			dynT := &Template{
				Name:           t.Name,
				Image:          image,
				Port:           t.Port,
				Pattern:        t.Pattern,
				Pool:           t.Pool,
				Type:           t.Type,
				NoStartupProbe: t.NoStartupProbe,
				Description:    t.Description,
			}
			return dynT, nil
		}
	}

	klog.Errorf("Template %s not found", name)
	return nil, fmt.Errorf("Template  %s not found", name)
}

// GetTemplatesForMCPTools return json string, but exclude image field
func GetTemplatesForMCPTools() string {
	type TplForTool struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var tpls []TplForTool
	for _, env := range Templates {
		if env.Type == "dynamic" {
			continue
		}
		tpls = append(tpls, TplForTool{
			Name:        env.Name,
			Description: env.Description,
		})
	}

	tplsJson, err := json.MarshalIndent(tpls, "", "  ")
	if err != nil {
		klog.Errorf("Failed to marshal Templates for MCP tools: %v", err)
		return ""
	}

	return string(tplsJson)
}
