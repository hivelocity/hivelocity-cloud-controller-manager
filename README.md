# Hivelocity Cloud Controller Manager

This code provides a [cloud controller manager](https://kubernetes.io/docs/tasks/administer-cluster/running-cloud-controller/) for
the provider [hivelocity.net](https://www.hivelocity.net/).

# Tests

To run the tests you need an API key in the file `.envrc`. See `.envrc-example`.

# Releasing

Increase the version in `charts/ccm-hivelocity/Chart.yaml`

On push to the main branch, the Github Workflow build.yml will trigger "actions/upload-artifact".

On push to the main branch, the Github Workflow release.yml will trigger "helm/chart-releaser-action".
