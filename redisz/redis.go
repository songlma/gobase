package redisz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

type Pool struct {
	redisPool   *redis.Pool
	addr, auth  string
	opentracing bool
}

func NewPool(ctx context.Context, addr, auth string) *Pool {
	redisPool := &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 10 * time.Minute,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", addr)
			if err != nil {
				errorLog(ctx, "redisNewPoolDial", err, "Dial tcp")
				return nil, err
			}
			if auth != "" {
				if _, err = conn.Do("AUTH", auth); err != nil {
					errorLog(ctx, "redisNewPoolDoAuTH", err)
					return nil, err
				}
			}
			return conn, nil
		},
	}
	return &Pool{
		redisPool: redisPool,
		addr:      addr,
		auth:      auth,
	}
}

func (p *Pool) Opentracing() {
	p.opentracing = true
}
func (p *Pool) GetConn() *Conn {
	return &Conn{redisConn: p.redisPool.Get(),
		opentracing: p.opentracing,
	}
}

func (p *Pool) Close(ctx context.Context) error {
	return p.redisPool.Close()
}

type Conn struct {
	redisConn   redis.Conn
	opentracing bool
}

func (conn *Conn) Close(ctx context.Context) error {
	err := conn.redisConn.Close()
	if err != nil {
		errorLog(ctx, "redisConnClose", err)
	}
	return err
}

func (conn *Conn) Ping(ctx context.Context) error {
	reply, err := String(conn.do(ctx, "ping"))
	if err != nil {
		return err
	}
	if reply != "PONG" {
		return errors.New("ping result not PONG")
	}
	return nil
}

/*
*
err示例：

	WRONGTYPE Operation against a key holding the wrong kind of value redis key 对应的Value 类型错误
*/
func (conn *Conn) do(ctx context.Context, commandName string, args ...interface{}) (reply interface{}, err error) {
	if !conn.opentracing {
		reply, err = conn.redisConn.Do(commandName, args...)
		if err != nil {
			errorLog(ctx, "redisConnDo", err, commandName, args)
		}
		return reply, err
	}
	operationName := fmt.Sprintf("Redis %s", commandName)
	span, _ := opentracing.StartSpanFromContext(ctx, operationName)
	defer span.Finish()
	ext.Component.Set(span, "redis")
	span.LogFields(log.Object("args", args))
	reply, err = conn.redisConn.Do(commandName, args...)
	if err != nil {
		errorLog(ctx, "redisConnDo", err, commandName, args)
		ext.Error.Set(span, true)
		span.LogKV("event", "error")
		span.LogKV("error.kind", "redis")
		span.LogKV("error.object", err.Error())
		span.LogKV("message", fmt.Sprintf("%v", err))
	}
	return reply, err
}

func (conn *Conn) Exists(ctx context.Context, key string) (reply int64, err error) {
	return Int64(conn.do(ctx, "EXISTS", key))
}

func (conn *Conn) Expire(ctx context.Context, key string, seconds int64) (reply int64, err error) {
	return Int64(conn.do(ctx, "EXPIRE", key, seconds))
}

/*
*
返回key所存储的value的数据结构类型，它可以返回string, list, set, zset 和 hash等不同的类型。
*/
func (conn *Conn) Type(ctx context.Context, key string) (reply string, err error) {
	return String(conn.do(ctx, "TYPE", key))
}

/*
*
返回key剩余的过期时间(秒数)
如果key不存在或者已过期，返回 -2
如果key存在并且没有设置过期时间（永久有效），返回 -1
*/
func (conn *Conn) TTL(ctx context.Context, key string) (reply int64, err error) {
	return Int64(conn.do(ctx, "TTL", key))
}

/*
*
将键key设定为指定的“字符串”值。
如果 key 已经保存了一个值，那么这个操作会直接覆盖原来的值，并且忽略原始类型。
当set命令执行成功之后，之前设置的过期时间都将失效
*/
func (conn *Conn) Set(ctx context.Context, key string, value interface{}) (reply interface{}, err error) {
	return conn.do(ctx, "SET", key, value)
}

