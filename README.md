# Go Swag CLI Setup

## 1. Add Go bin directory to PATH

To make Go binaries globally accessible, run:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Swagger docs
```
swag --version
swag.exe version v1.16.4
```
- Reads code comments Pattern to generate docs
swag init

# Makefile

```sh
make build   # compile binary
make run     # run locally
make test    # run unit tests
make docker  # build docker image
make clean   # clean build artifacts
```

- For passing build params
```
go build -ldflags "-X main.AppVersion=1.0.0 -X main.CommitHash=$(git rev-parse --short HEAD) -X 'main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" -o myapp main.go

```

# Build image
docker build -t go-service:latest .

# Run
```sh
docker run -d -p 80:8080 \
  go-service:latest
```
OR
```sh
docker-compose up --build
```
# Dockerhub
```sh
docker tag go-service:latest kunaljainwindock/go-service:1.0.0
docker push kunaljainwindock/go-service:1.0.0
****
bff7f7a9d443: Pushed
1.0: digest: sha256:c287abf22302c5f8985bb7ec0dfc23d170b1974500a6c67cdb165933d558f539 size: 3651
1.1: digest: sha256:0a590f01f64adb18a046dae5ea42b02617b69cbf436fef219bb4c55684f79fbd size: 3651
```

# Kubernetes Deployment

```
└── k8s/          <-- Kubernetes YAML files are here
    ├── namespace.yml
    ├── deployment.yml
    └── service.yml 
```
## 1. Apply Kubernetes Manifests
Apply all manifests in the `k8s/` directory:
```bash
kubectl apply -f k8s/
```

List all pods:
```
kubectl get pods -o wide
```

List services:
```
kubectl get svc
```

Check pod resource usage:
```
kubectl top pod
kubectl top pod -w   # watch continuously

```

Watch events in the cluster:
```
kubectl get events -w
```
# Kubernetes Load Balancing Test Report

## 1. Setup

- **Deployment**: `my-go-service`  
- **Replicas**: 2 Pods
- **Pods**:

| Pod Name                                | Pod IP       |
|----------------------------------------|-------------|
| my-go-service-5c789c4dd6-qcb7v         | 10.244.0.6 |
| my-go-service-5c789c4dd6-rc6ch         | 10.244.0.5 |

- **Service**: `my-go-service`  
  - Exposes port `8080`
  - Load balances traffic between the two pods

---

## 2. Test Script

Windows batch file used to send repeated requests:

```bat
@echo off
REM Loop from 1 to 50
for /l %%i in (1,1,50) do (
    echo Iteration %%i
    curl http://localhost/api/v1/readyz
)
pause
```
## 3. Observations
- Each response includes the pod hostname
```
Iteration 48
{"hostname":"my-go-service-5c789c4dd6-rc6ch"}
Iteration 49
{"hostname":"my-go-service-5c789c4dd6-rc6ch"}
Iteration 50
{"hostname":"my-go-service-5c789c4dd6-qcb7v"}
...
```

## 4. Conclusion

- The Deployment and Service are correctly configured.

- Kubernetes Service is successfully distributing traffic across multiple pods.

- The batch script demonstrates that the load is balanced.
