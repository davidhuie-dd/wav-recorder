.PHONY: server-image run-server

models/deepspeech/deepspeech-0.6.1-models:
	mkdir -p models/deepspeech
	cd models/deepspeech && \
		wget https://github.com/mozilla/DeepSpeech/releases/download/v0.6.1/deepspeech-0.6.1-models.tar.gz && \
		tar xvfz deepspeech-0.6.1-models.tar.gz && \
		rm deepspeech-0.6.1-models.tar.gz

server-image:
	docker build -t  mic2text .

run-server: models/deepspeech/deepspeech-0.6.1-models
	docker run --rm -p 3000:3000 -v `pwd`/models/deepspeech/deepspeech-0.6.1-models:/code  mic2text
