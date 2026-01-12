package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	linkvite "github.com/linkvite/go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	baseURL := os.Getenv("LINKVITE_BASE_URL")

	apiKey := os.Getenv("LINKVITE_API_KEY")
	if apiKey == "" {
		log.Fatal("LINKVITE_API_KEY not set in environment")
	}

	client, err := linkvite.NewClient(apiKey, linkvite.WithBaseURL(baseURL))
	if err != nil {
		log.Fatal(err)
	}

	// client, err := linkvite.NewClientWithTokens(
	// 	os.Getenv("LINKVITE_ACCESS_TOKEN"),
	// 	os.Getenv("LINKVITE_REFRESH_TOKEN"),
	// 	linkvite.WithBaseURL(baseURL),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	ctx := context.Background()

	user, err := client.User.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Hello, %s!\n", user.Name)

	collections, err := client.Collections.List(ctx, &linkvite.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, collection := range collections {
		fmt.Printf("Collection: %s : %s\n", collection.ID, collection.Name)
	}

	members, err := client.Collections.GetMembers(ctx, "coll_7VC2WwmGAhDj")
	if err != nil {
		log.Fatal(err)
	}
	for _, member := range members.Members {
		fmt.Printf("Member: %s : %s\n", member.ID, member.Name)
	}

	storage, err := client.User.GetDetailedStorageUsage(ctx, &linkvite.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bookmarks: %v\n", storage.Bookmarks)

	err = client.RefreshAccessToken(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Access token refreshed successfully")
}
