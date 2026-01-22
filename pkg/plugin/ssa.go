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

// buildHTTPRouteApplyFromRules creates an HTTPRouteApplyConfiguration from already-modified rules.
// This is useful when the rules have been modified in place (e.g., by experiment handling)
// and we want to preserve those modifications.
func buildHTTPRouteApplyFromRules(
	name, namespace string,
	rules []gatewayv1.HTTPRouteRule,
	labels map[string]string,
) *applyconfigv1.HTTPRouteApplyConfiguration {
	applyConfig := applyconfigv1.HTTPRoute(name, namespace)

	if labels != nil {
		applyConfig.WithLabels(labels)
	}

	var ruleConfigs []*applyconfigv1.HTTPRouteRuleApplyConfiguration
	for _, rule := range rules {
		ruleConfig := buildHTTPRouteRuleApplyFromRule(rule)
		ruleConfigs = append(ruleConfigs, ruleConfig)
	}

	applyConfig.WithSpec(applyconfigv1.HTTPRouteSpec().WithRules(ruleConfigs...))
	return applyConfig
}

// buildHTTPRouteRuleApplyFromRule creates an HTTPRouteRuleApplyConfiguration preserving all fields.
func buildHTTPRouteRuleApplyFromRule(rule gatewayv1.HTTPRouteRule) *applyconfigv1.HTTPRouteRuleApplyConfiguration {
	ruleConfig := applyconfigv1.HTTPRouteRule()

	if rule.Name != nil {
		ruleConfig.WithName(*rule.Name)
	}

	// Copy matches
	for _, match := range rule.Matches {
		matchConfig := buildHTTPRouteMatchApply(match)
		ruleConfig.WithMatches(matchConfig)
	}

	// Copy filters
	for _, filter := range rule.Filters {
		filterConfig := buildHTTPRouteFilterApply(filter)
		ruleConfig.WithFilters(filterConfig)
	}

	// Copy backend refs
	for _, backendRef := range rule.BackendRefs {
		backendConfig := buildHTTPBackendRefApplyFromRef(backendRef)
		ruleConfig.WithBackendRefs(backendConfig)
	}

	// Copy timeouts
	if rule.Timeouts != nil {
		timeoutsConfig := applyconfigv1.HTTPRouteTimeouts()
		if rule.Timeouts.Request != nil {
			timeoutsConfig.WithRequest(*rule.Timeouts.Request)
		}
		if rule.Timeouts.BackendRequest != nil {
			timeoutsConfig.WithBackendRequest(*rule.Timeouts.BackendRequest)
		}
		ruleConfig.WithTimeouts(timeoutsConfig)
	}

	return ruleConfig
}

// buildHTTPRouteMatchApply creates an HTTPRouteMatchApplyConfiguration from an HTTPRouteMatch.
func buildHTTPRouteMatchApply(match gatewayv1.HTTPRouteMatch) *applyconfigv1.HTTPRouteMatchApplyConfiguration {
	matchConfig := applyconfigv1.HTTPRouteMatch()

	if match.Path != nil {
		pathConfig := applyconfigv1.HTTPPathMatch()
		if match.Path.Type != nil {
			pathConfig.WithType(*match.Path.Type)
		}
		if match.Path.Value != nil {
			pathConfig.WithValue(*match.Path.Value)
		}
		matchConfig.WithPath(pathConfig)
	}

	for _, header := range match.Headers {
		headerConfig := applyconfigv1.HTTPHeaderMatch().
			WithName(header.Name)
		if header.Type != nil {
			headerConfig.WithType(*header.Type)
		}
		headerConfig.WithValue(header.Value)
		matchConfig.WithHeaders(headerConfig)
	}

	for _, query := range match.QueryParams {
		queryConfig := applyconfigv1.HTTPQueryParamMatch().
			WithName(query.Name).
			WithValue(query.Value)
		if query.Type != nil {
			queryConfig.WithType(*query.Type)
		}
		matchConfig.WithQueryParams(queryConfig)
	}

	if match.Method != nil {
		matchConfig.WithMethod(*match.Method)
	}

	return matchConfig
}

