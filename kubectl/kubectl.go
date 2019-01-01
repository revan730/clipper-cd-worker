package kubectl

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const manifestPathPrefix = "/tmp/clipper-tmp-manifest"

type Kubectl struct {
	configPath string
}

func NewKCtl(configPath string) *Kubectl {
	return &Kubectl{
		configPath: configPath,
	}
}

// CreateDeployment calls kubectl to create deployment from provided manifest
// first return parameter show if operation was successful, second is stdout combined with
// stderr
func (k Kubectl) CreateDeployment(manifest string) (bool, string) {
	// Create temporary manifest file
	manifestPath := manifestPathPrefix + strconv.FormatInt(time.Now().Unix(), 10)
	err := ioutil.WriteFile(manifestPath, []byte(manifest), 0644)
	if err != nil {
		return false, "Failed to write temporary file"
	}
	// execute kubectl
	ok := true
	out, err := exec.Command("kubectl", "create", "-f", manifestPath).CombinedOutput()
	if err != nil {
		ok = false
	}
	// remove temporary manifest file
	os.Remove(manifestPath)
	return ok, string(out)
}
