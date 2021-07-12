package errgroup

import "testing"

func TestGroup(t *testing.T) {
	eg := &Group{}
	arr := make([]int, 0)
	c := 10
	for i:=0; i<c; i++ {
		j := i
		eg.Go(func() error {
			arr = append(arr, j)
			return nil
		})
	}
	eg.Wait()
	if len(arr) == c{
		t.Log("errgroup ok", arr)
	}else {
		t.Fatal("errgroup fai", arr)
	}
	
}
