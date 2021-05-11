package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/dlclark/regexp2"
	"github.com/gookit/color"
)

var (
	client       = &http.Client{}
	cookieClient = &http.Client{}
	in           = color.HiBlue.Render
)

type Account struct {
	Email    string
	Username string
	Password string
}

func decode(toDecode []byte) map[string]string {
	var output map[string]string
	json.Unmarshal([]byte(toDecode), &output)
	return output
}

func getCSRF() string {
	regx := regexp2.MustCompile("(?<=\"csrf_token\":\")\\w+", 0)

	req, _ := http.NewRequest("GET", "https://www.instagram.com/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Brave Chrome/83.0.4103.116 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if m, _ := regx.FindStringMatch(string(body)); m != nil {
		return m.String()
	}
	return ""
}

func updateDetails(csrfToken string, email string, username string) {
	data := "first_name=&email=" + email + "&username=" + username + "&phone_number=&biography=" + "" + "&external_url=&chaining_enabled=on"
	req, _ := http.NewRequest("POST", "https://www.instagram.com/accounts/edit/", bytes.NewBuffer([]byte(data)))
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("x-csrftoken", csrfToken)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Brave Chrome/83.0.4103.116 Safari/537.36")
	resp, err := cookieClient.Do(req)

	if err != nil {
		color.Red.Println("An error has occured")
		fmt.Scanln()
	}

	body, _ := ioutil.ReadAll(resp.Body)
	response := string(body)
	defer resp.Body.Close()

	fmt.Printf("%v - Update Detail Resp: %v", resp.StatusCode, response)

	if strings.Contains(response, "Please wait a few minutes before you try again") {
		color.Red.Println("[+] Rate limited")
	}

}

func urlCheck(check string) bool {
	req, _ := http.NewRequest("GET", "https://www.instagram.com/"+check, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Brave Chrome/83.0.4103.116 Safari/537.36")
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		color.Green.Printf("[+] [%s] available\n", check)
		return true
	} else if resp.StatusCode == 200 && strings.Contains(string(body), "Login • Instagram") {
		color.Red.Println("[+] Failed to check username")
		return false
	} else {
		color.Yellow.Printf("[+] [%s] currently unavailable\n", check)
		return false
	}
}

func createCheck(check string) bool {
	csrfToken := getCSRF()
	data := "username=" + check + "&email=random@gmail.com&first_name=firstname&opt_into_one_tap=false&enc_password=#PWD_INSTAGRAM_BROWSER:0:0:password0000"
	req, _ := http.NewRequest("POST", "https://www.instagram.com/accounts/web_create_ajax/attempt/", bytes.NewBuffer([]byte(data)))
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("x-csrftoken", csrfToken)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Brave Chrome/83.0.4103.116 Safari/537.36")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error occured")
		fmt.Scanln()
	}

	body, _ := ioutil.ReadAll(resp.Body)
	response := string(body)
	defer resp.Body.Close()

	if strings.Contains(response, "spam") {
		color.Red.Println("[+] Spam detected")
		fmt.Println(response)
		return false
	}

	if strings.Contains(response, "try again") || resp.StatusCode == 403 {
		color.Red.Println("[+] Rate limited")
		return false
	} else if !strings.Contains(response, "\"username\":") && !strings.Contains(response, "username isn't available") && !strings.Contains(response, "username_is_taken") && !strings.Contains(response, "username_held_by_others") && resp.StatusCode != 403 {
		color.Green.Printf("[+] [%s] available\n", check)
		fmt.Println(response)
		return true

	} else {
		color.Yellow.Printf("[+] [%s] currently unavailable\n", check)
		return false
	}
}

func login(username string, password string) (*http.Response, string) {
	data := "username=" + username + "&enc_password=" + password + "&queryParams={}&optIntoOneTap=false"
	req, _ := http.NewRequest("POST", "https://www.instagram.com/accounts/login/ajax/", bytes.NewBuffer([]byte(data)))
	csrfToken := getCSRF()

	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("x-csrftoken", csrfToken)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Brave Chrome/83.0.4103.116 Safari/537.36")

	jar, _ := cookiejar.New(nil)
	cookieClient = &http.Client{Jar: jar}

	resp, err := cookieClient.Do(req)

	if err != nil {
		fmt.Println("Error occured when trying to login.")
		fmt.Scanln()
	}
	return resp, resp.Cookies()[0].Value
}

