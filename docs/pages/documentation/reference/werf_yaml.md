---
title: werf.yaml
permalink: documentation/reference/werf_yaml.html
sidebar: documentation
toc: false
---

{% include documentation/reference/werf_yaml/table.html %}

## Project name

`project` defines unique project name of your application. Project name affects build cache image names, Kubernetes Namespace, Helm Release name and other derived names (see [deploy to Kubernetes for detailed description]({{ site.baseurl }}/documentation/reference/configuration/deploy_into_kubernetes.html)). This is single required field of meta configuration.

Project name should be unique within group of projects that shares build hosts and deployed into the same Kubernetes cluster (i.e. unique across all groups within the same gitlab).

Project name must be maximum 50 chars, only lowercase alphabetic chars, digits and dashes are allowed.

**WARNING**. You should never change project name, once it has been set up, unless you know what you are doing.

Changing project name leads to issues:
1. Invalidation of build cache. New images must be built. Old images must be cleaned up from local host and Docker registry manually.
2. Creation of completely new Helm Release. So if you already had deployed your application, then changed project name and deployed it again, there will be created another instance of the same application.

werf cannot automatically resolve project name change. Described issues must be resolved manually.

## Deploy

### Release name

werf allows to define a custom release name template, which [used during deploy process]({{ site.baseurl }}/documentation/advanced/helm/basics.html#release-name) to generate a release name:

```yaml
project: PROJECT_NAME
configVersion: 1
deploy:
  helmRelease: TEMPLATE
  helmReleaseSlug: false
```

`deploy.helmRelease` is a Go template with `[[` and `]]` delimiters. There are `[[ project ]]`, `[[ env ]]` functions support. Default: `[[ project ]]-[[ env ]]`.

`deploy.helmReleaseSlug` defines whether to apply or not [slug]({{ site.baseurl }}/documentation/advanced/helm/basics.html#slugging-the-release-name) to generated helm release name. Default: `true`.

`TEMPLATE` as well as any value of the config can include [werf Go templates functions]({{ site.baseurl }}/documentation/reference/configuration/introduction.html#go-templates). E.g. you can mix the value with an environment variable:

{% raw %}
```yaml
deploy:
  helmRelease: >-
    [[ project ]]-{{ env "HELM_RELEASE_EXTRA" }}-[[ env ]]
```
{% endraw %}

### Kubernetes namespace

werf allows to define a custom Kubernetes namespace template, which [used during deploy process]({{ site.baseurl }}/documentation/advanced/helm/basics.html#kubernetes-namespace) to generate a Kubernetes Namespace:

```yaml
project: PROJECT_NAME
configVersion: 1
deploy:
  namespace: TEMPLATE
  namespaceSlug: true|false
```

`deploy.namespace` is a Go template with `[[` and `]]` delimiters. There are `[[ project ]]`, `[[ env ]]` functions support. Default: `[[ project ]]-[[ env ]]`.

`deploy.namespaceSlug` defines whether to apply or not [slug]({{ site.baseurl }}/documentation/advanced/helm/basics.html#slugging-kubernetes-namespace) to generated kubernetes namespace. Default: `true`.

## Cleanup

### Configuring cleanup policies

The cleanup configuration consists of a set of policies called `keepPolicies`. They are used to select relevant images using the git history. Thus, during a [cleanup]({{ site.baseurl }}/documentation/advanced/cleanup.html#git-history-based-cleanup-algorithm), __images not meeting the criteria of any policy are deleted__.

Each policy consists of two parts:

- `references` defines a set of references, git tags, or git branches to perform scanning on.
- `imagesPerReference` defines the limit on the number of images for each reference contained in the set.

Each policy should be linked to some set of git tags (`tag: <string|/REGEXP/>`) or git branches (`branch: <string|/REGEXP/>`). You can specify the name/group of a reference using the [Golang's regular expression syntax](https://golang.org/pkg/regexp/syntax/#hdr-Syntax).

```yaml
tag: v1.1.1
tag: /^v.*$/
branch: master
branch: /^(master|production)$/
```

> When scanning, werf searchs for the provided set of git branches in the origin remote references, but in the configuration, the  `origin/` prefix is omitted in branch names.

You can limit the set of references on the basis of the date when the git tag was created or the activity in the git branch. The `limit` group of parameters allows the user to define flexible and efficient policies for various workflows.

```yaml
- references:
    branch: /^features\/.*/
    limit:
      last: 10
      in: 168h
      operator: And
``` 

In the example above, werf selects no more than 10 latest branches that have the `features/` prefix in the name and have shown any activity during the last week.

- The `last: <int>` parameter allows you to select n last references from those defined in the `branch` / `tag`.
- The `in: <duration string>` parameter (you can learn more about the syntax in the [docs](https://golang.org/pkg/time/#ParseDuration)) allows you to select git tags that were created during the specified period or git branches that were active during the period. You can also do that for the specific set of `branches` / `tags`.
- The `operator: <And|Or>` parameter defines if references should satisfy both conditions or either of them (`And` is set by default).

When scanning references, the number of images is not limited by default. However, you can configure this behavior using the `imagesPerReference` set of parameters:

```yaml
imagesPerReference:
  last: <int>
  in: <duration string>
  operator: <And|Or>
```

- The `last: <int>` parameter defines the number of images to search for each reference. Their amount is unlimited by default (`-1`).
- The `in: <duration string>` parameter (you can learn more about the syntax in the [docs](https://golang.org/pkg/time/#ParseDuration)) defines the time frame in which werf searches for images.
- The `operator: <And|Or>` parameter defines what images will stay after applying the policy: those that satisfy both conditions or either of them (`And` is set by default).

> In the case of git tags, werf checks the HEAD commit only; the value of `last`>1 does not make any sense and is invalid

When describing a group of policies, you have to move from the general to the particular. In other words, `imagesPerReference` for a specific reference will match the latest policy it falls under:

```yaml
- references:
    branch: /.*/
  imagesPerReference:
    last: 1
- references:
    branch: master
  imagesPerReference:
    last: 5
```

In the above example, the _master_ reference matches both policies. Thus, when scanning the branch, the `last` parameter will equal to 5.

### Default policies

If there are no custom cleanup policies defined in `werf.yaml`, werf uses default policies configured as follows:

```yaml
cleanup:
  keepPolicies:
  - references:
      tag: /.*/
      limit:
        last: 10
  - references:
      branch: /.*/
      limit:
        last: 10
        in: 168h
        operator: And
    imagesPerReference:
      last: 2
      in: 168h
      operator: And
  - references:  
      branch: /^(master|staging|production)$/
    imagesPerReference:
      last: 10
``` 

Let us examine each policy individually:

1. Keep an image for the last 10 tags (by date of creation).
2. Keep no more than two images published over the past week, for no more than 10 branches active over the past week.
3. Keep the 10 latest images for master, staging, and production branches.

## Image section

Building image from Dockerfile is the easiest way to start using werf in an existing project.
Minimal `werf.yaml` below describes an image named `example` related with a project `Dockerfile`:

```yaml
project: my-project
configVersion: 1
---
image: example
dockerfile: Dockerfile
```

To specify some images from one Dockerfile:

```yaml
image: backend
dockerfile: Dockerfile
target: backend
---
image: frontend
dockerfile: Dockerfile
target: frontend
```

And also from different Dockerfiles:

```yaml
image: backend
dockerfile: dockerfiles/DockerfileBackend
---
image: frontend
dockerfile: dockerfiles/DockerfileFrontend
```

### Naming

{% include /configuration/stapel_image/naming.md %}
