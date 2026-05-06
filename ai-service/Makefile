.PHONY: dev run

test:
	python3 -m unittest discover -s tests

dev:
	python3 -m uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload

run:
	python3 -m uvicorn app.main:app --host 0.0.0.0 --port 8000
