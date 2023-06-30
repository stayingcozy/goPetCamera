package main  

import (
	"context"
	"cloud.google.com/go/firestore"
	"log"
)

func GetDoc(callRef *firestore.DocumentRef, ctx context.Context, callSnap chan *firestore.DocumentSnapshot) {
	passtorch, err := callRef.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to create document: %v", err)
	}
	callSnap <- passtorch
}