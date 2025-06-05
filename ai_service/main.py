from insightface.app import FaceAnalysis
import aio_pika
import asyncio
import json
import numpy as np
from PIL import Image
import aiohttp

# เตรียม InsightFace
app = FaceAnalysis(name="buffalo_l")
app.prepare(ctx_id=0, det_size=(640, 640))

# URL ของ API ที่จะส่งผลลัพธ์ไป
API_URL = "http://your-api-service/api/face-result"  # เปลี่ยน URL ตรงนี้ตามของคุณ

async def send_to_api(data: dict):
    async with aiohttp.ClientSession() as session:
        try:
            async with session.post(API_URL, json=data) as response:
                if response.status == 200:
                    print("✅ ส่งผลลัพธ์ไป API สำเร็จ")
                else:
                    print(f"⚠️ API ตอบกลับ: {response.status}")
        except Exception as e:
            print(f"❌ ส่ง API ไม่สำเร็จ: {str(e)}")

async def on_message(message: aio_pika.IncomingMessage):
    async with message.process():
        try:
            payload = json.loads(message.body.decode())
            image_id = payload["imageId"]
            image_path = payload["path"]

            print(f"📥 รับงาน: imageId={image_id}, path={image_path}")

            image = Image.open(image_path).convert("RGB")
            image_np = np.array(image)

            faces = app.get(image_np)
            print(f"🧠 พบ {len(faces)} ใบหน้า")

            results = []
            for face in faces:
                result = {
                    "bbox": face.bbox.tolist(),
                    "embedding": face.embedding.tolist(),
                    "landmark": face.kps.tolist(),
                }
                results.append(result)

            # เตรียมข้อมูลส่ง API
            output = {
                "imageId": image_id,
                "face_count": len(faces),
                "faces": results
            }

            # ส่งไปยัง API
            await send_to_api(output)

        except Exception as e:
            print(f"❌ เกิดข้อผิดพลาดกับ imageId={payload.get('imageId')}: {str(e)}")

async def main():
    try:
        connection = await aio_pika.connect_robust("amqp://skko:skkospiderman@rabbitmq/")
        channel = await connection.channel()
        queue = await channel.declare_queue("face_job", durable=True)

        await queue.consume(on_message)
        print("✅ AI Service พร้อมทำงานแล้ว...")
        await asyncio.Future()
    except Exception as e:
        print("❌ เชื่อมต่อ RabbitMQ ไม่สำเร็จ:", str(e))

if __name__ == "__main__":
    asyncio.run(main())
