FROM registry.fedoraproject.org/fedora:34 AS builder
RUN dnf install -y golang-bin
WORKDIR /opt/
COPY . .
RUN go build -o bin/releases ./cmd/releases

FROM registry.fedoraproject.org/fedora:34
#RUN dnf install -y ca-certificates golang-bin yarnpkg maven rubygem-bundler postgresql npm rpm-build tini && dnf clean all
#TODO get oc client
COPY --from=builder /opt/bin/releases /app/
WORKDIR /app
ENTRYPOINT ["./releases"]