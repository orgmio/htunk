# 配置文件

htunk使用 [SOML](https://github.com/soml-lang/soml) 作为配置文件格式。

使用 SOML 而不是 JSON 或 YAML 是因为 SOML 拥有一样甚至更易读的特性，但是
- 跟YAML比，不用把空格敲报废
- 跟JSON比，不用拿放大镜看我哪里少了个/多了结束
- 跟TOML比，SOML在复杂配置上的直观程度是巨大的提升
于是htunk选择了SOML作为配置文件格式。

## 结构

运行 `htunk run` 时，默认读取当前目录下的 `config.soml`。可以用 `htunk run -c 路径` 指定其他文件。

配置文件的顶层可以直接写键值对。日志配置是挂在 `log` 键下的对象，字段用冒号赋值，每行写一个字段即可。

```soml
log: {}
```

## 日志

HTUNK 已能读取并解析以下日志字段，日志输出仍使用 CLI 的默认设置。

```soml
log: {
  enable: true
  level: 'info'
  output: 'htunk.log'
  timestamp: true
}
```

- `enable`：是否启用日志，默认为 `true`。`true`表示启用，`false`表示禁用。
  如果定义为`false`，则其他选项均是无效的。
- `level`：日志级别，默认为 `'info'`，可选值为`'debug'`、`'info'`、`'warn'`、`'error'`。
- `output`：日志文件的保存路径，留空则输出到标准错误。
- `timestamp`：是否带完整时间戳，默认`false`。
  为`true`时每条日志以`+0800 2026-11-15 19:00:37`这样的时间开头；
  为`false`时，输出到终端的日志带启动以来的秒数（如`INFO[0012]`），输出到文件的日志不带时间。
