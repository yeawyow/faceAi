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
    wget \
    unzip \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements.txt เข้า container
COPY requirements.txt ./

# ติดตั้ง python packages
RUN pip install --no-cache-dir -r requirements.txt
RUN pip install --no-cache-dir --upgrade insightface

# ดาวน์โหลดโมเดลล่วงหน้าเพื่อหลีกเลี่ยงปัญหา RuntimeError
RUN mkdir -p /root/.insightface/models && \
    wget -O /root/.insightface/models/retinaface_r50_v1.zip https://github.com/deepinsight/insightface/releases/download/v0.7/retinaface_r50_v1.zip && \
    unzip /root/.insightface/models/retinaface_r50_v1.zip -d /root/.insightface/models/retinaface_r50_v1

# Copy source code เข้า container
COPY . .

# สั่งรันแอป
CMD ["python", "main.py"]
