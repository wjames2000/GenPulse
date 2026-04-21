package history

import (
	"context"
	"fmt"
	"time"

	"GenPulse/internal/monitoring/storage"
)

type TraceStoreAdapter struct {
	traceStore *storage.TraceStore
}

func NewTraceStoreAdapter(traceStore *storage.TraceStore) *TraceStoreAdapter {
	return &TraceStoreAdapter{
		traceStore: traceStore,
	}
}

func (tsa *TraceStoreAdapter) GetTraceTree(ctx context.Context, traceID string) ([]interface{}, error) {
	spans, err := tsa.traceStore.GetTraceTree(ctx, traceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trace tree: %w", err)
	}

	// 转换为通用接口类型
	result := make([]interface{}, len(spans))
	for i, span := range spans {
		result[i] = span
	}

	return result, nil
}

func (tsa *TraceStoreAdapter) QuerySpans(ctx context.Context, query interface{}) ([]interface{}, error) {
	// 将通用查询转换为TraceQuery
	traceQuery, ok := query.(storage.TraceQuery)
	if !ok {
		// 尝试从map转换
		if queryMap, ok := query.(map[string]interface{}); ok {
			traceQuery = tsa.convertMapToTraceQuery(queryMap)
		} else {
			return nil, fmt.Errorf("invalid query type: %T", query)
		}
	}

	spans, err := tsa.traceStore.QuerySpans(ctx, traceQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query spans: %w", err)
	}

	// 转换为通用接口类型
	result := make([]interface{}, len(spans))
	for i, span := range spans {
		result[i] = span
	}

	return result, nil
}

func (tsa *TraceStoreAdapter) convertMapToTraceQuery(queryMap map[string]interface{}) storage.TraceQuery {
	traceQuery := storage.TraceQuery{}

	if traceID, ok := queryMap["trace_id"].(string); ok {
		traceQuery.TraceID = traceID
	}
	if spanName, ok := queryMap["span_name"].(string); ok {
		traceQuery.SpanName = spanName
	}
	if agentName, ok := queryMap["agent_name"].(string); ok {
		traceQuery.AgentName = agentName
	}
	if toolName, ok := queryMap["tool_name"].(string); ok {
		traceQuery.ToolName = toolName
	}
	if pipelineID, ok := queryMap["pipeline_id"].(string); ok {
		traceQuery.PipelineID = pipelineID
	}
	if status, ok := queryMap["status"].(string); ok {
		traceQuery.Status = status
	}
	if startTime, ok := queryMap["start_time"].(time.Time); ok {
		traceQuery.StartTime = startTime
	}
	if endTime, ok := queryMap["end_time"].(time.Time); ok {
		traceQuery.EndTime = endTime
	}
	if limit, ok := queryMap["limit"].(int); ok {
		traceQuery.Limit = limit
	}
	if offset, ok := queryMap["offset"].(int); ok {
		traceQuery.Offset = offset
	}

	return traceQuery
}

func (tsa *TraceStoreAdapter) GetSpan(ctx context.Context, id string) (interface{}, error) {
	span, err := tsa.traceStore.GetSpan(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get span: %w", err)
	}
	return span, nil
}

func (tsa *TraceStoreAdapter) StoreSpan(ctx context.Context, span interface{}) error {
	traceSpan, ok := span.(*storage.TraceSpan)
	if !ok {
		return fmt.Errorf("invalid span type: %T", span)
	}

	return tsa.traceStore.StoreSpan(ctx, traceSpan)
}

func (tsa *TraceStoreAdapter) GetAgentSpans(ctx context.Context, agentName string, limit int) ([]interface{}, error) {
	spans, err := tsa.traceStore.GetAgentSpans(ctx, agentName, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent spans: %w", err)
	}

	result := make([]interface{}, len(spans))
	for i, span := range spans {
		result[i] = span
	}

	return result, nil
}

func (tsa *TraceStoreAdapter) GetPipelineSpans(ctx context.Context, pipelineID string) ([]interface{}, error) {
	spans, err := tsa.traceStore.GetPipelineSpans(ctx, pipelineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline spans: %w", err)
	}

	result := make([]interface{}, len(spans))
	for i, span := range spans {
		result[i] = span
	}

	return result, nil
}

func (tsa *TraceStoreAdapter) GetRecentSpans(ctx context.Context, limit int) ([]interface{}, error) {
	spans, err := tsa.traceStore.GetRecentSpans(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent spans: %w", err)
	}

	result := make([]interface{}, len(spans))
	for i, span := range spans {
		result[i] = span
	}

	return result, nil
}

func (tsa *TraceStoreAdapter) GetStatistics(ctx context.Context, startTime, endTime time.Time) (map[string]interface{}, error) {
	return tsa.traceStore.GetStatistics(ctx, startTime, endTime)
}

func (tsa *TraceStoreAdapter) CleanupOldData(ctx context.Context, retentionDays int) error {
	return tsa.traceStore.CleanupOldData(ctx, retentionDays)
}
