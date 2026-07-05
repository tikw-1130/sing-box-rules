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
- 生成逻辑参考了 [lyc8503/sing-box-rules](https://github.com/lyc8503/sing-box-rules)

建议在正式公开或长期分发前补充仓库自己的 `LICENSE`，并在说明文件中保留上游来源链接。若上游许可证、数据来源或分发规则发生变化，应以对应上游仓库的最新说明为准。

本项目仅提供规则转换与自动发布脚本，不对上游数据的完整性、准确性、可用性或合规性作保证。使用者应自行承担使用、同步、再分发这些规则文件所产生的风险和责任。

