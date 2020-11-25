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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	redisv1 "gitlab.sysbind.biz/SRE/redis-operator/api/v1"
)

// RedisReconciler reconciles a Redis object
type RedisReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=redis.sysbind.co.il,resources=redis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=redis.sysbind.co.il,resources=redis/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=statefulsets/status,verbs=get

func (r *RedisReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("redis", req.NamespacedName)

	var redis redisv1.Redis
	if err := r.Get(ctx, req.NamespacedName, &redis); err != nil {
		log.Error(err, "unable to fetch Redis")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Reconciling ", "redis", req.NamespacedName.Name)
	log.Info("Structure", "masters", redis.Spec.Masters, "replicas", redis.Spec.Replicas)
	if redis.Spec.Replicas == nil || redis.Spec.Masters == nil {
		log.Info("Still no spec, re-queueing..")
		return ctrl.Result{Requeue: true}, nil
	}
	return ctrl.Result{}, nil
}

func (r *RedisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&redisv1.Redis{}).
		Owns()
		Complete(r)
}
