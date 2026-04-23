FROM alpine:3.23

RUN apk add --no-cache bash curl gcompat git less openssh-client-default ripgrep tzdata make
RUN ln -sf /usr/share/zoneinfo/hostlocaltime /etc/localtime

ARG USERNAME=tinfoil-proxy
RUN addgroup --gid 1000 $USERNAME && adduser --uid 1000 --ingroup $USERNAME --disabled-password $USERNAME
RUN mkdir -p /home/$USERNAME/persistent-config && chown $USERNAME:$USERNAME /home/$USERNAME /home/$USERNAME/persistent-config

USER $USERNAME
WORKDIR /home/$USERNAME

RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/home/$USERNAME/persistent-config/.bash_history" \
  && echo "$SNIPPET" >> "/home/$USERNAME/.bashrc"

RUN curl -L -o /tmp/go.tar.gz https://go.dev/dl/go1.26.2.linux-amd64.tar.gz \
  && mkdir -p "/home/$USERNAME/local" && tar -C "/home/$USERNAME/local" -xzf /tmp/go.tar.gz \
  && rm /tmp/go.tar.gz

ENV PATH="$PATH:/home/$USERNAME/local/go/bin:/home/$USERNAME/go/bin"

WORKDIR /workspaces/tinfoil-proxy
