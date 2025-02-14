FROM python:3.10.15-slim

RUN apt update && apt upgrade -y

WORKDIR /python-api-onnx

COPY ./python-api-onnx/ /python-api-onnx/
COPY ./models/logistic_regression_iris.onnx /models/

RUN pip install -r requirements.txt -c constraint.txt --upgrade pip

ENTRYPOINT [ "python3" ]
CMD [ "main.py" ]
