package kubernetesconfig

import (
	"os"
	"path/filepath"

	"github.com/probr/probr-sdk/config/setter"
)

// Kubernetes contains common variables needed when using the Kubernetes provider
type Kubernetes struct {
	KeepPods                 string `yaml:"KeepPods"`
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
	setter.SetVar(&ctx.ProbeNamespace, "PROBR_K8S_PROBE_NAMESPACE", "probr-general-test-ns")
}

func getDefaultKubeConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kube", "config")
}
