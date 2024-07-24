import json
from google.cloud import storage

CRED = """{
  ...
}"""


def upload_blob(bucket_name, source_file_name, destination_blob_name):
    """Uploads a file to the bucket."""
    storage_client = storage.Client.from_service_account_info(json.loads(CRED))
    bucket = storage_client.bucket(bucket_name)
    blob = bucket.blob(destination_blob_name)
    print(blob)
    blob.upload_from_filename(source_file_name)


if __name__ == "__main__":
    upload_blob("sd-uploads", "test.txt", "uploads/test3.txt")