/*
*
只有键key不存在的时候才会将键key设定为指定的“字符串”值。
seconds

	过期时间 单位秒

reply

	是否成功添加

err

	添加成功返回nil 添加失败返回具体err
*/
func (conn *Conn) SetNxEx(ctx context.Context, key string, value interface{}, seconds int64) (reply bool, err error) {
	re, err := String(conn.do(ctx, "SET", key, value, "nx", "ex", seconds))
	if err != nil {
		return false, err
	}
	if re == "OK" {
		return true, nil
	}
	errorLog(ctx, "SetNxEx", errors.New(fmt.Sprintf("re not ok is:%s", re)), "key:", key, "value:", value, "seconds:", seconds)
	return false, errors.New(re)
}

/*
*
将键key设定为指定的“字符串”值。
如果 key 已经保存了一个值，那么这个操作会直接覆盖原来的值，并且忽略原始类型
seconds

	过期时间 单位秒

reply

	是否成功添加

err

	添加成功返回nil 添加失败返回具体err
*/
func (conn *Conn) SetEx(ctx context.Context, key string, value interface{}, seconds int64) (reply bool, err error) {
	re, err := String(conn.do(ctx, "SET", key, value, "ex", seconds))
	if err != nil {
		return false, err
	}
	if re == "OK" {
		return true, nil
	}
	errorLog(ctx, "SetEx", errors.New(fmt.Sprintf("re not ok is:%s", re)), "key:", key, "value:", value, "seconds:", seconds)
	return false, errors.New(re)
}

/*
*
返回key的value。
err

	如果key不存在，返回 redisz.ErrNil
	如果key的value不是string(redis value的string 类型)，就返回错误，因为GET只处理string类型的values。
*/
func (conn *Conn) Get(ctx context.Context, key string) (reply interface{}, err error) {
	reply, err = conn.do(ctx, "GET", key)
	if err != nil {
		return reply, err
	}
	switch reply.(type) {
	case nil:
		return reply, redis.ErrNil
	}
	return reply, err
}

func (conn *Conn) GetString(ctx context.Context, key string) (reply string, err error) {
	return String(conn.Get(ctx, key))
}

func (conn *Conn) GetInt64(ctx context.Context, key string) (reply int64, err error) {
	return Int64(conn.Get(ctx, key))
}

/*
*
将key对应的数字减decrement。如果key不存在，操作之前，key就会被置为0。
reply

	返回一个数字：减少之后的value值

err

	如果key的value类型错误或者是个不能表示成数字的字符串，就返回错误
*/
func (conn *Conn) Decr(ctx context.Context, key string) (reply int64, err error) {
	return Int64(conn.do(ctx, "DECR", key))
}

/*
*
将key对应的数字减decrement。如果key不存在，操作之前，key就会被置为0。
reply

	返回一个数字：减少之后的value值

err

	如果key的value类型错误或者是个不能表示成数字的字符串，就返回错误
*/
func (conn *Conn) DecrBy(ctx context.Context, key string, decrement int64) (reply int64, err error) {
	return Int64(conn.do(ctx, "DECRBY", key, decrement))
}

/*
*
对存储在指定key的数值执行原子的加1操作。如果指定的key不存在，那么在执行incr操作之前，会先将它的值设定为0。
reply

	执行递增操作后key对应的值

err

	如果指定的key中存储的值不是字符串类型,返回的错误 :WRONGTYPE Operation against a key holding the wrong kind of value
	存储的字符串类型不能表示为一个整数,返回的错误:ERR value is not an integer or out of range
*/
func (conn *Conn) Incr(ctx context.Context, key string) (reply int64, err error) {
	return Int64(conn.do(ctx, "INCR", key))
}

/*
*
将key对应的数字加decrement。如果key不存在，操作之前，key就会被置为0
reply

	增加之后的value值

err

	如果指定的key中存储的值不是字符串类型,返回的错误 :WRONGTYPE Operation against a key holding the wrong kind of value
	存储的字符串类型不能表示为一个整数,返回的错误:ERR value is not an integer or out of range
*/
func (conn *Conn) IncrBy(ctx context.Context, key string, decrement int64) (reply int64, err error) {
	return Int64(conn.do(ctx, "INCRBY", key, decrement))
}

/*
*
添加一个或多个指定的member元素到集合的 key中.指定的一个或者多个元素member 如果已经在集合key中存在则忽略.如果集合key 不存在，则新建集合key,并添加member元素到集合key中
reply返回新成功添加到集合里元素的数量，不包括已经存在于集合中的元素.
err 如果key 的类型不是集合则返回错误
*/
func (conn *Conn) SAdd(ctx context.Context, key string, member ...interface{}) (reply int64, err error) {
	return Int64(conn.do(ctx, "SADD", append([]interface{}{key}, member...)...))
}

