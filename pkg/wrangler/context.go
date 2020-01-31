package wrangler

import (
	"context"

	"github.com/rancher/rancher/pkg/wrangler/generated/controllers/management.cattle.io"
	managementv3 "github.com/rancher/rancher/pkg/wrangler/generated/controllers/management.cattle.io/v3"
	"github.com/rancher/wrangler-api/pkg/generated/controllers/rbac"
	rbacv1 "github.com/rancher/wrangler-api/pkg/generated/controllers/rbac/v1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/start"
	"k8s.io/client-go/rest"
)

type WranglerContext struct {
	Apply    apply.Apply
	RBAC     rbacv1.Interface
	Mgmt     managementv3.Interface
	starters []start.Starter
}

func (w *WranglerContext) Start(ctx context.Context) error {
	return start.All(ctx, 5, w.starters...)
}

func NewContext(restConfig *rest.Config) (*WranglerContext, error) {
	rbac, err := rbac.NewFactoryFromConfig(restConfig)
	if err != nil {
		return nil, err
	}

	apply, err := apply.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	mgmt, err := management.NewFactoryFromConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return &WranglerContext{
		Apply: apply,
		RBAC:  rbac.Rbac().V1(),
		Mgmt:  mgmt.Management().V3(),
		starters: []start.Starter{
			rbac,
			mgmt,
		},
	}, nil
}
