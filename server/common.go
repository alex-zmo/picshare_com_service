package server

import (
	"encoding/json"
	"fmt"
	"github.com/gmo-personal/picshare_com_service/database"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Headers", "Authorization")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
}

func InitServer() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/send/", sendHandler)
	http.HandleFunc("/listComments/", listCommentsHandler)
	http.HandleFunc("/delete/", deleteCommentsHandler)


	log.Fatal(http.ListenAndServe(":8083", nil))
}

func getURLParam(r *http.Request, paramName string) string {
	keys, seen := r.URL.Query()[paramName]

	if seen && len(keys) > 0 {
		return keys[0]
	}
	return ""
}

func listCommentsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	link := getURLParam(r,"src")
	picId := database.GetPicIdFromLink(link)
	userId := GetUserIdFromRequest(w, r)

	var linkList = database.ListComments(picId, userId)

	result, _ := json.Marshal(linkList)
	//linkList = append(linkList, strconv.Itoa(userId))
	fmt.Fprintln(w, string(result))
}
func deleteCommentsHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	commentIdStr := getURLParam(r,"commentID")

	commentId, _ := strconv.Atoi(commentIdStr)

	database.DeleteComment(commentId)
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	link, comment := parseForLinkAndComment(w,r)
	if len(link) == 0 || len(comment) == 0 {
		fmt.Println("Error in Send Handler, Link or comment length is 0")
		return
	}
	userId := GetUserIdFromRequest(w,r)
	if userId == -1 {
		fmt.Println("Error in UserId Request Write is Favorite Handler")
	}
	picId := database.GetPicIdFromLink(link)
	err := database.InsertComment(userId, picId, comment)
	if err != "" {
		fmt.Println(err + " at Send Handler")
		w.WriteHeader(http.StatusConflict)
	}

}

var client = http.Client{}

func GetUserIdFromRequest(w http.ResponseWriter, r *http.Request ) int{
	enableCors(&w)
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) == 1 {
		fmt.Println("Token is Empty")
		return -1
	}
	reqToken = splitToken[1]

	req, err := http.NewRequest("GET", "http://localhost:8081/validate/", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Bearer " + reqToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return -1
	}
	userIdString := string(bodyBytes)

	userId, err := strconv.Atoi(userIdString)

	if err != nil {
		fmt.Println(err)
		return -1
	}
	return userId
}

func parseForLinkAndComment(w http.ResponseWriter, r *http.Request) (string, string){
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ERROR GETTING BODY")
	}
	form := strings.Split(string(body), "EXIST")
	if len(form) < 3 {
		return "",""
	}
	link := strings.TrimSpace(form[1])
	comment := form[3]

	return link, comment
}