func getLines(path string) []string {
	file, err := os.Open(path)

	if err != nil {
		color.Red.Println("\nUnable to open text file:", path)
		color.Red.Printf("Make sure path \"%v\" is available\n", path)
		fmt.Scanln()
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines

}

func readAcc() Account {
	var obj Account
	data, err := ioutil.ReadFile("config/account.json")
	if err != nil {
		color.Red.Println("\nUnable to read login.json file. Make sure that config/account.json exists")
		fmt.Scanln()
	}

	err = json.Unmarshal(data, &obj)
	if err != nil {
		color.Red.Println("\nError in json.Unmarshal")
		fmt.Scanln()
	}
	return obj
}
func changeTitle(title string) (int, error) {
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer syscall.FreeLibrary(handle)
	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")

	if err != nil {
		return 0, err
	}

	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	return int(r), err
}

func printLogo() {
	color.HiBlue.Println(`
			██╗ ██████╗     ███████╗███╗   ██╗██╗██████╗ ███████╗██████╗ 
			██║██╔════╝     ██╔════╝████╗  ██║██║██╔══██╗██╔════╝██╔══██╗
			██║██║  ███╗    ███████╗██╔██╗ ██║██║██████╔╝█████╗  ██████╔╝
			██║██║   ██║    ╚════██║██║╚██╗██║██║██╔═══╝ ██╔══╝  ██╔══██╗
			██║╚██████╔╝    ███████║██║ ╚████║██║██║     ███████╗██║  ██║
			╚═╝ ╚═════╝     ╚══════╝╚═╝  ╚═══╝╚═╝╚═╝     ╚══════╝╚═╝  ╚═╝
								   by NightfallGT					 
	`)
	fmt.Print("\n\n")
}

func main() {
	fmt.Print("\033[H\033[2J")
	changeTitle("[IG Sniper] | NightfallGT")
	printLogo()

	acc := readAcc()

	var emailLogin string
	var usernameLogin string
	var passwordLogin string

	accTargets := getLines("config/targets.txt")

	emailLogin = acc.Email
	usernameLogin = acc.Username
	passwordLogin = acc.Password

	if len(emailLogin) < 1 {
		return
	}

	fmt.Printf("[%s] Email: %s\n", in("-"), emailLogin)
	fmt.Printf("[%s] Username: %s\n", in("-"), usernameLogin)
	fmt.Printf("[%s] Password: %s\n\n", in("-"), strings.Repeat("*", len(passwordLogin)))
	fmt.Printf("[%s] Press enter to begin.\n", in("+"))
	fmt.Scanln()

	fmt.Print("\033[H\033[2J")
	printLogo()
	fmt.Printf("[%s] Attempting to login through Instagram API.. \n", in("+"))

	resp, csrf := login(usernameLogin, "#PWD_INSTAGRAM_BROWSER:0:0:"+passwordLogin)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if strings.Contains(string(body), "\"authenticated\":true") {
		usernameMethod := 1
		sleepTime := 0
		changeTitle("[IG Sniper] | Logged in as: " + usernameLogin)
		fmt.Printf("[%s] Successfully logged in\n", in("*"))
		fmt.Printf("[%s] Authenticated: True, userID: %s\n", in("*"), decode(body)["userId"])
		fmt.Printf("[%s] Select username checking method. \n", in("+"))
		fmt.Println(in("---------------------------"))
		fmt.Println("\n[1] Create Username Check")
		fmt.Println("[2] URL Username Check")

		fmt.Print("Enter number choice: ")
		fmt.Scan(&usernameMethod)
		fmt.Println()
		fmt.Print("Enter sleep time [in seconds](0 for none): ")
		fmt.Scan(&sleepTime)
		fmt.Println()

		fmt.Print("\033[H\033[2J")
		printLogo()

		var attemptCount int = 0

		for {
			for _, target := range accTargets {
				attemptCount += 1
				changeTitle("[IG Sniper] | Logged in as: " + usernameLogin + " | Target: " + target + " | Request: " + strconv.Itoa(attemptCount))

				switch usernameMethod {
				case 1:
					if createCheck(target) {
						updateDetails(csrf, emailLogin, target)
					}
				case 2:
					if urlCheck(target) {
						updateDetails(csrf, emailLogin, target)
					}
				}
			}
			changeTitle("[IG Sniper] | Logged in as: " + usernameLogin + " | Sleeping.. | Request: " + strconv.Itoa(attemptCount))
			if sleepTime != 0 {
				fmt.Printf("[%s] Sleeping for %v seconds..\n", in("+"), sleepTime)
			}
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}

	} else {
		fmt.Printf("[%s] Unable to log in. Status Code: %v\n", in("!"), resp.StatusCode)
		fmt.Println(string(body))
	}
	fmt.Scanln()
}
