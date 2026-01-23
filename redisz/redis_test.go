package redisz

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func GetPool(ctx context.Context) *Pool {
	return NewPool(ctx, "localhost:6379", "")
}

func TestConn_Set(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	reply, err := conn.Set(ctx, redisKey, 123)
	if err != nil {
		t.Error(err)
		return
	}
	if reply != "OK" {
		t.Log(reply)
	}

	reply, err = conn.Type(ctx, redisKey)
	if err != nil {
		t.Error(err)
		return
	}

	if reply != "string" {
		t.Error(reply)
	}

	ttl, err := conn.TTL(ctx, redisKey)
	if err != nil {
		t.Error(err)
		return
	}
	if ttl > 0 {
		t.Error(ttl)
	}
}

func TestConn_SetNX(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	do, err := conn.SetNxEx(ctx, redisKey, 123, 100)
	if err != nil {
		t.Error(err)
	}
	if !do {
		t.Error(do)
	}
	do, err = conn.SetNxEx(ctx, redisKey, 123, 100)
	if err == nil {
		t.Error("err must not nil")
	}
	if do {
		t.Error(do)
	}
	ttl, err := conn.TTL(ctx, redisKey)
	if err != nil {
		t.Error(err)
	}
	if ttl < 10 {
		t.Error(ttl)
	}

}

func TestConn_GetString(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	_, err := conn.Set(ctx, redisKey, "llal")
	if err != nil {
		t.Error(err)
		return
	}

	reply, err := conn.GetString(ctx, redisKey)
	if err != nil {
		t.Error(err)
		return
	}
	if reply != "llal" {
		t.Error(reply)
	}

}

func TestConn_Get(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	reply, err := conn.Get(ctx, redisKey)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(reply)
}

func TestConn_GetNotSetValue(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	reply2, err := conn.GetInt64(ctx, "not_set_value_key")
	if err != nil {
		if !errors.Is(err, redis.ErrNil) {
			t.Error(err)
			return
		}
	}
	t.Log(reply2)
}

func TestConn_GetKeyValueNotString(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	reply2, err := conn.GetInt64(ctx, redisKey)
	if err != nil {
		if !errors.Is(err, redis.ErrNil) {
			t.Error(err)
			return
		}
	}
	t.Log(reply2)
}

func TestConn_Decr(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	reply, err := conn.Decr(ctx, redisKey)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(reply)

}

func TestConn_DecrBy(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	reply, err := conn.DecrBy(ctx, redisKey, 2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(reply)

}

func TestConn_Incr(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	reply, err := conn.Incr(ctx, redisKey)
	if err != nil {
		t.Error(err)
		return
	}
	if reply != 1 {
		t.Error(reply)
	}

}

func TestConn_IncrBy(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	reply, err := conn.IncrBy(ctx, redisKey, 10)
	if err != nil {
		t.Error(err)
		return
	}
	if reply != 10 {
		t.Error(reply)
	}
}

func TestConn_SAdd(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	reply, err := conn.SAdd(ctx, redisKey, "123")
	if err != nil {
		t.Error(err)
		return
	}
	if reply != 1 {
		t.Error(reply)
	}

	reply, err = conn.SAdd(ctx, redisKey, "123", "223", "334")
	if err != nil {
		t.Error(err)
		return
	}
	if reply != 2 {
		t.Error(reply)
	}

	scard, err := conn.SCard(ctx, redisKey)
	if err != nil {
		t.Error(err)
	}
	if scard != 3 {
		t.Error(scard)
	}

}

func TestConn_SinterStore(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	redisKey1 := "test_redis_key1"
	redisKey2 := "test_redis_key2"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)
	defer conn.redisConn.Do("del", redisKey1)
	defer conn.redisConn.Do("del", redisKey2)

	_, err := conn.SAdd(ctx, redisKey1, 123, 223, 323)
	if err != nil {
		t.Error(err)
	}

	_, err = conn.SAdd(ctx, redisKey2, 323, 223, 523)
	if err != nil {
		t.Error(err)
	}

	store, err := conn.SInterStore(ctx, redisKey, redisKey1, redisKey2)
	if err != nil {
		t.Error(err)
	}
	if store != 2 {
		t.Error(store)
	}

	scard, err := conn.SCard(ctx, redisKey)
	if err != nil {
		t.Error(err)
	}
	if scard != 2 {
		t.Error(scard)
	}
}

