package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"os"
)

var kubecfg string
var colors map[string]LedColor

func main() {
	colors = map[string]LedColor{
		"red":    LedColor{Red: 255, Blue: 0, Green: 0, Led: 0},
		"brown":  LedColor{Red: 139, Blue: 19, Green: 69, Led: 1},
		"orange": LedColor{Red: 255, Blue: 0, Green: 69, Led: 2},
		"yellow": LedColor{Red: 255, Blue: 0, Green: 255, Led: 3},
		"green":  LedColor{Red: 0, Blue: 0, Green: 255, Led: 4},
		"blue":   LedColor{Red: 0, Blue: 255, Green: 0, Led: 5},
		"violet": LedColor{Red: 128, Blue: 128, Green: 0, Led: 6},
		"grey":   LedColor{Red: 255, Blue: 100, Green: 255, Led: 7},
	}
	var config *rest.Config

	flag.StringVar(&kubecfg, "kubeconfig", "", "Path to kubeconfig")
	flag.Parse()
	fmt.Println(kubecfg)

	if kubecfg == "" {
		config, _ = rest.InClusterConfig()
	} else {
		config, _ = clientcmd.BuildConfigFromFlags("", kubecfg)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	bc := BlinktController{
		client:  client,
		blinkts: make(map[string]string),
	}

	_, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return client.CoreV1().Pods(meta_v1.NamespaceAll).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return client.CoreV1().Pods(meta_v1.NamespaceAll).Watch(options)
			},
		},
		&core_v1.Pod{},
		0, //Skip resync
		cache.ResourceEventHandlerFuncs{
			AddFunc:    func(new interface{}) { bc.PodAdded(new) },
			UpdateFunc: func(old, new interface{}) { bc.PodDeleted(old); bc.PodAdded(new) },
			DeleteFunc: func(new interface{}) { bc.PodDeleted(new) },
		},
	)
	bc.controller = controller

	controller.Run(wait.NeverStop)
}

type BlinktController struct {
	client     kubernetes.Interface
	controller cache.Controller
	blinkts    map[string]string
}

func (c *BlinktController) PodDeleted(obj interface{}) {
	pod := obj.(*core_v1.Pod)

	val, ok := pod.Annotations["blinkt.furikuri.net/color"]
	if !ok {
		return
	}

	blinktIP, ok := c.blinkts[pod.Status.HostIP]
	if !ok {
		return
	}

	color := colors[val]
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(LedColor{Red: 0, Blue: 0, Green: 0, Led: color.Led})
	res, err := http.Post("http://"+blinktIP+":5000/set_color", "application/json; charset=utf-8", b)
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, res.Body)
}

func (c *BlinktController) PodAdded(obj interface{}) {
	pod := obj.(*core_v1.Pod)

	if pod.Labels["name"] == "http-blinkt" {
		fmt.Println("is http blinkt " + pod.Status.PodIP)
		c.blinkts[pod.Status.HostIP] = pod.Status.PodIP
	}

	val, ok := pod.Annotations["blinkt.furikuri.net/color"]
	if !ok {
		return
	}

	blinktIP, ok := c.blinkts[pod.Status.HostIP]
	if !ok {
		return
	}

	color := colors[val]
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(color)
	res, err := http.Post("http://"+blinktIP+":5000/set_color", "application/json; charset=utf-8", b)
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, res.Body)
}

type LedColor struct {
	Red   int32 `json:"red"`
	Green int32 `json:"green"`
	Blue  int32 `json:"blue"`
	Led   int32 `json:"led"`
}
