/*
Copyright 2019 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// StatePoweringOn indicates device is powering on
	StatePoweringOn = "powering_on"

	// StatePoweringOff indicates device is powering off
	StatePoweringOff = "powering_off"

	// StateActive indicates device is in active state
	StateActive = "active"

	// StateInactive indicates device is in inactive state
	StateInactive = "inactive"

	// StateProvisioning indicates device is in provisioning state
	StateProvisioning = "provisioning"

	// StateDeprovisioning indicates device is in deprovisioning state
	StateDeprovisioning = "deprovisioning"

	// StateReinstalling indicates device is in reinstalling state
	StateReinstalling = "reinstalling"

	// StateFailed indicates device is in a failed state
	StateFailed = "failed"

	// StateQueued indicates device is in a queued state
	StateQueued = "queued"
)

// TODO: make optional parameters pointers and add +optional

// DeviceSpec defines the desired state of Device
type DeviceSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       DeviceParameters `json:"forProvider"`
}

// DeviceStatus defines the observed state of Device
type DeviceStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          DeviceObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// Device is a managed resource that represents an Equinix Metal Device
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="STATE",type="string",JSONPath=".status.atProvider.state"
// +kubebuilder:printcolumn:name="ID",type="string",JSONPath=".status.atProvider.id"
// +kubebuilder:printcolumn:name="HOSTNAME",type="string",JSONPath=".spec.forProvider.hostname"
// +kubebuilder:printcolumn:name="FACILITY",type="string",JSONPath=".status.atProvider.facility"
// +kubebuilder:printcolumn:name="IPV4",type="string",JSONPath=".status.atProvider.ipv4"
// +kubebuilder:printcolumn:name="RECLAIM-POLICY",type="string",JSONPath=".spec.reclaimPolicy"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,equinix}
type Device struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeviceSpec   `json:"spec,omitempty"`
	Status DeviceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DeviceList contains a list of Devices
type DeviceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Device `json:"items"`
}

// IPAddress is a packngo.IPAddressCreateRequest used for managing IP addresses
// at Device, at creation and observer time.
type IPAddress struct {
	AddressFamily int      `json:"address_family"`
	Public        bool     `json:"public"`
	CIDR          int      `json:"cidr,omitempty"`
	Reservations  []string `json:"ip_reservations,omitempty"`
}

// NamespacedName represents a namespaced object name
type NamespacedName struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// DataKeySelector defines required spec to access a key of a configmap or secret
type DataKeySelector struct {
	NamespacedName `json:",inline,omitempty"`
	// +kubebuilder:validation:Enum=Secret;ConfigMap
	Kind     string `json:"kind"`
	Key      string `json:"key,omitempty"`
	Optional bool   `json:"optional,omitempty"`
}

// DeviceParameters define the desired state of an Equinix Metal device.
// https://metal.equinix.com/developers/api/#devices
//
// Reference values are used for optional parameters to determine if
// LateInitialization should update the parameter after creation.
type DeviceParameters struct {
	// +immutable
	// +required
	Plan string `json:"plan"`

	// +immutable
	Facility string `json:"facility,omitempty"`

	// +immutable
	Metro string `json:"metro,omitempty"`

	// +immutable
	// +required
	OS string `json:"operatingSystem"`

	// +optional
	Hostname *string `json:"hostname,omitempty"`

	// +optional
	Description *string `json:"description,omitempty"`

	// +optional
	BillingCycle *string `json:"billingCycle,omitempty"`

	// +optional
	UserData *string `json:"userdata,omitempty"`

	// +optional
	UserDataRef *DataKeySelector `json:"userdataRef,omitempty"`

	// +optional
	Tags []string `json:"tags,omitempty"`

	// +optional
	Locked *bool `json:"locked,omitempty"`

	// +optional
	IPXEScriptURL *string `json:"ipxeScriptUrl,omitempty"`

	// +immutable
	// +optional
	PublicIPv4SubnetSize *int `json:"publicIPv4SubnetSize,omitempty"`

	// +optional
	AlwaysPXE *bool `json:"alwaysPXE,omitempty"`

	// +immutable
	// +optional
	HardwareReservationID *string `json:"hardwareReservationID,omitempty"`

	// +optional
	CustomData *string `json:"customData,omitempty"`

	// +immutable
	// +optional
	UserSSHKeys []string `json:"userSSHKeys,omitempty"`

	// +immutable
	// +optional
	ProjectSSHKeys []string `json:"projectSSHKeys,omitempty"`

	// +optional
	// +kubebuilder:validation:Enum="hybrid";"layer2-individual";"layer2-bonded";"layer3"
	NetworkType *string `json:"networkType,omitempty"`

	// Features can be used to require or prefer devices with optional features:
	//
	// features:
	// - tpm: required
	// - tpm: preferred
	// +immutable
	// +optional
	Features map[string]string `json:"features,omitempty"`

	// IPAddresses will be attached to the device. These addresses can be drawn
	// from existing reservations.
	//
	// +immutable
	// +optional
	IPAddresses []IPAddress `json:"ipAddresses,omitempty"`
}

// DeviceObservation is used to reflect in the Kubernetes API, the observed
// state of the Device resource from the Equinix Metal API.
type DeviceObservation struct {
	ID   string `json:"id"`
	Href string `json:"href,omitempty"`

	// Facility is where the device is deployed. This field may differ from
	// spec.forProvider.facility when the "any" value was used.
	Facility            string            `json:"facility"`
	Metro               string            `json:"metro,omitempty"`
	State               string            `json:"state,omitempty"`
	ProvisionPercentage resource.Quantity `json:"provisionPercentage,omitempty"`
	IPv4                string            `json:"ipv4,omitempty"`
	Locked              bool              `json:"locked"`

	// +optional
	CreatedAt *metav1.Time `json:"createdAt,omitempty"`

	// +optional
	UpdatedAt *metav1.Time `json:"updatedAt,omitempty"`

	// IQN string is omitted
	// ImageURL *string is omitted
	// Hostname string is omitted (represented in ForProvider)
	// Tags []string is omitted (represented in ForProvider)
	// BillingCycle string is omitted (represented in ForProvider)
	// HardwareReservation map is omitted (represented in ForProvider by HardwareReservationID)
	// IPAddresses []map is omitted
	// NetworkPorts []map is omitted
	// OperatingSystem map is omitted
	// Plan map is omitted (represented in ForProvider by Plan)
	// Project map is omitted (represented through ProviderReference)
	// ShortID string is omitted
	// SSHKeys []map is omitted
	// Volumes []map is omitted

	// User string is omitted (written to Credentials)
	// RootPassword string is omitted (written to Credentials)
}
