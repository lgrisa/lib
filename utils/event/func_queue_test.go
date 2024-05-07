package event

//import (
//	"fmt"
//	. "github.com/onsi/gomega"
//	"testing"
//	"time"
//)
//
//func TestFuncQueue(t *testing.T) {
//	RegisterTestingT(t)
//
//	q := NewFuncQueue(20, "test")
//
//	q.TryFunc(func() {
//		time.Sleep(100*time.Millisecond)
//	})
//
//	cs := make(chan struct{})
//	for i := 1; i < 40; i++ {
//		idx := i
//		q.MustFunc(func() {
//			fmt.Println(idx)
//
//			if idx == 34 {
//				cs <- struct{}{}
//				time.Sleep(100*time.Millisecond)
//			}
//		})
//	}
//
//	Ω(q.TryFunc(func() {})).Should(BeFalse())
//
//	Ω(len(q.funcQueue)).Should(BeEquivalentTo(20))  // 1-19
//	Ω(q.funcCache.Len()).Should(BeEquivalentTo(20)) // 20 - 39
//
//	<-cs
//	// 20-35
//	// 36-39
//
//	Ω(len(q.funcQueue)).Should(BeEquivalentTo(1))  // 35
//	Ω(q.funcCache.Len()).Should(BeEquivalentTo(4)) // 36-39
//
//	for i := 40; i < 70; i++ {
//		idx := i
//		q.MustFunc(func() {
//			fmt.Println(idx)
//
//			if idx == 53 {
//				cs <- struct{}{}
//				time.Sleep(100*time.Millisecond)
//			}
//
//			if idx == 68 {
//				cs <- struct{}{}
//				time.Sleep(100*time.Millisecond)
//			}
//		})
//	}
//
//	Ω(len(q.funcQueue)).Should(BeEquivalentTo(20))  // 35-54
//	Ω(q.funcCache.Len()).Should(BeEquivalentTo(15)) // 55-69
//
//	<-cs
//	//Ω(len(q.funcQueue)).Should(BeEquivalentTo(5)) // 54
//	//Ω(q.funcCache.Len()).Should(BeEquivalentTo(15)) // 55-69
//
//	<-cs
//	Ω(len(q.funcQueue)).Should(BeEquivalentTo(1))  // 68
//	Ω(q.funcCache.Len()).Should(BeEquivalentTo(0)) //  69
//
//	q.Close(false)
//}
