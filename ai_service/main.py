import os
os.environ['ORT_DISABLE_CPU_AFFINITY'] = '1'  # ต้องอยู่ก่อน import insightface หรือ onnxruntime

import aio_pika
import asyncio
import json
import mysql.connector
from insightface.app import FaceAnalysis
from PIL import Image
import numpy as np

IMAGE_BASE_PATH = "/app/images"  # ปรับตามที่ mount ใน container

# ตั้งค่า InsightFace
app = FaceAnalysis(name="buffalo_l")
app.prepare(ctx_id=0, det_size=(2048, 2048))
app.threshold = 0.4

# ตั้งค่า MySQL
db_config = {
    'host': '192.168.0.121',
    'user': 'mysql_121',
    'password': 'hdcdatarit9esoydld]o8i',
    'database': 'officedd_photo'
}

async def save_to_db(image_id, embeddings, face_count):
    connection = None
    cursor = None
    try:
        connection = mysql.connector.connect(**db_config)
        cursor = connection.cursor()

        flat_embeddings = np.array(embeddings).flatten().tolist()
        embeddings_json = json.dumps(flat_embeddings)

        # บันทึก embeddings
        query = "INSERT INTO face_embeddings (image_id, embeddings) VALUES (%s, %s)"
        cursor.execute(query, (image_id, embeddings_json))

        # อัพเดตสถานะ process และจำนวนใบหน้า
        update_query = "UPDATE images SET process_status_id = %s, faces = %s WHERE images_id = %s"
        cursor.execute(update_query, (3, face_count, image_id))

        connection.commit()
        print(f"✅ บันทึก embeddings สำหรับ image_id={image_id} สำเร็จ")
    except mysql.connector.Error as err:
        print(f"❌ เกิดข้อผิดพลาดในการบันทึกข้อมูล: {err}")
    finally:
        if cursor:
            cursor.close()
        if connection:
            connection.close()

async def on_message(message: aio_pika.IncomingMessage):
    async with message.process():  # ack message เมื่อออกจากบล็อกนี้
        try:
            payload = json.loads(message.body.decode())
            image_id = payload.get("image_id")
            image_name = payload.get("image_name")

            if not image_id or not image_name:
                print("❌ ข้อมูลใน message ไม่ครบ")
                return

            image_path = os.path.join(IMAGE_BASE_PATH, image_name)
            print(f"📥 กำลังประมวลผล image_id={image_id} path={image_path}")

            # โหลดรูปและประมวลผลใบหน้า
            image = Image.open(image_path).convert("RGB")
            image_np = np.array(image)
            faces = app.get(image_np)
            face_count = len(faces)

            print(f"🧠 พบใบหน้า {face_count} ใบหน้า")

            embeddings = [face.embedding.tolist() for face in faces]

            await save_to_db(image_id, embeddings, face_count)

        except Exception as e:
            print(f"❌ เกิดข้อผิดพลาดขณะประมวลผล: {e}")

async def main():
    try:
        # เชื่อมต่อ RabbitMQ
        connection = await aio_pika.connect_robust("amqp://skko:skkospiderman@rabbitmq:5672/")
        channel = await connection.channel()

        # รับทีละ 1 งาน (prefetch=1)
        await channel.set_qos(prefetch_count=1)

        queue = await channel.declare_queue("images_queue", durable=True)
        await queue.consume(on_message)

        print("✅ AI Worker พร้อมทำงานแล้ว...")

        # รัน loop แบบไม่จบเพื่อรอรับ message ตลอดไป
        await asyncio.Future()

    except Exception as e:
        print(f"❌ เชื่อมต่อ RabbitMQ ไม่สำเร็จ: {e}")

if __name__ == "__main__":
    asyncio.run(main())
