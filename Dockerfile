# Aşama 1: Uygulamayı Derleme (Builder)
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Modül dosyalarını kopyala
COPY go.mod ./
# (Projemizde şu an dış bağımlılık olmadığı için go.sum dosyasına gerek yok, 
#  ileride eklerseniz COPY go.sum ./ satırını açabilirsiniz)

# Tüm proje kodunu kopyala
COPY . .

# Go uygulamasını statik bir binary olarak derle
# CGO_ENABLED=0 ile tamamen bağımsız bir çalıştırılabilir dosya elde ederiz
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o micro-es ./cmd/server

# Aşama 2: Çalıştırma (Production)
# Sadece çalıştırılabilir dosyayı alarak imaj boyutunu (image size) çok küçültüyoruz
FROM alpine:latest

WORKDIR /root/

# Derlenmiş dosyayı birinci aşamadan (builder) kopyala
COPY --from=builder /app/micro-es .

# API portumuzu dışarı açıyoruz
EXPOSE 8080

# Uygulamayı başlat
CMD ["./micro-es"]
