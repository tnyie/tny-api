package views

// import (
// 	"encoding/json"
// 	"io/ioutil"
// 	"log"
// 	"net/http"

// 	"github.com/gal/tny/authentication"
// )

// func PostUser(w http.ResponseWriter, r *http.Request) {
// 	data := make(map[string]interface{})

// 	bd, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		log.Println("Couldn't read json Body")
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	if err = json.Unmarshal(bd, &data); err != nil {
// 		log.Println("Couldn't parse json body\n", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	user, err := authentication.CreateUser(data)
// 	if err != nil {
// 		log.Println("Erorr creating user\n", err)
// 		if user.UID == "" {
// 			w.WriteHeader(http.StatusBadRequest)
// 		}
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
// 	user.Create()

// }
