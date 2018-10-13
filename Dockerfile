# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo:v0.12.6 as builder

RUN mkdir -p $GOPATH/src/buddhabowls
WORKDIR $GOPATH/src/buddhabowls/

# this will cache the npm install step, unless package.json changes
# ADD package.json .
# ADD yarn.lock .
# RUN yarn install --no-progress
ADD . .
RUN npm install
ENV GO_ENV=development
RUN buffalo build --static -o bin/app

# CMD ./bin/app
# RUN go get $(go list ./... | grep -v /vendor/)

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/
COPY --from=builder /go/src/buddhabowls/bin/app .

EXPOSE 3000
CMD /bin/app
# # Comment out to run the migrations before running the binary:
# CMD /bin/app migrate; sh /bin/app
# CMD sh /bin/app
