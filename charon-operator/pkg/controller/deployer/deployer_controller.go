package deployer

import (
	"context"

	deployerv1alpha1 "charon-operator/pkg/apis/deployer/v1alpha1"

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
	"k8s.io/client-go/kubernetes"
	rbac "k8s.io/api/rbac/v1beta1"
	"fmt"
)

var log = logf.Log.WithName("controller_deployer")

// Add creates a new Deployer Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	config := mgr.GetConfig()
	clientset,err := kubernetes.NewForConfig(config)
	if err != nil {
		panic("Failed to get node")
	}
	return &ReconcileDeployer{client: mgr.GetClient(), scheme: mgr.GetScheme(), clientset: clientset}
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
	client		client.Client
	scheme		*runtime.Scheme
	clientset	kubernetes.Interface
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
	pod := createPod(instance)

	// Set Deployer instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}

	// Check Service and RBAC existance
	errSvc := handleSvc(instance, r)
	if errSvc != nil {
		return reconcile.Result{}, errSvc
	}
	errRbac := handleRBAC(instance, r)
	if errRbac != nil {
		return reconcile.Result{}, errRbac
	}

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

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

func handleRBAC (cr *deployerv1alpha1.Deployer, r *ReconcileDeployer) error{
	fmt.Println("Handle RBAC")

	role := createRole(cr, r)
        rb := createRB(cr, r)
        sa := createSA(cr, r)

	foundRole := &rbac.Role{}
        err := r.client.Get(context.TODO(), types.NamespacedName{Name: role.Name, Namespace: role.Namespace}, foundRole)
        if err != nil && errors.IsNotFound(err) {
                err = r.client.Create(context.TODO(), role)
		if err != nil {
			return err
		}
        }

	foundRB := &rbac.RoleBinding{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: rb.Name, Namespace: rb.Namespace}, foundRB)
        if err != nil && errors.IsNotFound(err) {
                err = r.client.Create(context.TODO(), rb)
		if err != nil {
			return err
                }
        }

	foundSA := &corev1.ServiceAccount{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: sa.Name, Namespace: sa.Namespace}, foundSA)
        if err != nil && errors.IsNotFound(err) {
                err = r.client.Create(context.TODO(), sa)
		if err != nil {
                        return err
                }
        }
	return nil
}

func handleSvc (cr *deployerv1alpha1.Deployer, r *ReconcileDeployer) error {
	fmt.Println("Handle Svc")

	svc := createService(cr)

	found := &corev1.Service{}
        err := r.client.Get(context.TODO(), types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, found)
        if err != nil && errors.IsNotFound(err) {
                err = r.client.Create(context.TODO(), svc)
		if err != nil {
                        return err
		}
        }
	return nil
}

func createService (cr *deployerv1alpha1.Deployer) *corev1.Service {
	fmt.Println("Creating svc...")

	labels := map[string]string{
                "name": cr.Name,
        }
	var tport intstr.IntOrString
	tport.IntVal = 31337
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
                        Name:		cr.Name,
			Namespace:	cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector:	labels,
			Type:		"ClusterIP",
			Ports:		[]corev1.ServicePort{
				{
					Protocol:	"TCP",
					Port:		31337,
					TargetPort:	tport,
				},
			},
		},
	}
}

func createRole (cr *deployerv1alpha1.Deployer, r *ReconcileDeployer) *rbac.Role {
        fmt.Println("Creating role...")

        role := &rbac.Role{
                TypeMeta: metav1.TypeMeta{
                        APIVersion:     "rbac.authorization.k8s.io/v1",
                        Kind:           "Role",
                },
                ObjectMeta: metav1.ObjectMeta{
                        Name:           "charon-deployer-role",
                        Namespace:      "default",
                },
		Rules: []rbac.PolicyRule{
			{
				APIGroups:	[]string{"app.custom.cr", "apps"},
				Resources:	[]string{"apps", "pods"},
				Verbs:		[]string{"get", "create", "update", "delete", "patch", "list"},
			},
			{
				APIGroups:      []string{"apps"},
                                Resources:      []string{"*"},
                                Verbs:          []string{"get", "create", "update", "delete", "patch", "list"},
			},
			{
				APIGroups:      []string{""},
                                Resources:      []string{"*"},
                                Verbs:          []string{"get", "create", "update", "delete", "patch", "list"},
			},
		},
        }

        controllerutil.SetControllerReference(cr, role, r.scheme)
        return role
}

func createRB (cr *deployerv1alpha1.Deployer, r *ReconcileDeployer) *rbac.RoleBinding {
        fmt.Println("Creating rb...")

        rb := &rbac.RoleBinding{
                TypeMeta: metav1.TypeMeta{
                        APIVersion:     "rbac.authorization.k8s.io/v1",
                        Kind:           "ClusterRoleBinding",
                },
                ObjectMeta: metav1.ObjectMeta{
                        Name:           "charon-deployer-rb",
                        Namespace:      "default",
                },
		RoleRef: rbac.RoleRef{
			APIGroup:	"rbac.authorization.k8s.io",
			Kind:		"Role",
			Name:		"charon-deployer-role",
		},
		Subjects: []rbac.Subject{
			{
				Kind:		"ServiceAccount",
				Name:		"charon-deployer-sa",
				Namespace:	"default",
			},
		},
        }

        controllerutil.SetControllerReference(cr, rb, r.scheme)
        return rb
}

func createSA (cr *deployerv1alpha1.Deployer, r *ReconcileDeployer) *corev1.ServiceAccount {
	fmt.Println("Creating sa...")

	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion:	"v1",
			Kind:		"ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:		"charon-deployer-sa",
			Namespace:	"default",
		},
	}

	controllerutil.SetControllerReference(cr, sa, r.scheme)
	return sa
}

func createPod (cr *deployerv1alpha1.Deployer) *corev1.Pod {
	fmt.Println("Creating pod...")

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
			ServiceAccountName:		"charon-deployer-sa",
			Containers:			[]corev1.Container{
				{
					Name:	cr.Name,
					Image:	cr.Spec.Image,
				},
			},
		},
	}
}
