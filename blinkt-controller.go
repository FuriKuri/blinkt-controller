package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"log"
	"net/http"
	"os"
)

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

	color := Colors[val]
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

	color := Colors[val]
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(color)
	res, err := http.Post("http://"+blinktIP+":5000/set_color", "application/json; charset=utf-8", b)
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, res.Body)
}
