# sing-box-rules

基于上游规则自动生成 `sing-box` 远程 rule-set 文件。

## 生成内容

- `dist/geosite/*.srs`
- `dist/geoip/*.srs`
- `rule-set-geosite` 分支
- `rule-set-geoip` 分支

数据来源：

- `geosite`: [Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat)
- `geoip`: [Loyalsoldier/geoip](https://github.com/Loyalsoldier/geoip)

## 仓库结构

```text
cmd/build-rules/          构建入口
internal/githubapi/       GitHub release / asset 下载
internal/generator/       geosite / geoip 生成逻辑
scripts/publish-branches.sh
.github/workflows/release.yml
```

## sing-box 远程 rule-set 示例

```json
{
  "route": {
    "rule_set": [
      {
        "tag": "geosite-cn",
        "type": "remote",
        "format": "binary",
        "url": "https://cdn.jsdelivr.net/gh/tikw-1130/sing-box-rules@rule-set-geosite/geosite-cn.srs"
      },
      {
        "tag": "geoip-cn",
        "type": "remote",
        "format": "binary",
        "url": "https://cdn.jsdelivr.net/gh/tikw-1130/sing-box-rules@rule-set-geoip/geoip-cn.srs"
      }
    ]
  }
}
```

## 本地运行

```bash
go mod tidy
go run ./cmd/build-rules
```

构建结果会输出到 `dist/geosite/` 和 `dist/geoip/`。

## 许可提醒

本仓库生成的规则文件基于上游公开数据重新整理并分发，使用前请自行确认上游项目的许可证、免责声明和再分发要求：

- `geosite` 数据来自 [Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat)
- `geoip` 数据来自 [Loyalsoldier/geoip](https://github.com/Loyalsoldier/geoip)

