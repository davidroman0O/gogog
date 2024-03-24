# Middleware jobs

I'm also here to have fun so I'm building a job middleware for sqlite3, it is not intended to be fast just toensure i have enough control of the flow of operations that i will certainly have to retry!

If you're having database lock (with `:memory:`) with sqlite3 while using hooks, take a look at my gist https://gist.github.com/davidroman0O/8da76d79364c559b98ba969c29cc969f 

You can't re-use the same connection and have to leverage the `conn` value instead within your hooks while using `db.SetMaxOpenConns(1)`.

I will try to support both file and memory modes with `db` and `jobs`, maybe i will one day make a repo with multiple tools for sqlite3 