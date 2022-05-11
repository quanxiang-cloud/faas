package es

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/quanxiang-cloud/faas/internal/models"
	"github.com/quanxiang-cloud/faas/pkg/basic/k8s"

	"github.com/olivere/elastic/v7"
)

type builderLogRepo struct {
	db *elastic.Client
}

func NewBuildLogRepo(db *elastic.Client) models.BuilderLogRepo {
	return &builderLogRepo{
		db: db,
	}
}

func (b *builderLogRepo) index() string {
	return "builder*"
}

func (b *builderLogRepo) Search(ctx context.Context, id, step string, time time.Time,
	page, size int) ([]*models.BuilderLog, int64, error) {
	query := elastic.NewBoolQuery()

	boolQuery := make([]elastic.Query, 0, 2)
	if step != "" {
		boolQuery = append(boolQuery, elastic.NewTermQuery(k8s.Step+".keyword", step))
	}
	boolQuery = append(
		boolQuery,
		elastic.NewMatchPhrasePrefixQuery(k8s.ResourceRef, id),
		elastic.NewRangeQuery("time").Gte(time.UTC()),
	)

	query = query.Must(boolQuery...)

	searchResult, err := b.db.Search().
		Index(b.index()).
		Query(query).
		Sort("time", true).
		From(page).Size(size).
		Do(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}

	sh := searchResult.Hits.Hits
	bls := make([]*models.BuilderLog, 0, size)
	for _, v := range sh {
		bl := new(models.BuilderLog)

		err := json.Unmarshal(v.Source, bl)
		if err != nil {
			return nil, 0, err
		}
		bls = append(bls, bl)
	}
	return bls, int64(len(bls)), nil
}
