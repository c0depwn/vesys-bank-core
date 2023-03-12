package http_api

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type API struct {
	c *Controller
}

func NewAPI(c *Controller) *API {
	return &API{c: c}
}

func (api API) Listen(addr string) error {

	r := mux.NewRouter()

	// logger middleware
	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			log.Printf("[INFO] request received: method=%v path=%v", req.Method, req.URL.String())
			handler.ServeHTTP(res, req)
			log.Printf("[INFO] request handled")
		})
	})

	// routes
	bankRouter := r.PathPrefix("/bank").Subrouter()
	bankRouter.Path("/transfer").Methods(http.MethodPost).HandlerFunc(api.c.HandleTransfer)

	accountsRouter := bankRouter.PathPrefix("/accounts").Subrouter()
	accountsRouter.Path("").Methods(http.MethodGet).HandlerFunc(api.c.HandleGetAccountNumbers)
	accountsRouter.Path("").Methods(http.MethodPost).HandlerFunc(api.c.HandleCreateAccount)

	accountRouter := accountsRouter.PathPrefix("/{id}").Subrouter()
	accountRouter.Path("").Methods(http.MethodGet).HandlerFunc(api.c.HandleGetAccount)
	accountRouter.Path("").Methods(http.MethodDelete).HandlerFunc(api.c.HandleCloseAccount)
	accountRouter.Path("/withdraw").Methods(http.MethodPost).HandlerFunc(api.c.HandleWithdraw)
	accountRouter.Path("/deposit").Methods(http.MethodPost).HandlerFunc(api.c.HandleDeposit)

	log.Println("[INFO] listening on " + addr)
	return http.ListenAndServe(addr, r)
}
