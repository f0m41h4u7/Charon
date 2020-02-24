package deployer

import (
	"context"

	deployerv1alpha1 "deployer-operator/pkg/apis/deployer/v1alpha1"

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

	"k8s.io/apimachinery/pkg/util/intstr"
)

var log = logf.Log.WithName("controller_deployer")

// Add creates a new Deployer Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDeployer{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("deployer-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Deployer
	err = c.Watch(&source.Kind{Type: &deployerv1alpha1.Deployer{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner Deployer
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &deployerv1alpha1.Deployer{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileDeployer implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileDeployer{}

// ReconcileDeployer reconciles a Deployer object
type ReconcileDeployer struct {
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Deployer object and makes changes based on the state read
// and what is in the Deployer.Spec
func (r *ReconcileDeployer) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Deployer")

	// Fetch the Deployer instance
	instance := &deployerv1alpha1.Deployer{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	pod := newPodForCR(instance)
	// Define a new Service
	svc := newSvcForCr(instance)

	// Set Deployer instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Service already exists
	foundSvc := &corev1.Service{}
        err = r.client.Get(context.TODO(), types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, foundSvc)
        if err != nil && errors.IsNotFound(err) {
                reqLogger.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
                err = r.client.Create(context.TODO(), pod)
                if err != nil {
                        return reconcile.Result{}, err
                }

                // Service created successfully
                return reconcile.Result{}, nil
        } else if err != nil {
                return reconcile.Result{}, err
        }

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

func newSvcForCr(cr *deployerv1alpha1.Deployer) *corev1.Service {
	labels := map[string]string{
                "name": "deployer",
        }
	var tport intstr.IntOrString
	tport.IntVal = 31337
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
                        Name: cr.Name,
		},
		Spec: corev1.ServiceSpec{
			Type: "ClusterIP",
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Protocol: "TCP",
					Port: 31337,
					TargetPort: tport,
				},
			},
		},
	}
}

func newPodForCR(cr *deployerv1alpha1.Deployer) *corev1.Pod {
	labels := map[string]string{
		"name": "deployer",
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
					Name:    "deployer",
					Image:   cr.Spec.Image,
				},
			},
		},
	}
}
