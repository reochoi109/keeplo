package scheduler

// func TestScheduler(t *testing.T) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
// 	defer cancel()

// 	executed := 0
// 	doneCh := make(chan struct{})

// 	queue := NewInMemoryQueue()
// 	err := queue.Push(&Task{
// 		MonitorID:   "test-monitor",
// 		NextCheckAt: time.Now().Add(time.Second),
// 	})

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	scheduler := NewScheduler(queue, func(ctx context.Context, task *Task) {
// 		executed++
// 		fmt.Printf("[TEST CHECK] MonitorID: %s | Time: %s\n", task.MonitorID, time.Now().Format("15:04:05"))

// 		if executed >= 10 {
// 			// 일정 횟수 후 테스트 종료
// 			close(doneCh)
// 			return
// 		}

// 		task.NextCheckAt = time.Now().Add(3 * time.Second)
// 		if err := queue.Push(task); err != nil {
// 			t.Log(err)
// 		}
// 	})

// 	go scheduler.Start(ctx)

// 	select {
// 	case <-doneCh:
// 		t.Logf("%d 실행됨", executed)
// 	case <-ctx.Done():
// 		t.Fatal("timeout")
// 	}

// 	queue.Close()
// }
