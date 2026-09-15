module github.com/qianwj/typed/concurrency

go 1.27.1

require github.com/qianwj/typed/adt v0.0.0

require github.com/qianwj/typed/utils v0.0.0 // indirect

replace (
	github.com/qianwj/typed/adt => ../adt
	github.com/qianwj/typed/utils => ../utils
)
