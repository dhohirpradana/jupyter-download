FROM alpine:latest AS builder

RUN apk --no-cache add curl zip

RUN curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x kubectl \
    && mv kubectl /usr/local/bin/



FROM nginx

WORKDIR /app

COPY --from=builder /usr/local/bin/kubectl /usr/local/bin/kubectl

COPY jupyter-download /app

RUN apt update -y && apt install zip -y

EXPOSE 9090

CMD ./jupyter-download
