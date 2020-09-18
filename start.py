""" Custom RabbitMQ spammer """

import time
from pathlib import Path
from typing import Optional

import rabbitpy


def load_data(filename: str) -> str:
    return Path(filename).read_text()


def publish(channel: rabbitpy.Channel, queue_name: str, data: str):
    message = rabbitpy.Message(channel, data)
    message.publish("", queue_name, mandatory=True)


def main(filename: str, ampq_url: Optional[str], queue_name: str, num_messages: int):
    data = load_data(filename)

    with rabbitpy.Connection(ampq_url) as conn:
        print("Connected to RabbitMQ")

        with conn.channel() as channel:
            channel.enable_publisher_confirms()
            queue = rabbitpy.Queue(channel, queue_name)
            queue.durable = True
            queue.declare()

            print("Begin spamming...")
            start = time.time()
            for _ in range(num_messages):
                publish(channel, queue_name, data)

            duration = time.time() - start
            velocity = num_messages / duration
            print(f"Duration: {round(duration, 2)} secs | Velocity: {round(velocity, 2)} msgs/sec")


if __name__ == "__main__":
    from argparse import ArgumentParser

    parser = ArgumentParser()
    parser.add_argument("-f", "--filename", type=str, help="Filename of message to spam")
    parser.add_argument("-u", "--ampq_url", help="AMQP location")
    parser.add_argument("-q", "--queue_name", type=str, help="Queue name to post")
    parser.add_argument("-n", "--num_messages", type=int, help="Number of messages needed to spam Rabbit")
    kwargs = vars(parser.parse_args())
    main(**kwargs)
