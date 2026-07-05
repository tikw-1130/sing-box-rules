# sing-box-rules-starter

一个可直接上传到 GitHub 的 `sing-box` 规则仓库模板。

它会基于当前上游自动生成：

- `rule-set-geosite` 分支中的 `.srs`
- `rule-set-geoip` 分支中的 `.srs`
- Release 里的 `geosite-rule-set.zip`
- Release 里的 `geoip-rule-set.zip`

数据来源：

- `geosite`: [Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat)
- `geoip`: [Loyalsoldier/geoip](https://github.com/Loyalsoldier/geoip)

说明：

- 现代 `sing-box` 更推荐使用远程 `rule-set`。
- 这个模板已经按你的需求改成 `srs-only`，不再生成 `.db`。
- 当前模板按 `sing-box v1.13.14` 的 API 重写，`Go` 工作流使用 `1.26.4`。

## 仓库结构

```text
cmd/build-rules/          构建入口
internal/githubapi/       GitHub release / asset 下载
internal/generator/       geosite / geoip 生成逻辑
scripts/publish-branches.sh
.github/workflows/release.yml
```

## 上传到 GitHub

1. 新建一个空仓库。
2. 把这个目录里的文件全部推上去。
3. 到仓库 `Settings -> Actions -> General` 确认工作流具备写入仓库内容的权限。
4. 在 `Actions` 页面手动运行一次 `Release sing-box rules`。

第一次运行成功后，你会得到：

- Release 资产里的 `*.zip`
- `rule-set-geosite` 分支
- `rule-set-geoip` 分支

## 需要改的地方

当前模板已经按仓库 `tikw-1130/sing-box-rules` 填好了示例地址。

## 下载地址示例

把下面占位符替换成你的仓库名以后就能直接用。

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

## 构建行为

工作流会先读取你自己仓库的最新 release 中的 `manifest.json`，再和上游最新 tag 比较：

- 如果上游 tag 没变化，默认跳过发版
- 如果上游 tag 有变化，才会重新生成并发布
- 你也可以在 `workflow_dispatch` 里把 `force` 设为 `true` 强制发版

## 本地运行

```bash
go mod tidy
go run ./cmd/build-rules
```

构建结果会输出到 `dist/`。

## 许可提醒

这个模板会消费并重新分发上游规则数据。你正式公开前，建议你根据上游仓库的许可证要求补充自己的 `LICENSE` 和说明文件。
