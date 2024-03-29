package main

import (
	"fmt"
	config "hackaton/app/config"
	cookie "hackaton/app/cookies"
	"net/http"
	"os"
	"text/template"
)

var data = make(map[string]interface{})

// FilterEvents filtre les événements en fonction du titre et de la description
func FilterEvents(events []config.Event, title, description string) []config.Event {
	var filteredEvents []config.Event

	for _, event := range events {
		if title != "" && event.Title != title {
			continue // Si le titre ne correspond pas, passez à l'événement suivant
		}
		if description != "" && event.Description != description {
			continue // Si la description ne correspond pas, passez à l'événement suivant
		}

		// Si l'événement correspond aux critères de recherche, ajoutez-le à la liste filtrée
		filteredEvents = append(filteredEvents, event)
	}

	return filteredEvents
}

func main() {
	data["user"] = ""

	config.CreateDB()
	db, err := config.GetDB()
	if err != nil {
		panic(err)
	}
	defer db.CloseDB() // Assurez-vous que la connexion se ferme lorsque la fonction main() se termine

	db.CreateRoleTable()
	db.CreateUserTable()
	db.CreateEventTable()
	db.CreateParticipantsTable()
	db.CreateCommentTable()

	fmt.Println("Please connect to http://localhost:8000")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))) // Join Assets Directory to the server
	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/CreateEvent", CreateEvent)
	http.HandleFunc("/Event", Event)
	http.HandleFunc("/profil", profil)
	http.HandleFunc("/admin", admin)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Port par défaut si la variable d'environnement PORT n'est pas définie
	}

	// Initialisez votre application et vos routes ici

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}

// Le reste du code reste inchangé

func index(w http.ResponseWriter, r *http.Request) {
	config.Api()
	data_index := make(map[string]interface{})
	cookie.SetDataToSend(w, r, data_index, data, false, "")

	input_search := r.FormValue("search")

	config.DisplayEvents(data_index, input_search)

	t := template.New("index-template")
	t = template.Must(t.ParseFiles("public/index.html", "public/header.html", "public/head.html"))
	t.ExecuteTemplate(w, "index", data_index)
}

func register(w http.ResponseWriter, r *http.Request) {
	data_SignUp := make(map[string]interface{})
	cookie.SetDataToSend(w, r, data_SignUp, data, false, "")

	// Redirect if the user is already logged in
	if data_SignUp["cookieExist"] == true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	input_username := r.FormValue("username")
	input_email := r.FormValue("email")
	input_password := r.FormValue("password")

	if input_email != "" || input_username != "" || input_password != "" {
		config.AddUser(input_username, input_email, input_password, data_SignUp) // ./config/AddUser.go
		http.Redirect(w, r, "/login?user="+input_username, http.StatusSeeOther)
	}

	t := template.New("register-template")
	t = template.Must(t.ParseFiles("public/register.html", "public/header.html", "public/head.html"))
	t.ExecuteTemplate(w, "register", data_SignUp)
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	data_event := make(map[string]interface{})
	cookie.SetDataToSend(w, r, data_event, data, false, "")

	// Redirect si l'utilisateur n'est pas connecté
	if data_event["cookieExist"] != true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	input_title := r.FormValue("title")
	input_description := r.FormValue("description")
	input_date := r.FormValue("date")
	if input_title != "" || input_description != "" || input_date != "" {
		config.AddEvent(input_title, input_description, input_date, data_event)
	}

	t := template.New("CreateEvent-template")
	t = template.Must(t.ParseFiles("public/CreateEvent.html", "public/header.html", "public/head.html"))
	t.ExecuteTemplate(w, "CreateEvent", data_event)
}

func Event(w http.ResponseWriter, r *http.Request) {
	data_PageEvent := make(map[string]interface{})
	cookie.SetDataToSend(w, r, data_PageEvent, data, false, "")

	id_event := r.FormValue("Id")
	input_commentaire := r.FormValue("commentaire")
	isAskingToParticipate := r.FormValue("participate") == "true"
	userIdAskingToParticipate := r.FormValue("userID_toAccept")
	note := r.FormValue("note")
	remove_idParticipant := r.FormValue("userID_toRemove")

	// récupère l'évenement
	config.GetEvent(data_PageEvent, id_event)

	// Redirection si l'événement n'existe pas
	if data_PageEvent["event"] == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// récupère les commentaires de l'événement
	db, _ := config.GetDB() // Ignorer l'erreur ici car nous l'avons déjà gérée lors de l'initialisation du serveur
	config.GetComments(db, id_event, data_PageEvent)

	// nouveau commentaire
	if input_commentaire != "" {
		config.AddCommentOnEvent(id_event, input_commentaire, data_PageEvent)
		http.Redirect(w, r, "/Event?Id="+id_event, http.StatusSeeOther)
	}

	// demande de participation
	if isAskingToParticipate {
		config.AddParticipant(id_event, data_PageEvent)
		http.Redirect(w, r, "/Event?Id="+id_event, http.StatusSeeOther)
	}

	// récupère toutes les demandes de participations
	config.GetParticipants(db, id_event, data_PageEvent)

	// Accepte une demande de participation
	if userIdAskingToParticipate != "" {
		config.AcceptParticipation(db, id_event, userIdAskingToParticipate)
		http.Redirect(w, r, "/Event?Id="+id_event, http.StatusSeeOther)
	}

	// retire le participant
	if remove_idParticipant != "" {
		config.RemoveParticipant(db, id_event, remove_idParticipant)
		http.Redirect(w, r, "/Event?Id="+id_event, http.StatusSeeOther)
	}

	// vérifie si l'utilisateur demande déjà la participation
	data_PageEvent["hasRequested"] = config.HasRequestedParticipation(db, id_event, data_PageEvent)

	// vérifie si l'utilisateur connecté est déjà en participant
	data_PageEvent["isAccepted"] = config.IsParticipant(db, id_event, data_PageEvent)

	// modifie la note de l'événement
	if note != "" {
		config.UpdateNote(db, data_PageEvent, id_event, note)
		http.Redirect(w, r, "/Event?Id="+id_event, http.StatusSeeOther)
	}

	t := template.New("PageEvent-template")
	t = template.Must(t.ParseFiles("public/PageEvent.html", "public/header.html", "public/head.html"))
	t.ExecuteTemplate(w, "PageEvent", data_PageEvent)
}

