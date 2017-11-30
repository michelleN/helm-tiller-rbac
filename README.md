# Helm Secure-Tiller Plugin
This Helm plugin allows you to add an RBAC profile to a Tiller in a Kubernetes namespace. An **RBAC Profile** is a Helm chart that consists of a Kubernetes Role and RoleBinding definition. This plugin is designed to help the team of operators that set up multiple Tillers in their cluster (one Tiller per namespace) ensure that a Tiller is locked down to specific actions on specific Kubernetes resources in a given namespace.


## Usage
```console
$ helm secure-tiller [flags] RBAC_PROFILE
```

### Flags
```
     --namespace string   namespace of Tiller to apply profile (default "default")
```

## Install
```console
$ helm plugin install https://github.com/michelleN/helm-secure-tiller
```
The above will fetch the latest binary release of `helm secure-tiller` and install it.

### Developer (From Source) Install

If you would like to handle the build yourself, instead of fetching a binary, this is how recommend doing it.

First, set up your environment:

You need to have Go installed. Make sure to set $GOPATH
If you don't have Glide installed, this will install it into $GOPATH/bin for you.
Clone this repo into your $GOPATH using git.

```console
cd $GOPATH/src/github.com/michelleN # mkdir as needed
git clone https://github.com/michelleN/helm-secure-tiller
```
Then run the following to get running.

```
$ cd helm-secure-tiller
$ make bootstrap build
$ SKIP_BIN_INSTALL=1 helm plugin install $GOPATH/src/github.com/michelleN/helm-secure-tiller
```
That last command will skip fetching the binary install and use the one you built.


## Example Workflow
First, you'll need to set up Tiller in the namespace you want.
```console
$ kubectl create namespace dev-team
# install tiller in a namespace like this:
$ helm init --tiller-namespace dev-team
```

Then, use the secure-tiller plugin to apply the example dev-team RBAC profile in this repo.
```console
$ helm secure-tiller --namespace dev-team example-profiles/dev-team/
serviceaccount "dev-team-rbac-profile" created
rolebinding "dev-team-tiller-binding" created
role "dev-team" created

Congrats! Your Tiller in the dev-team namespace is secured with the dev-team-rbac-profile service account
You can verify your Tiller config by running this command
	$ kubectl -n dev-team get deployment tiller-deploy -o json
```

You should now have a Tiller in the `dev-team` namespace that is only allowed to do the things specified in `example-profiles/dev-team/templates/role-tiller.yaml`
