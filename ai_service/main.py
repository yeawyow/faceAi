import aio_pika
import asyncio

async def on_message(message: aio_pika.IncomingMessage):
    async with message.process():
        print("Message received:", message.body)

async def main():
    connection = await aio_pika.connect_robust("amqp://skko:skkospiderman@rabbitmq/")
    channel = await connection.channel()
    queue = await channel.declare_queue("face_jobs", durable=True)
    await queue.consume(on_message)
    print("Waiting for messages...")
    await asyncio.Future()  # รอไม่สิ้นสุด

asyncio.run(main())
