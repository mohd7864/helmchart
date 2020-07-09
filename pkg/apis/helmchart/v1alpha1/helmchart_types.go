package v1alpha1

import (
       
       "encoding/json"
       "fmt"
       "strings"
       "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HelmChartSpec defines the desired state of HelmChart
type HelmChartSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
        HelmVersion `json:"helmVersion,omitempty"`
	ChartSource `json:"chart"`
	ReleaseName string `json:"releaseName,omitempty"`
	MaxHistory *int `json:"maxHistory,omitempty"`
	ValueFileSecrets []LocalObjectReference `json:"valueFileSecrets,omitempty"`
	ValuesFrom       []ValuesFromSource     `json:"valuesFrom,omitempty"`
	TargetNamespace string `json:"targetNamespace,omitempty"`
	Timeout *int64 `json:"timeout,omitempty"`
	ResetValues *bool `json:"resetValues,omitempty"`
	SkipCRDs bool `json:"skipCRDs,omitempty"`
	Wait *bool `json:"wait,omitempty"`
	ForceUpgrade bool `json:"forceUpgrade,omitempty"`
	Values HelmValues `json:"values,omitempty"`
}

// HelmChartStatus defines the observed state of HelmChart
type HelmChartStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
        ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	Phase HelmChartPhase `json:"phase,omitempty"`
	ReleaseName string `json:"releaseName,omitempty"`
	ReleaseStatus string `json:"releaseStatus,omitempty"`
	Revision string `json:"revision,omitempty"`
	LastAttemptedRevision string `json:"lastAttemptedRevision,omitempty"`
	Conditions []HelmChartCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HelmChart is the Schema for the helmcharts API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=helmcharts,scope=Namespaced
type HelmChart struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelmChartSpec   `json:"spec,omitempty"`
	Status HelmChartStatus `json:"status,omitempty"`
}

func (hr HelmChart) GetReleaseName() string {
	if hr.Spec.ReleaseName == "" {
		namespace := hr.GetDefaultedNamespace()
		targetNamespace := hr.GetTargetNamespace()

		if namespace != targetNamespace {
			return fmt.Sprintf("%s-%s-%s", namespace, targetNamespace, hr.Name)
		}
		return fmt.Sprintf("%s-%s", targetNamespace, hr.Name)
	}

	return hr.Spec.ReleaseName
}

func (hr HelmChart) GetDefaultedNamespace() string {
	if hr.GetNamespace() == "" {
		return "default"
	}
	return hr.Namespace
}

func (hr HelmChart) GetTargetNamespace() string {
	if hr.Spec.TargetNamespace == "" {
		return hr.GetDefaultedNamespace()
	}
	return hr.Spec.TargetNamespace
}

func (hr HelmChart) GetHelmVersion(defaultVersion string) string {
	if hr.Spec.HelmVersion != "" {
		return string(hr.Spec.HelmVersion)
	}
	if defaultVersion != "" {
		return defaultVersion
	}
	return string(HelmV2)
}

func (hr HelmChart) GetTimeout() time.Duration {
	if hr.Spec.Timeout == nil {
		return 300 * time.Second
	}
	return time.Duration(*hr.Spec.Timeout) * time.Second
}

func (hr HelmChart) GetMaxHistory() int {
	if hr.Spec.MaxHistory == nil {
		return 10
	}
	return *hr.Spec.MaxHistory
}

func (hr HelmChart) GetReuseValues() bool {
	switch hr.Spec.ResetValues {
	case nil:
		return false
	default:
		return !*hr.Spec.ResetValues
	}
}

func (hr HelmChart) GetWait() bool {
	switch hr.Spec.Wait {
	case nil:
		return hr.Spec.Rollback.Enable || hr.Spec.Test.Enable
	default:
		return *hr.Spec.Wait
	}
}

func (hr HelmChart) GetValuesFromSources() []ValuesFromSource {
	valuesFrom := hr.Spec.ValuesFrom
	if hr.Spec.ValueFileSecrets != nil {
		var secretKeyRefs []ValuesFromSource
		for _, ref := range hr.Spec.ValueFileSecrets {
			s := &OptionalSecretKeySelector{}
			s.Name = ref.Name
			secretKeyRefs = append(secretKeyRefs, ValuesFromSource{SecretKeyRef: s})
		}
		valuesFrom = append(secretKeyRefs, valuesFrom...)
	}
	return valuesFrom
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HelmChartList contains a list of HelmChart
type HelmChartList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelmChart `json:"items"`
}

type ChartSource struct {
	// +optional
	*RepoChartSource `json:",inline"`
}

type RepoChartSource struct {
	RepoURL string `json:"repository"`
	Name string `json:"name"`
	Version string `json:"version"`
}

func (s RepoChartSource) CleanRepoURL() string {
	cleanURL := strings.TrimRight(s.RepoURL, "/")
	return cleanURL + "/"
}

type ValuesFromSource struct {
	ConfigMapKeyRef *OptionalConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
	SecretKeyRef *OptionalSecretKeySelector `json:"secretKeyRef,omitempty"`
	ExternalSourceRef *ExternalSourceSelector `json:"externalSourceRef,omitempty"`
	ChartFileRef *ChartFileSelector `json:"chartFileRef,omitempty"`
}

type ChartFileSelector struct {
	Path string `json:"path"`
	Optional *bool `json:"optional,omitempty"`
}

type ExternalSourceSelector struct {
	URL string `json:"url"`
	Optional *bool `json:"optional,omitempty"`
}

type HelmVersion string

const (
	HelmV2 HelmVersion = "v2"
	HelmV3 HelmVersion = "v3"
)

type HelmValues struct {

	Data map[string]interface{} `json:"-"`
}

func (v HelmValues) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Data)
}

func (v *HelmValues) UnmarshalJSON(data []byte) error {
	var out map[string]interface{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return err
	}
	v.Data = out
	return nil
}

func (in *HelmValues) DeepCopyInto(out *HelmValues) {
	b, err := json.Marshal(in.Data)
	if err != nil {
		panic(err)
	}
	var c map[string]interface{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(err)
	}
	out.Data = c
	return
}

type HelmChartConditionType string

const (
	HelmChartChartFetched HelmChartConditionType = "ChartFetched"
	HelmChartDeployed HelmChartConditionType = "Deployed"
	HelmChartReleased HelmChartConditionType = "Released"
)

type HelmChartCondition struct {
	Type HelmChartConditionType `json:"type"`
	Status ConditionStatus `json:"status"`
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`
	Reason string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

type HelmChartPhase string

const (
	HelmChartPhaseChartFetched HelmChartPhase = "ChartFetched"
	HelmChartPhaseChartFetchFailed HelmChartPhase = "ChartFetchFailed"
	HelmChartPhaseInstalling HelmChartPhase = "Installing"
	HelmChartPhaseMigrating HelmChartPhase = "Migrating"
	HelmChartPhaseUpgrading HelmChartPhase = "Upgrading"
	HelmChartPhaseDeployed HelmChartPhase = "Deployed"
	HelmChartPhaseDeployFailed HelmChartPhase = "DeployFailed"
	HelmChartPhaseSucceeded HelmChartPhase = "Succeeded"
        HelmChartPhaseFailed HelmChartPhase = "Failed"	
)



func init() {
	SchemeBuilder.Register(&HelmChart{}, &HelmChartList{})
}
