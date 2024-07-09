#!/bin/bash

# Eğer telefon numarası ve mesaj girilmemişse kullanıcıdan iste
if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Kullanım: ./sens.sh \"mesaj\" \"telefon_numarası\""
    exit 1
fi

message="$1"
DEST="$2"

# Doğru seri portu belirleyin
MD="/dev/ttyUSB0" # Bu portu bağlı cihazınıza göre ayarlayın

# Fonksiyon: Betiği kapat
cleanup() {
    echo "Interrupt signal yakalandı... Betik kapatılıyor."
    stty -F $MD hupcl  # Seri portu kapat
    exit 1
}

# Interrupt sinyalini yakala
trap cleanup INT

echo "Mesaj \"$message\" $DEST numarasına gönderiliyor."

# Seri port ayarları
if [ -e $MD ]; then
    if ! fuser -s $MD; then
        echo "Seri port ayarlanıyor..."
        stty -F $MD 9600 min 100 time 2 -hupcl brkint ignpar -opost -onlcr -isig -icanon -echo

        # SMS gönderme işlemi
        {
            sleep 0.5
            echo -e "AT+CMGF=1\r" > $MD
            sleep 0.5
            echo -e "AT+CMGS=\"$DEST\"\r" > $MD
            sleep 0.5
            echo -e "$message\x1A" > $MD
            sleep 0.5
        } < $MD > $MD

        # Modem yanıtını oku (SMS gönderme durumu)
        echo "Modemden yanıt bekleniyor..."
        response=$(timeout 10s cat < $MD)
        echo "Modem yanıtı: $response"

        # Başarılı olup olmadığını kontrol et
        if echo "$response" | grep -q "OK"; then
            echo "SMS başarıyla gönderildi."
        else
            echo "SMS gönderilemedi."
        fi

        # Seri portu kapat
        stty -F $MD hupcl

        # Betiği sonlandır
        exit 0
    else
        echo "$MD: Seri port başka bir işlem tarafından kullanılıyor."
        exit 1
    fi
else
    echo "$MD: Seri port bulunamadı."
    exit 1
fi
