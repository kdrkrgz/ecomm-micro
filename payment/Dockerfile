FROM golang:1.21 AS builder

ARG SSH_PRIVATE_KEY
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/

# for private package => docker build --build-arg SSH_PRIVATE_KEY="$(cat ~/.ssh/gh_rsa)" -t test .
RUN mkdir -p /root/.ssh && \
chmod 0700 /root/.ssh && \
ssh-keyscan github.com > /root/.ssh/known_hosts

WORKDIR /usr/src/app

COPY . .

RUN echo "$SSH_PRIVATE_KEY" > /root/.ssh/id_rsa && \
chmod 600 /root/.ssh/id_rsa && \
go mod download && \
rm -rf /root/.ssh/

RUN CGO_ENABLED=0 GOOS=linux go build -o payment ./cmd/main.go

FROM scratch
COPY --from=builder /usr/src/app/payment ./payment
CMD ["./payment"]