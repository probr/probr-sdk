package kubernetesconfig

import (
	"os"
	"path/filepath"

	"github.com/probr/probr-sdk/config/setter"
)

// Kubernetes contains common variables needed when using the Kubernetes provider
type Kubernetes struct {
	KeepPods                 string `yaml:"KeepPods"` // TODO: Change type to bool, this would allow us to remove logic from kubernetes.GetKeepPodsFromConfig()
	KubeConfigPath           string `yaml:"KubeConfig"`
	KubeContext              string `yaml:"KubeContext"`
	AuthorisedContainerImage string `yaml:"AuthorisedContainerImage"`
	ProbeNamespace           string `yaml:"ProbeNamespace"`
}

// SetEnvAndDefaults will set value from os.Getenv and default to the specified value
func (ctx *Kubernetes) SetEnvAndDefaults() {
	setter.SetVar(&ctx.KeepPods, "PROBR_KEEP_PODS", "false")
	setter.SetVar(&ctx.KubeConfigPath, "KUBE_CONFIG", getDefaultKubeConfigPath())
	setter.SetVar(&ctx.KubeContext, "KUBE_CONTEXT", "")
	setter.SetVar(&ctx.AuthorisedContainerImage, "PROBR_AUTHORISED_IMAGE", "")
}

func getDefaultKubeConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kube", "config")
}
