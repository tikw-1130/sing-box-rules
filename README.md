# sing-box-rules

基于上游规则自动生成 `sing-box` 远程 rule-set 文件。

参考仓库：[lyc8503/sing-box-rules](https://github.com/lyc8503/sing-box-rules)

本仓库只生成 `.srs` 二进制规则文件，不生成 `.db` 规则文件，也不生成 `.json` 规则文件。

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

## GitHub Actions

第一次运行前，到仓库 `Settings -> Actions -> General` 确认工作流具备写入仓库内容的权限。

然后在 `Actions` 页面手动运行一次 `Publish sing-box rules`。

运行成功后会得到：

- `rule-set-geosite` 分支，内容来自 `dist/geosite/*.srs`
- `rule-set-geoip` 分支，内容来自 `dist/geoip/*.srs`

工作流不会创建 GitHub Release，也不会上传压缩包。每次运行都会重新生成 `.srs` 并强制刷新上面两个规则分支。

## 下载地址示例

- `https://cdn.jsdelivr.net/gh/tikw-1130/sing-box-rules@rule-set-geosite/geosite-cn.srs`
- `https://cdn.jsdelivr.net/gh/tikw-1130/sing-box-rules@rule-set-geoip/geoip-cn.srs`

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

这个模板会消费并重新分发上游规则数据。正式公开前，建议根据上游仓库的许可证要求补充自己的 `LICENSE` 和说明文件。
