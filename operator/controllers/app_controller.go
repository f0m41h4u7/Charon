package controllers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/f0m41h4u7/Charon/operator/api/v1alpha2"
	charonv1alpha2 "github.com/f0m41h4u7/Charon/operator/api/v1alpha2"
)

// AppReconciler reconciles a App object
type AppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=charon.charon.cr,resources=apps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=charon.charon.cr,resources=apps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=charon.charon.cr,resources=apps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *AppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	lg := log.FromContext(ctx)
	lg.Info("Reconciling App")

	// Fetch the App instance
	instance := &v1alpha2.App{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Define a new Pod object
	pod := createPod(instance)

	// Set App instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}

	// Check Service
	errSvc := r.handleSvc(instance)
	if errSvc != nil {
		return ctrl.Result{}, errSvc
	}

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		lg.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.Client.Create(context.TODO(), pod)
		if err != nil {
			return ctrl.Result{}, err
		}

		// Pod created successfully - don't requeue
		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Pod already exists - don't requeue
	lg.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&charonv1alpha2.App{}).
		Complete(r)
}

func (r *AppReconciler) handleSvc(cr *v1alpha2.App) error {
	svc := getAppSvcConfig(cr)

	found := &corev1.Service{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(context.TODO(), svc)
		if err != nil {
			return err
		}
	}
	return nil
}

func getAppSvcConfig(cr *v1alpha2.App) *corev1.Service {
	labels := map[string]string{
		"name": cr.Name,
	}
	var tport intstr.IntOrString
	tport.IntVal = 1337
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "ClusterIP",
			Ports: []corev1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       1337,
					TargetPort: tport,
				},
			},
		},
	}
}

func createPod(cr *v1alpha2.App) *corev1.Pod {
	labels := map[string]string{
		"name": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  cr.Name,
					Image: cr.Spec.Image,
				},
			},
		},
	}
}
