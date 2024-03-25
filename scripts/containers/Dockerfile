FROM golang:1.20-bullseye AS build

RUN apt update && apt install build-essential -y

WORKDIR /terrad
COPY . /terrad

RUN make build

FROM golang:1.20-bullseye AS run
COPY ./scripts/containers/build-wrapper.sh /terrad/build-wrapper.sh
RUN chmod +x /terrad/build-wrapper.sh
RUN wget https://github.com/CosmWasm/wasmvm/releases/download/v1.5.2/libwasmvm.x86_64.so -P /usr/lib/

VOLUME /terrad
COPY --from=build /terrad/ /terrad/
WORKDIR /terrad

EXPOSE 26656 26657
ENTRYPOINT ["/terrad/build-wrapper.sh"]
CMD ["start", "--log_format", "plain"]
STOPSIGNAL SIGTERM
