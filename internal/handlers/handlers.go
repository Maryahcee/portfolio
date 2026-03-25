package handlers

import (
	"html/template"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

type App struct {
	template  *template.Template
	staticDir string
	store     SubmissionStore
}

type PageData struct {
	Error   string
	Success bool
	Name    string
	Email   string
	Message string
}

type Submission struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func NewApp(templatePath, staticDir, submissionsLog string) (*App, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	store, err := NewSubmissionStore(submissionsLog)
	if err != nil {
		return nil, err
	}

	return &App{
		template:  tmpl,
		staticDir: staticDir,
		store:     store,
	}, nil
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(a.staticDir+"/assets"))))
	mux.Handle("/styles.css", http.FileServer(http.Dir(a.staticDir)))
	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/contact", a.contact)
	return mux
}

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := PageData{
		Success: r.URL.Query().Get("status") == "sent",
	}

	a.render(w, http.StatusOK, data)
}

func (a *App) contact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	data := PageData{
		Name:    strings.TrimSpace(r.FormValue("name")),
		Email:   strings.TrimSpace(r.FormValue("email")),
		Message: strings.TrimSpace(r.FormValue("message")),
	}

	switch {
	case data.Name == "":
		data.Error = "Name is required."
	case data.Email == "":
		data.Error = "Email is required."
	case !validEmail(data.Email):
		data.Error = "Enter a valid email address."
	case data.Message == "":
		data.Error = "Message is required."
	}

	if data.Error != "" {
		a.render(w, http.StatusBadRequest, data)
		return
	}

	if err := a.store.Save(r.Context(), Submission{
		Name:      data.Name,
		Email:     data.Email,
		Message:   data.Message,
		CreatedAt: time.Now().UTC(),
	}); err != nil {
		data.Error = "Your message could not be saved right now. Please try again."
		a.render(w, http.StatusInternalServerError, data)
		return
	}

	http.Redirect(w, r, "/?status=sent", http.StatusSeeOther)
}

func (a *App) render(w http.ResponseWriter, status int, data PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = a.template.Execute(w, data)
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
