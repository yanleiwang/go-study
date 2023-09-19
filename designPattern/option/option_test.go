package additional

import "errors"

/*
场景： 类初始化， 实现默认参数的功能
本质还是因为go不支持函数重载和默认参数
*/

type MyStructOption func(myStruct *MyStruct)

type MyStructOptionErr func(myStruct *MyStruct) error

type MyStruct struct {
	// 第一个部分是必须用户输入的字段

	id   uint64
	name string

	// 第二个部分是可选的字段
	address string
	// 这里可以有很多字段
	//field1 int
	//field2 int
}

func WithAddressV1(address string) MyStructOptionErr {
	return func(myStruct *MyStruct) error {
		if address == "" {
			return errors.New("地址不能为空字符串")
		}
		myStruct.address = address
		return nil
	}
}

func WithAddressV2(address string) MyStructOption {
	return func(myStruct *MyStruct) {
		if address == "" {
			panic("地址不能为空字符串")
		}
		myStruct.address = address
	}
}

func WithAddress(address string) MyStructOption {
	return func(myStruct *MyStruct) {
		myStruct.address = address
	}
}

// var m =MyStruct{}

// NewMyStruct 参数包含所有的必须用户指定的字段
func NewMyStruct(id uint64, name string, opts ...MyStructOption) *MyStruct {
	// 构造必传的部分
	res := &MyStruct{
		id:   id,
		name: name,
	}

	// 假如 res 返回非指针
	// for _, opt := range opts {
	// 	opt(&res)
	// }

	for _, opt := range opts {
		opt(res)
	}
	return res
}

func NewMyStructV1(id uint64, name string, opts ...MyStructOptionErr) (*MyStruct, error) {
	res := &MyStruct{
		id:   id,
		name: name,
	}

	// 假如 res 返回非指针
	// for _, opt := range opts {
	// 	opt(&res)
	// }

	for _, opt := range opts {
		if err := opt(res); err != nil {
			return nil, err
		}
	}
	return res, nil
}
