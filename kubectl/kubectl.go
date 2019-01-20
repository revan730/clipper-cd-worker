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

func writeManifestFile(manifest string) (string, error) {
	manifestPath := manifestPathPrefix + strconv.FormatInt(time.Now().Unix(), 10)
	err := ioutil.WriteFile(manifestPath, []byte(manifest), 0644)
	return manifestPath, err
}

// CreateDeployment calls kubectl to create deployment from provided manifest
// first return parameter show if operation was successful, second is stdout combined with
// stderr
func (k Kubectl) CreateDeployment(manifest string) (bool, string) {
	manifestPath, err := writeManifestFile(manifest)
	if err != nil {
		return false, "Failed to write temporary file"
	}
	defer os.Remove(manifestPath)
	ok := true
	out, err := exec.Command(kubectlPath, "create", "-f", manifestPath).CombinedOutput()
	if err != nil {
		ok = false
	}
	return ok, string(out)
}

// DeleteDeployment calls kubectl to remove deployment using provided manifest
func (k Kubectl) DeleteDeployment(manifest string) (bool, string) {
	manifestPath, err := writeManifestFile(manifest)
	if err != nil {
		return false, "Failed to write temporary file"
	}
	defer os.Remove(manifestPath)
	ok := true
	out, err := exec.Command(kubectlPath, "delete", "-f", manifestPath).CombinedOutput()
	if err != nil {
		ok = false
	}
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
