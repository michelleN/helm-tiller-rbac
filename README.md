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
serviceaccount "dev-team-rbac-profile" created
rolebinding "dev-team-tiller-binding" created
role "dev-team" created

Congrats! Your Tiller in the dev-team namespace is secured with the dev-team-rbac-profile service account
You can verify your Tiller config by running this command
	$ kubectl -n dev-team get deployment tiller-deploy -o json

```

You should now have a Tiller in the `dev-team` namespace that is only allowed to do the things specified in `examples/dev-team-rbac-profile/role/role-tiller.yaml`
