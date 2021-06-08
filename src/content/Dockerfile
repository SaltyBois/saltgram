FROM dependencies

ENV GO111MODULE=on
ENV GOFLAGS=-mod=mod
ENV APP_HOME=/go/src/saltgram/content

WORKDIR $APP_HOME
COPY . .
#RUN go get -d -v ./...
RUN go install -v ./...
EXPOSE 443 80
CMD ["content"]
