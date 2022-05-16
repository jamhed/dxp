FROM node:16-alpine3.14 as frontend
COPY ui ui
RUN cd ui && npm i
RUN cd ui && npm run build

FROM golang:1.18.2-alpine3.14 as backend
COPY . /home/dxp
RUN apk add git && \
	cd /home/dxp && \
	export VERSION=$(git describe --tags) && \
	export COMMIT=$(git rev-parse --short HEAD) && \
	export BUILDDATE=$(date +%Y%m%d%H%M%S) && \
	go build -ldflags="-X main.version=$VERSION -X main.builddate=$BUILDDATE -X main.commit=$COMMIT"

FROM alpine:3.14
RUN apk update
COPY --from=backend /home/dxp/dxp dxp
COPY --from=frontend ui/dist ui/dist
CMD ["./dxp"]
