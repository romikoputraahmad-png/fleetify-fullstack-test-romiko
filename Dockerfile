# Menggunakan image golang versi alpine agar ringan
FROM golang:1.26-alpine

# Menentukan direktori kerja di dalam container
WORKDIR /app

# Menyalin file go.mod dan go.sum, lalu mengunduh dependencies
COPY go.mod go.sum ./
RUN go mod download

# Menyalin seluruh kode sumber ke dalam container
COPY . .

# Membangun aplikasi Go menjadi file binary bernama 'main'
RUN go build -o main .

# Membuka port 3000
EXPOSE 3000

# Perintah untuk menjalankan aplikasi
CMD ["./main"]