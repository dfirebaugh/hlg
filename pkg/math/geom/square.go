package geom

import "math"

type Square [4]Point

func MakeSquare(points [4]Point) Square {
	return Square(points)
}

func (s *Square) Centroid() Point {
	var sumX, sumY float32 = 0.0, 0.0
	for _, p := range s {
		sumX += p.X
		sumY += p.Y
	}
	return MakePoint(sumX/4, sumY/4)
}

func (s *Square) Translate(vector Vector) {
	for i := range s {
		s[i].X += vector[0]
		s[i].Y += vector[1]
	}
}

func (s *Square) Scale(factor float32) {
	centroid := s.Centroid()
	for i := range s {
		dir := s[i].ToVector().Subtract(centroid.ToVector())
		s[i] = Point{
			X: centroid.X + dir[0]*factor,
			Y: centroid.Y + dir[1]*factor,
		}
	}
}

func (s *Square) Rotate(angle float64) {
	cos := float32(math.Cos(angle))
	sin := float32(math.Sin(angle))
	centroid := s.Centroid()

	for i := range s {
		// Translate to origin
		x := s[i].X - centroid.X
		y := s[i].Y - centroid.Y

		// Rotate around origin
		newX := cos*x - sin*y
		newY := sin*x + cos*y

		// Translate back
		s[i] = Point{
			X: newX + centroid.X,
			Y: newY + centroid.Y,
		}
	}
}

func (s *Square) Area() float32 {
	sideLength := s[1].ToVector().GetDistance(s[0].ToVector())
	return float32(math.Pow(float64(sideLength), 2))
}
