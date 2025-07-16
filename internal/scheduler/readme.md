```go

queue := scheduler.NewInMemoryQueue()
ctx, cancel := context.WithCancel(context.Background())

go func() {
	for {
		task, err := queue.Pop(ctx)
		if err != nil {
			log.Println("Pop ended:", err)
			return
		}
		// 작업 처리
	}
}()

// 종료할 때
cancel()
queue.Close()

```