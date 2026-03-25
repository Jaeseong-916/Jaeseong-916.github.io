---
title: "CAPI + OpenStack + Cilium 기반 Kubernetes 클러스터 구축 가이드"
date: 2026-03-23
categories:
  - kubernetes
tags:
  - kubernetes
  - cilium
  - CAPI
  - openstack
  - CNI
  - kube-proxy-replacement
---

Cluster API(CAPI)와 OpenStack 인프라 프로바이더를 활용하여 **Cilium CNI 기반**의 Kubernetes 클러스터를 구축하는 과정을 정리합니다.

## 환경 정보

| 항목 | 버전/정보 |
|------|-----------|
| Management Cluster | Kind v0.27.0 (Kubernetes v1.32.2) |
| Workload Cluster | Kubernetes v1.32.4 |
| CAPI | v1.12.3 |
| CAPO (OpenStack Provider) | v0.14.1 |
| CNI | Cilium v1.19.1 (kube-proxy replacement 모드) |
| Node OS | Rocky Linux 9.7 |
| 구성 | Control Plane 3대 + Worker 3대 |

## 아키텍처

```
┌─────────────────────────────────────┐
│  CAPI Server (10.10.11.81)          │
│  ┌───────────────────────────────┐  │
│  │  Kind (Management Cluster)    │  │
│  │  ├─ CAPI Controller           │  │
│  │  ├─ CAPO Controller (v0.14.1) │  │
│  │  ├─ KubeAdm Bootstrap         │  │
│  │  └─ ORC CRDs                  │  │
│  └───────────────────────────────┘  │
└──────────────┬──────────────────────┘
               │ OpenStack API
               ▼
┌─────────────────────────────────────┐
│  OpenStack Infrastructure           │
│  ┌─────────────────────────────┐    │
│  │  k8s-cilium Cluster          │   │
│  │  CP-1  CP-2  CP-3           │   │
│  │  WK-1  WK-2  WK-3           │   │
│  │  CNI: Cilium v1.19.1        │   │
│  │  (kube-proxy replacement)   │   │
│  └─────────────────────────────┘    │
└─────────────────────────────────────┘
```

## 1. 사전 준비 - 도구 설치

```bash
# Docker
sudo dnf config-manager --add-repo https://download.docker.com/linux/rhel/docker-ce.repo
sudo dnf install -y docker-ce docker-ce-cli containerd.io
sudo systemctl enable --now docker

# kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl && sudo mv kubectl /usr/local/bin/

# Kind
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.27.0/kind-linux-amd64
chmod +x ./kind && sudo mv ./kind /usr/local/bin/

# clusterctl
curl -Lo ./clusterctl https://github.com/kubernetes-sigs/cluster-api/releases/latest/download/clusterctl-linux-amd64
chmod +x ./clusterctl && sudo mv ./clusterctl /usr/local/bin/

# Helm
curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Cilium CLI
CILIUM_CLI_VERSION=$(curl -s https://raw.githubusercontent.com/cilium/cilium-cli/main/stable.txt)
curl -L --fail --remote-name-all \
  "https://github.com/cilium/cilium-cli/releases/download/${CILIUM_CLI_VERSION}/cilium-linux-amd64.tar.gz{,.sha256sum}"
sha256sum --check cilium-linux-amd64.tar.gz.sha256sum
sudo tar xzvfC cilium-linux-amd64.tar.gz /usr/local/bin
```

## 2. Management Cluster 생성

```bash
kind create cluster --config kind-cluster-with-extramounts.yaml
```

## 3. CAPI Provider 초기화

```bash
clusterctl init --infrastructure openstack

# ORC CRD 설치 (CAPO v0.14.x 필수)
kubectl apply -f https://github.com/k-orc/openstack-resource-controller/releases/latest/download/install.yaml
```

> **주의**: ORC CRD를 설치하지 않으면 CAPO 컨트롤러가 `CrashLoopBackOff` 상태에 빠집니다.

## 4. OpenStack API 접근 설정

Kind 컨테이너 내부에서 OpenStack API에 접근할 수 있어야 합니다.

### CoreDNS에 커스텀 DNS 추가

```
.:53 {
    hosts {
       10.10.0.10 stack.osci.cloud    # OpenStack API 내부 IP 매핑
       fallthrough
    }
    kubernetes cluster.local in-addr.arpa ip6.arpa { ... }
}
```

```bash
kubectl rollout restart deployment/coredns -n kube-system
```

## 5. 핵심 설정 - kube-proxy 비활성화

