package deployment

import (
	"context"
	"fmt"
	traefikv1alpha1 "github.com/Education-Orga/traefik-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrl "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func ExecuteTraefikDeployment(ctx context.Context, k8sClient client.Client, scheme *runtime.Scheme, traefikInstance *traefikv1alpha1.TraefikInstance, namespace string) error {

	deployment := initTraefikDeployment(traefikInstance, scheme, namespace)

	// Set controller reference for CR TraefikInstance.
	if err := ctrl.SetControllerReference(traefikInstance, deployment, scheme); err != nil {
		return err
	}

	// Check if Traefik deployment exists.
	found := &appsv1.Deployment{}
	err := k8sClient.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create new deployment.
		fmt.Println("Creating a new Traefik deployment", "Namespace", deployment.Namespace, "Name", deployment.Name)
		err = k8sClient.Create(ctx, deployment)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func initTraefikDeployment(traefikInstance *traefikv1alpha1.TraefikInstance, scheme *runtime.Scheme, namespace string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      traefikInstance.Name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: traefikInstance.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "traefik"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "traefik"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "traefik",
							Image: traefikInstance.Spec.Image,
							Args:  []string{"--api.insecure=true", "--accesslog", "--entrypoints.web.Address=:80"},
						},
					},
				},
			},
		},
	}
}
