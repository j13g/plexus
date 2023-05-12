package config

import (
	"github.com/samber/do"
	"github.com/samber/mo"
)

func Invoke[T any]() mo.Result[T] {
	return mo.TupleToResult[T](do.Invoke[T](Get().Injector))
}

func ProvideValue[T any](val T) {
	do.ProvideValue[T](Get().Injector, val)
}

func Provide[T any](provider do.Provider[T]) {
	do.Provide[T](Get().Injector, provider)
}