// buildHTTPRouteFilterApply creates an HTTPRouteFilterApplyConfiguration from an HTTPRouteFilter.
func buildHTTPRouteFilterApply(filter gatewayv1.HTTPRouteFilter) *applyconfigv1.HTTPRouteFilterApplyConfiguration {
	filterConfig := applyconfigv1.HTTPRouteFilter().
		WithType(filter.Type)

	if filter.RequestHeaderModifier != nil {
		headerModConfig := applyconfigv1.HTTPHeaderFilter()
		for _, header := range filter.RequestHeaderModifier.Set {
			headerModConfig.WithSet(applyconfigv1.HTTPHeader().WithName(header.Name).WithValue(header.Value))
		}
		for _, header := range filter.RequestHeaderModifier.Add {
			headerModConfig.WithAdd(applyconfigv1.HTTPHeader().WithName(header.Name).WithValue(header.Value))
		}
		headerModConfig.WithRemove(filter.RequestHeaderModifier.Remove...)
		filterConfig.WithRequestHeaderModifier(headerModConfig)
	}

	if filter.ResponseHeaderModifier != nil {
		headerModConfig := applyconfigv1.HTTPHeaderFilter()
		for _, header := range filter.ResponseHeaderModifier.Set {
			headerModConfig.WithSet(applyconfigv1.HTTPHeader().WithName(header.Name).WithValue(header.Value))
		}
		for _, header := range filter.ResponseHeaderModifier.Add {
			headerModConfig.WithAdd(applyconfigv1.HTTPHeader().WithName(header.Name).WithValue(header.Value))
		}
		headerModConfig.WithRemove(filter.ResponseHeaderModifier.Remove...)
		filterConfig.WithResponseHeaderModifier(headerModConfig)
	}

	if filter.RequestMirror != nil {
		mirrorConfig := applyconfigv1.HTTPRequestMirrorFilter()
		backendRefConfig := applyconfigv1.BackendObjectReference()
		if filter.RequestMirror.BackendRef.Group != nil {
			backendRefConfig.WithGroup(*filter.RequestMirror.BackendRef.Group)
		}
		if filter.RequestMirror.BackendRef.Kind != nil {
			backendRefConfig.WithKind(*filter.RequestMirror.BackendRef.Kind)
		}
		backendRefConfig.WithName(filter.RequestMirror.BackendRef.Name)
		if filter.RequestMirror.BackendRef.Namespace != nil {
			backendRefConfig.WithNamespace(*filter.RequestMirror.BackendRef.Namespace)
		}
		if filter.RequestMirror.BackendRef.Port != nil {
			backendRefConfig.WithPort(int32(*filter.RequestMirror.BackendRef.Port))
		}
		mirrorConfig.WithBackendRef(backendRefConfig)
		if filter.RequestMirror.Percent != nil {
			mirrorConfig.WithPercent(*filter.RequestMirror.Percent)
		}
		filterConfig.WithRequestMirror(mirrorConfig)
	}

	if filter.RequestRedirect != nil {
		redirectConfig := applyconfigv1.HTTPRequestRedirectFilter()
		if filter.RequestRedirect.Scheme != nil {
			redirectConfig.WithScheme(*filter.RequestRedirect.Scheme)
		}
		if filter.RequestRedirect.Hostname != nil {
			redirectConfig.WithHostname(*filter.RequestRedirect.Hostname)
		}
		if filter.RequestRedirect.Path != nil {
			pathModConfig := applyconfigv1.HTTPPathModifier().
				WithType(filter.RequestRedirect.Path.Type)
			if filter.RequestRedirect.Path.ReplaceFullPath != nil {
				pathModConfig.WithReplaceFullPath(*filter.RequestRedirect.Path.ReplaceFullPath)
			}
			if filter.RequestRedirect.Path.ReplacePrefixMatch != nil {
				pathModConfig.WithReplacePrefixMatch(*filter.RequestRedirect.Path.ReplacePrefixMatch)
			}
			redirectConfig.WithPath(pathModConfig)
		}
		if filter.RequestRedirect.Port != nil {
			redirectConfig.WithPort(int32(*filter.RequestRedirect.Port))
		}
		if filter.RequestRedirect.StatusCode != nil {
			redirectConfig.WithStatusCode(*filter.RequestRedirect.StatusCode)
		}
		filterConfig.WithRequestRedirect(redirectConfig)
	}

	if filter.URLRewrite != nil {
		rewriteConfig := applyconfigv1.HTTPURLRewriteFilter()
		if filter.URLRewrite.Hostname != nil {
			rewriteConfig.WithHostname(*filter.URLRewrite.Hostname)
		}
		if filter.URLRewrite.Path != nil {
			pathModConfig := applyconfigv1.HTTPPathModifier().
				WithType(filter.URLRewrite.Path.Type)
			if filter.URLRewrite.Path.ReplaceFullPath != nil {
				pathModConfig.WithReplaceFullPath(*filter.URLRewrite.Path.ReplaceFullPath)
			}
			if filter.URLRewrite.Path.ReplacePrefixMatch != nil {
				pathModConfig.WithReplacePrefixMatch(*filter.URLRewrite.Path.ReplacePrefixMatch)
			}
			rewriteConfig.WithPath(pathModConfig)
		}
		filterConfig.WithURLRewrite(rewriteConfig)
	}

	return filterConfig
}

