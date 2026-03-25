package object

import (
	"slices"
	//"fmt"
	"sync"
)

var mid_point = Point{0,0}


func euclideanMod(a, b int) int {
    m := a % b
    if m < 0 {
        m += b
    }
    return m
}

func quadrant(p Point) int {
	if p.x >= 0 && p.y >= 0 {
		return 1
	} else if p.x <= 0 && p.y >= 0 {
		return 2
	} else if p.x <= 0 && p.y <= 0 {
		return 3
	}
	return 4
}



func orientation(a,b,c Point) int {
	res := (b.y - a.y) * (c.x - b.x) - (c.y - b.y) * (b.x - a.x)
	if res == 0 {
		return 0
	} else if res > 0 {
		return 1
	}
	return -1
}

func compare(p,q Point) int {
	p_ := Point{p.x - mid_point.x, p.y - mid_point.y}
	q_ := Point{q.x - mid_point.x, q.y - mid_point.y}
	one := quadrant(p_)
	two := quadrant(q_)
	if one != two {
		if  one < two {
			return -1
		}
		return 1
	}
	if p_.y*q_.x < q_.y*p_.x {
		return -1
	}
	return 1
}

func sortCounterClockwise(p1, p2 []Point) ([]Point) {
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
			if compare(p1[i_1], p2[i_2]) <= 0 {
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


func MergeSortPointsCounterClockwise(p []Point) []Point {
	if len(p) <= 1 {
		return p
	} else {
		half_point := int(len(p) / 2)
		f_h := p[0:half_point]
		s_h := p[half_point:]
		sorted_f_h := MergeSortPointsCounterClockwise(f_h)
		sorted_s_h := MergeSortPointsCounterClockwise(s_h)
		return sortCounterClockwise(sorted_f_h, sorted_s_h)
	}
}

func sortX(p1, p2 []Point) ([]Point) {
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
			if p1[i_1].x < p2[i_2].x {
				merged = append(merged, p1[i_1])
				i_1++
			} else if p1[i_1].x == p2[i_2].x {
				if p1[i_1].y < p2[i_2].y {
					merged = append(merged, p1[i_1])
					i_1++
				} else {
					merged = append(merged, p2[i_2])
					i_2++
				}
			} else {
				merged = append(merged, p2[i_2])
				i_2++
			}
		}
	}
	return merged
}

func MergeSortPointsX(p []Point) []Point {
	if len(p) <= 1 {
		return p
	} else {
		half_point := int(len(p) / 2)
		f_h := p[0:half_point]
		s_h := p[half_point:]
		var sorted_f_h, sorted_s_h []Point
		sorted_f_h = MergeSortPointsX(f_h)
		sorted_s_h = MergeSortPointsX(s_h)
		return sortX(sorted_f_h, sorted_s_h)
	}
}



