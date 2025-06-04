FROM python:3.10-slim

# ตั้ง working directory ใน container
WORKDIR /app

# ติดตั้ง dependencies ที่จำเป็นสำหรับ opencv, insightface และอื่น ๆ
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    libgl1-mesa-glx \
    libglib2.0-0 \
    libffi-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements.txt เข้า container
COPY requirements.txt ./

# ติดตั้ง python packages
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code เข้า container
COPY . .

# ถ้าคุณมีไฟล์ main.py หรือ app.py ที่จะรัน เช่น
CMD ["python", "main.py"]
