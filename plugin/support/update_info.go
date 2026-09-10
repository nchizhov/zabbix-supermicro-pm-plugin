package support

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return readDownloadedData(resp, "incorrect status code for MD5: %d")
}

func CheckUpdate(releaseUrl string, prefix string) (bool, string) {
	releaseInfo, err := getUpdateInfo(releaseUrl, prefix)
	if err != nil {
		return false, ""
	}

	programMD5, err := getProgramChecksum()
	if err != nil {
		return false, ""
	}
	if programMD5 == releaseInfo.MD5 {
		return false, ""
	}

	return true, fmt.Sprintf("has new release update. Url for download: %s", releaseInfo.Url)
}

func getUpdateInfo(url string, prefix string) (ReleaseInfo, error) {
	var releaseInfo ReleaseInfo

	data, err := getGithubData(url)
	if err != nil {
		return releaseInfo, err
	}

	err = releaseInfo.fillReleaseInfo(data, prefix)
	if err != nil {
		return releaseInfo, err
	}
	return releaseInfo, nil
}

func getGithubData(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := readDownloadedData(resp, "incorrect status code for latest release info: %d")
	if err != nil {
		return nil, err
	}
	return body, nil
}

func readDownloadedData(resp *http.Response, errFormat string) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(errFormat, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func getProgramChecksum() (string, error) {
	fl, err := os.ReadFile(os.Args[0])
	if err != nil {
		return "", err
	}

	md5SumTmp := md5.Sum(fl)
	md5Sum := hex.EncodeToString(md5SumTmp[:])
	return md5Sum, nil
}
