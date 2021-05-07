/*


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
	"errors"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	networkingv1alpha1 "github.com/rcanderson23/homelab-dynamic-dns/api/v1alpha1"
	"github.com/rcanderson23/homelab-dynamic-dns/networking/ip"
)

// IPLookupReconciler reconciles a IPLookup object
type IPLookupReconciler struct {
	client.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	RequeuePeriod time.Duration
}

// +kubebuilder:rbac:groups=networking.thehomelab.tech,resources=iplookups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.thehomelab.tech,resources=iplookups/status,verbs=get;update;patch

func (r *IPLookupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("iplookup", req.Name)

	var iplookup networkingv1alpha1.IPLookup
	if err := r.Get(ctx, req.NamespacedName, &iplookup); err != nil {
		log.Error(err, "unable to find IPLookup")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Reconciling iplookup object")
	var config ip.Lookup
	switch iplookup.Spec.Type {
	case "http":
		config = &ip.HttpLookup{Url: iplookup.Spec.Config.Http.Url}
	default:
		log.Info("Not an accepted IPLookup Type")
		return ctrl.Result{}, errors.New("not an accepted type")
	}

	ipAddr, err := config.GetIP()
	if err != nil {
		return ctrl.Result{}, err
	}
	if iplookup.Status.Address == ipAddr {
		return ctrl.Result{RequeueAfter: r.RequeuePeriod}, nil
	}
	iplookup.Status.Address = ipAddr
	err = r.Status().Update(ctx, &iplookup)
	if err != nil {
		log.Error(err, "failed to update status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: r.RequeuePeriod}, nil
}

func (r *IPLookupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1alpha1.IPLookup{}).
		Complete(r)
}
