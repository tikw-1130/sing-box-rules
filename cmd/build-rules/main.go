package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sing-box-rules/internal/generator"
	"sing-box-rules/internal/githubapi"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}

	token := firstNonEmpty(os.Getenv("GITHUB_TOKEN"), os.Getenv("ACCESS_TOKEN"))
	client := githubapi.New(token)

	geositeRelease, err := client.LatestRelease(ctx, generator.GeositeUpstreamRepo)
	if err != nil {
		return fmt.Errorf("fetch geosite release: %w", err)
	}

	geoIPRelease, err := client.LatestRelease(ctx, generator.GeoIPUpstreamRepo)
	if err != nil {
		return fmt.Errorf("fetch geoip release: %w", err)
	}

	force := strings.EqualFold(os.Getenv("FORCE"), "true")
	changed, err := shouldBuild(ctx, client, os.Getenv("GITHUB_REPOSITORY"), geositeRelease.TagName, geoIPRelease.TagName, force)
	if err != nil {
		return err
	}

	if err := setGitHubOutput("changed", boolString(changed)); err != nil {
		return err
	}
	if err := setGitHubOutput("geosite_tag", geositeRelease.TagName); err != nil {
		return err
	}
	if err := setGitHubOutput("geoip_tag", geoIPRelease.TagName); err != nil {
		return err
	}

	releaseTag := fmt.Sprintf("geosite-%s_geoip-%s", geositeRelease.TagName, geoIPRelease.TagName)
	releaseName := fmt.Sprintf("geosite %s | geoip %s", geositeRelease.TagName, geoIPRelease.TagName)
	if err := setGitHubOutput("release_tag", releaseTag); err != nil {
		return err
	}
	if err := setGitHubOutput("release_name", releaseName); err != nil {
		return err
	}

	if !changed {
		log.Println("upstream tags unchanged, skip build")
		return nil
	}

	if err := generator.PrepareDist(rootDir); err != nil {
		return err
	}

	geositeManifest, geositeCount, err := generator.BuildGeosite(ctx, client, rootDir, geositeRelease)
	if err != nil {
		return err
	}

	geoIPManifest, geoIPCount, err := generator.BuildGeoIP(ctx, client, rootDir, geoIPRelease)
	if err != nil {
		return err
	}

	manifest := generator.Manifest{
		GeneratedAt:    time.Now().UTC().Format(time.RFC3339),
		SingBoxVersion: generator.SingBoxVersion,
		RuleSetVersion: generator.RuleSetBinaryVersion,
		Geosite:        geositeManifest,
		GeoIP:          geoIPManifest,
		Artifacts: generator.ArtifactManifest{
			GeositeRuleSetCount: geositeCount,
			GeoIPRuleSetCount:   geoIPCount,
			GeositeBranch:       generator.GeositeBranch,
			GeoIPBranch:         generator.GeoIPBranch,
		},
	}

	if err := writeJSON(filepath.Join(rootDir, generator.ManifestPath), manifest); err != nil {
		return err
	}

	if err := writeText(filepath.Join(rootDir, generator.ReleaseNotesPath), releaseNotes(manifest)); err != nil {
		return err
	}

	return nil
}

func shouldBuild(
	ctx context.Context,
	client *githubapi.Client,
	repository string,
	geositeTag string,
	geoIPTag string,
	force bool,
) (bool, error) {
	if force || strings.TrimSpace(repository) == "" {
		return true, nil
	}

	release, err := client.LatestRelease(ctx, repository)
	if err != nil {
		if errors.Is(err, githubapi.ErrNotFound) {
			return true, nil
		}
		return false, fmt.Errorf("fetch current repository release: %w", err)
	}

	manifestAsset := githubapi.FindAsset(release, generator.ManifestAssetName)
	if manifestAsset == nil {
		return true, nil
	}

	manifestBytes, err := client.Download(ctx, manifestAsset.BrowserDownloadURL)
	if err != nil {
		return false, fmt.Errorf("download manifest asset: %w", err)
	}

	var manifest generator.Manifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return true, nil
	}

	return manifest.Geosite.ReleaseTag != geositeTag || manifest.GeoIP.ReleaseTag != geoIPTag, nil
}

func setGitHubOutput(name string, value string) error {
	outputPath := strings.TrimSpace(os.Getenv("GITHUB_OUTPUT"))
	if outputPath == "" {
		return nil
	}

	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s=%s\n", name, value)
	return err
}

func writeJSON(path string, value any) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	return encoder.Encode(value)
}

func writeText(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}

func releaseNotes(manifest generator.Manifest) string {
	return fmt.Sprintf(`# sing-box rules

- sing-box version: %s
- rule-set binary version: %d
- geosite upstream: %s @ %s
- geoip upstream: %s @ %s
- geosite categories: %d
- geoip categories: %d

Assets:

- geosite-rule-set.zip
- geoip-rule-set.zip
- manifest.json

Branches:

- %s
- %s
`,
		manifest.SingBoxVersion,
		manifest.RuleSetVersion,
		manifest.Geosite.Repository,
		manifest.Geosite.ReleaseTag,
		manifest.GeoIP.Repository,
		manifest.GeoIP.ReleaseTag,
		manifest.Artifacts.GeositeRuleSetCount,
		manifest.Artifacts.GeoIPRuleSetCount,
		manifest.Artifacts.GeositeBranch,
		manifest.Artifacts.GeoIPBranch,
	)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}

func boolString(value bool) string {
	if value {
		return "true"
	}

	return "false"
}
