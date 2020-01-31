package management

import (
	"context"

	"github.com/rancher/rancher/pkg/wrangler"
	"github.com/rancher/types/config"
)

func wranglerControllers(ctx context.Context, managementContext *config.ManagementContext) error {
	w, err := wrangler.NewContext(&managementContext.RESTConfig)
	if err != nil {
		return err
	}

	// Add controllers to register here

	return w.Start(ctx)
}
