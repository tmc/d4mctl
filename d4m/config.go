package d4m

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

var settingsPaths = []string{
	"~/Library/Group Containers/group.com.docker/settings-store.json",
	"~/Library/Group Containers/group.com.docker/settings.json",
}

// Settings describes the local docker-for-mac settings.
type Settings struct {
	// New format fields (4.0+)
	AutoStart                         bool   `json:"AutoStart,omitempty"`
	Cpus                              int    `json:"Cpus,omitempty"`
	DisplayedOnboarding               bool   `json:"DisplayedOnboarding,omitempty"`
	DockerAppLaunchPath               string `json:"DockerAppLaunchPath,omitempty"`
	DockerBinInstallPath              string `json:"DockerBinInstallPath,omitempty"`
	EnableDefaultDockerSocket         bool   `json:"EnableDefaultDockerSocket,omitempty"`
	LastContainerdSnapshotterEnable   int64  `json:"LastContainerdSnapshotterEnable,omitempty"`
	LicenseTermsVersion               int    `json:"LicenseTermsVersion,omitempty"`
	RequireVmnetd                     bool   `json:"RequireVmnetd,omitempty"`
	SettingsVersion                   int    `json:"SettingsVersion,omitempty"`
	ShowInstallScreen                 bool   `json:"ShowInstallScreen,omitempty"`
	UseContainerdSnapshotter          bool   `json:"UseContainerdSnapshotter,omitempty"`
	UseVirtualizationFrameworkRosetta bool   `json:"UseVirtualizationFrameworkRosetta,omitempty"`

	// Legacy format fields (pre-4.0)
	BuildNumber                       string  `json:"buildNumber,omitempty"`
	ChannelID                         string  `json:"channelID,omitempty"`
	LegacyCpus                        int     `json:"cpus,omitempty"` // Legacy field
	DiskPath                          string  `json:"diskPath,omitempty"`
	DiskSizeMiB                       int     `json:"diskSizeMiB,omitempty"`
	DisplayedWelcomeMessage           bool    `json:"displayedWelcomeMessage,omitempty"`
	DisplayedWelcomeWhale             bool    `json:"displayedWelcomeWhale,omitempty"`
	LegacyDockerAppLaunchPath         string  `json:"dockerAppLaunchPath,omitempty"` // Legacy field
	KubernetesEnabled                 bool    `json:"kubernetesEnabled,omitempty"`
	KubernetesInitialInstallPerformed bool    `json:"kubernetesInitialInstallPerformed,omitempty"`
	LinuxDaemonConfigCreationDate     string  `json:"linuxDaemonConfigCreationDate,omitempty"`
	MemoryMiB                         int     `json:"memoryMiB,omitempty"`
	ProxyHTTPMode                     string  `json:"proxyHttpMode,omitempty"`
	LegacySettingsVersion             float64 `json:"settingsVersion,omitempty"` // Legacy field
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

// Load attempts to load settings from the default locations on disk.
func Load() (*Settings, error) {
	var lastErr error
	for _, path := range settingsPaths {
		f, err := os.Open(expandHome(path))
		if err != nil {
			lastErr = err
			continue
		}
		defer f.Close()
		s := &Settings{}
		if err := json.NewDecoder(f).Decode(s); err != nil {
			lastErr = err
			continue
		}
		// Handle legacy fields if present
		if s.LegacyCpus != 0 && s.Cpus == 0 {
			s.Cpus = s.LegacyCpus
		}
		if s.LegacyDockerAppLaunchPath != "" && s.DockerAppLaunchPath == "" {
			s.DockerAppLaunchPath = s.LegacyDockerAppLaunchPath
		}
		return s, nil
	}
	return nil, lastErr
}

func loadAsMap() (map[string]interface{}, error) {
	var lastErr error
	for _, path := range settingsPaths {
		f, err := os.Open(expandHome(path))
		if err != nil {
			lastErr = err
			continue
		}
		defer f.Close()
		s := map[string]interface{}{}
		if err := json.NewDecoder(f).Decode(&s); err != nil {
			lastErr = err
			continue
		}
		return s, nil
	}
	return nil, lastErr
}

func toMap(s interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	d := map[string]interface{}{}
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, err
	}
	return d, nil
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
	m, err := toMap(s)
	if err != nil {
		return err
	}
	merge(d, m)
	buf, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(expandHome(settingsPaths[0]), buf, 0644)
}
