/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package v1

import (
	"reflect"

	"k8s.io/kubernetes/pkg/api"
	v1 "k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/experimental"
	"k8s.io/kubernetes/pkg/conversion"
	"k8s.io/kubernetes/pkg/util"
)

func addConversionFuncs() {
	// Add non-generated conversion functions
	err := api.Scheme.AddConversionFuncs(
		convert_api_PodSpec_To_v1_PodSpec,
		convert_v1_PodSpec_To_api_PodSpec,
		convert_api_VolumeSource_To_v1_VolumeSource,
		convert_experimental_DeploymentSpec_To_v1_DeploymentSpec,
		convert_v1_DeploymentSpec_To_experimental_DeploymentSpec,
		convert_experimental_DeploymentStrategy_To_v1_DeploymentStrategy,
		convert_v1_DeploymentStrategy_To_experimental_DeploymentStrategy,
		convert_experimental_RollingUpdateDeployment_To_v1_RollingUpdateDeployment,
		convert_v1_RollingUpdateDeployment_To_experimental_RollingUpdateDeployment,
		convert_v1_VolumeSource_To_api_VolumeSource,
	)
	if err != nil {
		// If one of the conversion functions is malformed, detect it immediately.
		panic(err)
	}
}

// The following two PodSpec conversions functions where copied from pkg/api/conversion.go
// for the generated functions to work properly.
// This should be fixed: https://github.com/kubernetes/kubernetes/issues/12977
func convert_api_PodSpec_To_v1_PodSpec(in *api.PodSpec, out *v1.PodSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*api.PodSpec))(in)
	}
	if in.Volumes != nil {
		out.Volumes = make([]v1.Volume, len(in.Volumes))
		for i := range in.Volumes {
			if err := convert_api_Volume_To_v1_Volume(&in.Volumes[i], &out.Volumes[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Volumes = nil
	}
	if in.Containers != nil {
		out.Containers = make([]v1.Container, len(in.Containers))
		for i := range in.Containers {
			if err := convert_api_Container_To_v1_Container(&in.Containers[i], &out.Containers[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Containers = nil
	}
	out.RestartPolicy = v1.RestartPolicy(in.RestartPolicy)
	if in.TerminationGracePeriodSeconds != nil {
		out.TerminationGracePeriodSeconds = new(int64)
		*out.TerminationGracePeriodSeconds = *in.TerminationGracePeriodSeconds
	} else {
		out.TerminationGracePeriodSeconds = nil
	}
	if in.ActiveDeadlineSeconds != nil {
		out.ActiveDeadlineSeconds = new(int64)
		*out.ActiveDeadlineSeconds = *in.ActiveDeadlineSeconds
	} else {
		out.ActiveDeadlineSeconds = nil
	}
	out.DNSPolicy = v1.DNSPolicy(in.DNSPolicy)
	if in.NodeSelector != nil {
		out.NodeSelector = make(map[string]string)
		for key, val := range in.NodeSelector {
			out.NodeSelector[key] = val
		}
	} else {
		out.NodeSelector = nil
	}
	out.ServiceAccountName = in.ServiceAccountName
	// DeprecatedServiceAccount is an alias for ServiceAccountName.
	out.DeprecatedServiceAccount = in.ServiceAccountName
	out.NodeName = in.NodeName
	out.HostNetwork = in.HostNetwork
	if in.ImagePullSecrets != nil {
		out.ImagePullSecrets = make([]v1.LocalObjectReference, len(in.ImagePullSecrets))
		for i := range in.ImagePullSecrets {
			if err := convert_api_LocalObjectReference_To_v1_LocalObjectReference(&in.ImagePullSecrets[i], &out.ImagePullSecrets[i], s); err != nil {
				return err
			}
		}
	} else {
		out.ImagePullSecrets = nil
	}
	return nil
}

func convert_v1_PodSpec_To_api_PodSpec(in *v1.PodSpec, out *api.PodSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*v1.PodSpec))(in)
	}
	if in.Volumes != nil {
		out.Volumes = make([]api.Volume, len(in.Volumes))
		for i := range in.Volumes {
			if err := convert_v1_Volume_To_api_Volume(&in.Volumes[i], &out.Volumes[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Volumes = nil
	}
	if in.Containers != nil {
		out.Containers = make([]api.Container, len(in.Containers))
		for i := range in.Containers {
			if err := convert_v1_Container_To_api_Container(&in.Containers[i], &out.Containers[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Containers = nil
	}
	out.RestartPolicy = api.RestartPolicy(in.RestartPolicy)
	if in.TerminationGracePeriodSeconds != nil {
		out.TerminationGracePeriodSeconds = new(int64)
		*out.TerminationGracePeriodSeconds = *in.TerminationGracePeriodSeconds
	} else {
		out.TerminationGracePeriodSeconds = nil
	}
	if in.ActiveDeadlineSeconds != nil {
		out.ActiveDeadlineSeconds = new(int64)
		*out.ActiveDeadlineSeconds = *in.ActiveDeadlineSeconds
	} else {
		out.ActiveDeadlineSeconds = nil
	}
	out.DNSPolicy = api.DNSPolicy(in.DNSPolicy)
	if in.NodeSelector != nil {
		out.NodeSelector = make(map[string]string)
		for key, val := range in.NodeSelector {
			out.NodeSelector[key] = val
		}
	} else {
		out.NodeSelector = nil
	}
	// We support DeprecatedServiceAccount as an alias for ServiceAccountName.
	// If both are specified, ServiceAccountName (the new field) wins.
	out.ServiceAccountName = in.ServiceAccountName
	if in.ServiceAccountName == "" {
		out.ServiceAccountName = in.DeprecatedServiceAccount
	}
	out.NodeName = in.NodeName
	out.HostNetwork = in.HostNetwork
	if in.ImagePullSecrets != nil {
		out.ImagePullSecrets = make([]api.LocalObjectReference, len(in.ImagePullSecrets))
		for i := range in.ImagePullSecrets {
			if err := convert_v1_LocalObjectReference_To_api_LocalObjectReference(&in.ImagePullSecrets[i], &out.ImagePullSecrets[i], s); err != nil {
				return err
			}
		}
	} else {
		out.ImagePullSecrets = nil
	}
	return nil
}

func convert_experimental_DeploymentSpec_To_v1_DeploymentSpec(in *experimental.DeploymentSpec, out *DeploymentSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*experimental.DeploymentSpec))(in)
	}
	out.Replicas = new(int)
	*out.Replicas = in.Replicas
	if in.Selector != nil {
		out.Selector = make(map[string]string)
		for key, val := range in.Selector {
			out.Selector[key] = val
		}
	} else {
		out.Selector = nil
	}
	if in.Template != nil {
		out.Template = new(v1.PodTemplateSpec)
		if err := convert_api_PodTemplateSpec_To_v1_PodTemplateSpec(in.Template, out.Template, s); err != nil {
			return err
		}
	} else {
		out.Template = nil
	}
	if err := convert_experimental_DeploymentStrategy_To_v1_DeploymentStrategy(&in.Strategy, &out.Strategy, s); err != nil {
		return err
	}
	out.UniqueLabelKey = new(string)
	*out.UniqueLabelKey = in.UniqueLabelKey
	return nil
}

func convert_v1_DeploymentSpec_To_experimental_DeploymentSpec(in *DeploymentSpec, out *experimental.DeploymentSpec, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*DeploymentSpec))(in)
	}
	if in.Replicas != nil {
		out.Replicas = *in.Replicas
	}
	if in.Selector != nil {
		out.Selector = make(map[string]string)
		for key, val := range in.Selector {
			out.Selector[key] = val
		}
	} else {
		out.Selector = nil
	}
	if in.Template != nil {
		out.Template = new(api.PodTemplateSpec)
		if err := convert_v1_PodTemplateSpec_To_api_PodTemplateSpec(in.Template, out.Template, s); err != nil {
			return err
		}
	} else {
		out.Template = nil
	}
	if err := convert_v1_DeploymentStrategy_To_experimental_DeploymentStrategy(&in.Strategy, &out.Strategy, s); err != nil {
		return err
	}
	if in.UniqueLabelKey != nil {
		out.UniqueLabelKey = *in.UniqueLabelKey
	}
	return nil
}

func convert_experimental_DeploymentStrategy_To_v1_DeploymentStrategy(in *experimental.DeploymentStrategy, out *DeploymentStrategy, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*experimental.DeploymentStrategy))(in)
	}
	out.Type = DeploymentType(in.Type)
	if in.RollingUpdate != nil {
		out.RollingUpdate = new(RollingUpdateDeployment)
		if err := convert_experimental_RollingUpdateDeployment_To_v1_RollingUpdateDeployment(in.RollingUpdate, out.RollingUpdate, s); err != nil {
			return err
		}
	} else {
		out.RollingUpdate = nil
	}
	return nil
}