func TestConn_SDiffStore(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	redisKey1 := "test_redis_key1"
	redisKey2 := "test_redis_key2"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)
	defer conn.redisConn.Do("del", redisKey1)
	defer conn.redisConn.Do("del", redisKey2)

	_, err := conn.SAdd(ctx, redisKey1, 1, 2, 3, 4, 5)
	if err != nil {
		t.Error(err)
	}
	_, err = conn.SAdd(ctx, redisKey2, 1, 2, 3, 223, 523)
	if err != nil {
		t.Error(err)
	}

	store, err := conn.SDiffStore(ctx, redisKey, redisKey1, redisKey2)
	if err != nil {
		t.Error(err)
	}
	if store != 2 {
		t.Error(store)
	}

	scard, err := conn.SCard(ctx, redisKey)
	if err != nil {
		t.Error(err)
	}
	if scard != 2 {
		t.Error(scard)
	}

	var iterator int64
	reply, err := conn.SScan(ctx, redisKey, iterator, 5)
	if err != nil {
		t.Error(err)
		return
	}
	iterator = reply.Iterator
	strings, err := Strings(reply.Reply, nil)
	t.Log(strings)

}

func TestConn_Scan(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_set_key"
	redisKey1 := "test_redis_set_key1"
	redisKey2 := "test_redis_set_key2"
	redisKey3 := "test_redis_list_key3"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)
	defer conn.redisConn.Do("del", redisKey1)
	defer conn.redisConn.Do("del", redisKey2)
	defer conn.redisConn.Do("del", redisKey3)

	_, err := conn.SAdd(ctx, redisKey, 1, 2, 3, 4, 5)
	if err != nil {
		t.Error(err)
	}
	_, err = conn.SAdd(ctx, redisKey1, 1, 2, 3, 4, 5)
	if err != nil {
		t.Error(err)
	}
	_, err = conn.SAdd(ctx, redisKey2, 1, 2, 3, 223, 523)
	if err != nil {
		t.Error(err)
	}
	_, err = conn.LPush(ctx, redisKey3, 1, 2, 3, 223, 523)
	if err != nil {
		t.Error(err)
	}

	var iterator int64
	var reply ScanReply
	var strings []string
	for {
		reply, err = conn.Scan(ctx, iterator, 50, "test_redis_set*")
		if err != nil {
			t.Error(err)
			return
		}
		iterator = reply.Iterator
		strings, err = Strings(reply.Reply, nil)
		t.Log(iterator)
		t.Log(strings)
		if iterator == 0 {
			break
		}
	}

}
func TestConn_SIsMember(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	member, err := conn.SIsMember(ctx, redisKey, 123)
	if err != nil {
		t.Error(err)
	}
	if member {
		t.Error(member)
	}
	_, err = conn.SAdd(ctx, redisKey, 123, 223, 342)
	if err != nil {
		t.Error(err)
	}
	member, err = conn.SIsMember(ctx, redisKey, 123)
	if err != nil {
		t.Error(err)
	}
	if !member {
		t.Error(member)
	}
}

func TestConn_Scard(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)
	scard, err := conn.SCard(ctx, redisKey)
	if err != nil {
		t.Error(err)
	}
	if scard != 0 {
		t.Error(scard)
	}

}

func TestConn_SRem(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	_, err := conn.SAdd(ctx, redisKey, 123, 223, 323)
	if err != nil {
		t.Error(err)
	}
	reply, err := conn.SRem(ctx, redisKey, 123, 323)
	if err != nil {
		t.Error(err)
	}
	if reply != 2 {
		t.Error(reply)
	}

	reply, err = conn.SRem(ctx, redisKey, 223)
	if err != nil {
		t.Error(err)
	}
	if reply != 1 {
		t.Error(reply)
	}

}