/*
*
count

	集合的基数(元素的数量),如果key不存在,则返回 0.
*/
func (conn *Conn) SCard(ctx context.Context, key string) (count int64, err error) {
	return Int64(conn.do(ctx, "SCARD", key))
}

/*
*
将多个set key 取交集 保存到 destinationKey 里
count

	结果集中成员的个数
*/
func (conn *Conn) SInterStore(ctx context.Context, destinationKey interface{}, targetKey ...interface{}) (count int64, err error) {
	return Int64(conn.do(ctx, "SINTERSTORE", append([]interface{}{destinationKey}, targetKey...)...))
}

// SDiffStore 该命令类似于 SDiff, 不同之处在于该命令不返回结果集，而是将结果存放在destination集合中.
// 如果destination已经存在, 则将其覆盖重写
// 返回结果集元素的个数
func (conn *Conn) SDiffStore(ctx context.Context, destinationKey interface{}, targetKey ...interface{}) (count int64, err error) {
	return Int64(conn.do(ctx, "SDIFFSTORE", append([]interface{}{destinationKey}, targetKey...)...))
}

/*
*
返回成员 member 是否是存储的集合 key的成员.
count

	如果member元素是集合key的成员，则返回true
	如果member元素不是key的成员，或者集合key不存在，则返回false
*/
func (conn *Conn) SIsMember(ctx context.Context, key string, member interface{}) (reply bool, err error) {
	re, err := Int64(conn.do(ctx, "SISMEMBER", key, member))
	if err != nil {
		return false, err
	}
	if re == 1 {
		return true, nil
	}
	return false, nil
}

/*
*
在key集合中移除指定的元素. 如果指定的元素不是key集合中的元素则忽略 如果key集合不存在则被视为一个空的集合，该命令返回0.
reply

	从集合中移除元素的个数，不包括不存在的成员

err

	如果key的类型不是一个集合,则返回错误.
*/
func (conn *Conn) SRem(ctx context.Context, key string, member ...interface{}) (reply int64, err error) {
	return Int64(conn.do(ctx, "SREM", append([]interface{}{key}, member...)...))
}

type ScanReply struct {
	Iterator int64
	Reply    []interface{}
}

func (conn *Conn) Scan(ctx context.Context, iterator int64, count int64, match string) (reply ScanReply, err error) {
	arg := []interface{}{iterator, "count", count}
	if match != "" {
		arg = append(arg, "match", match)
	}
	values, err := Values(conn.do(ctx, "SCAN", arg...))
	if err != nil {
		return ScanReply{}, err
	}
	ite, _ := Int64(values[0], nil)
	rep, _ := Values(values[1], nil)
	return ScanReply{
		Iterator: ite,
		Reply:    rep,
	}, nil
}

/*
*
reply

	iterator-命令返回了游标 0 ， 这表示迭代已经结束， 整个数据集已经被完整遍历过了。

注意事项:

	同一个元素可能会被返回多次。 处理重复元素的工作交由应用程序负责， 比如说， 可以考虑将迭代返回的元素仅仅用于可以安全地重复执行多次的操作上。
	如果一个元素是在迭代过程中被添加到数据集的， 又或者是在迭代过程中从数据集中被删除的， 那么这个元素可能会被返回， 也可能不会。
*/
func (conn *Conn) SScan(ctx context.Context, key string, iterator int64, count int64) (reply ScanReply, err error) {
	values, err := Values(conn.do(ctx, "SSCAN", key, iterator, "count", count))
	if err != nil {
		return ScanReply{}, err
	}
	ite, _ := Int64(values[0], nil)
	rep, _ := Values(values[1], nil)
	return ScanReply{
		Iterator: ite,
		Reply:    rep,
	}, nil
}

/*--------------------------------------------Lists--------------------------------------------*/

