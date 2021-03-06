/*
Copyright 2018 BlackRock, Inc.

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

package v1alpha1

import (
	"fmt"
	"hash/fnv"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"

	apicommon "github.com/argoproj/argo-events/pkg/apis/common"
)

// NotificationType represent a type of notifications that are handled by a sensor
type NotificationType string

const (
	// EventNotification is a notification for an event dependency (from a Gateway)
	EventNotification NotificationType = "Event"
	// ResourceUpdateNotification is a notification that an associated resource was updated
	ResourceUpdateNotification NotificationType = "ResourceUpdate"
)

// NodeType is the type of a node
type NodeType string

const (
	// NodeTypeEventDependency is a node that represents a single event dependency
	NodeTypeEventDependency NodeType = "EventDependency"
	// NodeTypeTrigger is a node that represents a single trigger
	NodeTypeTrigger NodeType = "Trigger"
	// NodeTypeDependencyGroup is a node that represents a group of event dependencies
	NodeTypeDependencyGroup NodeType = "DependencyGroup"
)

// NodePhase is the label for the condition of a node
type NodePhase string

// possible types of node phases
const (
	NodePhaseComplete NodePhase = "Complete" // the node has finished successfully
	NodePhaseActive   NodePhase = "Active"   // the node is active and waiting on dependencies to resolve
	NodePhaseError    NodePhase = "Error"    // the node has encountered an error in processing
	NodePhaseNew      NodePhase = ""         // the node is new
)

// TriggerCycleState is the label for the state of the trigger cycle
type TriggerCycleState string

// possible values of trigger cycle states
const (
	TriggerCycleSuccess TriggerCycleState = "Success" // all triggers are successfully executed
	TriggerCycleFailure TriggerCycleState = "Failure" // one or more triggers failed
)

// KubernetesResourceOperation refers to the type of operation performed on the K8s resource
type KubernetesResourceOperation string

// possible values for KubernetesResourceOperation
const (
	// deprecate create.
	Create KubernetesResourceOperation = "create" // create the resource
	Update KubernetesResourceOperation = "update" // updates the resource
	Patch  KubernetesResourceOperation = "patch"  // patch resource
)

// ArgoWorkflowOperation refers to the type of the operation performed on the Argo Workflow
type ArgoWorkflowOperation string

// possible values for ArgoWorkflowOperation
const (
	Submit   ArgoWorkflowOperation = "submit"   // submit a workflow
	Suspend  ArgoWorkflowOperation = "suspend"  // suspends a workflow
	Resubmit ArgoWorkflowOperation = "resubmit" // resubmit a workflow
	Retry    ArgoWorkflowOperation = "retry"    // retry a workflow
	Resume   ArgoWorkflowOperation = "resume"   // resume a workflow
)

// Comparator refers to the comparator operator for a data filter
type Comparator string

const (
	GreaterThanOrEqualTo Comparator = ">=" // Greater than or equal to value provided in data filter
	GreaterThan          Comparator = ">"  // Greater than value provided in data filter
	EqualTo              Comparator = "="  // Equal to value provided in data filter
	LessThan             Comparator = "<"  // Less than value provided in data filter
	LessThanOrEqualTo    Comparator = "<=" // Less than or equal to value provided in data filter
	EmptyComparator                 = ""   // Equal to value provided in data filter
)

// Sensor is the definition of a sensor resource
// +genclient
// +genclient:noStatus
// +kubebuilder:resource:shortName=sn
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Sensor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Spec              SensorSpec   `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	Status            SensorStatus `json:"status" protobuf:"bytes,3,opt,name=status"`
}

// SensorList is the list of Sensor resources
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SensorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Items []Sensor `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// SensorSpec represents desired sensor state
type SensorSpec struct {

	// Dependencies is a list of the events that this sensor is dependent on.
	Dependencies []EventDependency `json:"dependencies" protobuf:"bytes,1,rep,name=dependencies"`

	// Triggers is a list of the things that this sensor evokes. These are the outputs from this sensor.
	Triggers []Trigger `json:"triggers" protobuf:"bytes,2,rep,name=triggers"`
	// Template is the pod specification for the sensor
	// +optional
	Template Template `json:"template,omitempty" protobuf:"bytes,3,opt,name=template"`
	// Subscription refers to the modes of events subscriptions for the sensor.
	// At least one of the types of subscription must be defined in order for sensor to be meaningful.
	Subscription *Subscription `json:"subscription,omitempty" protobuf:"bytes,4,opt,name=subscription"`
	// Circuit is a boolean expression of dependency groups
	Circuit string `json:"circuit,omitempty" protobuf:"bytes,5,opt,name=circuit"`

	// DependencyGroups is a list of the groups of events.
	DependencyGroups []DependencyGroup `json:"dependencyGroups,omitempty" protobuf:"bytes,6,rep,name=dependencyGroups"`
	// ErrorOnFailedRound if set to true, marks sensor state as `error` if the previous trigger round fails.
	// Once sensor state is set to `error`, no further triggers will be processed.
	ErrorOnFailedRound bool `json:"errorOnFailedRound,omitempty" protobuf:"varint,7,opt,name=errorOnFailedRound"`
	// ServiceLabels to be set for the service generated
	ServiceLabels map[string]string `json:"serviceLabels,omitempty" protobuf:"bytes,8,rep,name=serviceLabels"`
	// ServiceAnnotations refers to annotations to be set
	// for the service generated
	ServiceAnnotations map[string]string `json:"serviceAnnotations,omitempty" protobuf:"bytes,9,rep,name=serviceAnnotations"`
}

// Template holds the information of a sensor deployment template
type Template struct {
	// ServiceAccountName is the name of the ServiceAccount to use to run gateway pod.
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty" protobuf:"bytes,1,opt,name=serviceAccountName"`
	// Container is the main container image to run in the gateway pod
	// +optional
	Container *corev1.Container `json:"container,omitempty" protobuf:"bytes,2,opt,name=container"`
	// Volumes is a list of volumes that can be mounted by containers in a workflow.
	// +patchStrategy=merge
	// +patchMergeKey=name
	// +optional
	Volumes []corev1.Volume `json:"volumes,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,3,rep,name=volumes"`
	// SecurityContext holds pod-level security attributes and common container settings.
	// Optional: Defaults to empty.  See type description for default values of each field.
	// +optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty" protobuf:"bytes,4,opt,name=securityContext"`
	// Spec holds the sensor deployment spec.
	// DEPRECATED: Use Container instead.
	Spec *corev1.PodSpec `json:"spec,omitempty" protobuf:"bytes,5,opt,name=spec"`
}

// Subscription holds different modes of subscription available for sensor to consume events.
type Subscription struct {
	// HTTP refers to the HTTP subscription of events for the sensor.
	// +optional
	HTTP *HTTPSubscription `json:"http,omitempty" protobuf:"bytes,1,opt,name=http"`
	// NATS refers to the NATS subscription of events for the sensor
	// +optional
	NATS *NATSSubscription `json:"nats,omitempty" protobuf:"bytes,2,opt,name=nats"`
}

// HTTPSubscription holds the context of the HTTP subscription of events for the sensor.
type HTTPSubscription struct {
	// Port on which sensor server should run.
	Port int32 `json:"port" protobuf:"varint,1,opt,name=port"`
}

// NATSSubscription holds the context of the NATS subscription of events for the sensor
type NATSSubscription struct {
	// ServerURL refers to NATS server url.
	ServerURL string `json:"serverURL" protobuf:"bytes,1,opt,name=serverURL"`
	// Subject refers to NATS subject name.
	Subject string `json:"subject" protobuf:"bytes,2,opt,name=subject"`
}

// EventDependency describes a dependency
type EventDependency struct {
	// Name is a unique name of this dependency
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// GatewayName is the name of the gateway from whom the event is received
	// DEPRECATED: Use EventSourceName instead.
	GatewayName string `json:"gatewayName" protobuf:"bytes,2,name=gatewayName"`
	// EventSourceName is the name of EventSource that Sensor depends on
	EventSourceName string `json:"eventSourceName" protobuf:"bytes,3,name=eventSourceName"`
	// EventName is the name of the event
	EventName string `json:"eventName" protobuf:"bytes,4,name=eventName"`
	// Filters and rules governing toleration of success and constraints on the context and data of an event
	Filters *EventDependencyFilter `json:"filters,omitempty" protobuf:"bytes,5,opt,name=filters"`
}

// DependencyGroup is the group of dependencies
type DependencyGroup struct {
	// Name of the group
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// Dependencies of events
	Dependencies []string `json:"dependencies" protobuf:"bytes,2,rep,name=dependencies"`
}

// EventDependencyFilter defines filters and constraints for a event.
type EventDependencyFilter struct {
	// Name is the name of event filter
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Time filter on the event with escalation
	Time *TimeFilter `json:"time,omitempty" protobuf:"bytes,2,opt,name=time"`
	// Context filter constraints
	Context *EventContext `json:"context,omitempty" protobuf:"bytes,3,opt,name=context"`

	// Data filter constraints with escalation
	Data []DataFilter `json:"data,omitempty" protobuf:"bytes,4,rep,name=data"`
}

// TimeFilter describes a window in time.
// DataFilters out event events that occur outside the time limits.
// In other words, only events that occur after Start and before Stop
// will pass this filter.
type TimeFilter struct {
	// Start is the beginning of a time window.
	// Before this time, events for this event are ignored and
	// format is hh:mm:ss
	Start string `json:"start,omitempty" protobuf:"bytes,1,opt,name=start"`
	// StopPattern is the end of a time window.
	// After this time, events for this event are ignored and
	// format is hh:mm:ss
	Stop string `json:"stop,omitempty" protobuf:"bytes,2,opt,name=stop"`
}

// JSONType contains the supported JSON types for data filtering
type JSONType string

// the various supported JSONTypes
const (
	JSONTypeBool   JSONType = "bool"
	JSONTypeNumber JSONType = "number"
	JSONTypeString JSONType = "string"
)

// DataFilter describes constraints and filters for event data
// Regular Expressions are purposefully not a feature as they are overkill for our uses here
// See Rob Pike's Post: https://commandcenter.blogspot.com/2011/08/regular-expressions-in-lexing-and.html
type DataFilter struct {
	// Path is the JSONPath of the event's (JSON decoded) data key
	// Path is a series of keys separated by a dot. A key may contain wildcard characters '*' and '?'.
	// To access an array value use the index as the key. The dot and wildcard characters can be escaped with '\\'.
	// See https://github.com/tidwall/gjson#path-syntax for more information on how to use this.
	Path string `json:"path" protobuf:"bytes,1,opt,name=path"`
	// Type contains the JSON type of the data
	Type JSONType `json:"type" protobuf:"bytes,2,opt,name=type,casttype=JSONType"`

	// Value is the allowed string values for this key
	// Booleans are passed using strconv.ParseBool()
	// Numbers are parsed using as float64 using strconv.ParseFloat()
	// Strings are taken as is
	// Nils this value is ignored
	Value []string `json:"value" protobuf:"bytes,3,rep,name=value"`
	// Comparator compares the event data with a user given value.
	// Can be ">=", ">", "=", "<", or "<=".
	// Is optional, and if left blank treated as equality "=".
	Comparator Comparator `json:"comparator,omitempty" protobuf:"bytes,4,opt,name=comparator,casttype=Comparator"`
}

// Trigger is an action taken, output produced, an event created, a message sent
type Trigger struct {
	// Template describes the trigger specification.
	Template *TriggerTemplate `json:"template,omitempty" protobuf:"bytes,1,opt,name=template"`

	// Parameters is the list of parameters applied to the trigger template definition
	Parameters []TriggerParameter `json:"parameters,omitempty" protobuf:"bytes,2,rep,name=parameters"`
	// Policy to configure backoff and execution criteria for the trigger
	Policy *TriggerPolicy `json:"policy,omitempty" protobuf:"bytes,3,opt,name=policy"`
}

// TriggerTemplate is the template that describes trigger specification.
type TriggerTemplate struct {
	// Name is a unique name of the action to take.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Switch is the condition to execute the trigger.
	// +optional
	Switch *TriggerSwitch `json:"switch,omitempty" protobuf:"bytes,2,opt,name=switch"`
	// StandardK8STrigger refers to the trigger designed to create or update a generic Kubernetes resource.
	// +optional
	K8s *StandardK8STrigger `json:"k8s,omitempty" protobuf:"bytes,3,opt,name=k8s"`
	// ArgoWorkflow refers to the trigger that can perform various operations on an Argo workflow.
	// +optional
	ArgoWorkflow *ArgoWorkflowTrigger `json:"argoWorkflow,omitempty" protobuf:"bytes,4,opt,name=argoWorkflow"`
	// HTTP refers to the trigger designed to dispatch a HTTP request with on-the-fly constructable payload.
	// +optional
	HTTP *HTTPTrigger `json:"http,omitempty" protobuf:"bytes,5,opt,name=http"`
	// AWSLambda refers to the trigger designed to invoke AWS Lambda function with with on-the-fly constructable payload.
	// +optional
	AWSLambda *AWSLambdaTrigger `json:"awsLambda,omitempty" protobuf:"bytes,6,opt,name=awsLambda"`
	// CustomTrigger refers to the trigger designed to connect to a gRPC trigger server and execute a custom trigger.
	// +optional
	CustomTrigger *CustomTrigger `json:"custom,omitempty" protobuf:"bytes,7,opt,name=custom"`
	// Kafka refers to the trigger designed to place messages on Kafka topic.
	// +optional.
	Kafka *KafkaTrigger `json:"kafka,omitempty" protobuf:"bytes,8,opt,name=kafka"`
	// NATS refers to the trigger designed to place message on NATS subject.
	// +optional.
	NATS *NATSTrigger `json:"nats,omitempty" protobuf:"bytes,9,opt,name=nats"`
	// Slack refers to the trigger designed to send slack notification message.
	// +optional
	Slack *SlackTrigger `json:"slack,omitempty" protobuf:"bytes,10,opt,name=slack"`
	// OpenWhisk refers to the trigger designed to invoke OpenWhisk action.
	// +optional
	OpenWhisk *OpenWhiskTrigger `json:"openWhisk,omitempty" protobuf:"bytes,11,opt,name=openWhisk"`
}

// TriggerSwitch describes condition which must be satisfied in order to execute a trigger.
// Depending upon condition type, status of dependency groups is used to evaluate the result.
type TriggerSwitch struct {

	// Any acts as a OR operator between dependencies
	Any []string `json:"any,omitempty" protobuf:"bytes,1,rep,name=any"`

	// All acts as a AND operator between dependencies
	All []string `json:"all,omitempty" protobuf:"bytes,2,rep,name=all"`
}

// StandardK8STrigger is the standard Kubernetes resource trigger
type StandardK8STrigger struct {
	// The unambiguous kind of this object - used in order to retrieve the appropriate kubernetes api client for this resource
	metav1.GroupVersionResource `json:",inline" protobuf:"bytes,1,opt,name=groupVersionResource"`
	// Source of the K8 resource file(s)
	Source *ArtifactLocation `json:"source,omitempty" protobuf:"bytes,2,opt,name=source"`
	// Operation refers to the type of operation performed on the k8s resource.
	// Default value is Create.
	// +optional
	Operation KubernetesResourceOperation `json:"operation,omitempty" protobuf:"bytes,3,opt,name=operation,casttype=KubernetesResourceOperation"`
	// Parameters is the list of parameters that is applied to resolved K8s trigger object.

	Parameters []TriggerParameter `json:"parameters,omitempty" protobuf:"bytes,4,rep,name=parameters"`
	// PatchStrategy controls the K8s object patching strategy when the trigger operation is specified as patch.
	// possible values:
	// "application/json-patch+json"
	// "application/merge-patch+json"
	// "application/strategic-merge-patch+json"
	// "application/apply-patch+yaml".
	// Defaults to "application/merge-patch+json"
	// +optional
	PatchStrategy k8stypes.PatchType `json:"patchStrategy,omitempty" protobuf:"bytes,5,opt,name=patchStrategy,casttype=k8s.io/apimachinery/pkg/types.PatchType"`
	// LiveObject specifies whether the resource should be directly fetched from K8s instead
	// of being marshaled from the resource artifact. If set to true, the resource artifact
	// must contain the information required to uniquely identify the resource in the cluster,
	// that is, you must specify "apiVersion", "kind" as well as "name" and "namespace" meta
	// data.
	// Only valid for operation type `update`
	// +optional
	LiveObject bool `json:"liveObject,omitempty" protobuf:"varint,6,opt,name=liveObject"`
}

// ArgoWorkflowTrigger is the trigger for the Argo Workflow
type ArgoWorkflowTrigger struct {
	// Source of the K8 resource file(s)
	Source *ArtifactLocation `json:"source,omitempty" protobuf:"bytes,1,opt,name=source"`
	// Operation refers to the type of operation performed on the argo workflow resource.
	// Default value is Submit.
	// +optional
	Operation ArgoWorkflowOperation `json:"operation,omitempty" protobuf:"bytes,2,opt,name=operation,casttype=ArgoWorkflowOperation"`
	// Parameters is the list of parameters to pass to resolved Argo Workflow object

	Parameters []TriggerParameter `json:"parameters,omitempty" protobuf:"bytes,3,rep,name=parameters"`
	// The unambiguous kind of this object - used in order to retrieve the appropriate kubernetes api client for this resource
	metav1.GroupVersionResource `json:",inline" protobuf:"bytes,4,opt,name=groupVersionResource"`
}

// HTTPTrigger is the trigger for the HTTP request
type HTTPTrigger struct {
	// URL refers to the URL to send HTTP request to.
	URL string `json:"url" protobuf:"bytes,1,opt,name=url"`
	// Payload is the list of key-value extracted from an event payload to construct the HTTP request payload.

	Payload []TriggerParameter `json:"payload" protobuf:"bytes,2,rep,name=payload"`
	// TLS configuration for the HTTP client.
	// +optional
	TLS *TLSConfig `json:"tls,omitempty" protobuf:"bytes,3,opt,name=tls"`
	// Method refers to the type of the HTTP request.
	// Refer https://golang.org/src/net/http/method.go for more info.
	// Default value is POST.
	// +optional
	Method string `json:"method,omitempty" protobuf:"bytes,4,opt,name=method"`
	// Parameters is the list of key-value extracted from event's payload that are applied to
	// the HTTP trigger resource.

	Parameters []TriggerParameter `json:"parameters,omitempty" protobuf:"bytes,5,rep,name=parameters"`
	// Timeout refers to the HTTP request timeout in seconds.
	// Default value is 60 seconds.
	// +optional
	Timeout int64 `json:"timeout,omitempty" protobuf:"varint,6,opt,name=timeout"`
	// BasicAuth configuration for the http request.
	// +optional
	BasicAuth *BasicAuth `json:"basicAuth,omitempty" protobuf:"bytes,7,opt,name=basicAuth"`
	// Headers for the HTTP request.
	// +optional
	Headers map[string]string `json:"headers,omitempty" protobuf:"bytes,8,rep,name=headers"`
}

// TLSConfig refers to TLS configuration for the HTTP client
type TLSConfig struct {
	// CACertPath refers the file path that contains the CA cert.
	CACertPath string `json:"caCertPath" protobuf:"bytes,1,opt,name=caCertPath"`
	// ClientCertPath refers the file path that contains client cert.
	ClientCertPath string `json:"clientCertPath" protobuf:"bytes,2,opt,name=clientCertPath"`
	// ClientKeyPath refers the file path that contains client key.
	ClientKeyPath string `json:"clientKeyPath" protobuf:"bytes,3,opt,name=clientKeyPath"`
}

// BasicAuth contains the reference to K8s secrets that holds the username and password
type BasicAuth struct {
	// Username refers to the Kubernetes secret that holds the username required for basic auth.
	Username *corev1.SecretKeySelector `json:"username,omitempty" protobuf:"bytes,1,opt,name=username"`
	// Password refers to the Kubernetes secret that holds the password required for basic auth.
	Password *corev1.SecretKeySelector `json:"password,omitempty" protobuf:"bytes,2,opt,name=password"`
	// Namespace to read the secrets from.
	// Defaults to sensor's namespace.
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
}

// AWSLambdaTrigger refers to specification of the trigger to invoke an AWS Lambda function
type AWSLambdaTrigger struct {
	// FunctionName refers to the name of the function to invoke.
	FunctionName string `json:"functionName" protobuf:"bytes,1,opt,name=functionName"`
	// AccessKey refers K8 secret containing aws access key
	AccessKey *corev1.SecretKeySelector `json:"accessKey,omitempty" protobuf:"bytes,2,opt,name=accessKey"`
	// SecretKey refers K8 secret containing aws secret key
	SecretKey *corev1.SecretKeySelector `json:"secretKey,omitempty" protobuf:"bytes,3,opt,name=secretKey"`
	// Namespace refers to Kubernetes namespace to read access related secret from.
	// Defaults to sensor's namespace.
	// +optional.
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,4,opt,name=namespace"`
	// Region is AWS region
	Region string `json:"region" protobuf:"bytes,5,opt,name=region"`
	// Payload is the list of key-value extracted from an event payload to construct the request payload.

	Payload []TriggerParameter `json:"payload" protobuf:"bytes,6,rep,name=payload"`
	// Parameters is the list of key-value extracted from event's payload that are applied to
	// the trigger resource.

	// +optional
	Parameters []TriggerParameter `json:"parameters,omitempty" protobuf:"bytes,7,rep,name=parameters"`
}

// KafkaTrigger refers to the specification of the Kafka trigger.
type KafkaTrigger struct {
	// URL of the Kafka broker.
	URL string `json:"url" protobuf:"bytes,1,opt,name=url"`
	// Name of the topic.
	// More info at https://kafka.apache.org/documentation/#intro_topics
	Topic string `json:"topic" protobuf:"bytes,2,opt,name=topic"`
	// Partition to write data to.
	Partition int32 `json:"partition" protobuf:"varint,3,opt,name=partition"`
	// Parameters is the list of parameters that is applied to resolved Kafka trigger object.

	Parameters []TriggerParameter `json:"parameters,omitempty" protobuf:"bytes,4,rep,name=parameters"`
	// RequiredAcks used in producer to tell the broker how many replica acknowledgements
	// Defaults to 1 (Only wait for the leader to ack).
	// +optional.
	RequiredAcks int32 `json:"requiredAcks,omitempty" protobuf:"varint,5,opt,name=requiredAcks"`
	// Compress determines whether to compress message or not.
	// Defaults to false.
	// If set to true, compresses message using snappy compression.
	// +optional
	Compress bool `json:"compress,omitempty" protobuf:"varint,6,opt,name=compress"`
	// FlushFrequency refers to the frequency in milliseconds to flush batches.
	// Defaults to 500 milliseconds.
	// +optional
	FlushFrequency int32 `json:"flushFrequency,omitempty" protobuf:"varint,7,opt,name=flushFrequency"`
	// TLS configuration for the Kafka producer.
	// +optional
	TLS *TLSConfig `json:"tls,omitempty" protobuf:"bytes,8,opt,name=tls"`
	// Payload is the list of key-value extracted from an event payload to construct the request payload.

	Payload []TriggerParameter `json:"payload" protobuf:"bytes,9,rep,name=payload"`
	// The partitioning key for the messages put on the Kafka topic.
	// Defaults to broker url.
	// +optional.
	PartitioningKey string `json:"partitioningKey,omitempty" protobuf:"bytes,10,opt,name=partitioningKey"`
}

// NATSTrigger refers to the specification of the NATS trigger.
type NATSTrigger struct {
	// URL of the NATS cluster.
	URL string `json:"url" protobuf:"bytes,1,opt,name=url"`
	// Name of the subject to put message on.
	Subject string `json:"subject" protobuf:"bytes,2,opt,name=subject"`
	// Payload is the list of key-value extracted from an event payload to construct the request payload.

	Payload []TriggerParameter `json:"payload" protobuf:"bytes,3,rep,name=payload"`
	// Parameters is the list of parameters that is applied to resolved NATS trigger object.

	Parameters []TriggerParameter `json:"parameters,omitempty" protobuf:"bytes,4,rep,name=parameters"`
	// TLS configuration for the NATS producer.
	// +optional
	TLS *TLSConfig `json:"tls,omitempty" protobuf:"bytes,5,opt,name=tls"`
}

// CustomTrigger refers to the specification of the custom trigger.
type CustomTrigger struct {
	// ServerURL is the url of the gRPC server that executes custom trigger
	ServerURL string `json:"serverURL" protobuf:"bytes,1,opt,name=serverURL"`
	// Secure refers to type of the connection between sensor to custom trigger gRPC
	Secure bool `json:"secure" protobuf:"varint,2,opt,name=secure"`
	// CertFilePath is path to the cert file within sensor for secure connection between sensor and custom trigger gRPC server.
	CertFilePath string `json:"certFilePath,omitempty" protobuf:"bytes,3,opt,name=certFilePath"`
	// ServerNameOverride for the secure connection between sensor and custom trigger gRPC server.
	ServerNameOverride string `json:"serverNameOverride,omitempty" protobuf:"bytes,4,opt,name=serverNameOverride"`
	// Spec is the custom trigger resource specification that custom trigger gRPC server knows how to interpret.
	Spec map[string]string `json:"spec" protobuf:"bytes,5,rep,name=spec"`
	// Parameters is the list of parameters that is applied to resolved custom trigger trigger object.

	Parameters []TriggerParameter `json:"parameters,omitempty" protobuf:"bytes,6,rep,name=parameters"`
	// Payload is the list of key-value extracted from an event payload to construct the request payload.

	Payload []TriggerParameter `json:"payload" protobuf:"bytes,7,rep,name=payload"`
}

// SlackTrigger refers to the specification of the slack notification trigger.
type SlackTrigger struct {
	// Parameters is the list of key-value extracted from event's payload that are applied to
	// the trigger resource.

	// +optional
	Parameters []TriggerParameter `json:"parameters,omitempty" protobuf:"bytes,1,rep,name=parameters"`
	// SlackToken refers to the Kubernetes secret that holds the slack token required to send messages.
	SlackToken *corev1.SecretKeySelector `json:"slackToken,omitempty" protobuf:"bytes,2,opt,name=slackToken"`
	// Namespace to read the password secret from.
	// This is required if the password secret selector is specified.
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
	// Channel refers to which Slack channel to send slack message.
	// +optional
	Channel string `json:"channel,omitempty" protobuf:"bytes,4,opt,name=channel"`
	// Message refers to the message to send to the Slack channel.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// OpenWhiskTrigger refers to the specification of the OpenWhisk trigger.
type OpenWhiskTrigger struct {
	// Host URL of the OpenWhisk.
	Host string `json:"host" protobuf:"bytes,1,opt,name=host"`
	// Version for the API.
	// Defaults to v1.
	// +optional
	Version string `json:"version,omitempty" protobuf:"bytes,2,opt,name=version"`
	// Namespace for the action.
	// Defaults to "_".
	// +optional.
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
	// AuthToken for authentication.
	// +optional
	AuthToken *corev1.SecretKeySelector `json:"authToken,omitempty" protobuf:"bytes,4,opt,name=authToken"`
	// Name of the action/function.
	ActionName string `json:"actionName" protobuf:"bytes,5,opt,name=actionName"`
	// Payload is the list of key-value extracted from an event payload to construct the request payload.

	Payload []TriggerParameter `json:"payload" protobuf:"bytes,6,rep,name=payload"`
	// Parameters is the list of key-value extracted from event's payload that are applied to
	// the trigger resource.

	// +optional
	Parameters []TriggerParameter `json:"parameters,omitempty" protobuf:"bytes,7,rep,name=parameters"`
}

// TriggerParameterOperation represents how to set a trigger destination
// resource key
type TriggerParameterOperation string

const (
	// TriggerParameterOpNone is the zero value of TriggerParameterOperation
	TriggerParameterOpNone TriggerParameterOperation = ""
	// TriggerParameterOpAppend means append the new value to the existing
	TriggerParameterOpAppend TriggerParameterOperation = "append"
	// TriggerParameterOpOverwrite means overwrite the existing value with the new
	TriggerParameterOpOverwrite TriggerParameterOperation = "overwrite"
	// TriggerParameterOpPrepend means prepend the new value to the existing
	TriggerParameterOpPrepend TriggerParameterOperation = "prepend"
)

// TriggerParameter indicates a passed parameter to a service template
type TriggerParameter struct {
	// Src contains a source reference to the value of the parameter from a dependency
	Src *TriggerParameterSource `json:"src,omitempty" protobuf:"bytes,1,opt,name=src"`
	// Dest is the JSONPath of a resource key.
	// A path is a series of keys separated by a dot. The colon character can be escaped with '.'
	// The -1 key can be used to append a value to an existing array.
	// See https://github.com/tidwall/sjson#path-syntax for more information about how this is used.
	Dest string `json:"dest" protobuf:"bytes,2,opt,name=dest"`
	// Operation is what to do with the existing value at Dest, whether to
	// 'prepend', 'overwrite', or 'append' it.
	Operation TriggerParameterOperation `json:"operation,omitempty" protobuf:"bytes,3,opt,name=operation,casttype=TriggerParameterOperation"`
}

// TriggerParameterSource defines the source for a parameter from a event event
type TriggerParameterSource struct {
	// DependencyName refers to the name of the dependency. The event which is stored for this dependency is used as payload
	// for the parameterization. Make sure to refer to one of the dependencies you have defined under Dependencies list.
	DependencyName string `json:"dependencyName" protobuf:"bytes,1,opt,name=dependencyName"`
	// ContextKey is the JSONPath of the event's (JSON decoded) context key
	// ContextKey is a series of keys separated by a dot. A key may contain wildcard characters '*' and '?'.
	// To access an array value use the index as the key. The dot and wildcard characters can be escaped with '\\'.
	// See https://github.com/tidwall/gjson#path-syntax for more information on how to use this.
	ContextKey string `json:"contextKey,omitempty" protobuf:"bytes,2,opt,name=contextKey"`
	// ContextTemplate is a go-template for extracting a string from the event's context.
	// If a ContextTemplate is provided with a ContextKey, the template will be evaluated first and fallback to the ContextKey.
	// The templating follows the standard go-template syntax as well as sprig's extra functions.
	// See https://pkg.go.dev/text/template and https://masterminds.github.io/sprig/
	ContextTemplate string `json:"contextTemplate,omitempty" protobuf:"bytes,3,opt,name=contextTemplate"`
	// DataKey is the JSONPath of the event's (JSON decoded) data key
	// DataKey is a series of keys separated by a dot. A key may contain wildcard characters '*' and '?'.
	// To access an array value use the index as the key. The dot and wildcard characters can be escaped with '\\'.
	// See https://github.com/tidwall/gjson#path-syntax for more information on how to use this.
	DataKey string `json:"dataKey,omitempty" protobuf:"bytes,4,opt,name=dataKey"`
	// DataTemplate is a go-template for extracting a string from the event's data.
	// If a DataTemplate is provided with a DataKey, the template will be evaluated first and fallback to the DataKey.
	// The templating follows the standard go-template syntax as well as sprig's extra functions.
	// See https://pkg.go.dev/text/template and https://masterminds.github.io/sprig/
	DataTemplate string `json:"dataTemplate,omitempty" protobuf:"bytes,5,opt,name=dataTemplate"`
	// Value is the default literal value to use for this parameter source
	// This is only used if the DataKey is invalid.
	// If the DataKey is invalid and this is not defined, this param source will produce an error.
	Value *string `json:"value,omitempty" protobuf:"bytes,6,opt,name=value"`
}

// TriggerPolicy dictates the policy for the trigger retries
type TriggerPolicy struct {
	// K8SResourcePolicy refers to the policy used to check the state of K8s based triggers using using labels
	K8s *K8SResourcePolicy `json:"k8s,omitempty" protobuf:"bytes,1,opt,name=k8s"`
	// Status refers to the policy used to check the state of the trigger using response status
	Status *StatusPolicy `json:"status,omitempty" protobuf:"bytes,2,opt,name=status"`
}

// K8SResourcePolicy refers to the policy used to check the state of K8s based triggers using using labels
type K8SResourcePolicy struct {
	// Labels required to identify whether a resource is in success state
	Labels map[string]string `json:"labels,omitempty" protobuf:"bytes,1,rep,name=labels"`
	// Backoff before checking resource state
	Backoff apicommon.Backoff `json:"backoff" protobuf:"bytes,2,opt,name=backoff"`
	// ErrorOnBackoffTimeout determines whether sensor should transition to error state if the trigger policy is unable to determine
	// the state of the resource
	ErrorOnBackoffTimeout bool `json:"errorOnBackoffTimeout" protobuf:"varint,3,opt,name=errorOnBackoffTimeout"`
}

// StatusPolicy refers to the policy used to check the state of the trigger using response status
type StatusPolicy struct {
	// Allow refers to the list of allowed response statuses. If the response status of the the trigger is within the list,
	// the trigger will marked as successful else it will result in trigger failure.

	Allow []int32 `json:"allow" protobuf:"varint,1,rep,name=allow"`
}

func (in *StatusPolicy) GetAllow() []int {
	statuses := make([]int, len(in.Allow))
	for i, s := range in.Allow {
		statuses[i] = int(s)
	}
	return statuses
}

// SensorResources holds the metadata of the resources created for the sensor
type SensorResources struct {
	// Deployment holds the metadata of the deployment for the sensor
	Deployment *metav1.ObjectMeta `json:"deployment,omitempty" protobuf:"bytes,1,opt,name=deployment"`
	// Service holds the metadata of the service for the sensor
	// +optional
	Service *metav1.ObjectMeta `json:"service,omitempty" protobuf:"bytes,2,opt,name=service"`
}

// SensorStatus contains information about the status of a sensor.
type SensorStatus struct {
	// Phase is the high-level summary of the sensor.
	Phase NodePhase `json:"phase" protobuf:"bytes,1,opt,name=phase,casttype=NodePhase"`
	// StartedAt is the time at which this sensor was initiated
	StartedAt metav1.Time `json:"startedAt,omitempty" protobuf:"bytes,2,opt,name=startedAt"`
	// CompletedAt is the time at which this sensor was completed
	CompletedAt metav1.Time `json:"completedAt,omitempty" protobuf:"bytes,3,opt,name=completedAt"`
	// Message is a human readable string indicating details about a sensor in its phase
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
	// Nodes is a mapping between a node ID and the node's status
	// it records the states for the FSM of this sensor.
	Nodes map[string]NodeStatus `json:"nodes,omitempty" protobuf:"bytes,5,rep,name=nodes"`
	// TriggerCycleCount is the count of sensor's trigger cycle runs.
	TriggerCycleCount int32 `json:"triggerCycleCount,omitempty" protobuf:"varint,6,opt,name=triggerCycleCount"`
	// TriggerCycleState is the status from last cycle of triggers execution.
	TriggerCycleStatus TriggerCycleState `json:"triggerCycleStatus" protobuf:"bytes,7,opt,name=triggerCycleStatus,casttype=TriggerCycleState"`
	// LastCycleTime is the time when last trigger cycle completed
	LastCycleTime metav1.Time `json:"lastCycleTime" protobuf:"bytes,8,opt,name=lastCycleTime"`
	// Resources refers to metadata of the resources created for the sensor
	Resources *SensorResources `json:"resources,omitempty" protobuf:"bytes,9,opt,name=resources"`
}

// NodeStatus describes the status for an individual node in the sensor's FSM.
// A single node can represent the status for event or a trigger.
type NodeStatus struct {
	// ID is a unique identifier of a node within a sensor
	// It is a hash of the node name
	ID string `json:"id" protobuf:"bytes,1,opt,name=id"`
	// Name is a unique name in the node tree used to generate the node ID
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`
	// DisplayName is the human readable representation of the node
	DisplayName string `json:"displayName" protobuf:"bytes,3,opt,name=displayName"`
	// Type is the type of the node
	Type NodeType `json:"type" protobuf:"bytes,4,opt,name=type,casttype=NodeType"`
	// Phase of the node
	Phase NodePhase `json:"phase" protobuf:"bytes,5,opt,name=phase,casttype=NodePhase"`
	// StartedAt is the time at which this node started
	StartedAt metav1.MicroTime `json:"startedAt,omitempty" protobuf:"bytes,6,opt,name=startedAt"`
	// CompletedAt is the time at which this node completed
	CompletedAt metav1.MicroTime `json:"completedAt,omitempty" protobuf:"bytes,7,opt,name=completedAt"`
	// store data or something to save for event notifications or trigger events
	Message string `json:"message,omitempty" protobuf:"bytes,8,opt,name=message"`
	// Event stores the last seen event for this node
	Event *Event `json:"event,omitempty" protobuf:"bytes,9,opt,name=event"`
	// UpdatedAt refers to the time at which the node was updated.
	UpdatedAt metav1.MicroTime `json:"updatedAt,omitempty" protobuf:"bytes,10,opt,name=updatedAt"`
	// ResolvedAt refers to the time at which the node was resolved.
	ResolvedAt metav1.MicroTime `json:"resolvedAt,omitempty" protobuf:"bytes,11,opt,name=resolvedAt"`
}

// ArtifactLocation describes the source location for an external artifact
type ArtifactLocation struct {
	// S3 compliant artifact
	S3 *apicommon.S3Artifact `json:"s3,omitempty" protobuf:"bytes,1,opt,name=s3"`
	// Inline artifact is embedded in sensor spec as a string
	Inline *string `json:"inline,omitempty" protobuf:"bytes,2,opt,name=inline"`
	// File artifact is artifact stored in a file
	File *FileArtifact `json:"file,omitempty" protobuf:"bytes,3,opt,name=file"`
	// URL to fetch the artifact from
	URL *URLArtifact `json:"url,omitempty" protobuf:"bytes,4,opt,name=url"`
	// Configmap that stores the artifact
	Configmap *ConfigmapArtifact `json:"configmap,omitempty" protobuf:"bytes,5,opt,name=configmap"`
	// Git repository hosting the artifact
	Git *GitArtifact `json:"git,omitempty" protobuf:"bytes,6,opt,name=git"`
	// Resource is generic template for K8s resource
	Resource *apicommon.Resource `json:"resource,omitempty" protobuf:"bytes,7,opt,name=resource"`
}

// ConfigmapArtifact contains information about artifact in k8 configmap
type ConfigmapArtifact struct {
	// Name of the configmap
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Namespace where configmap is deployed
	Namespace string `json:"namespace" protobuf:"bytes,2,opt,name=namespace"`
	// Key within configmap data which contains trigger resource definition
	Key string `json:"key" protobuf:"bytes,3,opt,name=key"`
}

// FileArtifact contains information about an artifact in a filesystem
type FileArtifact struct {
	Path string `json:"path,omitempty" protobuf:"bytes,1,opt,name=path"`
}

// URLArtifact contains information about an artifact at an http endpoint.
type URLArtifact struct {
	// Path is the complete URL
	Path string `json:"path" protobuf:"bytes,1,opt,name=path"`
	// VerifyCert decides whether the connection is secure or not
	VerifyCert bool `json:"verifyCert,omitempty" protobuf:"varint,2,opt,name=verifyCert"`
}

// GitArtifact contains information about an artifact stored in git
type GitArtifact struct {
	// Git URL
	URL string `json:"url" protobuf:"bytes,1,opt,name=url"`
	// Directory to clone the repository. We clone complete directory because GitArtifact is not limited to any specific Git service providers.
	// Hence we don't use any specific git provider client.
	CloneDirectory string `json:"cloneDirectory" protobuf:"bytes,2,opt,name=cloneDirectory"`
	// Creds contain reference to git username and password
	// +optional
	Creds *GitCreds `json:"creds,omitempty" protobuf:"bytes,3,opt,name=creds"`
	// Namespace where creds are stored.
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,4,opt,name=namespace"`
	// SSHKeyPath is path to your ssh key path. Use this if you don't want to provide username and password.
	// ssh key path must be mounted in sensor pod.
	// +optional
	SSHKeyPath string `json:"sshKeyPath,omitempty" protobuf:"bytes,5,opt,name=sshKeyPath"`
	// Path to file that contains trigger resource definition
	FilePath string `json:"filePath" protobuf:"bytes,6,opt,name=filePath"`
	// Branch to use to pull trigger resource
	// +optional
	Branch string `json:"branch,omitempty" protobuf:"bytes,7,opt,name=branch"`
	// Tag to use to pull trigger resource
	// +optional
	Tag string `json:"tag,omitempty" protobuf:"bytes,8,opt,name=tag"`
	// Ref to use to pull trigger resource. Will result in a shallow clone and
	// fetch.
	// +optional
	Ref string `json:"ref,omitempty" protobuf:"bytes,9,opt,name=ref"`
	// Remote to manage set of tracked repositories. Defaults to "origin".
	// Refer https://git-scm.com/docs/git-remote
	// +optional
	Remote *GitRemoteConfig `json:"remote,omitempty" protobuf:"bytes,10,opt,name=remote"`
}

// GitRemoteConfig contains the configuration of a Git remote
type GitRemoteConfig struct {
	// Name of the remote to fetch from.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// URLs the URLs of a remote repository. It must be non-empty. Fetch will
	// always use the first URL, while push will use all of them.
	URLS []string `json:"urls" protobuf:"bytes,2,rep,name=urls"`
}

// GitCreds contain reference to git username and password
type GitCreds struct {
	Username *corev1.SecretKeySelector `json:"username,omitempty" protobuf:"bytes,1,opt,name=username"`
	Password *corev1.SecretKeySelector `json:"password,omitempty" protobuf:"bytes,2,opt,name=password"`
}

// Event represents the cloudevent received from a gateway.
type Event struct {
	Context *EventContext `json:"context,omitempty" protobuf:"bytes,1,opt,name=context"`
	Data    []byte        `json:"data" protobuf:"bytes,2,opt,name=data"`
}

// EventContext holds the context of the cloudevent received from a gateway.
type EventContext struct {
	// ID of the event; must be non-empty and unique within the scope of the producer.
	ID string `json:"id" protobuf:"bytes,1,opt,name=id"`
	// Source - A URI describing the event producer.
	Source string `json:"source" protobuf:"bytes,2,opt,name=source"`
	// SpecVersion - The version of the CloudEvents specification used by the event.
	SpecVersion string `json:"specversion" protobuf:"bytes,3,opt,name=specversion"`
	// Type - The type of the occurrence which has happened.
	Type string `json:"type" protobuf:"bytes,4,opt,name=type"`
	// DataContentType - A MIME (RFC2046) string describing the media type of `data`.
	DataContentType string `json:"dataContentType" protobuf:"bytes,5,opt,name=dataContentType"`
	// Subject - The subject of the event in the context of the event producer
	Subject string `json:"subject" protobuf:"bytes,6,opt,name=subject"`
	// Time - A Timestamp when the event happened.
	Time metav1.Time `json:"time" protobuf:"bytes,7,opt,name=time"`
}

// HasLocation whether or not an artifact has a location defined
func (a *ArtifactLocation) HasLocation() bool {
	return a.S3 != nil || a.Inline != nil || a.File != nil || a.URL != nil
}

// IsComplete determines if the node has reached an end state
func (node NodeStatus) IsComplete() bool {
	return node.Phase == NodePhaseComplete ||
		node.Phase == NodePhaseError
}

// IsComplete determines if the sensor has reached an end state
func (s *Sensor) IsComplete() bool {
	if !(s.Status.Phase == NodePhaseComplete || s.Status.Phase == NodePhaseError) {
		return false
	}
	for _, node := range s.Status.Nodes {
		if !node.IsComplete() {
			return false
		}
	}
	return true
}

// AreAllNodesSuccess determines if all nodes of the given type have completed successfully
func (s *Sensor) AreAllNodesSuccess(nodeType NodeType) bool {
	for _, node := range s.Status.Nodes {
		if node.Type == nodeType && node.Phase != NodePhaseComplete {
			return false
		}
	}
	return true
}

// NodeID creates a deterministic node ID based on a node name
// we support 3 kinds of "nodes" - sensors, events, triggers
// each should pass it's name field
func (s *Sensor) NodeID(name string) string {
	if name == s.ObjectMeta.Name {
		return s.ObjectMeta.Name
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(name))
	return fmt.Sprintf("%s-%v", s.ObjectMeta.Name, h.Sum32())
}
