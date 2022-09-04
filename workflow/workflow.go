package workflow

import (
	"bytes"
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	kubernetes "k8s.io/client-go/kubernetes"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/restmapper"
)

type Operation string

const (
	Deploy   Operation = "deploy"
	Replace  Operation = "replace"
	Displace Operation = "displace"
)

type Interface interface {
	Apply(manifest []byte, op Operation) error
}

type Workflow struct {
	ctx       context.Context
	api       dynamic.Interface
	clientset kubernetes.Interface
}

func NewWorkflow(ctx context.Context, api dynamic.Interface, clientset kubernetes.Interface) *Workflow {
	return &Workflow{
		ctx:       ctx,
		api:       api,
		clientset: clientset,
	}
}

func (wf *Workflow) Apply(manifest []byte, op Operation) error {
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(manifest), 4096)
	for {
		var (
			rawObj runtime.RawExtension
			dri    dynamic.ResourceInterface
		)

		if err := decoder.Decode(&rawObj); err != nil {
			break
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return err
		}

		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(wf.clientset.Discovery())
		if err != nil {
			return err
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}

		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace("default")
			}
			dri = wf.api.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = wf.api.Resource(mapping.Resource)
		}

		switch op {
		case Deploy:
			if _, err := dri.Create(wf.ctx, unstructuredObj, metav1.CreateOptions{}); err != nil {
				return err
			}

		case Replace:
			if _, err := dri.Update(wf.ctx, unstructuredObj, metav1.UpdateOptions{}); err != nil {
				return err
			}

		case Displace:
			if err := dri.Delete(wf.ctx, unstructuredObj.GetName(), metav1.DeleteOptions{}); err != nil {
				return err
			}
		}
	}

	return nil
}