/*
*
将所有指定的值插入到存于 key 的列表的头部。如果 key 不存在，那么在进行 push 操作前会创建一个空列表。
可以把多个元素 push 进入列表，只需在命令末尾加上多个指定的参数。元素是从最左端的到最右端的、一个接一个被插入到 list 的头部。

	所以对于这个命令例子 LPUSH mylist a b c，返回的列表是 c 为第一个元素， b 为第二个元素， a 为第三个元素。

reply

	在 push 操作后的 list 长度

err

	如果 key 对应的值不是一个 list 的话，那么会返回一个错误。
*/
func (conn *Conn) LPush(ctx context.Context, key string, member ...interface{}) (reply int64, err error) {
	return Int64(conn.do(ctx, "LPUSH", append([]interface{}{key}, member...)...))
}

/*
*
向存于 key 的列表的尾部插入所有指定的值。如果 key 不存在，那么会创建一个空的列表然后再进行 push 操作。
可以把多个元素打入队列，只需要在命令后面指定多个参数。元素是从左到右一个接一个从列表尾部插入。

	比如命令 RPUSH mylist a b c 会返回一个列表，其第一个元素是 a ，第二个元素是 b ，第三个元素是 c。

reply

	在 push 操作后的 list 长度

err

	当 key 保存的不是一个列表，那么会返回一个错误。
*/
func (conn *Conn) RPush(ctx context.Context, key string, member ...interface{}) (reply int64, err error) {
	return Int64(conn.do(ctx, "RPUSH", append([]interface{}{key}, member...)...))
}

// LPop 移除并且返回 key 对应的 list 的第一个元素 list数据为空时 返回 err redis.ErrNil
func (conn *Conn) LPop(ctx context.Context, key string) (reply interface{}, err error) {
	return conn.do(ctx, "LPOP", key)
}

/*
*
移除并返回存于 key 的 list 的最后一个元素。
最后一个元素的值，或者当 key 不存在的时候返回 nil。
*/
func (conn *Conn) RPop(ctx context.Context, key string) (reply interface{}, err error) {
	return conn.do(ctx, "RPOP", key)
}

/*
*
移除并且返回 key 对应的 list 的第一个元素。当给定列表内没有任何元素可供弹出的时候， 连接将被阻塞timeout(秒)时长
timeout

	参数表示的是一个指定阻塞的最大秒数的整型值。当 timeout 为 0 是表示阻塞时间无限制。
	若经过了指定的 timeout 仍没有值,返回nil

reply

	当没有元素的时候会弹出一个 nil 的多批量值，并且 timeout 过期。
	当有元素弹出时会返回一个双元素的多批量值，其中第一个元素是弹出元素的 key，第二个元素是 value。
*/
func (conn *Conn) BLPop(ctx context.Context, key string, timeout int64) (reply []interface{}, err error) {
	return Values(conn.do(ctx, "BLPOP", key, timeout))
}

/*
*
移除并返回存于 key 的 list 的最后一个元素。当给定列表内没有任何元素可供弹出的时候， 连接将被阻塞timeout(秒)时长
timeout

	参数表示的是一个指定阻塞的最大秒数的整型值。当 timeout 为 0 是表示阻塞时间无限制。
	若经过了指定的 timeout 仍没有值,返回nil

reply

	当没有元素的时候会弹出一个 nil 的多批量值，并且 timeout 过期。
	当有元素弹出时会返回一个双元素的多批量值，其中第一个元素是弹出元素的 key，第二个元素是 value。
*/
func (conn *Conn) BRPop(ctx context.Context, key string, timeout int64) (reply []interface{}, err error) {
	return Values(conn.do(ctx, "BRPOP", key, timeout))
}

/*
*
返回存储在 key 里的list的长度。 如果 key 不存在，那么就被看作是空list，并且返回长度为 0。
err

	当存储在 key 里的值不是一个list的话，会返回error
*/
func (conn *Conn) LLen(ctx context.Context, key string) (reply int64, err error) {
	return Int64(conn.do(ctx, "LLEN", key))
}

/*
*
修剪(trim)一个已存在的 list，这样 list 就会只包含指定范围的指定元素。

start 和 stop 都是由0开始计数的， 这里的 0 是列表里的第一个元素（表头），1 是第二个元素，-1 表示列表里的最后一个元素， -2 表示倒数第二个，以此类推。

例如： LTRIM foobar 0 2 将会对存储在 foobar 的列表进行修剪，只保留列表里的前3个元素。

超过范围的下标并不会产生错误：如果 start 超过列表尾部，或者 start > stop，结果会是列表变成空表（即该 key 会被移除）。 如果 stop 超过列表尾部，Redis 会将其当作列表的最后一个元素。

reply

	返回"OK"
*/
func (conn *Conn) LTrim(ctx context.Context, key string, start, stop int64) (reply string, err error) {
	return String(conn.do(ctx, "LTRIM", key, start, stop))
}

