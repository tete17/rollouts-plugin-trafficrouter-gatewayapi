package mocks

import (
	"encoding/json"

	"github.com/argoproj-labs/rollouts-plugin-trafficrouter-gatewayapi/internal/defaults"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	k8stesting "k8s.io/client-go/testing"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"
	gwFake "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned/fake"
)

// NewFakeClientsetWithSSA creates a fake Gateway API clientset with proper SSA support.
// The standard fake clientset doesn't properly merge Apply patches, so this adds
// reactors that simulate SSA merge semantics for HTTPRoute, GRPCRoute, TCPRoute, and TLSRoute.
func NewFakeClientsetWithSSA(objects ...runtime.Object) *gwFake.Clientset {
	cs := gwFake.NewSimpleClientset(objects...)

	// Add reactor for HTTPRoute Apply
	cs.PrependReactor("patch", "httproutes", func(action k8stesting.Action) (bool, runtime.Object, error) {
		patchAction := action.(k8stesting.PatchAction)
		if patchAction.GetPatchType() != types.ApplyPatchType {
			return false, nil, nil
		}
		return handleHTTPRouteApply(cs, patchAction)
	})

	// Add reactor for GRPCRoute Apply
	cs.PrependReactor("patch", "grpcroutes", func(action k8stesting.Action) (bool, runtime.Object, error) {
		patchAction := action.(k8stesting.PatchAction)
		if patchAction.GetPatchType() != types.ApplyPatchType {
			return false, nil, nil
		}
		return handleGRPCRouteApply(cs, patchAction)
	})

	// Add reactor for TCPRoute Apply
	cs.PrependReactor("patch", "tcproutes", func(action k8stesting.Action) (bool, runtime.Object, error) {
		patchAction := action.(k8stesting.PatchAction)
		if patchAction.GetPatchType() != types.ApplyPatchType {
			return false, nil, nil
		}
		return handleTCPRouteApply(cs, patchAction)
	})

	// Add reactor for TLSRoute Apply
	cs.PrependReactor("patch", "tlsroutes", func(action k8stesting.Action) (bool, runtime.Object, error) {
		patchAction := action.(k8stesting.PatchAction)
		if patchAction.GetPatchType() != types.ApplyPatchType {
			return false, nil, nil
		}
		return handleTLSRouteApply(cs, patchAction)
	})

	return cs
}

func handleHTTPRouteApply(cs *gwFake.Clientset, action k8stesting.PatchAction) (bool, runtime.Object, error) {
	// Get existing object
	existing, err := cs.Tracker().Get(
		gatewayv1.SchemeGroupVersion.WithResource("httproutes"),
		action.GetNamespace(),
		action.GetName(),
	)
	if err != nil {
		return true, nil, err
	}

	route := existing.(*gatewayv1.HTTPRoute).DeepCopy()

	// Parse patch and merge
	var patch map[string]interface{}
	if err := json.Unmarshal(action.GetPatch(), &patch); err != nil {
		return true, nil, err
	}

	// Merge labels
	// In SSA, if the patch includes labels, we merge them.
	// If the patch doesn't include labels at all (the field is absent), we simulate
	// the expected behavior for our managed label: remove it if it exists.
	if metadata, ok := patch["metadata"].(map[string]interface{}); ok {
		if labels, ok := metadata["labels"].(map[string]interface{}); ok {
			if route.Labels == nil {
				route.Labels = make(map[string]string)
			}
			// For SSA: if labels map is present but empty, remove all managed labels
			if len(labels) == 0 {
				route.Labels = map[string]string{}
			} else {
				for k, v := range labels {
					if v == nil {
						delete(route.Labels, k)
					} else {
						route.Labels[k] = v.(string)
					}
				}
			}
		} else {
			// Labels field not included in patch - simulate SSA behavior:
			// Remove the managed in-progress label if it exists
			if route.Labels != nil {
				delete(route.Labels, defaults.InProgressLabelKey)
			}
		}
	} else {
		// No metadata in patch - remove managed label
		if route.Labels != nil {
			delete(route.Labels, defaults.InProgressLabelKey)
		}
	}

	// Merge spec.rules
	if spec, ok := patch["spec"].(map[string]interface{}); ok {
		if rules, ok := spec["rules"].([]interface{}); ok {
			route.Spec.Rules = mergeHTTPRouteRules(route.Spec.Rules, rules)
		}
	}

	// Update tracker
	if err := cs.Tracker().Update(
		gatewayv1.SchemeGroupVersion.WithResource("httproutes"),
		route,
		action.GetNamespace(),
	); err != nil {
		return true, nil, err
	}

	return true, route, nil
}

