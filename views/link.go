package views

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"

	"github.com/tnyie/tny-api/models"
	"github.com/tnyie/tny-api/util"
)

const (
	LEASE_TIME = time.Hour * 24 * 30 * 6
)

// CreateLink creates a new link
func CreateLink(w http.ResponseWriter, r *http.Request) {
	var link models.Link
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Couldn't read json body when creating link\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bd, &link)
	if err != nil {
		log.Println("malformed create query\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, authenticated, admin := util.CheckLogin(r, "")
	if !authenticated && !admin {
		link.OwnerID = ""
	} else {
		link.OwnerID = user.UID
	}

	if link.OwnerID == "" || link.Slug == "" {
	tryNew:
		link.Slug = petname.Generate(2, "-")
		// if link exists, and no error, try a new petname
		if err = link.GetBySlug(); err == nil {
			goto tryNew
		}
	}

	// set to zero, so GORM can set it to current time
	link.UpdatedAt, link.CreatedAt = 0, 0

	link.Lease = time.Now().Add(LEASE_TIME).Unix()

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

	log.Println("New link created by UID: ", link.OwnerID, " with Slug: ", link.Slug)
	respondJSON(w, encoded, http.StatusCreated)
}

// RedirectSlug takes the link's slug as a parameter, and redirects request
func RedirectSlug(w http.ResponseWriter, r *http.Request) {
	var link models.Link
	link.Slug = chi.URLParam(r, "slug")
	err := link.GetBySlug()
	if err != nil {
		log.Println("Search error\n", err)
		w.WriteHeader(http.StatusNotFound)
	}

	curr_time := time.Now().Unix()

	if link.URL != "" {
		if link.Lease == 0 || link.Lease > curr_time {
			if link.UnlockTime < curr_time {
				if link.Password != "" {
					log.Println("Link requires password, passing to other handler")
					http.Redirect(w, r,
						"https://"+viper.GetString("tny.ui.url")+"/redirect/"+link.Slug,
						http.StatusTemporaryRedirect,
					)
					return
				}
				log.Println("Redirecting to", link.URL)
				http.Redirect(w, r, link.URL, http.StatusTemporaryRedirect)
				setVisit(&link)
				return
			} else {
				log.Println("Link not unlocked yet")
				log.Println(link.UnlockTime, " < ", curr_time)
				w.WriteHeader(404)
				return
			}
		} else {
			log.Println("Link lease expired")
			w.WriteHeader(404)
		}
	} else {
		log.Println("Link corrupted")
		w.WriteHeader(404)
	}
}

// GetLink returns a link object
func GetLinkFromSlug(w http.ResponseWriter, r *http.Request) {
	var link models.Link
	link.Slug = chi.URLParam(r, "slug")

	err := link.GetBySlug()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Error getting link from id\n", err)
		return
	}

	_, authorized, admin := util.CheckLogin(r, link.OwnerID)
	if !authorized && !admin {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("user not authorized to access link object")
		return
	}

	encoded, err := json.Marshal(&link)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to marshal link to json\n", err)
		return
	}

	respondJSON(w, encoded, http.StatusOK)
}

func GetLink(w http.ResponseWriter, r *http.Request) {
	var link models.Link
	link.ID = chi.URLParam(r, "id")

	err := link.Get()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Error getting link from id\n", err)
		return
	}
	_, authorized, admin := util.CheckLogin(r, link.OwnerID)
	if !authorized && !admin {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("user not authorized to access link object")
		return
	}

	encoded, err := json.Marshal(&link)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to marshal link to json\n", err)
		return
	}

	respondJSON(w, encoded, http.StatusOK)
}

// GetAuthenticatedLink returns the
func GetAuthenticatedLink(w http.ResponseWriter, r *http.Request) {
	var link models.Link
	link.Slug = chi.URLParam(r, "slug")
	err := link.GetBySlug()
	if err != nil {
		log.Println("Search error\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	curr_time := time.Now().Unix()

	if link.URL != "" {
		if link.Lease == 0 || link.Lease > curr_time {
			if link.UnlockTime < curr_time {
				if link.Password != "" {
					jsonMap := make(map[string]interface{})
					bd, err := ioutil.ReadAll(r.Body)

					if err != nil {
						log.Println("Error reading json body")
						return
					}

					err = json.Unmarshal(bd, &jsonMap)
					if err != nil {
						log.Println("Couldn't parse json body")
						return
					}

					if jsonMap["password"] != "" {
						if jsonMap["password"] == link.Password {
							encoded, err := json.Marshal(&models.GenericResponse{Data: link.URL})
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								log.Println("Couldn't generate response")
								return
							}
							respondJSON(w, encoded, http.StatusOK)
							setVisit(&link)
							return
						}
						w.WriteHeader(http.StatusUnauthorized)
						log.Println("incorrect password")
						return
					}
					w.WriteHeader(http.StatusUnauthorized)
					log.Println("no password provided for authenticated link")
					return
				}
				w.WriteHeader(http.StatusUnauthorized)
				log.Println("Link has no password")
				return
			}
			log.Println("link not unlocked yet")
			w.WriteHeader(http.StatusFound)
			return
		}
		log.Println("Link lease exceeded")
	}
	log.Println("link is invalid")
	// if link is missing url or has expired, 404
	w.WriteHeader(http.StatusNotFound)
}

// setVisit updates visits for a given link object in database
func setVisit(link *models.Link) {
	visit := &models.Visit{LinkID: link.ID}
	err := visit.Create()
	if err != nil {
		log.Println("Couldn't create visit obj\n", err)
	}

	link.Visits += 1
	err = link.Update()
	if err != nil {
		log.Println("Couldnt update link\n", err)
	}
}

// GetLinksByUser checks for authorized user and returns all links owned by user
func GetLinksByUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, authorized, admin := util.CheckLogin(r, id)
	if !authorized && !admin {
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

// UpdateLinkLease updates link lease time if user is
// authenticated, link is expired, and user owns resource
func UpdateLinkLease(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Println("Updating link lease for link id ", id)

	link := &models.Link{ID: id}

	err := link.Get()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, authorized, admin := util.CheckLogin(r, link.OwnerID)
	if !authorized && !admin {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("User is not authorized to modify link lease")
		return
	}

	err = link.Put(link.OwnerID, "lease", time.Now().Add(time.Hour*24*30).Unix())
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		log.Println("Error updating link lease\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// PutLinkAttribute updates a given attribute
func PutLinkAttribute(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	log.Println("Updating link of id ", id)

	link := &models.Link{ID: id}
	err := link.Get()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, authorized, admin := util.CheckLogin(r, link.OwnerID)
	if !authorized && !admin {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Println("Updating link ", link)

	data := make(map[string]interface{})
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Couldn't ready json body\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(bd, &data); err != nil {
		log.Println(string(bd))
		log.Println("Couldn't unmarhsal json\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch path.Base(r.URL.Path) {
	case "url":
		err = link.Put(link.OwnerID, "url", data["url"])
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

// DeleteLink deletes a given link
func DeleteLink(w http.ResponseWriter, r *http.Request) {
	link := &models.Link{ID: chi.URLParam(r, "id")}
	link.Get()

	userAuth, authorized, admin := util.CheckLogin(r, link.OwnerID)

	if admin || authorized && link.OwnerID != "" {
		link.Delete()
		w.WriteHeader(http.StatusAccepted)
		return
	}

	if link.OwnerID == "" || link.OwnerID != userAuth.UID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}
