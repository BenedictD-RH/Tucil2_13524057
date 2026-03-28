package voxel

import (
	"go_project/object"
)

const centerPoint_x = 0
const centerPoint_y = 0
const MaxEdgeLength float32 = 10

type VoxelObject struct {
	VertexArray []*object.Vertex
	FaceArray   []*object.Face
}

type Voxel struct {
	center     object.Vertex
	edgeLength float32
}

func (VO *VoxelObject) Translate() *object.Object {
	O := object.NewObject()
	O.VertexArray = VO.VertexArray
	O.FaceArray = VO.FaceArray
	O.Dimension = object.FindObjectDimension(O.VertexArray)
	return O
}

func GetVoxelVertices(vox *Voxel) []*object.Vertex {
	verts := [8]*object.Vertex{}
	c := vox.center
	ex := vox.edgeLength / 2
	verts[0] = object.NewVertex(c.GetX()+ex, c.GetY()+ex, c.GetZ()+ex)
	verts[1] = object.NewVertex(c.GetX()+ex, c.GetY()+ex, c.GetZ()-ex)
	verts[2] = object.NewVertex(c.GetX()+ex, c.GetY()-ex, c.GetZ()+ex)
	verts[3] = object.NewVertex(c.GetX()+ex, c.GetY()-ex, c.GetZ()-ex)
	verts[4] = object.NewVertex(c.GetX()-ex, c.GetY()+ex, c.GetZ()+ex)
	verts[5] = object.NewVertex(c.GetX()-ex, c.GetY()+ex, c.GetZ()-ex)
	verts[6] = object.NewVertex(c.GetX()-ex, c.GetY()-ex, c.GetZ()+ex)
	verts[7] = object.NewVertex(c.GetX()-ex, c.GetY()-ex, c.GetZ()-ex)
	return verts[:]
}

func GetVoxelEdges(vox *Voxel) [][2]*object.Vertex {
	edge_list := [12][2]*object.Vertex{}
	corners := GetVoxelVertices(vox)
	edge_list[0] = [2]*object.Vertex{corners[0], corners[1]}
	edge_list[1] = [2]*object.Vertex{corners[0], corners[2]}
	edge_list[2] = [2]*object.Vertex{corners[0], corners[4]}
	edge_list[3] = [2]*object.Vertex{corners[1], corners[3]}
	edge_list[4] = [2]*object.Vertex{corners[1], corners[5]}
	edge_list[5] = [2]*object.Vertex{corners[7], corners[6]}
	edge_list[6] = [2]*object.Vertex{corners[7], corners[5]}
	edge_list[7] = [2]*object.Vertex{corners[7], corners[3]}
	edge_list[8] = [2]*object.Vertex{corners[4], corners[6]}
	edge_list[9] = [2]*object.Vertex{corners[2], corners[6]}
	edge_list[10] = [2]*object.Vertex{corners[2], corners[3]}
	edge_list[11] = [2]*object.Vertex{corners[4], corners[5]}
	return edge_list[:]
}

func VoxelToObject(vox *Voxel) *VoxelObject {
	O := &VoxelObject{[]*object.Vertex{}, []*object.Face{}}
	O.VertexArray = GetVoxelVertices(vox)
	O.FaceArray = append(O.FaceArray, object.NewFace(1, 2, 3))
	O.FaceArray = append(O.FaceArray, object.NewFace(2, 3, 4))
	O.FaceArray = append(O.FaceArray, object.NewFace(5, 6, 7))
	O.FaceArray = append(O.FaceArray, object.NewFace(6, 7, 8))
	O.FaceArray = append(O.FaceArray, object.NewFace(1, 5, 7))
	O.FaceArray = append(O.FaceArray, object.NewFace(1, 3, 7))
	O.FaceArray = append(O.FaceArray, object.NewFace(2, 6, 8))
	O.FaceArray = append(O.FaceArray, object.NewFace(2, 4, 8))
	O.FaceArray = append(O.FaceArray, object.NewFace(1, 5, 6))
	O.FaceArray = append(O.FaceArray, object.NewFace(1, 2, 6))
	O.FaceArray = append(O.FaceArray, object.NewFace(3, 7, 8))
	O.FaceArray = append(O.FaceArray, object.NewFace(3, 4, 8))
	return O
}

