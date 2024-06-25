package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type DeveloperLevel int

const (
	Junior DeveloperLevel = iota
	Middle
	Senior
)

type Developer struct {
	ID                    int            `json:"id"`
	Name                  string         `json:"name"`
	DeveloperLevel        DeveloperLevel `json:"developerLevel"`
	EfficiencyCoefficient float64        `json:"efficiencyCoefficient"`
	TaskList              []Task         `json:"taskList"`
}

type Sprint struct {
	ID        int       `json:"id"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	TaskList  []Task    `json:"taskList"`
}

type Team struct {
	ID         int           `json:"id"`
	Developers []*Developer  `json:"developers"`
	Sprints    []*Sprint     `json:"sprints"`
}

type Task struct {
	ID          int           `json:"id"`
	AverageTime time.Duration `json:"averageTime"`
	StartDate   *time.Time    `json:"startDate"`
	EndDate     *time.Time    `json:"endDate"`
	SprintID    int           `json:"sprintId"`
}

var teams = map[int]*Team{}

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/teams", listTeams)
	r.Post("/teams", createTeam)
	r.Put("/teams/{id}", editTeam)
	r.Get("/teams/{id}", getTeam)
	r.Post("/teams/{id}/sprints", createSprint)
	r.Post("/teams/{id}/developers", createDeveloper)

	// Init
	initDemoData()

	// Start the server
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}

func initDemoData() {
	// Create team 0
	team0 := &Team{
		ID:         0,
		Developers: []*Developer{
			{
				ID:                    0,
				Name:                  "John Doe",
				DeveloperLevel:        Middle,
				EfficiencyCoefficient: 0.8,
			},
			{
				ID:                    1,
				Name:                  "Jane Smith",
				DeveloperLevel:        Senior,
				EfficiencyCoefficient: 0.9,
			},
		},
		Sprints: []*Sprint{
			{
				ID:        0,
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(0, 0, 14),
				TaskList: []Task{
					{
						ID:          0,
						AverageTime: 20 * time.Hour,
						StartDate:   nil,
						EndDate:     nil,
						SprintID:    0,
					},
					{
						ID:          1,
						AverageTime: 30 * time.Hour,
						StartDate:   nil,
						EndDate:     nil,
						SprintID:    0,
					},
				},
			},
		},
	}

	// Create team 1
	team1 := &Team{
		ID:         1,
		Developers: []*Developer{
			{
				ID:                    2,
				Name:                  "Bob Johnson",
				DeveloperLevel:        Junior,
				EfficiencyCoefficient: 0.7,
			},
			{
				ID:                    3,
				Name:                  "Alice Williams",
				DeveloperLevel:        Middle,
				EfficiencyCoefficient: 0.8,
			},
		},
		Sprints: []*Sprint{
			{
				ID:        0,
				StartDate: time.Now().AddDate(0, 0, 15),
				EndDate:   time.Now().AddDate(0, 0, 29),
				TaskList: []Task{
					{
						ID:          2,
						AverageTime: 25 * time.Hour,
						StartDate:   nil,
						EndDate:     nil,
						SprintID:    0,
					},
					{
						ID:          3,
						AverageTime: 35 * time.Hour,
						StartDate:   nil,
						EndDate:     nil,
						SprintID:    0,
					},
				},
			},
		},
	}

	teams[team0.ID] = team0
	teams[team1.ID] = team1
}

func listTeams(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

func createTeam(w http.ResponseWriter, r *http.Request) {
	var newTeam Team
	err := json.NewDecoder(r.Body).Decode(&newTeam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTeam.ID = len(teams)
	teams[newTeam.ID] = &newTeam

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTeam)
}

func editTeam(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if _, ok := teams[id]; !ok {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	}

	var updatedTeam Team
	err = json.NewDecoder(r.Body).Decode(&updatedTeam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTeam.ID = id
	teams[id] = &updatedTeam

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeam)
}

func getTeam(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if team, ok := teams[id]; ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(team)
	} else {
		http.Error(w, "Team not found", http.StatusNotFound)
	}
}

func createSprint(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if team, ok := teams[id]; ok {
		var newSprint Sprint
		err := json.NewDecoder(r.Body).Decode(&newSprint)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		team.Sprints = append(team.Sprints, &newSprint)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newSprint)
	} else {
		http.Error(w, "Team not found", http.StatusNotFound)
	}
}

func createDeveloper(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if team, ok := teams[id]; ok {
		var newDeveloper Developer
		err := json.NewDecoder(r.Body).Decode(&newDeveloper)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		team.Developers = append(team.Developers, &newDeveloper)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newDeveloper)
	} else {
		http.Error(w, "Team not found", http.StatusNotFound)
	}
}
