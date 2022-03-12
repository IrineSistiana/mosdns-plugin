# mosdns-plugin

mosdns 的外围插件库。

## 摘要

这里是插件的概述。详细说明见各个插件内的 README.md。

- forward_must_edns0: 这个插件会转发包含 EDNS0 的 UDP 报文，然后过滤掉没有 EDNS0 回应的报文。

## 如何使用

将本项目的 `ext_plugin.go` 放在 mosdns 的 `dispatcher/plugin` 下，然后正常编译 mosdns 即可。
