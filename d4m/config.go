package d4m

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

var settingsPath = "~/Library/Group Containers/group.com.docker/settings.json"

// Settings describes the local docker-for-mac settings.
type Settings struct {
	BuildNumber                       string  `json:"buildNumber,omitempty"`
	ChannelID                         string  `json:"channelID,omitempty"`
	Cpus                              int     `json:"cpus,omitempty"`
	DiskPath                          string  `json:"diskPath,omitempty"`
	DiskSizeMiB                       int     `json:"diskSizeMiB,omitempty"`
	DisplayedWelcomeMessage           bool    `json:"displayedWelcomeMessage"`
	DisplayedWelcomeWhale             bool    `json:"displayedWelcomeWhale"`
	DockerAppLaunchPath               string  `json:"dockerAppLaunchPath,omitempty"`
	KubernetesEnabled                 bool    `json:"kubernetesEnabled"`
	KubernetesInitialInstallPerformed bool    `json:"kubernetesInitialInstallPerformed,omitempty"`
	LinuxDaemonConfigCreationDate     string  `json:"linuxDaemonConfigCreationDate,omitempty"`
	MemoryMiB                         int     `json:"memoryMiB,omitempty"`
	ProxyHTTPMode                     string  `json:"proxyHttpMode,omitempty"`
	SettingsVersion                   float64 `json:"settingsVersion,omitempty"`
	UseCredentialHelper               bool    `json:"useCredentialHelper,omitempty"`
	Version                           string  `json:"version,omitempty"`
}

func expandHome(path string) string {
	u, err := user.Current()
	if err != nil {
		return path
	}
	return strings.Replace(path, "~", u.HomeDir, 1)
}

// Load attempts to load settings from the default location on disk.
func Load() (*Settings, error) {
	f, err := os.Open(expandHome(settingsPath))
	if err != nil {
		return nil, err
	}
	s := &Settings{}
	err = json.NewDecoder(f).Decode(s)
	return s, err
}

func loadAsMap() (map[string]interface{}, error) {
	f, err := os.Open(expandHome(settingsPath))
	if err != nil {
		return nil, err
	}
	s := map[string]interface{}{}
	err = json.NewDecoder(f).Decode(&s)
	return s, err
}

func toMap(s interface{}) map[string]interface{} {
	b, _ := json.Marshal(s)
	d := map[string]interface{}{}
	json.Unmarshal(b, &d)
	return d
}

func merge(dst, src map[string]interface{}) {
	for k, v := range src {
		dst[k] = v
	}
}

func (s *Settings) Write() error {
	d, err := loadAsMap()
	if err != nil {
		return err
	}
	merge(d, toMap(s))
	buf, _ := json.MarshalIndent(d, "", "  ")
	return ioutil.WriteFile(expandHome(settingsPath), buf, 0644)
}
