# kubectl grep

A kubectl plugin performs a grep lookup to show the state of the pod when managing multiple workloads and namespaces

## Usage

```
go build -o kubectl-grep  ./...

mv kubectl-grep /usr/local/bin

kubectl grep nginx
```

