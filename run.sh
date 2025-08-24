#!/bin/bash

# Required
export EMAIL_FROM="unknow@gmail.com"
export EMAIL_PASSWORD="00000000000000"
export EMAIL_TO="unknow@qq.com"
export EMAIL_SMTP_HOST="smtp.gmail.com"
export EMAIL_SMTP_PORT="465"

BIN_PATH="bin/mailsender-windows-amd64.exe"
$BIN_PATH -s TEST -c TEST