func convex_hull_merge(a, b []Point) []Point {
	n1, n2 := len(a), len(b)
	ia, ib := 0,0
	for i, _ := range a {
		if (a[ia].x < a[i].x) {
			ia = i
		}
	}

	for i, _ := range b {
		if (b[ib].x > b[i].x) {
			ib = i
		}
	}
	var uppera, upperb, lowera, lowerb int
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		inda, indb := ia,ib
		curr_n := 0
		max_i := 0 
		done := false
		for {
			done = true
			curr_n = 0
			max_i = inda
			for {
				if (orientation(b[indb], a[inda], a[(inda + 1) % n1]) >= 0) && curr_n <= n1 {
					inda = (inda + 1) % n1
					if (a[inda].y > a[max_i].y) { max_i = inda }
					curr_n++
				} else { 
					if curr_n > n1 { inda = max_i }
					break 
				}
			}
			curr_n = 0
			max_i = indb
			for {
				if (orientation(a[inda], b[indb], b[euclideanMod((n2 + indb - 1), n2)]) <= 0) && curr_n <= n2 {
					indb = euclideanMod((indb - 1), n2)
					if (b[indb].y > b[max_i].y) { max_i = indb }
					done = false
					curr_n++
				} else { 
					if curr_n > n2 { 
						indb = max_i
						done = true 
					}
					break 
				}
			}
			if done { break }
		}
		uppera, upperb = inda, indb
		wg.Done()
	}()
	

	go func() {
		inda, indb := ia,ib
		curr_n := 0
		min_i := 0 
		done := false
		for {
			done = true
			curr_n = 0
			min_i = indb
			for {
				if (orientation(a[inda], b[indb], b[(indb + 1) % n2]) >= 0) && curr_n <= n2 {
					indb = (indb + 1) % n2
					if (b[indb].y < b[min_i].y) { min_i = indb }
					curr_n++
				} else { 
					if curr_n > n2 { indb = min_i }
					break 
				}
			}
			curr_n = 0
			min_i = inda
			for {
				if (orientation(b[indb], a[inda], a[euclideanMod(n1+inda-1, n1)]) <= 0) && curr_n <= n1 {
					inda = euclideanMod((inda - 1), n1)
					if (a[inda].y < a[min_i].y) { min_i = inda }
					done = false
					curr_n++
				} else { 
					if curr_n > n1 { 
						inda = min_i
						done = true 
					}
					break
				}
			}
			if done { break }
		}
		lowera, lowerb = inda, indb
		wg.Done()
	}()
	wg.Wait()
	ret := []Point{}
	
	ind := uppera
	ret = append(ret, a[uppera])
	for {
		if (ind != lowera) {
			ind = (ind + 1) % n1
			ret = append(ret, a[ind])
		} else { break }
	}
	ind = lowerb
	ret = append(ret, b[lowerb])
	for {
		if (ind != upperb) {
			ind = (ind + 1) % n2
			ret = append(ret, b[ind])
		} else { break }
	}
	return ret
}

func addPointToSet(p Point, s []Point) ([] Point) {
	if slices.Contains(s, p) { return s }
	return append(s, p)
}

func bruteHull(a []Point) []Point {
	s := []Point{}
	for i, _ := range a {
		for j := i + 1; j < len(a); j++ {
			x1, x2 := a[i].x, a[j].x
            y1, y2 := a[i].y, a[j].y
            a1, b1, c1 := y1 - y2, x2 -x1, x1*y2-y1*x2
            pos, neg := 0, 0
			for k, _ := range a {
				if (k == i) || (k == j) || (a1*a[k].x+b1*a[k].y+c1 <= 0) {
					neg += 1
				}
				if (k == i) || (k == j) || (a1*a[k].x+b1*a[k].y+c1 >= 0) {
                    pos += 1
				}
			}
			if (pos == len(a)) || (neg == len(a)) {
				s = addPointToSet(a[i], s)
				s = addPointToSet(a[j], s)
			}
		}
	}
	ret := s
	mid_point = Point{0,0}
	n := float32(len(ret))
	for i,_ := range ret {
		mid_point.x += ret[i].x
		mid_point.y += ret[i].y
		ret[i].x *= n
		ret[i].y *= n
	}
	ret = MergeSortPointsCounterClockwise(ret)
	for i,_ := range ret {
		ret[i].x = ret[i].x/n
		ret[i].y = ret[i].y/n
	}
	return ret
}

func convex_hull_divide(a []Point) []Point {
	//fmt.Println("Start : ", len(a))
	if len(a) <= 5 {
        return bruteHull(a)
	}
	left := a[:len(a)/2]
	right := a[len(a)/2:]
	left_hull, right_hull := []Point{}, []Point{}
	left_hull = convex_hull_divide(left)
	right_hull = convex_hull_divide(right)
	// var wg sync.WaitGroup
	// wg.Add(2)
	// go func() {
	// 	defer wg.Done()
	// 	left_hull = convex_hull_divide(left)
	// }()
	// go func() {
	// 	defer wg.Done()
	//  	right_hull = convex_hull_divide(right)
	// }()
	// wg.Wait()
	res := convex_hull_merge(left_hull, right_hull)
	//fmt.Println("Result : ", len(res))
	return res
}


