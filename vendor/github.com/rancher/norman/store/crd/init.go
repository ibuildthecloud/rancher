package crd

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"golang.org/x/sync/errgroup"

	"github.com/rancher/norman/store/proxy"
	"github.com/rancher/norman/types"
	"github.com/rancher/wrangler/pkg/crd"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

type Factory struct {
	clientGetter proxy.ClientGetter
	factory      *crd.Factory
	eg           errgroup.Group
}

func NewFactoryFromClientGetter(clientGetter proxy.ClientGetter) (*Factory, error) {
	return &Factory{
		factory: &crd.Factory{
			CRDClient: clientGetter.APIExtClient(),
		},
		clientGetter: clientGetter,
	}, nil
}

func (f *Factory) BatchWait() error {
	return f.eg.Wait()
}

func (f *Factory) BatchCreateCRDs(ctx context.Context, storageContext types.StorageContext, schemas *types.Schemas, version *types.APIVersion, schemaIDs ...string) {
	f.eg.Go(func() error {
		var schemasToCreate []*types.Schema

		for _, schemaID := range schemaIDs {
			s := schemas.Schema(version, schemaID)
			if s == nil {
				return fmt.Errorf("can not find schema %s", schemaID)
			}
			schemasToCreate = append(schemasToCreate, s)
		}

		err := f.assignStores(ctx, storageContext, schemasToCreate...)
		if err != nil {
			return fmt.Errorf("creating CRD store %s", err.Error())
		}

		return nil
	})
}

func (f *Factory) assignStores(ctx context.Context, storageContext types.StorageContext, schemas ...*types.Schema) error {
	schemaStatus, err := f.createCRDs(ctx, schemas...)
	if err != nil {
		return err
	}

	for _, schema := range schemas {
		crd, ok := schemaStatus[schema]
		if !ok {
			return fmt.Errorf("failed to create create/find CRD for %s", schema.ID)
		}

		schema.Store = proxy.NewProxyStore(ctx, f.clientGetter,
			storageContext,
			[]string{"apis"},
			crd.Spec.Group,
			crd.Spec.Version,
			crd.Status.AcceptedNames.Kind,
			crd.Status.AcceptedNames.Plural)
	}

	return nil
}

func (f *Factory) createCRDs(ctx context.Context, schemas ...*types.Schema) (map[*types.Schema]*apiext.CustomResourceDefinition, error) {
	gvkToSchema := map[schema.GroupVersionKind]*types.Schema{}
	var crds []crd.CRD

	for _, s := range schemas {
		if s.ReflectType == nil {
			continue
		}
		gvk := schema.GroupVersionKind{
			Group:   s.Version.Group,
			Version: s.Version.Version,
			Kind:    s.CodeName,
		}
		crd := crd.CRD{
			GVK:          gvk,
			PluralName:   strings.ToLower(s.PluralName),
			NonNamespace: s.Scope != types.NamespaceScope,
			Categories:   []string{"rancher"},
		}.WithColumnsFromStruct(s.ReflectType).
			WithSchemaFromStruct(s.ReflectType)

		crds = append(crds, crd)
		gvkToSchema[gvk] = s
	}

	created, err := f.factory.CreateCRDs(ctx, crds...)
	if err != nil {
		return nil, err
	}

	result := map[*types.Schema]*apiext.CustomResourceDefinition{}
	for gvk, crd := range created {
		result[gvkToSchema[gvk]] = crd
	}

	return result, nil
}
