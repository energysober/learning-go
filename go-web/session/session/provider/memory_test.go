package provider

import "testing"

func getMemoryProvider() Provider {
	return NewMemoryProvider()
}

func TestMemoryProvider_SessionInit(t *testing.T) {
	testKey := "testSessionId"
	testVal := "testValue"
	pr := getMemoryProvider()
	sess, err := pr.SessionInit(testKey)
	if err != nil {
		t.Fatal("MemoryProvider: session init error")
	}

	err = sess.Set(testKey, "testValue")
	if err != nil {
		t.Fatal("MemoryProvider: session set error")
	}

	gotValue, err := sess.Get(testKey)
	if err != nil || gotValue.(string) != testVal {
		t.Fatal("MemoryProvider: session get error")
	}

	gotKey := sess.SessionId()
	if gotKey != testKey {
		t.Fatal("MemoryProvider: session id error, want testSessionId not " + gotKey)
	}

	err = sess.Delete(testKey)
	if err != nil {
		t.Fatal("MemoryProvider: delete session error")
	}
}

func TestMemoryProvider_SessionRead(t *testing.T) {
	testKey := "testSessionId"
	testVal := "testValue"
	pr := getMemoryProvider()
	sess, err := pr.SessionInit(testKey)
	if err != nil {
		t.Fatal("MemoryProvider: session init error")
	}

	err = sess.Set(testKey, "testValue")
	if err != nil {
		t.Fatal("MemoryProvider: session set error")
	}
	sess, err = pr.SessionRead(testKey)
	if err != nil {
		t.Fatal("MemoryProvider: session read error")
	}
	gotValue, err := sess.Get(testKey)
	if err != nil || gotValue.(string) != testVal {
		t.Fatal("MemoryProvider: sess get error, want testSessionId but get " + gotValue.(string))
	}
}

func TestMemoryProvider_SessionDestroy(t *testing.T) {
	testKey := "testSessionId"
	pr := getMemoryProvider()
	sess, err := pr.SessionInit(testKey)
	if err != nil {
		t.Fatal("MemoryProvider: session init error")
	}
	if err := pr.SessionDestroy(testKey); err != nil {
		t.Fatal("MemoryProvider: session destroy error")
	}
	_, err = sess.Get(testKey)
	if err == nil {
		t.Fatal("MemoryProvider: session destroy failed, the session still exist")
	}
}
