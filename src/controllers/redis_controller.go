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
	"gitlab.sysbind.biz/operators/redis-operator/cluster"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	redisv1 "gitlab.sysbind.biz/operators/redis-operator/api/v1"
	kapps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	stsOwnerKey = ".metadata.controller"
	apiGVStr    = redisv1.GroupVersion.String()
)

// RedisReconciler reconciles a Redis object
type RedisReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=redis.sysbind.co.il,resources=redis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=redis.sysbind.co.il,resources=redis/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
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
	log.Info("Searching for existing Statefulset objects for", "Redis", req.Name)

	var childSTS kapps.StatefulSetList
	if err := r.List(ctx, &childSTS, client.InNamespace(req.Namespace), client.MatchingFields{stsOwnerKey: req.Name}); err != nil {
		log.Error(err, "unable to list child StatefulSets")
		return ctrl.Result{}, err
	}

	// Check if we already have StatefulSet for this Redis
	for i, sts := range childSTS.Items {
		log.Info("found statefulset", "idx", i, "name", sts.Name)
		return ctrl.Result{}, nil
	}

	newCluster := cluster.NewCluster(redis, r, r.Scheme, ctx, log)
	if err := newCluster.Boot(); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *RedisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Setup stsOwnerKey Index for better searching statefulsets owned by this controller.
	if err := mgr.GetFieldIndexer().IndexField(&kapps.StatefulSet{}, stsOwnerKey, func(rawObj runtime.Object) []string {
		// grab the job object, extract the owner...
		sts := rawObj.(*kapps.StatefulSet)
		owner := metav1.GetControllerOf(sts)
		if owner == nil {
			return nil
		}
		// ...make sure it's a CronJob...
		if owner.APIVersion != apiGVStr || owner.Kind != "Redis" {
			return nil
		}

		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&redisv1.Redis{}).
		Owns(&kapps.StatefulSet{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
