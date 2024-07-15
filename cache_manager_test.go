package cache

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type Article struct {
	Title   string
	Content string
}

var articleList = []Article{
	{
		Title:   "The Evolution of Lightsaber Battles in Star Wars",
		Content: "Lightsaber battles have evolved significantly since the original Star Wars trilogy. From the slow, methodical duels of the early films to the fast-paced, choreographed fights in the prequels and sequels, the depiction of lightsaber combat has changed dramatically. This article examines the stylistic differences and technological advancements that have shaped these iconic scenes over the decades.",
	},
	{
		Title:   "Iconic Star Wars Locations: From Tatooine to Exegol",
		Content: "The Star Wars universe is filled with iconic locations that have become legendary in popular culture. This article takes readers on a journey through these memorable settings, from the deserts of Tatooine to the mysterious Sith planet of Exegol, exploring their significance and role in the saga.",
	},
	{
		Title:   "Exploring the Mysteries of the Force in Star Wars",
		Content: "The Force is a central element of the Star Wars universe, embodying the struggle between the light and dark sides. This article explores the different interpretations and manifestations of the Force, from the teachings of the Jedi to the dark arts of the Sith, and how it has shaped the destiny of key characters.",
	},
}

const ARTIClE_LIST_KEY = "ARTIClE_LIST"

func Test_cacheManager_Get(t *testing.T) {
	type fields struct {
		redis         *redis.Client
		hitCounter    prometheus.Counter
		missedCounter prometheus.Counter
	}
	type args struct {
		ctx context.Context
		key string
	}

	db, mockRedis := redismock.NewClientMock()
	hitCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_hit",
		Help: "Cache Hit",
	})
	missedCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_missed",
		Help: "Cache Missed",
	})
	data := "CACHED_DATA"

	tests := []struct {
		name     string
		fields   fields
		args     args
		mockFunc func()
		want     *string
	}{
		{
			name: "Should success get cache from redis",
			fields: fields{
				redis:         db,
				hitCounter:    hitCounter,
				missedCounter: missedCounter,
			},
			args: args{
				ctx: context.Background(),
				key: ARTIClE_LIST_KEY,
			},
			mockFunc: func() {
				mockRedis.ExpectGet(ARTIClE_LIST_KEY).SetVal(data)
			},
			want: &data,
		},
		{
			name: "Should get no cache from redis",
			fields: fields{
				redis:         db,
				hitCounter:    hitCounter,
				missedCounter: missedCounter,
			},
			args: args{
				ctx: context.Background(),
				key: ARTIClE_LIST_KEY,
			},
			mockFunc: func() {
				mockRedis.ExpectGet(ARTIClE_LIST_KEY).RedisNil()
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		tt.mockFunc()
		t.Run(tt.name, func(t *testing.T) {
			c := cacheManager{
				redis:         tt.fields.redis,
				hitCounter:    tt.fields.hitCounter,
				missedCounter: tt.fields.missedCounter,
			}
			got := c.Get(tt.args.ctx, tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cacheManager_Set(t *testing.T) {
	type fields struct {
		redis         *redis.Client
		hitCounter    prometheus.Counter
		missedCounter prometheus.Counter
	}
	type args struct {
		ctx context.Context
		key string
		val any
	}
	db, mockRedis := redismock.NewClientMock()

	tests := []struct {
		name     string
		fields   fields
		args     args
		mockFunc func()
		wantErr  error
	}{
		{
			name: "Success success add cache",
			fields: fields{
				redis:         db,
				hitCounter:    nil,
				missedCounter: nil,
			},
			args: args{
				ctx: context.Background(),
				key: ARTIClE_LIST_KEY,
				val: articleList,
			},
			mockFunc: func() {
				result, _ := json.Marshal(articleList)
				mockRedis.ExpectSet(ARTIClE_LIST_KEY, string(result), 24*time.Hour).SetVal("")
			},
			wantErr: nil,
		},
		{
			name: "Failed to add cache",
			fields: fields{
				redis:         db,
				hitCounter:    nil,
				missedCounter: nil,
			},
			args: args{
				ctx: context.Background(),
				key: ARTIClE_LIST_KEY,
				val: articleList,
			},
			mockFunc: func() {
				result, _ := json.Marshal(articleList)
				mockRedis.ExpectSet(ARTIClE_LIST_KEY, string(result), 24*time.Hour).SetErr(errors.New("Invalid value"))
			},
			wantErr: errors.New("Invalid value"),
		},
	}
	for _, tt := range tests {
		tt.mockFunc()
		t.Run(tt.name, func(t *testing.T) {
			c := cacheManager{
				redis:         tt.fields.redis,
				hitCounter:    tt.fields.hitCounter,
				missedCounter: tt.fields.missedCounter,
			}

			err := c.Set(tt.args.ctx, tt.args.key, tt.args.val)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_cacheManager_Del(t *testing.T) {
	type fields struct {
		redis         *redis.Client
		hitCounter    prometheus.Counter
		missedCounter prometheus.Counter
	}
	type args struct {
		ctx context.Context
		key string
	}
	db, mockRedis := redismock.NewClientMock()
	tests := []struct {
		name     string
		fields   fields
		args     args
		mockFunc func()
		wantErr  error
	}{
		{
			name: "Success delete cache",
			fields: fields{
				redis:         db,
				hitCounter:    nil,
				missedCounter: nil,
			},
			args: args{
				ctx: context.Background(),
				key: ARTIClE_LIST_KEY,
			},
			mockFunc: func() {
				mockRedis.ExpectDel(ARTIClE_LIST_KEY).SetVal(0)
			},
			wantErr: nil,
		},
		{
			name: "Failed delete cache",
			fields: fields{
				redis:         db,
				hitCounter:    nil,
				missedCounter: nil,
			},
			args: args{
				ctx: context.Background(),
				key: ARTIClE_LIST_KEY,
			},
			mockFunc: func() {
				mockRedis.ExpectDel(ARTIClE_LIST_KEY).SetErr(errors.New("Not found"))
			},
			wantErr: errors.New("Not found"),
		},
	}
	for _, tt := range tests {
		tt.mockFunc()
		t.Run(tt.name, func(t *testing.T) {
			c := cacheManager{
				redis:         tt.fields.redis,
				hitCounter:    tt.fields.hitCounter,
				missedCounter: tt.fields.missedCounter,
			}
			err := c.Del(tt.args.ctx, tt.args.key)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
