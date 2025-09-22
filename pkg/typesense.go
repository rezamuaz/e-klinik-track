package pkg

import (
	"context"
	"e-klinik/config"
	"e-klinik/utils"
	"fmt"
	"log"

	"github.com/typesense/typesense-go/v3/typesense"
	"github.com/typesense/typesense-go/v3/typesense/api"
)

type TypeSense struct {
	Client *typesense.Client
}

func NewTypeSense(cfg *config.Config) *TypeSense {
	connectionString := fmt.Sprintf("%s:%s", cfg.TypeSense.Host, cfg.TypeSense.Port)
	client := typesense.NewClient(
		typesense.WithServer(connectionString),
		typesense.WithAPIKey(cfg.TypeSense.ApiKey),
	)

	return &TypeSense{Client: client}
}

// typesense.WithConnectionTimeout(5*time.Second),
// 		typesense.WithCircuitBreakerMaxRequests(50),
// 		typesense.WithCircuitBreakerInterval(2*time.Minute),
// 		typesense.WithCircuitBreakerTimeout(1*time.Minute),

// Method to ensure a collection exists or create it
func (ts *TypeSense) EnsureCollectionExists(c context.Context, schema *api.CollectionSchema) error {
	result, err := ts.Client.Collections().Retrieve(c)
	if err != nil {
		// log.Printf("Collection already exists: %s", schema.Name)
		log.Printf("Error indexing document: %v", err)
		return nil
	}
	// Check if the collection already exists
	for _, c := range result {
		if c.Name == schema.Name {
			log.Printf("Collection already exists: %s", schema.Name)
			return nil
		}
	}
	// If not found, create the collection
	newschame, err := ts.Client.Collections().Create(c, schema)
	if err != nil {
		return fmt.Errorf("failed to create collection %s: %w", schema.Name, err)
	}
	utils.DumpTest(newschame)
	log.Printf("Collection created: %s", schema.Name)
	return nil
}
