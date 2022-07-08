package v1

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hashicorp/terraform-provider-kubernetes/kubernetes/structures"
)

func DataSourceKubernetesConfigMap() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesConfigMapRead,

		Schema: map[string]*schema.Schema{
			"metadata": NamespacedMetadataSchema("config_map", false),
			"data": {
				Type:        schema.TypeMap,
				Description: "A map of the config map data.",
				Computed:    true,
			},
			"binary_data": {
				Type:        schema.TypeMap,
				Description: "A map of the config map binary data.",
				Computed:    true,
			},
		},
	}
}

func dataSourceKubernetesConfigMapRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	om := metav1.ObjectMeta{
		Namespace: d.Get("metadata.0.namespace").(string),
		Name:      d.Get("metadata.0.name").(string),
	}
	d.SetId(structures.BuildId(om))

	return resourceKubernetesConfigMapRead(ctx, d, meta)
}
