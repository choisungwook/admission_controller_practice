# 개요
kubernetes adminssion controller 연습

# 준비
* kind 쿠버네티스 클러스터 생성

```sh
# 생성
make up

# 삭제
make down
```

# 실행 방법

## self signed 인증서 생성
* [self-signed-certificated 문서 이동](./manifests/self-sigend-certificates/)

# admission controller 생성

* golang 컨테이너가 있는 pod생성

```sh
kubecatl apply -f ./manifests/golang-deployment.yaml
kubecatl apply -f ./manifests/service.yaml
```

* golang pod 쉘 접속

```sh
kubectl get pod -n default
kubectl exec -it {pod 이름} -- /bin/bash
```

* go 모듈 초기화와 외부모듈 다운로드

```sh
go mod init github.com/opsarena/admission-controller
go get k8s.io/api/admission/v1
go get k8s.io/apimachinery/pkg/apis/meta/v1
go get k8s.io/api/core/v1
```

* main.go 파일 생성: [main.go](./src/main.go)파일을 pod안에 복사

* main.go 실행

```sh
go run main.go
```

## webhook 생성

```sh
CA_BUNDLE=$(cat ./manifests/self-sigend-certificates/certs/ca.crt | base64 | tr -d '\n')
sed -e 's@${CA_BUNDLE}@'"$CA_BUNDLE"'@g' < ./manifests/webhook.yaml | kubectl apply -f -
```

# pod 생성 테스트

* buysbox 생성 후 pod로그 확인

```sh
kubectl apply -f ./manifests/busybox-pod.yaml
```

# 참고자료
* https://github.com/GoesToEleven/golang-web-dev
* https://koras02.tistory.com/236
* https://dev.to/ashokan/simple-validation-webhook-with-python-5bof
* https://github.com/slackhq/simple-kubernetes-webhook/blob/main/pkg/admission/admission.go#L83
