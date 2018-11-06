# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo:v0.12.6 as builder

ENV GO_ENV=development

RUN mkdir -p $GOPATH/src/buddhabowls
WORKDIR $GOPATH/src/buddhabowls/

# this will cache the npm install step, unless package.json changes
ADD ./package.json .
RUN npm install
ADD . .
RUN sed -i "s|host: .*$|host: postgres|g" database.yml
RUN buffalo build --static -o bin/app

# We need to use an older version of gobuffalo here for the migrations to succeed
FROM gobuffalo/buffalo:v0.11.0 as migrate

WORKDIR $GOPATH/src/buddhabowls/
ADD . .
ENV GO_ENV=development

# Run the compiled binary in a slim alpine image
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
