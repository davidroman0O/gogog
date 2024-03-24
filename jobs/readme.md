# Middleware jobs

I'm also here to have fun so I'm building a job middleware for sqlite3, it is not intended to be fast just toensure i have enough control of the flow of operations that i will certainly have to retry!

I just want a very easy to use job lib for that app without adding MORE dependencies

# Database lock fixies

If you're having database lock (with `:memory:`) with sqlite3 while using hooks, take a look at my gist https://gist.github.com/davidroman0O/8da76d79364c559b98ba969c29cc969f 

You can't re-use the same connection and have to leverage the `conn` value instead within your hooks while using `db.SetMaxOpenConns(1)`.

I will try to support both file and memory modes with `db` and `jobs`, maybe i will one day make a repo with multiple tools for sqlite3 

# TODO / wish list


- all dependency injection style configuration
- worker based handling: workers are consuming
    - consumers signature should be `func[T any](ctx context.Context, data T) error`
    - create a new worker and add consumers inside`jobs.NewWorker(opts ...workerOptions) error`
    - `jobs.WorkerWithConsumer[T any](func (data T) error) workerOptions` 
- client based creation:
    - producer `func Enqueue[T any](data T, opts ...enenqueOptions) error`, by default it will create into `default` queue
    - `jobs.EnqueueWithQueue(value string) enenqueOptions`

I also want to be able to have dependencies between jobs and a saga pattern to handle orchestration

