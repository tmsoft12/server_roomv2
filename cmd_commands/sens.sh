#!/bin/bash

# Eğer telefon numarası ve mesaj girilmemişse kullanıcıdan iste
if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Kullanım: ./sens.sh \"mesaj\" \"telefon_numarası\""
    exit 1
fi

message="$1"
DEST="$2"

MD="/dev/ttyUSB0"

cleanup() {
    echo "err"
    stty -F $MD hupcl  
    exit 1
}

trap cleanup INT

echo "SMS \"$message\" $DEST nomerine iberilyar."

if [ -e $MD ]; then
    if ! fuser -s $MD; then
        echo "Seri port gozlenyar..."
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

        echo "Modemden jogaba garasylyar..."
        response=$(timeout 10s cat < $MD)
        echo "Modem jogaby: $response"

        # Başarılı olup olmadığını kontrol et
        if echo "$response" | grep -q "OK"; then
            echo "SMS ussunlikli iberildi."
        else
            echo "SMS iberilmedi."
        fi

        stty -F $MD hupcl
        echo -e "ATZ\r" > $MD  
        sleep 0.5
        echo -e "AT&F\r" > $MD 

     
        exit 0
    else
        echo "$MD: Seri port başka bir isde ulanylyar."
        exit 1
    fi
else
    echo "$MD: Seri port tapylmady."
    exit 1
fi
