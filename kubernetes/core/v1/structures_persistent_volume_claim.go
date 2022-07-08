package v1

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"

	"github.com/hashicorp/terraform-provider-kubernetes/kubernetes/structures"
)

// Flatteners

func FlattenPersistentVolumeClaimSpec(in v1.PersistentVolumeClaimSpec) []interface{} {
	att := make(map[string]interface{})
	att["access_modes"] = structures.FlattenPersistentVolumeAccessModes(in.AccessModes)
	att["resources"] = flattenResourceRequirements(in.Resources)
	if in.Selector != nil {
		att["selector"] = structures.FlattenLabelSelector(in.Selector)
	}
	if in.VolumeName != "" {
		att["volume_name"] = in.VolumeName
	}
	if in.StorageClassName != nil {
		att["storage_class_name"] = *in.StorageClassName
	}
	return []interface{}{att}
}

func flattenResourceRequirements(in v1.ResourceRequirements) []interface{} {
	att := make(map[string]interface{})
	if len(in.Limits) > 0 {
		att["limits"] = structures.FlattenResourceList(in.Limits)
	}
	if len(in.Requests) > 0 {
		att["requests"] = structures.FlattenResourceList(in.Requests)
	}
	return []interface{}{att}
}

// Expanders

func ExpandPersistentVolumeClaim(p map[string]interface{}) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	if len(p) == 0 {
		return pvc, nil
	}
	m, ok := p["metadata"].([]interface{})
	if !ok {
		return pvc, errors.New("persistent_volume_claim: failed to expand 'metadata'")
	}
	pvc.ObjectMeta = structures.ExpandMetadata(m)
	s, ok := p["spec"].([]interface{})
	if !ok {
		return pvc, errors.New("persistent_volume_claim: failed to expand 'spec'")
	}
	spec, err := expandPersistentVolumeClaimSpec(s)
	if err != nil {
		return pvc, err
	}
	pvc.Spec = *spec
	return pvc, nil
}

func expandPersistentVolumeClaimSpec(l []interface{}) (*corev1.PersistentVolumeClaimSpec, error) {
	obj := &corev1.PersistentVolumeClaimSpec{}
	if len(l) == 0 || l[0] == nil {
		return obj, nil
	}
	in := l[0].(map[string]interface{})
	resourceRequirements, err := expandResourceRequirements(in["resources"].([]interface{}))
	if err != nil {
		return nil, err
	}
	obj.AccessModes = structures.ExpandPersistentVolumeAccessModes(in["access_modes"].(*schema.Set).List())
	obj.Resources = *resourceRequirements
	if v, ok := in["selector"].([]interface{}); ok && len(v) > 0 {
		obj.Selector = structures.ExpandLabelSelector(v)
	}
	if v, ok := in["volume_name"].(string); ok {
		obj.VolumeName = v
	}
	if v, ok := in["storage_class_name"].(string); ok && v != "" {
		obj.StorageClassName = structures.PtrToString(v)
	}
	return obj, nil
}

func expandResourceRequirements(l []interface{}) (*corev1.ResourceRequirements, error) {
	obj := &corev1.ResourceRequirements{}
	if len(l) == 0 || l[0] == nil {
		return obj, nil
	}
	in := l[0].(map[string]interface{})
	if v, ok := in["limits"].(map[string]interface{}); ok && len(v) > 0 {
		rl, err := structures.ExpandMapToResourceList(v)
		if err != nil {
			return obj, err
		}
		obj.Limits = *rl
	}
	if v, ok := in["requests"].(map[string]interface{}); ok && len(v) > 0 {
		rq, err := structures.ExpandMapToResourceList(v)
		if err != nil {
			return obj, err
		}
		obj.Requests = *rq
	}
	return obj, nil
}