Cilium이 kube-proxy를 대체하므로, `initConfiguration.skipPhases`에서 `addon/kube-proxy`를 추가합니다.

```yaml
apiVersion: controlplane.cluster.x-k8s.io/v1beta2
kind: KubeadmControlPlane
spec:
  kubeadmConfigSpec:
    initConfiguration:
      skipPhases:
        - addon/kube-proxy          # Cilium이 kube-proxy를 대체
```

> **Calico와의 차이점**: 기존 Calico 기반 매니페스트에는 `skipPhases`가 없습니다.

## 6. Cilium용 보안 그룹 규칙

Calico와 Cilium은 사용하는 포트가 다릅니다.

```yaml
managedSecurityGroups:
  allNodesSecurityGroupRules:
    - name: VXLAN (Cilium)
      protocol: udp
      portRangeMin: 8472
      portRangeMax: 8472
      direction: ingress
      remoteManagedGroups: [controlplane, worker]
    - name: Health (Cilium)
      protocol: tcp
      portRangeMin: 4240
      portRangeMax: 4240
      direction: ingress
      remoteManagedGroups: [controlplane, worker]
    - name: Hubble (Cilium)
      protocol: tcp
      portRangeMin: 4244
      portRangeMax: 4244
      direction: ingress
      remoteManagedGroups: [controlplane, worker]
```

## 7. Cilium CNI 설치

```bash
clusterctl get kubeconfig k8s-cilium > k8s-cilium-kubeconfig.yaml

API_SERVER_IP=$(kubectl get openstackcluster k8s-cilium \
  -o jsonpath='{.status.apiServerLoadBalancer.ip}')

cilium install \
  --kubeconfig k8s-cilium-kubeconfig.yaml \
  --set kubeProxyReplacement=true \
  --set k8sServiceHost=$API_SERVER_IP \
  --set k8sServicePort=6443 \
  --set ipam.operator.clusterPoolIPv4PodCIDRList=192.168.0.0/16
```

## 8. 최종 확인

```bash
$ kubectl --kubeconfig k8s-cilium-kubeconfig.yaml get nodes
NAME                             STATUS   ROLES           VERSION
k8s-cilium-control-plane-4xpqv   Ready    control-plane   v1.32.4
k8s-cilium-control-plane-k6zf8   Ready    control-plane   v1.32.4
k8s-cilium-control-plane-vcmz2   Ready    control-plane   v1.32.4
k8s-cilium-md-0-tf6bl-dw5gv      Ready    <none>          v1.32.4
k8s-cilium-md-0-tf6bl-ffx9l      Ready    <none>          v1.32.4
k8s-cilium-md-0-tf6bl-kxqtr      Ready    <none>          v1.32.4

$ cilium status --kubeconfig k8s-cilium-kubeconfig.yaml
    /¯¯\
 /¯¯\__/¯¯\    Cilium:             OK
 \__/¯¯\__/    Operator:           OK
 /¯¯\__/¯¯\    Envoy DaemonSet:    OK
 \__/¯¯\__/    Hubble Relay:       disabled
    \__/       ClusterMesh:        disabled
```

## Calico vs Cilium 비교 (CAPI 매니페스트 관점)

| 항목 | Calico | Cilium |
|------|--------|--------|
| kube-proxy | 유지 | `skipPhases: addon/kube-proxy`로 제거 |
| 보안 그룹 | BGP(TCP/179) + IP-in-IP(Protocol 4) | VXLAN(UDP/8472) + Health(TCP/4240) + Hubble(TCP/4244) |
| CNI 설치 | `kubectl apply` (Calico manifest) | `cilium install --set kubeProxyReplacement=true` |
| 추가 설정 | 없음 | `k8sServiceHost`, `k8sServicePort` 필요 |

## 트러블슈팅

### CAPO CrashLoopBackOff - ORC CRD 누락
```bash
kubectl apply -f https://github.com/k-orc/openstack-resource-controller/releases/latest/download/install.yaml
kubectl rollout restart deployment/capo-controller-manager -n capo-system
```

### OpenStack API 연결 실패
CoreDNS에 hosts 플러그인으로 내부 IP 매핑을 추가합니다.

### TLS 인증서 만료
`clouds.yaml`에 `verify: false` 추가 (임시 조치). 프로덕션에서는 인증서 갱신을 권장합니다.

### Floating IP 접근 불가
`externalNetwork.id`를 CAPI 서버와 동일한 네트워크로 변경합니다.
