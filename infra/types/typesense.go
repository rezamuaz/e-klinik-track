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

// NewTask instantiates the Task repository.
func NewIndexRepository(ts *pkg.TypeSense) *IndexRepositoryImpl {
	return &IndexRepositoryImpl{
		ts:    ts,
		index: "pvsave",
	}
}

// Index creates or updates a task in an index.
func (t *IndexRepositoryImpl) CreateIndex(c context.Context, col string, document any) error {
	_, err := t.ts.Client.Collection(col).Documents().Create(c, document, &api.DocumentIndexParameters{})
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "IndexRequest.Do")
	}

	return nil
}

// Index creates or updates a task in an index.
func (t *IndexRepositoryImpl) Upsert(c context.Context, col string, document any) error {

	_, err := t.ts.Client.Collection(col).Documents().Upsert(c, document, &api.DocumentIndexParameters{})
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "UpsertDocument.Do")
	}

	return nil
}

// Delete removes a task from the index.
func (t *IndexRepositoryImpl) DeleteIndex(c context.Context, col string, id string) error {
	_, err := t.ts.Client.Collection(col).Document(id).Delete(c)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "DeleteRequest.Do")
	}

	return nil
}

// Search returns tasks matching a query.
func (t *IndexRepositoryImpl) SearchIndex(c context.Context, args request.SearchParams) (any, any, error) {
	// ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Search")
	// defer span.End()

	// if args.IsZero() {
	// 	return internal.SearchResults{}, nil
	// }
	var should bytes.Buffer
	var stringsList []string

	if len(args.Genre) != 0 {
		stringsList = append(stringsList, "genres.name IN"+"["+strings.Join(args.Genre, ",")+"]")
	}

	if len(args.Status) != 0 {
		stringsList = append(stringsList, "status IN"+"["+strings.Join(args.Status, ",")+"]")
	}

	if len(args.Type) != 0 {
		stringsList = append(stringsList, "type IN"+"["+strings.Join(args.Status, ",")+"]")
	}

	should.WriteString(strings.Join(stringsList, " AND "))

	arg := &api.SearchCollectionParams{
		Q:       &args.Query,
		QueryBy: utils.StringPtr("title"),
	}
	resp, err := t.ts.Client.Collection("pvsave").Documents().Search(c, arg)
	if err != nil {
		return nil, nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "SearchRequest.Do")
	}

	return resp, dto.Pagination{CurrentPage: int64(*resp.Page), TotalPage: int64(resp.RequestParams.PerPage)}, nil
}
