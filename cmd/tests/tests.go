package main

import (
    "github.com/DGHeroin/LoadBalance"
    "log"
    "time"
)

func main()  {
    rb := LoadBalance.NewRoundRobin()
    rb.Put("",2, "hello")
    rb.Put("",2, "world")
    rb.Put("",2, "!")
    rb.Put("",1, "?")
    cnt := 0
    for {
        if idx, obj := rb.Get(); obj != nil {
            log.Printf("%v. get: %s => %d", cnt, obj, idx)
            cnt++
            if cnt % 3 == 0 {
                log.Println("删除: ", idx)
                rb.Remove(idx) // 每 3 秒删除一个
            }
        }
        time.Sleep(time.Second)
    }
}
