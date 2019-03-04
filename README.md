Zundoko Operator
=========================

Zundoko Operator is a [Kubernetes operator](https://github.com/operator-framework/awesome-operators) built with [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder), that plays ZundokoKiyoshi.

Deploy to Kubernetes cluster
-----------------------------------

You need [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) to deploy Zundoko Operator with it.

1. `git clone https://github.com/kaitoy/zundoko-operator.git`
2. `cd zundoko-operator`
3. `kubectl apply -f config/crds/zundokokiyoshi_v1beta1_hikawa.yaml`
4. `kubectl apply -f config/crds/zundokokiyoshi_v1beta1_kiyoshi.yaml`
5. `kubectl apply -f config/crds/zundokokiyoshi_v1beta1_zundoko.yaml`
6. `kubectl apply -f zundoko-operator.yaml`

Usage
-----

### Start a ZundokoKiyoshi

1. Write a Hikawa manifest.

    e.g.)

    ```yaml
    apiVersion: zundokokiyoshi.kaitoy.github.com/v1beta1
    kind: Hikawa
    metadata:
      labels:
        controller-tools.k8s.io: "1.0"
      name: hikawa-sample
    spec:
      intervalMillis: 500
    ```

2. Apply the manifest to start a ZundokoKiyoshi.
    * Zundoko Operator starts to create Zundokos one by one with the Say fields set to "Zun" or "Doko" randomly.
    * When 4 "Zun"s are created and then "Doko" is done, Zundoko Operator creates Kiyoshi.
3. See Zundokos by: `kubectl get zundoko`

    ```sh
    [root@k8s-master ~]# kubectl get zundoko
    NAME                        SAY    AGE
    hikawa-sample-zundoko-001   Doko   6m
    hikawa-sample-zundoko-002   Doko   6m
    hikawa-sample-zundoko-003   Doko   6m
    hikawa-sample-zundoko-004   Zun    6m
    hikawa-sample-zundoko-005   Doko   6m
    hikawa-sample-zundoko-006   Doko   6m
    ```

4. See Kiyoshi by: `kubectl get kiyoshi`

    ```sh
    [root@k8s-master ~]# kubectl get kiyoshi
    NAME                    SAY        AGE
    hikawa-sample-kiyoshi   Kiyoshi!   7m
    ```

5. Delete the manifest to clean a ZundokoKiyoshi.

Development
-----------

You need [Docker](https://www.docker.com/) to build Zundoko Operator container image, and [kustomize](https://github.com/kubernetes-sigs/kustomize) to generate a Kubernetes manifest file.

1. Build container image
    1. `git clone https://github.com/kaitoy/zundoko-operator.git`
    2. `cd zundoko-operator`
    3. `docker build .`
2. Generate Kubernetes manifest file
    1. `kustomize build config/default > zundoko-operator.yaml`
