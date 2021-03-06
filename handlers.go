package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusNotFound, "Not found")
}

func (a *App) getCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := getCategoriesFromDB(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

func (a *App) showCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["category"]

	boards, err := showCategoryFromDB(a.DB, name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category")
		return
	}

	respondWithJSON(w, http.StatusOK, boards)
}

func (a *App) showBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["board"]

	posts, err := showBoardFromDB(a.DB, name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid board")
		return
	}

	respondWithJSON(w, http.StatusOK, posts)
}

func (a *App) showThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	thread, err := strconv.Atoi(vars["thread"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid thread ID")
		return
	}

	posts, err := showThreadFromDB(a.DB, board, thread)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid board")
		return
	}

	respondWithJSON(w, http.StatusOK, posts)
}

func (a *App) getBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := getBoardsFromDB(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, boards)
}

func (a *App) postThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	r.ParseMultipartForm(32 << 20)
	comment := r.FormValue("comment")
	if comment == "" {
		respondWithError(w, http.StatusBadRequest, "Comment cannot be empty")
		return
	}

	post_num, err := makeThreadInDB(a.DB, board, comment)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error creating thread")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int{"post_num": post_num})
}

func (a *App) postReply(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	thread, err := strconv.Atoi(vars["thread"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid thread ID")
		return
	}

	is_op, err := checkIsOp(a.DB, board, thread)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error")
		return
	}

	if !is_op {
		respondWithError(w, http.StatusBadRequest, "Specified post is not OP")
		return
	}

	r.ParseMultipartForm(32 << 20)
	comment := r.FormValue("comment")
	if comment == "" {
		respondWithError(w, http.StatusBadRequest, "Comment cannot be empty")
		return
	}

	post_num, err := makePostInDB(a.DB, board, thread, comment)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error creating post")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int{"post_num": post_num})
}
