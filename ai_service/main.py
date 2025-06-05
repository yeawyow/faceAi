import aio_pika
import asyncio
import json
from insightface.app import FaceAnalysis
from PIL import Image
import numpy as np

app = FaceAnalysis(name="buffalo_l")
app.prepare(ctx_id=0, det_size=(640, 640))

async def on_message(message: aio_pika.IncomingMessage):
    async with message.process():
        payload = json.loads(message.body.decode())
        images = payload.get("images", [])

        for img in images:
            image_id = img.get("image_id")
            image_path = img.get("image_name")

            print(f"Processing image_id: {image_id} path: {image_path}")

            image = Image.open(image_path).convert("RGB")
            image_np = np.array(image)
            faces = app.get(image_np)
            print(f"Detected {len(faces)} faces")

            # ประมวลผลต่อ เช่น สร้าง embedding, บันทึกลง DB หรือส่งต่อ API

async def main():
    connection = await aio_pika.connect_robust("amqp://skko:skkospiderman@rabbitmq:5672/")
    channel = await connection.channel()
    queue = await channel.declare_queue("face_images_queue", durable=True)

    await queue.consume(on_message)

    print("AI service is listening on face_images_queue")
    await asyncio.Future()  # keep running

if __name__ == "__main__":
    asyncio.run(main())
