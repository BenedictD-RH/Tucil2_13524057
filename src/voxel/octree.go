package voxel

import (
	"go_project/object"
	"sync"
)

type Octree struct {
	value             bool
	object            *object.Object
	voxel             Voxel
	intersectingFaces []int
	children          [8]*Octree
}

func NewOctreeRoot(O *object.Object) *Octree {
	face_idx_list := []int{}
	for i := range O.FaceArray {
		face_idx_list = append(face_idx_list, i)
	}
	voxel := Voxel{*object.NewVertex(0, 0, 0), MaxEdgeLength}
	return &Octree{true,
		O,
		voxel,
		face_idx_list,
		[8]*Octree{nil, nil, nil, nil, nil, nil, nil, nil}}
}

func NewOctree(voxel Voxel, fl []int, O *object.Object) *Octree {
	return &Octree{true,
		O,
		voxel,
		fl,
		[8]*Octree{nil, nil, nil, nil, nil, nil, nil, nil}}
}

func isLeaf(Oct *Octree) bool {
	return !Oct.value || (Oct.children[0] == nil)
}

func amountOfEmptyNodesAtDepth(Oct *Octree, n, depth int) int {
	if Oct == nil {
		return 0
	} else if n == depth {
		if !Oct.value {
			return 1
		} else {
			return 0
		}
	} else {
		total := 0
		for _, c := range Oct.children {
			total += amountOfEmptyNodesAtDepth(c, n+1, depth)
		}
		return total
	}
}

func amountOfNodeAtDepth(Oct *Octree, n, depth int) int {
	if Oct == nil {
		return 0
	} else if n == depth {
		return 1
	} else {
		total := 0
		for _, c := range Oct.children {
			total += amountOfNodeAtDepth(c, n+1, depth)
		}
		return total
	}
}

func GetAmountOfNodeAtDepth(Oct *Octree, depth int) int {
	return amountOfNodeAtDepth(Oct, 0, depth)
}

func GetAmountOfEmptyNodesAtDepth(Oct *Octree, depth int) int {
	return amountOfEmptyNodesAtDepth(Oct, 0, depth)
}

func GetTreeMaxDepth(Oct *Octree) int {
	if isLeaf(Oct) {
		return 1
	} else {
		maxDepth := 0
		for _, p := range Oct.children {
			p_depth := 1 + GetTreeMaxDepth(p)
			if p_depth > maxDepth {
				maxDepth = p_depth
			}
		}
		return maxDepth
	}
}

func getOctreeVoxelEdgeLength(Oct *Octree) float32 {
	if isLeaf(Oct) {
		return Oct.voxel.edgeLength
	} else {
		minEdgeLength := MaxEdgeLength
		for _, p := range Oct.children {
			p_el := getOctreeVoxelEdgeLength(p)
			if p_el < minEdgeLength {
				minEdgeLength = p_el
			}
		}
		return minEdgeLength
	}
}

func IncreaseDepth(Oct *Octree) {
	if isLeaf(Oct) && Oct.value {
		voxelCorners := GetVoxelVertices(&Oct.voxel)
		for i := range Oct.children {
			subVox := Voxel{*object.MidVertex(&Oct.voxel.center, voxelCorners[i]), Oct.voxel.edgeLength / 2}
			new_i_f := []int{}
			for _, f_idx := range Oct.intersectingFaces {
				if isPolygonClippingVoxel(subVox, Oct.object, f_idx) {
					new_i_f = append(new_i_f, f_idx)
				}
			}
			Oct.children[i] = NewOctree(subVox, new_i_f, Oct.object)
			Oct.children[i].value = (len(new_i_f) > 0)
		}
	} else if Oct.value {
		var wg sync.WaitGroup
		for _, p := range Oct.children {
			wg.Add(1)
			go func() {
				defer wg.Done()
				IncreaseDepth(p)
			}()
		}
		wg.Wait()
	}
}

func OctreeVoxelAmount(Oct *Octree) int {
	if !Oct.value {
		return 0
	} else if isLeaf(Oct) {
		return 1
	} else {
		total := 0
		for _, o := range Oct.children {
			total += OctreeVoxelAmount(o)
		}
		return total
	}
}

func OctreeToVoxelObject(Oct *Octree) *VoxelObject {
	if !Oct.value {
		return &VoxelObject{[]*object.Vertex{}, []*object.Face{}}
	} else if isLeaf(Oct) {
		return VoxelToObject(&Oct.voxel)
	} else {
		O := &VoxelObject{[]*object.Vertex{}, []*object.Face{}}
		for _, p := range Oct.children {
			O_p := OctreeToVoxelObject(p)
			O = CombineVoxelObjects(O, O_p)
		}
		return O
	}
}
