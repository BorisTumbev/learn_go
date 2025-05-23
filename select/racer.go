package racer

import (
	"fmt"
	"net/http"
	"time"
)

// func Racer(a, b string) (winner string) {
// 	aDuration := measureResponseTime(a)
// 	bDuration := measureResponseTime(b)

// 	if aDuration < bDuration {
// 		return a
// 	}

// 	return b
// }

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
}

var tenSecondTimeout = 10 * time.Second

func Racer(a, b string) (winner string, error error) {
	return ConfigurableRacer(a, b, tenSecondTimeout)
}

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}

func ping(url string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		http.Get(url)
		close(ch)
	}()
	return ch
}

// Python equivalent, it's interesting one

// import asyncio
// import aiohttp

// async def ping(session, url):
//     try:
//         async with session.get(url) as response:
//             await response.read()  # ensure the response completes
//     except Exception:
//         pass  # treat errors like Go's version: just ignore
//     return url

// async def configurable_racer(a, b, timeout):
//     async with aiohttp.ClientSession() as session:
//         try:
//             # Schedule both tasks concurrently
//             tasks = [ping(session, a), ping(session, b)]
//             done, pending = await asyncio.wait(tasks, timeout=timeout, return_when=asyncio.FIRST_COMPLETED)

//             if not done:
//                 raise asyncio.TimeoutError(f"Timed out waiting for {a} and {b}")

//             # Cancel the slower task
//             for task in pending:
//                 task.cancel()

//             # Return the result of the first completed task
//             winner = list(done)[0].result()
//             return winner

//         except asyncio.TimeoutError as e:
//             return str(e)

// def racer(a, b):
//     return asyncio.run(configurable_racer(a, b, timeout=10))
