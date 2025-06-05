from insightface.app import FaceAnalysis
import aio_pika
import asyncio
import json
import numpy as np
from PIL import Image
import os

# โหลดโมเดล InsightFace
app = FaceAnalysis()
app.prepare(ctx_id=0, det_size=(640, 640))  # ctx_id=-1 ถ้าไม่มี GPU

async def on_message(message: aio_pika.IncomingMessage):
    async with message.process():
        try:
            payload = json.loads(message.body.decode())
            image_id = payload.get("imageId")
            image_path = payload.get("path")

            print(f"📥 รับงาน: imageId={image_id}, path={image_path}")

            # ตรวจสอบว่าภาพมีอยู่จริง
            if not os.path.exists(image_path):
                print(f"❌ ไม่พบไฟล์ภาพที่ path: {image_path}")
                return

            # โหลดและแปลงภาพ
            image = Image.open(image_path).convert("RGB")
            faces = app.get(np.array(image))

            print(f"🧠 imageId={image_id} เจอใบหน้าทั้งหมด: {len(faces)}")

            # >>> ถ้าคุณต้องการส่งผลลัพธ์กลับหรือบันทึก embedding ให้เพิ่มตรงนี้ <<<

        except Exception as e:
            print(f"❌ เกิดข้อผิดพลาด imageId={payload.get('imageId')}: {str(e)}")

async def main():
    try:
        connection = await aio_pika.connect_robust(
            "amqp://skko:skkospiderman@rabbitmq/",
            heartbeat=60,
            timeout=60
        )
        channel = await connection.channel()
        queue = await channel.declare_queue("face_jobs", durable=True)
        await queue.consume(on_message)
        print("✅ AI Service พร้อมทำงานแล้ว รอรับงาน...")
        await asyncio.Future()  # รอไม่สิ้นสุด
    except Exception as e:
        print("❌ เชื่อมต่อ RabbitMQ ไม่สำเร็จ:", str(e))
