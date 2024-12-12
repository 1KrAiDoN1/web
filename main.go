package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	handleRequest()
	// fmt.Println(os.Getenv("MYSQL_PASSWORD"))
	// pswd := os.Getenv("MYSQL_PASSWORD")
	// db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/test_db")
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// Установка данных
	// mysql -u root -p команда для запуска хоста

	// // УДАЛЕНИЕ ЭЛЕМЕНТА В БАЗЕ ДАННЫХ
	// // Выполнение запроса
	// result, err := db.Exec("DELETE FROM `test_db`.`articles` WHERE id = ?", 5)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// rowsAffected1, err := result.RowsAffected()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("Удалено строк: %d\n", rowsAffected1)

	// rowsAffected, err := insert.RowsAffected()
	// if err != nil {
	// 	log.Fatalf("Ошибка при получении количества затронутых строк: %v", err)
	// }
	// fmt.Printf("Затронуто строк: %d\n", rowsAffected)

}

func handleRequest() {
	router := mux.NewRouter()
	router.HandleFunc("/home_page", home_page).Methods("GET")
	router.HandleFunc("/create", create).Methods("GET")
	router.HandleFunc("/save_articles", save_articles).Methods("POST")
	router.HandleFunc("/contacts/", contactInformation).Methods("GET")
	router.HandleFunc("/more_information", more_information).Methods("GET")
	router.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")
	http.Handle("/", router)
	router.HandleFunc("/our_error", ourError)
	http.ListenAndServe(":8080", nil)
}

func show_post(t http.ResponseWriter, s *http.Request) {
	vars := mux.Vars(s)
	tmpl, err := template.ParseFiles("templates/show_post.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/test_db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM `test_db`.`articles` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}
	showPost = Articles{}
	for res.Next() {
		var post Articles
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_Text)
		if err != nil {
			panic(err)
		}
		showPost = post
	}
	tmpl.ExecuteTemplate(t, "show_post", showPost)

}

var posts = []Articles{}
var showPost = Articles{}

func home_page(page http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index1.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(page, err.Error())
	}

	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/test_db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	res, err := db.Query("SELECT * FROM `test_db`.`articles`")
	if err != nil {
		panic(err)
	}
	posts = []Articles{}
	for res.Next() {
		var post Articles
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_Text)
		if err != nil {
			panic(err.Error())
		}
		posts = append(posts, post)
	}
	tmpl.ExecuteTemplate(page, "index", posts)

}

func create(t http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(t, err.Error())
	}
	tmpl.ExecuteTemplate(t, "create", posts)

}

func contactInformation(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/contacts.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmpl.ExecuteTemplate(w, "contacts", nil)

}

func more_information(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/more_information.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	tmpl.ExecuteTemplate(w, "more_information", nil)

}

func ourError(t http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/ourError.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(t, "ourError", nil)
}

func save_articles(h http.ResponseWriter, a *http.Request) {
	title := a.FormValue("title")
	anons := a.FormValue("anons")
	full_text := a.FormValue("full_text")
	if title == "" || anons == "" || full_text == "" {
		http.Redirect(h, a, "/our_error", http.StatusSeeOther)
		return
	} else {
		db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/test_db")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		err = db.Ping()
		if err != nil {
			panic(err)
		}
		// ДОБАВЛЕНИЕ ЭЛЕМЕНТА В БАЗУ ДАННЫХ
		insert, err := db.Exec(fmt.Sprintf("INSERT INTO `test_db`.`articles` (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
		if err != nil {
			log.Fatalf("Ошибка при вставке данных: %v", err)
		}
		http.Redirect(h, a, "/home_page", http.StatusSeeOther)
		fmt.Println(insert.LastInsertId())

	}
}

type User struct {
	Name        string
	Age         uint
	PhoneNumber string
	Email       string
}

type Articles struct {
	Id        uint16
	Title     string
	Anons     string
	Full_Text string
}