/*--------------------------------------------sorted sets--------------------------------------------*/

/*
*
将指定成员添加到键为key有序集合（sorted set）里面。 本方法只支持个分数/成员（score/member）对。

如果指定添加的成员已经是有序集合里面的成员，则会更新改成员的分数（scrore）并更新到正确的排序位置。
如果key不存在，将会创建一个新的有序集合（sorted set）并将分数/成员（score/member）对添加到有序集合，就像原来存在一个空的有序集合一样。

score

	分数值是一个双精度的浮点型数字字符串。+inf和-inf都是有效值。

options 支持的参数

	XX: 仅仅更新存在的成员，不添加新成员。
	NX: 不更新存在的成员。只添加新成员。
	CH: 修改返回值为发生变化的成员总数，原始是返回新添加成员的总数 (CH 是 changed 的意思)。更改的元素是新添加的成员，已经存在的成员更新分数。
		所以在命令中指定的成员有相同的分数将不被计算在内。注：在通常情况下，ZADD返回值只计算新添加成员的数量。
	INCR: 当ZADD指定这个选项时，成员的操作就等同ZINCRBY命令，对成员的分数进行递增操作。

reply

	默认返回：添加到有序集合的成员数量，不包括已经存在更新分数的成员。
	options存在CH时，返回所有发生变化的成员数量，包括已经存在更新分数的成员。

err

	如果key存在，但是类型不是有序集合，将会返回一个错误应答。
*/
func (conn *Conn) ZAdd(ctx context.Context, key string, score int64, member interface{}, options ...interface{}) (reply interface{}, err error) {
	arg := []interface{}{key}
	if len(options) > 0 {
		arg = append(arg, options...)
	}
	arg = append(arg, score, member)

	return conn.do(ctx, "ZADD", arg...)
}

/*
*
返回有序集key中，score值在min和max之间(默认包括score值等于min或max)的成员个数
*/
func (conn *Conn) ZCount(ctx context.Context, key string, min, max int64) (reply int64, err error) {
	return Int64(conn.do(ctx, "ZCOUNT", key, min, max))
}

/*
*
返回key的有序集元素个数。
*/
func (conn *Conn) ZCard(ctx context.Context, key string) (reply int64, err error) {
	return Int64(conn.do(ctx, "ZCARD", key))
}

/*
*
为有序集key的成员member的score值加上增量increment。
如果key中不存在member，就在key中添加一个member，score是increment（就好像它之前的score是0.0）
如果key不存在，就创建一个只含有指定member成员的有序集合。
reply

	member成员的新score值，以字符串形式表示。
*/
func (conn *Conn) ZIncrBy(ctx context.Context, key string, increment, member interface{}) (reply interface{}, err error) {
	return conn.do(ctx, "ZINCRBY", key, increment, member)
}

/*
*
返回存储在有序集合key中的指定范围的元素。 返回的元素可以认为是按得分从最低到最高排列。 如果得分相同，将按字典排序。
参数start和stop都是基于零的索引，即0是第一个元素，1是第二个元素，以此类推。 它们也可以是负数，表示从有序集合的末尾的偏移量，其中-1是有序集合的最后一个元素，-2是倒数第二个元素，等等。
start和stop都是全包含的区间，因此例如ZRANGE myzset 0 1将会返回有序集合的第一个和第二个元素。
options

	支持WITHSCORES,以便将元素的分数与元素一起返回。这样，返回的列表将包含value1,score1,...,valueN,scoreN

reply

	默认：给定范围内的元素列表 value1,value2,...,valueN 切片 可以使用redis.Strings 或者redis.Int64s 方法 返回对应类型的切片
	options传入WITHSCORES,返回value1,score1,...,valueN,scoreN 切片，可以使用redis.Int64Map 或者redis.StringMap 等方法 把返回值 转换成map 类型  key 为成员 value 为 分数
*/
func (conn *Conn) ZRange(ctx context.Context, key string, start, stop int64, options ...interface{}) (reply interface{}, err error) {
	return conn.do(ctx, "ZRANGE", append([]interface{}{key, start, stop}, options...)...)
}

