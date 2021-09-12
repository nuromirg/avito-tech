package server

import (
	"avito_task/config"
	"avito_task/internal/model"
	"avito_task/internal/store"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type Resp struct {
	Message interface{}
	Error string
}

type Server struct {
	config  *config.Config
	logger  *logrus.Logger
	router  *mux.Router
	store   *store.Storage
}

func New(_config *config.Config) *Server {
	return &Server{
		config: _config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	 if err := s.configureLogger(); err != nil {
		return err
	 }

	 err := s.configureStore()
	 if err != nil {
		return err
	 }
	 s.configureRouter()
	 s.logger.Info("Starting httpserver")
	 return http.ListenAndServe(s.config.Bindaddr, s.router)
}

func (s *Server) configureLogger() error {
	lvl, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(lvl)
	return nil
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/users/{id}", s.handleUsers())				//	get existing user by id

	s.router.HandleFunc("/users/", s.handleUsers())					//	empty query parameter line means
																			//	that you want to create new user
	s.router.HandleFunc("/transaction/", s.handleTransaction,
		).Queries("from", "{from}", "to", "{to}", "sum", "{sum}",	//	example of transaction handler
		).Methods("GET")											//	http://localhost:8081/transaction/?from=12&to=13&sum=100000
																			//	URL that transfer 10000.00RUB from id 12 to id 13
	s.router.HandleFunc("/transaction/", s.handleTransaction,
		).Queries("to", "{to}", "sum", "{sum}",
		).Methods("GET")
}


func (s *Server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}

func (s *Server) configureStore() error {
	st := store.CreateStorage(*s.config)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}

func (s *Server) handleUsers() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.handlerGetUser(w, r)
		} else if r.Method == http.MethodPost {
			s.handlerAddUser(w, r)
		}
	}
}

func (s *Server) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	var response Resp
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		logrus.Error(err)
		response.Error = err.Error()
		responseJson, _ := json.Marshal(response)
		w.Write(responseJson)
		return
	}
	u, err := s.store.User().FindById(id)
	if err != nil {
		logrus.Error(err)
		response.Error = err.Error()
		responseJson, _ := json.Marshal(response)
		w.Write(responseJson)
		return
	}
	response.Message = u
	responseJson, _ := json.Marshal(response)
	w.Write(responseJson)
	return
}

func (s *Server) handleTransaction(w http.ResponseWriter, r *http.Request) {
	var response Resp
	params := mux.Vars(r)
	toId, _ := strconv.Atoi(params["to"])
	sum, _ := strconv.ParseInt(params["sum"], 10, 64)
	fromId, _ := strconv.Atoi(params["from"])
	if fromId == 0 {
		u, err := s.store.User().ChangeFunds(toId, sum)
		if err != nil {
			logrus.Error(err)
			response.Error = err.Error()
			responseJson, _ := json.Marshal(response)
			w.Write(responseJson)
			return
		}
		response.Message = u
		responseJson, _ := json.Marshal(response)
		w.Write(responseJson)
		return
	}
	err := s.store.User().TransactionFunds(fromId, toId, sum)
	if err != nil {
		logrus.Error(err)
		response.Error = err.Error()
		responseJson, _ := json.Marshal(response)
		w.Write(responseJson)
		return
	}
	balance := make([]int64, 2)
	balance[0], _ = s.store.User().GetBalanceById(fromId)
	balance[1], _ = s.store.User().GetBalanceById(toId)
	str := fmt.Sprintf("from %v, balance: %v; to: %v, balance: %v",
		fromId, balance[0], toId, balance[1])
	response.Message = str
	responseJson, _ := json.Marshal(response)
	w.Write(responseJson)
}


func (s *Server) handlerAddUser(w http.ResponseWriter, r *http.Request) {
	var (
		response 	Resp
		u			*model.User
	)
	u, err := s.store.User().CreateUser()
	if err != nil {
		logrus.Error(err)
		response.Error = err.Error()
		responseJson, _ := json.Marshal(response)
		w.Write(responseJson)
		return
	}
	response.Message = u
	responseJson, _ := json.Marshal(response)
	w.Write(responseJson)
}