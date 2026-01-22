package plugin

import (
	"context"
	"errors"

	"github.com/argoproj-labs/rollouts-plugin-trafficrouter-gatewayapi/internal/defaults"
	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	pluginTypes "github.com/argoproj/argo-rollouts/utils/plugin/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *RpcPlugin) setTLSRouteWeight(rollout *v1alpha1.Rollout, desiredWeight int32, gatewayAPIConfig *GatewayAPITrafficRouting) pluginTypes.RpcError {
	ctx := context.TODO()
	tlsRouteClient := r.TLSRouteClient
	if !r.IsTest {
		gatewayClientV1alpha2 := r.GatewayAPIClientset.GatewayV1alpha2()
		tlsRouteClient = gatewayClientV1alpha2.TLSRoutes(gatewayAPIConfig.Namespace)
	}
	tlsRoute, err := tlsRouteClient.Get(ctx, gatewayAPIConfig.TLSRoute, metav1.GetOptions{})
	if err != nil {
		return pluginTypes.RpcError{
			ErrorString: err.Error(),
		}
	}
	canaryServiceName := rollout.Spec.Strategy.Canary.CanaryService
	stableServiceName := rollout.Spec.Strategy.Canary.StableService
	stableWeight := 100 - desiredWeight

	// Modify labels in-place using existing logic, then copy to apply config
	ensureInProgressLabel(tlsRoute, desiredWeight, gatewayAPIConfig)
	applyConfig := buildTLSRouteApply(
		tlsRoute.Name,
		tlsRoute.Namespace,
		tlsRoute.Spec.Rules,
		canaryServiceName,
		stableServiceName,
		desiredWeight,
		stableWeight,
		tlsRoute.Labels,
	)

	updatedTLSRoute, err := tlsRouteClient.Apply(ctx, applyConfig, metav1.ApplyOptions{
		FieldManager: defaults.FieldManager,
		Force:        true,
	})
	if r.IsTest {
		r.UpdatedTLSRouteMock = updatedTLSRoute
	}
	if err != nil {
		return pluginTypes.RpcError{
			ErrorString: err.Error(),
		}
	}
	return pluginTypes.RpcError{}
}

func (r *TLSRouteRule) Iterator() (GatewayAPIRouteRuleIterator[*TLSBackendRef], bool) {
	backendRefList := r.BackendRefs
	index := 0
	next := func() (*TLSBackendRef, bool) {
		if len(backendRefList) == index {
			return nil, false
		}
		backendRef := (*TLSBackendRef)(&backendRefList[index])
		index = index + 1
		return backendRef, len(backendRefList) > index
	}
	return next, len(backendRefList) > index
}

func (r TLSRouteRuleList) Iterator() (GatewayAPIRouteRuleListIterator[*TLSBackendRef, *TLSRouteRule], bool) {
	routeRuleList := r
	index := 0
	next := func() (*TLSRouteRule, bool) {
		if len(routeRuleList) == index {
			return nil, false
		}
		routeRule := (*TLSRouteRule)(&routeRuleList[index])
		index = index + 1
		return routeRule, len(routeRuleList) > index
	}
	return next, len(routeRuleList) > index
}

func (r TLSRouteRuleList) Error() error {
	return errors.New(BackendRefListWasNotFoundInTLSRouteError)
}

func (r *TLSBackendRef) GetName() string {
	return string(r.Name)
}

func (r TLSRoute) GetName() string {
	return r.Name
}
