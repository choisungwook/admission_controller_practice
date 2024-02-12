# 개요
* admission controller을 사용하기 위한 ```self-signed 인증서``` 생성

# 생성 방법

```sh
rm -rf ./certs
mkdir certs

openssl req -x509 -newkey rsa:4096 -nodes -out certs/ca.crt -keyout certs/ca.key -days 365 -config ./ext.cnf -extensions req_ext
kubectl create secret tls webhook-certs --cert=certs/ca.crt --key=certs/ca.key --namespace=default
```
