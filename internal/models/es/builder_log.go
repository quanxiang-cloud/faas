package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/quanxiang-cloud/faas/internal/models"
	"time"

	"github.com/olivere/elastic/v7"
)

// es query name
const (
	KuberRef      = "kubernetes.labels.faas.module/ref.keyword"
	KuberPipeline = "kubernetes.labels.tekton.dev/pipelineTask"
	KuberTask     = "kubernetes.labels.tekton.dev/task"
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

func (b *builderLogRepo) Search(ctx context.Context, ID string, time time.Time,
	page, size int) ([]*models.BuilderLog, int64, error) {
	query := elastic.NewBoolQuery()

	boolQuery := make([]elastic.Query, 0, 2)
	boolQuery = append(
		boolQuery,
		elastic.NewTermQuery(KuberRef, ID),
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
