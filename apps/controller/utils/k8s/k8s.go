package k8s

import (
	"context"
	"controller/utils/auth"
	"fmt"
	"os"
	"time"

	"github.com/carlmjohnson/requests"
	"helm.sh/helm/v3/pkg/release"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset
var eigenDbChartPath string = "./infra/helmcharts/eigendb"

// Initializes the clientset in the k8s package
func Init(devMode bool) error {
	var config *rest.Config
	var err error
	if devMode {
		config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return err
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	return nil
}

func SpawnInstance(customerId string) (*release.Release, string, error) {
	namespace := generateNamespaceString(customerId)
	releaseName := generateReleaseName(customerId)
	//instanceUrl := getInstanceUrl(customerId)
	apiKey, err := auth.GenerateApiKey()
	if err != nil {
		return nil, "", err
	}

	actionCfg, err := newHelmActionCfg(namespace)
	if err != nil {
		return nil, "", err
	}

	chartValues := map[string]any{
		"customerId": customerId,
		"apiKey":     apiKey,
	}
	release, err := helmInstall(
		actionCfg,
		namespace,
		releaseName,
		eigenDbChartPath,
		chartValues,
	)
	if err != nil {
		return nil, "", err
	}

	//if err := pollUntilInstanceReady(customerId, time.Second*60); err != nil { // wait until the instance is ready
	//	return nil, "", err
	//}
	//if err := testApiKey(apiKey, instanceUrl); err != nil { // testing the validity of the generated API key
	//	return nil, "", err
	//}

	return release, apiKey, nil
}

func TerminateInstance(customerId string) (*release.UninstallReleaseResponse, error) {
	namespace := generateNamespaceString(customerId)
	releaseName := generateReleaseName(customerId)
	actionCfg, err := newHelmActionCfg(namespace)
	if err != nil {
		return nil, err
	}
	res, err := helmUninstall(
		actionCfg,
		releaseName,
	)
	if err != nil {
		return nil, err
	}

	// Delete the instance's namespace. This is done to clear any lingering PVCs and/or PVs as those get left behind by Helm.
	// https://github.com/helm/helm/issues/5156
	ctx := context.Background()
	deletionPolicy := metav1.DeletePropagationBackground
	corev1 := clientset.CoreV1()
	if err := corev1.Namespaces().Delete(
		ctx,
		namespace,
		metav1.DeleteOptions{
			PropagationPolicy: &deletionPolicy,
		},
	); err != nil {
		return nil, err
	}

	return res, nil
}

// used for testing the API key. this function (in theory) should create a delay until a customer instance is fully created, before sending a /test-auth request
// testing the API key has been planned to do later so this function is currently un-used.
func pollUntilInstanceReady(customerId string, timeout time.Duration) error {
	namespace := generateNamespaceString(customerId)
	statefulSetName := "eigen-" + customerId
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(context.Context) (bool, error) {
		statefulSet, err := clientset.AppsV1().StatefulSets(namespace).Get(ctx, statefulSetName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		// Check if all pods in StatefulSet are at "Ready" status
		if statefulSet.Status.ReadyReplicas == *statefulSet.Spec.Replicas {
			// fmt.Printf("All %d replicas are ready for StatefulSet %s\n", statefulSet.Status.ReadyReplicas, name)
			return true, nil
		}
		fmt.Printf("Waiting for StatefulSet %s: %d/%d replicas are ready\n", statefulSetName, statefulSet.Status.ReadyReplicas, *statefulSet.Spec.Replicas)
		return false, nil
	})
}

func testApiKey(key string, instanceUrl string) error {
	type resSchema struct {
		Status  int            `json:"status"`
		Message string         `json:"message"`
		Data    map[string]any `json:"data"`
		Error   struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	}

	var res resSchema
	err := requests.
		URL(instanceUrl+"/test-auth").
		Header("X-Eigen-API-Key", key).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		return err
	}

	if res.Status != 200 {
		return fmt.Errorf("ERROR TESTING API KEY: %s - %s", res.Error.Code, res.Error.Description)
	}

	return nil
}

func generateNamespaceString(customerId string) string {
	return fmt.Sprintf("eigendb-%s", customerId)
}

func generateReleaseName(customerId string) string {
	return fmt.Sprintf("eigendb-%s-release", customerId)
}

func getInstanceUrl(customerId string) string {
	return fmt.Sprintf("http://%s.127.0.0.1.nip.io", customerId)
}
