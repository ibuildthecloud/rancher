package deploy

import (
	"context"
	"encoding/json"

	controllerv1 "github.com/rancher/rancher/pkg/wrangler/generated/controllers/management.cattle.io/v3"
	v3 "github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Register(ctx context.Context, apply apply.Apply, deployment controllerv1.AgentDeploymentController) {
	handler := &handler{}

	controllerv1.RegisterAgentDeploymentGeneratingHandler(ctx,
		deployment,
		apply,
		condition.Cond(v3.AgentDeploymentConditionDeployed),
		"agent-deploy",
		handler.Handle,
		&generic.GeneratingHandlerOptions{
			DynamicLookup: true,
		})
}

type handler struct {
}

func (h *handler) Handle(agentDep *v3.AgentDeployment, status v3.AgentDeploymentStatus) ([]runtime.Object, v3.AgentDeploymentStatus, error) {
	var objs []runtime.Object
	for i, obj := range agentDep.Spec.Objects {
		if obj.Object != nil {
			objs = append(objs, obj.Object)
		} else {
			var u unstructured.Unstructured
			if err := json.Unmarshal(obj.Raw, &u); err != nil {
				return nil, status, err
			}

			agentDep.Spec.Objects[i].Object = &u
			objs = append(objs, &u)
		}
	}

	status.AppliedHash = agentDep.Spec.Hash
	return objs, status, nil
}
