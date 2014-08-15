package mphf

import "testing"

func TestMPHF(t *testing.T) {
	testEntries := make([]KeyValue, 0, 10)
	testEntries = append(testEntries, KeyValue{[]byte("test1"), "Testing 1"})
	testEntries = append(testEntries, KeyValue{[]byte("test2"), "Testing 2"})
	testEntries = append(testEntries, KeyValue{[]byte("test3"), "Testing 3"})
	testEntries = append(testEntries, KeyValue{[]byte("test4"), "Testing 4"})
	testEntries = append(testEntries, KeyValue{[]byte("test5"), "Testing 5"})
	testEntries = append(testEntries, KeyValue{[]byte("test6"), "Testing 6"})
	testEntries = append(testEntries, KeyValue{[]byte("test7"), "Testing 7"})
	testEntries = append(testEntries, KeyValue{[]byte("test8"), "Testing 8"})
	testEntries = append(testEntries, KeyValue{[]byte("test9"), "Testing 9"})
	testEntries = append(testEntries, KeyValue{[]byte("test10"), "Testing 10"})

	mphf, ok := BuildMPHF(testEntries)
	if ! ok {
		t.Errorf("Unable to build an MPHF")
	}

	val, ok := mphf.Get([]byte("test5"))
	if !ok {
		t.Errorf("Unable to get test5 out")
	}

	if val != "Testing 5" {
		t.Errorf("Val should've been \"Testing 5\", got %v instead", val)
	}

	val, ok = mphf.Get([]byte("test8"))
	t.Logf("Item test8: %v", val)
}
