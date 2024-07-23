package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func watchDeployment(namespace, deploymentName string, client *kubernetes.Clientset, sleep int) {
    // loop until break condition
	for {
        // get the deployments details
		deployment, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				panic(fmt.Sprintf("Deployment %s not found in %s namespace", deploymentName, namespace))
			} else {
				panic(err.Error())
			}
		}

        // if deployment has more than 0 replicas and desired replicas are ready, break the loop
		if deployment.Status.Replicas > 0 && deployment.Status.ReadyReplicas == deployment.Status.Replicas {
			break
		} else {
            // report status of deployment
            fmt.Printf("%d of %d replicas ready\n", deployment.Status.ReadyReplicas, deployment.Status.Replicas)
        }

        // wait a bit before we check again
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

func main() {
    // arg parsing
    var namespace string
    var deploymentName string
    var sleep int

    // -namespace
    flag.StringVar(&namespace, "namespace", "", "Namespace of deployment")

    // -deployment
    flag.StringVar(&deploymentName, "deployment", "", "Deployment name")

    // -sleep
    flag.IntVar(&sleep, "sleep", 10, "Number of seconds between checking the deployment")

    flag.Parse()

    if namespace == "" || deploymentName == "" {
        flag.Usage()
        os.Exit(1)
    }

    fmt.Printf("Monitoring deployment %s in the %s namespace\n", deploymentName, namespace)


	// creates the in-cluster config from service account
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the k8s api client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

    // https://github.com/kubernetes/client-go/blob/master/examples/in-cluster-client-configuration/main.go

    watchDeployment(namespace, deploymentName, client, sleep)

    fmt.Println("Ready to go!")

}
