#!/bin/python3

import json

from command_runner import command_runner

from google.api_core import retry
from google.cloud import pubsub_v1, storage


#Log to console, it will be captured by systemd
import logging
from sys import stdout
logFormatter = logging.Formatter("%(asctime)s [%(levelname)-5.5s]  %(message)s")
consoleHandler = logging.StreamHandler(stdout)
consoleHandler.setFormatter(logFormatter)
log = logging.getLogger("worker_manager")
log.addHandler(consoleHandler)
log.setLevel(level=logging.INFO)


TIMEOUT = 30 #Run for 2 minutes maximum. Currently 30s, because that should be enough for 2-3 iterations of the flag.

project_id = "dicectf-2022"

subscription_id = "cache_challenge_entry-sub"
subscriber = pubsub_v1.SubscriberClient()
subscription_path = subscriber.subscription_path(project_id, subscription_id)

result_bucket_name = "cache-on-the-side"
storage_client = storage.Client()
result_bucket = storage_client.bucket(result_bucket_name)
MAXFILESIZE = 1024*1024*10

def get_next_attempt():
    # Returns a single message, blocks indefinitely until message is received.

    while True: #Breaks early with a return
        response = subscriber.pull(
            request={"subscription": subscription_path, "max_messages": 1},
            retry=retry.Retry(deadline=300),
        )
        assert len(response.received_messages) <= 1 #Die if otherwise, shouldn't be possible tho!

        if len(response.received_messages) == 1:
            message = response.received_messages[0]
            entry_id = json.loads(message.message.data)["entry_id"]
            log.info(f"\n****\nReceived entry: {entry_id}. Attributes: ")
            for key in message.message.attributes:
                value = message.message.attributes.get(key)
                log.info(f"\t{key}: {value}")
            return message


def upload_result_blob(result_blob_name, contents):
    """Uploads a file to the bucket."""
    
    if len(contents) > MAXFILESIZE:
        contents = contents[:MAXFILESIZE]

    blob = result_bucket.blob(result_blob_name)
    blob.upload_from_string(contents)
    metadata = {
        "Cache-Control":"public, max-age=3600",
        "Content-Disposition":"inline",
        "Content-Type":"application/octet-stream"
    }
    blob.metadata=metadata
    blob.patch()


def return_results(message, entry_id, result:str):
    """Return results by uploading to google cloud storage, then acknowledge the entry from the queue """
    
    log.info(f"Entry {entry_id} | Uploading results. ")
    upload_result_blob(entry_id, result)
    log.info(f"Entry {entry_id} | Result Uploaded")

    #Finally acknowledge that script has run and is done. Removes request from the queue.
    ack_ids = [message.ack_id]
    subscriber.acknowledge(
        request={"subscription": subscription_path, "ack_ids": ack_ids}
    )

    log.info(f"Entry {entry_id} | Acknowledged original request.")


def run_attempt(code_to_test, entry_id):
    try:
        with open("attack/attack.c","w") as f:
            f.write(code_to_test)

        exit_code_build, output_build = command_runner(["docker", "build", "-t", "attack", "."], timeout = 10)
        if exit_code_build != 0:
            log.error(f"Entry {entry_id} | docker build: {exit_code_build}\n{output_build}")
            return "There was an error preparing your code. Please contact the challenge author."
        exit_code_run, result = command_runner(["docker", "run", "-t", "--network", "none", "--cpuset-cpus=1", "--name", "attack_container", "attack"], timeout = TIMEOUT)
        if exit_code_run == -254: #Timeout: Strip " ... for command {} execution. Original output was: " from log.
            location = result.find("Original output was:")
            result = "Timeout expired.\n" + result[location+21 :]

    except Exception as err:
        log.critical(f"Entry {entry_id} | {err}")
        result = "There was an unknown error. Please contact the challenge author"
    finally:
        command_runner(["docker","rm","-f","attack_container"]) #Make absolutely sure the docker container has stopped.
    return result



# Wrap the subscriber in a 'with' block to automatically call close() to
# close the underlying gRPC channel when done.
with subscriber:
    log.info("Started worker_manager.py server")
    while True: #Run forever!
        message = get_next_attempt()
        j = json.loads(message.message.data)
        code_to_test = j["code_to_test"]
        entry_id = j["entry_id"]

        result = run_attempt(code_to_test, entry_id)

        return_results(message, entry_id, result)


