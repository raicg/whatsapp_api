package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"strconv"
	"strings"
	config "./config"

	altcsv "github.com/tushar2708/altcsv"
	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
)

func main() {
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(20 * time.Second)
	if err != nil {
		log.Fatalf("error creating connection: %v\n", err)
	}

	//Add handler
	wac.AddHandler(&waHandler{wac})

	//login or restore
	if err := login(wac); err != nil {
		log.Fatalf("error logging in: %v\n", err)
	}

	go func() {
		var text, phone string
		fmt.Print(`Intructions: You will enter the text and the phone number (like: "text example" "5584998765432") if you want to send a message;`+"\n")
		fmt.Print("phone number is just number, so dont have the + sign.\n")
		for {
			fmt.Print("Enter the text and the phone number: \n")
			_, err := fmt.Scanf("%q", &text)
			_, err = fmt.Scanf("%q \n", &phone)
			msg := whatsapp.TextMessage{
				Info: whatsapp.MessageInfo{
					RemoteJid:     phone + "@s.whatsapp.net",
				},
				Text: text,
			}
		
			msgId, err := wac.Send(msg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error sending message: %v", err)
				os.Exit(1)
			} else {
				fmt.Println("Message Sent -> ID : " + msgId)
			}
		}
	}()

	//verifies phone connectivity
	pong, err := wac.AdminTest()

	if !pong || err != nil {
		log.Fatalf("error pinging in: %v\n", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

type waHandler struct {
	c *whatsapp.Conn
}

//HandleError needs to be implemented to be a valid WhatsApp handler
func (h *waHandler) HandleError(err error) {

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.c.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
	
}

//Optional to be implemented. Implement HandleXXXMessage for the types you need.
func (w *waHandler) HandleTextMessage(message whatsapp.TextMessage) {
	phone := ""
	if strings.Contains(message.Info.RemoteJid, "@s.whatsapp.net"){
		phone = strings.Replace(message.Info.RemoteJid, "@s.whatsapp.net", "", 1)
	} else {
		group_data, err := w.c.GetGroupMetaData(message.Info.RemoteJid)
		_ = err
		createGroupCsv(<- group_data, message)
	}
	
	createTextCsv(phone, message)
}

// Example for media handling. Video, Audio, Document are also possible in the same way
func (h *waHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	data, err := message.Download()
	if err != nil {
		if err != whatsapp.ErrMediaDownloadFailedWith410 && err != whatsapp.ErrMediaDownloadFailedWith410 {
			return
		}
		if _, err = h.c.LoadMediaInfo(message.Info.SenderJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe)); err == nil {
			data, err = message.Download()
			if err != nil {
				return
			}
		}
	}

	filename := message.Info.Id + "." + strings.Split(message.Type, "/")[1]
	file, err := os.Create(config.ImagesFolder + filename)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.Write(data)
	if err != nil {
		return
	}
	createImageCsv(message)
}

func (h *waHandler) HandleVideoMessage(message whatsapp.VideoMessage) {
	data, err := message.Download()
	if err != nil {
		if err != whatsapp.ErrMediaDownloadFailedWith410 && err != whatsapp.ErrMediaDownloadFailedWith410 {
			return
		}
		if _, err = h.c.LoadMediaInfo(message.Info.SenderJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe)); err == nil {
			data, err = message.Download()
			if err != nil {
				return
			}
		}
	}

	filename := message.Info.Id + "." + strings.Split(message.Type, "/")[1]
	file, err := os.Create(config.VideosFolder + filename)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.Write(data)
	if err != nil {
		return
	}
	createVideoCsv(message)
}

func (h *waHandler) HandleAudioMessage(message whatsapp.AudioMessage) {
	data, err := message.Download()
	if err != nil {
		if err != whatsapp.ErrMediaDownloadFailedWith410 && err != whatsapp.ErrMediaDownloadFailedWith410 {
			return
		}
		if _, err = h.c.LoadMediaInfo(message.Info.SenderJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe)); err == nil {
			data, err = message.Download()
			if err != nil {
				return
			}
		}
	}

	filename := message.Info.Id + "." + strings.Split(strings.Split(message.Type, "/")[1], ";")[0]
	file, err := os.Create(config.AudiosFolder + filename)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.Write(data)
	if err != nil {
		return
	}
	createAudioCsv(message)
}

