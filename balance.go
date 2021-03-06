package main

import (
	// "encoding/json"

	"fmt"
	"hash/fnv"
	"math"
	"sort"
	"strconv"

	"github.com/Freeaqingme/GoConsistentHash"
	"github.com/spaolacci/murmur3"
	"github.com/tenfyzhong/cityhash"
)

var (
	nodes    = 10
	requests = 500000
)

func stats(m map[string]int, src string, replica int) {
	nums := make([]int, 0, len(m))
	for _, v := range m {
		nums = append(nums, v)
	}

	sort.Ints(nums)

	var (
		sum float64
		avg float64
	)

	for i := 0; i < len(nums); i++ {
		avg += float64(nums[i])
	}

	avg = avg / float64(len(nums))

	for i := 0; i < len(nums); i++ {
		sum = sum + math.Pow(float64(nums[i])-avg, 2)
	}

	varince := math.Sqrt(sum)

	fmt.Printf("%s\treplica:%d\tvar:%.0f\tmax:%d\tmin:%d\tdiff:%d\tratio:%f\n", src, replica, varince, nums[nodes-1], nums[0], nums[nodes-1]-nums[0], float64(nums[nodes-1]-nums[0])*100/avg)
}

func testMurMurConsistHash(replica int) {
	ch := GoConsistentHash.New(replica, murmur3.Sum32)
	for i := 0; i < nodes; i++ {
		err := ch.AddString(fmt.Sprintf("nodes-%2d", i))
		if err != nil {
			panic(err)
		}
	}

	nums := make(map[string]int, nodes)
	for i := 0; i < requests; i++ {
		node := ch.Get(strconv.Itoa(i))
		if _, exists := nums[node]; exists {
			nums[node]++
		} else {
			nums[node] = 1
		}
	}

	stats(nums, "MurMur", replica)
}

func testCrc32ConsistHash(replica int) {
	ch := GoConsistentHash.New(replica, nil)
	for i := 0; i < nodes; i++ {
		err := ch.AddString(fmt.Sprintf("nodes-%d", i))
		if err != nil {
			panic(err)
		}
	}

	nums := make(map[string]int, nodes)
	for i := 0; i < requests; i++ {
		node := ch.Get(strconv.Itoa(i))
		if _, exists := nums[node]; exists {
			nums[node]++
		} else {
			nums[node] = 1
		}
	}

	stats(nums, "Crc32", replica)
}

func testFnv1ConsistHash(replica int) {
	ch := GoConsistentHash.New(replica, fnv1)
	for i := 0; i < nodes; i++ {
		err := ch.AddString(fmt.Sprintf("nodes-%d", i))
		if err != nil {
			panic(err)
		}
	}

	nums := make(map[string]int, nodes)
	for i := 0; i < requests; i++ {
		node := ch.Get(strconv.Itoa(i))
		if _, exists := nums[node]; exists {
			nums[node]++
		} else {
			nums[node] = 1
		}
	}

	stats(nums, "Fnv1", replica)
}

func fnv1(data []byte) uint32 {
	h := fnv.New32a()
	h.Write(data)
	return h.Sum32()
}

func testCityConsistHash(replica int) {
	ch := GoConsistentHash.New(replica, cityhash.CityHash32)
	for i := 0; i < nodes; i++ {
		err := ch.AddString(fmt.Sprintf("nodes-%d", i))
		if err != nil {
			panic(err)
		}
	}

	nums := make(map[string]int, nodes)
	for i := 0; i < requests; i++ {
		node := ch.Get(strconv.Itoa(i))
		if _, exists := nums[node]; exists {
			nums[node]++
		} else {
			nums[node] = 1
		}
	}

	stats(nums, "Cityhash", replica)
}

func main() {
	replicas := []int{3, 10, 50, 100, 200, 400, 600, 800, 1000}
	for _, replica := range replicas {
		testMurMurConsistHash(replica)
	}

	for _, replica := range replicas {
		testCrc32ConsistHash(replica)
	}

	for _, replica := range replicas {
		testFnv1ConsistHash(replica)
	}

	for _, replica := range replicas {
		testCityConsistHash(replica)
	}
}
