package e2e

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	// Mysteriously required to support GCP auth (required by k8s libs).
	// Apparently just importing it is enough. @_@ side effects @_@.
	// https://github.com/kubernetes/client-go/issues/242
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/google/go-cmp/cmp"
	perrors "github.com/pkg/errors"
	"knative.dev/pkg/system"
	pkgTest "knative.dev/pkg/test"
	"knative.dev/serving/pkg/apis/networking"
	"knative.dev/serving/pkg/autoscaler"
	"knative.dev/serving/test"
	v1a1test "knative.dev/serving/test/v1alpha1"
)

// Setup creates the client objects needed in the e2e tests.
func Setup(t *testing.T) *test.Clients {
	return SetupWithNamespace(t, test.ServingNamespace)
}

// SetupAlternativeNamespace creates the client objects needed in e2e tests
// under the alternative namespace.
func SetupAlternativeNamespace(t *testing.T) *test.Clients {
	return SetupWithNamespace(t, test.AlternativeServingNamespace)
}

// SetupWithNamespace creates the client objects needed in the e2e tests under the specified namespace.
func SetupWithNamespace(t *testing.T, namespace string) *test.Clients {
	clients, err := test.NewClients(
		pkgTest.Flags.Kubeconfig,
		pkgTest.Flags.Cluster,
		namespace)
	if err != nil {
		t.Fatalf("Couldn't initialize clients: %v", err)
	}
	return clients
}

// autoscalerCM returns the current autoscaler config map deployed to the
// test cluster.
func autoscalerCM(clients *test.Clients) (*autoscaler.Config, error) {
	autoscalerCM, err := clients.KubeClient.Kube.CoreV1().ConfigMaps("knative-serving").Get(
		autoscaler.ConfigName,
		metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return autoscaler.NewConfigFromMap(autoscalerCM.Data)
}

// WaitForScaleToZero will wait for the specified deployment to scale to 0 replicas.
// Will wait up to 6 times the configured ScaleToZeroGracePeriod before failing.
func WaitForScaleToZero(t *testing.T, deploymentName string, clients *test.Clients) error {
	t.Helper()
	t.Logf("Waiting for %q to scale to zero", deploymentName)

	cfg, err := autoscalerCM(clients)
	if err != nil {
		return perrors.Wrap(err, "failed to get autoscaler configmap")
	}

	return pkgTest.WaitForDeploymentState(
		clients.KubeClient,
		deploymentName,
		test.DeploymentScaledToZeroFunc,
		"DeploymentIsScaledDown",
		test.ServingNamespace,
		cfg.ScaleToZeroGracePeriod*6,
	)
}

// waitForActivatorEndpoints waits for the Service endpoints to match that of activator.
func waitForActivatorEndpoints(resources *v1a1test.ResourceObjects, clients *test.Clients) error {
	return wait.Poll(250*time.Millisecond, time.Minute, func() (bool, error) {
		// We need to fetch the activator endpoints at every check, since it can change.
		aeps, err := clients.KubeClient.Kube.CoreV1().Endpoints(
			system.Namespace()).Get(networking.ActivatorServiceName, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}
		svcEps, err := clients.KubeClient.Kube.CoreV1().Endpoints(test.ServingNamespace).Get(
			resources.Revision.Status.ServiceName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if len(svcEps.Subsets) != len(aeps.Subsets) {
			return false, nil
		}
		for i, ss := range svcEps.Subsets {
			if !cmp.Equal(ss.Addresses, aeps.Subsets[i].Addresses) {
				return false, nil
			}
		}
		return true, nil
	})
}
