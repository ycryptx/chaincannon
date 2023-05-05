FROM ignitehq/cli:0.23.0

RUN ignite scaffold chain example

COPY config.yml example/

RUN cd example; \
    ignite chain init;

ENTRYPOINT exampled start