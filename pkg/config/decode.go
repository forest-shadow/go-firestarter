package config

import "github.com/go-viper/mapstructure/v2"

func DecodeHook() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToWeakSliceHookFunc(","),
		mapstructure.TextUnmarshallerHookFunc(),
	)
}
