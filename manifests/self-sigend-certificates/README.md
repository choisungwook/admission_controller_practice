# 개요
* admission controller을 사용하기 위한 ```self-signed 인증서``` 생성

# 생성 방법

```sh
rm -rf ./certs
mkdir certs

openssl req -x509 -newkey rsa:4096 -nodes -out certs/ca.crt -keyout certs/ca.key -days 365 -config ./ext.cnf -extensions req_ext

kubectl create secret tls webhook-certs --cert=certs/ca.crt --key=certs/ca.key --namespace=default

```


```sh
CA_BUNDLE=$(cat ./certs/ca.crt | base64 | tr -d '\n')
sed -e 's@${CA_BUNDLE}@'"$CA_BUNDLE"'@g' < ../webhook.yaml | kubectl apply -f -
```



```sh
# Clean up existing certificates
rm -rf ./certs
mkdir certs

# Generate a new self-signed CA certificate
openssl req -nodes -new -x509 -keyout certs/ca.key -out certs/ca.crt -subj "/CN=Admission Controller CA"

# Generate a new private key for the server
openssl genrsa -out certs/tls.key 2048

# Generate the certificate signing request (CSR) with SANs
openssl req -new -key certs/tls.key -out certs/tls.csr -subj "/CN=admission-server.default.svc" -config san.cnf

# Sign the CSR with the CA certificate to generate the server certificate
openssl x509 -req -in certs/tls.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial \
    -out certs/tls.crt -days 365 -extensions v3_req -extfile san.cnf
```

```sh

openssl req -x509 -sha256 -newkey rsa:2048 -keyout .key -out webhook.crt -days 1024 -nodes -addext "subjectAltName = DNS.1:validate.default.svc"
```
