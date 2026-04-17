package main

import (
	"Assessment_gift-redemption/internal/handler"
	"Assessment_gift-redemption/internal/repository"
	"Assessment_gift-redemption/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// returns the value of the env var key or fallback val if var not set
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	//configuration
	staffMappingFile := getEnv("STAFF_MAPPING_FILE", "data/staffmapping.csv")
	redempFile := getEnv("REDEMP_FILE", "data/redemptions.csv")
	port := getEnv("PORT", "8080")

	//initialise the data repo
	staffRepo, err := repository.NewStaffRepository(staffMappingFile)
	if err != nil {
		//kills the program if fails to load csv
		log.Fatalf("failed to load staff mappings from %s: %v", staffMappingFile, err)
	}
	log.Printf("staff mappings loaded from %s", staffMappingFile)

	redempRepo, err := repository.NewRedempRepo(redempFile)
	if err != nil {
		log.Fatalf("failed to load redemptions from %s: %v", redempFile, err)
	}
	log.Printf("Redemption data loaded from %s", redempFile)

	//initialise service
	svc := service.NewRedempService(staffRepo, redempRepo)

	//initialise handler
	h := handler.NewHandler(svc)

	//set up router
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	//start web server
	addr := ":" + port
	log.Printf("Gift Redemption API running on http://localhost%s", addr)
	log.Printf("Routes:")
	log.Printf("/health")
	log.Printf("/redeem")

	//cleanup after ending session (for easier testing)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	//start web server in goroutine to prevent ListenAndServe from blocking the rest of the code
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- http.ListenAndServe(addr, mux)
	}()

	//block and wait for server crash or close signal
	select {
	case err := <-serverErrors:
		//if server crashes
		log.Fatalf("server failed: %v", err)
	case sig := <-shutdown:
		//user press ctrl c
		log.Printf("\nreceived signal: %v", sig)
		log.Printf("Deleting test data file: %s", redempFile)
		if err := os.Remove(redempFile); err != nil {
			//if file doesnt exist
			log.Printf("could not delete %s: %v", redempFile, err)
		} else {
			log.Printf("successfully deleted %s", redempFile)
		}
		log.Println("Shutdown complete! byebye :D")
	}

	//ListenAndServe blocks main thread forvever, waiting for incoming traffic
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
