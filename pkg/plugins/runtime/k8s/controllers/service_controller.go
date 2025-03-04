package controllers

import (
	"context"
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	kube_core "k8s.io/api/core/v1"
	kube_apierrs "k8s.io/apimachinery/pkg/api/errors"
	kube_types "k8s.io/apimachinery/pkg/types"
	kube_ctrl "sigs.k8s.io/controller-runtime"
	kube_client "sigs.k8s.io/controller-runtime/pkg/client"

	k8scnicncfio "github.com/kumahq/kuma/pkg/plugins/runtime/k8s/apis/k8s.cni.cncf.io"
	"github.com/kumahq/kuma/pkg/plugins/runtime/k8s/metadata"
)

// ServiceReconciler reconciles a Service object
type ServiceReconciler struct {
	kube_client.Client
	Log logr.Logger
}

// Reconcile is in charge of injecting "ingress.kubernetes.io/service-upstream" annotation to the Services
// that are in Kuma enabled namespaces
func (r *ServiceReconciler) Reconcile(req kube_ctrl.Request) (kube_ctrl.Result, error) {
	log := r.Log.WithValues("service", req.NamespacedName)
	ctx := context.Background()

	svc := &kube_core.Service{}
	if err := r.Get(ctx, req.NamespacedName, svc); err != nil {
		if kube_apierrs.IsNotFound(err) {
			return kube_ctrl.Result{}, nil
		}
		return kube_ctrl.Result{}, errors.Wrapf(err, "unable to fetch Service %s", req.NamespacedName.Name)
	}

	namespace := &kube_core.Namespace{}
	if err := r.Get(ctx, kube_types.NamespacedName{Name: svc.GetNamespace()}, namespace); err != nil {
		if kube_apierrs.IsNotFound(err) {
			return kube_ctrl.Result{}, nil
		}
		return kube_ctrl.Result{}, errors.Wrapf(err, "unable to fetch Service %s", req.NamespacedName.Name)
	}

	injected, _, err := metadata.Annotations(namespace.Annotations).GetEnabled(metadata.KumaSidecarInjectionAnnotation)
	if err != nil {
		return kube_ctrl.Result{}, errors.Wrapf(err, "unable to check sidecar injection annotation on namespace %s", namespace.Name)
	}
	if !injected {
		log.V(1).Info(req.NamespacedName.String() + "is not part of the mesh")
		return kube_ctrl.Result{}, nil
	}

	log.Info("annotating service which is part of the mesh", "annotation", fmt.Sprintf("%s=%s", metadata.IngressServiceUpstream, metadata.AnnotationTrue))
	if svc.Annotations == nil {
		svc.Annotations = map[string]string{}
	}
	svc.Annotations[metadata.IngressServiceUpstream] = metadata.AnnotationTrue

	if err = r.Update(ctx, svc); err != nil {
		return kube_ctrl.Result{}, errors.Wrapf(err, "unable to update ingress service upstream annotation on service %s", svc.Name)
	}

	return kube_ctrl.Result{}, nil
}

func (r *ServiceReconciler) SetupWithManager(mgr kube_ctrl.Manager) error {
	if err := kube_core.AddToScheme(mgr.GetScheme()); err != nil {
		return errors.Wrapf(err, "could not add %q to scheme", kube_core.SchemeGroupVersion)
	}
	if err := k8scnicncfio.AddToScheme(mgr.GetScheme()); err != nil {
		return errors.Wrapf(err, "could not add %q to scheme", k8scnicncfio.GroupVersion)
	}
	if err := apiextensionsv1.AddToScheme(mgr.GetScheme()); err != nil {
		return errors.Wrapf(err, "could not add %q to scheme", apiextensionsv1.SchemeGroupVersion)
	}
	return kube_ctrl.NewControllerManagedBy(mgr).
		For(&kube_core.Service{}, builder.WithPredicates(serviceEvents)).
		Complete(r)
}

// we only want create and update events
var serviceEvents = predicate.Funcs{
	CreateFunc: func(event event.CreateEvent) bool {
		return true
	},
	DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
		return false
	},
	UpdateFunc: func(updateEvent event.UpdateEvent) bool {
		return true
	},
	GenericFunc: func(genericEvent event.GenericEvent) bool {
		return false
	},
}
