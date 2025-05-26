FROM golang:1.18 as build
WORKDIR /ics-proxy
COPY go.mod go.sum ./
COPY cmd/ cmd/
COPY internal/ internal/
RUN go build -tags lambda.norpc cmd/ics-proxy-lambda/ics-proxy-lambda.go
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /ics-proxy/ics-proxy-lambda ./ics-proxy-lambda
COPY res/ res/
ENTRYPOINT [ "./ics-proxy-lambda" ]
