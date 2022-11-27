FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY xyz2ogc /bin/xyz2ogc
ENTRYPOINT ["/bin/xyz2ogc"]
