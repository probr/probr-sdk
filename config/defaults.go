package config

import (
	"os"
	"path/filepath"

	"github.com/citihub/probr-sdk/config/setter"
)

// setEnvOrDefaults will set value from os.Getenv and default to the specified value
func setFromEnvOrDefaults(e *VarOptions) {

	setter.SetVar(&e.Tags, "PROBR_TAGS", "")
	setter.SetVar(&e.AuditEnabled, "PROBR_AUDIT_ENABLED", "true")
	setter.SetVar(&e.OutputType, "PROBR_OUTPUT_TYPE", "IO")
	setter.SetVar(&e.WriteDirectory, "PROBR_WRITE_DIRECTORY", "probr_output")
	setter.SetVar(&e.LogLevel, "PROBR_LOG_LEVEL", "ERROR")
	setter.SetVar(&e.OverwriteHistoricalAudits, "OVERWRITE_AUDITS", "true")
	setter.SetVar(&e.WriteConfig, "PROBR_LOG_CONFIG", "true")
	setter.SetVar(&e.ResultsFormat, "PROBR_RESULTS_FORMAT", "cucumber")

	setter.SetVar(&e.ServicePacks.Kubernetes.KeepPods, "PROBR_KEEP_PODS", "false")
	setter.SetVar(&e.ServicePacks.Kubernetes.KubeConfigPath, "KUBE_CONFIG", getDefaultKubeConfigPath())
	setter.SetVar(&e.ServicePacks.Kubernetes.KubeContext, "KUBE_CONTEXT", "")
	setter.SetVar(&e.ServicePacks.Kubernetes.SystemClusterRoles, "", []string{"system:", "aks", "cluster-admin", "policy-agent"})
	setter.SetVar(&e.ServicePacks.Kubernetes.AuthorisedContainerRegistry, "PROBR_AUTHORISED_REGISTRY", "")
	setter.SetVar(&e.ServicePacks.Kubernetes.UnauthorisedContainerRegistry, "PROBR_UNAUTHORISED_REGISTRY", "")
	setter.SetVar(&e.ServicePacks.Kubernetes.ProbeImage, "PROBR_PROBE_IMAGE", "citihub/probr-probe")
	setter.SetVar(&e.ServicePacks.Kubernetes.ContainerRequiredDropCapabilities, "PROBR_REQUIRED_DROP_CAPABILITIES", []string{"NET_RAW"})
	setter.SetVar(&e.ServicePacks.Kubernetes.ContainerAllowedAddCapabilities, "PROBR_ALLOWED_ADD_CAPABILITIES", []string{""})
	setter.SetVar(&e.ServicePacks.Kubernetes.ApprovedVolumeTypes, "PROBR_APPROVED_VOLUME_TYPES", []string{"configmap", "emptydir", "persistentvolumeclaim"})
	setter.SetVar(&e.ServicePacks.Kubernetes.UnapprovedHostPort, "PROBR_UNAPPROVED_HOSTPORT", "22")
	setter.SetVar(&e.ServicePacks.Kubernetes.SystemNamespace, "PROBR_K8S_SYSTEM_NAMESPACE", "kube-system")
	setter.SetVar(&e.ServicePacks.Kubernetes.DashboardPodNamePrefix, "PROBR_K8S_DASHBOARD_PODNAMEPREFIX", "kubernetes-dashboard")
	setter.SetVar(&e.ServicePacks.Kubernetes.ProbeNamespace, "PROBR_K8S_PROBE_NAMESPACE", "probr-general-test-ns")
	setter.SetVar(&e.ServicePacks.Kubernetes.Azure.DefaultNamespaceAIB, "DEFAULT_NS_AZURE_IDENTITY_BINDING", "probr-aib")
	setter.SetVar(&e.ServicePacks.Kubernetes.Azure.IdentityNamespace, "PROBR_K8S_AZURE_IDENTITY_NAMESPACE", "kube-system")

	setter.SetVar(&e.CloudProviders.Azure.TenantID, "AZURE_TENANT_ID", "")
	setter.SetVar(&e.CloudProviders.Azure.SubscriptionID, "AZURE_SUBSCRIPTION_ID", "")
	setter.SetVar(&e.CloudProviders.Azure.ClientID, "AZURE_CLIENT_ID", "")
	setter.SetVar(&e.CloudProviders.Azure.ClientSecret, "AZURE_CLIENT_SECRET", "")
	setter.SetVar(&e.CloudProviders.Azure.ResourceGroup, "AZURE_RESOURCE_GROUP", "")
	setter.SetVar(&e.CloudProviders.Azure.ResourceLocation, "AZURE_RESOURCE_LOCATION", "")
}

func getDefaultKubeConfigPath() string {
	return filepath.Join(homeDir(), ".kube", "config")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