// buildHTTPBackendRefApplyFromRef creates an HTTPBackendRefApplyConfiguration preserving the existing weight.
func buildHTTPBackendRefApplyFromRef(backendRef gatewayv1.HTTPBackendRef) *applyconfigv1.HTTPBackendRefApplyConfiguration {
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
	if backendRef.Weight != nil {
		backendConfig.WithWeight(*backendRef.Weight)
	}

	return backendConfig
}

// buildGRPCRouteApplyFromRules creates a GRPCRouteApplyConfiguration from already-modified rules.
func buildGRPCRouteApplyFromRules(
	name, namespace string,
	rules []gatewayv1.GRPCRouteRule,
	labels map[string]string,
) *applyconfigv1.GRPCRouteApplyConfiguration {
	applyConfig := applyconfigv1.GRPCRoute(name, namespace)

	if labels != nil {
		applyConfig.WithLabels(labels)
	}

	var ruleConfigs []*applyconfigv1.GRPCRouteRuleApplyConfiguration
	for _, rule := range rules {
		ruleConfig := buildGRPCRouteRuleApplyFromRule(rule)
		ruleConfigs = append(ruleConfigs, ruleConfig)
	}

	applyConfig.WithSpec(applyconfigv1.GRPCRouteSpec().WithRules(ruleConfigs...))
	return applyConfig
}

// buildGRPCRouteRuleApplyFromRule creates a GRPCRouteRuleApplyConfiguration preserving all fields.
func buildGRPCRouteRuleApplyFromRule(rule gatewayv1.GRPCRouteRule) *applyconfigv1.GRPCRouteRuleApplyConfiguration {
	ruleConfig := applyconfigv1.GRPCRouteRule()

	if rule.Name != nil {
		ruleConfig.WithName(*rule.Name)
	}

	// Copy matches
	for _, match := range rule.Matches {
		matchConfig := buildGRPCRouteMatchApply(match)
		ruleConfig.WithMatches(matchConfig)
	}

	// Copy filters
	for _, filter := range rule.Filters {
		filterConfig := buildGRPCRouteFilterApply(filter)
		ruleConfig.WithFilters(filterConfig)
	}

	// Copy backend refs
	for _, backendRef := range rule.BackendRefs {
		backendConfig := buildGRPCBackendRefApplyFromRef(backendRef)
		ruleConfig.WithBackendRefs(backendConfig)
	}

	return ruleConfig
}

