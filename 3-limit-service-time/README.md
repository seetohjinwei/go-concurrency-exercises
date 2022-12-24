# Limit Service Time for Free-tier Users

Your video processing service has a freemium model. Everyone has 10
sec of free processing time on your service. After that, the
service will kill your process, unless you are a paid premium user.

Beginner Level: 10s max per request
Advanced Level: 10s max per user (accumulated)

**Advanced Output**

```
❯ go run -race .
UserID: 0       Process 1 started.
UserID: 1       Process 2 started.
UserID: 0       Process 3 started.
UserID: 0       Process 4 started.
UserID: 1       Process 5 started.
UserID: 0       Process 1 done.
UserID: 0       Process 3 killed. (No quota left)
UserID: 0       Process 4 killed. (No quota left)
UserID: 1       Process 5 done.
UserID: 1       Process 2 done.
```