func convert_v1_DeploymentStrategy_To_experimental_DeploymentStrategy(in *DeploymentStrategy, out *experimental.DeploymentStrategy, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*DeploymentStrategy))(in)
	}
	out.Type = experimental.DeploymentType(in.Type)
	if in.RollingUpdate != nil {
		out.RollingUpdate = new(experimental.RollingUpdateDeployment)
		if err := convert_v1_RollingUpdateDeployment_To_experimental_RollingUpdateDeployment(in.RollingUpdate, out.RollingUpdate, s); err != nil {
			return err
		}
	} else {
		out.RollingUpdate = nil
	}
	return nil
}

func convert_experimental_RollingUpdateDeployment_To_v1_RollingUpdateDeployment(in *experimental.RollingUpdateDeployment, out *RollingUpdateDeployment, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*experimental.RollingUpdateDeployment))(in)
	}
	if out.MaxUnavailable == nil {
		out.MaxUnavailable = &util.IntOrString{}
	}
	if err := s.Convert(&in.MaxUnavailable, out.MaxUnavailable, 0); err != nil {
		return err
	}
	if out.MaxSurge == nil {
		out.MaxSurge = &util.IntOrString{}
	}
	if err := s.Convert(&in.MaxSurge, out.MaxSurge, 0); err != nil {
		return err
	}
	out.MinReadySeconds = in.MinReadySeconds
	return nil
}

