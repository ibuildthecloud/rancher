package controllers

import (
	"context"
	"net/url"

	"github.com/rancher/rancher/pkg/agent/controllers/deploy"
	"github.com/rancher/rancher/pkg/wrangler/generated/controllers/management.cattle.io"
	apply2 "github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/start"
	"k8s.io/client-go/rest"
)

func StartControllers(ctx context.Context, token, url, namespace string) error {
	cloudCfg, err := getCloudCfg(token, url)
	if err != nil {
		return err
	}

	clusterCfg, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	apply, err := apply2.NewForConfig(clusterCfg)
	if err != nil {
		return err
	}

	mgmt, err := management.NewFactoryFromConfigWithNamespace(cloudCfg, namespace)
	if err != nil {
		return err
	}

	deploy.Register(ctx, apply, mgmt.Management().V3().AgentDeployment())

	return start.All(ctx, 5, mgmt)
}

func getCloudCfg(token, cloudURL string) (*rest.Config, error) {
	u, err := url.Parse(cloudURL)
	if err != nil {
		return nil, err
	}
	u.Path = "/"

	return &rest.Config{
		Host:        u.String(),
		BearerToken: token,
	}, nil
}
