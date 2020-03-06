package clusterauthtoken

import (
	"context"

	v3 "github.com/rancher/types/apis/cluster.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/rancher/wrangler/pkg/crd"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func CRDSetup(ctx context.Context, apiContext *config.UserOnlyContext) error {
	crds := []crd.CRD{
		crd.CRD{
			GVK: v3.ClusterAuthTokenGroupVersionKind,
		}.WithSchemaFromStruct(v3.ClusterAuthToken{}).
			WithColumn("Username", ".userName").
			WithCustomColumn(v1beta1.CustomResourceColumnDefinition{
				Name:     "Expires",
				Type:     "date",
				JSONPath: ".expiresAt",
			}).
			WithColumn("LastRefresh", ".lastRefresh").
			WithColumn("Enabled", ".enabled"),
		crd.CRD{
			GVK: v3.ClusterUserAttributeGroupVersionKind,
		}.WithSchemaFromStruct(v3.ClusterUserAttribute{}).
			WithColumn("Groups", ".groups").
			WithColumn("LastRefresh", ".lastRefresh").
			WithColumn("Enabled", ".enabled"),
	}

	f, err := crd.NewFactoryFromClient(&apiContext.RESTConfig)
	if err != nil {
		return err
	}
	_, err = f.CreateCRDs(ctx, crds...)
	return err
}
