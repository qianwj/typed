module github.com/qianwj/typed/reactivex

go 1.27.1

require github.com/qianwj/typed/adt v0.0.3

replace (
	github.com/qianwj/typed/adt => ../adt
	github.com/qianwj/typed/utils => ../utils
)
