FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-panther"]
COPY baton-panther /