/*
*
返回有序集key中，指定区间内的成员。其中成员的位置按score值递减(从大到小)来排列。具有相同score值的成员按字典序的反序排列。
除了成员按score值递减的次序排列这一点外，ZRevRange命令的其他方面和ZRange命令一样。
*/
func (conn *Conn) ZRevRange(ctx context.Context, key string, start, stop int64, options ...interface{}) (reply interface{}, err error) {
	return conn.do(ctx, "ZREVRANGE", append([]interface{}{key, start, stop}, options...)...)
}

/*
*
删除一个或者多个成员
reply

	从有序集合中删除的成员个数，不包括不存在的成员
*/
func (conn *Conn) ZRem(ctx context.Context, key string, member ...interface{}) (reply int64, err error) {
	return Int64(conn.do(ctx, "ZREM", append([]interface{}{key}, member...)...))
}

/*
*
返回有序集key中成员member的排名，其中有序集成员按score值从大到小排列。排名以0为底，也就是说，score值最大的成员排名为0。
*/
func (conn *Conn) ZRevRank(ctx context.Context, key string, member interface{}) (reply int64, err error) {
	return Int64(conn.do(ctx, "ZREVRANK", key, member))
}

/*
*
返回有序集key中成员member的排名。其中有序集成员按score值递增(从小到大)顺序排列。排名以0为底，也就是说，score值最小的成员排名为0。
*/
func (conn *Conn) ZRank(ctx context.Context, key string, member interface{}) (reply int64, err error) {
	return Int64(conn.do(ctx, "ZRANK", key, member))
}

/*
*
返回有序集key中，成员member的score值。
err

	如果member元素不是有序集key的成员，或key不存在，返回redis.ErrNil
*/
func (conn *Conn) ZScore(ctx context.Context, key string, member interface{}) (reply int64, err error) {
	return Int64(conn.do(ctx, "ZSCORE", key, member))
}

/*--------------------------------------------hashes--------------------------------------------*/
/**
设置 key 指定的哈希集中指定字段的值。

如果 key 指定的哈希集不存在，会创建一个新的哈希集并与 key 关联。

如果字段在哈希集中存在，它将被重写。
reply
	1如果field是一个新的字段
	0如果field原来在map里面已经存在
*/
func (conn *Conn) HSet(ctx context.Context, key string, field, value interface{}) (reply int64, err error) {
	return Int64(conn.do(ctx, "HSET", key, field, value))
}

/*
*
只在 key 指定的哈希集中不存在指定的字段时，设置字段的值。如果 key 指定的哈希集不存在，会创建一个新的哈希集并与 key 关联。如果字段已存在，该操作无效果。
reply

	1：如果字段是个新的字段，并成功赋值
	0：如果哈希集中已存在该字段，没有操作被执行
*/
func (conn *Conn) HSetNX(ctx context.Context, key string, field, value interface{}) (reply int64, err error) {
	return Int64(conn.do(ctx, "HSETNX", key, field, value))
}

/*
*
返回 key 指定的哈希集中该字段所关联的值
reply

	当字段不存在或者 key 不存在时返回nil
*/
func (conn *Conn) HGet(ctx context.Context, key string, member interface{}) (reply interface{}, err error) {
	return conn.do(ctx, "HGET", key, member)
}

/*
*
增加 key 指定的哈希集中指定字段的数值。如果 key 不存在，会创建一个新的哈希集并与 key 关联。如果字段不存在，则字段的值在该操作执行前被设置为 0
HIncrBy 支持的值的范围限定在 64位 有符号整数
reply

	增值操作执行后的该字段的值
*/
func (conn *Conn) HIncrBy(ctx context.Context, key string, field interface{}, increment int64) (reply int64, err error) {
	return Int64(conn.do(ctx, "HINCRBY", key, field, increment))
}

/*
*
返回 key 指定的哈希集包含的字段的数量。

	哈希集中字段的数量，当 key 指定的哈希集不存在时返回 0
*/
func (conn *Conn) HLen(ctx context.Context, key string) (reply int64, err error) {
	return Int64(conn.do(ctx, "HLEN", key))
}
