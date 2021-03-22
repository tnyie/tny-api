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

	"github.com/tnyie/tny-api/middleware"
	"github.com/tnyie/tny-api/models"
	"github.com/tnyie/tny-api/util"
)

// GetLink returns a link object
func GetLink(w http.ResponseWriter, r *http.Request) {
	var link models.Link
	link.Slug = chi.URLParam(r, "slug")
	err := link.Get()
	if err != nil {
		log.Println("Search error\n", err)
		w.WriteHeader(http.StatusNotFound)
	}

	curr_time := time.Now().Unix()

	if link.URL != "" {
		if link.Lease > curr_time {
			if link.UnlockTime != 0 || link.UnlockTime <= curr_time {
				http.Redirect(w, r, link.URL, http.StatusTemporaryRedirect)
				visit := &models.Visit{LinkID: link.ID}
				err := visit.Create()
				if err != nil {
					log.Println("Couldn't create visit obj\n", err)
					return
				}
				link.Visits += 1
				err = link.Update()
				if err != nil {
					log.Println("Couldnt update link\n", err)
				}
				return
			}
		}
	}
	w.WriteHeader(404)
}

// GetLinksByUser checks for authorized user and returns all links owned by user
func GetLinksByUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, authorized := util.CheckLogin(r, id)

	if !authorized {
		log.Println("User unauthorized to access resource")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	links, err := models.GetLinksByUser(id)
	if err != nil {
		log.Println("Error getting user links\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encoded, err := json.Marshal(links)
	if err != nil {
		log.Println("Error marshalling data\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respondJSON(w, encoded, http.StatusOK)
}

// PutLinkAttribute updates a given attribute
func PutLinkAttribute(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, authorized := util.CheckLogin(r, id)

	if !authorized {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	link, err := getLink(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(404)
		return
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
		err = link.Put(id, "url", data["url"])
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
	if err != nil {
		log.Println("Couldn't read json body when creating link\n", err)
	}
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
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
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

// DeleteLink deletes a given link
func DeleteLink(w http.ResponseWriter, r *http.Request) {
	link := &models.Link{ID: chi.URLParam(r, "id")}
	link.Get()

	userAuth, authorized := util.CheckLogin(r, link.OwnerID)

	if !authorized || link.OwnerID == "" || link.OwnerID != userAuth.UID {
		log.Println("Link's Owner ID : ", link.OwnerID)
		log.Println("User's ID ", userAuth.UID)
		w.WriteHeader(http.StatusUnauthorized)
	}

	err := link.Delete()
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