// buildGRPCRouteMatchApply creates a GRPCRouteMatchApplyConfiguration from a GRPCRouteMatch.
func buildGRPCRouteMatchApply(match gatewayv1.GRPCRouteMatch) *applyconfigv1.GRPCRouteMatchApplyConfiguration {
	matchConfig := applyconfigv1.GRPCRouteMatch()

	if match.Method != nil {
		methodConfig := applyconfigv1.GRPCMethodMatch()
		if match.Method.Type != nil {
			methodConfig.WithType(*match.Method.Type)
		}
		if match.Method.Service != nil {
			methodConfig.WithService(*match.Method.Service)
		}
		if match.Method.Method != nil {
			methodConfig.WithMethod(*match.Method.Method)
		}
		matchConfig.WithMethod(methodConfig)
	}

	for _, header := range match.Headers {
		headerConfig := applyconfigv1.GRPCHeaderMatch().
			WithName(header.Name).
			WithValue(header.Value)
		if header.Type != nil {
			headerConfig.WithType(*header.Type)
		}
		matchConfig.WithHeaders(headerConfig)
	}

	return matchConfig
}

// buildGRPCRouteFilterApply creates a GRPCRouteFilterApplyConfiguration from a GRPCRouteFilter.
func buildGRPCRouteFilterApply(filter gatewayv1.GRPCRouteFilter) *applyconfigv1.GRPCRouteFilterApplyConfiguration {
	filterConfig := applyconfigv1.GRPCRouteFilter().
		WithType(filter.Type)

	if filter.RequestHeaderModifier != nil {
		headerModConfig := applyconfigv1.HTTPHeaderFilter()
		for _, header := range filter.RequestHeaderModifier.Set {
			headerModConfig.WithSet(applyconfigv1.HTTPHeader().WithName(header.Name).WithValue(header.Value))
		}
		for _, header := range filter.RequestHeaderModifier.Add {
			headerModConfig.WithAdd(applyconfigv1.HTTPHeader().WithName(header.Name).WithValue(header.Value))
		}
		headerModConfig.WithRemove(filter.RequestHeaderModifier.Remove...)
		filterConfig.WithRequestHeaderModifier(headerModConfig)
	}

	if filter.ResponseHeaderModifier != nil {
		headerModConfig := applyconfigv1.HTTPHeaderFilter()
		for _, header := range filter.ResponseHeaderModifier.Set {
			headerModConfig.WithSet(applyconfigv1.HTTPHeader().WithName(header.Name).WithValue(header.Value))
		}
		for _, header := range filter.ResponseHeaderModifier.Add {
			headerModConfig.WithAdd(applyconfigv1.HTTPHeader().WithName(header.Name).WithValue(header.Value))
		}
		headerModConfig.WithRemove(filter.ResponseHeaderModifier.Remove...)
		filterConfig.WithResponseHeaderModifier(headerModConfig)
	}

	if filter.RequestMirror != nil {
		mirrorConfig := applyconfigv1.HTTPRequestMirrorFilter()
		backendRefConfig := applyconfigv1.BackendObjectReference()
		if filter.RequestMirror.BackendRef.Group != nil {
			backendRefConfig.WithGroup(*filter.RequestMirror.BackendRef.Group)
		}
		if filter.RequestMirror.BackendRef.Kind != nil {
			backendRefConfig.WithKind(*filter.RequestMirror.BackendRef.Kind)
		}
		backendRefConfig.WithName(filter.RequestMirror.BackendRef.Name)
		if filter.RequestMirror.BackendRef.Namespace != nil {
			backendRefConfig.WithNamespace(*filter.RequestMirror.BackendRef.Namespace)
		}
		if filter.RequestMirror.BackendRef.Port != nil {
			backendRefConfig.WithPort(int32(*filter.RequestMirror.BackendRef.Port))
		}
		mirrorConfig.WithBackendRef(backendRefConfig)
		if filter.RequestMirror.Percent != nil {
			mirrorConfig.WithPercent(*filter.RequestMirror.Percent)
		}
		filterConfig.WithRequestMirror(mirrorConfig)
	}

	return filterConfig
}

// buildGRPCBackendRefApplyFromRef creates a GRPCBackendRefApplyConfiguration preserving the existing weight.
func buildGRPCBackendRefApplyFromRef(backendRef gatewayv1.GRPCBackendRef) *applyconfigv1.GRPCBackendRefApplyConfiguration {
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
	if backendRef.Weight != nil {
		backendConfig.WithWeight(*backendRef.Weight)
	}

	return backendConfig
}
