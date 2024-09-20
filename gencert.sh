#!/usr/bin/env bash

mkdir -p certs
pushd certs || exit 1

if [ ! -f ca.key ]; then
    openssl genrsa -out ca.key 4096
    openssl req -x509 -new -nodes -sha512 -days 3650 -subj "/C=CN/ST=Beijing/L=Beijing/O=example/OU=Personal/CN=CA" -key ca.key -out ca.crt
fi
#################################################################

openssl genrsa -out server.key 4096
openssl req -sha512 -new -subj "/C=CN/ST=Beijing/L=Beijing/O=example/OU=Personal/CN=server" -key server.key -out server.csr

cat > v3.ext <<-EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1=server
DNS.2=localhost
IP.1=127.0.0.1
IP.2=192.168.6.171
EOF

openssl x509 -req -sha512 -days 3650 -extfile v3.ext -CA ca.crt -CAkey ca.key -CAcreateserial -in server.csr -out server.crt

#################################################################
openssl genrsa -out client.key 4096
openssl req -sha512 -new -subj "/C=CN/ST=Beijing/L=Beijing/O=example/OU=Personal/CN=client" -key client.key -out client.csr

cat > v3_client.ext <<-EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = clientAuth
subjectAltName = @alt_names

[alt_names]
DNS.1=client
DNS.2=localhost
IP.1=127.0.0.1
IP.2=192.168.6.171
EOF

openssl x509 -req -sha512 -days 3650 -extfile v3_client.ext -CA ca.crt -CAkey ca.key -CAcreateserial -in client.csr -out client.crt

openssl x509 -in server.crt -noout -text |grep -A2 Alternative
openssl x509 -in client.crt -noout -text |grep -A2 Alternative
openssl verify -CAfile ca.crt server.crt
openssl verify -CAfile ca.crt client.crt
# certs/client.crt: OK
popd || exit 1