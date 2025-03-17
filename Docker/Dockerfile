FROM ubuntu:22.04

ENV TZ=Europe/Warsaw

RUN apt-get update && apt-get install -y \
    python3.10 \
    openjdk-8-jdk \
    curl \
    zip \
    unzip \
    && ln -sf /usr/bin/python3.10 /usr/bin/python3 \
    && rm -rf /var/lib/apt/lists/*



RUN curl -s "https://get.sdkman.io" | bash && \
    bash -c "source $HOME/.sdkman/bin/sdkman-init.sh && sdk install kotlin && sdk install gradle"

ENV SDKMAN_DIR="/root/.sdkman"
ENV PATH="${SDKMAN_DIR}/bin:${SDKMAN_DIR}/candidates/gradle/current/bin:$PATH"
    
WORKDIR /app

COPY . .

RUN gradle dependencies --no-daemon

RUN gradle build --no-daemon


CMD ["gradle", "run", "--no-daemon"]



