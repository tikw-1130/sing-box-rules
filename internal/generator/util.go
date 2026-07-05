package generator

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sagernet/sing-box/common/srs"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"

	"sing-box-rules/internal/githubapi"
)

func PrepareDist(rootDir string) error {
	distDir := filepath.Join(rootDir, "dist")
	if err := os.RemoveAll(distDir); err != nil {
		return err
	}

	return os.MkdirAll(distDir, 0o755)
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

func sortedKeys[T any](input map[string]T) []string {
	keys := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func newPlainRuleSet(rule option.DefaultHeadlessRule) option.PlainRuleSet {
	return option.PlainRuleSet{
		Rules: []option.HeadlessRule{
			{
				Type:           C.RuleTypeDefault,
				DefaultOptions: rule,
			},
		},
	}
}

func writeRuleSetFile(srsPath string, ruleSet option.PlainRuleSet) error {
	if err := ensureDir(filepath.Dir(srsPath)); err != nil {
		return err
	}

	file, err := os.Create(srsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return srs.Write(file, ruleSet, RuleSetBinaryVersion)
}

func sha256Hex(payload []byte) string {
	checksum := sha256.Sum256(payload)
	return hex.EncodeToString(checksum[:])
}

func parseSHA256File(content []byte) (string, error) {
	fields := strings.Fields(string(content))
	if len(fields) == 0 {
		return "", fmt.Errorf("empty checksum file")
	}

	value := strings.ToLower(strings.TrimSpace(fields[0]))
	if len(value) != 64 {
		return "", fmt.Errorf("invalid checksum value %q", value)
	}

	return value, nil
}

func downloadVerifiedReleaseAsset(
	ctx context.Context,
	client *githubapi.Client,
	release *githubapi.Release,
	assetName string,
	checksumAssetName string,
) ([]byte, string, error) {
	asset := githubapi.FindAsset(release, assetName)
	if asset == nil {
		return nil, "", fmt.Errorf("asset %q not found in %s", assetName, release.TagName)
	}

	payload, err := client.Download(ctx, asset.BrowserDownloadURL)
	if err != nil {
		return nil, "", err
	}

	actualSHA := sha256Hex(payload)
	if checksumAssetName == "" {
		return payload, actualSHA, nil
	}

	checksumAsset := githubapi.FindAsset(release, checksumAssetName)
	if checksumAsset == nil {
		return nil, "", fmt.Errorf("asset %q not found in %s", checksumAssetName, release.TagName)
	}

	checksumContent, err := client.Download(ctx, checksumAsset.BrowserDownloadURL)
	if err != nil {
		return nil, "", err
	}

	expectedSHA, err := parseSHA256File(checksumContent)
	if err != nil {
		return nil, "", err
	}

	if actualSHA != expectedSHA {
		return nil, "", fmt.Errorf("checksum mismatch for %s: expected %s, got %s", assetName, expectedSHA, actualSHA)
	}

	return payload, actualSHA, nil
}
