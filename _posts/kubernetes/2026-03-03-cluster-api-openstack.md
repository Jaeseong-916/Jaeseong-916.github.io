---
title: "Cluster API(CAPI)로 OpenStack에 Kubernetes 클러스터 배포하기"
date: 2026-03-03
categories:
  - kubernetes
tags:
  - kubernetes
  - cluster-api
  - CAPI
  - openstack
  - kind
---

Cluster API(CAPI)와 OpenStack 인프라 프로바이더를 활용하여 Kubernetes 클러스터를 구축하는 과정을 정리합니다.

## 환경 구성

- Management Cluster: Kind
- Infrastructure Provider: OpenStack (CAPO)
- 사내 KVM 서버의 VM에 Management Cluster를 구성하여 OpenStack에 배포하는 방식

### 사전 준비

- kubectl
- Kind + Docker
- Helm
- clusterctl

## 도구 설치

### kubectl

```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/bin/kubectl
```

### Kind

```bash
[ $(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.30.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/bin/kind
```

### Helm

```bash
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-4
chmod 700 get_helm.sh
./get_helm.sh
```

## Kind 클러스터 구성

```yaml
cat > kind-cluster-with-extramounts.yaml <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  ipFamily: dual
nodes:
- role: control-plane
  extraMounts:
    - hostPath: /var/run/docker.sock
      containerPath: /var/run/docker.sock
EOF

kind create cluster --config kind-cluster-with-extramounts.yaml
```

## Clusterctl 설치 및 초기화

```bash
curl -L https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.11.3/clusterctl-linux-amd64 -o clusterctl
sudo install -o root -g root -m 0755 clusterctl /usr/local/bin/clusterctl

# ORC 설치 (CAPO >=v0.12 필수)
kubectl apply -f https://github.com/k-orc/openstack-resource-controller/releases/latest/download/install.yaml

# OpenStack Infrastructure Provider 초기화
clusterctl init --infrastructure openstack
```

## 클러스터 생성

### 이미지 준비

Cluster API를 사용하여 cluster를 배포하려면 클러스터가 배포될 OS 이미지가 필요합니다. [image-builder](https://image-builder.sigs.k8s.io/capi/providers/openstack.html)를 사용하여 이미지를 만듭니다.

### 환경 변수 설정

```bash
wget https://raw.githubusercontent.com/kubernetes-sigs/cluster-api-provider-openstack/master/templates/env.rc -O /tmp/env.rc
source /tmp/env.rc /root/clouds.yaml infra

export OPENSTACK_DNS_NAMESERVERS="8.8.8.8"
export OPENSTACK_FAILURE_DOMAIN="nova"
export OPENSTACK_CONTROL_PLANE_MACHINE_FLAVOR="m1.large"
export OPENSTACK_NODE_MACHINE_FLAVOR="m1.large"
export OPENSTACK_IMAGE_NAME="rockylinux-9-kube-v1.32.4"
export OPENSTACK_SSH_KEY_NAME="jsshin"
export OPENSTACK_EXTERNAL_NETWORK_ID="<network-id>"
```

### 클러스터 매니페스트 생성 및 배포

```bash
clusterctl generate cluster capi-quickstart \
  --kubernetes-version v1.32.4 \
  --control-plane-machine-count=3 \
  --worker-machine-count=3 \
  > capi-quickstart.yaml

kubectl apply -f capi-quickstart.yaml
```

## 클러스터 확인

```bash
kubectl get cluster
clusterctl describe cluster capi-quickstart
clusterctl get kubeconfig capi-quickstart > capi-quickstart.kubeconfig
```

## Cloud Provider 설치

### cloud.conf 구성 및 배포

```bash
kubectl --kubeconfig=./capi-quickstart.kubeconfig -n kube-system \
  create secret generic cloud-config --from-file=cloud.conf

kubectl apply --kubeconfig=./capi-quickstart.kubeconfig \
  -f https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/controller-manager/cloud-controller-manager-roles.yaml
kubectl apply --kubeconfig=./capi-quickstart.kubeconfig \
  -f https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/controller-manager/cloud-controller-manager-role-bindings.yaml
kubectl apply --kubeconfig=./capi-quickstart.kubeconfig \
  -f https://raw.githubusercontent.com/kubernetes/cloud-provider-openstack/master/manifests/controller-manager/openstack-cloud-controller-manager-ds.yaml
```

## CNI 배포 (Calico)

```bash
kubectl --kubeconfig=./capi-quickstart.kubeconfig \
  apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/calico.yaml
```

## 클러스터 제거

```bash
kubectl delete cluster <cluster_name>
```
