{% if include.header %}
{% assign header = include.header %}
{% else %}
{% assign header = "###" %}
{% endif %}

Search the Helm Hub or an instance of Monocular for Helm charts.

The Helm Hub provides a centralized search for publicly available distributed
charts. It is maintained by the Helm project. It can be visited at
[https://hub.helm.sh](https://hub.helm.sh)

Monocular is a web-based application that enables the search and discovery of
charts from multiple Helm Chart repositories. It is the codebase that powers the
Helm Hub. You can find it at [https://github.com/helm/monocular](https://github.com/helm/monocular)


{{ header }} Syntax

```shell
werf helm-v3 search hub [keyword] [flags] [options]
```

{{ header }} Options

```shell
      --endpoint='https://hub.helm.sh':
            monocular instance to query for charts
  -h, --help=false:
            help for hub
      --max-col-width=50:
            maximum column width for output table
  -o, --output=table:
            prints the output in the specified format. Allowed values: table, json, yaml
```

{{ header }} Options inherited from parent commands

```shell
      --hooks-status-progress-period=5:
            Hooks status progress period in seconds. Set 0 to stop showing hooks status progress.   
            Defaults to $WERF_HOOKS_STATUS_PROGRESS_PERIOD_SECONDS or status progress period value
      --kube-config='':
            Kubernetes config file path (default $WERF_KUBE_CONFIG or $WERF_KUBECONFIG or           
            $KUBECONFIG)
      --kube-config-base64='':
            Kubernetes config data as base64 string (default $WERF_KUBE_CONFIG_BASE64 or            
            $WERF_KUBECONFIG_BASE64 or $KUBECONFIG_BASE64)
      --kube-context='':
            Kubernetes config context (default $WERF_KUBE_CONTEXT)
  -n, --namespace='':
            namespace scope for this request
      --status-progress-period=5:
            Status progress period in seconds. Set -1 to stop showing status progress. Defaults to  
            $WERF_STATUS_PROGRESS_PERIOD_SECONDS or 5 seconds
```

