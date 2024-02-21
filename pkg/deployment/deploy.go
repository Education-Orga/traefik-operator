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
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrl "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func ControlTraefikDeployment(ctx context.Context, k8sClient client.Client, scheme *runtime.Scheme, traefikInstance *traefikv1alpha1.TraefikInstance, namespace string) error {

	desiredDeployment := initTraefikDeployment(traefikInstance, scheme, namespace)

	// Set controller reference for CR TraefikInstance.
	if err := ctrl.SetControllerReference(traefikInstance, desiredDeployment, scheme); err != nil {
		return err
	}

	// Check if Traefik deployment exists.
	actualDeployment := &appsv1.Deployment{}
	err := k8sClient.Get(ctx, types.NamespacedName{Name: desiredDeployment.Name, Namespace: desiredDeployment.Namespace}, actualDeployment)
	if err != nil && errors.IsNotFound(err) {
		// Create new deployment.
		fmt.Println("Creating a new Traefik deployment", "Namespace", desiredDeployment.Namespace, "Name", desiredDeployment.Name)
		if err = k8sClient.Create(ctx, desiredDeployment); err != nil {
			return err
		}
	} else if err == nil {
		// Update deployment.
		if !reflect.DeepEqual(desiredDeployment.Spec, actualDeployment.Spec) {
			// Aktualisiere das aktuelle Deployment mit den Spezifikationen des gewünschten Deployments
			actualDeployment.Spec = desiredDeployment.Spec
			if err = k8sClient.Update(ctx, desiredDeployment); err != nil {
				return err
			}
		}
	}

	return nil
}

func initTraefikDeployment(traefikInstance *traefikv1alpha1.TraefikInstance, scheme *runtime.Scheme, namespace string) *appsv1.Deployment {

	// base configuration args for Traefik
	traefikArgs := initTraefikBaseConfiguration(traefikInstance)

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
					ServiceAccountName: "traefik-service-account",
					Containers: []corev1.Container{
						{
							Name:  "traefik",
							Image: traefikInstance.Spec.Image,
							Args:  traefikArgs,
						},
					},
				},
			},
		},
	}
}

func initTraefikBaseConfiguration(traefikInstance *traefikv1alpha1.TraefikInstance) []string {

	// base configuration args for Traefik
	traefikArgs := []string{
		"--api.insecure=true",
		"--accesslog",
		"--entrypoints.web.Address=:80",
		"--providers.kubernetescrd", // activate Kubernetes CRD provider for interacting with CR IngressRoute from traefik.io
		"--log.level=DEBUG",
	}

	// add additional args specified in CR
	if len(traefikInstance.Spec.AdditionalArgs) > 0 {
		traefikArgs = append(traefikArgs, traefikInstance.Spec.AdditionalArgs...)
	}

	return traefikArgs
}
