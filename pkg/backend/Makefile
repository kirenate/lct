PYTHON ?= .venv/bin/python

install:
	python3.12 -m venv .venv
	$(PYTHON) -m pip install poetry
	$(PYTHON) -m poetry update

run:
	$(PYTHON) web_entrypoint.py
