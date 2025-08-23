# timeline-fan-out-worker-go-exmpl
An example timeline fan out service with folanf using my-sql and redis

## Run the app

### Use dokcer compose :
```bash
#Debug
docker-compose up --build app0


# Production
docker-compose up --build -d 
```

### Run tests: 
```bash
go test ./...
```

### Important Env Variables : 

```yaml
      - DB_HOST=<database host address in network>
      - DB_USER=admin
      - DB_PASSWORD=admin
      - DB_NAME=timeline_db
      - REDIS_ADDR=redis:6379
      - APP_PORT=:8080 
```

## Api Refrence

### `POST /posts/`
get's a json in body : 
```json
{
    "sender_id": <SENDER_USER_ID>,
    "content": "test content"
}
```
example curl test api : 
```bash
curl -X POST http://localhost:8080/posts   -H "Content-Type: application/json"   -d '{
    "sender_id": 101,
    "content": "test content"
  }'

```

> این api باید با یه روش authorization در ابتدا کلاینتی که پست ثبت میکرد را احراز میکرد و فیلد sender_id را مستقیما از کاربر نمیگرفت.   
> منتهی چون پروژه تستی بود و در شرح پروژه نیامده بود، این مسئله رو فقط همینجا بیان میکنیم و بعدا در صورت نیاز به پروژه اضافه میکنیم.

---

### `GET /timeline`   
get's query strings and return's a json data in this schema : 

- `userId`: the subscriber user id
- `limit`: limit for pagination (default value is 10) 
- `offset`: index for start the pagination (default value is 0) 

```json
[
    ...

    {
        "id":1897,
        "sender_id":101,
        "content":"test content",
        "created_at":"2025-08-22T23:25:44.573678Z"
    },
    
    ...

]
```

example curl test:
```bash
curl -X GET http://localhost:8080/timeline?userId=202 -H "Content-Type: application/json"
```

> در اینجا هم باید با توکنی چیزی کاربر را احراز هویت میکردیم نه اینکه id کاربر سابسکرایبر را در query ها میگرفتیم منتهی چون پروژه تستی بود و در شرح پروژه نیامده بود فعلا همینجا بیان میکنیم و بعدا در صورت نیاز به آن اضافه میکنیم


---

## Structure and architecture

The core architecture combines a Fan-Out Worker pattern with a Worker Pool. Since scalability was a key requirement, I decided to implement a simple Worker Pool, which allows us in the future to run multiple instances of app0 and improve service delivery. Additionally, by adjusting the maximum number of workers in the background, we can increase the number of posts processed for each user’s timeline concurrently.    
I also focused on minimizing coupling in the project’s design and maintaining strong package isolation.


Thank you.



