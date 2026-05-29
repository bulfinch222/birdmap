package main

import (
	"birdmap/internal"
	"birdmap/internal/lifelist"
	"birdmap/internal/maps"
	"birdmap/internal/taxonomy"
	"context"
	"database/sql"
	"log"
 	_ "modernc.org/sqlite"
	"github.com/wailsapp/wails/v2/pkg/runtime" 
	"os"
	"github.com/joho/godotenv"

)

// App struct
type App struct {
	ctx             context.Context
	db              *sql.DB
	taxonomyService *taxonomy.Service
	userText 		string
	ready 			chan struct{}
	key 			string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.ready = make(chan struct{})
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл")
	}
	apiKey := os.Getenv("EBIRD_API_KEY")
	if apiKey == "" {
    	log.Fatal("EBIRD_API_KEY not set")
	}
	log.Println(apiKey)
	a.key = apiKey
	db, err := sql.Open("sqlite", "./bird.db")
	if err != nil {
		log.Fatal(err)
	}
	a.db = db
	a.taxonomyService = taxonomy.NewService(a.db)
	err = a.taxonomyService.SyncTaxonomy(a.key, a.ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = lifelist.CreateLifelist(a.taxonomyService)
	if err != nil {
		panic(err)	
	}
	close(a.ready)
}

func (a *App) shutdown(ctx context.Context) {
	a.db.Close()
}

func (a *App) getCurrentLocation() {
	<- a.ready
	maps.GetCurrentLocation()
}

func (a *App) GetLocations() []maps.Location {
	<- a.ready
	locations, err := maps.GetLocations(a.key, a.taxonomyService)
	if err != nil {
		panic(err)
	}
	return locations
}

func (a *App) GetUserText(input string) {

	a.userText = input
}

func (a *App) SearchBirds(input string) ([]internal.Bird) {
	<- a.ready
	birds, err := a.taxonomyService.SearchBirds(input)	
	if err != nil {
		panic(err)
	}
	return birds
}

func (a *App) GetLifelist() ([]internal.Bird){
	<- a.ready
	birds, err := lifelist.GetLifelist(a.taxonomyService)
	if err != nil {
		panic(err)
	}
	return birds
}

func (a *App) AddBird(bird internal.Bird) {
	<- a.ready
	err := lifelist.AddBird(bird, a.taxonomyService)
	if err != nil {
		panic(err)
	}
}

func (a *App) IsReady() bool {
	select {
		case <- a.ready:
			return true
		default:
			return false
	}
}

func (a *App) ImportLifelist(){
	<- a.ready
	filepath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Выберите импортируемый файл",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "CSV Files (*.csv)",
				Pattern:     "*.csv",
			},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
	err = lifelist.ImportLifelist(a.taxonomyService, filepath)
	if err != nil {
		log.Fatal(err)
	}
}