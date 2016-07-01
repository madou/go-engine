package renderer

var cubeIndicies = []uint32{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17,
	18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35,
}

var skyboxVerticies = []float32{
	-1, -1, -1, 1, 1, 1, 0.25, 0.334, 1, 1, 1, 1,
	-1, 1, -1, 1, 1, 1, 0.25, 0.6662, 1, 1, 1, 1,
	-1, 1, 1, 1, 1, 1, 0.5, 0.6662, 1, 1, 1, 1,
	-1, 1, 1, 1, 1, 1, 0.5, 0.6662, 1, 1, 1, 1,
	-1, -1, 1, 1, 1, 1, 0.5, 0.334, 1, 1, 1, 1,
	-1, -1, -1, 1, 1, 1, 0.25, 0.334, 1, 1, 1, 1,
	1, -1, -1, 1, 1, 1, 1, 0.334, 1, 1, 1, 1,
	1, -1, 1, 1, 1, 1, 0.75, 0.334, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 0.75, 0.666, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 0.75, 0.666, 1, 1, 1, 1,
	1, 1, -1, 1, 1, 1, 1, 0.666, 1, 1, 1, 1,
	1, -1, -1, 1, 1, 1, 1, 0.334, 1, 1, 1, 1,
	-1, -1, -1, 1, 1, 1, 0.25, 0.333, 1, 1, 1, 1,
	-1, -1, 1, 1, 1, 1, 0.5, 0.333, 1, 1, 1, 1,
	1, -1, 1, 1, 1, 1, 0.5, 0, 1, 1, 1, 1,
	1, -1, 1, 1, 1, 1, 0.5, 0, 1, 1, 1, 1,
	1, -1, -1, 1, 1, 1, 0.25, 0, 1, 1, 1, 1,
	-1, -1, -1, 1, 1, 1, 0.25, 0.333, 1, 1, 1, 1,
	-1, -1, 1, 1, 1, 1, 0.5, 0.334, 1, 1, 1, 1,
	-1, 1, 1, 1, 1, 1, 0.5, 0.666, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 0.75, 0.666, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 0.75, 0.666, 1, 1, 1, 1,
	1, -1, 1, 1, 1, 1, 0.75, 0.334, 1, 1, 1, 1,
	-1, -1, 1, 1, 1, 1, 0.5, 0.334, 1, 1, 1, 1,
	-1, 1, 1, 1, 1, 1, 0.5, 0.667, 1, 1, 1, 1,
	-1, 1, -1, 1, 1, 1, 0.25, 0.667, 1, 1, 1, 1,
	1, 1, -1, 1, 1, 1, 0.25, 1, 1, 1, 1, 1,
	1, 1, -1, 1, 1, 1, 0.25, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 0.5, 1, 1, 1, 1, 1,
	-1, 1, 1, 1, 1, 1, 0.5, 0.667, 1, 1, 1, 1,
	1, -1, -1, 1, 1, 1, 0, 0.333, 1, 1, 1, 1,
	1, 1, -1, 1, 1, 1, 0, 0.666, 1, 1, 1, 1,
	-1, 1, -1, 1, 1, 1, 0.25, 0.666, 1, 1, 1, 1,
	-1, 1, -1, 1, 1, 1, 0.25, 0.666, 1, 1, 1, 1,
	-1, -1, -1, 1, 1, 1, 0.25, 0.333, 1, 1, 1, 1,
	1, -1, -1, 1, 1, 1, 0, 0.333, 1, 1, 1, 1,
}