func (h *waHandler) HandleDocumentMessage(message whatsapp.DocumentMessage) {
	data, err := message.Download()
	if err != nil {
		if err != whatsapp.ErrMediaDownloadFailedWith410 && err != whatsapp.ErrMediaDownloadFailedWith410 {
			return
		}
		if _, err = h.c.LoadMediaInfo(message.Info.SenderJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe)); err == nil {
			data, err = message.Download()
			if err != nil {
				return
			}
		}
	}

	filename := message.Info.Id + "." + strings.Split(message.Type, "/")[1]
	file, err := os.Create(config.DocumentsFolder + filename)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.Write(data)
	if err != nil {
		return
	}
	createDocumentCsv(message)
}

func login(wac *whatsapp.Conn) error {
	qr := make(chan string)
	go func() {
		terminal := qrcodeTerminal.New()
		terminal.Get(<-qr).Print()
	}()
	session, err := wac.Login(qr)
	if err != nil {
		return fmt.Errorf("error during login: %v\n", err)
	}
	_ = session
	return nil
}

func createTextCsv(phone string, message whatsapp.TextMessage) {
	file, _ := os.Create(config.TextsFolder + message.Info.Id + ".csv")
	csvWtr := altcsv.NewWriter(file)
	csvWtr.Quote = '\''
	csvWtr.AllQuotes = true
	csvWtr.Write([]string{message.Info.Id, strconv.FormatUint(message.Info.Timestamp, 10), message.Info.RemoteJid, phone, message.Info.QuotedMessageID, message.Text})
	csvWtr.Flush()
	file.Close()
}

func createGroupCsv(group_data string, message whatsapp.TextMessage) {
	file, _ := os.Create(config.GroupsFolder + message.Info.RemoteJid + ".csv")
	csvWtr := altcsv.NewWriter(file)
	csvWtr.Quote = '\''
	csvWtr.AllQuotes = true
	csvWtr.Write([]string{message.Info.RemoteJid, "{" + group_data[1 : len(group_data)-1] + "}"})
	csvWtr.Flush()
	file.Close()
}

func createImageCsv(message whatsapp.ImageMessage) {
	file, _ := os.Create(config.ImagesFolder + message.Info.Id + ".csv")
	csvWtr := altcsv.NewWriter(file)
	csvWtr.Quote = '\''
	csvWtr.AllQuotes = true
	csvWtr.Write([]string{message.Info.Id, strconv.FormatUint(message.Info.Timestamp, 10), message.Info.RemoteJid, message.Caption})
	csvWtr.Flush()
	file.Close()
}

func createVideoCsv(message whatsapp.VideoMessage) {
	file, _ := os.Create(config.VideosFolder + message.Info.Id + ".csv")
	csvWtr := altcsv.NewWriter(file)
	csvWtr.Quote = '\''
	csvWtr.AllQuotes = true
	csvWtr.Write([]string{message.Info.Id, strconv.FormatUint(message.Info.Timestamp, 10), message.Info.RemoteJid, message.Caption})
	csvWtr.Flush()
	file.Close()
}

func createAudioCsv(message whatsapp.AudioMessage) {
	file, _ := os.Create(config.AudiosFolder + message.Info.Id + ".csv")
	csvWtr := altcsv.NewWriter(file)
	csvWtr.Quote = '\''
	csvWtr.AllQuotes = true
	csvWtr.Write([]string{message.Info.Id, strconv.FormatUint(message.Info.Timestamp, 10), message.Info.RemoteJid})
	csvWtr.Flush()
	file.Close()
}

func createDocumentCsv(message whatsapp.DocumentMessage) {
	file, _ := os.Create(config.DocumentsFolder + message.Info.Id + ".csv")
	csvWtr := altcsv.NewWriter(file)
	csvWtr.Quote = '\''
	csvWtr.AllQuotes = true
	csvWtr.Write([]string{message.Info.Id, strconv.FormatUint(message.Info.Timestamp, 10), message.Info.RemoteJid, message.Title, message.FileName})
	csvWtr.Flush()
	file.Close()
}