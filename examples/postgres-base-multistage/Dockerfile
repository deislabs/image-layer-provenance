##############################################################################################
# Multistage Build Stage (Builder 1) – Multistage Base Image: node
##############################################################################################
FROM docker.io/library/node:18-bullseye

RUN echo "builder1Foo" > /builder1Foo.txt

EXPOSE 8080/tcp
EXPOSE 8080/udp
ENV NODE_ENV_1="node_env_val_1"

RUN echo "builder1Bar" > /builder1Bar.txt

##############################################################################################
# Multistage Build Stage (Builder 2) – Multistage Base Image: python
##############################################################################################
FROM docker.io/library/python:3-bullseye as builder2

RUN echo "builder2Foo" > /builder2Foo.txt

EXPOSE 5000/tcp
EXPOSE 5000/udp
ENV PYTHON_ENV_1="python_env_val_1"

RUN echo "builder2Bar" > /builder2Bar.txt

##############################################################################################
# Primary Stage – Primary Base Image: postgres
##############################################################################################
FROM docker.io/library/postgres:14-bullseye

ADD https://github.com/kubernetes/client-go/archive/master.tar.gz /kubernetes/client-go.tar.gz
RUN echo "hello" > /hello.txt

EXPOSE 5432/tcp
EXPOSE 5432/udp
ENV POSTGRES_ENV_1="postgres_env_val_1"

# `COPY --from=0` is the same as `COPY --from=builder1`
COPY --from=0 /builder1Foo.txt /builder1Foo.txt
COPY --from=0 /builder1Bar.txt /builder1Bar.txt
ADD https://github.com/kubernetes/kubectl/archive/master.tar.gz /kubernetes/kubectl.tar.gz
RUN echo "breakfast" > /breakfast.txt

EXPOSE 5433/tcp
EXPOSE 5433/udp
ENV POSTGRES_ENV_2="postgres_env_val_2"

COPY --from=builder2 /builder2Foo.txt /builder2Foo.txt
COPY --from=builder2 /builder2Bar.txt /builder2Bar.txt
ADD https://github.com/kubernetes/dashboard/archive/master.tar.gz /kubernetes/dashboard.tar.gz
RUN echo "goodbye" > /goodbye.txt

EXPOSE 5434/tcp
EXPOSE 5434/udp
ENV POSTGRES_ENV_3="postgres_env_val_3"
