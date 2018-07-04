package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type questionID struct {
	ID string `json:"id"`
}

type feedBack struct {
	Hit   int    `json:"hit"`
	Blow  int    `json:"blow"`
	Round int    `json:"round"`
	Msg   string `json:"message"`
}

const token = "localで設定してください"

func questionRequest() string {
	url := "https://apiv2.twitcasting.tv/internships/2018/games?level=10"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := new(http.Client)
	resp, err := client.Do(req)

	if err != nil {
		return "Invalid token"
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	jsonBytes := ([]byte)(byteArray)
	data := new(questionID)

	if err := json.Unmarshal(jsonBytes, data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return "Can't catch id"
	}

	fmt.Println(data.ID)
	return data.ID
}

func solveProblem(id, ans string) (int, int) {
	url := "https://apiv2.twitcasting.tv/internships/2018/games/" + id
	jsonStr := `{"answer":"` + ans + `"}`
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return 0, 0
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	jsonBytes := ([]byte)(byteArray)
	data := new(feedBack)

	if err := json.Unmarshal(jsonBytes, data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return 0, 0
	}

	fmt.Printf("hit: %d, blow: %d, message: %s", data.Hit, data.Blow, data.Msg)
	fmt.Println("========")
	return data.Hit, data.Blow
}

func allAlgorithm(id string) {
	arr := []string{"9", "8", "7", "6", "5", "4", "3", "2", "1", "0"}
	rm := []int{0}

	ih, _ := solveProblem(id, "9876543210")
	if ih == 10 {
		return
	}

	for i := 0; i < 10; i++ {
		fmt.Println("===========", i, "===========")

		for j := i + 1; j < 10; j++ {

			flag := false
			for fact := 1; fact < len(rm); fact++ {
				if i == rm[fact] {
					flag = true
				}
			}
			if flag == true {
				break
			}

			arr[i], arr[j] = arr[j], arr[i]
			h, _ := solveProblem(id, order(arr))
			if h == 10 {
				return
			}

			if h == ih-2 { //ooからxxになった場合
				rm = append(rm, j)
				arr[i], arr[j] = arr[j], arr[i]
				break
			} else if h == ih+2 { //xxからooになった場合
				rm = append(rm, j)
				ih = h
				break
			} else if h == ih+1 { //xx=> xo or xx => ox
				// i = 0でここに振り分けられると失敗
				arr[i], arr[rm[0]] = arr[rm[0]], arr[i]
				c, _ := solveProblem(id, order(arr))

				if c == h-2 { //arr[i] =oが決定
					arr[i], arr[rm[0]] = arr[rm[0]], arr[i]
					ih++
					break
				} else if c == h-1 { //arr[i] = xが決定
					arr[i], arr[rm[0]] = arr[rm[0]], arr[i]
					arr[i], arr[j] = arr[j], arr[i]

				}
			} else if h == ih-1 { //ox => xx or xo => xx

				arr[i], arr[j] = arr[j], arr[i]
				for k := j + 1; k < 10; k++ {
					arr[i], arr[k] = arr[k], arr[i]
					c, _ := solveProblem(id, order(arr))

					if c == ih+1 || c == ih+2 { //xo =>常に-1 or xxがきたら0or+1or+2
						ih = c
						break //for k をbreak
					} else if c == ih-2 { //ox =>常に-1 or ooがきたら-2
					}
					arr[i], arr[k] = arr[k], arr[i]
				}
				break //for j をbreak
			}
		}
	}
}

func order(arr []string) string {
	return arr[0] + arr[1] + arr[2] + arr[3] + arr[4] + arr[5] + arr[6] + arr[7] + arr[8] + arr[9]
}

func main() {
	id := questionRequest()
	allAlgorithm(id)
}
