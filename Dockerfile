FROM ubuntu:19.10

RUN apt update -y && apt install -y git sox libc++-dev libpthread-stubs0-dev \
    python3-pip software-properties-common

RUN pip3 install --upgrade pip
RUN pip3 install deepspeech-gpu

RUN git clone -b v0.6.1 --depth=1 https://github.com/mozilla/DeepSpeech.git /DeepSpeech
RUN cd /DeepSpeech && pip3 install -r requirements.txt && pip3 install $(python3 util/taskcluster.py --decoder) \
    && python3 util/taskcluster.py --target .

RUN add-apt-repository -y ppa:longsleep/golang-backports && apt update -y && \
    apt install -y golang-1.13

ENV CGO_LDFLAGS "-L/DeepSpeech/"
ENV CGO_CXXFLAGS "-I/DeepSpeech/native_client/"
ENV LD_LIBRARY_PATH "/DeepSpeech/:$LD_LIBRARY_PATH"

COPY . /code
WORKDIR /code

RUN /usr/lib/go-1.13/bin/go build -o /usr/local/bin/mic2text-server ./cmd/mic2text-server/

ENTRYPOINT ["/usr/local/bin/mic2text-server"]
CMD [""]