func convert_v1_RollingUpdateDeployment_To_experimental_RollingUpdateDeployment(in *RollingUpdateDeployment, out *experimental.RollingUpdateDeployment, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*RollingUpdateDeployment))(in)
	}
	if err := s.Convert(in.MaxUnavailable, &out.MaxUnavailable, 0); err != nil {
		return err
	}
	if err := s.Convert(in.MaxSurge, &out.MaxSurge, 0); err != nil {
		return err
	}
	out.MinReadySeconds = in.MinReadySeconds
	return nil
}

// This will convert our internal represantation of VolumeSource to its v1 representation
// Used for keeping backwards compatibility for the Metadata field
func convert_api_VolumeSource_To_v1_VolumeSource(in *api.VolumeSource, out *v1.VolumeSource, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*api.VolumeSource))(in)
	}

	if err := s.DefaultConvert(in, out, conversion.IgnoreMissingFields); err != nil {
		return err
	}

	if in.DownwardAPI != nil {
		out.DownwardAPI = new(v1.DownwardAPIVolumeSource)
		if err := convert_api_DownwardAPIVolumeSource_To_v1_DownwardAPIVolumeSource(in.DownwardAPI, out.DownwardAPI, s); err != nil {
			return err
		}

		// also copy to Metadata
		out.Metadata = new(v1.MetadataVolumeSource)
		if err := convert_api_DownwardAPIVolumeSource_To_v1_MetadataVolumeSource(in.DownwardAPI, out.Metadata, s); err != nil {
			return err
		}
	} else {
		out.DownwardAPI = nil
	}
	return nil
}

// downward -> metadata (api -> v1)
func convert_api_DownwardAPIVolumeSource_To_v1_MetadataVolumeSource(in *api.DownwardAPIVolumeSource, out *v1.MetadataVolumeSource, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*api.DownwardAPIVolumeSource))(in)
	}
	if in.Items != nil {
		out.Items = make([]v1.MetadataFile, len(in.Items))
		for i := range in.Items {
			if err := convert_api_DownwardAPIVolumeFile_To_v1_MetadataFile(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	}
	return nil
}

func convert_api_DownwardAPIVolumeFile_To_v1_MetadataFile(in *api.DownwardAPIVolumeFile, out *v1.MetadataFile, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*api.DownwardAPIVolumeFile))(in)
	}
	out.Name = in.Path
	if err := convert_api_ObjectFieldSelector_To_v1_ObjectFieldSelector(&in.FieldRef, &out.FieldRef, s); err != nil {
		return err
	}
	return nil
}

// This will convert the v1 representation of VolumeSource to our internal representation
// Used for keeping backwards compatibility for the Metadata field
func convert_v1_VolumeSource_To_api_VolumeSource(in *v1.VolumeSource, out *api.VolumeSource, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*v1.VolumeSource))(in)
	}

	if err := s.DefaultConvert(in, out, conversion.IgnoreMissingFields); err != nil {
		return err
	}

	// if specified Metadata will stomp DownwardAPI
	if in.Metadata != nil {
		out.DownwardAPI = new(api.DownwardAPIVolumeSource)
		if err := convert_v1_MetadataVolumeSource_To_api_DownwardAPIVolumeSource(in.Metadata, out.DownwardAPI, s); err != nil {
			return err
		}
	} else {
		if in.DownwardAPI != nil {
			out.DownwardAPI = new(api.DownwardAPIVolumeSource)
			if err := convert_v1_DownwardAPIVolumeSource_To_api_DownwardAPIVolumeSource(in.DownwardAPI, out.DownwardAPI, s); err != nil {
				return err
			}
		} else {
			out.DownwardAPI = nil
		}
	}
	return nil
}

// metadata -> downward (v1 -> api)
func convert_v1_MetadataVolumeSource_To_api_DownwardAPIVolumeSource(in *v1.MetadataVolumeSource, out *api.DownwardAPIVolumeSource, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*v1.MetadataVolumeSource))(in)
	}
	if in.Items != nil {
		out.Items = make([]api.DownwardAPIVolumeFile, len(in.Items))
		for i := range in.Items {
			if err := convert_v1_MetadataFile_To_api_DownwardAPIVolumeFile(&in.Items[i], &out.Items[i], s); err != nil {
				return err
			}
		}
	}
	return nil
}

func convert_v1_MetadataFile_To_api_DownwardAPIVolumeFile(in *v1.MetadataFile, out *api.DownwardAPIVolumeFile, s conversion.Scope) error {
	if defaulting, found := s.DefaultingInterface(reflect.TypeOf(*in)); found {
		defaulting.(func(*v1.MetadataFile))(in)
	}
	out.Path = in.Name
	if err := convert_v1_ObjectFieldSelector_To_api_ObjectFieldSelector(&in.FieldRef, &out.FieldRef, s); err != nil {
		return err
	}
	return nil
}
