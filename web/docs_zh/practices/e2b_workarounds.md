---
icon: lucide/wrench
---

# E2B 兼容性解决方案

由于 E2B SDK 需要 **HTTPS** 和 **通配符域名** 支持沙箱路由，如果你的环境不支持这些特性，可以使用以下解决方案适配你的环境。

E2B 默认沙箱 URL 格式如下：  
`https://{port}-{sandboxID}.your-domain.com`   

例如：  
`https://6080-294bef011f1e4567b4c5d02593e2e90e.example.com`

‼️ 如果你没有 `https` 或通配符域名（*.example.com）支持，请配置并使用以下 `no_wildcard()` 函数适配你的环境。

```python
import os

def no_wildcard():
    os.environ['E2B_API_URL'] = 'http://localhost:10000/e2b/v1'
    os.environ['E2B_API_KEY'] = 'testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f'
    os.environ['E2B_DOMAIN'] = 'localhost:10000'

    def __connection_config_get_host(_, sandbox_id: str, sandbox_domain: str, port: int) -> str:
        return f"{sandbox_domain}/sandboxes/router/{sandbox_id}/{port}"
    from e2b import ConnectionConfig
    ConnectionConfig.get_host = __connection_config_get_host
```
### 无 HTTPS
`os.environ['E2B_DEBUG'] = "true"` 会对 API 请求使用 HTTP，适用于本地开发环境没有 HTTPS 的情况。

## Kubernetes Ingress 通配符域名配置
你可以设置以下环境变量来配置 E2B SDK：

```python
import os
os.environ['E2B_DOMAIN'] = 'example.domain.com'
os.environ['E2B_API_URL'] = 'https://example.domain.com/e2b/v1'
os.environ['E2B_API_KEY'] = 'testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f'
```

Agent-Sandbox 的 Ingress 也需要配置支持通配符域名。

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: agent-sandbox
  annotations:
    nginx.ingress.kubernetes.io/proxy-body-size: 1024M
    nginx.ingress.kubernetes.io/proxy-connect-timeout: "300"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "300"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "300"
spec:
  ingressClassName: ingress-controller
  rules:
    - host: "*.example.domain.com"
      http:
        paths:
          - backend:
              service:
                name: agent-sandbox
                port:
                  number: 80
            path: /
            pathType: ImplementationSpecific
```


## 更多信息

- [E2B 官方文档](https://docs.e2b.dev/)
- [e2b-code-interpreter SDK](https://github.com/e2b-dev/e2b-code-interpreter)
- [E2B Desktop SDK](https://github.com/e2b-dev/desktop)