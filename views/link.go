package views

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"gorm.io/gorm"

	"github.com/gal/tny/middleware"
	"github.com/gal/tny/models"
)

// GetLink returns a link object
func GetLink(w http.ResponseWriter, r *http.Request) {
	var link models.Link
	link.Slug = chi.URLParam(r, "slug")
	err := link.Search()
	if err != nil {
		log.Println("Search error\n", err)
		w.WriteHeader(http.StatusNotFound)
	}
	if link.URL != "" {
		httpRedirect(w, r, link.URL)
		visit := *&models.Visit{LinkID: link.ID, Time: time.Now().UTC().Unix()}
		err := visit.Create()
		if err != nil {
			log.Println("Couldn't create visit obj\n", err)
			return
		}
		err = link.Update()
		if err != nil {
			log.Println("Couldnt update link\n", err)
		}
		return
	}
	w.WriteHeader(404)
}

func GetLinksByUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id != "" {
		links, err := models.GetLinksByUser(id)
		if err != nil {
			log.Println("Error getting user links\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		encoded, err := json.Marshal(links)
		if err != nil {
			// TODO
		}
		respondJSON(w, encoded, http.StatusOK)
		return
	}
}

// GetLinkAttribute returns given attribute
func GetLinkAttribute(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	link, err := getLink(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(404)
	}
	var data []byte
	switch chi.URLParam(r, "attr") {
	case "slug":
		response["data"] = link.Slug
	case "url":
		response["data"] = link.URL
	default:
		w.WriteHeader(403)
		return
	}
	respondJSON(w, data, http.StatusAccepted)
}

// PutLinkAttribute updates a given attribute
func PutLinkAttribute(w http.ResponseWriter, r *http.Request) {
	link, err := getLink(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(404)
		return
	}

	uid := ""

	if claims, ok := r.Context().Value(middleware.AuthCtx{}).(jwt.MapClaims); ok {
		if userID, ok := claims["UserID"].(string); ok {
			uid = userID
		}
	} else {
		// user not allowed to put link attributes
		w.WriteHeader(http.StatusUnauthorized)
	}

	data := make(map[string]interface{})
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Couldn't ready json body\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(bd, &data); err != nil {
		log.Println("Couldn't unmarhsal json\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch path.Base(r.URL.Path) {
	case "url":
		err = link.Put(uid, "url", data["url"])
	}
	if err != nil {
		log.Println(err)
	}
	encoded, err := json.Marshal(link)
	if err != nil {
		log.Println("Error updating link")
		w.WriteHeader(http.StatusInternalServerError)
	}
	respondJSON(w, encoded, http.StatusOK)
}

// CreateLink creates a new link
func CreateLink(w http.ResponseWriter, r *http.Request) {
	var link models.Link
	bd, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(bd, &link)
	if err != nil {
		log.Println("malformed create query\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	if claims, ok := r.Context().Value(middleware.AuthCtx{}).(jwt.MapClaims); ok {
		if userID, ok := claims["UserID"].(string); ok {
			log.Println(claims)
			link.OwnerID = userID
		} else {
			link.OwnerID = ""
		}
	}

	if link.OwnerID == "" || link.Slug == "" {
		link.Slug = ""
		chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-"
		for i := 0; i < 6; i++ {
			link.Slug = link.Slug + string(chars[rand.Intn(62)])
			log.Println(link.Slug)
		}
	}

	link.Lease = time.Now().Add(time.Hour * 24 * 30).Unix()
	err = link.Create()
	if err != nil {
		http.Error(w, http.StatusText(409), 409)
		log.Println(err)
		return
	}
	encoded, err := json.Marshal(link)
	if err != nil {
		log.Println("Couldn't parse link as json\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Println(string(encoded))
	respondJSON(w, encoded, http.StatusCreated)
}

func DeleteLink(w http.ResponseWriter, r *http.Request) {
	link := &models.Link{ID: chi.URLParam(r, "id")}
	uid := ""
	if claims, ok := r.Context().Value(middleware.AuthCtx{}).(jwt.MapClaims); ok {
		uid = claims["UserID"].(string)
	}

	err := link.Delete(uid)
	if err != nil {
		log.Println("error deleting link\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// SearchLink returns link object from given slug
func SearchLink(w http.ResponseWriter, r *http.Request) {
	var link *models.Link
	link.Slug = chi.URLParam(r, "query")
	if err := link.Search(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(404)
			log.Println("Couldn't find slug\n", err)
			return
		}
		w.WriteHeader(502)
		return
	}
	encoded, err := json.Marshal(link)
	if err != nil {
		// do something
		w.WriteHeader(502)
		log.Println("Couldn't encode json response")
		return
	}
	respondJSON(w, encoded, http.StatusAccepted)
}

func getLink(id string) (*models.Link, error) {
	link := &models.Link{ID: id}
	if err := link.Get(); err != nil {
		log.Println("error getting link\n", err)
		return nil, err
	}
	return link, nil
}

func httpRedirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, 302)
}
