# wg-k8s

Minimal helm chart for deploying wireguard in Kubernetes.

```
helm repo add wg-k8s https://serrichard.github.io/wg-k8s
helm repo update
```

```
version=2025.10.1-rc1
helm upgrade --install wg wg-k8s/wg-k8s -f values.yaml --version $version --namespace wireguard
```

```
helm uninstall wgk8s
```