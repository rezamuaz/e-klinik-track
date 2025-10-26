package types

import (
	"bytes"
	"context"
	"e-klinik/config"
	"e-klinik/internal/domain/dto"
	"e-klinik/internal/domain/request"
	"e-klinik/pkg"
	"e-klinik/utils"

	"strings"

	"github.com/typesense/typesense-go/v3/typesense/api"
)

type IndexRepository interface {
	CreateIndex(c context.Context, col string, document any) error
	DeleteIndex(c context.Context, col string, id string) error
	SearchIndex(c context.Context, args request.SearchParams) (any, any, error)
	Upsert(c context.Context, col string, document any) error
}

type IndexRepositoryImpl struct {
	ts    *pkg.TypeSense
	index string
	cfg   *config.Config
}

// NewIndexRepository instantiates the index repository
func NewIndexRepository(ts *pkg.TypeSense) *IndexRepositoryImpl {
	return &IndexRepositoryImpl{
		ts:    ts,
		index: "pvsave",
	}
}

// CreateIndex inserts a new document into a Typesense collection
func (t *IndexRepositoryImpl) CreateIndex(ctx context.Context, col string, document any) error {
	_, err := t.ts.Client.Collection(col).Documents().Create(ctx, document, &api.DocumentIndexParameters{})
	if err != nil {
		return pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to create document index")
	}
	return nil
}

// Upsert inserts or updates a document in a Typesense collection
func (t *IndexRepositoryImpl) Upsert(ctx context.Context, col string, document any) error {
	_, err := t.ts.Client.Collection(col).Documents().Upsert(ctx, document, &api.DocumentIndexParameters{})
	if err != nil {
		return pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to upsert document index")
	}
	return nil
}

// DeleteIndex removes a document from the Typesense collection
func (t *IndexRepositoryImpl) DeleteIndex(ctx context.Context, col string, id string) error {
	_, err := t.ts.Client.Collection(col).Document(id).Delete(ctx)
	if err != nil {
		return pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to delete document index")
	}
	return nil
}

// SearchIndex performs a search query on the Typesense collection
func (t *IndexRepositoryImpl) SearchIndex(ctx context.Context, args request.SearchParams) (any, any, error) {
	var should bytes.Buffer
	var filters []string

	// Build filter conditions dynamically
	if len(args.Genre) > 0 {
		filters = append(filters, "genres.name:=["+strings.Join(args.Genre, ",")+"]")
	}
	if len(args.Status) > 0 {
		filters = append(filters, "status:=["+strings.Join(args.Status, ",")+"]")
	}
	if len(args.Type) > 0 {
		filters = append(filters, "type:=["+strings.Join(args.Type, ",")+"]")
	}

	if len(filters) > 0 {
		should.WriteString(strings.Join(filters, " && "))
	}

	// Search parameters
	arg := &api.SearchCollectionParams{
		Q:       &args.Query,
		QueryBy: utils.StringPtr("title"),
		FilterBy: func() *string {
			if should.Len() > 0 {
				s := should.String()
				return &s
			}
			return nil
		}(),
	}

	resp, err := t.ts.Client.Collection(t.index).Documents().Search(ctx, arg)
	if err != nil {
		return nil, nil, pkg.WrapError(err, pkg.ErrorCodeInternal, "failed to execute search query")
	}

	pagination := dto.Pagination{
		CurrentPage: int64(*resp.Page),
		TotalPage:   int64(resp.RequestParams.PerPage),
	}

	return resp, pagination, nil
}
