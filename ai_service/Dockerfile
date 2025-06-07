FROM python:3.10-slim

WORKDIR /app

# ติดตั้ง dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    libgl1-mesa-glx \
    libglib2.0-0 \
    libffi-dev \
    && rm -rf /var/lib/apt/lists/*

# คัดลอก requirements แล้วติดตั้งก่อน (cache stage)
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

# คัดลอกโมเดลไปยังตำแหน่ง insightface ใช้งาน
RUN mkdir -p /root/.insightface/models/buffalo_l
COPY models/buffalo_l /root/.insightface/models/buffalo_l

# คัดลอก source code
COPY . .

CMD ["python", "main.py"]
