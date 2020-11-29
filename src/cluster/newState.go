package cluster

import (
	"fmt"
	redisv1 "gitlab.sysbind.biz/SRE/redis-operator/api/v1"
	kapps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

type newState struct {
	cluster *Cluster
}

func constructStatefulSetForRedis(redis *redisv1.Redis, scheme *runtime.Scheme) (*kapps.StatefulSet, error) {
	spec := kapps.StatefulSetSpec{}
	//spec.Selector = metav1.LabelSelector{}
	spec.Selector = &metav1.LabelSelector{MatchLabels: map[string]string{"redis": redis.Name}}
	spec.Template.Labels = map[string]string{"redis": redis.Name}
	spec.Replicas = redis.Spec.Masters
	spec.ServiceName = redis.Name // Must be same as the Headless service name

	// Configuring each POD with same ports for the containers will cause
	// Anti-Affinity to happen automatically
	var i int32
	for i = 0; i <= (*redis.Spec.Replicas); i++ {
		port := 6379 + i
		cluster_port := 16379 + i
		RedisContainer := corev1.Container{Name: fmt.Sprintf("redis-%d", i), Image: "redis:6.0.9",
			Command: []string{"/usr/local/bin/redis-server",
				"--port", fmt.Sprintf("%d", port),
				"--cluster-enabled", "yes"},
			Ports: []corev1.ContainerPort{{
				Name:          fmt.Sprintf("redis-%d", i),
				HostPort:      port,
				ContainerPort: port,
				Protocol:      "TCP",
			},
				{
					Name:          fmt.Sprintf("redis-cluster-%d", i),
					HostPort:      cluster_port,
					ContainerPort: cluster_port,
					Protocol:      "TCP",
				},
			},
		}
		spec.Template.Spec.Containers = append(spec.Template.Spec.Containers, RedisContainer)
	}
	sts := &kapps.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Name:        redis.Name,
			Namespace:   redis.Namespace,
		},
		Spec: spec,
	}

	if err := ctrl.SetControllerReference(redis, sts, scheme); err != nil {
		return nil, err
	}

	return sts, nil
} // constructStatefulSetForRedis

func constructHeadlessServiceForRedis(redis *redisv1.Redis, scheme *runtime.Scheme) (*corev1.Service, error) {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redis.Name,
			Namespace: redis.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector:  map[string]string{"redis": redis.Name},
			ClusterIP: corev1.ClusterIPNone,
		},
	}
	if err := ctrl.SetControllerReference(redis, svc, scheme); err != nil {
		return nil, err
	}
	return svc, nil
}

func (state newState) boot() error {
	log := state.cluster.logger
	// New Redis - Create Statefulset
	if newsts, err := constructStatefulSetForRedis(&state.cluster.spec, state.cluster.scheme); err != nil {
		log.Error(err, "unable to construct statefulset for redis")
		return err
	} else {
		if err := state.cluster.client.Create(state.cluster.ctx, newsts); err != nil {
			log.Error(err, "unable to create Statefulset for Redis", "statefuleset", newsts)
			return err
		}
		log.Info("Created Statefulset for Redis")

		// New Redis - Create Headless Service
		if newsvc, err := constructHeadlessServiceForRedis(&state.cluster.spec, state.cluster.scheme); err != nil {
			log.Error(err, "unable to construct headless service for redis")
			return err
		} else {
			if err := state.cluster.client.Create(state.cluster.ctx, newsvc); err != nil {
				log.Error(err, "unable to create Headless Service for Redis", "service", newsvc)
				return err
			}
		}
		log.Info("Created Headless Service for Redis")
	}
	return nil
}

func (state newState) create() error {
	return nil
}

func (state newState) scale(count int) error {
	return nil
}

func (state newState) destroy() error {
	return nil
}
