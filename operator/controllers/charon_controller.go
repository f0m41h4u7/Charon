package controllers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/f0m41h4u7/Charon/operator/api/v1alpha2"
	charonv1alpha2 "github.com/f0m41h4u7/Charon/operator/api/v1alpha2"
)

// CharonReconciler reconciles a Charon object
type CharonReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=charon.charon.cr,resources=charons,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=charon.charon.cr,resources=charons/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=charon.charon.cr,resources=charons/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *CharonReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	lg := log.FromContext(ctx)
	lg.Info("Reconciling Charon")

	// Fetch the Deployer instance
	instance := &v1alpha2.Charon{}
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
	pod := createCharonPod(instance)

	// Set Deployer instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}

	// Check Service and RBAC
	errSvc := r.handleSvc(instance)
	if errSvc != nil {
		return ctrl.Result{}, errSvc
	}
	errRbac := r.handleRBAC(instance)
	if errRbac != nil {
		return ctrl.Result{}, errRbac
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
func (r *CharonReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&charonv1alpha2.Charon{}).
		Complete(r)
}

// handleSvc creates Charon service.
func (r *CharonReconciler) handleSvc(cr *v1alpha2.Charon) error {
	svc := getCharonSvcConfig(cr)

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

func getCharonSvcConfig(cr *v1alpha2.Charon) *corev1.Service {
	labels := map[string]string{
		"name": cr.Name,
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "NodePort",
			Ports: []corev1.ServicePort{
				{
					Protocol: "TCP",
					Port:     31337,
					NodePort: 31337,
				},
			},
		},
	}
}

func createCharonPod(cr *v1alpha2.Charon) *corev1.Pod {
	labels := map[string]string{
		"name": cr.Name,
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "charon-deployer-sa",
			Containers: []corev1.Container{
				{
					Name:  cr.Name,
					Image: cr.Spec.DeployerImage,
					EnvFrom: []corev1.EnvFromSource{
						{
							ConfigMapRef: &corev1.ConfigMapEnvSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "env-configmap",
								},
							},
						},
					},
				},
				{
					Name:  cr.Spec.Analyzer,
					Image: cr.Spec.AnalyzerImage,
					EnvFrom: []corev1.EnvFromSource{
						{
							ConfigMapRef: &corev1.ConfigMapEnvSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "env-configmap",
								},
							},
						},
					},
				},
			},
		},
	}

	return pod
}
