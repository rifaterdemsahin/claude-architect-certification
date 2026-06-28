package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

// ── Supabase config ───────────────────────────────────────────────────────────

const (
	supabaseURL = "https://rmekfsdhglyiralxvkwc.supabase.co"
	pageID      = "00000000-0000-0000-0000-000000000001"
)

func anonKey() string {
	if k := os.Getenv("SUPABASE_ANON_KEY"); k != "" {
		return k
	}
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJtZWtmc2RoZ2x5aXJhbHh2a3djIiwicm9sZSI6ImFub24iLCJpYXQiOjE3ODA3NzYzODgsImV4cCI6MjA5NjM1MjM4OH0.ay6HsmO_CbMGDXL_Dv4DT-paAIAOqlQPfmqqx_mexSU"
}

// ── Data models ───────────────────────────────────────────────────────────────

type ProblemPage struct {
	ID          string `json:"id"`
	StageNumber int    `json:"stage_number"`
	Title       string `json:"title"`
	Headline    string `json:"headline"`
	Subheadline string `json:"subheadline"`
}

type Persona struct {
	ID           string `json:"id"`
	Emoji        string `json:"emoji"`
	RoleTitle    string `json:"role_title"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order"`
}

type Challenge struct {
	ID              string `json:"id"`
	ChallengeNumber int    `json:"challenge_number"`
	Title           string `json:"title"`
	Description     string `json:"description"`
}

type ExamDomain struct {
	ID            string   `json:"id"`
	DomainName    string   `json:"domain_name"`
	TopicsCovered []string `json:"topics_covered"`
	WeightPercent int      `json:"weight_percent"`
	DisplayOrder  int      `json:"display_order"`
}

type Solution struct {
	ID          string `json:"id"`
	StepNumber  int    `json:"step_number"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type PageData struct {
	Page       ProblemPage
	Personas   []Persona
	Challenges []Challenge
	Domains    []ExamDomain
	Solutions  []Solution
	AdminKey   string
}

// ── Supabase fetch ────────────────────────────────────────────────────────────

func sbGet(table, query string, dest any) error {
	url := fmt.Sprintf("%s/rest/v1/%s?%s", supabaseURL, table, query)
	req, _ := http.NewRequest("GET", url, nil)
	key := anonKey()
	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s fetch: %w", table, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s: %d %s", table, resp.StatusCode, body)
	}
	return json.NewDecoder(resp.Body).Decode(dest)
}

func fetchPageData() (PageData, error) {
	var (
		page       []ProblemPage
		personas   []Persona
		challenges []Challenge
		domains    []ExamDomain
		solutions  []Solution
		mu         sync.Mutex
		wg         sync.WaitGroup
		errs       []error
	)

	fetch := func(table, query string, dest any) {
		defer wg.Done()
		if err := sbGet(table, query, dest); err != nil {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		}
	}

	const pid = "page_id=eq." + pageID
	filterAsc := pid + "&order=display_order.asc"

	wg.Add(5)
	go fetch("problem_pages", "id=eq."+pageID+"&select=*", &page)
	go fetch("target_personas", filterAsc, &personas)
	go fetch("core_challenges", pid+"&order=challenge_number.asc", &challenges)
	go fetch("exam_domains", filterAsc, &domains)
	go fetch("course_solutions", pid+"&order=step_number.asc", &solutions)
	wg.Wait()

	if len(errs) > 0 {
		return PageData{}, errs[0]
	}

	var p ProblemPage
	if len(page) > 0 {
		p = page[0]
	}

	return PageData{
		Page:       p,
		Personas:   personas,
		Challenges: challenges,
		Domains:    domains,
		Solutions:  solutions,
		AdminKey:   os.Getenv("ADMIN_KEY"),
	}, nil
}

// ── HTTP handlers ─────────────────────────────────────────────────────────────

var tmpl = template.Must(template.ParseFiles("templates/problem.html"))

func handleProblem(w http.ResponseWriter, r *http.Request) {
	data, err := fetchPageData()
	if err != nil {
		http.Error(w, "Failed to load data: "+err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("template error: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/", handleProblem)

	log.Printf("problem-server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
