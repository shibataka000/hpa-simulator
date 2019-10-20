# HPA Simulator

## Description
This tool simulate HorizontalPodAutoscaler and output some information to know it's internal behavior.
This tool is developed to debug HorizontalPodAutoscaler behavior.

This tool support only resource metrics. This tool doesn't support custom metrics and external metrics.

## Requirement
Go 1.13.3 or later

## Usage
```
$ hpasimulator --help
NAME:
   hpasimulator - Simulate HorizontalPodAutoscaler and output some information to know it's internal behavior

USAGE:
   hpasimulator [flags]

VERSION:
   v0.0.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --kubeconfig value           Path to the kubeconfig file to use for CLI requests. (default: "/home/ubuntu/.kube/config")
   --context value              The name of the kubeconfig context to use
   --namespace value, -n value  If present, the namespace scope for this CLI request (default: "default")
   --selector value, -l value   Selector (label query) to filter on, supports '='.(e.g. -l key1=value1,key2=value2)
   --help, -h                   show help
   --version, -v                print the version
```

For example,

```
$ hpasimulator -l app=ubuntu
2019/10/29 07:17:53 [Current Replicas] 1
2019/10/29 07:17:53 [Container Metrics] ubuntu-5f77f46cf4-drnsr ubuntu1804: 0
2019/10/29 07:17:53 [Pod Resource Utilization] ubuntu-5f77f46cf4-drnsr: 0 / 100 = 0 %
2019/10/29 07:17:53 [Deployment Resource Utilization] 0 %
2019/10/29 07:17:53 [Scale] 1 -> 0
2019/10/29 07:17:53 ========================================
2019/10/29 07:17:53 [Current Replicas] 0
2019/10/29 07:17:53 [Container Metrics] ubuntu-5f77f46cf4-drnsr ubuntu1804: 0
2019/10/29 07:17:53 [Pod Resource Utilization] ubuntu-5f77f46cf4-drnsr: 0 / 100 = 0 %
2019/10/29 07:17:53 [Deployment Resource Utilization] 0 %
2019/10/29 07:17:53 ========================================
```

### Notes
- The only way to narrow down target pods of HorizontalPodAutoscaler in simulation is specify pod labels by `-l` or `--selector` options. This tool doesn't consider Deployment and HorizontalPodAutoscaler. 
- This tool doesn't simulate HorizontalPodAutoscaler fully.
   - This tool doesn't support `minReplicas` and `maxReplicas`. This tool scale unlimitedly in simulation.
   - This tool doesn't get `currentReplicas` from Deployment. This tool assume `currentReplicas` is `1` at beggining.
   - HorizontalPodAutoscaler stabilize scale-in, but this tool doesn't. So Deployment scale in more quickly than actual.

## Install
```bash
git clone https://github.com/shibataka000/hpa-simulator
cd hpa-simulator
make install
```
