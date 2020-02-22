package main

import (
	"fmt"
	"math"
)

type Beacon struct {
	name string
	x, y float64
	p0   float64
}

type Obstacle struct {
	id, typeOfObstacle int
	x, y               float64
}

type Object struct {
	lines []int
	x, y  float64
	dist  []float64
}

func DistToObj(ob Object, b Beacon) float64 {
	return math.Sqrt(math.Pow(ob.x-b.x, 2.) + math.Pow(ob.y-b.y, 2.))
}

func FillDistToObj(objs []Object, beacons []Beacon) {
	for i, obj := range objs {
		for _, b := range beacons {
			obj.dist = append(obj.dist, DistToObj(obj, b))
		}
		objs[i] = obj
	}
}

func CalcCalibrationPowers(beacons []Beacon, objs []Object, RSSI [][]int) {
	objInd := map[int]int{} // line index -> Object index
	for i, obj := range objs {
		for _, j := range obj.lines {
			objInd[j] = i
		}
	}
	for i, b := range beacons {
		var validbeaconsignal float64
		for j, vals := range RSSI {
			rssi := vals[i]
			if rssi == 100 {
				continue
			}
			obj := objs[objInd[j]]
			b.p0 += float64(rssi) + 10*2*math.Log10(obj.dist[i])
			validbeaconsignal++
		}
		b.p0 /= validbeaconsignal
		beacons[i] = b
	}
}

func (b *Beacon) Dist(rssi float64) float64 {
	return math.Pow(10, (rssi-b.p0)/-10/2)
}

func (obj *Object) Error(xc, yc float64) float64 {
	err := math.Sqrt((xc-obj.x)*(xc-obj.x) + (yc-obj.y)*(yc-obj.y))
	fmt.Printf("(%7.3f %7.3f) (%7.3f %7.3f) %7.3f\n",
		xc, yc, obj.x, obj.y, err)
	return err
}

func (obj *Object) Coordinate(rss [][]int, beacons []Beacon,
) (float64, float64, float64) {
	var d []float64
	var x []float64
	var y []float64

	for _, numberline := range obj.lines {
		for j, rssi := range rss[numberline] {
			if rssi == 100 {
				continue

			}
			b := beacons[j]
			d = append(d, b.Dist(float64(rssi)))
			x = append(x, b.x)
			y = append(y, b.y)
		}
	}
	//fmt.Println(len(d))
	var xc, yc float64
	for i := 0; i < 10000; i++ {
		var dfdx, dfdy float64
		for i, xi := range x {
			yi, di := y[i], d[i]
			k := 1 / math.Sqrt((xc-xi)*(xc-xi)+(yc-yi)*(yc-yi))
			dfdx += 2 * (xc - xi) * (1 - di*k)
			dfdy += 2 * (yc - yi) * (1 - di*k)
		}
		dd := math.Sqrt(dfdx*dfdx + dfdy*dfdy) // normalize gradient
		xc -= 0.01 * dfdx / dd
		yc -= 0.01 * dfdy / dd
		//fmt.Printf("(%7.3f %7.3f)\n", xc, yc)
	}
	err := obj.Error(xc, yc)
	return xc, yc, err
}

func main() {
	beacons := BeaconsValue("geoBeacons.csv")
	RSSI := RSSIMeasurements("geo_rss.csv")
	objs := ObjectValue("geo_crd.csv")
	FillDistToObj(objs, beacons)
	CalcCalibrationPowers(beacons, objs, RSSI)
	var err float64
	for _, obj := range objs {
		_, _, e := obj.Coordinate(RSSI, beacons)
		err += e
	}
	fmt.Println(err / float64(len(objs)))
}

//for j, b := range beacons {
//	var (
//		cnt int
//		v   float64
//	)
//	for _, numberline := range obj.lines {
//		rssi := rss[numberline][j]
//		if rssi == 100 {
//			continue
//		}
//
//		v += float64(rssi)
//		cnt++
//	}
//	if cnt == 0 {
//		continue
//	}
//	v /= float64(cnt)
//	d = append(d, b.Dist(v))
//	x = append(x, b.x)
//	y = append(y, b.y)
//}
