package api

import (
	"NestedSetsStorage/configs"
	"NestedSetsStorage/treestorage"
	"encoding/json"
	"errors"
	"net/http"
)

// Server starts storage
type Server struct {
	Config      *configs.Config
	Storage     *treestorage.NestedSetsStorage
	apiKeyCache string
}

// Start starts the api server
func (s *Server) Start() error {
	http.HandleFunc("/", s.startFace())
	http.HandleFunc("/all", s.all())
	http.HandleFunc("/children", s.children())
	http.HandleFunc("/parents", s.parents())
	http.HandleFunc("/add", s.add())
	http.HandleFunc("/move", s.move())
	http.HandleFunc("/remove", s.remove())
	http.HandleFunc("/rename", s.rename())

	s.apiKeyCache = s.Config.APIKey
	return http.ListenAndServe(s.Config.APIPort, nil)
}

func (s *Server) startFace() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Nested sets storage started\n"))
	}
}

func (s *Server) all() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.FormValue("key")
		err := s.checkKey(key)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		data, err := s.Storage.GetWholeTree()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		j, _ := json.Marshal(data)
		w.Write([]byte(j))
	}
}

func (s *Server) parents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.FormValue("key")
		err := s.checkKey(key)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		data, err := s.Storage.GetParents(r.FormValue("name"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		j, _ := json.Marshal(data)
		w.Write([]byte(j))
	}
}

func (s *Server) children() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.FormValue("key")
		err := s.checkKey(key)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		data, err := s.Storage.GetChildren(r.FormValue("name"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		j, _ := json.Marshal(data)
		w.Write([]byte(j))
	}
}

func (s *Server) add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.FormValue("key")
		err := s.checkKey(key)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		err = s.Storage.AddNode(r.FormValue("name"), r.FormValue("parent"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func (s *Server) move() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.FormValue("key")
		err := s.checkKey(key)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		err = s.Storage.MoveNode(r.FormValue("name"), r.FormValue("parent"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func (s *Server) remove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.FormValue("key")
		err := s.checkKey(key)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		err = s.Storage.RemoveNode(r.FormValue("name"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func (s *Server) rename() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.FormValue("key")
		err := s.checkKey(key)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		err = s.Storage.RenameNode(r.FormValue("name"), r.FormValue("new_name"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func (s *Server) checkKey(key string) error {
	if key != s.apiKeyCache {
		return errors.New("access denied")
	}
	return nil
}