func TestConn_SScan(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)
	for i := 0; i < 100; i++ {
		conn.SAdd(ctx, redisKey, fmt.Sprintf("klalla_%d", i))
	}
	var iterator int64
	var allResult []string
	for {
		reply, err := conn.SScan(ctx, redisKey, iterator, 5)
		if err != nil {
			t.Error(err)
			break
		}
		iterator = reply.Iterator
		strings, err := Strings(reply.Reply, nil)
		allResult = append(allResult, strings...)
		if iterator == 0 {
			break
		}
	}
	if len(allResult) != 100 {
		t.Error(len(allResult))
	}

}

func TestConn_LPush(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	reply, err := conn.LPush(ctx, redisKey, 123, 223, 323)
	if err != nil {
		t.Error(err)
	}
	if reply != 3 {
		t.Error(reply)
	}
	lPop, err := Int64(conn.LPop(ctx, redisKey))

	if err != nil {
		t.Error(err)
	}
	if lPop != 323 {
		t.Error(lPop)
	}
}

func TestConn_RPush(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	reply, err := conn.LPush(ctx, redisKey, 123, 223, 323)
	if err != nil {
		t.Error(err)
	}
	if reply != 3 {
		t.Error(reply)
	}
	lPop, err := Int64(conn.LPop(ctx, redisKey))

	if err != nil {
		t.Error(err)
	}
	if lPop != 323 {
		t.Error(lPop)
	}
}

func TestConn_LPop(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	conn.LPush(ctx, redisKey, "lala")

	pop, err := conn.LPop(ctx, redisKey)
	if err != nil {
		t.Error(err)
	}
	t.Log(String(pop, err))
}

func TestConn_RPop(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	conn.LPush(ctx, redisKey, "lala")

	pop, err := conn.RPop(ctx, redisKey)
	if err != nil {
		t.Error(err)
	}
	t.Log(String(pop, err))
}

func TestConn_BLPopString(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	time.Sleep(time.Second)
	_, err := conn.RPush(ctx, redisKey, "123")
	if err != nil {
		t.Error(err)
	}

	pop, err := conn.BLPop(ctx, redisKey, 10)
	if err != nil {
		t.Error(err)
		return
	}
	s, err := String(pop[0], nil)
	if s != redisKey {
		t.Error(s)
	}
	s2, err := String(pop[1], nil)
	if s2 != "123" {
		t.Error(s2)
	}
}

func TestConn_BRPopString(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	time.Sleep(time.Second)
	_, err := conn.RPush(ctx, redisKey, "123")
	if err != nil {
		t.Error(err)
	}

	pop, err := conn.BRPop(ctx, redisKey, 10)
	if err != nil {
		t.Error(err)
		return
	}
	s, err := String(pop[0], nil)
	if s != redisKey {
		t.Error(s)
	}
	s2, err := String(pop[1], nil)
	if s2 != "123" {
		t.Error(s2)
	}
}

func TestConn_LLen(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	lLen, err := conn.LLen(ctx, redisKey)
	if err != nil {
		t.Error(err)
	}
	if lLen != 0 {
		t.Error(lLen)
	}
	conn.RPush(ctx, redisKey, 1231)
	lLen, err = conn.LLen(ctx, redisKey)
	if err != nil {
		t.Error(err)
	}
	if lLen != 1 {
		t.Error(lLen)
	}
}

func TestConn_LTrim(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	for i := 0; i < 100; i++ {
		conn.LPush(ctx, redisKey, fmt.Sprintf("klalla_%d", i))
	}
	trim, err := conn.LTrim(ctx, redisKey, 0, 19)
	if err != nil {
		t.Error(err)
	}

	if trim != "OK" {
		t.Error(trim)
	}

	lLen, err := conn.LLen(ctx, redisKey)
	if err != nil {
		t.Error(err)
	}
	if lLen != 20 {
		t.Error(lLen)
	}

}

