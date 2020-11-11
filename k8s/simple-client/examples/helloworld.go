package main
import (
	"context"
	"flag"
	"fmt"
	_ "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	_ "time"
	_ "k8s.io/client-go/tools/metrics"

	core "k8s.io/api/core/v1"

)


func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("" , *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	sa, err := clientset.CoreV1().ServiceAccounts("default").List(context.TODO(),
		metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}
	for _, e := range sa.Items {
		fmt.Println(e.Name)
	}
	kapi := clientset.CoreV1()
	nodes, err := kapi.Nodes().List(context.TODO(),
		metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, e := range nodes.Items {
		// fmt.Printf("%d. %s :: [Status : %s] \n",i, e.Name, e.Status.String())
		fmt.Println(e.Status.Capacity.Cpu().String())
	}

	created, err := kapi.Pods("dev").Create(context.TODO(), getPodObject(), metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Created Pod Successfully")

	fmt.Println(created.Name)
	/*for {
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(),
			metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("There are %d pods in the cluster\n", len(pods.Items))

		namespace := "default"
		pod := "example-xxxxx"
		_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(),
			pod, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		}
		time.Sleep(3 * time.Second)
	}*/
}

func getPodObject() *core.Pod {
	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "firstly-pod",
			Namespace: "dev",
			Labels: map[string]string {
				"app": "demo",
			},
		},
		Spec : core.PodSpec {
			NodeName: "dn2",
			Containers: []core.Container {
				{
					Name: "busybox-hello",
					Image: "busybox",
					ImagePullPolicy: core.PullIfNotPresent,
					Command: []string {
						"sleep",
						"600",
					},
				},
			},
		},
	}
}