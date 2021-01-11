## Performance

There is a benchmark test which compares using JSONry and `encoding/json` to
marshal and unmarshal the same data. The value of the benchmark is in the
relative performance between the two. In the benchamark test:

- Marshal
  - JSONry takes 3 times as long as `encoding/json`
  - JSONry allocates 7 times as much memory as `enconding/json`

- Unmarshal
  - JSONry takes 13 times as long as `encoding/json`
  - JSONry allocates 24 times as much memory as `enconding/json`

The command used to run the benchmark tests is:
```bash
go test -run none -bench . -benchmem -benchtime 10s
```