func TestConn_ZAdd(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	//正常添加 返回值 1
	reply, err := Int64(conn.ZAdd(ctx, redisKey, 1, "al"))
	if err != nil {
		t.Error(err)
	}
	if reply != 1 {
		t.Error(reply)
	}
	//返回所有有变更的成员数量
	reply, err = Int64(conn.ZAdd(ctx, redisKey, 2, "al", "CH"))
	if err != nil {
		t.Error(err)
	}
	if reply != 1 {
		t.Error(reply)
	}

	//不更新存在的成员。只添加新成员 返回值 0
	reply, err = Int64(conn.ZAdd(ctx, redisKey, 2, "al", "NX"))
	if err != nil {
		t.Error(err)
	}
	if reply != 0 {
		t.Error(reply)
	}
	//仅更新已经存在的成员 返回1
	reply, err = Int64(conn.ZAdd(ctx, redisKey, 3, "al", "XX", "CH"))
	if err != nil {
		t.Error(err)
	}
	if reply != 1 {
		t.Error(reply)
	}
	//仅更新已经存在的成员 返回0
	reply, err = Int64(conn.ZAdd(ctx, redisKey, 3, "alal", "XX"))
	if err != nil {
		t.Error(err)
	}
	if reply != 0 {
		t.Error(reply)
	}
	//仅更新已经存在的成员 返回0
	rep, err := Float64(conn.ZAdd(ctx, redisKey, 3, "al", "XX", "INCR"))
	if err != nil {
		t.Error(err)
	}
	if rep != 0.6 {
		t.Error(rep)
	}

	do, err := Values(conn.redisConn.Do("ZRANGE", redisKey, 0, -1, "WITHSCORES"))
	if err != nil {
		t.Error(err)
	}

	t.Log(String(do[0], nil))
	t.Log(Float64(do[1], nil))
}

func TestConn_ZCount(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	conn.ZAdd(ctx, redisKey, 1, "lla")
	conn.ZAdd(ctx, redisKey, 2, "lalal2")
	conn.ZAdd(ctx, redisKey, 3, "lalal3")

	count, err := conn.ZCount(ctx, redisKey, 1, 2)
	if err != nil {
		t.Error(err)
	}
	if count != 2 {
		t.Error(count)
	}

}

func TestConn_ZCard(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	conn.ZAdd(ctx, redisKey, 1, "lla")
	conn.ZAdd(ctx, redisKey, 2, "lalal2")
	conn.ZAdd(ctx, redisKey, 3, "lalal3")

	count, err := conn.ZCard(ctx, redisKey)
	if err != nil {
		t.Error(err)
	}
	if count != 3 {
		t.Error(count)
	}

}

func TestConn_ZIncrBy(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	conn.ZAdd(ctx, redisKey, 1, "lala")

	zIncrBy, err := conn.ZIncrBy(ctx, redisKey, 1, "lala")
	if err != nil {
		t.Error(err)
	}
	t.Log(Int64(zIncrBy, nil))

}

func TestConn_ZRange(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)
	for i := 0; i < 100; i++ {
		conn.ZAdd(ctx, redisKey, int64(i), fmt.Sprintf("klalla_%d", i))
	}
	zRange, err := Int64Map(conn.ZRange(ctx, redisKey, 10, 11, "WITHSCORES"))
	if err != nil {
		t.Error(err)
	}
	t.Log(zRange)
	for k, i := range zRange {
		t.Log(k)
		t.Log(i)
	}
	if len(zRange) != 2 {
		t.Error(len(zRange))
	}

	stings, err := Strings(conn.ZRange(ctx, redisKey, 10, 11))
	for _, sting := range stings {
		t.Log(sting)
	}
}

func TestConn_ZRevRange(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)
	for i := 0; i < 100; i++ {
		conn.ZAdd(ctx, redisKey, int64(i), fmt.Sprintf("klalla_%d", i))
	}
	zRange, err := Int64Map(conn.ZRevRange(ctx, redisKey, 10, 11, "WITHSCORES"))
	if err != nil {
		t.Error(err)
	}
	t.Log(zRange)
	for k, i := range zRange {
		t.Log(k)
		t.Log(i)
	}
	if len(zRange) != 2 {
		t.Error(len(zRange))
	}

	stings, err := Strings(conn.ZRevRange(ctx, redisKey, 10, 11))
	for _, sting := range stings {
		t.Log(sting)
	}
}

func TestConn_ZRem(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)
	for i := 0; i < 100; i++ {
		conn.ZAdd(ctx, redisKey, int64(i), fmt.Sprintf("klalla_%d", i))
	}
	reply, err := conn.ZRem(ctx, redisKey, "klalla_1", "klalla_2")
	if err != nil {
		t.Error(err)
	}
	if reply != 2 {
		t.Error(reply)
	}
	reply, err = conn.ZCard(ctx, redisKey)
	if err != nil {
		t.Error(err)
	}
	if reply != 98 {
		t.Error(reply)
	}
}

