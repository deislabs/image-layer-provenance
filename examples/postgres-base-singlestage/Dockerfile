##############################################################################################
# Primary Stage – Primary Base Image: postgres
##############################################################################################
FROM docker.io/library/postgres:14-bullseye

ADD https://github.com/kubernetes/client-go/archive/master.tar.gz /kubernetes/client-go.tar.gz
RUN echo "hello" > /hello.txt

EXPOSE 5432/tcp
EXPOSE 5432/udp
ENV POSTGRES_ENV_1="postgres_env_val_1"

ADD https://github.com/kubernetes/kubectl/archive/master.tar.gz /kubernetes/kubectl.tar.gz
RUN echo "breakfast" > /breakfast.txt

EXPOSE 5433/tcp
EXPOSE 5433/udp
ENV POSTGRES_ENV_2="postgres_env_val_2"

ADD https://github.com/kubernetes/dashboard/archive/master.tar.gz /kubernetes/dashboard.tar.gz
RUN echo "goodbye" > /goodbye.txt

EXPOSE 5434/tcp
EXPOSE 5434/udp
ENV POSTGRES_ENV_3="postgres_env_val_3"
