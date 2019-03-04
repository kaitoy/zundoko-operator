/*
Copyright 2019 kaitoy.

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

package hikawa

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	zundokokiyoshiv1beta1 "github.com/kaitoy/zundoko-operator/pkg/apis/zundokokiyoshi/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")

// Add creates a new Hikawa Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	rand.Seed(time.Now().UnixNano())
	return &ReconcileHikawa{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("hikawa-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Hikawa
	err = c.Watch(&source.Kind{Type: &zundokokiyoshiv1beta1.Hikawa{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &zundokokiyoshiv1beta1.Zundoko{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &zundokokiyoshiv1beta1.Hikawa{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileHikawa{}

// ReconcileHikawa reconciles a Hikawa object
type ReconcileHikawa struct {
	client.Client
	scheme *runtime.Scheme
}

const wordZun = "Zun"
const wordDoko = "Doko"
const wordKiyoshi = "Kiyoshi!"

// Reconcile reads that state of the cluster for a Hikawa object and makes changes based on the state read
// and what is in the Hikawa.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Zundokos
// +kubebuilder:rbac:groups=zundokokiyoshi.kaitoy.github.com,resources=hikawas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=zundokokiyoshi.kaitoy.github.com,resources=hikawas/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=zundokokiyoshi.kaitoy.github.com,resources=zundokos,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=zundokokiyoshi.kaitoy.github.com,resources=zundokos/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=zundokokiyoshi.kaitoy.github.com,resources=kiyoshis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=zundokokiyoshi.kaitoy.github.com,resources=kiyoshis/status,verbs=get;update;patch
func (r *ReconcileHikawa) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	instanceName := request.NamespacedName.String()
	log.Info("Reconciling a Hikawa: " + instanceName)

	// Fetch the Hikawa instance
	instance := &zundokokiyoshiv1beta1.Hikawa{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if instance.Status.Kiyoshied {
		log.Info(instanceName + " has kiyoshied.")
		return reconcile.Result{}, nil
	}

	zundokoList := &zundokokiyoshiv1beta1.ZundokoList{}
	if err := r.List(context.TODO(), &client.ListOptions{Namespace: instance.Namespace}, zundokoList); err != nil {
		log.Error(err, "Failed to list zundokos for: ", instanceName)
		return reconcile.Result{}, err
	}

	var dependents []zundokokiyoshiv1beta1.Zundoko
	for _, zundoko := range zundokoList.Items {
		for _, owner := range zundoko.GetOwnerReferences() {
			if owner.Name == instance.Name {
				dependents = append(dependents, zundoko)
			}
		}
	}
	numZundokosSaid := len(dependents)

	if instance.Spec.NumZundokos > numZundokosSaid {
		log.Info(instanceName + " wants " + strconv.Itoa(instance.Spec.NumZundokos-numZundokosSaid) + " more zundoko(s).")
		time.Sleep(instance.Spec.IntervalMillis * time.Millisecond)
		word := getZundoko()
		if err := createZundoko(instance, r, fmt.Sprintf("-zundoko-%03d", numZundokosSaid+1), word); err != nil {
			return reconcile.Result{}, err
		}
	} else if instance.Status.NumZundokosSaid != numZundokosSaid {
		log.Info(instanceName + " has said " + strconv.Itoa(numZundokosSaid) + " zundoko(s). Updating the status.")
		instance.Status.NumZundokosSaid = numZundokosSaid
		if err := r.Update(context.Background(), instance); err != nil {
			log.Error(err, "Failed to update "+instanceName)
			return reconcile.Result{}, err
		}
	} else if instance.Spec.SayKiyoshi {
		log.Info(instanceName + " is going to say " + wordKiyoshi)
		time.Sleep(instance.Spec.IntervalMillis * time.Millisecond)
		if err := createKiyoshi(instance, r); err != nil {
			return reconcile.Result{}, err
		}

		instance.Status.Kiyoshied = true
		if err := r.Update(context.Background(), instance); err != nil {
			log.Error(err, "Failed to update "+instanceName)
			return reconcile.Result{}, err
		}
	} else if isReadyToKiyoshi(dependents) {
		log.Info(instanceName + " is ready to say " + wordKiyoshi)
		instance.Spec.SayKiyoshi = true
		if err := r.Update(context.Background(), instance); err != nil {
			log.Error(err, "Failed to update "+instanceName)
			return reconcile.Result{}, err
		}
	} else {
		log.Info(instanceName + " keeps going on ZUNDOKO.")
		instance.Spec.NumZundokos++
		if err := r.Update(context.Background(), instance); err != nil {
			log.Error(err, "Failed to update "+instanceName)
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func createZundoko(instance *zundokokiyoshiv1beta1.Hikawa, r *ReconcileHikawa, nameSuffix, word string) error {
	zundoko := &zundokokiyoshiv1beta1.Zundoko{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + nameSuffix,
			Namespace: instance.Namespace,
		},
		Spec: zundokokiyoshiv1beta1.ZundokoSpec{
			Say: word,
		},
	}
	if err := controllerutil.SetControllerReference(instance, zundoko, r.scheme); err != nil {
		log.Error(err, "An error occurred in SetControllerReference", "instance", instance.Name, "namespace", zundoko.Namespace, "name", zundoko.Name)
		return err
	}

	log.Info("Creating Zundoko", "namespace", zundoko.Namespace, "name", zundoko.Name)
	if err := r.Create(context.TODO(), zundoko); err != nil {
		log.Error(err, "Failed to create Zundoko", "namespace", zundoko.Namespace, "name", zundoko.Name)
		return err
	}

	return nil
}

func createKiyoshi(instance *zundokokiyoshiv1beta1.Hikawa, r *ReconcileHikawa) error {
	kiyoshi := &zundokokiyoshiv1beta1.Kiyoshi{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-kiyoshi",
			Namespace: instance.Namespace,
		},
		Spec: zundokokiyoshiv1beta1.KiyoshiSpec{
			Say: wordKiyoshi,
		},
	}
	if err := controllerutil.SetControllerReference(instance, kiyoshi, r.scheme); err != nil {
		log.Error(err, "An error occurred in SetControllerReference", "instance", instance.Name, "namespace", kiyoshi.Namespace, "name", kiyoshi.Name)
		return err
	}

	log.Info("Creating Zundoko", "namespace", kiyoshi.Namespace, "name", kiyoshi.Name)
	if err := r.Create(context.TODO(), kiyoshi); err != nil {
		log.Error(err, "Failed to create Kiyoshi", "namespace", kiyoshi.Namespace, "name", kiyoshi.Name)
		return err
	}

	return nil
}

func getZundoko() string {
	word := wordZun
	if rand.Intn(10) > 4 {
		word = wordDoko
	}
	return word
}

func isReadyToKiyoshi(zundokos []zundokokiyoshiv1beta1.Zundoko) bool {
	numZundokos := len(zundokos)
	if numZundokos < 5 {
		return false
	}

	sort.Slice(zundokos, func(i, j int) bool {
		return zundokos[j].GetCreationTimestamp().Time.After(zundokos[i].GetCreationTimestamp().Time)
	})

	if numZundokos > 5 {
		if zundokos[numZundokos-6].Spec.Say == wordZun {
			return false
		}
	}

	return zundokos[numZundokos-5].Spec.Say == wordZun &&
		zundokos[numZundokos-4].Spec.Say == wordZun &&
		zundokos[numZundokos-3].Spec.Say == wordZun &&
		zundokos[numZundokos-2].Spec.Say == wordZun &&
		zundokos[numZundokos-1].Spec.Say == wordDoko
}
