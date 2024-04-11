# syntax=docker/dockerfile:1

FROM alpine:3.19 as downloader
ADD https://github.com/blackstork-io/fabric/releases/latest/download/fabric_linux_x86_64.tar.gz /
RUN tar -xf fabric_linux_x86_64.tar.gz


FROM gcr.io/distroless/static-debian12:latest
COPY --from=downloader /fabric /fabric
ENTRYPOINT [ "/fabric" ]
CMD [ "--help" ]