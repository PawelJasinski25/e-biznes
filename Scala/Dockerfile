FROM ubuntu:22.04

ENV TZ="Europe/Warsaw"

RUN apt-get update && apt-get install -y \
    curl \
    unzip \
    zip \
    software-properties-common \
    && rm -rf /var/lib/apt/lists/*

RUN curl -s "https://get.sdkman.io" | bash
RUN bash -c "source $HOME/.sdkman/bin/sdkman-init.sh && \
    sdk install java 11.0.26-tem && \
    sdk install sbt 1.10.9 && \
    sdk install scala 3.3.1"

ENV PATH=$PATH:/root/.sdkman/candidates/java/current/bin
ENV PATH=$PATH:/root/.sdkman/candidates/sbt/current/bin
ENV PATH=$PATH:/root/.sdkman/candidates/scala/current/bin

RUN mkdir -p /scala-app
WORKDIR /scala-app

COPY . .

EXPOSE 9000

CMD ["sbt", "run"]

