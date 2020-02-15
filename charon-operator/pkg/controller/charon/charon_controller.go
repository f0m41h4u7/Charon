package charon

import (
	"context"

	charonv1alpha1 "charon-operator/pkg/apis/charon/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_charon")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCharon{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("charon-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &charonv1alpha1.Charon{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &charonv1alpha1.Charon{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileCharon{}

type ReconcileCharon struct {
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileCharon) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Charon")

	instance := &charonv1alpha1.Charon{}
        err := r.client.Get(context.TODO(), request.NamespacedName, instance)
        if err != nil {
                if errors.IsNotFound(err) {
                        // Request object not found, could have been deleted after reconcile request.
                        // Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
                        // Return and don't requeue
                        return reconcile.Result{}, nil
                }
                // Error reading the object - requeue the request.
                return reconcile.Result{}, err
        }

        // If pod doesn't exist, create one
        pod := &corev1.Pod{}
        err = r.client.Get(context.Background(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, pod)
        if err != nil && errors.IsNotFound(err) {
                reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
                pod := newPodForCR(instance)

                // Set Test instance as the owner and controller                                                 
                if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
                       return reconcile.Result{}, err
                }

//        found := &corev1.Pod{}                                                                           
                err = r.client.Create(context.TODO(), pod)
                if err != nil {
                        return reconcile.Result{}, err
                }

                // Pod created successfully - don't requeue                                               
                return reconcile.Result{}, nil
        }
        // If exists, check for updates                                                                                                        
//      if instanse.Status.VersionChanged {
        //      instance.Status.VersionChanged = false

        // ......

//              err = r.client.Status().Update(context.Background(), instance)
//              if err != nil {
//                      reqLogger.Error(err, "Failed to update pod status", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
//                      return reconcile.Result{}, err
//              }
//              return reconcile.Result{Requeue: true}, nil
//      }
        reqLogger.Info("Skip reconcile: Pod already exists and up-to-date", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name) 
        return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *charonv1alpha1.Charon) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    cr.Name, //Command: []string{"sleep", "3600"},
					Image:   cr.Spec.Image,
					Command: []string{"go", "run", "$GOPATH/server.go"},
				},
			},
		},
	}
}
