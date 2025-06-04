from insightface.app import FaceAnalysis
import aio_pika
import asyncio
import json
import numpy as np
from PIL import Image

# โหลดโมเดล
app = FaceAnalysis()
app.prepare(ctx_id=0, det_size=(640, 640))

async def on_message(message: aio_pika.IncomingMessage):
    async with message.process():
        try:
            payload = json.loads(message.body.decode())
            image_id = payload["imageId"]
            path = payload["path"]
            print(f"📥 รับงาน: imageId={image_id}, path={path}")

            # อ่านภาพจาก path โดยตรง
            image = Image.open(path).convert("RGB")
            faces = app.get(np.array(image))

            print(f"🧠 ค้นหาใบหน้า imageId={image_id} เจอทั้งหมด {len(faces)} ใบหน้า")

            # คุณสามารถเก็บ embedding หรือตอบกลับต่อได้ที่นี่

        except Exception as e:
            print(f"❌ เกิดข้อผิดพลาดกับ imageId={payload.get('imageId')}: {str(e)}")

async def main():
    try:
        connection = await aio_pika.connect_robust("amqp://skko:skkospiderman@rabbitmq/")
        channel = await connection.channel()
        queue = await channel.declare_queue("face_job", durable=True)
        await queue.consume(on_message)
        print("✅ AI Service พร้อมทำงานแล้ว...")
        await asyncio.Future()  # รอไม่สิ้นสุด
    except Exception as e:
        print("❌ เชื่อมต่อ RabbitMQ ไม่สำเร็จ:", str(e))

asyncio.run(main())
