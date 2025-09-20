from dataclasses import dataclass

from shared.settings import app_settings
import boto3
from botocore.exceptions import ClientError
from loguru import logger
import os
import threading
import sys

@dataclass
class MinioRepository:
    def __post_init__(self) -> None:
        self.client = boto3.client(
            "s3",
            aws_access_key_id=app_settings.minio.aws_access_key_id,
            aws_secret_access_key=app_settings.minio.aws_secret_access_key,
        )

    async def health(self) -> None:
        if not await self.client.ping():
            raise Exception("non true ping")

    async def create_bucket(self, bucket_name: str) -> bool:
        try:
            self.client.create_bucket(Bucket=bucket_name)
        except ClientError as e:
            logger.error(e)
            return False
        return True

    async def upload_file(self, file_name: str, bucket: str):
        # If S3 object_name was not specified, use file_name
        try:
            response = self.client.upload_file(file_name, bucket)
        except ClientError as e:
            logger.error(e)
            return False
        return True

class ProgressPercentage(object):

    def __init__(self, filename):
        self._filename = filename
        self._size = float(os.path.getsize(filename))
        self._seen_so_far = 0
        self._lock = threading.Lock()

    def __call__(self, bytes_amount):
        # To simplify, assume this is hooked up to a single filename
        with self._lock:
            self._seen_so_far += bytes_amount
            percentage = (self._seen_so_far / self._size) * 100
            sys.stdout.write(
                "\r%s  %s / %s  (%.2f%%)" % (
                    self._filename, self._seen_so_far, self._size,
                    percentage))
            sys.stdout.flush()