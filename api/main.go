package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Course model
type course struct {
	CourseID string  `json:"course_id"`
	Title    string  `json:"title"`
	Author   *author `json:"author"`
	Price    float64 `json:"price"`
}

// Author model
type author struct {
	AuthorID int    `json:"author_id"`
	Name     string `json:"name"`
}

// Fake DB
var courses []course

// Helper function to check if course is empty
func isEmpty(c *course) bool {
	return c.CourseID == "" || c.Title == ""
}

// Main function
func main() {
	// Seed data
	courses = append(courses, course{
		CourseID: "1",
		Title:    "First Course",
		Author:   &author{AuthorID: 1, Name: "John Doe"},
		Price:    10.00,
	})
	courses = append(courses, course{
		CourseID: "2",
		Title:    "Second Course",
		Author:   &author{AuthorID: 2, Name: "Jane Smith"},
		Price:    20.00,
	})

	log.Println("🚀 Server is running on http://localhost:3000")

	// Router setup
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/create", createCourse).Methods("POST")
	r.HandleFunc("/getall", getAllCourses).Methods("GET")
	r.HandleFunc("/getone/{id}", getOneCourse).Methods("GET")
	r.HandleFunc("/update/{id}", updateCourse).Methods("PUT")
	r.HandleFunc("/delete/{id}", deleteCourse).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", r))
}

// Home route
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("📚 Welcome to the Course API"))
}

// GET all courses
func getAllCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// GET one course by ID
func getOneCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	for _, item := range courses {
		if item.CourseID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Course not found"})
}

// POST create a new course
func createCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Request body is empty"})
		return
	}

	var c course
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	if isEmpty(&c) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing course ID or title"})
		return
	}

	courses = append(courses, c)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

// PUT update a course by ID
func updateCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	for i, item := range courses {
		if item.CourseID == id {
			var updated course
			if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
				return
			}

			if isEmpty(&updated) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Missing fields in course"})
				return
			}

			courses[i] = updated
			json.NewEncoder(w).Encode(updated)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Course not found"})
}

// DELETE a course by ID
func deleteCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	for i, item := range courses {
		if item.CourseID == id {
			courses = append(courses[:i], courses[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string{"message": "Course deleted successfully"})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Course not found"})
}
