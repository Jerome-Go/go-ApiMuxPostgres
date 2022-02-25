package main
// Créer une application API REST à l'aide de la base de données Golang et PostgreSQL
//
//https://www.section.io/engineering-education/build-a-rest-api-application-using-golang-and-postgresql-database/
//
// https://linuxize.com/post/curl-post-request/   -> POST
//
// ----------TERMINAL---------------
//
// ----- connexion a tb
// psql -U postgres -d jerome_movie       -> connexion a la dab
// 
// ----- creation table
// jerome_movie=# CREATE TABLE movies(
// jerome_movie(# id SERIAL PRIMARY KEY, 
// jerome_movie(# movieID TEXT NOT NULL,
// jerome_movie(# movieName TEXT NOT NULL);
// CREATE TABLE
//
// ----- insert enr
// jerome_movie=# INSERT INTO movies(movieID, movieName)
// jerome_movie-# VALUES('1', 'movie3'), 
// jerome_movie-# ('2', 'movie2'), 
// jerome_movie-# ('3', 'movie1');
// INSERT 0 3
//
// ----- POST
// curl -X POST -d 'movieid=1&moviename=movie01' localhost:8000/movies/ -> Post 
// curl --data "@mydata"                     -> post fichier
// curl --data "{'count':'5'}" -H "Content-Type:application/json" -> post json
// 
// select * from movie where id=6;
// select * from movie;
//
//
//
// -------- PACKAGES -------
// log                    -> pour consigner les erreurs
// encoding/json          -> gestion des données json
// database/sql           -> gestion de la communication db sur SQL.
// net/http               -> gerer les requetes http
// github.com/gorilla/mux -> il aide à implémenter des routeurs de requêtes
//                           et à faire correspondre chaque requête entrante
//                           avec son gestionnaire correspondant.
// github.com/lib/pq      -> pilote Go PostgreSQL

