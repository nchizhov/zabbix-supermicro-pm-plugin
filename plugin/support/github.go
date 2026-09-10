package support

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type ReleaseInfo struct {
	Name string
	Url  string
	MD5  string
}

func (releaseInfo *ReleaseInfo) fillReleaseInfo(data []byte, prefix string) error {
	var dataInfo map[string]any
	if err := json.Unmarshal(data, &dataInfo); err != nil {
		return err
	}

	releaseInfo.Name = getRequiredReleaseName(dataInfo["tag_name"].(string), prefix)
	assets, ok := dataInfo["assets"]
	if !ok {
		return errors.New("assets in release not found")
	}
	return releaseInfo.parseReleaseInfo(assets.([]any))
}

func (releaseInfo *ReleaseInfo) parseReleaseInfo(assetInfo []any) error {
	hasMD5 := false
	hasRelease := false
	for _, baseAsset := range assetInfo {
		asset := baseAsset.(map[string]any)
		assetName := asset["name"].(string)
		if !strings.HasPrefix(assetName, releaseInfo.Name) {
			continue
		}
		url := asset["browser_download_url"].(string)
		if strings.HasSuffix(assetName, ".md5") {
			md5Hash, err := downloadFile(url)
			if err != nil {
				return err
			}
			releaseInfo.MD5 = strings.TrimSpace(string(md5Hash))
			hasMD5 = true
			continue
		}
		releaseInfo.Url = url
		hasRelease = true
	}
	if !hasRelease {
		return errors.New("release file not exists")
	}
	if !hasMD5 {
		return errors.New("md5 checksum not exists")
	}
	return nil
}

func getRequiredReleaseName(tag string, prefix string) string {
	postfix := ""
	if runtime.GOOS == "windows" {
		postfix = ".exe"
	}
	return fmt.Sprintf("%s-%s-%s-%s%s", prefix, tag, runtime.GOOS, runtime.GOARCH, postfix)
}