func CombineVoxelObjects(O1, O2 *VoxelObject) *VoxelObject {
	O_combined := &VoxelObject{O1.VertexArray, O1.FaceArray}
	for _, vert := range O2.VertexArray {
		O_combined.VertexArray = append(O_combined.VertexArray, vert)
	}
	for _, face := range O2.FaceArray {
		O_combined.FaceArray = append(O_combined.FaceArray, face.Increment(len(O1.VertexArray)))
	}
	return O_combined
}

func EraseUnusedVertices(O *VoxelObject) *VoxelObject {
	new_vertices := []*object.Vertex{}
	vertex_map := make(map[object.Vertex]int)
	for _, f := range O.FaceArray {
		v_list := f.ListOfIdx()
		for _, v_idx := range v_list {
			_, exists := vertex_map[*O.VertexArray[v_idx-1]]
			if !exists {
				new_vertices = append(new_vertices, O.VertexArray[v_idx-1])
				vertex_map[*O.VertexArray[v_idx-1]] = len(new_vertices)
			}
		}
	}
	reordered_faces := []*object.Face{}
	for _, face := range O.FaceArray {
		v_list := face.ListOfIdx()
		new_v_list := [3]int{}
		for i, v_idx := range v_list {
			new_v_list[i] = vertex_map[*O.VertexArray[v_idx-1]]
		}
		new_face := object.NewFace(new_v_list[0], new_v_list[1], new_v_list[2])
		reordered_faces = append(reordered_faces, new_face)
	}

	return &VoxelObject{new_vertices, reordered_faces}
}

func EraseInnerFaces(O *VoxelObject) *VoxelObject {
	unique_vertices := []*object.Vertex{}
	vertex_map := make(map[object.Vertex]int)
	for _, vert := range O.VertexArray {
		_, exists := vertex_map[*vert]
		if !exists {
			unique_vertices = append(unique_vertices, vert)
			vertex_map[*vert] = len(unique_vertices)
		}
	}
	unique_faces := []*object.Face{}
	face_map := make(map[object.Face]bool)
	for _, f := range O.FaceArray {
		v_list := f.ListOfIdx()
		new_v_list := [3]int{}
		for i, v_idx := range v_list {
			new_v_list[i] = vertex_map[*O.VertexArray[v_idx-1]]
		}
		new_face := object.NewFace(new_v_list[0], new_v_list[1], new_v_list[2])
		val, exists := face_map[*new_face]
		if !exists {
			face_map[*new_face] = true
		} else if val {
			face_map[*new_face] = false
		}
	}
	for key, val := range face_map {
		if val {
			unique_faces = append(unique_faces, &key)
		}
	}
	return EraseUnusedVertices(&VoxelObject{unique_vertices, unique_faces})
}

func EraseDuplicates(O *VoxelObject) *VoxelObject {
	unique_vertices := []*object.Vertex{}
	vertex_map := make(map[object.Vertex]int)
	for _, vert := range O.VertexArray {
		_, exists := vertex_map[*vert]
		if !exists {
			unique_vertices = append(unique_vertices, vert)
			vertex_map[*vert] = len(unique_vertices)
		}
	}
	unique_faces := []*object.Face{}
	face_map := make(map[object.Face]bool)
	for _, f := range O.FaceArray {
		v_list := f.ListOfIdx()
		new_v_list := [3]int{}
		for i, v_idx := range v_list {
			new_v_list[i] = vertex_map[*O.VertexArray[v_idx-1]]
		}
		new_face := object.NewFace(new_v_list[0], new_v_list[1], new_v_list[2])
		_, exists := face_map[*new_face]
		if !exists {
			face_map[*new_face] = true
			unique_faces = append(unique_faces, new_face)
		}
	}
	return &VoxelObject{unique_vertices, unique_faces}
}

func inBetween(a, b, x float32) bool {
	return (a <= x && x <= b) || (b <= x && x <= a)
}

func isVertexInVoxel(vox Voxel, v *object.Vertex) bool {
	voxCorner := GetVoxelVertices(&vox)
	c_1, c_2 := voxCorner[0], voxCorner[7]
	inside := inBetween(c_1.GetX(), c_2.GetX(), v.GetX())
	inside = inside && inBetween(c_1.GetY(), c_2.GetY(), v.GetY())
	inside = inside && inBetween(c_1.GetZ(), c_2.GetZ(), v.GetZ())
	return inside
}

