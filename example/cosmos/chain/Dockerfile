FROM ignitehq/cli:0.23.0

RUN ignite scaffold chain example

COPY config.yml example/

RUN cd example; \
    ignite chain init; \
    exampled config chain-id example; 

EXPOSE 26656 26657 1317 9090

ENTRYPOINT exampled start