func login(w http.ResponseWriter, r *http.Request) {
	create_cookie := false

	// Initialize the data that will be send to html
	data_logIn := make(map[string]interface{})
	cookie.SetDataToSend(w, r, data_logIn, data, false, "")

	// Redirect if the user is already logged in
	if data_logIn["cookieExist"] == true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// Get user input to log in
	username := r.FormValue("username")
	password := r.FormValue("password")

	create_cookie = cookie.SearchUserToLog(data_logIn, username, password, data) // ./config/user.go

	// Créé un cookie si user bien authentifié
	if create_cookie {
		cookie.CreateCookie(w, r) // ./cookies/createCookie.go
		data_logIn["wrongPassword"] = false
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	data_logIn["user_registered"] = r.FormValue(("user"))

	t := template.New("login-template")
	t = template.Must(t.ParseFiles("public/login.html", "public/header.html", "public/head.html"))
	t.ExecuteTemplate(w, "login", data_logIn)
}

func logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	cookie.DeleteCookie(w)

	// reset the login from FB
	http.Redirect(w, r, "/", http.StatusSeeOther)
	delete(data, "user")
}

func profil(w http.ResponseWriter, r *http.Request) {
	data_profil := make(map[string]interface{})
	cookie.SetDataToSend(w, r, data_profil, data, false, "")

	userID := r.FormValue("id")
	db, err := config.GetDB()
	if err != nil {
		panic(err)
	}
	config.GetUser(db, data_profil, userID)

	t := template.New("profil-template")
	t = template.Must(t.ParseFiles("public/profil.html", "public/header.html", "public/head.html"))
	t.ExecuteTemplate(w, "profil", data_profil)
}

func admin(w http.ResponseWriter, r *http.Request) {
	data_admin := make(map[string]interface{})
	cookie.SetDataToSend(w, r, data_admin, data, false, "")

	if data_admin["role"] != "ADMIN" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	// Get all users
	dbUsers, err := config.GetDB()
	if err != nil {
		panic(err)
	}
	data_admin["allUsers"] = config.GetAllUsers(dbUsers)

	// Get all events
	dbEvents, err := config.GetDB()
	if err != nil {
		panic(err)
	}
	data_admin["allEvents"], err = config.GetAllEvents(dbEvents)
	if err != nil {
		panic(err)
	}

	// Get all comments
	dbComments, err := config.GetDB()
	if err != nil {
		panic(err)
	}
	data_admin["allComments"], err = config.GetAllComments(dbComments)
	if err != nil {
		panic(err)
	}

	// delete User
	user_toDelete := r.FormValue("delUserId")
	if user_toDelete != "" {
		db, err := config.GetDB()
		if err != nil {
			panic(err)
		}
		config.DeleteUser(db, user_toDelete)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}

	// delete Event
	event_toDelete := r.FormValue("delEventId")
	if event_toDelete != "" {
		db, err := config.GetDB()
		if err != nil {
			panic(err)
		}
		config.DeleteEvent(db, event_toDelete)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}

	// make a user an Admin
	user_toAdmin := r.FormValue("makeAdmin")
	if user_toAdmin != "" {
		db, err := config.GetDB()
		if err != nil {
			panic(err)
		}
		config.MakeAdmin(db, user_toAdmin)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}

	// delete Comment
	comment_toDelete := r.FormValue("delCommentId")
	if comment_toDelete != "" {
		db, err := config.GetDB()
		if err != nil {
			panic(err)
		}
		config.DeleteComment(db, comment_toDelete)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}

	t := template.New("admin-template")
	t = template.Must(t.ParseFiles("public/admin.html", "public/header.html", "public/head.html"))
	t.ExecuteTemplate(w, "admin", data_admin)
}
