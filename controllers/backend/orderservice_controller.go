/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package backend

import (
	"context"

	backendv1 "github.com/ManojDhanorkar/restaurant-mgmt-operator/apis/backend/v1"
	resources "github.com/ManojDhanorkar/restaurant-mgmt-operator/resources"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// OrderServiceReconciler reconciles a OrderService object
type OrderServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=backend.xyzcompany.com,resources=orderservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=backend.xyzcompany.com,resources=orderservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=backend.xyzcompany.com,resources=orderservices/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps;secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the OrderService object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *OrderServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)
	log := ctrllog.FromContext(ctx)
	log.Info("Reconciling Order Service")

	orderService := &backendv1.OrderService{}

	err := r.Client.Get(ctx, req.NamespacedName, orderService)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("orderService resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get orderService.")
		return ctrl.Result{}, err
	}

	log.Info(orderService.Name, orderService.Namespace, orderService.Spec.Size)

	//log.Info("Deployment Details", "Deployment.Image", orderService.Spec.Image, "Deployment.Size", orderService.Spec.Size)

	// Check if the Deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: orderService.Name, Namespace: orderService.Namespace}, found)
	//log.Info(*found.)
	if err != nil && errors.IsNotFound(err) {
		// Define a new DeploymentJob
		//Deployment := r.DeploymentForOrderService(ctx, req, orderService)
		Deployment := resources.CreateNewDeployment(ctx, req, orderService, r.Scheme)
		log.Info("Tried Creating a new Deployment", "Deployment.Namespace", Deployment.Namespace, "Deployment.Name", Deployment.Name)
		err = r.Client.Create(ctx, Deployment)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", Deployment.Namespace, "Deployment.Name", Deployment.Name)
			return ctrl.Result{}, err
		}
		// Deploymentjob created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// // Deployment Spec
// func (r *OrderServiceReconciler) DeploymentForOrderService(ctx context.Context, req ctrl.Request, orderService *backendv1.OrderService) *appsv1.Deployment {

// 	log := ctrllog.FromContext(ctx)

// 	var replicas = orderService.Spec.Size

// 	log.Info(strconv.Itoa(int(replicas)))

// 	var labels = map[string]string{
// 		"app": req.NamespacedName.Name,
// 	}

// 	Deployment := &appsv1.Deployment{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      orderService.Name,
// 			Namespace: orderService.Namespace,
// 			Labels:    labels,
// 		},
// 		Spec: appsv1.DeploymentSpec{
// 			Replicas: &replicas,
// 			Selector: &metav1.LabelSelector{
// 				MatchLabels: labels,
// 			},
// 			Template: corev1.PodTemplateSpec{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Labels: labels,
// 				},
// 				Spec: corev1.PodSpec{
// 					Containers: []corev1.Container{{
// 						Name:  orderService.Name,
// 						Image: orderService.Spec.Image,
// 						Env: []corev1.EnvVar{
// 							{
// 								Name: "MONGODB_URL",
// 								ValueFrom: &corev1.EnvVarSource{
// 									SecretKeyRef: &corev1.SecretKeySelector{
// 										LocalObjectReference: corev1.LocalObjectReference{
// 											Name: "order-service-secrets",
// 										},
// 										Key: "mongodb_url",
// 									},
// 								},
// 							},
// 							{
// 								Name: "PORT",
// 								ValueFrom: &corev1.EnvVarSource{
// 									SecretKeyRef: &corev1.SecretKeySelector{
// 										LocalObjectReference: corev1.LocalObjectReference{
// 											Name: "order-service-secrets",
// 										},
// 										Key: "port",
// 									},
// 								},
// 							}},
// 						//RestartPolicy: "Always",
// 					}}, // Container
// 				}, // PodSec
// 			}, // PodTemplateSpec
// 		}, // Spec
// 	} // Deployment

// 	// Set AWSManager instance as the owner and controller
// 	ctrl.SetControllerReference(orderService, Deployment, r.Scheme)
// 	return Deployment
// }

// SetupWithManager sets up the controller with the Manager.
func (r *OrderServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&backendv1.OrderService{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
