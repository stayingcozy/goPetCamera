package main 

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
)

func UpdateDoc(callRef *firestore.DocumentRef, ctx context.Context, answer Answer, ud chan *firestore.WriteResult) {

	passtorch, Uerr := callRef.Update(ctx, []firestore.Update{
		{
			Path:  "answer",
			Value: answer,
		},
	})
	if Uerr != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", Uerr)
	}

	ud <- passtorch

}