func handleGRPCRouteApply(cs *gwFake.Clientset, action k8stesting.PatchAction) (bool, runtime.Object, error) {
	existing, err := cs.Tracker().Get(
		gatewayv1.SchemeGroupVersion.WithResource("grpcroutes"),
		action.GetNamespace(),
		action.GetName(),
	)
	if err != nil {
		return true, nil, err
	}

	route := existing.(*gatewayv1.GRPCRoute).DeepCopy()

	var patch map[string]interface{}
	if err := json.Unmarshal(action.GetPatch(), &patch); err != nil {
		return true, nil, err
	}

	// Merge labels (same logic as HTTPRoute)
	if metadata, ok := patch["metadata"].(map[string]interface{}); ok {
		if labels, ok := metadata["labels"].(map[string]interface{}); ok {
			if route.Labels == nil {
				route.Labels = make(map[string]string)
			}
			if len(labels) == 0 {
				route.Labels = map[string]string{}
			} else {
				for k, v := range labels {
					if v == nil {
						delete(route.Labels, k)
					} else {
						route.Labels[k] = v.(string)
					}
				}
			}
		} else {
			if route.Labels != nil {
				delete(route.Labels, defaults.InProgressLabelKey)
			}
		}
	} else {
		if route.Labels != nil {
			delete(route.Labels, defaults.InProgressLabelKey)
		}
	}

	// Merge spec.rules
	if spec, ok := patch["spec"].(map[string]interface{}); ok {
		if rules, ok := spec["rules"].([]interface{}); ok {
			route.Spec.Rules = mergeGRPCRouteRules(route.Spec.Rules, rules)
		}
	}

	if err := cs.Tracker().Update(
		gatewayv1.SchemeGroupVersion.WithResource("grpcroutes"),
		route,
		action.GetNamespace(),
	); err != nil {
		return true, nil, err
	}

	return true, route, nil
}

func handleTCPRouteApply(cs *gwFake.Clientset, action k8stesting.PatchAction) (bool, runtime.Object, error) {
	existing, err := cs.Tracker().Get(
		v1alpha2.SchemeGroupVersion.WithResource("tcproutes"),
		action.GetNamespace(),
		action.GetName(),
	)
	if err != nil {
		return true, nil, err
	}

	route := existing.(*v1alpha2.TCPRoute).DeepCopy()

	var patch map[string]interface{}
	if err := json.Unmarshal(action.GetPatch(), &patch); err != nil {
		return true, nil, err
	}

	// Merge labels (same logic as HTTPRoute)
	if metadata, ok := patch["metadata"].(map[string]interface{}); ok {
		if labels, ok := metadata["labels"].(map[string]interface{}); ok {
			if route.Labels == nil {
				route.Labels = make(map[string]string)
			}
			if len(labels) == 0 {
				route.Labels = map[string]string{}
			} else {
				for k, v := range labels {
					if v == nil {
						delete(route.Labels, k)
					} else {
						route.Labels[k] = v.(string)
					}
				}
			}
		} else {
			if route.Labels != nil {
				delete(route.Labels, defaults.InProgressLabelKey)
			}
		}
	} else {
		if route.Labels != nil {
			delete(route.Labels, defaults.InProgressLabelKey)
		}
	}

	// Merge spec.rules
	if spec, ok := patch["spec"].(map[string]interface{}); ok {
		if rules, ok := spec["rules"].([]interface{}); ok {
			route.Spec.Rules = mergeTCPRouteRules(route.Spec.Rules, rules)
		}
	}

	if err := cs.Tracker().Update(
		v1alpha2.SchemeGroupVersion.WithResource("tcproutes"),
		route,
		action.GetNamespace(),
	); err != nil {
		return true, nil, err
	}

	return true, route, nil
}

