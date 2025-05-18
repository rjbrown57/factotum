package handlers

import (
	"testing"

	"github.com/rjbrown57/factotum/api/v1alpha1"
	"github.com/rjbrown57/factotum/pkg/factotum/config"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestMetaDataHandler_Update(t *testing.T) {
	handler := &MetaDataHandler{}

	t.Run("Valid Object with Annotations and Labels", func(t *testing.T) {
		obj := &v1.Node{}
		obj.SetAnnotations(map[string]string{"key1": "value1"})
		obj.SetLabels(map[string]string{"label1": "value1"})

		nodeConfig := &v1alpha1.NodeConfig{
			Spec: v1alpha1.NodeConfigSpec{
				CommonSpec: config.CommonSpec{
					Labels: map[string]string{
						"label1": "newValue1",
						"label2": "value2",
					},
					Annotations: map[string]string{
						"key1": "newValue1",
						"key2": "value2",
					},
				},
			},
		}

		updatedObj := handler.Update(obj, nodeConfig)
		assert.NotNil(t, updatedObj)

		updatedAnnotations := updatedObj.(*v1.Node).GetAnnotations()
		updatedLabels := updatedObj.(*v1.Node).GetLabels()

		assert.Equal(t, "newValue1", updatedAnnotations["key1"])
		assert.Equal(t, "value2", updatedAnnotations["key2"])
		assert.Equal(t, "newValue1", updatedLabels["label1"])
		assert.Equal(t, "value2", updatedLabels["label2"])
	})
}