func TestConn_ZRevRank(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)
	for i := 0; i < 100; i++ {
		conn.ZAdd(ctx, redisKey, int64(i), fmt.Sprintf("klalla_%d", i))
	}

	rank, err := conn.ZRevRank(ctx, redisKey, "klalla_29")
	if err != nil {
		t.Error(err)
	}
	if rank != 70 {
		t.Error(rank)
	}
}

func TestConn_ZRank(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)
	for i := 0; i < 100; i++ {
		conn.ZAdd(ctx, redisKey, int64(i), fmt.Sprintf("klalla_%d", i))
	}

	rank, err := conn.ZRank(ctx, redisKey, "klalla_29")
	if err != nil {
		t.Error(err)
	}
	if rank != 29 {
		t.Error(rank)
	}
}

func TestConn_ZScore(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)
	for i := 0; i < 100; i++ {
		conn.ZAdd(ctx, redisKey, int64(i), fmt.Sprintf("klalla_%d", i))
	}

	rank, err := conn.ZScore(ctx, redisKey, "klalla_29")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(rank)
	}

	rank, err = conn.ZScore(ctx, redisKey, "klalla_101")
	if err != nil {
		if !errors.Is(err, ErrNil) {
			t.Error(err)
		}
		t.Log(err)
	} else {
		t.Error(rank)
	}
	rank, err = conn.ZScore(ctx, "123123", "klalla_10")
	if err != nil {
		if !errors.Is(err, ErrNil) {
			t.Error(err)
		}
		t.Log(err)
	} else {
		t.Error(rank)
	}

}

func TestConn_HSet(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	reply, err := conn.HSet(ctx, redisKey, "123", "123")
	if err != nil {
		t.Error(err)
	}
	if reply != 1 {
		t.Error(reply)
	}
	reply, err = conn.HSet(ctx, redisKey, "123", "223")
	if err != nil {
		t.Error(err)
	}
	if reply != 0 {
		t.Error(reply)
	}

	get, err := String(conn.HGet(ctx, redisKey, "123"))
	if err != nil {
		t.Error(err)
	}
	if get != "223" {
		t.Error(get)
	}

	hGet, err := conn.HGet(ctx, "12314", "223")
	if err != nil {
		t.Error(err)
	}
	t.Log(hGet)
}

func TestConn_HSetNX(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	reply, err := conn.HSetNX(ctx, redisKey, "123", "123")
	if err != nil {
		t.Error(err)
	}
	if reply != 1 {
		t.Error(reply)
	}
	reply, err = conn.HSetNX(ctx, redisKey, "123", "223")
	if err != nil {
		t.Error(err)
	}
	if reply != 0 {
		t.Error(reply)
	}

	get, err := String(conn.HGet(ctx, redisKey, "123"))
	if err != nil {
		t.Error(err)
	}
	if get != "123" {
		t.Error(get)
	}

}

func TestConn_HIncrBy(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)

	_, err := conn.HSet(ctx, redisKey, "123", 100)
	if err != nil {
		t.Error(err)
	}

	reply, err := conn.HIncrBy(ctx, redisKey, "123", 5)

	if err != nil {
		t.Error(err)
	}

	if reply != 105 {
		t.Error(reply)
	}
}

func TestConn_HLen(t *testing.T) {
	ctx := context.Background()
	pool := GetPool(ctx)
	redisKey := "test_redis_key"
	defer pool.Close(ctx)
	conn := pool.GetConn()
	defer conn.Close(ctx)
	defer conn.redisConn.Do("del", redisKey)
	_, err := conn.HSet(ctx, redisKey, "123", 100)
	if err != nil {
		t.Error(err)
	}
	_, err = conn.HSet(ctx, redisKey, "223", 100)
	if err != nil {
		t.Error(err)
	}

	hLen, err := conn.HLen(ctx, redisKey)
	if err != nil {
		t.Error(err)
	}
	if hLen != 2 {
		t.Error(hLen)
	}
}
