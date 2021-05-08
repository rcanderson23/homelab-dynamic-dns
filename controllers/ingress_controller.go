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
	"github.com/go-logr/logr"
	"strings"
	"time"

	namev1alpha1 "github.com/rcanderson23/homelab-dynamic-dns/api/v1alpha1"
	dnsp "github.com/rcanderson23/homelab-dynamic-dns/networking/dnsproviders"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const externalIngress = "thehomelab.tech/ip-address"
const nsName = "thehomelab.tech/resolver"
const ipLookupName = "thehomelab.tech/iplookup"

// IngressReconciler reconciles a Ingress object
type IngressReconciler struct {
	client.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	RequeuePeriod time.Duration
}

// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get
// +kubebuilder:rbac:groups=networking.thehomelab.tech,resources=iplookups,verbs=get
// +kubebuilder:rbac:groups=networking.thehomelab.tech,resources=dnsproviders,verbs=get;list;watch

func (r *IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("ingress", req.NamespacedName)

	var ing netv1.Ingress
	if err := r.Get(ctx, req.NamespacedName, &ing); err != nil {
		log.Error(err, "unable to find Ingress")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if this ingress has the necessary annotations to proceed
	annotated := hasAnnotations(ing)
	if !annotated {
		return ctrl.Result{}, nil
	}
	log.Info("Reconciling Ingress object")

	// Get the IPLookup defined if it set to external
	var iplookup namev1alpha1.IPLookup
	ipName, present := ing.Annotations[ipLookupName]
	if present && strings.ToLower(ipName) == "external" {
		if err := r.Get(ctx, client.ObjectKey{Name: ipName}, &iplookup); err != nil {
			log.Error(err, "IPLookup not found", ipName)
			return ctrl.Result{}, err
		}
	}
	var ipAddrs []string
	if strings.ToLower(ipName) == "external" {
		ipAddrs = append(ipAddrs, iplookup.Status.Address)
	} else {
		ipAddrs = getIngressIPs(ing.Status)
	}

	// Get all the hosts in the ingress that don't match the desired DNS entry
	var hosts []string
	for _, rule := range ing.Spec.Rules {
		if rule.Host != "" {
			for _, ip := range ipAddrs {
				ep, err := dnsp.IsCurrentEndpoint(rule.Host, ip)
				log.Info("Endpoint lookup", "host", rule.Host, "endpoint", ep)
				if err != nil {
					return ctrl.Result{}, err
				}
				if !ep {
					log.Info("Host does not match endpoint", "host", rule.Host)
					hosts = append(hosts, rule.Host)
				}
			}
		}
	}

	// if the length of host updates is 0, return early
	if len(hosts) == 0 {
		log.Info("Ingress endpoint(s) do not require update")
		return ctrl.Result{RequeueAfter: r.RequeuePeriod}, nil
	}

	// Get resolver object to update nameserver if necessary
	var p namev1alpha1.DNSProvider
	if err := r.Get(ctx, client.ObjectKey{Name: ing.Annotations[nsName]}, &p); err != nil {
		log.Info("Resolver not found", "Nameserver", ing.Annotations[nsName])
		return ctrl.Result{}, err
	}

	// Create the nameserver and ensure A record is present
	ns, err := dnsp.NewNameserver(p.Spec.Type, p.Spec.Config)
	if err != nil {
		log.Error(err, "Failed to create nameserver", ing.Annotations[nsName])
		return ctrl.Result{}, err
	}
	for _, host := range hosts {
		for _, ip := range ipAddrs {
			err := ns.EnsureRecordA(ctx, host, ip)
			if err != nil {
				return ctrl.Result{}, err
			}
			log.Info("DNS entry created", "record", host, "value", ip)
		}
	}
	return ctrl.Result{RequeueAfter: r.RequeuePeriod}, nil
}

func (r *IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&netv1.Ingress{}).
		Complete(r)
}

func getIngressIPs(status netv1.IngressStatus) []string {
	var ips []string
	for _, lbStatus := range status.LoadBalancer.Ingress {
		ips = append(ips, lbStatus.IP)
	}
	return ips
}

func hasAnnotations(ing netv1.Ingress) bool {
	var present bool
	_, present = ing.Annotations[externalIngress]
	if !present {
		return false
	}
	_, present = ing.Annotations[nsName]
	if !present {
		return false
	}
	return true
}
