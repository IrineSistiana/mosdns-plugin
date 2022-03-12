# forward_must_edns0

这个插件会转发包含 EDNS0 的 UDP 报文，然后过滤掉没有 EDNS0 回应的报文。

参数:

```yaml
tag: ''
type: 'forward_must_edns0'
args:
  upstream:
    - addr: 8.8.8.8   # 上游 UDP 地址。
    - addr: 1.1.1.1   # 可配置多个。
```