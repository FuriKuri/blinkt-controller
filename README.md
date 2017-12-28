# blinkt-controller

![rpi cluster with led](https://raw.githubusercontent.com/FuriKuri/blinkt-controller/master/rpi.jpg)

This is a simple kubernetes controller, which uses the Blinkt! LEDs to visualize which pod was scheduled on which node.

## Deploy

First deploy the http blinkt container, which controlls the LEDs of the raspberry pi.
After this deploy the controller, which listen on Pod changes of the kubernetes cluster.

```
kubectl apply -f https://raw.githubusercontent.com/FuriKuri/blinkt-controller/master/server.yaml
kubectl apply -f https://raw.githubusercontent.com/FuriKuri/blinkt-controller/master/controller.yaml
```

## Usage

If a pod has an annotation `blinkt.furikuri.net/color: color`, the corresponding raspberry pi will activate the LED `color`.

The following example will deploy a pod with a nginx container. If the pod will scheduled on a node, the one LED will be appear green:

```
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  annotations:
    blinkt.furikuri.net/color: green
spec:
  containers:
  - name: nginx
    image: nginx
```

Valid colors:
* blue
* green
* red
* yellow
* orange
* violet
* brown
* grey

## Example

The following yaml file contains deployments with a replica set of 2 for each color.
It just starts a pause container.

```
kubectl apply -f https://raw.githubusercontent.com/FuriKuri/blinkt-controller/master/colors-example.yaml
```