func projectPolygon(v_list []*object.Vertex, axis string) []object.Point {
	projections := []object.Point{}
	projection_map := make(map[object.Point]bool)
	for _, v := range v_list {
		var pr object.Point
		if axis == "x" {
			pr = *object.NewPoint(v.GetY(), v.GetZ())
		} else if axis == "y" {
			pr = *object.NewPoint(v.GetX(), v.GetZ())
		} else if axis == "z" {
			pr = *object.NewPoint(v.GetX(), v.GetY())
		}
		_, exists := projection_map[pr]
		if !exists {
			projections = append(projections, pr)
			projection_map[pr] = true
		}
	}
	return projections
}

func projectionCheck(vox Voxel, O *object.Object, f_idx int) bool {
	p := O.FaceArray[f_idx]
	axis_check := 0
	for _, a := range []string{"x", "y", "z"} {
		face_projection := projectPolygon(object.GetFaceVertices(O, p), a)
		vox_center_pr := projectPolygon([]*object.Vertex{&vox.center}, a)
		vox_projection := object.MergeSortPointsClockwise(projectPolygon(GetVoxelVertices(&vox), a), vox_center_pr[0])
		if object.IsPolygonInPolygon(vox_projection, face_projection) || object.IsPolygonInPolygon(face_projection, vox_projection) {
			axis_check++
		} else {
			intersect := false
			for i := range face_projection {
				p_point_1 := face_projection[i]
				p_point_2 := face_projection[(i+1)%len(face_projection)]
				for j := range vox_projection {
					v_point_1 := vox_projection[j]
					v_point_2 := vox_projection[(j+1)%len(vox_projection)]
					o1 := object.Orientation(p_point_1, p_point_2, v_point_1)
					o2 := object.Orientation(p_point_1, p_point_2, v_point_2)
					o3 := object.Orientation(v_point_1, v_point_2, p_point_1)
					o4 := object.Orientation(v_point_1, v_point_2, p_point_2)
					if o1 != o2 && o3 != o4 {
						intersect = true
					}
					if o1 == 0 && object.OnSegment(p_point_1, p_point_2, v_point_1) {
						intersect = true
					}
					if o2 == 0 && object.OnSegment(p_point_1, p_point_2, v_point_2) {
						intersect = true
					}
					if o3 == 0 && object.OnSegment(v_point_1, v_point_2, p_point_1) {
						intersect = true
					}
					if o4 == 0 && object.OnSegment(v_point_1, v_point_2, p_point_2) {
						intersect = true
					}
					if intersect {
						break
					}
				}
				if intersect {
					break
				}
			}
			if intersect {
				axis_check++
			}
		}
	}
	if axis_check >= 3 {
		return true
	}
	return false
}

func isAPolygonEdgeInVoxel(vox Voxel, O *object.Object, f_idx int) bool {
	p := O.FaceArray[f_idx]
	vert := object.GetFaceVertices(O, p)
	voxCorner := GetVoxelVertices(&vox)
	c_1, c_2 := voxCorner[0], voxCorner[7]
	big_n := 0
	for i := range vert {
		n := 0
		if inBetween(vert[i].GetX(), vert[(i+1)%len(vert)].GetX(), c_1.GetX()) ||
			inBetween(vert[i].GetX(), vert[(i+1)%len(vert)].GetX(), c_2.GetX()) {
			n++
		}
		if inBetween(vert[i].GetY(), vert[(i+1)%len(vert)].GetY(), c_1.GetY()) ||
			inBetween(vert[i].GetY(), vert[(i+1)%len(vert)].GetY(), c_2.GetY()) {
			n++
		}
		if inBetween(vert[i].GetZ(), vert[(i+1)%len(vert)].GetZ(), c_1.GetZ()) ||
			inBetween(vert[i].GetZ(), vert[(i+1)%len(vert)].GetZ(), c_2.GetZ()) {
			n++
		}
		if n >= 3 {
			big_n++
		}
	}
	if big_n >= 3 {
		return true
	}
	return false
}

func isPolygonClippingVoxel(vox Voxel, O *object.Object, f_idx int) bool {
	p := O.FaceArray[f_idx]
	vert := object.GetFaceVertices(O, p)
	for _, v := range vert {
		if isVertexInVoxel(vox, v) {
			return true
		}
	}
	if projectionCheck(vox, O, f_idx) {
		return true
	}
	return false
}

