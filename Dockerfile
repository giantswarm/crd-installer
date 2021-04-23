FROM scratch

COPY ./scripts /scripts
COPY ./crd-installer /usr/local/bin/crd-installer

ENTRYPOINT ["/usr/local/bin/crd-installer"]
