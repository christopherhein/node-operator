# Node Operator

This operator runs as a DaemonSet and configures the nodes
`/.ssh/authorized_keys` file so that as employees get added all that needs to
happen to give them access to the cluster is to create a new `AuthorizedKey`
manifest. 

> This is only the first thing implemented in the node operator, expect that 
> this functionality could grow.

## Get Started

First thing you need to do is apply the operator.

```bash
kubectl apply -f https://raw.githubusercontent.com/christopherhein/node-operator/master/node-operator.yaml
```

To test this create a file like this

```yaml
# authorized_key.yml
apiVersion: node.io/v1
kind: AuthorizedKeys
metadata:
  name: christopherhein-key
data:
  key: ssh-rsa AAAAB3N... me@christopherhein.com
```

Then install like you would any other manifest file.

```bash
kubectl apply -f authorized_key.yml
```

Now you can view the installed keys.

```bash
$ kubectl get authorizedkeys
NAME                  KIND
christopherhein-key   AuthorizedKey.v1.node.io
```

It also responds to `authkey`, `authorizedkeys`, `authorized-keys`,
`authorizedkey` and `authorized-key`

## Development

Make sure you first have all the dependencies installed.

```bash
brew install dep
brew upgrade dep
dep ensure
```

After you make any changes to the specs under `pkg/api/*/v1alpha1/*.go` make
sure to rerun the codegen script.

```bash
./codegen.sh
```

After have made sure the codegen is run you can build the container, push to a
registry and update or deploy the operator.

```bash
export VERSION=0.0.4
CGO_ENABLED=0 GOOS=linux go build
docker build -t christopherhein/node-operator:$VERSION .
docker push christopherhein/node-operator:$VERSION
kubectl apply -f node-operator.yaml
```

