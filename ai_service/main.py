import aio_pika
import asyncio
import json
import mysql.connector
from insightface.app import FaceAnalysis
from PIL import Image
import numpy as np
import os
os.environ['ORT_DISABLE_CPU_AFFINITY'] = '1'  # 👈 ต้องอยู่ตรงนี้ก่อน import onnxruntime หรือ insightface

IMAGE_BASE_PATH = "/app/images"  # ปรับตาม path ที่ mount จริงใน container

# ตั้งค่า InsightFace
app = FaceAnalysis(name="buffalo_l")
app.prepare(ctx_id=0, det_size=(640, 640))

# ตั้งค่า MySQL
db_config = {
    'host': '192.168.0.121',
    'user': 'mysql_121',
    'password': 'hdcdatarit9esoydld]o8i',
    'database': 'officedd_photo'
}

async def save_to_db(image_id, embeddings):
    try:
        connection = mysql.connector.connect(**db_config)
        cursor = connection.cursor()

        # แปลง embeddings เป็น JSON string
        embeddings_json = json.dumps(embeddings)

        # SQL query สำหรับบันทึกข้อมูล
        query = "INSERT INTO face_embeddings (image_id, embeddings) VALUES (%s, %s)"
        cursor.execute(query, (image_id, embeddings_json))
          # 2. Update process_status_id ในตาราง images เป็น 2
        update_query = "UPDATE images SET process_status_id = %s WHERE id = %s"
        cursor.execute(update_query, (2, image_id))

        connection.commit()
        cursor.close()
        connection.close()
        print(f"✅ บันทึก embeddings สำหรับ image_id={image_id} สำเร็จ")
    except mysql.connector.Error as err:
        print(f"❌ เกิดข้อผิดพลาดในการบันทึกข้อมูล: {err}")

async def on_message(message: aio_pika.IncomingMessage):
    async with message.process():
        try:
            payload = json.loads(message.body.decode())
            images = payload.get("images", [])
            for img in images:
                image_id = img.get("image_id")
                image_filename = img.get("image_name")
                image_path = os.path.join(IMAGE_BASE_PATH, image_filename)

                print(f"📥 กำลังประมวลผล image_id={image_id} path={image_path}")

                try:
                    image = Image.open(image_path).convert("RGB")
                    image_np = np.array(image)
                    faces = app.get(image_np)

                    print(f"🧠 พบ {len(faces)} ใบหน้า")

                    embeddings = [face.embedding.tolist() for face in faces]
                    await save_to_db(image_id, embeddings)

                except Exception as e:
                    print(f"❌ เกิดข้อผิดพลาดกับ image_id={image_id}: {str(e)}")
        except Exception as e:
            print(f"❌ เกิดข้อผิดพลาดในการประมวลผลข้อความ: {str(e)}")
async def main():
    try:
        connection = await aio_pika.connect_robust("amqp://skko:skkospiderman@rabbitmq:5672/")
        channel = await connection.channel()
        queue = await channel.declare_queue("face_images_queue", durable=True)

        await queue.consume(on_message)
        print("✅ AI Service พร้อมทำงานแล้ว...")
        await asyncio.Future()
    except Exception as e:
        print(f"❌ เชื่อมต่อ RabbitMQ ไม่สำเร็จ: {str(e)}")

if __name__ == "__main__":
    asyncio.run(main())
