package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	examplecomv1 "github.com/trewolff/operator-example/api/v1"
	examplecomv1alpha1 "github.com/trewolff/operator-example/api/v1alpha1"
)

// MyAppReconciler reconciles a MyApp object

type MyAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.com,resources=myapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.com,resources=myapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.com,resources=myapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile

func (r *MyAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// 1) Load primary resource
	myapp := &examplecomv1alpha1.MyApp{}
	if err := r.Get(ctx, req.NamespacedName, myapp); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2) Desired Deployment skeleton
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      myapp.Name,
			Namespace: myapp.Namespace,
		},
	}

	// 3) CreateOrUpdate to be idempotent
	op, err := controllerutil.CreateOrUpdate(ctx, r.Client, dep, func() error {
		// set owner
		if err := controllerutil.SetControllerReference(myapp, dep, r.Scheme); err != nil {
			return err
		}

		// mutate spec from MyApp
		replicas := myapp.Spec.Replicas
		dep.Spec = appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": myapp.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": myapp.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "app",
						Image: myapp.Spec.Image,
					}},
				},
			},
		}
		return nil
	})
	if err != nil {
		log.Error(err, "failed to create/update Deployment")
		return ctrl.Result{}, err
	}
	log.Info("Deployment reconciled", "operation", op)

	// 4) Update status (example)
	myapp.Status.Phase = "Ready"
	if err := r.Status().Update(ctx, myapp); err != nil {
		log.Error(err, "failed to update MyApp status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplecomv1alpha1.MyApp{}).
		Named("myapp").
		Complete(r)
}

type ClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.com,resources=clusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.com,resources=clusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.com,resources=clusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Cluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = ctrl.LoggerFrom(ctx)

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplecomv1.Cluster{}).
		Named("cluster").
		Complete(r)
}
