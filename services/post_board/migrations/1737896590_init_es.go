package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"strings"

	post_board_config "github.com/a179346/robert-go-monorepo/services/post_board/config"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up1737896590, Down1737896590)
}

const policy = "post-board-api-30-days"
const indexTemplate = "post-board-api"
const indexPattern = "post-board-api-*"

func Up1737896590(ctx context.Context, tx *sql.Tx) error {
	cfg := elasticsearch.Config{
		Addresses: []string{post_board_config.GetLoggingConfig().ElasticSearchAddress},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("elasticsearch.NewClient error: %w", err)
	}

	body := `{"policy":{"phases":{"hot":{"min_age":"0ms","actions":{}},"warm":{"min_age":"2d","actions":{}},"delete":{"min_age":"30d","actions":{"delete":{"delete_searchable_snapshot":true}}}},"deprecated":false}}`
	upsertILMLifecycleReq := esapi.ILMPutLifecycleRequest{
		Policy: policy,
		Body:   strings.NewReader(body),
	}
	res, err := upsertILMLifecycleReq.Do(ctx, es)
	if err != nil || res.IsError() {
		return fmt.Errorf("upsertILMLifecycleReq.Do error: %w", err)
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(res.Body)

	body = fmt.Sprintf(
		`{"template":{"settings":{"index":{"lifecycle":{"name":"%s"}}}},"index_patterns":["%s"],"data_stream":{}}`,
		policy,
		indexPattern,
	)
	upsertIndexTemplateReq := esapi.IndicesPutIndexTemplateRequest{
		Name: indexTemplate,
		Body: strings.NewReader(body),
	}
	res, err = upsertIndexTemplateReq.Do(ctx, es)
	if err != nil || res.IsError() {
		return fmt.Errorf("upsertIndexTemplateReq.Do error: %w", err)
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(res.Body)

	return nil
}

func Down1737896590(ctx context.Context, tx *sql.Tx) error {
	cfg := elasticsearch.Config{
		Addresses: []string{post_board_config.GetLoggingConfig().ElasticSearchAddress},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("elasticsearch.NewClient error: %w", err)
	}

	deleteILMLifecycleReq := esapi.ILMDeleteLifecycleRequest{
		Policy: policy,
	}
	res, err := deleteILMLifecycleReq.Do(ctx, es)
	if err != nil || res.IsError() {
		return fmt.Errorf("deleteILMLifecycleReq.Do error: %w", err)
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(res.Body)

	deleteIndexTemplateReq := esapi.IndicesDeleteIndexTemplateRequest{
		Name: indexTemplate,
	}
	res, err = deleteIndexTemplateReq.Do(ctx, es)
	if err != nil || res.IsError() {
		return fmt.Errorf("deleteIndexTemplateReq.Do error: %w", err)
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(res.Body)

	return nil
}
