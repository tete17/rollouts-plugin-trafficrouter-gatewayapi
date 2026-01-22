package plugin

import (
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
	applyconfigv1 "sigs.k8s.io/gateway-api/applyconfiguration/apis/v1"
	applyconfigv1alpha2 "sigs.k8s.io/gateway-api/applyconfiguration/apis/v1alpha2"
)

// buildHTTPRouteApply creates an HTTPRouteApplyConfiguration for setting backend weights.
// It rebuilds the spec.rules with updated backend ref weights for the canary and stable services.
func buildHTTPRouteApply(
	name, namespace string,
	rules []gatewayv1.HTTPRouteRule,
	canaryService, stableService string,
	canaryWeight, stableWeight int32,
	labels map[string]string,
) *applyconfigv1.HTTPRouteApplyConfiguration {
	applyConfig := applyconfigv1.HTTPRoute(name, namespace)

	if labels != nil {
		applyConfig.WithLabels(labels)
	}

	var ruleConfigs []*applyconfigv1.HTTPRouteRuleApplyConfiguration
	for _, rule := range rules {
		ruleConfig := buildHTTPRouteRuleApply(rule, canaryService, stableService, canaryWeight, stableWeight)
		ruleConfigs = append(ruleConfigs, ruleConfig)
	}

	applyConfig.WithSpec(applyconfigv1.HTTPRouteSpec().WithRules(ruleConfigs...))
	return applyConfig
}

// buildHTTPRouteRuleApply creates an HTTPRouteRuleApplyConfiguration with updated weights.
func buildHTTPRouteRuleApply(
	rule gatewayv1.HTTPRouteRule,
	canaryService, stableService string,
	canaryWeight, stableWeight int32,
) *applyconfigv1.HTTPRouteRuleApplyConfiguration {
	ruleConfig := applyconfigv1.HTTPRouteRule()

	if rule.Name != nil {
		ruleConfig.WithName(*rule.Name)
	}

	for _, backendRef := range rule.BackendRefs {
		backendConfig := buildHTTPBackendRefApply(backendRef, canaryService, stableService, canaryWeight, stableWeight)
		ruleConfig.WithBackendRefs(backendConfig)
	}

	return ruleConfig
}

// buildHTTPBackendRefApply creates an HTTPBackendRefApplyConfiguration with the appropriate weight.
func buildHTTPBackendRefApply(
	backendRef gatewayv1.HTTPBackendRef,
	canaryService, stableService string,
	canaryWeight, stableWeight int32,
) *applyconfigv1.HTTPBackendRefApplyConfiguration {
	backendConfig := applyconfigv1.HTTPBackendRef()

	if backendRef.Group != nil {
		backendConfig.WithGroup(*backendRef.Group)
	}
	if backendRef.Kind != nil {
		backendConfig.WithKind(*backendRef.Kind)
	}
	backendConfig.WithName(backendRef.Name)
	if backendRef.Namespace != nil {
		backendConfig.WithNamespace(*backendRef.Namespace)
	}
	if backendRef.Port != nil {
		backendConfig.WithPort(int32(*backendRef.Port))
	}

	// Set weight based on service name
	serviceName := string(backendRef.Name)
	switch serviceName {
	case canaryService:
		backendConfig.WithWeight(canaryWeight)
	case stableService:
		backendConfig.WithWeight(stableWeight)
	default:
		// Preserve original weight for other backends (e.g., experiment services)
		if backendRef.Weight != nil {
			backendConfig.WithWeight(*backendRef.Weight)
		}
	}

	return backendConfig
}

// buildGRPCRouteApply creates a GRPCRouteApplyConfiguration for setting backend weights.
func buildGRPCRouteApply(
	name, namespace string,
	rules []gatewayv1.GRPCRouteRule,
	canaryService, stableService string,
	canaryWeight, stableWeight int32,
	labels map[string]string,
) *applyconfigv1.GRPCRouteApplyConfiguration {
	applyConfig := applyconfigv1.GRPCRoute(name, namespace)

	if labels != nil {
		applyConfig.WithLabels(labels)
	}

	var ruleConfigs []*applyconfigv1.GRPCRouteRuleApplyConfiguration
	for _, rule := range rules {
		ruleConfig := buildGRPCRouteRuleApply(rule, canaryService, stableService, canaryWeight, stableWeight)
		ruleConfigs = append(ruleConfigs, ruleConfig)
	}

	applyConfig.WithSpec(applyconfigv1.GRPCRouteSpec().WithRules(ruleConfigs...))
	return applyConfig
}

// buildGRPCRouteRuleApply creates a GRPCRouteRuleApplyConfiguration with updated weights.
func buildGRPCRouteRuleApply(
	rule gatewayv1.GRPCRouteRule,
	canaryService, stableService string,
	canaryWeight, stableWeight int32,
) *applyconfigv1.GRPCRouteRuleApplyConfiguration {
	ruleConfig := applyconfigv1.GRPCRouteRule()

	if rule.Name != nil {
		ruleConfig.WithName(*rule.Name)
	}

	for _, backendRef := range rule.BackendRefs {
		backendConfig := buildGRPCBackendRefApply(backendRef, canaryService, stableService, canaryWeight, stableWeight)
		ruleConfig.WithBackendRefs(backendConfig)
	}

	return ruleConfig
}

