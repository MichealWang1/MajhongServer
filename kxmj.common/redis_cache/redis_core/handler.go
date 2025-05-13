package redis_core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"kxmj.common/log"
	"time"
)

// The action name for sync.
const (
	UpdateAction = "update"
	InsertAction = "insert"
	DeleteAction = "delete"
)

type Event struct {
	Action  string `json:"action" redis:"action"`
	Schema  string `json:"schema" redis:"schema"`
	Table   string `json:"table" redis:"table"`
	Data    string `json:"data" redis:"data"`
	Updates string `json:"updates" redis:"updates"`
}

type Handler struct {
	client *redis.Client
	caches map[string][]ICache
}

func NewHandler(client *redis.Client) *Handler {
	return &Handler{
		client: client,
		caches: make(map[string][]ICache, 0),
	}
}

func (h *Handler) Register(schema string, table string, cache ICache) {
	field := fmt.Sprintf("%s:%s", schema, table)
	tableCaches, has := h.caches[field]
	if has == false {
		tableCaches = make([]ICache, 0)
	}
	tableCaches = append(tableCaches, cache)
	h.caches[field] = tableCaches
}

func (h *Handler) Consumer() {
	for field, caches := range h.caches {
		if len(field) <= 0 {
			continue
		}

		if len(caches) <= 0 {
			continue
		}

		go h.pop(field)
	}
}

func (h *Handler) pop(field string) {
	for {
		ctx := context.Background()
		key := h.generateKey(field)

		list, err := h.client.BRPop(ctx, time.Second*60, key).Result()
		if err != nil {
			log.Sugar().Error(fmt.Sprintf("pop field:%s err:%v", field, err))
			continue
		}

		caches, has := h.caches[field]
		if has == false {
			return
		}

		if len(list) <= 1 {
			log.Sugar().Error(fmt.Sprintf("pop field:%s list:%v err:%v", field, list, err))
			continue
		}

		if key != list[0] {
			log.Sugar().Error(fmt.Sprintf("pop field:%s list:%v err:%v", field, list, err))
			continue
		}

		result := list[1]
		event := &Event{}
		err = json.Unmarshal([]byte(result), event)
		if err != nil {
			log.Sugar().Error(fmt.Sprintf("Decode Event field:%s data:%s err:%v", field, result, err))
			continue
		}

		template := caches[0].GetTableTemplate()
		updates := make(map[string]interface{})
		err = json.Unmarshal([]byte(event.Data), template)
		if err != nil {
			log.Sugar().Error(fmt.Sprintf("Decode TableTemplate field:%s data:%s err:%v", field, result, err))
			continue
		}

		if len(event.Updates) > 0 {
			err = json.Unmarshal([]byte(event.Updates), &updates)
			if err != nil {
				log.Sugar().Error(fmt.Sprintf("Decode Updates field:%s data:%s err:%v", field, result, err))
				continue
			}
		}

		for _, cache := range caches {
			cache.EventHandler(ctx, &EventParams{
				Action:  event.Action,
				Data:    template,
				Updates: updates,
			})
		}
	}
}

func (h *Handler) Publish(ctx context.Context, event *Event, field string) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	key := h.generateKey(field)
	return h.client.LPush(ctx, key, data).Err()
}

func (h *Handler) generateKey(field string) string {
	return fmt.Sprintf(EventFormatKey, field)
}
