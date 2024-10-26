package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/fatih/color"
)

func cc() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func art() {
	color.Set(color.FgHiMagenta)
	fmt.Println(`

	░░░     ░     ░    ░ ░   ░     ░   ░   ░ ░░               ░      
 ░                        ░                    ░░ ░         ░  ░ 
 ░          ░░       ░     ░                  ░    ░ ░    ░     ░
   ░               ░   ░                     ░        ░ ░        
 ░                ░ ░  ░           ░        ░░           ░   ░░  
  ░                      ░      ░    ███████  ░                ░ 
             ░ ░█▓▒███▓     ░  ░ ██ ░  █████▓░█    ░             
    ░░░░       ░  ▒██████▒█░    ██  ██▓      ▒█░  ░              
░           ░         ░ █   ░    ░     ░ ▒█    ░  ▓██    █ ░     
  ░░    ▒░   ██  ▓   ░  █ ░░   ░     ░ ░     ▓█▓░          █     
    ░       ░         ██ ░       ░▒  ░    ░░    ░░ ▒███▒   ░█░  ░
  ░        █░  ░   ░██           ░  █    ░   ░░███    ████ █     
          ███         █░░░  ░    ▓ ██     █████▒   ░ ██   ░  ░░  
         ░██████       ░██▓  ░  ░░  ████ ░    █▒ ▒████ ░         
    ░     ██ █   ████████████████░  ░▓   ░  █████░ █░ ░          
░         ██░█░  █     █▒    █  ░    ▒███████ █   █    ░         
          ██████████████████████████████▒  ░   ███  ░ ░  ░       
      ░  ░████████████████████████▒  ░▒█      ██         ░░░     
          ██ █░██░▓█▓░▒██░    ░░░      █   ▒██      ░   ░     ░  
  ░  ░░  ░ ▓██  █  █▒   █     ▓█        ███░ ░         ░     ░░░ 
        ░  ░ ░████▓░█▒ ░▒█   ░░█▓▒████▒         ░             ░  
  ░        ░                                 ░        ░       ░  
         ░    ░░                   ░          ░     ░            
    ░░ ░   ░    ░ ░░    ░     ░░      ░         ░      ░ ░       
        ░        ░    ░               ░  ░           ░          ░
	`)
	color.Unset()
}

func vw(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		color.Red("[!] Error Checking Webhook: %v", err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent
}

func send(url, jd string) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(jd)))
	if err != nil {
		color.Red("[!] Error Sending message: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		color.Green("[+] Message Sent!")
	} else {
		color.Red("[!] Failed To Send Message: %d", resp.StatusCode)
		if resp.StatusCode == 429 {
			color.Red("[!] Rate limited, pausing for 6 seconds...")
			time.Sleep(6 * time.Second)
		}
	}
}

func delw(url string) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		color.Red("[!] Error Sending 'DELETE' Request: %v", err)
		return
	}
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		color.Red("[!] Error Deleting Webhook!: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		color.Green("[+] Webhook Deleted!")
	}
}

func main() {
	cc()
	art()

	var wh string
	fmt.Print(color.CyanString("Enter A Webhook URL: "))
	fmt.Scanln(&wh)

	if !vw(wh) {
		color.Red("[!] Invalid Webhook URL!")
		return
	}

	for {
		cc()
		art()
		color.Cyan("Pick One:")
		color.Cyan("1. Send Message To Webhook")
		color.Cyan("2. Delete Webhook")
		fmt.Print(color.CyanString("Enter 1 or 2: "))

		var choice int
		fmt.Scanln(&choice)

		if choice == 1 {
			var embedH, embedMSG, embedF, embedP string
			var di string
			var delay time.Duration

			fmt.Print(color.CyanString("Delay Per Message (in seconds): "))
			fmt.Scanln(&di)
			fmt.Sscanf(di, "%d", &delay)

			color.Cyan("Choose the format:")
			color.Cyan("1. Message With Embed")
			color.Cyan("2. Embed Only")
			color.Cyan("3. Message Only")

			var fco int
			fmt.Print(color.CyanString("Enter 1, 2, or 3: "))
			fmt.Scanln(&fco)

			if fco == 1 || fco == 2 {
				fmt.Print(color.CyanString("Enter Embed Header: "))
				fmt.Scanln(&embedH)

				fmt.Print(color.CyanString("Enter Embed Middle Message: "))
				fmt.Scanln(&embedMSG)

				fmt.Print(color.CyanString("Enter Embed Picture URL (or skip): "))
				fmt.Scanln(&embedP)

				fmt.Print(color.CyanString("Enter Embed Footer: "))
				fmt.Scanln(&embedF)
			} else if fco == 3 {
				fmt.Print(color.CyanString("Enter your message: "))
				fmt.Scanln(&embedMSG)
			} else {
				color.Red("[!] Invalid choice. Please select 1, 2, or 3.")
				continue
			}

			var mainMessage string
			if fco == 1 {
				fmt.Print(color.CyanString("Enter your message (outside embed): "))
				fmt.Scanln(&mainMessage)
			}

			for {
				var jd string

				switch fco {
				case 1:
					jd = fmt.Sprintf(`{
						"content": "%s",
						"embeds": [{
							"title": "%s",
							"description": "%s",
							"footer": {"text": "%s"},
							"image": {"url": "%s"},
							"color": 1752220
						}]
					}`, mainMessage, embedH, embedMSG, embedF, embedP)

				case 2:
					jd = fmt.Sprintf(`{
						"embeds": [{
							"title": "%s",
							"description": "%s",
							"footer": {"text": "%s"},
							"image": {"url": "%s"},
							"color": 1752220
						}]
					}`, embedH, embedMSG, embedF, embedP)

				case 3:
					jd = fmt.Sprintf(`{
						"content": "%s"
					}`, embedMSG)

				default:
					color.Red("[!] Invalid choice.")
					continue
				}

				send(wh, jd)
				time.Sleep(delay * time.Second)
			}
		} else if choice == 2 {
			delw(wh)
			break
		} else {
			color.Red("[!] Pick 1 or 2")
		}
	}
}
