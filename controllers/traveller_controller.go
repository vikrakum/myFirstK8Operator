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

package controllers

import (
	"context"
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"operators/api/v1alpha1"
	mydomainv1alpha1 "operators/api/v1alpha1"
)

// TravellerReconciler reconciles a Traveller object
type TravellerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=my.domain,resources=travellers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=my.domain,resources=travellers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=my.domain,resources=travellers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Traveller object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *TravellerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// _ = log.FromContext(ctx)

	log.Log.Info("Reconsile called", "Request.Namespace", req.Namespace, "Request.Name", req.Name)
	traveller := &v1alpha1.Traveller{}
	err := r.Get(ctx, req.NamespacedName, traveller)
	if err != nil {
		return ctrl.Result{}, nil
	}
	startTime := traveller.Spec.Start
	endTime := traveller.Spec.End
	replicas := traveller.Spec.Replicas
	curHour := time.Now().UTC().Hour()

	if curHour >= startTime && curHour <= endTime {
		for _, deploy := range traveller.Spec.Deployments {
			deployment := &v1.Deployment{}
			err := r.Get(ctx, types.NamespacedName{
				Name:      deploy.Name,
				Namespace: deploy.Namespace,
			}, deployment)
			if err != nil {
				return ctrl.Result{}, nil
			}
			if deployment.Spec.Replicas != &replicas {
				deployment.Spec.Replicas = &replicas
				err := r.Update(ctx, deployment)
				if err != nil {
					return ctrl.Result{}, nil
				}
			}
		}
	}

	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TravellerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mydomainv1alpha1.Traveller{}).
		Complete(r)
}
