FROM alpine:3.17

ARG TARGETOS
ARG TARGETARCH

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>" \
  org.label-schema.name="CodeGPT" \
  org.label-schema.vendor="Bo-Yi Wu" \
  org.label-schema.schema-version="1.0"

LABEL org.opencontainers.image.source=https://github.com/appleboy/CodeGPT
LABEL org.opencontainers.image.description="A CLI written in Go language that writes git commit messages or do a code review brief for you using ChatGPT AI."
LABEL org.opencontainers.image.licenses=MIT

RUN apk add --no-cache ca-certificates=20220614-r4 git=2.38.5-r0 && \
  rm -rf /var/cache/apk/*

COPY release/${TARGETOS}/${TARGETARCH}/codegpt /bin/

CMD ["/bin/codegpt"]
