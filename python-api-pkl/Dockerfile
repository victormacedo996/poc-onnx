FROM python:3.10.15-slim

RUN apt update && apt upgrade -y

WORKDIR /python-api-pkl

COPY ./python-api-pkl/ /python-api-pkl/
COPY ./models/logistic_regression_iris.pkl /models/

RUN pip install -r requirements.txt -c constraint.txt --upgrade pip

ENTRYPOINT [ "python3" ]
CMD [ "main.py" ]
