
A minimal helm chart for deploying wireguard in Kubernetes.

## Installation

#### Pre-requisites

Refer to the [Key Generation](https://www.wireguard.com/quickstart/#key-generation) in the WireGuard QuickStart. After generating the private key, you will need to populate a secret named `wireguard-secret` in the namespace you intend to install the chart. The key name in the secret needs to be named `privatekey`. The helm chart will populate the wireguard configuration with this value.

```
helm repo add wg-k8s https://serrichard.github.io/wg-k8s
helm repo update
```

```
version=""
helm upgrade --install wg wg-k8s/wg-k8s -f values.yaml --version $version --namespace wireguard
```

```
# For local
helm upgrade --install wgk8s ./wg-k8s -f .values.yaml --namespace wireguard
```

```
helm uninstall wgk8s
```

## Helm Values

| Key                                          | Description                                                                            | Type     | Default                       |
| -------------------------------------------- | -------------------------------------------------------------------------------------- | -------- | ----------------------------- |
| `image.repository`                           | Container image repository                                                             | `string` | `linuxserver/wireguard`       |                 |
| `image.pullPolicy`                           | Image pull policy                                                                      | `string` | `IfNotPresent`                |                 |
| `image.tag`                                  | Image tag                                                                              | `string` | `1.0.20250521`                |                 |
| `config.interface.address`                   | IP address of the interface                                                            | `string` | `10.3.0.1/24`                 |                 |
| `config.interface.dns`                       | DNS server for the interface                                                           | `string` | `10.152.183.10`               |                 |
| `config.peers`                               | List of peers. Each peer can define `publicKey`, `allowedIPs`, and an optional `endpoint` | `array`  | `[]`                          |                 |
| `gateway.create`                             | Create gateway deployment                                                              | `bool`   | `true`                        |                 |
| `gateway.image.repository`                   | Gateway image repository                                                               | `string` | `vimagick/tinyproxy`          |                 |
| `gateway.image.pullPolicy`                   | Gateway image pull policy                                                              | `string` | `IfNotPresent`                |                 |
| `gateway.image.tag`                          | Gateway image tag                                                                      | `string` | `latest`                      |                 |
| `gateway.service.type`                       | Gateway service type                                                                   | `string` | `ClusterIP`                   |                 |
| `gateway.service.ports`                      | Gateway service ports                                                                  | `array`  | `[]`                          |                 |
| `replicaCount`                               | Number of replicas for the deployment                                                  | `int`    | `1`                           |                 |
| `imagePullSecrets`                           | Secrets for pulling private images                                                     | `array`  | `[]`                          |                 |
| `nameOverride`                               | Override chart name                                                                    | `string` | `""`                          |                 |
| `fullnameOverride`                           | Override full resource names                                                           | `string` | `""`                          |                 |
| `serviceAccount.create`                      | Whether to create a ServiceAccount                                                     | `bool`   | `true`                        |                 |
| `serviceAccount.automount`                   | Automatically mount ServiceAccount token                                               | `bool`   | `true`                        |                 |
| `serviceAccount.annotations`                 | Annotations for ServiceAccount                                                         | `object` | `{}`                          |                 |
| `serviceAccount.name`                        | Name of ServiceAccount                                                                 | `string` | `"wg-k8s"`                    |                 |
| `podAnnotations`                             | Additional pod annotations                                                             | `object` | `{}`                          |                 |
| `podLabels`                                  | Additional pod labels                                                                  | `object` | `{}`                          |                 |
| `podSecurityContext`                         | Pod-level security context                                                             | `object` | `{}`                          |                 |
| `securityContext.privileged`                 | Run container in privileged mode                                                       | `bool`   | `true`                        |                 |
| `securityContext.capabilities.add`           | Linux capabilities to add                                                              | `array`  | `["NET_ADMIN", "SYS_MODULE"]` |                 |
| `service.type`                               | Service type                                                                           | `string` | `ClusterIP`                   |                 |
| `service.port`                               | Service port                                                                           | `int`    | `4500`                        |                 |
| `resources`                                  | Resource requests and limits                                                           | `object` | `{}`                          |                 |
| `livenessProbe.exec.command`                 | Command for liveness probe                                                             | `array`  | `["/bin/sh", "-c", "ss -lnu   | grep -q 4500"]` |
| `livenessProbe.initialDelaySeconds`          | Delay before liveness probe starts                                                     | `int`    | `5`                           |                 |
| `livenessProbe.periodSeconds`                | Liveness probe interval                                                                | `int`    | `10`                          |                 |
| `livenessProbe.failureThreshold`             | Failure threshold for liveness probe                                                   | `int`    | `3`                           |                 |
| `readinessProbe.exec.command`                | Command for readiness probe                                                            | `array`  | `["/bin/sh", "-c", "ss -lnu   | grep -q 4500"]` |
| `readinessProbe.initialDelaySeconds`         | Delay before readiness probe starts                                                    | `int`    | `5`                           |                 |
| `readinessProbe.periodSeconds`               | Readiness probe interval                                                               | `int`    | `10`                          |                 |
| `readinessProbe.failureThreshold`            | Failure threshold for readiness probe                                                  | `int`    | `3`                           |                 |
| `volumes`                                    | Additional volumes for pods                                                            | `array`  | `[]`                          |                 |
| `volumeMounts`                               | Additional volume mounts for pods                                                      | `array`  | `[]`                          |                 |
| `nodeSelector`                               | Node selector for pods                                                                 | `object` | `{}`                          |                 |
| `tolerations`                                | Tolerations for pods                                                                   | `array`  | `[]`                          |                 |
| `affinity`                                   | Affinity rules for pods                                                                | `object` | `{}`                          |                 |