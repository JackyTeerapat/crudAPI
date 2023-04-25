// package initializers

// import (
// 	"log"

// 	"github.com/joho/godotenv"
// )

// func LoadEnvVariables() {

// 	err := godotenv.Load(".env")

//		if err != nil {
//			log.Fatal("Error loading .env file")
//		}
//		// Load environment variables
//		// LoadEnvVariables()
//	}
package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
}
