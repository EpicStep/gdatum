ARG IMG_TAG=nonroot
FROM gcr.io/distroless/static-debian12:${IMG_TAG}

COPY --chown=nonroot:nonroot /gdatum ./gdatum

EXPOSE 8080 8081

USER 65532:65532

ENTRYPOINT ["./gdatum"]
