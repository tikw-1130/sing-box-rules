package generator

import (
	C "github.com/sagernet/sing-box/constant"
)

const (
	SingBoxVersion = "v1.13.14"

	GeositeUpstreamRepo = "Loyalsoldier/v2ray-rules-dat"
	GeoIPUpstreamRepo   = "Loyalsoldier/geoip"

	GeositeBranch = "rule-set-geosite"
	GeoIPBranch   = "rule-set-geoip"

	GeositeOutputDir = "geosite"
	GeoIPOutputDir   = "geoip"

	ManifestAssetName = "manifest.json"
	ReleaseNotesPath  = "dist/release-notes.md"
	ManifestPath      = "dist/manifest.json"

	RuleSetBinaryVersion = C.RuleSetVersion2
)

type Manifest struct {
	GeneratedAt    string           `json:"generated_at"`
	SingBoxVersion string           `json:"sing_box_version"`
	RuleSetVersion uint8            `json:"rule_set_version"`
	Geosite        DatasetManifest  `json:"geosite"`
	GeoIP          DatasetManifest  `json:"geoip"`
	Artifacts      ArtifactManifest `json:"artifacts"`
}

type DatasetManifest struct {
	Repository string `json:"repository"`
	ReleaseTag string `json:"release_tag"`
	AssetName  string `json:"asset_name"`
	SourceSHA  string `json:"source_sha256"`
	ItemCount  int    `json:"item_count"`
}

type ArtifactManifest struct {
	GeositeRuleSetCount int    `json:"geosite_rule_set_count"`
	GeoIPRuleSetCount   int    `json:"geoip_rule_set_count"`
	GeositeDirectory    string `json:"geosite_directory"`
	GeoIPDirectory      string `json:"geoip_directory"`
	GeositeBranch       string `json:"geosite_branch"`
	GeoIPBranch         string `json:"geoip_branch"`
}

type BuildSummary struct {
	Manifest    Manifest
	ReleaseTag  string
	ReleaseName string
}