func handleTLSRouteApply(cs *gwFake.Clientset, action k8stesting.PatchAction) (bool, runtime.Object, error) {
	existing, err := cs.Tracker().Get(
		v1alpha2.SchemeGroupVersion.WithResource("tlsroutes"),
		action.GetNamespace(),
		action.GetName(),
	)
	if err != nil {
		return true, nil, err
	}

	route := existing.(*v1alpha2.TLSRoute).DeepCopy()

	var patch map[string]interface{}
	if err := json.Unmarshal(action.GetPatch(), &patch); err != nil {
		return true, nil, err
	}

	// Merge labels (same logic as HTTPRoute)
	if metadata, ok := patch["metadata"].(map[string]interface{}); ok {
		if labels, ok := metadata["labels"].(map[string]interface{}); ok {
			if route.Labels == nil {
				route.Labels = make(map[string]string)
			}
			if len(labels) == 0 {
				route.Labels = map[string]string{}
			} else {
				for k, v := range labels {
					if v == nil {
						delete(route.Labels, k)
					} else {
						route.Labels[k] = v.(string)
					}
				}
			}
		} else {
			if route.Labels != nil {
				delete(route.Labels, defaults.InProgressLabelKey)
			}
		}
	} else {
		if route.Labels != nil {
			delete(route.Labels, defaults.InProgressLabelKey)
		}
	}

	// Merge spec.rules
	if spec, ok := patch["spec"].(map[string]interface{}); ok {
		if rules, ok := spec["rules"].([]interface{}); ok {
			route.Spec.Rules = mergeTLSRouteRules(route.Spec.Rules, rules)
		}
	}

	if err := cs.Tracker().Update(
		v1alpha2.SchemeGroupVersion.WithResource("tlsroutes"),
		route,
		action.GetNamespace(),
	); err != nil {
		return true, nil, err
	}

	return true, route, nil
}

// mergeHTTPRouteRules replaces existing rules with patch rules.
// For SSA, the entire rules array is managed, so we replace it.
func mergeHTTPRouteRules(existing []gatewayv1.HTTPRouteRule, patchRules []interface{}) []gatewayv1.HTTPRouteRule {
	if len(patchRules) == 0 {
		return existing
	}

	// Convert patch rules to HTTPRouteRule
	rulesJSON, _ := json.Marshal(patchRules)
	var rules []gatewayv1.HTTPRouteRule
	if err := json.Unmarshal(rulesJSON, &rules); err != nil {
		return existing
	}
	return rules
}

func mergeGRPCRouteRules(existing []gatewayv1.GRPCRouteRule, patchRules []interface{}) []gatewayv1.GRPCRouteRule {
	if len(patchRules) == 0 {
		return existing
	}

	rulesJSON, _ := json.Marshal(patchRules)
	var rules []gatewayv1.GRPCRouteRule
	if err := json.Unmarshal(rulesJSON, &rules); err != nil {
		return existing
	}
	return rules
}

func mergeTCPRouteRules(existing []v1alpha2.TCPRouteRule, patchRules []interface{}) []v1alpha2.TCPRouteRule {
	if len(patchRules) == 0 {
		return existing
	}

	rulesJSON, _ := json.Marshal(patchRules)
	var rules []v1alpha2.TCPRouteRule
	if err := json.Unmarshal(rulesJSON, &rules); err != nil {
		return existing
	}
	return rules
}

func mergeTLSRouteRules(existing []v1alpha2.TLSRouteRule, patchRules []interface{}) []v1alpha2.TLSRouteRule {
	if len(patchRules) == 0 {
		return existing
	}

	rulesJSON, _ := json.Marshal(patchRules)
	var rules []v1alpha2.TLSRouteRule
	if err := json.Unmarshal(rulesJSON, &rules); err != nil {
		return existing
	}
	return rules
}
