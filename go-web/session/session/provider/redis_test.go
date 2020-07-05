package provider

import "testing"

var redisProvider Provider

func getRedisProvider() Provider {
	if redisProvider == nil {
		redisProvider = NewRedisProvider("127.0.0.1:6379", "")
	}
	return redisProvider
}

func TestRedisProvider_SessionInit(t *testing.T) {
	testKey := "testKey"
	testVal := "testValue"
	pr := getRedisProvider()
	if pr == nil {
		t.Fatal("RedisProvider: new redis provider error")
	}

	sess, err := pr.SessionInit(testVal)
	if err != nil {
		t.Fatal("RedisProvider: session init error")
	}
	err = sess.Set(testKey, testVal)
	if err != nil {
		t.Fatal("RedisProvider: session set error when session init")
	}

	gotVal, err := sess.Get(testKey)
	if err != nil || gotVal.(string) != testVal {
		t.Fatal("RedisProvider: session get error when session init")
	}

	err = sess.Delete(testKey)
	if err != nil {
		t.Fatal("RedisProvider: session delete key error when session init")
	}
}

func TestRedisProvider_SessionRead(t *testing.T) {
	testKey := "testKey"
	testVal := "testValue"
	pr := getRedisProvider()
	if pr == nil {
		t.Fatal("RedisProvider: new redis provider error")
	}

	sess, err := pr.SessionInit(testKey)
	if err != nil {
		t.Fatal("RedisProvider: session init error")
	}
	err = sess.Set(testKey, testVal)
	if err != nil {
		t.Fatal("RedisProvider: session set error when session read")
	}
	sess, err = pr.SessionRead(testKey)
	if err != nil {
		t.Fatal("RedisProvider: session read error")
	}
	gotVal, err := sess.Get(testKey)
	if err != nil || gotVal.(string) != testVal {
		t.Fatal("RedisProvider: session get error when session read")
	}
}

func TestRedisProvider_SessionDestroy(t *testing.T) {
	testKey := "testKey"
	pr := getRedisProvider()
	if pr == nil {
		t.Fatal("RedisProvider: new redis provider error")
	}

	_, err := pr.SessionInit(testKey)
	if err != nil {
		t.Fatal("RedisProvider: session init error")
	}
	if err := pr.SessionDestroy(testKey); err != nil {
		t.Fatal("RedisProvider: session destroy error")
	}
}
