// +build !windows,foo

package cluster

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/rancher/dynamiclistener/server"
	"github.com/rancher/rancher/manager/modules/steve"
	"github.com/rancher/rancher/manager/pkg/ccontext"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

const (
	nsFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

func runSteve(ctx context.Context, webhookURL string) error {
	logrus.Info("Starting steve")
	c, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	ns, err := ioutil.ReadFile(nsFile)
	if err != nil {
		return err
	}

	//u, err := url.Parse(webhookURL)
	//if err != nil {
	//	return err
	//}

	cfg := ccontext.WranglerControllerContextConfig{
		Namespace:     strings.TrimSpace(string(ns)),
		ClusterConfig: c,
		AuthConfig:    c,
		//WebhookConfig: cli.WebhookConfig{
		//	WebhookAuthentication: true,
		//	WebhookURL:            fmt.Sprintf("https://%s/v3/tokenreview", u.Host),
		//},
	}

	config, err := ccontext.NewControllerContext(cfg)
	if err != nil {
		return err
	}

	if err := steve.Register(ctx, config); err != nil {
		return err
	}

	if err := config.Start(ctx); err != nil {
		return err
	}

	return server.ListenAndServe(ctx, 8443, 8080, config.Handler)
}
