package plugin

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/argoproj-labs/rollouts-plugin-trafficrouter-gatewayapi/internal/defaults"
)

// buildInProgressLabels returns a map containing the in-progress label when desiredWeight > 0,
// or nil when the label should be removed (desiredWeight == 0).
// This function is designed for use with Server-Side Apply (SSA).
func buildInProgressLabels(desiredWeight int32, config *GatewayAPITrafficRouting) map[string]string {
	if config == nil || config.DisableInProgressLabel {
		return nil
	}

	key := config.inProgressLabelKey()
	if key == "" {
		return nil
	}

	if desiredWeight == 0 {
		// Return empty map to remove the label from our field manager's ownership
		return map[string]string{}
	}

	return map[string]string{
		key: config.inProgressLabelValue(),
	}
}

func ensureInProgressLabel(obj metav1.Object, desiredWeight int32, config *GatewayAPITrafficRouting) bool {
	if obj == nil || config == nil || config.DisableInProgressLabel {
		return false
	}

	key := config.inProgressLabelKey()
	if key == "" {
		return false
	}

	labels := obj.GetLabels()
	if desiredWeight == 0 {
		if labels == nil {
			return false
		}
		if _, ok := labels[key]; ok {
			delete(labels, key)
			obj.SetLabels(labels)
			return true
		}
		return false
	}

	value := config.inProgressLabelValue()
	if labels == nil {
		labels = make(map[string]string)
	}
	if current, ok := labels[key]; ok && current == value {
		return false
	}
	labels[key] = value
	obj.SetLabels(labels)
	return true
}

func (c *GatewayAPITrafficRouting) inProgressLabelKey() string {
	if c.InProgressLabelKey != "" {
		return c.InProgressLabelKey
	}
	return defaults.InProgressLabelKey
}

func (c *GatewayAPITrafficRouting) inProgressLabelValue() string {
	if c.InProgressLabelValue != "" {
		return c.InProgressLabelValue
	}
	return defaults.InProgressLabelValue
}
