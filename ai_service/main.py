import aio_pika
import asyncio
import json
import mysql.connector
from insightface.app import FaceAnalysis
from PIL import Image
import numpy as np
import os

# ⚠️ ต้องอยู่ก่อน import onnxruntime
os.environ['ORT_DISABLE_CPU_AFFINITY'] = '1'

# พาธรูปภาพที่ mount ไว้ใน container
IMAGE_BASE_PATH = "/app/images"

# 🔍 โหลดโมเดล InsightFace
app = FaceAnalysis(name="buffalo_l")
app.prepare(ctx_id=0, det_size=(640, 640))  # ctx_id=0 = GPU

# ⚙️ การตั้งค่า MySQL
db_config = {
    'host': '192.168.0.121',
    'user': 'mysql_121',
    'password': 'hdcdatarit9esoydld]o8i',
    'database': 'officedd_photo'
}

# 💾 ฟังก์ชันบันทึก embeddings ลง MySQL (1 ใบหน้า = 1 record)
async def save_to_db(image_id, faces):
    try:
        connection = mysql.connector.connect(**db_config)
        cursor = connection.cursor()

        for idx, face in enumerate(faces):
            emb_json = json.dumps(face.embedding.tolist())

            insert_query = """
                INSERT INTO face_embeddings (image_id, face_index, embedding)
                VALUES (%s, %s, %s)
            """
            cursor.execute(insert_query, (image_id, idx, emb_json))

        # ✅ update process_status_id และจำนวนใบหน้า
        update_query = "UPDATE images SET process_status_id = %s, faces = %s WHERE images_id = %s"
        cursor.execute(update_query, (3, len(faces), image_id))

        connection.commit()
        cursor.close()
        connection.close()
        print(f"✅ บันทึก {len(faces)} ใบหน้าสำหรับ image_id={image_id} สำเร็จ")

    except mysql.connector.Error as err:
        print(f"❌ เกิดข้อผิดพลาดในการบันทึกข้อมูล: {err}")

# 📩 ฟังก์ชันประมวลผลเมื่อมีข้อความเข้า RabbitMQ
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

                    # 🔍 ตรวจจับใบหน้า
                    faces = app.get(image_np)

                    print(f"🧠 พบ {len(faces)} ใบหน้าในภาพนี้")

                    await save_to_db(image_id, faces)

                except Exception as e:
                    print(f"❌ เกิดข้อผิดพลาดกับ image_id={image_id}: {str(e)}")

        except Exception as e:
            print(f"❌ เกิดข้อผิดพลาดในการประมวลผลข้อความ: {str(e)}")

# 🚀 main() รอรับข้อความจาก RabbitMQ
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
