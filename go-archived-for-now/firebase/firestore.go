package firebase

import (
	"context"
	"errors"
	"reflect"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FirestoreService provides Firestore database functionality
type FirestoreService struct {
	client *firestore.Client
}

// NewFirestoreService creates a new FirestoreService
func NewFirestoreService(client *firestore.Client) *FirestoreService {
	return &FirestoreService{
		client: client,
	}
}

// Collection returns a reference to a collection
func (s *FirestoreService) Collection(path string) *firestore.CollectionRef {
	if s.client == nil {
		return nil
	}
	return s.client.Collection(path)
}

// Document returns a reference to a document
func (s *FirestoreService) Document(collectionPath, documentID string) *firestore.DocumentRef {
	if s.client == nil {
		return nil
	}
	return s.client.Collection(collectionPath).Doc(documentID)
}

// Create adds a new document to the specified collection
func (s *FirestoreService) Create(ctx context.Context, collectionPath string, data interface{}) (string, error) {
	if s.client == nil {
		return "", errors.New("firestore client not initialized")
	}

	ref, _, err := s.client.Collection(collectionPath).Add(ctx, data)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}

// Set creates or overwrites a document
func (s *FirestoreService) Set(ctx context.Context, collectionPath, documentID string, data interface{}) error {
	if s.client == nil {
		return errors.New("firestore client not initialized")
	}

	_, err := s.client.Collection(collectionPath).Doc(documentID).Set(ctx, data)
	return err
}

// Update updates specific fields of a document
func (s *FirestoreService) Update(ctx context.Context, collectionPath, documentID string, updates map[string]interface{}) error {
	if s.client == nil {
		return errors.New("firestore client not initialized")
	}

	// Convert map to slice of firestore.Update
	var updateFields []firestore.Update
	for key, value := range updates {
		updateFields = append(updateFields, firestore.Update{
			Path:  key,
			Value: value,
		})
	}

	_, err := s.client.Collection(collectionPath).Doc(documentID).Update(ctx, updateFields)
	return err
}

// Get retrieves a document
func (s *FirestoreService) Get(ctx context.Context, collectionPath, documentID string, dest interface{}) error {
	if s.client == nil {
		return errors.New("firestore client not initialized")
	}

	docRef := s.client.Collection(collectionPath).Doc(documentID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errors.New("document not found")
		}
		return err
	}

	return docSnap.DataTo(dest)
}

// Delete removes a document
func (s *FirestoreService) Delete(ctx context.Context, collectionPath, documentID string) error {
	if s.client == nil {
		return errors.New("firestore client not initialized")
	}

	_, err := s.client.Collection(collectionPath).Doc(documentID).Delete(ctx)
	return err
}

// List retrieves all documents in a collection
func (s *FirestoreService) List(ctx context.Context, collectionPath string) ([]map[string]interface{}, error) {
	if s.client == nil {
		return nil, errors.New("firestore client not initialized")
	}

	var results []map[string]interface{}
	iter := s.client.Collection(collectionPath).Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		data := doc.Data()
		data["id"] = doc.Ref.ID
		results = append(results, data)
	}

	return results, nil
}

// Query executes a query on a collection
func (s *FirestoreService) Query(ctx context.Context, collectionPath string, dest interface{}, queries ...firestore.Query) error {
	if s.client == nil {
		return errors.New("firestore client not initialized")
	}

	// Verify dest is a pointer to a slice
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr || destVal.Elem().Kind() != reflect.Slice {
		return errors.New("destination must be a pointer to a slice")
	}

	// Create a slice to store the results
	sliceType := destVal.Elem().Type()
	elemType := sliceType.Elem()
	results := reflect.MakeSlice(sliceType, 0, 0)

	// Start with the collection reference
	query := s.client.Collection(collectionPath).Query

	// Apply all query filters
	for _, q := range queries {
		query = q
	}

	// Execute the query
	iter := query.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		// Create a new element of the slice's element type
		elem := reflect.New(elemType).Elem()
		elemAddr := elem.Addr().Interface()

		// If the element is a struct, use DataTo
		if elemType.Kind() == reflect.Struct {
			if err := doc.DataTo(elemAddr); err != nil {
				return err
			}
		} else if elemType.Kind() == reflect.Map {
			// If the element is a map, use Data
			data := doc.Data()
			data["id"] = doc.Ref.ID
			elem.Set(reflect.ValueOf(data))
		} else {
			return errors.New("unsupported element type in slice")
		}

		// Append the element to the results
		results = reflect.Append(results, elem)
	}

	// Set the results back to the destination
	destVal.Elem().Set(results)
	return nil
}
