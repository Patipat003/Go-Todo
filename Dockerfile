# ใช้ภาพพื้นฐานของ Go
FROM golang:1.20-alpine

# ตั้งค่าไดเรกทอรีทำงานในคอนเทนเนอร์
WORKDIR /app

# คัดลอกไฟล์ go.mod และ go.sum ไปยังไดเรกทอรีทำงาน
COPY go.mod go.sum ./

# ดาวน์โหลดและติดตั้ง dependencies
RUN go mod download

# คัดลอกไฟล์โค้ดทั้งหมดไปยังไดเรกทอรีทำงาน
COPY . .

# สร้างแอปพลิเคชัน
RUN go build -o main .

# ระบุพอร์ตที่แอปพลิเคชันจะใช้
EXPOSE 8080

# คำสั่งเพื่อรันแอปพลิเคชัน
CMD ["./main"]