package generator

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sagernet/sing-box/common/geosite"
	"github.com/sagernet/sing-box/option"

	"github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"

	"sing-box-rules/internal/githubapi"
)

const (
	geositeAssetName         = "geosite.dat"
	geositeChecksumAssetName = "geosite.dat.sha256sum"
)

func BuildGeosite(
	ctx context.Context,
	client *githubapi.Client,
	rootDir string,
	release *githubapi.Release,
) (DatasetManifest, int, error) {
	payload, sourceSHA, err := downloadVerifiedReleaseAsset(
		ctx,
		client,
		release,
		geositeAssetName,
		geositeChecksumAssetName,
	)
	if err != nil {
		return DatasetManifest{}, 0, err
	}

	domainMap, err := parseGeosite(payload)
	if err != nil {
		return DatasetManifest{}, 0, err
	}

	if err := writeGeositeOutputs(rootDir, domainMap); err != nil {
		return DatasetManifest{}, 0, err
	}

	return DatasetManifest{
		Repository: GeositeUpstreamRepo,
		ReleaseTag: release.TagName,
		AssetName:  geositeAssetName,
		SourceSHA:  sourceSHA,
		ItemCount:  len(domainMap),
	}, len(domainMap), nil
}

func parseGeosite(payload []byte) (map[string][]geosite.Item, error) {
	var geositeList routercommon.GeoSiteList
	if err := proto.Unmarshal(payload, &geositeList); err != nil {
		return nil, err
	}

	domainMap := make(map[string][]geosite.Item)
	for _, entry := range geositeList.Entry {
		code := strings.ToLower(strings.TrimSpace(entry.CountryCode))
		if code == "" {
			continue
		}

		items := make([]geosite.Item, 0, len(entry.Domain)*2)
		attributes := make(map[string][]*routercommon.Domain)
		for _, domain := range entry.Domain {
			items = append(items, expandGeositeDomain(domain)...)
			for _, attribute := range domain.Attribute {
				key := strings.ToLower(strings.TrimSpace(attribute.Key))
				if key == "" {
					continue
				}
				attributes[key] = append(attributes[key], domain)
			}
		}

		domainMap[code] = dedupeGeositeItems(items)
		for attribute, domains := range attributes {
			attributeItems := make([]geosite.Item, 0, len(domains)*2)
			for _, domain := range domains {
				attributeItems = append(attributeItems, expandGeositeDomain(domain)...)
			}
			domainMap[code+"@"+attribute] = dedupeGeositeItems(attributeItems)
		}
	}

	return domainMap, nil
}

func writeGeositeOutputs(rootDir string, domainMap map[string][]geosite.Item) error {
	outputDir := filepath.Join(rootDir, "dist", GeositeOutputDir)

	if err := ensureDir(outputDir); err != nil {
		return err
	}

	for _, code := range sortedKeys(domainMap) {
		defaultRule := geosite.Compile(domainMap[code])
		rule := option.DefaultHeadlessRule{
			Domain:        defaultRule.Domain,
			DomainSuffix:  defaultRule.DomainSuffix,
			DomainKeyword: defaultRule.DomainKeyword,
			DomainRegex:   defaultRule.DomainRegex,
		}
		ruleSet := newPlainRuleSet(rule)
		baseName := "geosite-" + code
		if err := writeRuleSetFile(filepath.Join(outputDir, baseName+".srs"), ruleSet); err != nil {
			return err
		}
	}

	return nil
}

func expandGeositeDomain(domain *routercommon.Domain) []geosite.Item {
	items := make([]geosite.Item, 0, 2)

	switch domain.Type {
	case routercommon.Domain_Plain:
		items = append(items, geosite.Item{
			Type:  geosite.RuleTypeDomainKeyword,
			Value: domain.Value,
		})
	case routercommon.Domain_Regex:
		items = append(items, geosite.Item{
			Type:  geosite.RuleTypeDomainRegex,
			Value: domain.Value,
		})
	case routercommon.Domain_RootDomain:
		if strings.Contains(domain.Value, ".") {
			items = append(items, geosite.Item{
				Type:  geosite.RuleTypeDomain,
				Value: domain.Value,
			})
		}
		items = append(items, geosite.Item{
			Type:  geosite.RuleTypeDomainSuffix,
			Value: "." + domain.Value,
		})
	case routercommon.Domain_Full:
		items = append(items, geosite.Item{
			Type:  geosite.RuleTypeDomain,
			Value: domain.Value,
		})
	}

	return items
}

func dedupeGeositeItems(items []geosite.Item) []geosite.Item {
	result := make([]geosite.Item, 0, len(items))
	seen := make(map[string]struct{}, len(items))

	for _, item := range items {
		key := fmt.Sprintf("%d:%s", item.Type, item.Value)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, item)
	}

	return result
}
