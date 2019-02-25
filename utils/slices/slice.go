package slices

import (
	"math/rand"
	"reflect"
	"time"
)

// 数据操作
type reducetype func(interface{}) interface{}

// 条件过滤
type filtertype func(interface{}) bool

func Slice_randList(min, max int) []int {
	if max < min {
		min, max = max, min
	}
	length := max - min + 1
	t0 := time.Now()
	rand.Seed(int64(t0.Nanosecond()))
	list := rand.Perm(length)
	for index, _ := range list {
		list[index] += min
	}
	return list
}

// 去除（有序）列表重复元素
func Slice_duplicate(slice interface{}) (ret []interface{}) {
	va := reflect.ValueOf(slice)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface())
	}
	return ret
}

// 切片拼接
func Slice_merge(slice1, slice2 []interface{}) []interface{} {
	return append(slice1, slice2...)
}

// 元素val是否在切片中
func Slice_in(val interface{}, slice []interface{}) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// 对切片中的数据进行操作（例如 加减乘除...）
func Slice_reduce(slice []interface{}, a reducetype) (dslice []interface{}) {
	for _, v := range slice {
		dslice = append(dslice, a(v))
	}
	return
}

// 随机生成切片
func Slice_rand(a []interface{}) (b interface{}) {
	randnum := rand.Intn(len(a))
	b = a[randnum]
	return
}

func Slice_sum(intslice []int64) (sum int64) {
	for _, v := range intslice {
		sum += v
	}
	return
}

// 对切片中的数据进行过滤
func Slice_filter(slice []interface{}, a filtertype) (ftslice []interface{}) {
	for _, v := range slice {
		if a(v) {
			ftslice = append(ftslice, v)
		}
	}
	return
}

// 返回slice1 不在 slice2 中的所有值
func Slice_diff(slice1, slice2 []interface{}) (diffslice []interface{}) {
	for _, v := range slice1 {
		if !Slice_in(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}

// 交集
func Slice_intersect(slice1, slice2 []interface{}) (intersectslice []interface{}) {
	for _, v := range slice1 {
		if Slice_in(v, slice2) {
			intersectslice = append(intersectslice, v)
		}
	}
	return
}

// 以size为列进行切换。（比如切片长度为20，size为5，则切成4行5列。
func Slice_chunk(slice []interface{}, size int) (chunkslice [][]interface{}) {
	if size >= len(slice) {
		chunkslice = append(chunkslice, slice)
		return
	}
	end := size
	for i := 0; i <= (len(slice) - size); i += size {
		chunkslice = append(chunkslice, slice[i:end])
		end += size
	}
	return
}

// 生成切片
func Slice_range(start, end, step int64) (intslice []int64) {
	for i := start; i <= end; i += step {
		intslice = append(intslice, i)
	}
	return
}

// 如果切片长度小于size，则在size后面补上 切片长度减去size 个val。
func Slice_pad(slice []interface{}, size int, val interface{}) []interface{} {
	if size <= len(slice) {
		return slice
	}
	for i := 0; i < (size - len(slice)); i++ {
		slice = append(slice, val)
	}
	return slice
}

// 切片去重
func Slice_unique(slice []interface{}) (uniqueslice []interface{}) {
	for _, v := range slice {
		if !Slice_in(v, uniqueslice) {
			uniqueslice = append(uniqueslice, v)
		}
	}
	return
}
