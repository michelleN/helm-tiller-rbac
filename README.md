# helm-secure-tiller
enable rbac profiles for Tiller

# example usage
```console

# grab repo & build project ... will turn into helm plugin soon...
$ git clone git@github.com:michelleN/helm-secure-tiller.git
$ cd helm-secure-tiller
$ make bootstrap build

# install tiller in a namespace like this:
$ kubectl create namespace dev-team
$ helm init --tiller-namespace dev-team

# creates service account, role, and rolebinding needed and attaches service account to tiller
$ ./secure-tiller examples/dev-team-rbac-profile/ --namespace dev-team

```

You should now have a Tiller in the dev-team namespace that is only allowed to do the things specified in examples/dev-team-rbac-profile/role/role-tiller.yaml
