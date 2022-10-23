# kubectl grep

A kubectl plugin performs basic string lookup. If the pod is present, the outputs show the state of the pod and event errors. This tool comes in handy when managing multiple workloads in different namespaces.

Example:

![Screenshot](/assests/img/screenshot.png)

## Usage

```
go build -o kubectl-grep  ./...

mv kubectl-grep /usr/local/bin

kubectl grep nginx
```

## Credits:
Inspired from [kubectl-tree](https://github.com/ahmetb/kubectl-tree) 
