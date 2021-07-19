package controllers

import (
	"context"

	"github.com/f0m41h4u7/Charon/operator/api/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *CharonReconciler) handleRBAC(cr *v1alpha2.Charon) error {
	role := r.createRole(cr)
	rb := r.createRB(cr)
	sa := r.createSA(cr)

	foundRole := &rbac.Role{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: role.Name, Namespace: role.Namespace}, foundRole)
	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(context.TODO(), role)
		if err != nil {
			return err
		}
	}

	foundRB := &rbac.RoleBinding{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: rb.Name, Namespace: rb.Namespace}, foundRB)
	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(context.TODO(), rb)
		if err != nil {
			return err
		}
	}

	foundSA := &corev1.ServiceAccount{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: sa.Name, Namespace: sa.Namespace}, foundSA)
	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(context.TODO(), sa)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *CharonReconciler) createRole(cr *v1alpha2.Charon) *rbac.Role {
	role := &rbac.Role{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "Role",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "charon-deployer-role",
			Namespace: "charon",
		},
		Rules: []rbac.PolicyRule{
			{
				APIGroups: []string{"app.custom.cr", "apps"},
				Resources: []string{"apps", "pods"},
				Verbs:     []string{"get", "create", "update", "delete", "patch", "list"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"*"},
				Verbs:     []string{"get", "create", "update", "delete", "patch", "list"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"*"},
				Verbs:     []string{"get", "create", "update", "delete", "patch", "list"},
			},
		},
	}

	_ = controllerutil.SetControllerReference(cr, role, r.Scheme)
	return role
}

func (r *CharonReconciler) createRB(cr *v1alpha2.Charon) *rbac.RoleBinding {
	rb := &rbac.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "ClusterRoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "charon-deployer-rb",
			Namespace: "charon",
		},
		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     "charon-deployer-role",
		},
		Subjects: []rbac.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "charon-deployer-sa",
				Namespace: "charon",
			},
		},
	}

	_ = controllerutil.SetControllerReference(cr, rb, r.Scheme)
	return rb
}

func (r *CharonReconciler) createSA(cr *v1alpha2.Charon) *corev1.ServiceAccount {
	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "charon-deployer-sa",
			Namespace: "charon",
		},
	}

	_ = controllerutil.SetControllerReference(cr, sa, r.Scheme)
	return sa
}
