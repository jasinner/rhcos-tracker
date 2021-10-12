FROM registry.fedoraproject.org/fedora:34 AS builder
RUN dnf install -y golang-bin
WORKDIR /opt/
COPY . .
RUN go build -o bin/releases ./cmd/releases

FROM registry.fedoraproject.org/fedora:34
COPY wait-for-postgres.sh /app/
COPY krb5.conf /app/
#RUN dnf install -y ca-certificates golang-bin yarnpkg maven rubygem-bundler postgresql npm rpm-build tini && dnf clean all
RUN curl -o /etc/pki/tls/certs/RH-IT-Root-CA.crt https://password.corp.redhat.com/RH-IT-Root-CA.crt && update-ca-trust
COPY --from=builder /opt/bin/releases /app/
WORKDIR /app
ENTRYPOINT ["./releases"]