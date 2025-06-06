FROM python:3.10-slim

WORKDIR /app

# ติดตั้งเครื่องมือและไลบรารีที่จำเป็น
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    libgl1-mesa-glx \
    libglib2.0-0 \
    libffi-dev \
    wget \
    unzip \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# ติดตั้ง Python packages
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

# ติดตั้ง insightface ที่ใหม่ที่สุด
RUN pip install --no-cache-dir --upgrade insightface

# คัดลอกโค้ดเข้า container
COPY . .

# สั่งรันแอป
CMD ["python", "main.py"]
