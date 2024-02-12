<!-- TOC -->

- [개요](#개요)
- [준비](#준비)
- [실행 방법](#실행-방법)
  - [self signed 인증서 생성](#self-signed-인증서-생성)
  - [admission controller를 실행할 golang pod 생성](#admission-controller를-실행할-golang-pod-생성)
  - [webhook 생성](#webhook-생성)
- [admission controller 테스트](#admission-controller-테스트)
- [참고자료](#참고자료)

<!-- /TOC -->
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

```sh
mkdir certs

openssl req -x509 -newkey rsa:4096 -nodes -out certs/ca.crt -keyout certs/ca.key -days 365 -config ./cert.cnf -extensions req_ext

kubectl create secret tls webhook-certs --cert=certs/ca.crt --key=certs/ca.key --namespace=default
```

## admission controller를 실행할 golang pod 생성

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
CA_BUNDLE=$(cat ./certs/ca.crt | base64 | tr -d '\n')
sed -e 's@${CA_BUNDLE}@'"$CA_BUNDLE"'@g' < ./manifests/webhook.yaml | kubectl apply -f -
```

# admission controller 테스트

* buysbox 생성 후 admission controller pod로그 확인

```sh
kubectl apply -f ./manifests/busybox-pod.yaml
```

# 참고자료
* https://github.com/GoesToEleven/golang-web-dev
* https://koras02.tistory.com/236
* https://dev.to/ashokan/simple-validation-webhook-with-python-5bof
* https://github.com/slackhq/simple-kubernetes-webhook/blob/main/pkg/admission/admission.go#L83
* https://sphong0417.tistory.com/54
