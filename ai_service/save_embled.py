import mysql.connector
import json

def connect_db():
    return mysql.connector.connect(
        host="192.168.0.121",        # เปลี่ยนตามของคุณ
        user="mysql_121",
        password="hdcdatarit9esoydld]o8i",
        database="officedd_photo"
    )

def save_face_embedding(conn, image_id, face_id, embedding, bbox, landmark):
    cursor = conn.cursor()
    sql = """
    INSERT INTO face_embeddings (image_id, face_id, embedding, bbox, landmark)
    VALUES (%s, %s, %s, %s, %s)
    """
    # แปลงข้อมูลเป็น JSON string ก่อนบันทึก
    cursor.execute(sql, (
        image_id,
        face_id,
        json.dumps(embedding),
        json.dumps(bbox),
        json.dumps(landmark),
    ))
    conn.commit()
    cursor.close()

# ตัวอย่างใช้บันทึกผลลัพธ์หลายใบหน้า
def save_faces(conn, image_id, faces):
    for i, face in enumerate(faces):
        embedding = face["embedding"]
        bbox = face["bbox"]
        landmark = face["landmark"]
        save_face_embedding(conn, image_id, i, embedding, bbox, landmark)

# ตัวอย่างเรียกใช้งาน
if __name__ == "__main__":
    conn = connect_db()

    # ตัวอย่างข้อมูลจาก InsightFace
    faces = [
        {
            "embedding": [0.1, 0.2, 0.3, 0.4],
            "bbox": [10, 20, 110, 120],
            "landmark": [[15, 25], [30, 40], [45, 55], [60, 70], [75, 85]]
        },
        # ใบหน้าที่สองถ้ามี
    ]
    image_id = 1234

    save_faces(conn, image_id, faces)
    conn.close()
