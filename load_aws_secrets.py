# Use this code snippet in your app.
# If you need more information about configurations
# or implementing the sample code, visit the AWS docs:
# https://aws.amazon.com/developer/language/python/

import boto3
import json
import os
from botocore.exceptions import ClientError

def get_secret():
    secret_name = "simple_bank"
    region_name = "sa-east-1"

    # Create a Secrets Manager client
    session = boto3.session.Session()
    client = session.client(
        service_name='secretsmanager',
        region_name=region_name
    )

    try:
        get_secret_value_response = client.get_secret_value(
            SecretId=secret_name
        )
    except ClientError as e:
        # For a list of exceptions thrown, see
        # https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
        raise e

    # Decrypts secret using the associated KMS key.
    return get_secret_value_response['SecretString']

def write_secrets(filename, secret):
    # Parse the secret string as JSON
    secret_dict = json.loads(secret)

    # Create the app.env file and write the key-value pairs
    with open(filename, "w") as file:
        for key, value in secret_dict.items():
            file.write(f"{key}={value}\n")

def load_env_variables(file_path):
    with open(file_path, "r") as file:
        for line in file:
            line = line.strip()
            if line and not line.startswith("#"):
                key, value = line.split("=")
                os.environ[key] = value

write_secrets("app.env", get_secret())
load_env_variables("app.env")