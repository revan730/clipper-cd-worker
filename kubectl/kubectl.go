package kubectl

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const manifestPathPrefix = "/tmp/clipper-tmp-manifest"
const kubectlPath = "kubectl"

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
	out, err := exec.Command(kubectlPath, "create", "-f", manifestPath).CombinedOutput()
	if err != nil {
		ok = false
	}
	// remove temporary manifest file
	os.Remove(manifestPath)
	return ok, string(out)
}

// ChangeImage calls kubectl to change deployment's container image
func (k Kubectl) ChangeImage(deployment, imageURL string) (bool, string) {
	ok := true
	depParameter := fmt.Sprintf("deployment/%s", deployment)
	container := fmt.Sprintf("%s=%s", deployment, imageURL)
	out, err := exec.Command(kubectlPath, "set", "image", depParameter, container).CombinedOutput()
	if err != nil {
		ok = false
	}
	return ok, string(out)
}

// ScaleDeployment calls kubectl to scale deployment to provided size
func (k Kubectl) ScaleDeployment(deployment string, replicas int64) (bool, string) {
	ok := true
	depParameter := fmt.Sprintf("deployment/%s", deployment)
	replicasParameter := fmt.Sprintf("--replicas=%d", replicas)
	out, err := exec.Command(kubectlPath, "scale", depParameter, replicasParameter).CombinedOutput()
	if err != nil {
		ok = false
	}
	return ok, string(out)
}
