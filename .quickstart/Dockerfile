FROM docker:git
COPY .  /animal-store
RUN apk add curl
RUN apk add --no-cache make
COPY --from=golang:alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"
WORKDIR /animal-store/.quickstart