package object 

import (
	"math"
	"sync"
)

var c_hull_set map[Point]bool
var rwmu sync.RWMutex

func findSide(p1, p2, p Point) (int) {
	val := (p.y - p1.y)*(p2.x - p1.x) - (p.x - p1.x)*(p2.y - p1.y)
	if val > 0 {
		return 1
	} else if val < 0 {
		return -1
	}
	return 0
}

func lineDist(p1, p2, p Point) (float32) {
	return float32(math.Abs(float64((p.y - p1.y)*(p2.x - p1.x) - (p.x - p1.x)*(p2.y - p1.y))))
}

func joinPointSets(s1, s2 map[Point]bool) (map[Point]bool) {
	j_s := s1
	for key := range s2 {
		_, exists := j_s[key]
		if !exists {
			j_s[key] = true
		}
	}
	return j_s
}

func quickHull(a []Point, n int, p1, p2 Point, side int) {
	idx := -1
	var max_dist float32 = 0

	for i := range n {
		temp := lineDist(p1, p2, a[i])
		if (findSide(p1, p2, a[i]) == side) && (temp > max_dist) {
			idx = i
			max_dist = temp
		}
	}

	if idx == -1 {
		rwmu.Lock()
		defer rwmu.Unlock()
		_, exists := c_hull_set[p1]
		if !exists {
			c_hull_set[p1] = true
		}

		_, exists = c_hull_set[p2]
		if !exists {
			c_hull_set[p2] = true
		}
		return
	}
	rwmu.Lock()
	_, exists := c_hull_set[a[idx]]
	rwmu.Unlock()
	if exists { return }
	
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		quickHull(a, n, a[idx], p1, -findSide(a[idx], p1, p2))
		wg.Done()
	}()

	go func() {
		quickHull(a, n, a[idx], p2, -findSide(a[idx], p2, p1))
		wg.Done()
	}()
	wg.Wait()
}

func quickHullStart(a []Point) ([]Point){
	c_hull_set = make(map[Point]bool)
	min_x := 0
    max_x := len(a) - 1
    
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		quickHull(a, len(a), a[min_x], a[max_x], 1)
		wg.Done()
	}()

	go func() {
		quickHull(a, len(a), a[min_x], a[max_x], -1)
		wg.Done()
	}()
		
	
	wg.Wait()
	
	
	c_hull := []Point{}
	for key := range c_hull_set {
		c_hull = append(c_hull, key)
	}
	return MergeSortPointsClockwise(c_hull, GetMidPoint(a[min_x], a[max_x]))
}

func sortClockwise(p1, p2 []Point, center Point) ([]Point) {
	i_1 := 0
	i_2 := 0
	merged := []Point{}
	for {
		if i_1 == len(p1) && i_2 == len(p2) {
			break
		} else if i_1 == len(p1) {
			merged = append(merged, p2[i_2])
			i_2++
		} else if i_2 == len(p2) {
			merged = append(merged, p1[i_1])
			i_1++
		} else {
			if ComparePoints(center, p1[i_1], p2[i_2]) {
				merged = append(merged, p1[i_1])
				i_1++
			} else {
				merged = append(merged, p2[i_2])
				i_2++
			}
		}
	}
	return merged
}


func MergeSortPointsClockwise(p []Point, center Point) ([]Point) {
	if len(p) <= 1 {
		return p
	} else {
		half_point := int(len(p) / 2)
		f_h := p[0:half_point]
		s_h := p[half_point:]
		sorted_f_h, sorted_s_h := []Point{}, []Point{}
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			sorted_f_h = MergeSortPointsClockwise(f_h, center)
			wg.Done()
		}()
		go func() {
		 	sorted_s_h = MergeSortPointsClockwise(s_h, center)
			wg.Done()
		}()
		wg.Wait()
		return sortClockwise(sorted_f_h, sorted_s_h, center)
	}
}


func ComparePoints(center, p1, p2 Point) (bool) {
	angle_1, angle_2 := GetAngle(p1, center), GetAngle(p2, center)
	if (angle_1 < angle_2) {
		return true
	} else if (angle_1 == angle_2) {
		dist1, dist2 := GetPointDistance(center, p1), GetPointDistance(center, p2)
		if dist1 < dist2 {
			return true
		}
		return false
	}
	return false
}

func GetAngle(p, center Point) (float32) {
	x := float64(p.x - center.x)
	y := float64(p.y - center.y)
	angle := float32(math.Atan2(y,x))
	if angle <= 0 {
		angle += 2*math.Pi
	}
	return angle
}

func GetPointDistance(p1, p2 Point) (float32) {
	return float32(math.Sqrt(float64((p1.x - p2.x)*(p1.x - p2.x) + (p1.y - p2.y)*(p1.y - p2.y))))
}

func GetMidPoint(p1, p2 Point) (Point) {
	return Point{(p1.x + p2.x)/2, (p1.y + p2.y)/2}
}