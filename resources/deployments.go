package resources

import (
	"context"
	"strconv"

	backendv1 "github.com/ManojDhanorkar/restaurant-mgmt-operator/apis/backend/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// Deployment Spec
func CreateNewDeployment(ctx context.Context, req ctrl.Request, orderService *backendv1.OrderService, scheme *runtime.Scheme) *appsv1.Deployment {

	log := ctrllog.FromContext(ctx)

	var replicas = orderService.Spec.Size

	log.Info(strconv.Itoa(int(replicas)))

	var labels = map[string]string{
		"app": req.NamespacedName.Name,
	}

	Deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      orderService.Name,
			Namespace: orderService.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  orderService.Name,
						Image: orderService.Spec.Image,
						Env: []corev1.EnvVar{
							{
								Name: "MONGODB_URL",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "order-service-secrets",
										},
										Key: "mongodb_url",
									},
								},
							},
							{
								Name: "PORT",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "order-service-secrets",
										},
										Key: "port",
									},
								},
							}},
						//RestartPolicy: "Always",
					}}, // Container
				}, // PodSec
			}, // PodTemplateSpec
		}, // Spec
	} // Deployment

	// Set AWSManager instance as the owner and controller
	ctrl.SetControllerReference(orderService, Deployment, scheme)
	return Deployment

}
