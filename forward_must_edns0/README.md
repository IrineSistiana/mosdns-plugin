# forward_must_edns0

这个插件会发送包含 EDNS0 的 UDP 报文，然后过滤掉没有 EDNS0 回应的报文。

如果服务器支持 EDNS0，则会回应 EDNS0。过滤掉没有 EDNS0 报文可以过滤掉某些有问题的回应。

参数:

```yaml
tag: ''
type: 'forward_must_edns0'
args:
  upstream:
    - addr: 8.8.8.8   # 上游 UDP 地址。省略端口号会用 53 默认值。
    - addr: 1.1.1.1   # 可配置多个。
```