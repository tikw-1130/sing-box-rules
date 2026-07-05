package generator

import (
	"context"
	"net"
	"path/filepath"
	"strings"

	"github.com/oschwald/maxminddb-golang"
	"github.com/sagernet/sing-box/option"

	"sing-box-rules/internal/githubapi"
)

const (
	geoIPAssetName         = "Country.mmdb"
	geoIPChecksumAssetName = "Country.mmdb.sha256sum"
)

type geoIPCountryRecord struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
	RegisteredCountry struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"registered_country"`
	RepresentedCountry struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"represented_country"`
	Continent struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"continent"`
}

func BuildGeoIP(
	ctx context.Context,
	client *githubapi.Client,
	rootDir string,
	release *githubapi.Release,
) (DatasetManifest, int, error) {
	payload, sourceSHA, err := downloadVerifiedReleaseAsset(
		ctx,
		client,
		release,
		geoIPAssetName,
		geoIPChecksumAssetName,
	)
	if err != nil {
		return DatasetManifest{}, 0, err
	}

	countryMap, err := parseGeoIP(payload)
	if err != nil {
		return DatasetManifest{}, 0, err
	}

	if err := writeGeoIPOutputs(rootDir, countryMap); err != nil {
		return DatasetManifest{}, 0, err
	}

	return DatasetManifest{
		Repository: GeoIPUpstreamRepo,
		ReleaseTag: release.TagName,
		AssetName:  geoIPAssetName,
		SourceSHA:  sourceSHA,
		ItemCount:  len(countryMap),
	}, len(countryMap), nil
}

func parseGeoIP(payload []byte) (map[string][]*net.IPNet, error) {
	reader, err := maxminddb.FromBytes(payload)
	if err != nil {
		return nil, err
	}

	countryMap := make(map[string][]*net.IPNet)
	networks := reader.Networks(maxminddb.SkipAliasedNetworks)
	var record geoIPCountryRecord

	for networks.Next() {
		ipNet, err := networks.Network(&record)
		if err != nil {
			return nil, err
		}

		code := selectGeoIPCode(record)
		if code == "" {
			continue
		}

		countryMap[code] = append(countryMap[code], ipNet)
	}

	if err := networks.Err(); err != nil {
		return nil, err
	}

	return countryMap, nil
}

func writeGeoIPOutputs(rootDir string, countryMap map[string][]*net.IPNet) error {
	outputDir := filepath.Join(rootDir, "dist", GeoIPOutputDir)

	if err := ensureDir(outputDir); err != nil {
		return err
	}

	allCodes := sortedKeys(countryMap)
	for _, code := range allCodes {
		rule := option.DefaultHeadlessRule{
			IPCIDR: make([]string, 0, len(countryMap[code])),
		}
		for _, cidr := range countryMap[code] {
			rule.IPCIDR = append(rule.IPCIDR, cidr.String())
		}

		ruleSet := newPlainRuleSet(rule)
		baseName := "geoip-" + code
		if err := writeRuleSetFile(filepath.Join(outputDir, baseName+".srs"), ruleSet); err != nil {
			return err
		}
	}

	return nil
}

func selectGeoIPCode(record geoIPCountryRecord) string {
	switch {
	case record.Country.ISOCode != "":
		return strings.ToLower(record.Country.ISOCode)
	case record.RegisteredCountry.ISOCode != "":
		return strings.ToLower(record.RegisteredCountry.ISOCode)
	case record.RepresentedCountry.ISOCode != "":
		return strings.ToLower(record.RepresentedCountry.ISOCode)
	case record.Continent.Code != "":
		return strings.ToLower(record.Continent.Code)
	default:
		return ""
	}
}
