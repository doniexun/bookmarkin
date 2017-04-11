package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strconv"
)

type BookmarkHandlers struct {
	dbConnection *DBConnection
}

func (bHandlers *BookmarkHandlers) CreateBookmarkGroup(w http.ResponseWriter, r *http.Request) {
	groupName := r.FormValue("group_name")
	
	if groupName != "" {

		bookmark := &Bookmark{}

		result := bookmark.CreateNewGroup(bHandlers.dbConnection, groupName)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(result); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

func (bHandlers *BookmarkHandlers) UpdateBookmarkGroup(w http.ResponseWriter, r *http.Request) {
	groupName := r.FormValue("group_name")
	groupID := r.FormValue("group_id")
	
	if groupName != "" {

		bookmark := &Bookmark{}

		result := bookmark.UpdateExistingGroup(bHandlers.dbConnection, groupID, groupName)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(result); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

func (bHandlers *BookmarkHandlers) DeleteBookmarkGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.FormValue("group_id")
	forceDelete := r.FormValue("force")

	bForceDelete, err := strconv.ParseBool(forceDelete)
	if err != nil {
		bForceDelete = false
	}
	
	if groupID != "" {

		bookmark := &Bookmark{}

		result := bookmark.ListAllBookmarksInGroup(bHandlers.dbConnection, groupID)

		if result != nil && len(result) > 0 {
			if bForceDelete == true {
				result := bookmark.DeleteBookmarksInGroup(bHandlers.dbConnection, groupID)

				if result {
					
					result := bookmark.DeleteBookmarkGroupByID(bHandlers.dbConnection, groupID)

					if result {
						w.Header().Set("Content-Type", "application/json; charset=UTF-8")
						w.WriteHeader(http.StatusOK)

						if err := json.NewEncoder(w).Encode(result); err != nil {
							panic(err)
						}
						return
					} else {
						w.Header().Set("Content-Type", "application/json; charset=UTF-8")
						w.WriteHeader(http.StatusNotFound)
						if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Error deleting group."}); err != nil {
							panic(err)
						}
					}
				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusNotFound)
					if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Error deleting bookmarks from group."}); err != nil {
						panic(err)
					}
				}
			} else {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusNotFound)
				if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Bookmarks exists in the group."}); err != nil {
					panic(err)
				}
			}
		} else {
			result := bookmark.DeleteBookmarkGroupByID(bHandlers.dbConnection, groupID)

			if result {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)

				if err := json.NewEncoder(w).Encode(result); err != nil {
					panic(err)
				}
				return
			} else {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusNotFound)
				if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Error deleting group."}); err != nil {
					panic(err)
				}
			}
		}

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

func (bHandlers *BookmarkHandlers) ListBookmarkGroups(w http.ResponseWriter, r *http.Request) {
	bookmark := &Bookmark{}

	result := bookmark.ListAllGroups(bHandlers.dbConnection)

	if result != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(result); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}

}

func (bHandlers *BookmarkHandlers) ReadPageTitle(w http.ResponseWriter, r *http.Request) {
	bookmarkURL := r.FormValue("bookmark_url")
	bookmarkTitle := ""
	
	resp, err := http.Get(bookmarkURL)
	
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if title, ok := GetHtmlTitle(resp.Body); ok {
		bookmarkTitle = title
	} else {
		println("Fail to get HTML title")
	}

	if bookmarkTitle != "" {
		fmt.Println(bookmarkTitle)
		
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(bookmarkTitle); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

func (bHandlers *BookmarkHandlers) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkURL := r.FormValue("bookmark_url")
	bookmarkGroup := r.FormValue("bookmark_group")
	bookmarkTitle := ""
	result := false
	
	resp, err := http.Get(bookmarkURL)
	
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if title, ok := GetHtmlTitle(resp.Body); ok {
		bookmarkTitle = title
	} else {
		println("Fail to get HTML title")
	}

	if bookmarkTitle != "" {
		fmt.Println(bookmarkTitle)

		bookmark := &Bookmark{
			BookmarkUrl:bookmarkURL, 
			BookmarkTitle:bookmarkTitle, 
			BookmarkGroup:bookmarkGroup}

		result = bookmark.CreateNewBookmark(bHandlers.dbConnection)

		if result {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(result); err != nil {
				panic(err)
			}
			return
		}

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

func (bHandlers *BookmarkHandlers) SaveBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkURL := r.FormValue("bookmark_url")
	bookmarkTitle := r.FormValue("bookmark_title")
	bookmarkGroup := r.FormValue("bookmark_group")
	result := false

	if bookmarkTitle != "" {
		fmt.Println(bookmarkTitle)

		bookmark := &Bookmark{
			BookmarkUrl:bookmarkURL, 
			BookmarkTitle:bookmarkTitle, 
			BookmarkGroup:bookmarkGroup}

		result = bookmark.CreateNewBookmark(bHandlers.dbConnection)

		if result {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(result); err != nil {
				panic(err)
			}
			return
		}

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

func (bHandlers *BookmarkHandlers) ListBookmarks(w http.ResponseWriter, r *http.Request) {
	bookmark := &Bookmark{}

	result := bookmark.ListAllBookmarks(bHandlers.dbConnection)

		if result != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(result); err != nil {
				panic(err)
			}
			return
		}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

func (bHandlers *BookmarkHandlers) ListBookmarksInGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.FormValue("group_id")

	bookmark := &Bookmark{}

	result := bookmark.ListAllBookmarksInGroup(bHandlers.dbConnection, groupID)

		if result != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(result); err != nil {
				panic(err)
			}
			return
		}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

func (bHandlers *BookmarkHandlers) UpdateBookmarks(w http.ResponseWriter, r *http.Request) {
	bookmarkTitle := r.FormValue("bookmark_title")
	bookmarkID := r.FormValue("bookmark_id")

	result := false

	if bookmarkTitle != "" {
		fmt.Println(bookmarkTitle)

		bookmark := &Bookmark{}

		result = bookmark.UpdateExistingBookmark(bHandlers.dbConnection, bookmarkID, bookmarkTitle)

		if result {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(result); err != nil {
				panic(err)
			}
			return
		}

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

func (bHandlers *BookmarkHandlers) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	bookmarkID := r.FormValue("bookmark_id")
	result := false
	
	if bookmarkID != "" {
		bookmark := &Bookmark{}

		result = bookmark.DeleteBookmarkByID(bHandlers.dbConnection, bookmarkID)

		if result {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(result); err != nil {
				panic(err)
			}
			return
		}

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}