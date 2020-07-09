package v1alpha1

// +kubebuilder:validation:Enum="True";"False";"Unknown"
type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

type LocalObjectReference struct {
	Name string `json:"name"`
}

type ObjectReference struct {
	LocalObjectReference `json:",inline"`
	Namespace            string `json:"namespace,omitempty"`
}

type ConfigMapKeySelector struct {
	LocalObjectReference `json:",inline"`
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// +optional
	Key string `json:"key,omitempty"`
}

type OptionalConfigMapKeySelector struct {
	ConfigMapKeySelector `json:",inline"`
	// +optional
	Optional bool `json:"optional,omitempty"`
}

type SecretKeySelector struct {
	LocalObjectReference `json:",inline"`
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// +optional
	Key string `json:"key,omitempty"`
}

type OptionalSecretKeySelector struct {
	SecretKeySelector `json:",inline"`
	// +optional
	Optional bool `json:"optional,omitempty"`
}