// buildGRPCBackendRefApply creates a GRPCBackendRefApplyConfiguration with the appropriate weight.
func buildGRPCBackendRefApply(
	backendRef gatewayv1.GRPCBackendRef,
	canaryService, stableService string,
	canaryWeight, stableWeight int32,
) *applyconfigv1.GRPCBackendRefApplyConfiguration {
	backendConfig := applyconfigv1.GRPCBackendRef()

	if backendRef.Group != nil {
		backendConfig.WithGroup(*backendRef.Group)
	}
	if backendRef.Kind != nil {
		backendConfig.WithKind(*backendRef.Kind)
	}
	backendConfig.WithName(backendRef.Name)
	if backendRef.Namespace != nil {
		backendConfig.WithNamespace(*backendRef.Namespace)
	}
	if backendRef.Port != nil {
		backendConfig.WithPort(int32(*backendRef.Port))
	}

	serviceName := string(backendRef.Name)
	switch serviceName {
	case canaryService:
		backendConfig.WithWeight(canaryWeight)
	case stableService:
		backendConfig.WithWeight(stableWeight)
	default:
		if backendRef.Weight != nil {
			backendConfig.WithWeight(*backendRef.Weight)
		}
	}

	return backendConfig
}

// buildTCPRouteApply creates a TCPRouteApplyConfiguration for setting backend weights.
func buildTCPRouteApply(
	name, namespace string,
	rules []v1alpha2.TCPRouteRule,
	canaryService, stableService string,
	canaryWeight, stableWeight int32,
	labels map[string]string,
) *applyconfigv1alpha2.TCPRouteApplyConfiguration {
	applyConfig := applyconfigv1alpha2.TCPRoute(name, namespace)

	if labels != nil {
		applyConfig.WithLabels(labels)
	}

	var ruleConfigs []*applyconfigv1alpha2.TCPRouteRuleApplyConfiguration
	for _, rule := range rules {
		ruleConfig := buildTCPRouteRuleApply(rule, canaryService, stableService, canaryWeight, stableWeight)
		ruleConfigs = append(ruleConfigs, ruleConfig)
	}

	applyConfig.WithSpec(applyconfigv1alpha2.TCPRouteSpec().WithRules(ruleConfigs...))
	return applyConfig
}

// buildTCPRouteRuleApply creates a TCPRouteRuleApplyConfiguration with updated weights.
func buildTCPRouteRuleApply(
	rule v1alpha2.TCPRouteRule,
	canaryService, stableService string,
	canaryWeight, stableWeight int32,
) *applyconfigv1alpha2.TCPRouteRuleApplyConfiguration {
	ruleConfig := applyconfigv1alpha2.TCPRouteRule()

	if rule.Name != nil {
		ruleConfig.WithName(*rule.Name)
	}

	for _, backendRef := range rule.BackendRefs {
		backendConfig := buildBackendRefApply(backendRef, canaryService, stableService, canaryWeight, stableWeight)
		ruleConfig.WithBackendRefs(backendConfig)
	}

	return ruleConfig
}

// buildTLSRouteApply creates a TLSRouteApplyConfiguration for setting backend weights.
func buildTLSRouteApply(
	name, namespace string,
	rules []v1alpha2.TLSRouteRule,
	canaryService, stableService string,
	canaryWeight, stableWeight int32,
	labels map[string]string,
) *applyconfigv1alpha2.TLSRouteApplyConfiguration {
	applyConfig := applyconfigv1alpha2.TLSRoute(name, namespace)

	if labels != nil {
		applyConfig.WithLabels(labels)
	}

	var ruleConfigs []*applyconfigv1alpha2.TLSRouteRuleApplyConfiguration
	for _, rule := range rules {
		ruleConfig := buildTLSRouteRuleApply(rule, canaryService, stableService, canaryWeight, stableWeight)
		ruleConfigs = append(ruleConfigs, ruleConfig)
	}

	applyConfig.WithSpec(applyconfigv1alpha2.TLSRouteSpec().WithRules(ruleConfigs...))
	return applyConfig
}

// buildTLSRouteRuleApply creates a TLSRouteRuleApplyConfiguration with updated weights.
func buildTLSRouteRuleApply(
	rule v1alpha2.TLSRouteRule,
	canaryService, stableService string,
	canaryWeight, stableWeight int32,
) *applyconfigv1alpha2.TLSRouteRuleApplyConfiguration {
	ruleConfig := applyconfigv1alpha2.TLSRouteRule()

	if rule.Name != nil {
		ruleConfig.WithName(*rule.Name)
	}

	for _, backendRef := range rule.BackendRefs {
		backendConfig := buildBackendRefApply(backendRef, canaryService, stableService, canaryWeight, stableWeight)
		ruleConfig.WithBackendRefs(backendConfig)
	}

	return ruleConfig
}

// buildBackendRefApply creates a BackendRefApplyConfiguration (used by TCP/TLS routes).
func buildBackendRefApply(
	backendRef gatewayv1.BackendRef,
	canaryService, stableService string,
	canaryWeight, stableWeight int32,
) *applyconfigv1.BackendRefApplyConfiguration {
	backendConfig := applyconfigv1.BackendRef()

	if backendRef.Group != nil {
		backendConfig.WithGroup(*backendRef.Group)
	}
	if backendRef.Kind != nil {
		backendConfig.WithKind(*backendRef.Kind)
	}
	backendConfig.WithName(backendRef.Name)
	if backendRef.Namespace != nil {
		backendConfig.WithNamespace(*backendRef.Namespace)
	}
	if backendRef.Port != nil {
		backendConfig.WithPort(int32(*backendRef.Port))
	}

	serviceName := string(backendRef.Name)
	switch serviceName {
	case canaryService:
		backendConfig.WithWeight(canaryWeight)
	case stableService:
		backendConfig.WithWeight(stableWeight)
	default:
		if backendRef.Weight != nil {
			backendConfig.WithWeight(*backendRef.Weight)
		}
	}

	return backendConfig
}
