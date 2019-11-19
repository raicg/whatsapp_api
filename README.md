---------------------------

# Setup

## Pre-requisites

* GoLang

## Need help to install golang?
try: https://github.com/canha/golang-tools-install-script

## Setup project

1. `$ go get github.com/Rhymen/go-whatsapp`
1. `$ go get github.com/Baozisoftware/qrcode-terminal-go`
1. `$ go get github.com/tushar2708/altcsv`
1. Then edit the variables in `config/config.go` with the correct path to your folders. (follow the instructions inside the file)


## Running

1. `$ go run main.go`

## Login
1. To login you just need to read the QR code using the whatsapp on your device, you will connect it through the whatsapp web function.

## Receiving Messages
1. To receive the messages you just need to configure the variables in `config/config.go`, as requested on the setup project section, so you will receive the messages in csv format and the attachments.

## Sending Messages
1. To send a message you just need to enter the text, "contact" or "group" and the phone number without + sign or group id without @g.us, like: `"text example" "contact" "5584998765432"` or `"text example" "group" "557392152628-1538763567"` on the terminal.