import (
	"encoding/json"
	"fmt"
	//"log"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//
// coller json
type Movie struct {
	Id int `json:"ID"`
	MovieID string `json:"filmID"`
	MovieName string `json:"filmName"`
}
//
// reponse json
type reponseJson struct{
	Message string `json:"message"`
}
//
// Authentification a la db
const(
	host = "localhost"
	port = 5432
	user = "postgres"
	dbname = "jerome_movie"
)
//
// gestion des erreures
func gestionErr(err error)  {
	if err != nil {
		panic(err)
	}
}
//
// print message
func printMessage(message string)  {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")	
}
//
//
// --------1 DEBUT connexion db 
func connexionDB() *sql.DB {
	//
	// authentification string
	connexionDB := fmt.Sprintf("host=%v port=%d user=%v dbname=%v sslmode=disable", host, port, user, dbname)
	//
	// configuration ouverture
	db ,err := sql.Open("postgres", connexionDB)
	if err != nil {
		gestionErr(err)
	}
	//
	// connexion a la db
	err = db.Ping()
	if err != nil{
		gestionErr(err)
	}
	//
	fmt.Println("connexion a la db ok")
	//
	return db
}
// --------1 FIN connexion db 
//
//
//
// --------2 DEBUT incerer enr 
var id int
func incertEnr(db *sql.DB)  {
	//
	// var pour ligne de commande
	incertEnr := `INSERT INTO movies(movieID, movieName)
	VALUES ($1, $2)
	RETURNING id;`
	//
	// demande et colle le resultat dans un var
	row := db.QueryRow(incertEnr,"12", "movie12" )
	//
	// traduit en go et colle dans un var
	err := row.Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("incerte enr avec id= ", id)
}
// --------2 FIN incerer enr
//
//
// --------3 DEBUT lire enr avec struc
func selectNewEnr(db *sql.DB){
	//
	// ligne de commande
	selectNewEnr := `SELECT * FROM movies WHERE id=$1;`
	//
	var newMovie Movie
	//
	// demande et colle dans un var
	row := db.QueryRow(selectNewEnr, id)
	//
	// traduire en go et coller dans un var
	err := row.Scan(&newMovie.Id, &newMovie.MovieID, &newMovie.MovieName)
	if err != nil{
		panic(err)
	}
	//
	// print resultat
	fmt.Println("newMovie = ", newMovie)
}
// --------3 fin lire enr avec struc
//
//
// --------4 DEBUT GetMovies
func GetMovies(w http.ResponseWriter, r *http.Request)  {
	//
	// 1 connexion db
	db := connexionDB()
	//
	// ligne de commande
	getMovies := `SELECT * FROM movies;`
	//
	// demande et coller dans var
	 rows , err := db.Query(getMovies)
	 gestionErr(err)
	//
	var moviesList []Movie
	//
	//faire passer dans un boucle chaque enr.
	for rows.Next() {
		row := Movie{}
		//
		// traduire en go et coller dans var
		err := rows.Scan(&row.Id, &row.MovieID, &row.MovieName)
		gestionErr(err)
		//
		// additioner le resultat à la var list
		moviesList = append(moviesList, row)
	}
	//
	//
	fmt.Println(moviesList)	
	//
	// 1 premiere facon de faire un jsson 	
	//jsonResultat, err := json.Marshal(moviesList)
	//gestionErr(err)
	//fmt.Println(string(jsonResultat))
	 //
	 // 2 deuxième facon de faire un jsson 
	 json.NewEncoder(w).Encode(moviesList)
	 //
	 fmt.Println(moviesList)
}
// --------4 FIN GetMovie
//
//
// --------5 DEBUT incert nouvelle enr
func CreationMovie(w http.ResponseWriter, r *http.Request)  {
	//  curl -X POST -d 'movieid=1&moviename=movie01' localhost:8000/movies/
	// curl --data "movieid=...&moviename=..."  -> request = post
	// la demande r *http.Request -> 1ere request movieid
	//                             -> 2meme request moviename
	movieID01 := r.FormValue("movieid")
    movieName01 := r.FormValue("moviename")
	//
	//
	var reponseJson01 = reponseJson{}
	//
	if movieID01 == "" || movieName01 == ""{
		reponseJson01 = reponseJson{Message: "movieID ou movieName est vide"}
	} else {
	      //
	      // connexion a la db
	      db := connexionDB()
	      //
	      fmt.Println("Insert un enr avec L'id: " + movieID01 + " et le nom: " + movieName01 )
	      //
	      var lastInsertID int
	      // ligne de commande
	      creationMovie01 := `INSERT INTO movies(movieID, movieName)
	      VALUES ($1, $2)
	      RETURNING id;`
	      //
	      // demande et colle dans var
	      row := db.QueryRow(creationMovie01, movieID01, movieName01)
	      //
	      // traduit en go et colle dans un var id
	      err := row.Scan(&lastInsertID )
	      gestionErr(err)
	      //
	      reponseJson01 = reponseJson{Message: "insertion de l'enr ok"}
	      fmt.Println("l'id de l'enr est: ", lastInsertID)
	      }
	//
	// message sur chrome
	 json.NewEncoder(w).Encode(reponseJson01)
}
// --------5 FIN incet nouvelle enr
//
//
//
// --------6 DEBUT DELETE single enr
func DeleteEnr(w http.ResponseWriter, r *http.Request)  {
	//
	// mettre le mux var 'movieid' dans une var= Map[]
	// router.HandleFunc("/movies/{movieid}"
	mapMovieIdSupp := mux.Vars(r)
	fmt.Println("mapMovieIdSupp = ", mapMovieIdSupp)
	//
	// movieIdSupp est un Map[string]string
	 movieIdSupp := mapMovieIdSupp["movieid"]
	 fmt.Println("movieIdSupp = ", movieIdSupp)
	//
	var reponseJson01 = reponseJson{}
	if movieIdSupp == "" {
		reponseJson01 = reponseJson{Message: "aucun enr a supp est indiqué"}
	}else {
	      //
	      // se connecter à la db
	      db := connexionDB()
	      //
	      // auth string
	      authString := `DELETE FROM movies WHERE id=$1;`
	      //
	      // selectionné le ligne
	       _, err := db.Exec(authString, movieIdSupp)
	      gestionErr(err)
	      //
		  // imprimer terminal
	      fmt.Println("supprimé l'enr avec l'id: ", movieIdSupp)
		  //
		  // imprimer dans chrome
		  reponseJson01 = reponseJson{Message: "supprimer l'enr ok"}
	      }
	//
	// imprimer dans chrome
	json.NewEncoder(w).Encode(reponseJson01)
	
}
// --------6 DEBUT DELETE single enr
//
//
//
// --------7 DEBUT DELETE tout les enc
func DeleteToutEnr(w http.ResponseWriter, r *http.Request)  {
	//
	// connexion a la db
	db := connexionDB()
	//
	// ligne de commande
	deleteToutEnr := `DELETE FROM movies;`
	//
	// executer le ligne de commande dans le terminal
	_, err := db.Exec(deleteToutEnr)
	gestionErr(err)	
	//
	// imprime dans chrome
	reponse02 := reponseJson{Message: "tous les enr. sont supprimées"}
	//
	// traduire en json
	json.NewEncoder(w).Encode(reponse02)
}

//
//
//
func main()  {
//
// 1 connexion a la db
//db, err := connexionDB()
//gestionErr(err)
//
// 2 inserer un enr
//incertEnr(db)
//
//initialiser le mux router
router := mux.NewRouter()
//
// 4. get tout les enr
// curl localhost:8000/movies/
router.HandleFunc("/movies/", GetMovies).Methods("GET")
//
// 5 creer un enr
// curl -X POST -d 'movieid=1&moviename=movie01' localhost:8000/movies/
router.HandleFunc("/movies/", CreationMovie).Methods("POST")
//
// 6 delete un enr
// curl -X DELETE localhost:8000/movies/21
router.HandleFunc("/movies/{movieid}", DeleteEnr).Methods("DELETE")
//
// 7 delete tout les enr
// curl -X DELETE localhost:8000/movies/
router.HandleFunc("/movies/", DeleteToutEnr).Methods("DELETE")
//
// ouvrir le serveur
http.ListenAndServe(":8000", router)	
}